package meilisearch

import (
	"encoding/json"
	"time"
)

type SearchRulesRequest struct {
	Description string      `json:"description,omitempty"`
	Priority    *int        `json:"priority,omitempty"`
	Active      *bool       `json:"active,omitempty"`
	Conditions  []Condition `json:"conditions,omitempty"`
	Actions     []Action    `json:"actions,omitempty"`
}

type SearchRulesResults struct {
	Results []SearchRule `json:"results"`
	Offset  int64        `json:"offset"`
	Limit   int64        `json:"limit"`
	Total   int64        `json:"total"`
}

type SearchRulesParams struct {
	Offset int64              `json:"offset"`
	Limit  int64              `json:"limit"`
	Filter *SearchRulesFilter `json:"filter,omitempty"`
}

type SearchRulesFilter struct {
	AttributePatterns []string `json:"attributePatterns,omitempty"`
	Active            *bool    `json:"active,omitempty"`
}

type SearchRule struct {
	Uid         string      `json:"uid"`
	Description string      `json:"description"`
	Priority    int         `json:"priority"`
	Active      bool        `json:"active"`
	Conditions  []Condition `json:"conditions"`
	Actions     []Action    `json:"actions"`
}

type Condition struct {
	Scope   string     `json:"scope"`
	IsEmpty *bool      `json:"isEmpty,omitempty"`
	Start   *time.Time `json:"start,omitempty"`
	End     *time.Time `json:"end,omitempty"`
}

type Action struct {
	Selector Selector  `json:"selector"`
	Action   ActionDef `json:"action"`
}

type Selector struct {
	IndexUid string `json:"indexUid"`
	ID       string `json:"id,omitempty"`
}

type ActionDef struct {
	Type     string `json:"type"`
	Position int    `json:"position"`
}

// SearchRequest is the request url param needed for a search query.
// This struct will be converted to url param before sent.
//
// Documentation: https://www.meilisearch.com/docs/reference/api/search#search-parameters
type SearchRequest struct {
	Offset                  int64                     `json:"offset,omitempty"`
	Limit                   int64                     `json:"limit,omitempty"`
	AttributesToRetrieve    []string                  `json:"attributesToRetrieve,omitempty"`
	AttributesToSearchOn    []string                  `json:"attributesToSearchOn,omitempty"`
	AttributesToCrop        []string                  `json:"attributesToCrop,omitempty"`
	CropLength              int64                     `json:"cropLength,omitempty"`
	CropMarker              string                    `json:"cropMarker,omitempty"`
	AttributesToHighlight   []string                  `json:"attributesToHighlight,omitempty"`
	HighlightPreTag         string                    `json:"highlightPreTag,omitempty"`
	HighlightPostTag        string                    `json:"highlightPostTag,omitempty"`
	MatchingStrategy        MatchingStrategy          `json:"matchingStrategy,omitempty"`
	Filter                  interface{}               `json:"filter,omitempty"`
	ShowMatchesPosition     bool                      `json:"showMatchesPosition,omitempty"`
	ShowRankingScore        bool                      `json:"showRankingScore,omitempty"`
	ShowRankingScoreDetails bool                      `json:"showRankingScoreDetails,omitempty"`
	ShowPerformanceDetails  bool                      `json:"showPerformanceDetails,omitempty"`
	Facets                  []string                  `json:"facets,omitempty"`
	Sort                    []string                  `json:"sort,omitempty"`
	Vector                  []float32                 `json:"vector,omitempty"`
	HitsPerPage             *int64                    `json:"hitsPerPage,omitempty"`
	Page                    int64                     `json:"page,omitempty"`
	IndexUID                string                    `json:"indexUid,omitempty"`
	Query                   string                    `json:"q"`
	Distinct                string                    `json:"distinct,omitempty"`
	Hybrid                  *SearchRequestHybrid      `json:"hybrid"`
	RetrieveVectors         bool                      `json:"retrieveVectors,omitempty"`
	RankingScoreThreshold   float64                   `json:"rankingScoreThreshold,omitempty"`
	FederationOptions       *SearchFederationOptions  `json:"federationOptions,omitempty"`
	Locales                 []string                  `json:"locales,omitempty"`
	Media                   map[string]any            `json:"media,omitempty"`
	Personalize             *SearchRequestPersonalize `json:"personalize,omitempty"`
}

type SearchFederationOptions struct {
	Weight float64 `json:"weight,omitempty"`
	Remote string  `json:"remote,omitempty"`
}

// SearchRequestPersonalize configures personalized search.
// When set, Meilisearch re-ranks results based on the user profile described in UserContext.
//
// Requires the experimental search personalization feature to be enabled on the server.
// Documentation: https://www.meilisearch.com/docs/capabilities/personalization/getting_started/personalized_search
type SearchRequestPersonalize struct {
	// UserContext is a free-text, natural-language description of the user's
	// preferences, behavior, or intent (e.g. "Prefers compact mechanical keyboards
	// from Keychron, mid-range budget"). The re-ranking model only processes
	// positive signals, so state preferences affirmatively.
	UserContext string `json:"userContext"`
}

type SearchRequestHybrid struct {
	SemanticRatio float64 `json:"semanticRatio,omitempty"`
	Embedder      string  `json:"embedder"`
}

type MultiSearchRequest struct {
	Federation *MultiSearchFederation `json:"federation,omitempty"`
	Queries    []*SearchRequest       `json:"queries"`
}

type MultiSearchFederation struct {
	Offset        int64                             `json:"offset,omitempty"`
	Limit         int64                             `json:"limit,omitempty"`
	FacetsByIndex map[string][]string               `json:"facetsByIndex,omitempty"`
	MergeFacets   *MultiSearchFederationMergeFacets `json:"mergeFacets,omitempty"`
	Distinct      string                            `json:"distinct,omitempty"`
}

type MultiSearchFederationMergeFacets struct {
	MaxValuesPerFacet int `json:"maxValuesPerFacet,omitempty"`
}

// SearchResponse is the response body for search method
type SearchResponse struct {
	Hits               Hits            `json:"hits"`
	EstimatedTotalHits int64           `json:"estimatedTotalHits,omitempty"`
	Offset             int64           `json:"offset,omitempty"`
	Limit              int64           `json:"limit,omitempty"`
	ProcessingTimeMs   int64           `json:"processingTimeMs"`
	Query              string          `json:"query"`
	FacetDistribution  json.RawMessage `json:"facetDistribution,omitempty"`
	TotalHits          int64           `json:"totalHits,omitempty"`
	HitsPerPage        int64           `json:"hitsPerPage,omitempty"`
	Page               int64           `json:"page,omitempty"`
	TotalPages         int64           `json:"totalPages,omitempty"`
	FacetStats         json.RawMessage `json:"facetStats,omitempty"`
	IndexUID           string          `json:"indexUid,omitempty"`
	QueryVector        *[]float32      `json:"queryVector,omitempty"`
	PerformanceDetails map[string]any  `json:"performanceDetails,omitempty"`
}

type MultiSearchResponse struct {
	Results            []SearchResponse           `json:"results,omitempty"`
	Hits               Hits                       `json:"hits,omitempty"`
	ProcessingTimeMs   int64                      `json:"processingTimeMs,omitempty"`
	Offset             int64                      `json:"offset,omitempty"`
	Limit              int64                      `json:"limit,omitempty"`
	EstimatedTotalHits int64                      `json:"estimatedTotalHits,omitempty"`
	SemanticHitCount   int64                      `json:"semanticHitCount,omitempty"`
	FacetDistribution  map[string]json.RawMessage `json:"facetDistribution,omitempty"`
	FacetStats         map[string]json.RawMessage `json:"facetStats,omitempty"`
	RemoteErrors       map[string]*RemoteError    `json:"remoteErrors,omitempty"`
}

type RemoteError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Type    string `json:"type"`
	Link    string `json:"link"`
}

type FacetSearchRequest struct {
	FacetName            string      `json:"facetName,omitempty"`
	FacetQuery           string      `json:"facetQuery,omitempty"`
	Q                    string      `json:"q,omitempty"`
	Filter               interface{} `json:"filter,omitempty"`
	MatchingStrategy     string      `json:"matchingStrategy,omitempty"`
	AttributesToSearchOn []string    `json:"attributesToSearchOn,omitempty"`
	ExhaustiveFacetCount bool        `json:"exhaustiveFacetCount,omitempty"`
}

type FacetSearchResponse struct {
	FacetHits        Hits   `json:"facetHits"`
	FacetQuery       string `json:"facetQuery"`
	ProcessingTimeMs int64  `json:"processingTimeMs"`
}

func (s *SearchRequest) validate() {
	if s.Hybrid != nil && s.Hybrid.Embedder == "" {
		s.Hybrid.Embedder = "default"
	}
}
