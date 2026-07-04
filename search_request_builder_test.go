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

func TestSearchRequestBuilder_WithVector(t *testing.T) {
	t.Parallel()

	vec := []float32{0.1, 0.2, 0.3}
	req := NewSearchRequestBuilder().
		WithVector(vec).
		Build()

	require.Equal(t, vec, req.Vector)
}

func TestSearchRequestBuilder_WithMedia(t *testing.T) {
	t.Parallel()

	media := map[string]any{"query": "cat", "content": "data:image/png;base64,..."}
	req := NewSearchRequestBuilder().
		WithMedia(media).
		Build()

	require.Equal(t, media, req.Media)
}

func TestSearchRequestBuilder_FullChainIncludesAllSetters(t *testing.T) {
	t.Parallel()

	vec := []float32{0.5, 0.5}
	media := map[string]any{"query": "dog"}
	req := NewSearchRequestBuilder().
		WithQuery("prince").
		WithLimit(5).
		WithVector(vec).
		WithMedia(media).
		WithLocales("eng").
		Build()

	require.Equal(t, "prince", req.Query)
	require.Equal(t, int64(5), req.Limit)
	require.Equal(t, vec, req.Vector)
	require.Equal(t, media, req.Media)
	require.Equal(t, []string{"eng"}, req.Locales)
}
