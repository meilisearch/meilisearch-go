package meilisearch

import (
	"time"

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
)

// Task indicate information about a task is returned for asynchronous method
//
// Documentation: https://docs.meilisearch.com/learn/advanced/asynchronous_operations.html
type Task struct {
	Status     TaskStatus          `json:"status"`
	UID        int64               `json:"uid"`
	IndexUID   string              `json:"indexUid"`
	Type       string              `json:"type"`
	Error      meilisearchApiError `json:"error,omitempty"`
	Duration   string              `json:"duration,omitempty"`
	EnqueuedAt time.Time           `json:"enqueuedAt"`
	StartedAt  time.Time           `json:"startedAt,omitempty"`
	FinishedAt time.Time           `json:"finishedAt,omitempty"`
	Details    Details             `json:"details,omitempty"`
}

type Details struct {
	ReceivedDocuments    int                 `json:"receivedDocuments,omitempty"`
	IndexedDocuments     int                 `json:"indexedDocuments,omitempty"`
	DeletedDocuments     int                 `json:"deletedDocuments,omitempty"`
	PrimaryKey           string              `json:"primaryKey,omitempty"`
	RankingRules         []string            `json:"rankingRules,omitempty"`
	DistinctAttribute    *string             `json:"distinctAttribute,omitempty"`
	SearchableAttributes []string            `json:"searchableAttributes,omitempty"`
	DisplayedAttributes  []string            `json:"displayedAttributes,omitempty"`
	StopWords            []string            `json:"stopWords,omitempty"`
	Synonyms             map[string][]string `json:"synonyms,omitempty"`
	FilterableAttributes []string            `json:"filterableAttributes,omitempty"`
	SortableAttributes   []string            `json:"sortableAttributes,omitempty"`
}

type ResultTask struct {
	Results []Task `json:"results"`
}

// Keys allow the user to connect to the Meilisearch instance
//
// Documentation: https://docs.meilisearch.com/learn/advanced/security.html#protecting-a-meilisearch-instance
type Key struct {
	Description string    `json:"description"`
	Key         string    `json:"key,omitempty"`
	Actions     []string  `json:"actions,omitempty"`
	Indexes     []string  `json:"indexes,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

// This structure is used to send the exact ISO-8601 time format managed by Meilisearch
type KeyParsed struct {
	Description string    `json:"description"`
	Key         string    `json:"key,omitempty"`
	Actions     []string  `json:"actions,omitempty"`
	Indexes     []string  `json:"indexes,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
	ExpiresAt   *string   `json:"expiresAt"`
}

type ResultKey struct {
	Results []Key `json:"results"`
}

// DumpStatus is the status of a dump.
type DumpStatus string

const (
	// DumpStatusInProgress means the server is processing the dump
	DumpStatusInProgress DumpStatus = "in_progress"
	// DumpStatusFailed means the server failed to create a dump
	DumpStatusFailed DumpStatus = "failed"
	// DumpStatusDone means the server completed the dump
	DumpStatusDone DumpStatus = "done"
)

// Dump indicate information about an dump
//
// Documentation: https://docs.meilisearch.com/reference/api/dump.html
type Dump struct {
	UID        string     `json:"uid"`
	Status     DumpStatus `json:"status"`
	StartedAt  time.Time  `json:"startedAt"`
	FinishedAt time.Time  `json:"finishedAt"`
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
// Documentation: https://docs.meilisearch.com/reference/features/search_parameters.html
type SearchRequest struct {
	Offset                int64
	Limit                 int64
	AttributesToRetrieve  []string
	AttributesToCrop      []string
	CropLength            int64
	AttributesToHighlight []string
	Filter                interface{}
	Matches               bool
	FacetsDistribution    []string
	PlaceholderSearch     bool
	Sort                  []string
}

// SearchResponse is the response body for search method
type SearchResponse struct {
	Hits                  []interface{} `json:"hits"`
	NbHits                int64         `json:"nbHits"`
	Offset                int64         `json:"offset"`
	Limit                 int64         `json:"limit"`
	ExhaustiveNbHits      bool          `json:"exhaustiveNbHits"`
	ProcessingTimeMs      int64         `json:"processingTimeMs"`
	Query                 string        `json:"query"`
	FacetsDistribution    interface{}   `json:"facetsDistribution,omitempty"`
	ExhaustiveFacetsCount interface{}   `json:"exhaustiveFacetsCount,omitempty"`
}

// DocumentsRequest is the request body for list documents method
type DocumentsRequest struct {
	Offset               int64    `json:"offset,omitempty"`
	Limit                int64    `json:"limit,omitempty"`
	AttributesToRetrieve []string `json:"attributesToRetrieve,omitempty"`
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
