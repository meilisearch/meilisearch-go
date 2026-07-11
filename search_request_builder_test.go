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

// TestSearchRequestBuilder_AllSetters is a single table-driven suite that covers
// every With* setter on SearchRequestBuilder (maintainer request on PR #799).
func TestSearchRequestBuilder_AllSetters(t *testing.T) {
	t.Parallel()

	hybrid := &SearchRequestHybrid{SemanticRatio: 0.3, Embedder: "custom"}
	federation := &SearchFederationOptions{Weight: 0.7, Remote: "remote-a"}
	personalize := &SearchRequestPersonalize{UserContext: "Prefers mechanical keyboards"}
	vec := []float32{0.1, 0.2, 0.3}
	media := map[string]any{"query": "cat", "content": "data:image/png;base64,..."}

	tests := []struct {
		name   string
		apply  func(*SearchRequestBuilder) *SearchRequestBuilder
		assert func(*testing.T, *SearchRequest)
	}{
		{
			name:  "WithOffset",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder { return b.WithOffset(15) },
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, int64(15), req.Offset)
			},
		},
		{
			name:  "WithLimit",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder { return b.WithLimit(25) },
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, int64(25), req.Limit)
			},
		},
		{
			name: "WithAttributesToRetrieve",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder {
				return b.WithAttributesToRetrieve("title", "overview")
			},
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, []string{"title", "overview"}, req.AttributesToRetrieve)
			},
		},
		{
			name: "WithAttributesToSearchOn",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder {
				return b.WithAttributesToSearchOn("title")
			},
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, []string{"title"}, req.AttributesToSearchOn)
			},
		},
		{
			name: "WithAttributesToCrop",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder {
				return b.WithAttributesToCrop("description")
			},
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, []string{"description"}, req.AttributesToCrop)
			},
		},
		{
			name:  "WithCropLength",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder { return b.WithCropLength(80) },
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, int64(80), req.CropLength)
			},
		},
		{
			name:  "WithCropMarker",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder { return b.WithCropMarker("…") },
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, "…", req.CropMarker)
			},
		},
		{
			name: "WithAttributesToHighlight",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder {
				return b.WithAttributesToHighlight("title")
			},
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, []string{"title"}, req.AttributesToHighlight)
			},
		},
		{
			name: "WithHighlightPreTag",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder {
				return b.WithHighlightPreTag("<em>")
			},
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, "<em>", req.HighlightPreTag)
			},
		},
		{
			name: "WithHighlightPostTag",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder {
				return b.WithHighlightPostTag("</em>")
			},
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, "</em>", req.HighlightPostTag)
			},
		},
		{
			name: "WithMatchingStrategy",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder {
				return b.WithMatchingStrategy(All)
			},
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, All, req.MatchingStrategy)
			},
		},
		{
			name: "WithFilter",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder {
				return b.WithFilter("genres = Action")
			},
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, "genres = Action", req.Filter)
			},
		},
		{
			name: "WithShowMatchesPosition",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder {
				return b.WithShowMatchesPosition(true)
			},
			assert: func(t *testing.T, req *SearchRequest) {
				require.True(t, req.ShowMatchesPosition)
			},
		},
		{
			name: "WithShowRankingScore",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder {
				return b.WithShowRankingScore(true)
			},
			assert: func(t *testing.T, req *SearchRequest) {
				require.True(t, req.ShowRankingScore)
			},
		},
		{
			name: "WithShowRankingScoreDetails",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder {
				return b.WithShowRankingScoreDetails(true)
			},
			assert: func(t *testing.T, req *SearchRequest) {
				require.True(t, req.ShowRankingScoreDetails)
			},
		},
		{
			name: "WithShowPerformanceDetails",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder {
				return b.WithShowPerformanceDetails(true)
			},
			assert: func(t *testing.T, req *SearchRequest) {
				require.True(t, req.ShowPerformanceDetails)
			},
		},
		{
			name:  "WithFacets",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder { return b.WithFacets("genres") },
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, []string{"genres"}, req.Facets)
			},
		},
		{
			name:  "WithSort",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder { return b.WithSort("price:asc") },
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, []string{"price:asc"}, req.Sort)
			},
		},
		{
			name:  "WithHitsPerPage",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder { return b.WithHitsPerPage(12) },
			assert: func(t *testing.T, req *SearchRequest) {
				require.NotNil(t, req.HitsPerPage)
				require.Equal(t, int64(12), *req.HitsPerPage)
			},
		},
		{
			name:  "WithPage",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder { return b.WithPage(3) },
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, int64(3), req.Page)
			},
		},
		{
			name:  "WithIndexUID",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder { return b.WithIndexUID("movies") },
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, "movies", req.IndexUID)
			},
		},
		{
			name:  "WithQuery",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder { return b.WithQuery("prince") },
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, "prince", req.Query)
			},
		},
		{
			name:  "WithDistinct",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder { return b.WithDistinct("sku") },
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, "sku", req.Distinct)
			},
		},
		{
			name:  "WithHybrid",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder { return b.WithHybrid(hybrid) },
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, hybrid, req.Hybrid)
			},
		},
		{
			name: "WithRetrieveVectors",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder {
				return b.WithRetrieveVectors(true)
			},
			assert: func(t *testing.T, req *SearchRequest) {
				require.True(t, req.RetrieveVectors)
			},
		},
		{
			name: "WithRankingScoreThreshold",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder {
				return b.WithRankingScoreThreshold(0.42)
			},
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, 0.42, req.RankingScoreThreshold)
			},
		},
		{
			name: "WithFederationOptions",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder {
				return b.WithFederationOptions(federation)
			},
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, federation, req.FederationOptions)
			},
		},
		{
			name:  "WithLocales",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder { return b.WithLocales("eng", "fra") },
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, []string{"eng", "fra"}, req.Locales)
			},
		},
		{
			name:  "WithVector",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder { return b.WithVector(vec) },
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, vec, req.Vector)
			},
		},
		{
			name:  "WithMedia",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder { return b.WithMedia(media) },
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, media, req.Media)
			},
		},
		{
			name: "WithPersonalize",
			apply: func(b *SearchRequestBuilder) *SearchRequestBuilder {
				return b.WithPersonalize(personalize)
			},
			assert: func(t *testing.T, req *SearchRequest) {
				require.Equal(t, personalize, req.Personalize)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := tt.apply(NewSearchRequestBuilder()).Build()
			tt.assert(t, req)
		})
	}
}

func TestSearchRequestBuilder_FullChainSetsAllFields(t *testing.T) {
	t.Parallel()

	hybrid := &SearchRequestHybrid{SemanticRatio: 0.5, Embedder: "default"}
	federation := &SearchFederationOptions{Weight: 1.0}
	personalize := &SearchRequestPersonalize{UserContext: "likes sci-fi"}
	vec := []float32{0.5, 0.5}
	media := map[string]any{"query": "dog"}

	req := NewSearchRequestBuilder().
		WithOffset(1).
		WithLimit(5).
		WithAttributesToRetrieve("title").
		WithAttributesToSearchOn("title").
		WithAttributesToCrop("overview").
		WithCropLength(40).
		WithCropMarker("...").
		WithAttributesToHighlight("title").
		WithHighlightPreTag("<b>").
		WithHighlightPostTag("</b>").
		WithMatchingStrategy(Last).
		WithFilter("year > 2000").
		WithShowMatchesPosition(true).
		WithShowRankingScore(true).
		WithShowRankingScoreDetails(true).
		WithShowPerformanceDetails(true).
		WithFacets("genres").
		WithSort("title:asc").
		WithHitsPerPage(10).
		WithPage(2).
		WithIndexUID("books").
		WithQuery("prince").
		WithDistinct("isbn").
		WithHybrid(hybrid).
		WithRetrieveVectors(true).
		WithRankingScoreThreshold(0.1).
		WithFederationOptions(federation).
		WithLocales("eng").
		WithVector(vec).
		WithMedia(media).
		WithPersonalize(personalize).
		Build()

	require.Equal(t, int64(1), req.Offset)
	require.Equal(t, int64(5), req.Limit)
	require.Equal(t, []string{"title"}, req.AttributesToRetrieve)
	require.Equal(t, []string{"title"}, req.AttributesToSearchOn)
	require.Equal(t, []string{"overview"}, req.AttributesToCrop)
	require.Equal(t, int64(40), req.CropLength)
	require.Equal(t, "...", req.CropMarker)
	require.Equal(t, []string{"title"}, req.AttributesToHighlight)
	require.Equal(t, "<b>", req.HighlightPreTag)
	require.Equal(t, "</b>", req.HighlightPostTag)
	require.Equal(t, Last, req.MatchingStrategy)
	require.Equal(t, "year > 2000", req.Filter)
	require.True(t, req.ShowMatchesPosition)
	require.True(t, req.ShowRankingScore)
	require.True(t, req.ShowRankingScoreDetails)
	require.True(t, req.ShowPerformanceDetails)
	require.Equal(t, []string{"genres"}, req.Facets)
	require.Equal(t, []string{"title:asc"}, req.Sort)
	require.NotNil(t, req.HitsPerPage)
	require.Equal(t, int64(10), *req.HitsPerPage)
	require.Equal(t, int64(2), req.Page)
	require.Equal(t, "books", req.IndexUID)
	require.Equal(t, "prince", req.Query)
	require.Equal(t, "isbn", req.Distinct)
	require.Equal(t, hybrid, req.Hybrid)
	require.True(t, req.RetrieveVectors)
	require.Equal(t, 0.1, req.RankingScoreThreshold)
	require.Equal(t, federation, req.FederationOptions)
	require.Equal(t, []string{"eng"}, req.Locales)
	require.Equal(t, vec, req.Vector)
	require.Equal(t, media, req.Media)
	require.Equal(t, personalize, req.Personalize)
}
