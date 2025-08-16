package meilisearch

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func httpRespWithBody(t *testing.T, body string, contentType string) *http.Response {
	t.Helper()
	if contentType == "" {
		contentType = "text/event-stream"
	}
	return &http.Response{
		StatusCode: 200,
		Header:     http.Header{"Content-Type": []string{contentType}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestDecoder_SingleEventsAndDone(t *testing.T) {
	data := "" +
		"event: message\n" +
		"data: {\"key\": \"value\"}\n" +
		"\n" +
		"event: done\n" +
		"data: [DONE]\n" +
		"\n"

	res := httpRespWithBody(t, data, "text/event-stream")
	decoder := NewDecoder(res)
	require.NotNil(t, decoder)

	require.True(t, decoder.Next(), "expected Next() to return true for first event")
	event := decoder.Event()
	assert.Equal(t, "message", event.Type, "unexpected event type")
	// New decoder DOES NOT append a trailing newline after the last data line
	assert.Equal(t, []byte("{\"key\": \"value\"}"), event.Data, "unexpected event data")

	require.True(t, decoder.Next(), "expected Next() to return true for second event")
	event = decoder.Event()
	assert.Equal(t, "done", event.Type, "unexpected event type")
	assert.Equal(t, []byte("[DONE]"), event.Data, "unexpected event data")

	assert.False(t, decoder.Next(), "expected Next() to return false (EOF)")
	assert.NoError(t, decoder.Err(), "unexpected error")
	assert.NoError(t, decoder.Close(), "unexpected error on Close")
}

func TestStream_BasicFlow(t *testing.T) {
	type Message struct {
		Key string `json:"key"`
	}

	data := "" +
		"event: message\n" +
		"data: {\"key\": \"value\"}\n" +
		"\n" +
		"event: done\n" +
		"data: [DONE]\n" +
		"\n"

	res := httpRespWithBody(t, data, "text/event-stream")
	decoder := NewDecoder(res)
	stream := NewStream[Message](decoder, nil)

	require.True(t, stream.Next(), "expected Next() to return true")
	current := stream.Current()
	assert.Equal(t, "value", current.Key, "unexpected current.Key")

	// [DONE] should be consumed internally and stop the stream
	assert.False(t, stream.Next(), "expected Next() to return false after [DONE]")
	assert.NoError(t, stream.Err(), "unexpected error")
	assert.NoError(t, stream.Close(), "unexpected error on Close")
}

func TestStream_ErrorOnInvalidJSON(t *testing.T) {
	data := "" +
		"event: message\n" +
		"data: invalid-json\n" +
		"\n"

	res := httpRespWithBody(t, data, "text/event-stream")
	decoder := NewDecoder(res)
	stream := NewStream[map[string]string](decoder, nil)

	assert.False(t, stream.Next(), "expected Next() to return false on invalid JSON")
	require.Error(t, stream.Err(), "expected an error")
	assert.NoError(t, stream.Close(), "unexpected error on Close")
}

func TestDecoder_ReadErrorPropagates(t *testing.T) {
	res := &http.Response{
		StatusCode: 200,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:       io.NopCloser(&errorReader{}),
	}

	decoder := NewDecoder(res)
	assert.False(t, decoder.Next(), "expected Next() to return false on read error")
	require.Error(t, decoder.Err(), "expected an error")
	assert.NoError(t, decoder.Close(), "unexpected error on Close")
}

func TestStream_MultipleMessagesThenDone(t *testing.T) {
	type Message struct {
		Key string `json:"key"`
	}

	data := "" +
		"event: message\n" +
		"data: {\"key\": \"value1\"}\n" +
		"\n" +
		"event: message\n" +
		"data: {\"key\": \"value2\"}\n" +
		"\n" +
		"event: done\n" +
		"data: [DONE]\n" +
		"\n"

	res := httpRespWithBody(t, data, "text/event-stream")
	decoder := NewDecoder(res)
	stream := NewStream[Message](decoder, nil)

	require.True(t, stream.Next(), "expected Next() true for first message")
	current := stream.Current()
	assert.Equal(t, "value1", current.Key)

	require.True(t, stream.Next(), "expected Next() true for second message")
	current = stream.Current()
	assert.Equal(t, "value2", current.Key)

	assert.False(t, stream.Next(), "expected Next() false after [DONE]")
	assert.NoError(t, stream.Err())
	assert.NoError(t, stream.Close())
}

// Ensures we handle media type parameters (e.g., charset)
func TestDecoder_ContentTypeWithParams(t *testing.T) {
	data := "event: message\n" +
		"data: {\"ok\": true}\n\n"

	res := httpRespWithBody(t, data, "text/event-stream; charset=utf-8")
	decoder := NewDecoder(res)
	require.NotNil(t, decoder)

	require.True(t, decoder.Next())
	ev := decoder.Event()
	assert.Equal(t, "message", ev.Type)
	assert.Equal(t, []byte("{\"ok\": true}"), ev.Data)

	assert.False(t, decoder.Next())
	assert.NoError(t, decoder.Err())
	assert.NoError(t, decoder.Close())
}

// If stream ends without a trailing blank line, last event should still dispatch.
func TestDecoder_FlushLastEventWithoutTrailingBlankLine(t *testing.T) {
	data := "" +
		"event: message\n" +
		"data: {\"final\": true}\n" // no terminating blank line

	res := httpRespWithBody(t, data, "text/event-stream")
	decoder := NewDecoder(res)

	require.True(t, decoder.Next())
	ev := decoder.Event()
	assert.Equal(t, "message", ev.Type)
	assert.Equal(t, []byte("{\"final\": true}"), ev.Data)

	assert.False(t, decoder.Next())
	assert.NoError(t, decoder.Err())
	assert.NoError(t, decoder.Close())
}

// Multiple data lines must be joined with a single '\n' between lines (no trailing newline).
func TestDecoder_MultiDataLinesJoinWithLF(t *testing.T) {
	data := "" +
		"event: chunk\n" +
		"data: hello\n" +
		"data: world\n" +
		"\n"

	res := httpRespWithBody(t, data, "text/event-stream")
	decoder := NewDecoder(res)

	require.True(t, decoder.Next())
	ev := decoder.Event()
	assert.Equal(t, "chunk", ev.Type)
	assert.Equal(t, []byte("hello\nworld"), ev.Data)

	assert.False(t, decoder.Next())
	assert.NoError(t, decoder.Err())
	assert.NoError(t, decoder.Close())
}

// Comments and unknown fields are ignored per spec.
func TestDecoder_CommentsAndUnknownFieldsIgnored(t *testing.T) {
	data := "" +
		": this is a comment\n" +
		"retry: 10000\n" + // ignored
		"id: 42\n" + // ignored by our parser
		"event: ping\n" +
		"data: {}\n" +
		"\n"

	res := httpRespWithBody(t, data, "text/event-stream")
	decoder := NewDecoder(res)

	require.True(t, decoder.Next())
	ev := decoder.Event()
	assert.Equal(t, "ping", ev.Type)
	assert.Equal(t, []byte("{}"), ev.Data)

	assert.False(t, decoder.Next())
	assert.NoError(t, decoder.Err())
	assert.NoError(t, decoder.Close())
}

func TestRegisterDecoder_InvalidContentType(t *testing.T) {
	// Force the fallback branch in RegisterDecoder when mime.ParseMediaType fails
	key := "%%%INVALID%%%"
	called := false
	RegisterDecoder(key, func(rc io.ReadCloser) Decoder {
		called = true
		return newSSEDecoder(rc)
	})
	// The decoderTypes map should now contain the raw lower-cased key
	_, ok := decoderTypes[strings.ToLower(key)]
	assert.True(t, ok, "expected decoderTypes to contain fallback key")
	// Just sanity: constructing a response with this invalid content-type won't match
	// (NewDecoder parsing fails and defaults to SSE) so we only validate registration here.
	assert.False(t, called, "decoder factory should not be invoked in this test")
}

func TestNewDecoder_NilResponse(t *testing.T) {
	dec := NewDecoder(nil)
	assert.False(t, dec.Next())
	assert.Error(t, dec.Err())
	assert.NoError(t, dec.Close())
	// Validate errDecoder accessors explicitly
	assert.Equal(t, Event{}, dec.Event())
}

func TestNewDecoder_DefaultsToSSEWhenUnregistered(t *testing.T) {
	data := "event: msg\n" +
		"data: {\"x\":1}\n\n"
	// Use an unregistered (but valid) content type so NewDecoder falls back to SSE.
	res := httpRespWithBody(t, data, "application/json")
	dec := NewDecoder(res)
	require.True(t, dec.Next())
	assert.Equal(t, "msg", dec.Event().Type)
	assert.Equal(t, []byte("{\"x\":1}"), dec.Event().Data)
	assert.False(t, dec.Next())
	assert.NoError(t, dec.Err())
	assert.NoError(t, dec.Close())
}

func TestSSEDecoder_NextAfterCloseAndDoubleClose(t *testing.T) {
	data := "event: a\n" +
		"data: {\"k\":\"v\"}\n\n"
	res := httpRespWithBody(t, data, "text/event-stream")
	dec := NewDecoder(res)
	require.True(t, dec.Next())
	assert.NoError(t, dec.Close()) // first close
	assert.NoError(t, dec.Close()) // second close hits early-return branch
	assert.False(t, dec.Next(), "Next after close must be false")
}

func TestSSEDecoder_EventFieldWithoutColon(t *testing.T) {
	// Line with just 'event' (no colon) should create an event with empty type & data
	data := "event\n\n"
	res := httpRespWithBody(t, data, "text/event-stream")
	dec := NewDecoder(res)
	require.True(t, dec.Next())
	ev := dec.Event()
	assert.Equal(t, "", ev.Type)
	assert.Equal(t, 0, len(ev.Data))
	assert.False(t, dec.Next())
	assert.NoError(t, dec.Err())
}

func TestStream_NewStreamNilDecoder(t *testing.T) {
	stream := NewStream[struct{}](nil, nil)
	assert.Error(t, stream.Err())
	assert.False(t, stream.Next())    // early return path
	assert.NoError(t, stream.Close()) // decoder nil branch
}

// fakeDecoder lets us simulate an underlying decoder that ends with an error.
type fakeDecoder struct {
	events []Event
	idx    int
	err    error
	cur    Event
	closed bool
}

func (f *fakeDecoder) Next() bool {
	if f.idx < len(f.events) {
		f.cur = f.events[f.idx]
		f.idx++
		return true
	}
	return false
}
func (f *fakeDecoder) Event() Event { return f.cur }
func (f *fakeDecoder) Err() error   { return f.err }
func (f *fakeDecoder) Close() error { f.closed = true; return nil }

func TestStream_DecoderErrorPropagatesAfterIteration(t *testing.T) {
	sentinel := errors.New("decoder failed")
	fd := &fakeDecoder{err: sentinel}
	stream := NewStream[map[string]any](fd, nil)
	assert.False(t, stream.Next()) // no events; should look at decoder.Err()
	assert.ErrorIs(t, stream.Err(), sentinel)
	assert.NoError(t, stream.Close())
}

func TestSSEDecoder_LeadingEmptyLinesIgnored(t *testing.T) {
	data := "\n\nevent: message\n" +
		"data: {\"a\":1}\n\n"
	res := httpRespWithBody(t, data, "text/event-stream")
	dec := NewDecoder(res)
	// First Next should skip the leading blank lines and then parse the event.
	require.True(t, dec.Next())
	ev := dec.Event()
	assert.Equal(t, "message", ev.Type)
	assert.Equal(t, []byte("{\"a\":1}"), ev.Data)
	assert.False(t, dec.Next())
	assert.NoError(t, dec.Err())
}
