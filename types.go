package meilisearch

import (
	"bytes"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
)

var arp fastjson.ArenaPool

var (
	bf bytes.Buffer
	mu sync.Mutex
)

//
// Internal types to Meilisearch
//

// Client is a structure that give you the power for interacting with an high-level api with meilisearch.
type Client struct {
	config     ClientConfig
	httpClient *fasthttp.Client
}

// Index is the type that represent an index in MeiliSearch
type Index struct {
	UID        string    `json:"uid"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	PrimaryKey string    `json:"primaryKey,omitempty"`
	client     *Client
}

// Settings is the type that represents the settings in MeiliSearch
type Settings struct {
	RankingRules         []string            `json:"rankingRules,omitempty"`
	DistinctAttribute    *string             `json:"distinctAttribute,omitempty"`
	SearchableAttributes []string            `json:"searchableAttributes,omitempty"`
	DisplayedAttributes  []string            `json:"displayedAttributes,omitempty"`
	StopWords            []string            `json:"stopWords,omitempty"`
	Synonyms             map[string][]string `json:"synonyms,omitempty"`
	FilterableAttributes []string            `json:"filterableAttributes,omitempty"`
}

// Version is the type that represents the versions in MeiliSearch
type Version struct {
	CommitSha  string `json:"commitSha"`
	CommitDate string `json:"commitDate"`
	PkgVersion string `json:"pkgVersion"`
}

// StatsIndex is the type that represent the stats of an index in MeiliSearch
type StatsIndex struct {
	NumberOfDocuments int64            `json:"numberOfDocuments"`
	IsIndexing        bool             `json:"isIndexing"`
	FieldDistribution map[string]int64 `json:"fieldDistribution"`
}

// Stats is the type that represent all stats
type Stats struct {
	DatabaseSize int64                 `json:"database_size"`
	LastUpdate   time.Time             `json:"last_update"`
	Indexes      map[string]StatsIndex `json:"indexes"`
}

// UpdateStatus is the status of an update.
type UpdateStatus string

const (
	// UpdateStatusUnknown is the default UpdateStatus, should not exist
	UpdateStatusUnknown UpdateStatus = "unknown"
	// UpdateStatusEnqueued means the server know the update but didn't handle it yet
	UpdateStatusEnqueued UpdateStatus = "enqueued"
	// UpdateStatusProcessing means the server is processing the update and all went well
	UpdateStatusProcessing UpdateStatus = "processing"
	// UpdateStatusProcessed means the server has processed the update and all went well
	UpdateStatusProcessed UpdateStatus = "processed"
	// UpdateStatusFailed means the server has processed the update and an error has been reported
	UpdateStatusFailed UpdateStatus = "failed"
)

// Update indicate information about an update
type Update struct {
	Status      UpdateStatus `json:"status"`
	UpdateID    int64        `json:"updateId"`
	Type        Unknown      `json:"type"`
	Error       string       `json:"error"`
	EnqueuedAt  time.Time    `json:"enqueuedAt"`
	ProcessedAt time.Time    `json:"processedAt"`
}

// AsyncUpdateID is returned for asynchronous method
//
// Documentation: https://docs.meilisearch.com/learn/advanced/asynchronous_updates.html
type AsyncUpdateID struct {
	UpdateID int64 `json:"updateId"`
}

// Keys allow the user to connect to the MeiliSearch instance
//
// Documentation: https://docs.meilisearch.com/learn/advanced/asynchronous_updates.html
type Keys struct {
	Public  string `json:"public,omitempty"`
	Private string `json:"private,omitempty"`
}

// Dump indicate information about an dump
//
// Documentation: https://docs.meilisearch.com/reference/api/dump.html
type Dump struct {
	UID        string    `json:"uid"`
	Status     string    `json:"status"`
	StartedAt  time.Time `json:"startedAt"`
	FinishedAt time.Time `json:"finishedAt"`
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
