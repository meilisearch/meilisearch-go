package meilisearch

import "time"

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

// CreateIndexRequest is the request body for create index method
type CreateIndexRequest struct {
	UID        string `json:"uid,omitempty"`
	PrimaryKey string `json:"primaryKey,omitempty"`
}

type SwapIndexesParams struct {
	Indexes []string `json:"indexes"`
	Rename  bool     `json:"rename"`
}

// UpdateIndexRequestParams is the request body for update Index primary key and renaming IndexUid
type UpdateIndexRequestParams struct {
	PrimaryKey string `json:"primaryKey,omitempty"`
	UID        string `json:"uid,omitempty"`
}

type StatsParams struct {
	ShowInternalDatabaseSizes bool   `json:"showInternalDatabaseSizes,omitempty"`
	SizeFormat                string `json:"sizeFormat,omitempty"` // human, raw
}

// StatsIndex is the type that represent the stats of an index in meilisearch
type StatsIndex struct {
	NumberOfDocuments         int64            `json:"numberOfDocuments"`
	IsIndexing                bool             `json:"isIndexing"`
	FieldDistribution         map[string]int64 `json:"fieldDistribution"`
	RawDocumentDbSize         any              `json:"rawDocumentDbSize"`
	AvgDocumentSize           any              `json:"avgDocumentSize"`
	NumberOfEmbeddedDocuments int64            `json:"numberOfEmbeddedDocuments"`
	NumberOfEmbeddings        int64            `json:"numberOfEmbeddings"`
	InternalDatabaseSizes     map[string]any   `json:"internalDatabaseSizes"`
}

// Stats is the type that represent all stats
type Stats struct {
	DatabaseSize     any                   `json:"databaseSize"`
	UsedDatabaseSize any                   `json:"usedDatabaseSize"`
	LastUpdate       time.Time             `json:"lastUpdate"`
	Indexes          map[string]StatsIndex `json:"indexes"`
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
