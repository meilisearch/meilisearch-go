package meilisearch

// SearchRequestBuilder constructs a SearchRequest with a fluent API.
// Optional pointer fields (e.g. HitsPerPage) are managed internally.
type SearchRequestBuilder struct {
	req *SearchRequest
}

// NewSearchRequestBuilder starts a new search request builder.
func NewSearchRequestBuilder() *SearchRequestBuilder {
	return &SearchRequestBuilder{req: &SearchRequest{}}
}

// Build returns the configured SearchRequest after applying SDK validation defaults.
func (b *SearchRequestBuilder) Build() *SearchRequest {
	if b == nil || b.req == nil {
		return &SearchRequest{}
	}
	b.req.validate()
	return b.req
}

func (b *SearchRequestBuilder) WithOffset(v int64) *SearchRequestBuilder {
	b.req.Offset = v
	return b
}

func (b *SearchRequestBuilder) WithLimit(v int64) *SearchRequestBuilder {
	b.req.Limit = v
	return b
}

func (b *SearchRequestBuilder) WithAttributesToRetrieve(attrs ...string) *SearchRequestBuilder {
	b.req.AttributesToRetrieve = attrs
	return b
}

func (b *SearchRequestBuilder) WithAttributesToSearchOn(attrs ...string) *SearchRequestBuilder {
	b.req.AttributesToSearchOn = attrs
	return b
}

func (b *SearchRequestBuilder) WithAttributesToCrop(attrs ...string) *SearchRequestBuilder {
	b.req.AttributesToCrop = attrs
	return b
}

func (b *SearchRequestBuilder) WithCropLength(v int64) *SearchRequestBuilder {
	b.req.CropLength = v
	return b
}

func (b *SearchRequestBuilder) WithCropMarker(v string) *SearchRequestBuilder {
	b.req.CropMarker = v
	return b
}

func (b *SearchRequestBuilder) WithAttributesToHighlight(attrs ...string) *SearchRequestBuilder {
	b.req.AttributesToHighlight = attrs
	return b
}

func (b *SearchRequestBuilder) WithHighlightPreTag(v string) *SearchRequestBuilder {
	b.req.HighlightPreTag = v
	return b
}

func (b *SearchRequestBuilder) WithHighlightPostTag(v string) *SearchRequestBuilder {
	b.req.HighlightPostTag = v
	return b
}

func (b *SearchRequestBuilder) WithMatchingStrategy(v MatchingStrategy) *SearchRequestBuilder {
	b.req.MatchingStrategy = v
	return b
}

func (b *SearchRequestBuilder) WithFilter(v interface{}) *SearchRequestBuilder {
	b.req.Filter = v
	return b
}

func (b *SearchRequestBuilder) WithShowMatchesPosition(v bool) *SearchRequestBuilder {
	b.req.ShowMatchesPosition = v
	return b
}

func (b *SearchRequestBuilder) WithShowRankingScore(v bool) *SearchRequestBuilder {
	b.req.ShowRankingScore = v
	return b
}

func (b *SearchRequestBuilder) WithShowRankingScoreDetails(v bool) *SearchRequestBuilder {
	b.req.ShowRankingScoreDetails = v
	return b
}

func (b *SearchRequestBuilder) WithShowPerformanceDetails(v bool) *SearchRequestBuilder {
	b.req.ShowPerformanceDetails = v
	return b
}

func (b *SearchRequestBuilder) WithFacets(facets ...string) *SearchRequestBuilder {
	b.req.Facets = facets
	return b
}

func (b *SearchRequestBuilder) WithSort(sort ...string) *SearchRequestBuilder {
	b.req.Sort = sort
	return b
}

func (b *SearchRequestBuilder) WithHitsPerPage(v int64) *SearchRequestBuilder {
	b.req.HitsPerPage = Int64Ptr(v)
	return b
}

func (b *SearchRequestBuilder) WithPage(v int64) *SearchRequestBuilder {
	b.req.Page = v
	return b
}

func (b *SearchRequestBuilder) WithIndexUID(v string) *SearchRequestBuilder {
	b.req.IndexUID = v
	return b
}

func (b *SearchRequestBuilder) WithQuery(v string) *SearchRequestBuilder {
	b.req.Query = v
	return b
}

func (b *SearchRequestBuilder) WithDistinct(v string) *SearchRequestBuilder {
	b.req.Distinct = v
	return b
}

func (b *SearchRequestBuilder) WithHybrid(h *SearchRequestHybrid) *SearchRequestBuilder {
	b.req.Hybrid = h
	return b
}

func (b *SearchRequestBuilder) WithRetrieveVectors(v bool) *SearchRequestBuilder {
	b.req.RetrieveVectors = v
	return b
}

func (b *SearchRequestBuilder) WithRankingScoreThreshold(v float64) *SearchRequestBuilder {
	b.req.RankingScoreThreshold = v
	return b
}

func (b *SearchRequestBuilder) WithFederationOptions(o *SearchFederationOptions) *SearchRequestBuilder {
	b.req.FederationOptions = o
	return b
}

func (b *SearchRequestBuilder) WithLocales(locales ...string) *SearchRequestBuilder {
	b.req.Locales = locales
	return b
}