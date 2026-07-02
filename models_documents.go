package meilisearch

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
	ShowPerformanceDetails  bool        `json:"showPerformanceDetails,omitempty"`
	RankingScoreThreshold   float64     `json:"rankingScoreThreshold,omitempty"`
	RetrieveVectors         bool        `json:"retrieveVectors,omitempty"`
}

type SimilarDocumentResult struct {
	Hits               Hits           `json:"hits,omitempty"`
	ID                 string         `json:"id,omitempty"`
	ProcessingTimeMS   int64          `json:"processingTimeMs,omitempty"`
	Limit              int64          `json:"limit,omitempty"`
	Offset             int64          `json:"offset,omitempty"`
	EstimatedTotalHits int64          `json:"estimatedTotalHits,omitempty"`
	PerformanceDetails map[string]any `json:"performanceDetails,omitempty"`
}

// DocumentOptions is the options struct for adding or updating documents (JSON/NDJSON)
// and deleting documents.
type DocumentOptions struct {
	PrimaryKey   *string `json:"primaryKey,omitempty"`
	SkipCreation bool    `json:"skipCreation,omitempty"`
	// TaskCustomMetadata is the custom metadata to add to the task.
	// This string will be associated with the task and visible in the task details.
	// It is optional.
	TaskCustomMetadata string `json:"-"`
}

type CsvDocumentsQuery struct {
	PrimaryKey   string `json:"primaryKey,omitempty"`
	CsvDelimiter string `json:"csvDelimiter,omitempty"`
	SkipCreation bool   `json:"skipCreation,omitempty"`
	// TaskCustomMetadata is the custom metadata to add to the task.
	// This string will be associated with the task and visible in the task details.
	// It is optional.
	TaskCustomMetadata string `json:"-"`
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
	// TaskCustomMetadata is the custom metadata to add to the task.
	// This string will be associated with the task and visible in the task details.
	// It is optional.
	TaskCustomMetadata string `json:"-"`
}

type RenderTemplateParams struct {
	Template Template        `json:"template,omitempty"`
	Input    *TempelateInput `json:"input,omitempty"`
}

type TemplateKind string

const (
	DocumentTemplate       TemplateKind = "documentTemplate"
	ChatDocumentTemplate   TemplateKind = "chatDocumentTemplate"
	IndexingFragment       TemplateKind = "indexingFragment"
	SearchFragment         TemplateKind = "searchFragment"
	InlineDocumentTemplate TemplateKind = "inlineDocumentTemplate"
	InlineFragment         TemplateKind = "inlineFragment"
)

type Template struct {
	Kind                     TemplateKind `json:"kind"`
	IndexUID                 *string      `json:"indexUid"`
	Embedder                 *string      `json:"embedder"`
	Fragment                 *string      `json:"fragment"`
	Inline                   any          `json:"inline"`
	DocumentTemplateMaxBytes *int64       `json:"documentTemplateMaxBytes"`
}

type InputKind string

const (
	IndexDocument  InputKind = "indexDocument"
	InlineDocument InputKind = "inlineDocument"
	InlineSearch   InputKind = "inlineSearch"
)

type TempelateInput struct {
	Kind     InputKind `json:"kind"`
	IndexUID *string   `json:"indexUid"`
	ID       *string   `json:"id"`
	Inline   any       `json:"inline"`
}

type RenderTemplateResponse struct {
	Template any `json:"template"`
	Rendered any `json:"rendered"`
}
