package meilisearch

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"
)

//
// Internal types to Meilisearch
//

// Client is a structure that give you the power for interacting with an high-level api with Meilisearch.
type Client struct {
	config     ClientConfig
	httpClient *fasthttp.Client
}

// Index is the type that represent an index in Meilisearch
type Index struct {
	UID        string    `json:"uid"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	PrimaryKey string    `json:"primaryKey,omitempty"`
	client     *Client
}

// Return of multiple indexes is wrap in a IndexesResults
type IndexesResults struct {
	Results []Index `json:"results"`
	Offset  int64   `json:"offset"`
	Limit   int64   `json:"limit"`
	Total   int64   `json:"total"`
}

type IndexesQuery struct {
	Limit  int64
	Offset int64
}

// Settings is the type that represents the settings in Meilisearch
type Settings struct {
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
}

// TypoTolerance is the type that represents the typo tolerance setting in Meilisearch
type TypoTolerance struct {
	Enabled             bool                `json:"enabled,omitempty"`
	MinWordSizeForTypos MinWordSizeForTypos `json:"minWordSizeForTypos,omitempty"`
	DisableOnWords      []string            `json:"disableOnWords,omitempty"`
	DisableOnAttributes []string            `json:"disableOnAttributes,omitempty"`
}

// MinWordSizeForTypos is the type that represents the minWordSizeForTypos setting in the typo tolerance setting in Meilisearch
type MinWordSizeForTypos struct {
	OneTypo  int64 `json:"oneTypo,omitempty"`
	TwoTypos int64 `json:"twoTypos,omitempty"`
}

// Pagination is the type that represents the pagination setting in Meilisearch
type Pagination struct {
	MaxTotalHits int64 `json:"maxTotalHits"`
}

// Faceting is the type that represents the faceting setting in Meilisearch
type Faceting struct {
	MaxValuesPerFacet int64 `json:"maxValuesPerFacet"`
}

// Version is the type that represents the versions in Meilisearch
type Version struct {
	CommitSha  string `json:"commitSha"`
	CommitDate string `json:"commitDate"`
	PkgVersion string `json:"pkgVersion"`
}

// StatsIndex is the type that represent the stats of an index in Meilisearch
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

// TaskStatus is the status of a task.
type TaskStatus string

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

// TaskType is the type of a task
type TaskType string

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
}

// Return of multiple tasks is wrap in a TaskResult
type TaskResult struct {
	Results []Task `json:"results"`
	Limit   int64  `json:"limit"`
	From    int64  `json:"from"`
	Next    int64  `json:"next"`
	Total   int64  `json:"total"`
}

// Keys allow the user to connect to the Meilisearch instance
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

// This structure is used to send the exact ISO-8601 time format managed by Meilisearch
type KeyParsed struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	UID         string   `json:"uid,omitempty"`
	Actions     []string `json:"actions,omitempty"`
	Indexes     []string `json:"indexes,omitempty"`
	ExpiresAt   *string  `json:"expiresAt"`
}

// This structure is used to update a Key
type KeyUpdate struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Return of multiple keys is wrap in a KeysResults
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

// Information to create a tenant token
//
// ExpiresAt is a time.Time when the key will expire.
// Note that if an ExpiresAt value is included it should be in UTC time.
// ApiKey is the API key parent of the token.
type TenantTokenOptions struct {
	APIKey    string
	ExpiresAt time.Time
}

// Custom Claims structure to create a Tenant Token
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
	Offset                int64
	Limit                 int64
	AttributesToRetrieve  []string
	AttributesToSearchOn  []string
	AttributesToCrop      []string
	CropLength            int64
	CropMarker            string
	AttributesToHighlight []string
	HighlightPreTag       string
	HighlightPostTag      string
	MatchingStrategy      string
	Filter                interface{}
	ShowMatchesPosition   bool
	ShowRankingScore      bool
	Facets                []string
	PlaceholderSearch     bool
	Sort                  []string
	Vector                []float64
	HitsPerPage           int64
	Page                  int64
	IndexUID              string
	Query                 string
}

type MultiSearchRequest struct {
	Queries []SearchRequest `json:"queries"`
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

// Health is the request body for set Meilisearch health
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
