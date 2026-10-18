package meilisearch

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearch_IncludeMetadataHeader(t *testing.T) {
	var got string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("Meili-Include-Metadata")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"hits":[]}`))
	}))
	defer srv.Close()

	idx := New(srv.URL).Index("movies")

	_, err := idx.Search(context.Background(), "q", &SearchRequest{})
	require.NoError(t, err)
	require.Empty(t, got)

	_, err = idx.Search(context.Background(), "q", &SearchRequest{IncludeMetadata: true})
	require.NoError(t, err)
	require.Equal(t, "true", got)
}
