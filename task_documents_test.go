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

		body, err := newEncoding(GzipEncoding, DefaultCompression).Encode(strings.NewReader("{\"id\":\"1\",\"name\":\"Alice\"}\n"))
		require.NoError(t, err)

		resp := taskDocumentResponse(http.StatusOK, "application/x-ndjson", body)
		resp.Header.Set("Content-Encoding", GzipEncoding.String())
		return resp, nil
	}, WithContentEncoding(GzipEncoding, DefaultCompression))

	var docs []taskDocumentTest
	err := client.GetTaskDocuments(1, &docs)
	require.NoError(t, err)
	require.Equal(t, []taskDocumentTest{{ID: "1", Name: "Alice"}}, docs)
}
