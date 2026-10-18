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

func TestMultiSearch_FederatedMetadata(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "true", r.Header.Get("Meili-Include-Metadata"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"hits":[],"metadata":[{"queryUid":"abc","indexUid":"movies","primaryKey":"id"}]}`))
	}))
	defer srv.Close()

	resp, err := New(srv.URL).MultiSearch(context.Background(), &MultiSearchRequest{
		Federation:      &MultiSearchFederation{},
		Queries:         []*SearchRequest{{IndexUID: "movies", Query: "q"}},
		IncludeMetadata: true,
	})
	require.NoError(t, err)
	require.Len(t, resp.Metadata, 1)
	require.Equal(t, "movies", resp.Metadata[0].IndexUID)
	require.Equal(t, "id", resp.Metadata[0].PrimaryKey)
}
