package meilisearch

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	contentTypeJSON   string = "application/json"
	contentTypeNDJSON string = "application/x-ndjson"
	contentTypeCSV    string = "text/csv"
	nullBody                 = "null"
)

// Network represents the Meilisearch network configuration.
// Each field is wrapped in an Opt so it can be explicitly included,
// set to JSON null, or omitted entirely.
type Network struct {
	Self     Opt[string]                 `json:"self,omitempty"`
	Remotes  Opt[map[string]Opt[Remote]] `json:"remotes,omitempty"`
	Sharding Opt[bool]                   `json:"sharding,omitempty"`
}

func (n Network) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)

	if n.Self.Valid() {
		m["self"] = n.Self.Value
	} else if n.Self.Null() {
		m["self"] = nil
	}

	if n.Remotes.Valid() {
		m["remotes"] = n.Remotes.Value
	} else if n.Remotes.Null() {
		m["remotes"] = nil
	}

	if n.Sharding.Valid() {
		m["sharding"] = n.Sharding.Value
	} else if n.Sharding.Null() {
		m["sharding"] = nil
	}

	return json.Marshal(m)
}

// Remote describes a single remote Meilisearch node.
// Each field is wrapped in an Opt so it can be explicitly included,
// set to JSON null, or omitted entirely.
type Remote struct {
	URL          Opt[string] `json:"url"`
	SearchAPIKey Opt[string] `json:"searchApiKey,omitempty"`
	WriteAPIKey  Opt[string] `json:"writeApiKey,omitempty"`
}

func (r Remote) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)

	if r.URL.Valid() {
		m["url"] = r.URL.Value
	} else if r.URL.Null() {
		m["url"] = nil
	}

	if r.SearchAPIKey.Valid() {
		m["searchApiKey"] = r.SearchAPIKey.Value
	} else if r.SearchAPIKey.Null() {
		m["searchApiKey"] = nil
	}

	if r.WriteAPIKey.Valid() {
		m["writeApiKey"] = r.WriteAPIKey.Value
	} else if r.WriteAPIKey.Null() {
		m["writeApiKey"] = nil
	}

	return json.Marshal(m)
}

type UpdateWebhookRequest struct {
	URL     string            `json:"url,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

type Webhook struct {
	UUID       string            `json:"uuid"`
	IsEditable bool              `json:"isEditable"`
	URL        string            `json:"url"`
	Headers    map[string]string `json:"headers"`
}

type WebhookResults struct {
	Result []*Webhook `json:"results"`
}

type AddWebhookRequest struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
}

type IndexConfig struct {
	// Uid is the unique identifier of a given index.
	Uid string
	// PrimaryKey is optional
	PrimaryKey string
}

type IndexResult struct {
	UID        string    `json:"uid"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	PrimaryKey string    `json:"primaryKey,omitempty"`
	IndexManager
}

// IndexesResults return of multiple indexes is wrap in a IndexesResults
type IndexesResults struct {
	Results []*IndexResult `json:"results"`
	Offset  int64          `json:"offset"`
	Limit   int64          `json:"limit"`
	Total   int64          `json:"total"`
}

type IndexesQuery struct {
	Limit  int64
	Offset int64
}

type AttributeRule struct {
	AttributePatterns []string          `json:"attributePatterns"`
	Features          AttributeFeatures `json:"features"`
}

type AttributeFeatures struct {
	FacetSearch bool           `json:"facetSearch"`
	Filter      FilterFeatures `json:"filter"`
}

type FilterFeatures struct {
	Equality   bool `json:"equality"`
	Comparison bool `json:"comparison"`
}

// Settings is the type that represents the settings in meilisearch
type Settings struct {
	RankingRules         []string               `json:"rankingRules,omitempty"`
	DistinctAttribute    *string                `json:"distinctAttribute,omitempty"`
	SearchableAttributes []string               `json:"searchableAttributes,omitempty"`
	Dictionary           []string               `json:"dictionary,omitempty"`
	SearchCutoffMs       int64                  `json:"searchCutoffMs,omitempty"`
	ProximityPrecision   ProximityPrecisionType `json:"proximityPrecision,omitempty"`
	SeparatorTokens      []string               `json:"separatorTokens,omitempty"`
	NonSeparatorTokens   []string               `json:"nonSeparatorTokens,omitempty"`
	DisplayedAttributes  []string               `json:"displayedAttributes,omitempty"`
	StopWords            []string               `json:"stopWords,omitempty"`
	Synonyms             map[string][]string    `json:"synonyms,omitempty"`
	FilterableAttributes []string               `json:"filterableAttributes,omitempty"`
	SortableAttributes   []string               `json:"sortableAttributes,omitempty"`
	LocalizedAttributes  []*LocalizedAttributes `json:"localizedAttributes,omitempty"`
	TypoTolerance        *TypoTolerance         `json:"typoTolerance,omitempty"`
	Pagination           *Pagination            `json:"pagination,omitempty"`
	Faceting             *Faceting              `json:"faceting,omitempty"`
	Embedders            map[string]Embedder    `json:"embedders,omitempty"`
	PrefixSearch         *string                `json:"prefixSearch,omitempty"`
	FacetSearch          bool                   `json:"facetSearch,omitempty"`
	Chat                 *Chat                  `json:"chat,omitempty"`
}

type LocalizedAttributes struct {
	Locales           []string `json:"locales,omitempty"`
	AttributePatterns []string `json:"attributePatterns,omitempty"`
}

// TypoTolerance is the type that represents the typo tolerance setting in meilisearch
type TypoTolerance struct {
	Enabled             bool                `json:"enabled"`
	MinWordSizeForTypos MinWordSizeForTypos `json:"minWordSizeForTypos"`
	DisableOnWords      []string            `json:"disableOnWords"`
	DisableOnAttributes []string            `json:"disableOnAttributes"`
	DisableOnNumbers    bool                `json:"disableOnNumbers"`
}

// MinWordSizeForTypos is the type that represents the minWordSizeForTypos setting in the typo tolerance setting in meilisearch
type MinWordSizeForTypos struct {
	OneTypo  int64 `json:"oneTypo"`
	TwoTypos int64 `json:"twoTypos"`
}

// Pagination is the type that represents the pagination setting in meilisearch
type Pagination struct {
	MaxTotalHits int64 `json:"maxTotalHits"`
}

// Faceting is the type that represents the faceting setting in meilisearch
type Faceting struct {
	MaxValuesPerFacet int64 `json:"maxValuesPerFacet"`
	// SortFacetValuesBy index_name: alpha|count
	SortFacetValuesBy map[string]SortFacetType `json:"sortFacetValuesBy"`
}

type Chat struct {
	Description              string            `json:"description,omitempty"`
	DocumentTemplate         string            `json:"documentTemplate,omitempty"`
	DocumentTemplateMaxBytes int               `json:"documentTemplateMaxBytes,omitempty"`
	SearchParameters         *SearchParameters `json:"searchParameters,omitempty"`
}

type SearchParameters struct {
	Limit                 int64                `json:"limit,omitempty"`
	AttributesToSearchOn  []string             `json:"attributesToSearchOn,omitempty"`
	MatchingStrategy      MatchingStrategy     `json:"matchingStrategy,omitempty"`
	Sort                  []string             `json:"sort,omitempty"`
	Distinct              string               `json:"distinct,omitempty"`
	Hybrid                *SearchRequestHybrid `json:"hybrid,omitempty"`
	RankingScoreThreshold float64              `json:"rankingScoreThreshold,omitempty"`
}

// Embedder represents a unified configuration for various embedder types.
//
// Specs: https://www.meilisearch.com/docs/reference/api/settings#body-21
type Embedder struct {
	Source EmbedderSource `json:"source"` // The type of embedder: "openAi", "huggingFace", "userProvided", "rest", "ollama"
	// URL Meilisearch queries url to generate vector embeddings for queries and documents.
	// url must point to a REST-compatible embedder. You may also use url to work with proxies, such as when targeting openAi from behind a proxy.
	URL    string `json:"url,omitempty"`    // Optional for "openAi", "rest", "ollama"
	APIKey string `json:"apiKey,omitempty"` // Optional for "openAi", "rest", "ollama"
	// Model The model your embedder uses when generating vectors.
	// These are the officially supported models Meilisearch supports:
	//
	// - openAi: text-embedding-3-small, text-embedding-3-large, openai-text-embedding-ada-002
	//
	// - huggingFace: BAAI/bge-base-en-v1.5
	//
	// Other models, such as HuggingFace’s BERT models or those provided by Ollama and REST embedders
	// may also be compatible with Meilisearch.
	//
	// HuggingFace’s BERT models: https://huggingface.co/models?other=bert
	Model string `json:"model,omitempty"` // Optional for "openAi", "huggingFace", "ollama"
	// DocumentTemplate is a string containing a Liquid template. Meillisearch interpolates the template for each
	// document and sends the resulting text to the embedder. The embedder then generates document vectors based on this text.
	DocumentTemplate string `json:"documentTemplate,omitempty"` // Optional for most embedders
	// DocumentTemplateMaxBytes The maximum size of a rendered document template.
	//Longer texts are truncated to fit the configured limit.
	//
	// documentTemplateMaxBytes must be an integer. It defaults to 400.
	DocumentTemplateMaxBytes int                    `json:"documentTemplateMaxBytes,omitempty"`
	Dimensions               int                    `json:"dimensions,omitempty"`   // Optional for "openAi", "rest", "userProvided", "ollama"
	Revision                 string                 `json:"revision,omitempty"`     // Optional for "huggingFace"
	Distribution             *Distribution          `json:"distribution,omitempty"` // Optional for all embedders
	Request                  map[string]interface{} `json:"request,omitempty"`      // Optional for "rest"
	Response                 map[string]interface{} `json:"response,omitempty"`     // Optional for "rest"
	Headers                  map[string]string      `json:"headers,omitempty"`      // Optional for "rest"
	BinaryQuantized          bool                   `json:"binaryQuantized,omitempty"`
	Pooling                  EmbedderPooling        `json:"pooling,omitempty"`
	IndexingEmbedder         *Embedder              `json:"indexingEmbedder,omitempty"` // For Composite
	SearchEmbedder           *Embedder              `json:"searchEmbedder,omitempty"`   // For Composite
	IndexingFragments        map[string]Fragment    `json:"indexingFragments,omitempty"`
	SearchFragments          map[string]Fragment    `json:"searchFragments,omitempty"`
}

type Fragment struct {
	Value map[string]any `json:"value,omitempty"`
}

// Distribution represents a statistical distribution with mean and standard deviation (sigma).
type Distribution struct {
	Mean  float64 `json:"mean"`  // Mean of the distribution
	Sigma float64 `json:"sigma"` // Sigma (standard deviation) of the distribution
}

// Version is the type that represents the versions in meilisearch
type Version struct {
	CommitSha  string `json:"commitSha"`
	CommitDate string `json:"commitDate"`
	PkgVersion string `json:"pkgVersion"`
}

// StatsIndex is the type that represent the stats of an index in meilisearch
type StatsIndex struct {
	NumberOfDocuments         int64            `json:"numberOfDocuments"`
	IsIndexing                bool             `json:"isIndexing"`
	FieldDistribution         map[string]int64 `json:"fieldDistribution"`
	RawDocumentDbSize         int64            `json:"rawDocumentDbSize"`
	AvgDocumentSize           int64            `json:"avgDocumentSize"`
	NumberOfEmbeddedDocuments int64            `json:"numberOfEmbeddedDocuments"`
	NumberOfEmbeddings        int64            `json:"numberOfEmbeddings"`
}

// Stats is the type that represent all stats
type Stats struct {
	DatabaseSize     int64                 `json:"databaseSize"`
	UsedDatabaseSize int64                 `json:"usedDatabaseSize"`
	LastUpdate       time.Time             `json:"lastUpdate"`
	Indexes          map[string]StatsIndex `json:"indexes"`
}

// Task indicates information about a task resource
//
// Documentation: https://www.meilisearch.com/docs/learn/advanced/asynchronous_operations
type Task struct {
	Status      TaskStatus          `json:"status"`
	UID         int64               `json:"uid,omitempty"`
	TaskUID     int64               `json:"taskUid,omitempty"`
	IndexUID    string              `json:"indexUid"`
	Type        TaskType            `json:"type"`
	Error       meilisearchApiError `json:"error,omitempty"`
	TaskNetwork TaskNetwork         `json:"network,omitempty"`
	Duration    string              `json:"duration,omitempty"`
	EnqueuedAt  time.Time           `json:"enqueuedAt"`
	StartedAt   time.Time           `json:"startedAt,omitempty"`
	FinishedAt  time.Time           `json:"finishedAt,omitempty"`
	Details     Details             `json:"details,omitempty"`
	CanceledBy  int64               `json:"canceledBy,omitempty"`
}

// TaskNetwork indicates information about a task network
//
// Documentation: https://www.meilisearch.com/docs/reference/api/tasks#network
type TaskNetwork struct {
	Origin  *Origin                `json:"origin,omitempty"`
	Remotes map[string]*TaskRemote `json:"remotes,omitempty"`
}

type Origin struct {
	RemoteName string `json:"remoteName,omitempty"`
	TaskUID    string `json:"taskUid,omitempty"`
}

type TaskRemote struct {
	TaskUID *string `json:"task_uid,omitempty"`
	Error   *string `json:"error,omitempty"`
}

// TaskInfo indicates information regarding a task returned by an asynchronous method
//
// Documentation: https://www.meilisearch.com/docs/reference/api/tasks#tasks
type TaskInfo struct {
	Status     TaskStatus `json:"status"`
	TaskUID    int64      `json:"taskUid"`
	IndexUID   string     `json:"indexUid"`
	Type       TaskType   `json:"type"`
	EnqueuedAt time.Time  `json:"enqueuedAt"`
}

// TasksQuery is a list of filter available to send as query parameters
type TasksQuery struct {
	UIDS             []int64
	Limit            int64
	From             int64
	IndexUIDS        []string
	Statuses         []TaskStatus
	Types            []TaskType
	CanceledBy       []int64
	BeforeEnqueuedAt time.Time
	AfterEnqueuedAt  time.Time
	BeforeStartedAt  time.Time
	AfterStartedAt   time.Time
	BeforeFinishedAt time.Time
	AfterFinishedAt  time.Time
	Reverse          bool
}

// CancelTasksQuery is a list of filter available to send as query parameters
type CancelTasksQuery struct {
	UIDS             []int64
	IndexUIDS        []string
	Statuses         []TaskStatus
	Types            []TaskType
	BeforeEnqueuedAt time.Time
	AfterEnqueuedAt  time.Time
	BeforeStartedAt  time.Time
	AfterStartedAt   time.Time
}

// DeleteTasksQuery is a list of filter available to send as query parameters
type DeleteTasksQuery struct {
	UIDS             []int64
	IndexUIDS        []string
	Statuses         []TaskStatus
	Types            []TaskType
	CanceledBy       []int64
	BeforeEnqueuedAt time.Time
	AfterEnqueuedAt  time.Time
	BeforeStartedAt  time.Time
	AfterStartedAt   time.Time
	BeforeFinishedAt time.Time
	AfterFinishedAt  time.Time
}

type Details struct {
	ReceivedDocuments    int64               `json:"receivedDocuments,omitempty"`
	IndexedDocuments     int64               `json:"indexedDocuments,omitempty"`
	DeletedDocuments     int64               `json:"deletedDocuments,omitempty"`
	PrimaryKey           string              `json:"primaryKey,omitempty"`
	ProvidedIds          int64               `json:"providedIds,omitempty"`
	RankingRules         []string            `json:"rankingRules,omitempty"`
	DistinctAttribute    *string             `json:"distinctAttribute,omitempty"`
	SearchableAttributes []string            `json:"searchableAttributes,omitempty"`
	DisplayedAttributes  []string            `json:"displayedAttributes,omitempty"`
	StopWords            []string            `json:"stopWords,omitempty"`
	Synonyms             map[string][]string `json:"synonyms,omitempty"`
	FilterableAttributes []interface{}       `json:"filterableAttributes,omitempty"`
	SortableAttributes   []string            `json:"sortableAttributes,omitempty"`
	TypoTolerance        *TypoTolerance      `json:"typoTolerance,omitempty"`
	Pagination           *Pagination         `json:"pagination,omitempty"`
	Faceting             *Faceting           `json:"faceting,omitempty"`
	MatchedTasks         int64               `json:"matchedTasks,omitempty"`
	CanceledTasks        int64               `json:"canceledTasks,omitempty"`
	DeletedTasks         int64               `json:"deletedTasks,omitempty"`
	OriginalFilter       string              `json:"originalFilter,omitempty"`
	Swaps                []SwapIndexesParams `json:"swaps,omitempty"`
	DumpUid              string              `json:"dumpUid,omitempty"`
}

// TaskResult return of multiple tasks is wrap in a TaskResult
type TaskResult struct {
	Results []Task `json:"results"`
	Limit   int64  `json:"limit"`
	From    int64  `json:"from"`
	Next    int64  `json:"next"`
	Total   int64  `json:"total"`
}

// Key allow the user to connect to the meilisearch instance
//
// Documentation: https://www.meilisearch.com/docs/learn/security/master_api_keys#protecting-a-meilisearch-instance
type Key struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Key         string    `json:"key,omitempty"`
	UID         string    `json:"uid,omitempty"`
	Actions     []string  `json:"actions,omitempty"`
	Indexes     []string  `json:"indexes,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

// KeyParsed this structure is used to send the exact ISO-8601 time format managed by meilisearch
type KeyParsed struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	UID         string   `json:"uid,omitempty"`
	Actions     []string `json:"actions,omitempty"`
	Indexes     []string `json:"indexes,omitempty"`
	ExpiresAt   *string  `json:"expiresAt"`
}

// KeyUpdate this structure is used to update a Key
type KeyUpdate struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// KeysResults return of multiple keys is wrap in a KeysResults
type KeysResults struct {
	Results []Key `json:"results"`
	Offset  int64 `json:"offset"`
	Limit   int64 `json:"limit"`
	Total   int64 `json:"total"`
}

type KeysQuery struct {
	Limit  int64
	Offset int64
}

// TenantTokenOptions information to create a tenant token
//
// ExpiresAt is a time.Time when the key will expire.
// Note that if an ExpiresAt value is included it should be in UTC time.
// ApiKey is the API key parent of the token.
type TenantTokenOptions struct {
	APIKey    string
	ExpiresAt time.Time
}

// TenantTokenClaims custom Claims structure to create a Tenant Token
type TenantTokenClaims struct {
	APIKeyUID   string      `json:"apiKeyUid"`
	SearchRules interface{} `json:"searchRules"`
	jwt.RegisteredClaims
}

//
// Request/Response
//

// CreateIndexRequest is the request body for create index method
type CreateIndexRequest struct {
	UID        string `json:"uid,omitempty"`
	PrimaryKey string `json:"primaryKey,omitempty"`
}

// SearchRequest is the request url param needed for a search query.
// This struct will be converted to url param before sent.
//
// Documentation: https://www.meilisearch.com/docs/reference/api/search#search-parameters
type SearchRequest struct {
	Offset                  int64                    `json:"offset,omitempty"`
	Limit                   int64                    `json:"limit,omitempty"`
	AttributesToRetrieve    []string                 `json:"attributesToRetrieve,omitempty"`
	AttributesToSearchOn    []string                 `json:"attributesToSearchOn,omitempty"`
	AttributesToCrop        []string                 `json:"attributesToCrop,omitempty"`
	CropLength              int64                    `json:"cropLength,omitempty"`
	CropMarker              string                   `json:"cropMarker,omitempty"`
	AttributesToHighlight   []string                 `json:"attributesToHighlight,omitempty"`
	HighlightPreTag         string                   `json:"highlightPreTag,omitempty"`
	HighlightPostTag        string                   `json:"highlightPostTag,omitempty"`
	MatchingStrategy        MatchingStrategy         `json:"matchingStrategy,omitempty"`
	Filter                  interface{}              `json:"filter,omitempty"`
	ShowMatchesPosition     bool                     `json:"showMatchesPosition,omitempty"`
	ShowRankingScore        bool                     `json:"showRankingScore,omitempty"`
	ShowRankingScoreDetails bool                     `json:"showRankingScoreDetails,omitempty"`
	Facets                  []string                 `json:"facets,omitempty"`
	Sort                    []string                 `json:"sort,omitempty"`
	Vector                  []float32                `json:"vector,omitempty"`
	HitsPerPage             int64                    `json:"hitsPerPage,omitempty"`
	Page                    int64                    `json:"page,omitempty"`
	IndexUID                string                   `json:"indexUid,omitempty"`
	Query                   string                   `json:"q"`
	Distinct                string                   `json:"distinct,omitempty"`
	Hybrid                  *SearchRequestHybrid     `json:"hybrid"`
	RetrieveVectors         bool                     `json:"retrieveVectors,omitempty"`
	RankingScoreThreshold   float64                  `json:"rankingScoreThreshold,omitempty"`
	FederationOptions       *SearchFederationOptions `json:"federationOptions,omitempty"`
	Locales                 []string                 `json:"locales,omitempty"`
	Media                   map[string]any           `json:"media,omitempty"`
}

type SearchFederationOptions struct {
	Weight float64 `json:"weight,omitempty"`
	Remote string  `json:"remote,omitempty"`
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
	FacetName            string   `json:"facetName,omitempty"`
	FacetQuery           string   `json:"facetQuery,omitempty"`
	Q                    string   `json:"q,omitempty"`
	Filter               string   `json:"filter,omitempty"`
	MatchingStrategy     string   `json:"matchingStrategy,omitempty"`
	AttributesToSearchOn []string `json:"attributesToSearchOn,omitempty"`
	ExhaustiveFacetCount bool     `json:"exhaustiveFacetCount,omitempty"`
}

type FacetSearchResponse struct {
	FacetHits        Hits   `json:"facetHits"`
	FacetQuery       string `json:"facetQuery"`
	ProcessingTimeMs int64  `json:"processingTimeMs"`
}

// DocumentQuery is the request body get one documents method
type DocumentQuery struct {
	Fields          []string `json:"fields,omitempty"`
	RetrieveVectors bool     `json:"retrieveVectors,omitempty"`
}

// DocumentsQuery is the request body for list documents method
type DocumentsQuery struct {
	Offset          int64       `json:"offset,omitempty"`
	Limit           int64       `json:"limit,omitempty"`
	Fields          []string    `json:"fields,omitempty"`
	Filter          interface{} `json:"filter,omitempty"`
	RetrieveVectors bool        `json:"retrieveVectors,omitempty"`
	Ids             []string    `json:"ids,omitempty"`
	Sort            []string    `json:"sort,omitempty"`
}

// SimilarDocumentQuery is query parameters of similar documents
type SimilarDocumentQuery struct {
	Id                      interface{} `json:"id,omitempty"`
	Embedder                string      `json:"embedder"`
	AttributesToRetrieve    []string    `json:"attributesToRetrieve,omitempty"`
	Offset                  int64       `json:"offset,omitempty"`
	Limit                   int64       `json:"limit,omitempty"`
	Filter                  string      `json:"filter,omitempty"`
	ShowRankingScore        bool        `json:"showRankingScore,omitempty"`
	ShowRankingScoreDetails bool        `json:"showRankingScoreDetails,omitempty"`
	RankingScoreThreshold   float64     `json:"rankingScoreThreshold,omitempty"`
	RetrieveVectors         bool        `json:"retrieveVectors,omitempty"`
}

type SimilarDocumentResult struct {
	Hits               Hits   `json:"hits,omitempty"`
	ID                 string `json:"id,omitempty"`
	ProcessingTimeMS   int64  `json:"processingTimeMs,omitempty"`
	Limit              int64  `json:"limit,omitempty"`
	Offset             int64  `json:"offset,omitempty"`
	EstimatedTotalHits int64  `json:"estimatedTotalHits,omitempty"`
}

type CsvDocumentsQuery struct {
	PrimaryKey   string `json:"primaryKey,omitempty"`
	CsvDelimiter string `json:"csvDelimiter,omitempty"`
}

type DocumentsResult struct {
	Results Hits  `json:"results"`
	Limit   int64 `json:"limit"`
	Offset  int64 `json:"offset"`
	Total   int64 `json:"total"`
}

type UpdateDocumentByFunctionRequest struct {
	Filter   string                 `json:"filter,omitempty"`
	Function string                 `json:"function"`
	Context  map[string]interface{} `json:"context,omitempty"`
}

// ExperimentalFeaturesResult represents the experimental features result from the API.
type ExperimentalFeaturesBase struct {
	LogsRoute               *bool `json:"logsRoute,omitempty"`
	Metrics                 *bool `json:"metrics,omitempty"`
	EditDocumentsByFunction *bool `json:"editDocumentsByFunction,omitempty"`
	ContainsFilter          *bool `json:"containsFilter,omitempty"`
	Network                 *bool `json:"network,omitempty"`
	CompositeEmbedders      *bool `json:"compositeEmbedders,omitempty"`
	ChatCompletions         *bool `json:"chatCompletions,omitempty"`
	MultiModal              *bool `json:"multimodal,omitempty"`
}

type ExperimentalFeaturesResult struct {
	LogsRoute               bool `json:"logsRoute"`
	Metrics                 bool `json:"metrics"`
	EditDocumentsByFunction bool `json:"editDocumentsByFunction"`
	ContainsFilter          bool `json:"containsFilter"`
	Network                 bool `json:"network"`
	CompositeEmbedders      bool `json:"compositeEmbedders"`
	ChatCompletions         bool `json:"chatCompletions"`
	MultiModal              bool `json:"multimodal,omitempty"`
}

type SwapIndexesParams struct {
	Indexes []string `json:"indexes"`
	Rename  bool     `json:"rename"`
}

// Health is the request body for set meilisearch health
type Health struct {
	Status string `json:"status"`
}

// UpdateIndexRequest is the request body for update Index primary key and renaming IndexUid
type UpdateIndexRequestParams struct {
	PrimaryKey string `json:"primaryKey,omitempty"`
	UID        string `json:"uid,omitempty"`
}

func (s *SearchRequest) validate() {
	if s.Hybrid != nil && s.Hybrid.Embedder == "" {
		s.Hybrid.Embedder = "default"
	}
}

// JSONMarshal returns the JSON encoding of v.
type JSONMarshal func(v interface{}) ([]byte, error)

// JSONUnmarshal parses the JSON-encoded data and stores the result
// in the value pointed to by v. If v is nil or not a pointer,
// Unmarshal returns an InvalidUnmarshalError.
type JSONUnmarshal func(data []byte, v interface{}) error

// Batch gives information about the progress of batch of asynchronous operations.
type Batch struct {
	UID           int                    `json:"uid"`
	Progress      *BatchProgress         `json:"progress,omitempty"`
	Details       map[string]interface{} `json:"details,omitempty"`
	Stats         *BatchStats            `json:"stats,omitempty"`
	Duration      string                 `json:"duration,omitempty"`
	StartedAt     time.Time              `json:"startedAt,omitempty"`
	FinishedAt    time.Time              `json:"finishedAt,omitempty"`
	BatchStrategy string                 `json:"batchStrategy,omitempty"`
}

type BatchProgress struct {
	Steps      []*BatchProgressStep `json:"steps"`
	Percentage float64              `json:"percentage"`
}

type BatchProgressStep struct {
	CurrentStep string `json:"currentStep"`
	Finished    int    `json:"finished"`
	Total       int    `json:"total"`
}

type BatchStats struct {
	TotalNbTasks           int                               `json:"totalNbTasks"`
	Status                 map[string]int                    `json:"status"`
	Types                  map[string]int                    `json:"types"`
	IndexedUIDs            map[string]int                    `json:"indexUids"`
	ProgressTrace          map[string]string                 `json:"progressTrace"`
	WriteChannelCongestion *BatchStatsWriteChannelCongestion `json:"writeChannelCongestion"`
	InternalDatabaseSizes  *BatchStatsInternalDatabaseSize   `json:"internalDatabaseSizes"`
}

type BatchStatsWriteChannelCongestion struct {
	Attempts         int     `json:"attempts"`
	BlockingAttempts int     `json:"blocking_attempts"`
	BlockingRatio    float64 `json:"blocking_ratio"`
}

type BatchStatsInternalDatabaseSize struct {
	ExternalDocumentsIDs    string `json:"externalDocumentsIds"`
	WordDocIDs              string `json:"wordDocids"`
	WordPairProximityDocIDs string `json:"wordPairProximityDocids"`
	WordPositionDocIDs      string `json:"wordPositionDocids"`
	WordFidDocIDs           string `json:"wordFidDocids"`
	FieldIdWordCountDocIDs  string `json:"fieldIdWordCountDocids"`
	Documents               string `json:"documents"`
}

type BatchesResults struct {
	Results []*Batch `json:"results"`
	Total   int64    `json:"total"`
	Limit   int64    `json:"limit"`
	From    int64    `json:"from"`
	Next    int64    `json:"next"`
}

// BatchesQuery represents the query parameters for listing batches.
type BatchesQuery struct {
	UIDs             []int64
	BatchUIDs        []int64
	IndexUIDs        []string
	Statuses         []string
	Types            []string
	Limit            int64
	From             int64
	Reverse          bool
	BeforeEnqueuedAt time.Time
	BeforeStartedAt  time.Time
	BeforeFinishedAt time.Time
	AfterEnqueuedAt  time.Time
	AfterStartedAt   time.Time
	AfterFinishedAt  time.Time
}

type ExportParams struct {
	URL         string                        `json:"url,omitempty"`
	APIKey      string                        `json:"apiKey,omitempty"`
	PayloadSize string                        `json:"payloadSize,omitempty"`
	Indexes     map[string]IndexExportOptions `json:"indexes,omitempty"`
}

type IndexExportOptions struct {
	Filter           string `json:"filter,omitempty"`
	OverrideSettings bool   `json:"overrideSettings,omitempty"`
}

type ChatWorkspace struct {
	UID string `json:"uid"`
}

type ChatWorkspaceSettings struct {
	Source       ChatSource                    `json:"source"`
	OrgId        string                        `json:"orgId"`
	ProjectId    string                        `json:"projectId"`
	ApiVersion   string                        `json:"apiVersion"`
	DeploymentId string                        `json:"deploymentId"`
	BaseUrl      string                        `json:"baseUrl"`
	ApiKey       string                        `json:"apiKey,omitempty"`
	Prompts      *ChatWorkspaceSettingsPrompts `json:"prompts"`
}

type ChatWorkspaceSettingsPrompts struct {
	System              string `json:"system"`
	SearchDescription   string `json:"searchDescription"`
	SearchQParam        string `json:"searchQParam"`
	SearchFilterParam   string `json:"searchFilterParam"`
	SearchIndexUidParam string `json:"searchIndexUidParam"`
}

type ListChatWorkspace struct {
	Results []*ChatWorkspace `json:"results"`
	Offset  int64            `json:"offset"`
	Limit   int64            `json:"limit"`
	Total   int64            `json:"total"`
}

type ListChatWorkSpaceQuery struct {
	Limit  int64
	Offset int64
}

type ChatCompletionQuery struct {
	Model    string                   `json:"model"`
	Messages []*ChatCompletionMessage `json:"messages"`
	Stream   bool                     `json:"stream"`
}

type ChatCompletionMessage struct {
	Role    ChatRole `json:"role"`
	Content string   `json:"content"`
}

type ChatCompletionStreamChunk struct {
	ID                string                  `json:"id"`
	Object            *string                 `json:"object,omitempty"`
	Created           Timestampz              `json:"created,omitempty"`
	Model             string                  `json:"model,omitempty"`
	Choices           []*ChatCompletionChoice `json:"choices"`
	ServiceTier       *string                 `json:"service_tier,omitempty"`
	SystemFingerprint *string                 `json:"system_fingerprint,omitempty"`
	Usage             any                     `json:"usage,omitempty"`
}

type ChatCompletionChoice struct {
	Index        int64                      `json:"index"`
	Delta        *ChatCompletionChoiceDelta `json:"delta"`
	FinishReason *string                    `json:"finish_reason,omitempty"`
	Logprobs     any                        `json:"logprobs"`
}

type ChatCompletionChoiceDelta struct {
	Content      *string   `json:"content,omitempty"`
	Role         *ChatRole `json:"role,omitempty"`
	Refusal      *string   `json:"refusal,omitempty"`
	FunctionCall *string   `json:"function_call,omitempty"`
	ToolCalls    *string   `json:"tool_calls,omitempty"`
}

type Timestampz int64

func (t Timestampz) String() string {
	return time.Unix(int64(t), 0).UTC().Format(time.RFC3339)
}

func (t Timestampz) ToTime() time.Time {
	return time.Unix(int64(t), 0).UTC()
}
