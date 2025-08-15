package meilisearch

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
)

var decoderTypes = map[string]func(io.ReadCloser) Decoder{}

func RegisterDecoder(contentType string, decoder func(io.ReadCloser) Decoder) {
	// Normalize key to the media type only (e.g., "text/event-stream")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil || mediaType == "" {
		mediaType = strings.ToLower(strings.TrimSpace(contentType))
	}
	decoderTypes[strings.ToLower(mediaType)] = decoder
}

type Decoder interface {
	Next() bool   // advances to next event; false on EOF or error
	Event() Event // returns the current event
	Err() error   // returns the terminal error (if any)
	Close() error // closes the underlying stream
}

type Event struct {
	Type string
	Data []byte
}

func NewDecoder(res *http.Response) Decoder {
	if res == nil || res.Body == nil {
		return &errDecoder{err: errors.New("nil http response/body")}
	}

	mediaType, _, _ := mime.ParseMediaType(res.Header.Get("Content-Type"))
	mediaType = strings.ToLower(mediaType)

	if f, ok := decoderTypes[mediaType]; ok {
		return f(res.Body)
	}

	// Default: assume SSE if server didn’t set proper header (many do this in dev)
	// We still use our SSE decoder because Meilisearch chat streaming uses SSE.
	return newSSEDecoder(res.Body)
}

// sseDecoder implements parsing per WHATWG SSE spec basics:
// - Fields: "event:" and "data:" are supported (others ignored).
// - An empty line terminates one event.
// - Multiple "data:" lines are joined with "\n" between lines.
type sseDecoder struct {
	rc      io.ReadCloser
	scn     *bufio.Scanner
	cur     Event
	err     error
	closed  bool
	bufSize int
}

func newSSEDecoder(rc io.ReadCloser) *sseDecoder {
	scn := bufio.NewScanner(rc)

	// Increase max token size to safely handle large SSE chunks.
	// Default is 64K; we raise to 8MB to be safe with LLM chunks.
	const maxToken = 8 << 20
	scn.Buffer(make([]byte, 0, 64<<10), maxToken)

	// Use a split that returns full lines WITHOUT the trailing newline,
	// and handles LF/CRLF. bufio.ScanLines already does this well.
	// If you need binary safety, you could write a custom SplitFunc.
	scn.Split(bufio.ScanLines)

	return &sseDecoder{
		rc:      rc,
		scn:     scn,
		bufSize: maxToken,
	}
}

func (d *sseDecoder) Next() bool {
	if d.err != nil || d.closed {
		return false
	}

	var eventType string
	var data bytes.Buffer
	haveAny := false

	for d.scn.Scan() {
		line := d.scn.Text()

		// Empty line => dispatch accumulated event (if any)
		if len(line) == 0 {
			if haveAny {
				d.cur = Event{Type: eventType, Data: data.Bytes()}
				return true
			}
			// If we see multiple empty lines, just keep scanning.
			continue
		}

		// Comments (lines starting with ':') are ignored per spec.
		if line[0] == ':' {
			continue
		}

		name, val, found := strings.Cut(line, ":")
		if !found {
			// Entire line is the field name with empty value; spec allows it.
			// We only care about "event" and "data". Others ignored.
			if name == "event" {
				eventType = ""
				haveAny = true
			}
			continue
		}

		// Trim one optional leading space after colon.
		if len(val) > 0 && val[0] == ' ' {
			val = val[1:]
		}

		switch name {
		case "event":
			eventType = val
			haveAny = true
		case "data":
			if data.Len() > 0 {
				_ = data.WriteByte('\n')
			}
			if _, wErr := data.WriteString(val); wErr != nil {
				d.err = wErr
				return false
			}
			haveAny = true
		default:
			// Ignore other fields like id:, retry:, etc.
		}
	}

	// Flush last event if stream ended without trailing blank line
	if haveAny && d.scn.Err() == nil {
		d.cur = Event{Type: eventType, Data: data.Bytes()}
		return true
	}

	// Propagate scanner error (including io.EOF as nil)
	if scanErr := d.scn.Err(); scanErr != nil {
		d.err = scanErr
	}
	return false
}

func (d *sseDecoder) Event() Event { return d.cur }

func (d *sseDecoder) Err() error { return d.err }

func (d *sseDecoder) Close() error {
	if d.closed {
		return nil
	}
	d.closed = true
	return d.rc.Close()
}

type errDecoder struct {
	err error
}

func (e *errDecoder) Next() bool   { return false }
func (e *errDecoder) Event() Event { return Event{} }
func (e *errDecoder) Err() error   { return e.err }
func (e *errDecoder) Close() error { return nil }

func init() {
	RegisterDecoder("text/event-stream", func(rc io.ReadCloser) Decoder {
		return newSSEDecoder(rc)
	})
}

type Stream[T any] struct {
	decoder   Decoder
	current   T
	err       error
	done      bool
	unmarshal JSONUnmarshal
}

func NewStream[T any](decoder Decoder, unmarshal JSONUnmarshal) *Stream[T] {
	if decoder == nil {
		return &Stream[T]{err: errors.New("nil decoder")}
	}
	if unmarshal == nil {
		unmarshal = json.Unmarshal
	}
	return &Stream[T]{decoder: decoder, unmarshal: unmarshal}
}

func (s *Stream[T]) Next() bool {
	if s.err != nil || s.decoder == nil {
		return false
	}

	for s.decoder.Next() {
		ev := s.decoder.Event()

		// Handle [DONE] (ignore trailing newline if present)
		data := bytes.TrimSpace(ev.Data)
		if bytes.Equal(data, []byte("[DONE]")) {
			s.done = true
			continue
		}

		// Meilisearch returns OpenAI-style chunks directly in data lines.
		// Optionally: detect error envelope if your server sends one.
		var next T
		if err := s.unmarshal(data, &next); err != nil {
			s.err = fmt.Errorf("unmarshal stream chunk: %w", err)
			return false
		}
		s.current = next
		return true
	}

	// Drain terminal error from decoder (nil on clean EOF)
	if derr := s.decoder.Err(); derr != nil {
		s.err = derr
	}
	return false
}

func (s *Stream[T]) Current() T { return s.current }
func (s *Stream[T]) Err() error { return s.err }
func (s *Stream[T]) Close() error {
	if s.decoder == nil {
		return nil
	}
	return s.decoder.Close()
}
