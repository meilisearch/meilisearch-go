package meilisearch

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchRequestBuilder_BuildExampleFromIssue796(t *testing.T) {
	t.Parallel()

	req := NewSearchRequestBuilder().
		WithLimit(20).
		WithFilter("genres = Action").
		WithAttributesToRetrieve("title", "overview").
		WithHitsPerPage(10).
		Build()

	require.Equal(t, int64(20), req.Limit)
	require.Equal(t, "genres = Action", req.Filter)
	require.Equal(t, []string{"title", "overview"}, req.AttributesToRetrieve)
	require.NotNil(t, req.HitsPerPage)
	require.Equal(t, int64(10), *req.HitsPerPage)
}

func TestSearchRequestBuilder_HybridEmbedderDefault(t *testing.T) {
	t.Parallel()

	req := NewSearchRequestBuilder().
		WithHybrid(&SearchRequestHybrid{SemanticRatio: 0.5}).
		Build()

	require.NotNil(t, req.Hybrid)
	require.Equal(t, "default", req.Hybrid.Embedder)
}

func TestSearchRequestBuilder_NilSafeBuild(t *testing.T) {
	t.Parallel()

	var b *SearchRequestBuilder
	req := b.Build()
	require.NotNil(t, req)
}