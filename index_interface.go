package meilisearch

import (
	"context"
	"encoding/json"
	"io"
)

type IndexManager interface {
	IndexReader
	TaskReader
	DocumentManager
	SettingsManager
	SearchReader

	GetIndexReader() IndexReader
	GetTaskReader() TaskReader
	GetDocumentManager() DocumentManager
	GetDocumentReader() DocumentReader
	GetSettingsManager() SettingsManager
	GetSettingsReader() SettingsReader
	GetSearch() SearchReader

	// UpdateIndex updates the primary key of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/update-index
	UpdateIndex(params *UpdateIndexRequestParams) (*TaskInfo, error)

	// UpdateIndexWithContext updates the primary key of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/update-index
	UpdateIndexWithContext(ctx context.Context, params *UpdateIndexRequestParams) (*TaskInfo, error)

	// Delete removes the index identified by the given UID.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/delete-index
	Delete(uid string) (bool, error)

	// DeleteWithContext removes the index identified by the given UID using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/delete-index
	DeleteWithContext(ctx context.Context, uid string) (bool, error)

	// Compact compacts the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/compact-index
	Compact() (*TaskInfo, error)

	// CompactWithContext compacts the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/compact-index
	CompactWithContext(ctx context.Context) (*TaskInfo, error)
}

type IndexReader interface {
	// FetchInfo retrieves information about the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/get-index
	FetchInfo() (*IndexResult, error)

	// FetchInfoWithContext retrieves information about the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/get-index
	FetchInfoWithContext(ctx context.Context) (*IndexResult, error)

	// FetchPrimaryKey retrieves the primary key of the index.
	FetchPrimaryKey() (*string, error)

	// FetchPrimaryKeyWithContext retrieves the primary key of the index using the provided context for cancellation.
	FetchPrimaryKeyWithContext(ctx context.Context) (*string, error)

	// GetStats retrieves statistical information about the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/stats/get-stats-of-index
	GetStats() (*StatsIndex, error)

	// GetStatsWithContext retrieves statistical information about the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/stats/get-stats-of-index
	GetStatsWithContext(ctx context.Context) (*StatsIndex, error)
}

type DocumentManager interface {
	DocumentReader

	// AddDocuments adds multiple documents to the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	AddDocuments(documentsPtr interface{}, opts *DocumentOptions) (*TaskInfo, error)

	// AddDocumentsWithContext adds multiple documents to the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	AddDocumentsWithContext(ctx context.Context, documentsPtr interface{}, opts *DocumentOptions) (*TaskInfo, error)

	// AddDocumentsInBatches adds documents to the index in batches of specified size.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	AddDocumentsInBatches(documentsPtr interface{}, batchSize int, opts *DocumentOptions) ([]TaskInfo, error)

	// AddDocumentsInBatchesWithContext adds documents to the index in batches of specified size using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	AddDocumentsInBatchesWithContext(ctx context.Context, documentsPtr interface{}, batchSize int, opts *DocumentOptions) ([]TaskInfo, error)

	// AddDocumentsCsv adds documents from a CSV byte array to the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	AddDocumentsCsv(documents []byte, options *CsvDocumentsQuery) (*TaskInfo, error)

	// AddDocumentsCsvWithContext adds documents from a CSV byte array to the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	AddDocumentsCsvWithContext(ctx context.Context, documents []byte, options *CsvDocumentsQuery) (*TaskInfo, error)

	// AddDocumentsCsvInBatches adds documents from a CSV byte array to the index in batches of specified size.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	AddDocumentsCsvInBatches(documents []byte, batchSize int, options *CsvDocumentsQuery) ([]TaskInfo, error)

	// AddDocumentsCsvInBatchesWithContext adds documents from a CSV byte array to the index in batches of specified size using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	AddDocumentsCsvInBatchesWithContext(ctx context.Context, documents []byte, batchSize int, options *CsvDocumentsQuery) ([]TaskInfo, error)

	// AddDocumentsCsvFromReaderInBatches adds documents from a CSV reader to the index in batches of specified size.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	AddDocumentsCsvFromReaderInBatches(documents io.Reader, batchSize int, options *CsvDocumentsQuery) ([]TaskInfo, error)

	// AddDocumentsCsvFromReaderInBatchesWithContext adds documents from a CSV reader to the index in batches of specified size using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	AddDocumentsCsvFromReaderInBatchesWithContext(ctx context.Context, documents io.Reader, batchSize int, options *CsvDocumentsQuery) ([]TaskInfo, error)

	// AddDocumentsCsvFromReader adds documents from a CSV reader to the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	AddDocumentsCsvFromReader(documents io.Reader, options *CsvDocumentsQuery) (*TaskInfo, error)

	// AddDocumentsCsvFromReaderWithContext adds documents from a CSV reader to the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	AddDocumentsCsvFromReaderWithContext(ctx context.Context, documents io.Reader, options *CsvDocumentsQuery) (*TaskInfo, error)

	// AddDocumentsNdjson adds documents from a NDJSON byte array to the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	AddDocumentsNdjson(documents []byte, opts *DocumentOptions) (*TaskInfo, error)

	// AddDocumentsNdjsonWithContext adds documents from a NDJSON byte array to the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	AddDocumentsNdjsonWithContext(ctx context.Context, documents []byte, opts *DocumentOptions) (*TaskInfo, error)

	// AddDocumentsNdjsonInBatches adds documents from a NDJSON byte array to the index in batches of specified size.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	AddDocumentsNdjsonInBatches(documents []byte, batchSize int, opts *DocumentOptions) ([]TaskInfo, error)

	// AddDocumentsNdjsonInBatchesWithContext adds documents from a NDJSON byte array to the index in batches of specified size using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	AddDocumentsNdjsonInBatchesWithContext(ctx context.Context, documents []byte, batchSize int, opts *DocumentOptions) ([]TaskInfo, error)

	// AddDocumentsNdjsonFromReader adds documents from a NDJSON reader to the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	AddDocumentsNdjsonFromReader(documents io.Reader, opts *DocumentOptions) (*TaskInfo, error)

	// AddDocumentsNdjsonFromReaderWithContext adds documents from a NDJSON reader to the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	AddDocumentsNdjsonFromReaderWithContext(ctx context.Context, documents io.Reader, opts *DocumentOptions) (*TaskInfo, error)

	// AddDocumentsNdjsonFromReaderInBatches adds documents from a NDJSON reader to the index in batches of specified size.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	AddDocumentsNdjsonFromReaderInBatches(documents io.Reader, batchSize int, opts *DocumentOptions) ([]TaskInfo, error)

	// AddDocumentsNdjsonFromReaderInBatchesWithContext adds documents from a NDJSON reader to the index in batches of specified size using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	AddDocumentsNdjsonFromReaderInBatchesWithContext(ctx context.Context, documents io.Reader, batchSize int, opts *DocumentOptions) ([]TaskInfo, error)

	// UpdateDocuments updates multiple documents in the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	UpdateDocuments(documentsPtr interface{}, opts *DocumentOptions) (*TaskInfo, error)

	// UpdateDocumentsWithContext updates multiple documents in the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	UpdateDocumentsWithContext(ctx context.Context, documentsPtr interface{}, opts *DocumentOptions) (*TaskInfo, error)

	// UpdateDocumentsInBatches updates documents in the index in batches of specified size.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	UpdateDocumentsInBatches(documentsPtr interface{}, batchSize int, opts *DocumentOptions) ([]TaskInfo, error)

	// UpdateDocumentsInBatchesWithContext updates documents in the index in batches of specified size using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	UpdateDocumentsInBatchesWithContext(ctx context.Context, documentsPtr interface{}, batchSize int, opts *DocumentOptions) ([]TaskInfo, error)

	// UpdateDocumentsCsv updates documents in the index from a CSV byte array.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	UpdateDocumentsCsv(documents []byte, options *CsvDocumentsQuery) (*TaskInfo, error)

	// UpdateDocumentsCsvWithContext updates documents in the index from a CSV byte array using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	UpdateDocumentsCsvWithContext(ctx context.Context, documents []byte, options *CsvDocumentsQuery) (*TaskInfo, error)

	// UpdateDocumentsCsvInBatches updates documents in the index from a CSV byte array in batches of specified size.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	UpdateDocumentsCsvInBatches(documents []byte, batchsize int, options *CsvDocumentsQuery) ([]TaskInfo, error)

	// UpdateDocumentsCsvInBatchesWithContext updates documents in the index from a CSV byte array in batches of specified size using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	UpdateDocumentsCsvInBatchesWithContext(ctx context.Context, documents []byte, batchsize int, options *CsvDocumentsQuery) ([]TaskInfo, error)

	// UpdateDocumentsNdjson updates documents in the index from a NDJSON byte array.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	UpdateDocumentsNdjson(documents []byte, opts *DocumentOptions) (*TaskInfo, error)

	// UpdateDocumentsNdjsonWithContext updates documents in the index from a NDJSON byte array using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	UpdateDocumentsNdjsonWithContext(ctx context.Context, documents []byte, opts *DocumentOptions) (*TaskInfo, error)

	// UpdateDocumentsNdjsonInBatches updates documents in the index from a NDJSON byte array in batches of specified size.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	UpdateDocumentsNdjsonInBatches(documents []byte, batchsize int, opts *DocumentOptions) ([]TaskInfo, error)

	// UpdateDocumentsNdjsonInBatchesWithContext updates documents in the index from a NDJSON byte array in batches of specified size using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	UpdateDocumentsNdjsonInBatchesWithContext(ctx context.Context, documents []byte, batchsize int, opts *DocumentOptions) ([]TaskInfo, error)

	// UpdateDocumentsByFunction update documents by using function
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	UpdateDocumentsByFunction(req *UpdateDocumentByFunctionRequest) (*TaskInfo, error)

	// UpdateDocumentsByFunctionWithContext update documents by using function then provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-update-documents
	UpdateDocumentsByFunctionWithContext(ctx context.Context, req *UpdateDocumentByFunctionRequest) (*TaskInfo, error)

	// DeleteDocument deletes a single document from the index by identifier.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/delete-document
	DeleteDocument(identifier string, opts *DocumentOptions) (*TaskInfo, error)

	// DeleteDocumentWithContext deletes a single document from the index by identifier using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/delete-document
	DeleteDocumentWithContext(ctx context.Context, identifier string, opts *DocumentOptions) (*TaskInfo, error)

	// DeleteDocuments deletes multiple documents from the index by identifiers.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/delete-document
	DeleteDocuments(identifiers []string, opts *DocumentOptions) (*TaskInfo, error)

	// DeleteDocumentsWithContext deletes multiple documents from the index by identifiers using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/delete-document
	DeleteDocumentsWithContext(ctx context.Context, identifiers []string, opts *DocumentOptions) (*TaskInfo, error)

	// DeleteDocumentsByFilter deletes documents from the index by filter.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/delete-documents-by-filter
	DeleteDocumentsByFilter(filter interface{}, opts *DocumentOptions) (*TaskInfo, error)

	// DeleteDocumentsByFilterWithContext deletes documents from the index by filter using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/delete-documents-by-filter
	DeleteDocumentsByFilterWithContext(ctx context.Context, filter interface{}, opts *DocumentOptions) (*TaskInfo, error)

	// DeleteAllDocuments deletes all documents from the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/delete-all-documents
	DeleteAllDocuments(opts *DocumentOptions) (*TaskInfo, error)

	// DeleteAllDocumentsWithContext deletes all documents from the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/delete-all-documents
	DeleteAllDocumentsWithContext(ctx context.Context, opts *DocumentOptions) (*TaskInfo, error)
}

type DocumentReader interface {
	// GetDocument retrieves a single document from the index by identifier.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/get-document
	GetDocument(identifier string, request *DocumentQuery, documentPtr interface{}) error

	// GetDocumentWithContext retrieves a single document from the index by identifier using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/get-document
	GetDocumentWithContext(ctx context.Context, identifier string, request *DocumentQuery, documentPtr interface{}) error

	// GetDocuments retrieves multiple documents from the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/list-documents-with-get
	GetDocuments(param *DocumentsQuery, resp *DocumentsResult) error

	// GetDocumentsWithContext retrieves multiple documents from the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/list-documents-with-get
	GetDocumentsWithContext(ctx context.Context, param *DocumentsQuery, resp *DocumentsResult) error
}

type SearchReader interface {
	// Search performs a search query on the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/search/search-with-get
	Search(query string, request *SearchRequest) (*SearchResponse, error)

	// SearchWithContext performs a search query on the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/search/search-with-get
	SearchWithContext(ctx context.Context, query string, request *SearchRequest) (*SearchResponse, error)

	// SearchRaw performs a raw search query on the index, returning a JSON response.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/search/search-with-get
	SearchRaw(query string, request *SearchRequest) (*json.RawMessage, error)

	// SearchRawWithContext performs a raw search query on the index using the provided context for cancellation, returning a JSON response.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/search/search-with-get
	SearchRawWithContext(ctx context.Context, query string, request *SearchRequest) (*json.RawMessage, error)

	// FacetSearch performs a facet search query on the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/facet-search/search-in-facets
	FacetSearch(request *FacetSearchRequest) (*json.RawMessage, error)

	// FacetSearchWithContext performs a facet search query on the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/facet-search/search-in-facets
	FacetSearchWithContext(ctx context.Context, request *FacetSearchRequest) (*json.RawMessage, error)

	// SearchSimilarDocuments performs a search for similar documents.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/similar-documents/get-similar-documents-with-get
	SearchSimilarDocuments(param *SimilarDocumentQuery, resp *SimilarDocumentResult) error

	// SearchSimilarDocumentsWithContext performs a search for similar documents using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/similar-documents/get-similar-documents-with-get
	SearchSimilarDocumentsWithContext(ctx context.Context, param *SimilarDocumentQuery, resp *SimilarDocumentResult) error
}

type SettingsManager interface {
	SettingsReader

	// UpdateSettings updates the settings of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-all-settings
	UpdateSettings(request *Settings) (*TaskInfo, error)

	// UpdateSettingsWithContext updates the settings of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-all-settings
	UpdateSettingsWithContext(ctx context.Context, request *Settings) (*TaskInfo, error)

	// ResetSettings resets the settings of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-all-settings
	ResetSettings() (*TaskInfo, error)

	// ResetSettingsWithContext resets the settings of the index to default values using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-all-settings
	ResetSettingsWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateRankingRules updates the ranking rules of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-rankingrules
	UpdateRankingRules(request *[]string) (*TaskInfo, error)

	// UpdateRankingRulesWithContext updates the ranking rules of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-rankingrules
	UpdateRankingRulesWithContext(ctx context.Context, request *[]string) (*TaskInfo, error)

	// ResetRankingRules resets the ranking rules of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-rankingrules
	ResetRankingRules() (*TaskInfo, error)

	// ResetRankingRulesWithContext resets the ranking rules of the index to default values using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-rankingrules
	ResetRankingRulesWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateDistinctAttribute updates the distinct attribute of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-distinctattribute
	UpdateDistinctAttribute(request string) (*TaskInfo, error)

	// UpdateDistinctAttributeWithContext updates the distinct attribute of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-distinctattribute
	UpdateDistinctAttributeWithContext(ctx context.Context, request string) (*TaskInfo, error)

	// ResetDistinctAttribute resets the distinct attribute of the index to default value.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-distinctattribute
	ResetDistinctAttribute() (*TaskInfo, error)

	// ResetDistinctAttributeWithContext resets the distinct attribute of the index to default value using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-distinctattribute
	ResetDistinctAttributeWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateSearchableAttributes updates the searchable attributes of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-searchableattributes
	UpdateSearchableAttributes(request *[]string) (*TaskInfo, error)

	// UpdateSearchableAttributesWithContext updates the searchable attributes of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-searchableattributes
	UpdateSearchableAttributesWithContext(ctx context.Context, request *[]string) (*TaskInfo, error)

	// ResetSearchableAttributes resets the searchable attributes of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-searchableattributes
	ResetSearchableAttributes() (*TaskInfo, error)

	// ResetSearchableAttributesWithContext resets the searchable attributes of the index to default values using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-searchableattributes
	ResetSearchableAttributesWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateDisplayedAttributes updates the displayed attributes of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-displayedattributes
	UpdateDisplayedAttributes(request *[]string) (*TaskInfo, error)

	// UpdateDisplayedAttributesWithContext updates the displayed attributes of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-displayedattributes
	UpdateDisplayedAttributesWithContext(ctx context.Context, request *[]string) (*TaskInfo, error)

	// ResetDisplayedAttributes resets the displayed attributes of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-displayedattributes
	ResetDisplayedAttributes() (*TaskInfo, error)

	// ResetDisplayedAttributesWithContext resets the displayed attributes of the index to default values using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-displayedattributes
	ResetDisplayedAttributesWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateStopWords updates the stop words of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-stopwords
	UpdateStopWords(request *[]string) (*TaskInfo, error)

	// UpdateStopWordsWithContext updates the stop words of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-stopwords
	UpdateStopWordsWithContext(ctx context.Context, request *[]string) (*TaskInfo, error)

	// ResetStopWords resets the stop words of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-stopwords
	ResetStopWords() (*TaskInfo, error)

	// ResetStopWordsWithContext resets the stop words of the index to default values using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-stopwords
	ResetStopWordsWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateSynonyms updates the synonyms of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-synonyms
	UpdateSynonyms(request *map[string][]string) (*TaskInfo, error)

	// UpdateSynonymsWithContext updates the synonyms of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-synonyms
	UpdateSynonymsWithContext(ctx context.Context, request *map[string][]string) (*TaskInfo, error)

	// ResetSynonyms resets the synonyms of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-synonyms
	ResetSynonyms() (*TaskInfo, error)

	// ResetSynonymsWithContext resets the synonyms of the index to default values using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-synonyms
	ResetSynonymsWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateFilterableAttributes updates the filterable attributes of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-filterableattributes
	UpdateFilterableAttributes(request *[]interface{}) (*TaskInfo, error)

	// UpdateFilterableAttributesWithContext updates the filterable attributes of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-filterableattributes
	UpdateFilterableAttributesWithContext(ctx context.Context, request *[]interface{}) (*TaskInfo, error)

	// ResetFilterableAttributes resets the filterable attributes of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-filterableattributes
	ResetFilterableAttributes() (*TaskInfo, error)

	// ResetFilterableAttributesWithContext resets the filterable attributes of the index to default values using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-filterableattributes
	ResetFilterableAttributesWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateSortableAttributes updates the sortable attributes of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-sortableattributes
	UpdateSortableAttributes(request *[]string) (*TaskInfo, error)

	// UpdateSortableAttributesWithContext updates the sortable attributes of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-sortableattributes
	UpdateSortableAttributesWithContext(ctx context.Context, request *[]string) (*TaskInfo, error)

	// ResetSortableAttributes resets the sortable attributes of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-sortableattributes
	ResetSortableAttributes() (*TaskInfo, error)

	// ResetSortableAttributesWithContext resets the sortable attributes of the index to default values using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-sortableattributes
	ResetSortableAttributesWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateTypoTolerance updates the typo tolerance settings of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-typotolerance
	UpdateTypoTolerance(request *TypoTolerance) (*TaskInfo, error)

	// UpdateTypoToleranceWithContext updates the typo tolerance settings of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-typotolerance
	UpdateTypoToleranceWithContext(ctx context.Context, request *TypoTolerance) (*TaskInfo, error)

	// ResetTypoTolerance resets the typo tolerance settings of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-typotolerance
	ResetTypoTolerance() (*TaskInfo, error)

	// ResetTypoToleranceWithContext resets the typo tolerance settings of the index to default values using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-typotolerance
	ResetTypoToleranceWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdatePagination updates the pagination settings of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-pagination
	UpdatePagination(request *Pagination) (*TaskInfo, error)

	// UpdatePaginationWithContext updates the pagination settings of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-pagination
	UpdatePaginationWithContext(ctx context.Context, request *Pagination) (*TaskInfo, error)

	// ResetPagination resets the pagination settings of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-pagination
	ResetPagination() (*TaskInfo, error)

	// ResetPaginationWithContext resets the pagination settings of the index to default values using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-pagination
	ResetPaginationWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateFaceting updates the faceting settings of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-faceting
	UpdateFaceting(request *Faceting) (*TaskInfo, error)

	// UpdateFacetingWithContext updates the faceting settings of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-faceting
	UpdateFacetingWithContext(ctx context.Context, request *Faceting) (*TaskInfo, error)

	// ResetFaceting resets the faceting settings of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-faceting
	ResetFaceting() (*TaskInfo, error)

	// ResetFacetingWithContext resets the faceting settings of the index to default values using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-faceting
	ResetFacetingWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateEmbedders updates the embedders of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-embedders
	UpdateEmbedders(request map[string]Embedder) (*TaskInfo, error)

	// UpdateEmbeddersWithContext updates the embedders of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-embedders
	UpdateEmbeddersWithContext(ctx context.Context, request map[string]Embedder) (*TaskInfo, error)

	// ResetEmbedders resets the embedders of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-embedders
	ResetEmbedders() (*TaskInfo, error)

	// ResetEmbeddersWithContext resets the embedders of the index to default values using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-embedders
	ResetEmbeddersWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateSearchCutoffMs updates the search cutoff time in milliseconds.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-searchcutoffms
	UpdateSearchCutoffMs(request int64) (*TaskInfo, error)

	// UpdateSearchCutoffMsWithContext updates the search cutoff time in milliseconds using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-searchcutoffms
	UpdateSearchCutoffMsWithContext(ctx context.Context, request int64) (*TaskInfo, error)

	// ResetSearchCutoffMs resets the search cutoff time in milliseconds to default value.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-searchcutoffms
	ResetSearchCutoffMs() (*TaskInfo, error)

	// ResetSearchCutoffMsWithContext resets the search cutoff time in milliseconds to default value using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-searchcutoffms
	ResetSearchCutoffMsWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateSeparatorTokens update separator tokens
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-separatortokens
	UpdateSeparatorTokens(tokens []string) (*TaskInfo, error)

	// UpdateSeparatorTokensWithContext update separator tokens and support parent context
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-separatortokens
	UpdateSeparatorTokensWithContext(ctx context.Context, tokens []string) (*TaskInfo, error)

	// ResetSeparatorTokens reset separator tokens
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-separatortokens
	ResetSeparatorTokens() (*TaskInfo, error)

	// ResetSeparatorTokensWithContext reset separator tokens and support parent context
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-separatortokens
	ResetSeparatorTokensWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateNonSeparatorTokens update non-separator tokens
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-nonseparatortokens
	UpdateNonSeparatorTokens(tokens []string) (*TaskInfo, error)

	// UpdateNonSeparatorTokensWithContext update non-separator tokens and support parent context
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-nonseparatortokens
	UpdateNonSeparatorTokensWithContext(ctx context.Context, tokens []string) (*TaskInfo, error)

	// ResetNonSeparatorTokens reset non-separator tokens
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-nonseparatortokens
	ResetNonSeparatorTokens() (*TaskInfo, error)

	// ResetNonSeparatorTokensWithContext reset non-separator tokens and support parent context
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-nonseparatortokens
	ResetNonSeparatorTokensWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateDictionary update user dictionary
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-dictionary
	UpdateDictionary(words []string) (*TaskInfo, error)

	// UpdateDictionaryWithContext update user dictionary and support parent context
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-dictionary
	UpdateDictionaryWithContext(ctx context.Context, words []string) (*TaskInfo, error)

	// ResetDictionary reset user dictionary
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-dictionary
	ResetDictionary() (*TaskInfo, error)

	// ResetDictionaryWithContext reset user dictionary and support parent context
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-dictionary
	ResetDictionaryWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateProximityPrecision set ProximityPrecision value ByWord or ByAttribute
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-proximityprecision
	UpdateProximityPrecision(proximityType ProximityPrecisionType) (*TaskInfo, error)

	// UpdateProximityPrecisionWithContext set ProximityPrecision value ByWord or ByAttribute and support parent context
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-proximityprecision
	UpdateProximityPrecisionWithContext(ctx context.Context, proximityType ProximityPrecisionType) (*TaskInfo, error)

	// ResetProximityPrecision reset ProximityPrecision to default ByWord
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-proximityprecision
	ResetProximityPrecision() (*TaskInfo, error)

	// ResetProximityPrecisionWithContext reset ProximityPrecision to default ByWord and support parent context
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-proximityprecision
	ResetProximityPrecisionWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateLocalizedAttributes update the localized attributes settings of an index
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-localizedattributes
	UpdateLocalizedAttributes(request []*LocalizedAttributes) (*TaskInfo, error)

	// UpdateLocalizedAttributesWithContext update the localized attributes settings of an index using the provided context for cancellation
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-localizedattributes
	UpdateLocalizedAttributesWithContext(ctx context.Context, request []*LocalizedAttributes) (*TaskInfo, error)

	// ResetLocalizedAttributes reset the localized attributes settings
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-localizedattributes
	ResetLocalizedAttributes() (*TaskInfo, error)

	// ResetLocalizedAttributesWithContext reset the localized attributes settings using the provided context for cancellation
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-localizedattributes
	ResetLocalizedAttributesWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdatePrefixSearch updates the prefix search setting of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-prefixsearch
	UpdatePrefixSearch(request string) (*TaskInfo, error)

	// UpdatePrefixSearchWithContext updates the prefix search setting of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-prefixsearch
	UpdatePrefixSearchWithContext(ctx context.Context, request string) (*TaskInfo, error)

	// ResetPrefixSearch resets the prefix search setting of the index to default value.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-prefixsearch
	ResetPrefixSearch() (*TaskInfo, error)

	// ResetPrefixSearchWithContext resets the prefix search setting of the index to default value using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-prefixsearch
	ResetPrefixSearchWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateFacetSearch updates the facet search setting of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-facetsearch
	UpdateFacetSearch(request bool) (*TaskInfo, error)

	// UpdateFacetSearchWithContext updates the facet search setting of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-facetsearch
	UpdateFacetSearchWithContext(ctx context.Context, request bool) (*TaskInfo, error)

	// ResetFacetSearch resets the facet search setting of the index to default value.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-facetsearch
	ResetFacetSearch() (*TaskInfo, error)

	// ResetFacetSearchWithContext resets the facet search setting of the index to default value using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-facetsearch
	ResetFacetSearchWithContext(ctx context.Context) (*TaskInfo, error)
}

type SettingsReader interface {
	// GetSettings retrieves the settings of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/list-all-settings
	GetSettings() (*Settings, error)

	// GetSettingsWithContext retrieves the settings of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/list-all-settings
	GetSettingsWithContext(ctx context.Context) (*Settings, error)

	// GetRankingRules retrieves the ranking rules of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-rankingrules
	GetRankingRules() (*[]string, error)

	// GetRankingRulesWithContext retrieves the ranking rules of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-rankingrules
	GetRankingRulesWithContext(ctx context.Context) (*[]string, error)

	// GetDistinctAttribute retrieves the distinct attribute of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-distinctattribute
	GetDistinctAttribute() (*string, error)

	// GetDistinctAttributeWithContext retrieves the distinct attribute of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-distinctattribute
	GetDistinctAttributeWithContext(ctx context.Context) (*string, error)

	// GetSearchableAttributes retrieves the searchable attributes of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-searchableattributes
	GetSearchableAttributes() (*[]string, error)

	// GetSearchableAttributesWithContext retrieves the searchable attributes of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-searchableattributes
	GetSearchableAttributesWithContext(ctx context.Context) (*[]string, error)

	// GetDisplayedAttributes retrieves the displayed attributes of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-displayedattributes
	GetDisplayedAttributes() (*[]string, error)

	// GetDisplayedAttributesWithContext retrieves the displayed attributes of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-displayedattributes
	GetDisplayedAttributesWithContext(ctx context.Context) (*[]string, error)

	// GetStopWords retrieves the stop words of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-stopwords
	GetStopWords() (*[]string, error)

	// GetStopWordsWithContext retrieves the stop words of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-stopwords
	GetStopWordsWithContext(ctx context.Context) (*[]string, error)

	// GetSynonyms retrieves the synonyms of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-synonyms
	GetSynonyms() (*map[string][]string, error)

	// GetSynonymsWithContext retrieves the synonyms of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-synonyms
	GetSynonymsWithContext(ctx context.Context) (*map[string][]string, error)

	// GetFilterableAttributes retrieves the filterable attributes of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-filterableattributes
	GetFilterableAttributes() (*[]interface{}, error)

	// GetFilterableAttributesWithContext retrieves the filterable attributes of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-filterableattributes
	GetFilterableAttributesWithContext(ctx context.Context) (*[]interface{}, error)

	// GetSortableAttributes retrieves the sortable attributes of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-sortableattributes
	GetSortableAttributes() (*[]string, error)

	// GetSortableAttributesWithContext retrieves the sortable attributes of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-sortableattributes
	GetSortableAttributesWithContext(ctx context.Context) (*[]string, error)

	// GetTypoTolerance retrieves the typo tolerance settings of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-typotolerance
	GetTypoTolerance() (*TypoTolerance, error)

	// GetTypoToleranceWithContext retrieves the typo tolerance settings of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-typotolerance
	GetTypoToleranceWithContext(ctx context.Context) (*TypoTolerance, error)

	// GetPagination retrieves the pagination settings of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-pagination
	GetPagination() (*Pagination, error)

	// GetPaginationWithContext retrieves the pagination settings of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-pagination
	GetPaginationWithContext(ctx context.Context) (*Pagination, error)

	// GetFaceting retrieves the faceting settings of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-faceting
	GetFaceting() (*Faceting, error)

	// GetFacetingWithContext retrieves the faceting settings of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-faceting
	GetFacetingWithContext(ctx context.Context) (*Faceting, error)

	// GetEmbedders retrieves the embedders of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-embedders
	GetEmbedders() (map[string]Embedder, error)

	// GetEmbeddersWithContext retrieves the embedders of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-embedders
	GetEmbeddersWithContext(ctx context.Context) (map[string]Embedder, error)

	// GetSearchCutoffMs retrieves the search cutoff time in milliseconds.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-searchcutoffms
	GetSearchCutoffMs() (int64, error)

	// GetSearchCutoffMsWithContext retrieves the search cutoff time in milliseconds using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-searchcutoffms
	GetSearchCutoffMsWithContext(ctx context.Context) (int64, error)

	// GetSeparatorTokens returns separators tokens
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-separatortokens
	GetSeparatorTokens() ([]string, error)

	// GetSeparatorTokensWithContext returns separator tokens and support parent context
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-separatortokens
	GetSeparatorTokensWithContext(ctx context.Context) ([]string, error)

	// GetNonSeparatorTokens returns non-separator tokens
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-nonseparatortokens
	GetNonSeparatorTokens() ([]string, error)

	// GetNonSeparatorTokensWithContext returns non-separator tokens and support parent context
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-nonseparatortokens
	GetNonSeparatorTokensWithContext(ctx context.Context) ([]string, error)

	// GetDictionary returns user dictionary
	//
	//Allows users to instruct Meilisearch to consider groups of strings as a
	//single term by adding a supplementary dictionary of user-defined terms.
	//This is particularly useful when working with datasets containing many domain-specific
	//words, and in languages where words are not separated by whitespace such as Japanese.
	//Custom dictionaries are also useful in a few use-cases for space-separated languages,
	//such as datasets with names such as "J. R. R. Tolkien" and "W. E. B. Du Bois".
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-dictionary
	GetDictionary() ([]string, error)

	// GetDictionaryWithContext returns user dictionary and support parent context
	//
	//Allows users to instruct Meilisearch to consider groups of strings as a
	//single term by adding a supplementary dictionary of user-defined terms.
	//This is particularly useful when working with datasets containing many domain-specific
	//words, and in languages where words are not separated by whitespace such as Japanese.
	//Custom dictionaries are also useful in a few use-cases for space-separated languages,
	//such as datasets with names such as "J. R. R. Tolkien" and "W. E. B. Du Bois".
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-dictionary
	GetDictionaryWithContext(ctx context.Context) ([]string, error)

	// GetProximityPrecision returns ProximityPrecision configuration value
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-proximityprecision
	GetProximityPrecision() (ProximityPrecisionType, error)

	// GetProximityPrecisionWithContext returns ProximityPrecision configuration value and support parent context
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-proximityprecision
	GetProximityPrecisionWithContext(ctx context.Context) (ProximityPrecisionType, error)

	// GetLocalizedAttributes get the localized attributes settings of an index
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-localizedattributes
	GetLocalizedAttributes() ([]*LocalizedAttributes, error)

	// GetLocalizedAttributesWithContext get the localized attributes settings of an index using the provided context for cancellation
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-localizedattributes
	GetLocalizedAttributesWithContext(ctx context.Context) ([]*LocalizedAttributes, error)

	// GetPrefixSearch retrieves the prefix search setting of the index.
	//
	// https://www.meilisearch.com/docs/reference/api/settings/get-prefixsearch
	GetPrefixSearch() (*string, error)

	// GetPrefixSearchWithContext retrieves the prefix search setting of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-prefixsearch
	GetPrefixSearchWithContext(ctx context.Context) (*string, error)

	// GetFacetSearch retrieves the facet search setting of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-facetsearch
	GetFacetSearch() (bool, error)

	// GetFacetSearchWithContext retrieves the facet search setting of the index using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-facetsearch
	GetFacetSearchWithContext(ctx context.Context) (bool, error)
}
