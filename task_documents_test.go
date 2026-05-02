package meilisearch

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type taskDocumentTest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/tasks/1/documents", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results":[]}`))
	}))
	defer server.Close()

	client := New(server.URL)
	var docs []taskDocumentTest
	err := client.GetTaskDocuments(1, &docs)
	require.ErrorContains(t, err, `unexpected Content-Type "application/json"`)
}

func TestGetTaskDocumentsDecodesNDJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/tasks/42/documents", r.URL.Path)
		require.Empty(t, r.URL.RawQuery)
		w.Header().Set("Content-Type", "application/x-ndjson; charset=utf-8")
		_, _ = w.Write([]byte("{\"id\":\"1\",\"name\":\"Alice\"}\n{\"id\":\"2\",\"name\":\"Bob\"}\n"))
	}))
	defer server.Close()

	client := New(server.URL)
	var docs []taskDocumentTest
	err := client.GetTaskDocuments(42, &docs)
	require.NoError(t, err)
	require.Equal(t, []taskDocumentTest{
		{ID: "1", Name: "Alice"},
		{ID: "2", Name: "Bob"},
	}, docs)
}

func TestGetTaskDocumentsDecodesConcatenatedNDJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/tasks/42/documents", r.URL.Path)
		require.Empty(t, r.URL.RawQuery)
		w.Header().Set("Content-Type", "application/x-ndjson; charset=utf-8")
		_, _ = w.Write([]byte("{\"id\":\"1\",\"name\":\"Alice\"}{\"id\":\"2\",\"name\":\"Bob\"}{\"id\":\"3\",\"name\":\"Carol\"}"))
	}))
	defer server.Close()

	client := New(server.URL)
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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, GzipEncoding.String(), r.Header.Get("Accept-Encoding"))
		w.Header().Set("Content-Type", "application/x-ndjson")
		w.Header().Set("Content-Encoding", GzipEncoding.String())

		body, err := newEncoding(GzipEncoding, DefaultCompression).Encode(strings.NewReader("{\"id\":\"1\",\"name\":\"Alice\"}\n"))
		require.NoError(t, err)
		defer func() {
			_ = body.Close()
		}()

		_, _ = io.Copy(w, body)
	}))
	defer server.Close()

	client := New(server.URL, WithContentEncoding(GzipEncoding, DefaultCompression))
	var docs []taskDocumentTest
	err := client.GetTaskDocuments(1, &docs)
	require.NoError(t, err)
	require.Equal(t, []taskDocumentTest{{ID: "1", Name: "Alice"}}, docs)
}
