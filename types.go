package meilisearch

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	DefaultLimit int64 = 20

	contentTypeJSON   string = "application/json"
	contentTypeNDJSON string = "application/x-ndjson"
	contentTypeCSV    string = "text/csv"
)

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
	TypoTolerance        *TypoTolerance         `json:"typoTolerance,omitempty"`
	Pagination           *Pagination            `json:"pagination,omitempty"`
	Faceting             *Faceting              `json:"faceting,omitempty"`
	Embedders            map[string]Embedder    `json:"embedders,omitempty"`
}

// TypoTolerance is the type that represents the typo tolerance setting in meilisearch
type TypoTolerance struct {
	Enabled             bool                `json:"enabled"`
	MinWordSizeForTypos MinWordSizeForTypos `json:"minWordSizeForTypos,omitempty"`
	DisableOnWords      []string            `json:"disableOnWords,omitempty"`
	DisableOnAttributes []string            `json:"disableOnAttributes,omitempty"`
}

// MinWordSizeForTypos is the type that represents the minWordSizeForTypos setting in the typo tolerance setting in meilisearch
type MinWordSizeForTypos struct {
	OneTypo  int64 `json:"oneTypo,omitempty"`
	TwoTypos int64 `json:"twoTypos,omitempty"`
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

type Embedder struct {
	Source           string `json:"source"`
	ApiKey           string `json:"apiKey,omitempty"`
	Model            string `json:"model,omitempty"`
	Dimensions       int    `json:"dimensions,omitempty"`
	DocumentTemplate string `json:"documentTemplate,omitempty"`
}

// Version is the type that represents the versions in meilisearch
type Version struct {
	CommitSha  string `json:"commitSha"`
	CommitDate string `json:"commitDate"`
	PkgVersion string `json:"pkgVersion"`
}

// StatsIndex is the type that represent the stats of an index in meilisearch
type StatsIndex struct {
	NumberOfDocuments int64            `json:"numberOfDocuments"`
	IsIndexing        bool             `json:"isIndexing"`
	FieldDistribution map[string]int64 `json:"fieldDistribution"`
}

// Stats is the type that represent all stats
type Stats struct {
	DatabaseSize int64                 `json:"databaseSize"`
	LastUpdate   time.Time             `json:"lastUpdate"`
	Indexes      map[string]StatsIndex `json:"indexes"`
}

type (
	TaskType               string // TaskType is the type of a task
	SortFacetType          string // SortFacetType is type of facet sorting, alpha or count
	TaskStatus             string // TaskStatus is the status of a task.
	ProximityPrecisionType string // ProximityPrecisionType accepts one of the ByWord or ByAttribute
	MatchingStrategy       string // MatchingStrategy one of the Last, All, Frequency
)

const (
	// Last returns documents containing all the query terms first. If there are not enough results containing all
	// query terms to meet the requested limit, Meilisearch will remove one query term at a time,
	// starting from the end of the query.
	Last MatchingStrategy = "last"
	// All only returns documents that contain all query terms. Meilisearch will not match any more documents even
	// if there aren't enough to meet the requested limit.
	All MatchingStrategy = "all"
	// Frequency returns documents containing all the query terms first. If there are not enough results containing
	//all query terms to meet the requested limit, Meilisearch will remove one query term at a time, starting
	//with the word that is the most frequent in the dataset. frequency effectively gives more weight to terms
	//that appear less frequently in a set of results.
	Frequency MatchingStrategy = "frequency"
)

const (
	// ByWord calculate the precise distance between query terms. Higher precision, but may lead to longer
	// indexing time. This is the default setting
	ByWord ProximityPrecisionType = "byWord"
	// ByAttribute determine if multiple query terms are present in the same attribute.
	// Lower precision, but shorter indexing time
	ByAttribute ProximityPrecisionType = "byAttribute"
)

const (
	// TaskStatusUnknown is the default TaskStatus, should not exist
	TaskStatusUnknown TaskStatus = "unknown"
	// TaskStatusEnqueued the task request has been received and will be processed soon
	TaskStatusEnqueued TaskStatus = "enqueued"
	// TaskStatusProcessing the task is being processed
	TaskStatusProcessing TaskStatus = "processing"
	// TaskStatusSucceeded the task has been successfully processed
	TaskStatusSucceeded TaskStatus = "succeeded"
	// TaskStatusFailed a failure occurred when processing the task, no changes were made to the database
	TaskStatusFailed TaskStatus = "failed"
	// TaskStatusCanceled the task was canceled
	TaskStatusCanceled TaskStatus = "canceled"
)

const (
	SortFacetTypeAlpha SortFacetType = "alpha"
	SortFacetTypeCount SortFacetType = "count"
)

const (
	// TaskTypeIndexCreation represents an index creation
	TaskTypeIndexCreation TaskType = "indexCreation"
	// TaskTypeIndexUpdate represents an index update
	TaskTypeIndexUpdate TaskType = "indexUpdate"
	// TaskTypeIndexDeletion represents an index deletion
	TaskTypeIndexDeletion TaskType = "indexDeletion"
	// TaskTypeIndexSwap represents an index swap
	TaskTypeIndexSwap TaskType = "indexSwap"
	// TaskTypeDocumentAdditionOrUpdate represents a document addition or update in an index
	TaskTypeDocumentAdditionOrUpdate TaskType = "documentAdditionOrUpdate"
	// TaskTypeDocumentDeletion represents a document deletion from an index
	TaskTypeDocumentDeletion TaskType = "documentDeletion"
	// TaskTypeSettingsUpdate represents a settings update
	TaskTypeSettingsUpdate TaskType = "settingsUpdate"
	// TaskTypeDumpCreation represents a dump creation
	TaskTypeDumpCreation TaskType = "dumpCreation"
	// TaskTypeTaskCancelation represents a task cancelation
	TaskTypeTaskCancelation TaskType = "taskCancelation"
	// TaskTypeTaskDeletion represents a task deletion
	TaskTypeTaskDeletion TaskType = "taskDeletion"
	// TaskTypeSnapshotCreation represents a snapshot creation
	TaskTypeSnapshotCreation TaskType = "snapshotCreation"
)

// Task indicates information about a task resource
//
// Documentation: https://www.meilisearch.com/docs/learn/advanced/asynchronous_operations
type Task struct {
	Status     TaskStatus          `json:"status"`
	UID        int64               `json:"uid,omitempty"`
	TaskUID    int64               `json:"taskUid,omitempty"`
	IndexUID   string              `json:"indexUid"`
	Type       TaskType            `json:"type"`
	Error      meilisearchApiError `json:"error,omitempty"`
	Duration   string              `json:"duration,omitempty"`
	EnqueuedAt time.Time           `json:"enqueuedAt"`
	StartedAt  time.Time           `json:"startedAt,omitempty"`
	FinishedAt time.Time           `json:"finishedAt,omitempty"`
	Details    Details             `json:"details,omitempty"`
	CanceledBy int64               `json:"canceledBy,omitempty"`
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
	FilterableAttributes []string            `json:"filterableAttributes,omitempty"`
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
	Offset                  int64                `json:"offset,omitempty"`
	Limit                   int64                `json:"limit,omitempty"`
	AttributesToRetrieve    []string             `json:"attributesToRetrieve,omitempty"`
	AttributesToSearchOn    []string             `json:"attributesToSearchOn,omitempty"`
	AttributesToCrop        []string             `json:"attributesToCrop,omitempty"`
	CropLength              int64                `json:"cropLength,omitempty"`
	CropMarker              string               `json:"cropMarker,omitempty"`
	AttributesToHighlight   []string             `json:"attributesToHighlight,omitempty"`
	HighlightPreTag         string               `json:"highlightPreTag,omitempty"`
	HighlightPostTag        string               `json:"highlightPostTag,omitempty"`
	MatchingStrategy        MatchingStrategy     `json:"matchingStrategy,omitempty"`
	Filter                  interface{}          `json:"filter,omitempty"`
	ShowMatchesPosition     bool                 `json:"showMatchesPosition,omitempty"`
	ShowRankingScore        bool                 `json:"showRankingScore,omitempty"`
	ShowRankingScoreDetails bool                 `json:"showRankingScoreDetails,omitempty"`
	Facets                  []string             `json:"facets,omitempty"`
	Sort                    []string             `json:"sort,omitempty"`
	Vector                  []float32            `json:"vector,omitempty"`
	HitsPerPage             int64                `json:"hitsPerPage,omitempty"`
	Page                    int64                `json:"page,omitempty"`
	IndexUID                string               `json:"indexUid,omitempty"`
	Query                   string               `json:"q"`
	Distinct                string               `json:"distinct,omitempty"`
	Hybrid                  *SearchRequestHybrid `json:"hybrid,omitempty"`
	RetrieveVectors         bool                 `json:"retrieveVectors,omitempty"`
	RankingScoreThreshold   float64              `json:"rankingScoreThreshold,omitempty"`
}

type SearchRequestHybrid struct {
	SemanticRatio float64 `json:"semanticRatio,omitempty"`
	Embedder      string  `json:"embedder,omitempty"`
}

type MultiSearchRequest struct {
	Queries []*SearchRequest `json:"queries"`
}

// SearchResponse is the response body for search method
type SearchResponse struct {
	Hits               []interface{} `json:"hits"`
	EstimatedTotalHits int64         `json:"estimatedTotalHits,omitempty"`
	Offset             int64         `json:"offset,omitempty"`
	Limit              int64         `json:"limit,omitempty"`
	ProcessingTimeMs   int64         `json:"processingTimeMs"`
	Query              string        `json:"query"`
	FacetDistribution  interface{}   `json:"facetDistribution,omitempty"`
	TotalHits          int64         `json:"totalHits,omitempty"`
	HitsPerPage        int64         `json:"hitsPerPage,omitempty"`
	Page               int64         `json:"page,omitempty"`
	TotalPages         int64         `json:"totalPages,omitempty"`
	FacetStats         interface{}   `json:"facetStats,omitempty"`
	IndexUID           string        `json:"indexUid,omitempty"`
}

type MultiSearchResponse struct {
	Results []SearchResponse `json:"results"`
}

type FacetSearchRequest struct {
	FacetName            string   `json:"facetName,omitempty"`
	FacetQuery           string   `json:"facetQuery,omitempty"`
	Q                    string   `json:"q,omitempty"`
	Filter               string   `json:"filter,omitempty"`
	MatchingStrategy     string   `json:"matchingStrategy,omitempty"`
	AttributesToSearchOn []string `json:"attributesToSearchOn,omitempty"`
}

type FacetSearchResponse struct {
	FacetHits        []interface{} `json:"facetHits"`
	FacetQuery       string        `json:"facetQuery"`
	ProcessingTimeMs int64         `json:"processingTimeMs"`
}

// DocumentQuery is the request body get one documents method
type DocumentQuery struct {
	Fields []string `json:"fields,omitempty"`
}

// DocumentsQuery is the request body for list documents method
type DocumentsQuery struct {
	Offset int64       `json:"offset,omitempty"`
	Limit  int64       `json:"limit,omitempty"`
	Fields []string    `json:"fields,omitempty"`
	Filter interface{} `json:"filter,omitempty"`
}

// SimilarDocumentQuery is query parameters of similar documents
type SimilarDocumentQuery struct {
	Id                      interface{} `json:"id,omitempty"`
	Embedder                string      `json:"embedder,omitempty"`
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
	Hits               []interface{} `json:"hits,omitempty"`
	ID                 string        `json:"id,omitempty"`
	ProcessingTimeMS   int64         `json:"processingTimeMs,omitempty"`
	Limit              int64         `json:"limit,omitempty"`
	Offset             int64         `json:"offset,omitempty"`
	EstimatedTotalHits int64         `json:"estimatedTotalHits,omitempty"`
}

type CsvDocumentsQuery struct {
	PrimaryKey   string `json:"primaryKey,omitempty"`
	CsvDelimiter string `json:"csvDelimiter,omitempty"`
}

type DocumentsResult struct {
	Results []map[string]interface{} `json:"results"`
	Limit   int64                    `json:"limit"`
	Offset  int64                    `json:"offset"`
	Total   int64                    `json:"total"`
}

type SwapIndexesParams struct {
	Indexes []string `json:"indexes"`
}

// RawType is an alias for raw byte[]
type RawType []byte

// Health is the request body for set meilisearch health
type Health struct {
	Status string `json:"status"`
}

// UpdateIndexRequest is the request body for update Index primary key
type UpdateIndexRequest struct {
	PrimaryKey string `json:"primaryKey"`
}

// Unknown is unknown json type
type Unknown map[string]interface{}

// UnmarshalJSON supports json.Unmarshaler interface
func (b *RawType) UnmarshalJSON(data []byte) error {
	*b = data
	return nil
}

// MarshalJSON supports json.Marshaler interface
func (b RawType) MarshalJSON() ([]byte, error) {
	return b, nil
}

func (s *SearchRequest) validate() {
	if s.Limit == 0 {
		s.Limit = DefaultLimit
	}
	if s.Hybrid != nil && s.Hybrid.Embedder == "" {
		s.Hybrid.Embedder = "default"
	}
}
