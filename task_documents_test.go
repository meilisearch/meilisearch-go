package meilisearch

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type taskDocumentTest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type taskDocumentRoundTripFunc func(*http.Request) (*http.Response, error)

func (f taskDocumentRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type taskDocumentErrorReader struct{}

func (taskDocumentErrorReader) Read([]byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func newTaskDocumentTestClient(fn taskDocumentRoundTripFunc, options ...Option) ServiceManager {
	options = append([]Option{WithCustomClient(&http.Client{Transport: fn})}, options...)
	return New("http://meilisearch.test", options...)
}

func taskDocumentResponse(statusCode int, contentType string, body io.Reader) *http.Response {
	resp := &http.Response{
		StatusCode: statusCode,
		Header:     http.Header{},
		Body:       io.NopCloser(body),
	}
	resp.Header.Set("Content-Type", contentType)
	return resp
}

func encodedTaskDocumentBody(t *testing.T, encoding ContentEncoding, body string) io.Reader {
	t.Helper()

	encoded, err := newEncoding(encoding, DefaultCompression).Encode(strings.NewReader(body))
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = encoded.Close()
	})
	return encoded
}

func TestGetTaskDocumentsDestinationValidation(t *testing.T) {
	client := New("http://127.0.0.1:1")

	t.Run("dst is not a pointer", func(t *testing.T) {
		var docs []taskDocumentTest
		err := client.GetTaskDocuments(1, docs)
		require.ErrorContains(t, err, "dst must be a non-nil pointer to a slice")
	})

	t.Run("dst is a nil pointer", func(t *testing.T) {
		var docs *[]taskDocumentTest
		err := client.GetTaskDocuments(1, docs)
		require.ErrorContains(t, err, "dst must be a non-nil pointer to a slice")
	})

	t.Run("dst is a pointer to a non-slice", func(t *testing.T) {
		var doc taskDocumentTest
		err := client.GetTaskDocuments(1, &doc)
		require.ErrorContains(t, err, "dst must point to a slice")
	})
}

func TestGetTaskDocumentsRequiresNDJSONContentType(t *testing.T) {
	client := newTaskDocumentTestClient(func(r *http.Request) (*http.Response, error) {
		require.Equal(t, "/tasks/1/documents", r.URL.Path)
		return taskDocumentResponse(http.StatusOK, "application/json", strings.NewReader(`{"results":[]}`)), nil
	})

	var docs []taskDocumentTest
	err := client.GetTaskDocuments(1, &docs)
	require.ErrorContains(t, err, `unexpected Content-Type "application/json"`)
	var meiliErr *Error
	require.ErrorAs(t, err, &meiliErr)
	require.Equal(t, ErrCodeResponseUnmarshalBody, meiliErr.ErrCode)
}

func TestGetTaskDocumentsDecodesNDJSON(t *testing.T) {
	client := newTaskDocumentTestClient(func(r *http.Request) (*http.Response, error) {
		require.Equal(t, "/tasks/42/documents", r.URL.Path)
		require.Empty(t, r.URL.RawQuery)
		return taskDocumentResponse(http.StatusOK, "application/x-ndjson; charset=utf-8", strings.NewReader("{\"id\":\"1\",\"name\":\"Alice\"}\n{\"id\":\"2\",\"name\":\"Bob\"}\n")), nil
	})

	var docs []taskDocumentTest
	err := client.GetTaskDocuments(42, &docs)
	require.NoError(t, err)
	require.Equal(t, []taskDocumentTest{
		{ID: "1", Name: "Alice"},
		{ID: "2", Name: "Bob"},
	}, docs)
}

func TestGetTaskDocumentsDecodesConcatenatedNDJSON(t *testing.T) {
	client := newTaskDocumentTestClient(func(r *http.Request) (*http.Response, error) {
		require.Equal(t, "/tasks/42/documents", r.URL.Path)
		require.Empty(t, r.URL.RawQuery)
		return taskDocumentResponse(http.StatusOK, "application/x-ndjson; charset=utf-8", strings.NewReader("{\"id\":\"1\",\"name\":\"Alice\"}{\"id\":\"2\",\"name\":\"Bob\"}{\"id\":\"3\",\"name\":\"Carol\"}")), nil
	})

	var docs []taskDocumentTest
	err := client.GetTaskDocuments(42, &docs)
	require.NoError(t, err)
	require.Equal(t, []taskDocumentTest{
		{ID: "1", Name: "Alice"},
		{ID: "2", Name: "Bob"},
		{ID: "3", Name: "Carol"},
	}, docs)
}

func TestGetTaskDocumentsDecodesEncodedNDJSON(t *testing.T) {
	client := newTaskDocumentTestClient(func(r *http.Request) (*http.Response, error) {
		require.Equal(t, GzipEncoding.String(), r.Header.Get("Accept-Encoding"))

		resp := taskDocumentResponse(http.StatusOK, "application/x-ndjson", encodedTaskDocumentBody(t, GzipEncoding, "{\"id\":\"1\",\"name\":\"Alice\"}\n"))
		resp.Header.Set("Content-Encoding", GzipEncoding.String())
		return resp, nil
	}, WithContentEncoding(GzipEncoding, DefaultCompression))

	var docs []taskDocumentTest
	err := client.GetTaskDocuments(1, &docs)
	require.NoError(t, err)
	require.Equal(t, []taskDocumentTest{{ID: "1", Name: "Alice"}}, docs)
}

func TestGetTaskDocumentsDecodesResponseContentEncodings(t *testing.T) {
	tests := []struct {
		name     string
		encoding ContentEncoding
	}{
		{name: "deflate", encoding: DeflateEncoding},
		{name: "brotli", encoding: BrotliEncoding},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTaskDocumentTestClient(func(r *http.Request) (*http.Response, error) {
				require.Empty(t, r.Header.Get("Accept-Encoding"))

				resp := taskDocumentResponse(http.StatusOK, "application/x-ndjson", encodedTaskDocumentBody(t, tt.encoding, "{\"id\":\"1\",\"name\":\"Alice\"}\n"))
				resp.Header.Set("Content-Encoding", tt.encoding.String())
				return resp, nil
			})

			var docs []taskDocumentTest
			err := client.GetTaskDocuments(1, &docs)
			require.NoError(t, err)
			require.Equal(t, []taskDocumentTest{{ID: "1", Name: "Alice"}}, docs)
		})
	}
}

func TestGetTaskDocumentsDecodesEmptyNDJSON(t *testing.T) {
	client := newTaskDocumentTestClient(func(r *http.Request) (*http.Response, error) {
		require.Equal(t, "/tasks/42/documents", r.URL.Path)
		return taskDocumentResponse(http.StatusOK, "application/x-ndjson", strings.NewReader("")), nil
	})

	docs := []taskDocumentTest{{ID: "stale", Name: "Stale"}}
	err := client.GetTaskDocuments(42, &docs)
	require.NoError(t, err)
	require.NotNil(t, docs)
	require.Empty(t, docs)
}

func TestGetTaskDocumentsResponseDecoderError(t *testing.T) {
	tests := []struct {
		name     string
		encoding ContentEncoding
		body     string
	}{
		{name: "gzip", encoding: GzipEncoding, body: "not gzip"},
		{name: "deflate", encoding: DeflateEncoding, body: "not deflate"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTaskDocumentTestClient(func(r *http.Request) (*http.Response, error) {
				require.Equal(t, "/tasks/42/documents", r.URL.Path)
				resp := taskDocumentResponse(http.StatusOK, "application/x-ndjson", strings.NewReader(tt.body))
				resp.Header.Set("Content-Encoding", tt.encoding.String())
				return resp, nil
			})

			var docs []taskDocumentTest
			err := client.GetTaskDocuments(42, &docs)
			require.ErrorContains(t, err, "failed to create response decoder")
		})
	}
}

func TestGetTaskDocumentsUnsupportedContentEncoding(t *testing.T) {
	client := newTaskDocumentTestClient(func(r *http.Request) (*http.Response, error) {
		require.Equal(t, "/tasks/42/documents", r.URL.Path)
		resp := taskDocumentResponse(http.StatusOK, "application/x-ndjson", strings.NewReader("{}"))
		resp.Header.Set("Content-Encoding", "compress")
		return resp, nil
	})

	var docs []taskDocumentTest
	err := client.GetTaskDocuments(42, &docs)
	require.ErrorContains(t, err, `failed to create response decoder: unsupported Content-Encoding "compress"`)
}

func TestGetTaskDocumentsAPIError(t *testing.T) {
	tests := []struct {
		name     string
		encoding ContentEncoding
		options  []Option
	}{
		{name: "plain"},
		{name: "gzip", encoding: GzipEncoding, options: []Option{WithContentEncoding(GzipEncoding, DefaultCompression)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTaskDocumentTestClient(func(r *http.Request) (*http.Response, error) {
				if tt.encoding.IsZero() {
					require.Empty(t, r.Header.Get("Accept-Encoding"))
				} else {
					require.Equal(t, tt.encoding.String(), r.Header.Get("Accept-Encoding"))
				}

				body := io.Reader(strings.NewReader(`{"message":"bad request","code":"bad_request","type":"invalid_request","link":"https://docs.meilisearch.com/errors#bad_request"}`))
				resp := taskDocumentResponse(http.StatusBadRequest, "application/json", body)
				if !tt.encoding.IsZero() {
					resp.Body = io.NopCloser(encodedTaskDocumentBody(t, tt.encoding, `{"message":"bad request","code":"bad_request","type":"invalid_request","link":"https://docs.meilisearch.com/errors#bad_request"}`))
					resp.Header.Set("Content-Encoding", tt.encoding.String())
				}
				return resp, nil
			}, tt.options...)

			var docs []taskDocumentTest
			err := client.GetTaskDocuments(42, &docs)
			require.Error(t, err)

			var meiliErr *Error
			require.ErrorAs(t, err, &meiliErr)
			require.Equal(t, MeilisearchApiError, meiliErr.ErrCode)
			require.Equal(t, "bad_request", meiliErr.MeilisearchApiError.Code)
		})
	}
}

func TestGetTaskDocumentsStreamingHelpers(t *testing.T) {
	t.Run("nil accepted status codes", func(t *testing.T) {
		c := &client{}
		err := c.handleStreamingStatusCode(&internalRequest{}, taskDocumentResponse(http.StatusOK, "application/x-ndjson", strings.NewReader("")), &Error{})
		require.NoError(t, err)
	})

	t.Run("body read error", func(t *testing.T) {
		c := &client{}
		resp := taskDocumentResponse(http.StatusBadRequest, "application/json", taskDocumentErrorReader{})
		err := c.handleStreamingStatusCode(&internalRequest{acceptedStatusCodes: []int{http.StatusOK}}, resp, &Error{})
		require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	})

	t.Run("empty accepted content type", func(t *testing.T) {
		c := &client{}
		err := c.handleContentType(&internalRequest{}, taskDocumentResponse(http.StatusOK, "", strings.NewReader("")), &Error{})
		require.NoError(t, err)
	})

	t.Run("nil destination", func(t *testing.T) {
		_, _, err := validateNDJSONDestination("GetTaskDocuments", nil)
		require.ErrorContains(t, err, "dst must be a non-nil pointer to a slice")
	})

	t.Run("invalid destination in handler", func(t *testing.T) {
		c := &client{}
		var doc taskDocumentTest
		err := c.handleNDJSONResponse(&internalRequest{functionName: "GetTaskDocuments", withResponse: &doc}, taskDocumentResponse(http.StatusOK, "application/x-ndjson", strings.NewReader("")), &Error{})
		require.ErrorContains(t, err, "dst must point to a slice")
	})
}

func TestGetTaskDocumentsDecodeError(t *testing.T) {
	client := newTaskDocumentTestClient(func(r *http.Request) (*http.Response, error) {
		require.Equal(t, "/tasks/42/documents", r.URL.Path)
		return taskDocumentResponse(http.StatusOK, "application/x-ndjson", strings.NewReader("{\"id\":\"a\"}{\"id\":")), nil
	})

	var docs []taskDocumentTest
	err := client.GetTaskDocuments(42, &docs)
	require.ErrorContains(t, err, "failed to decode NDJSON")
}

func TestGetTaskDocumentsUnmarshalError(t *testing.T) {
	client := newTaskDocumentTestClient(func(r *http.Request) (*http.Response, error) {
		require.Equal(t, "/tasks/42/documents", r.URL.Path)
		return taskDocumentResponse(http.StatusOK, "application/x-ndjson", strings.NewReader(`"not-a-document"`)), nil
	})

	var docs []taskDocumentTest
	err := client.GetTaskDocuments(42, &docs)
	require.ErrorContains(t, err, "failed to unmarshal NDJSON response")
}
