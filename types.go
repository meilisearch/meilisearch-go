package meilisearch

import (
	"bytes"
	"github.com/valyala/fastjson"
	"sync"
	"time"
)

var arp fastjson.ArenaPool

var (
	bf bytes.Buffer
	mu sync.Mutex
)

//
// Internal types to Meilisearch
//

// Index is the type that represent an index in MeiliSearch
type Index struct {
	Name       string    `json:"name"`
	UID        string    `json:"uid"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	PrimaryKey string    `json:"primaryKey,omitempty"`
}

// Settings is the type that represents the settings in MeiliSearch
type Settings struct {
	RankingRules          []string            `json:"rankingRules,omitempty"`
	DistinctAttribute     *string             `json:"distinctAttribute,omitempty"`
	SearchableAttributes  []string            `json:"searchableAttributes,omitempty"`
	DisplayedAttributes   []string            `json:"displayedAttributes,omitempty"`
	StopWords             []string            `json:"stopWords,omitempty"`
	Synonyms              map[string][]string `json:"synonyms,omitempty"`
	AttributesForFaceting []string            `json:"attributesForFaceting,omitempty"`
}

// Version is the type that represents the versions in MeiliSearch
type Version struct {
	CommitSha  string    `json:"commitSha"`
	BuildDate  time.Time `json:"buildDate"`
	PkgVersion string    `json:"pkgVersion"`
}

// StatsIndex is the type that represent the stats of an index in MeiliSearch
type StatsIndex struct {
	NumberOfDocuments int64            `json:"numberOfDocuments"`
	IsIndexing        bool             `json:"isIndexing"`
	FieldsFrequency   map[string]int64 `json:"fieldsFrequency"`
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
// Documentation: https://docs.meilisearch.com/guides/advanced_guides/asynchronous_updates.html
type AsyncUpdateID struct {
	UpdateID int64 `json:"updateId"`
}

// Keys allow the user to connect to the MeiliSearch instance
//
// Documentation: https://docs.meilisearch.com/guides/advanced_guides/asynchronous_updates.html
type Keys struct {
	Public  string `json:"public,omitempty"`
	Private string `json:"private,omitempty"`
}

//
// Request/Response
//

// CreateIndexRequest is the request body for create index method
type CreateIndexRequest struct {
	Name       string `json:"name,omitempty"`
	UID        string `json:"uid,omitempty"`
	PrimaryKey string `json:"primaryKey,omitempty"`
}

// CreateIndexResponse is the response body for create index method
type CreateIndexResponse struct {
	Name       string    `json:"name"`
	UID        string    `json:"uid"`
	UpdateID   int64     `json:"updateID,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	PrimaryKey string    `json:"primaryKey,omitempty"`
}

// SearchRequest is the request url param needed for a search query.
// This struct will be converted to url param before sent.
//
// Documentation: https://docs.meilisearch.com/guides/advanced_guides/search_parameters.html
type SearchRequest struct {
	Query                 string
	Offset                int64
	Limit                 int64
	AttributesToRetrieve  []string
	AttributesToCrop      []string
	CropLength            int64
	AttributesToHighlight []string
	Filters               string
	Matches               bool
	FacetsDistribution    []string
	FacetFilters          interface{}
	PlaceholderSearch     bool
}

// SearchResponse is the response body for search method
type SearchResponse struct {
	Hits                  []interface{} `json:"hits"`
	NbHits                int64         `json:"nbHits"`
	Offset                int64         `json:"offset"`
	Limit                 int64         `json:"limit"`
	ProcessingTimeMs      int64         `json:"processingTimeMs"`
	Query                 string        `json:"query"`
	FacetsDistribution    interface{}   `json:"facetsDistribution,omitempty"`
	ExhaustiveFacetsCount interface{}   `json:"exhaustiveFacetsCount,omitempty"`
}

// ListDocumentsRequest is the request body for list documents method
type ListDocumentsRequest struct {
	Offset               int64    `json:"offset,omitempty"`
	Limit                int64    `json:"limit,omitempty"`
	AttributesToRetrieve []string `json:"attributesToRetrieve,omitempty"`
}

// RawType is an alias for raw byte[]
type RawType []byte

// Health is the request body for set Meilisearch health
type Health struct {
	Health bool `json:"health"`
}

// Name is the request body for set Index name
type Name struct {
	Name string `json:"name"`
}

// PrimaryKey is the request body for set Index primary key
type PrimaryKey struct {
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
