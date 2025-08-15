package meilisearch

import (
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
