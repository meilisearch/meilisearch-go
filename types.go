package meilisearch

import "time"

// Unknown is unknown json type
type Unknown map[string]interface{}

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

// Settings is the type that represent the settings in MeiliSearch
type Settings struct {
	RankingRules         []string            `json:"rankingRules,omitempty"`
	DistinctAttribute    *string             `json:"distinctAttribute,omitempty"`
	SearchableAttributes []string            `json:"searchableAttributes,omitempty"`
	DisplayedAttributes  []string            `json:"displayedAttributes,omitempty"`
	StopWords            []string            `json:"stopWords,omitempty"`
	Synonyms             map[string][]string `json:"synonyms,omitempty"`
	AcceptNewFields      bool                `json:"acceptNewFields,omitempty"`
}

// Version is the type that represent the versions in MeiliSearch
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

// SystemInformation is the type that represent the information system in MeiliSearch
type SystemInformation struct {
	MemoryUsage    float64   `json:"memoryUsage"`
	ProcessorUsage []float64 `json:"processorUsage"`
	Global         struct {
		TotalMemory int64 `json:"totalMemory"`
		UsedMemory  int64 `json:"usedMemory"`
		UsedSwap    int64 `json:"usedSwap"`
		InputData   int64 `json:"inputData"`
		OutputData  int64 `json:"outputData"`
	} `json:"global"`
	Process struct {
		Memory int64 `json:"memory"`
		CPU    int64 `json:"cpu"`
	} `json:"process"`
}

// SystemInformationPretty is the type that represent the information system (human readable) in MeiliSearch
type SystemInformationPretty struct {
	MemoryUsage    string   `json:"memoryUsage"`
	ProcessorUsage []string `json:"processorUsage"`
	Global         struct {
		TotalMemory string `json:"totalMemory"`
		UsedMemory  string `json:"usedMemory"`
		UsedSwap    string `json:"usedSwap"`
		InputData   string `json:"inputData"`
		OutputData  string `json:"outputData"`
	} `json:"global"`
	Process struct {
		Memory string `json:"memory"`
		CPU    string `json:"cpu"`
	} `json:"process"`
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
	UpdateID    int64        `json:"updateID"`
	Type        Unknown      `json:"type"`
	Error       string       `json:"error"`
	EnqueuedAt  time.Time    `json:"enqueuedAt"`
	ProcessedAt time.Time    `json:"processedAt"`
}

// AsyncUpdateID is returned for asynchronous method
//
// Documentation: https://docs.meilisearch.com/guides/advanced_guides/asynchronous_updates.html
type AsyncUpdateID struct {
	UpdateID int64 `json:"updateID"`
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
	FacetFilters          []string
}

// SearchResponse is the response body for search method
type SearchResponse struct {
	Hits               []interface{} `json:"hits"`
	Offset             int64         `json:"offset"`
	Limit              int64         `json:"limit"`
	ProcessingTimeMs   int64         `json:"processingTimeMs"`
	Query              string        `json:"query"`
	FacetsDistribution interface{}   `json:"facetsDistribution"`
}

// ListDocumentsRequest is the request body for list documents method
type ListDocumentsRequest struct {
	Offset               int64    `json:"offset,omitempty"`
	Limit                int64    `json:"limit,omitempty"`
	AttributesToRetrieve []string `json:"attributesToRetrieve,omitempty"`
}
