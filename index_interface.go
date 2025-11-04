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
	UpdateIndex(params *UpdateIndexRequestParams) (*TaskInfo, error)

	// UpdateIndexWithContext updates the primary key of the index using the provided context for cancellation.
	UpdateIndexWithContext(ctx context.Context, params *UpdateIndexRequestParams) (*TaskInfo, error)

	// Delete removes the index identified by the given UID.
	Delete(uid string) (bool, error)

	// DeleteWithContext removes the index identified by the given UID using the provided context for cancellation.
	DeleteWithContext(ctx context.Context, uid string) (bool, error)

	// Compact compacts the index.
	Compact() (*TaskInfo, error)

	// CompactWithContext compacts the index using the provided context for cancellation.
	CompactWithContext(ctx context.Context) (*TaskInfo, error)
}

type IndexReader interface {
	// FetchInfo retrieves information about the index.
	FetchInfo() (*IndexResult, error)

	// FetchInfoWithContext retrieves information about the index using the provided context for cancellation.
	FetchInfoWithContext(ctx context.Context) (*IndexResult, error)

	// FetchPrimaryKey retrieves the primary key of the index.
	FetchPrimaryKey() (*string, error)

	// FetchPrimaryKeyWithContext retrieves the primary key of the index using the provided context for cancellation.
	FetchPrimaryKeyWithContext(ctx context.Context) (*string, error)

	// GetStats retrieves statistical information about the index.
	GetStats() (*StatsIndex, error)

	// GetStatsWithContext retrieves statistical information about the index using the provided context for cancellation.
	GetStatsWithContext(ctx context.Context) (*StatsIndex, error)
}

type DocumentManager interface {
	DocumentReader

	// AddDocuments adds multiple documents to the index.
	AddDocuments(documentsPtr interface{}, primaryKey *string) (*TaskInfo, error)

	// AddDocumentsWithContext adds multiple documents to the index using the provided context for cancellation.
	AddDocumentsWithContext(ctx context.Context, documentsPtr interface{}, primaryKey *string) (*TaskInfo, error)

	// AddDocumentsInBatches adds documents to the index in batches of specified size.
	AddDocumentsInBatches(documentsPtr interface{}, batchSize int, primaryKey *string) ([]TaskInfo, error)

	// AddDocumentsInBatchesWithContext adds documents to the index in batches of specified size using the provided context for cancellation.
	AddDocumentsInBatchesWithContext(ctx context.Context, documentsPtr interface{}, batchSize int, primaryKey *string) ([]TaskInfo, error)

	// AddDocumentsCsv adds documents from a CSV byte array to the index.
	AddDocumentsCsv(documents []byte, options *CsvDocumentsQuery) (*TaskInfo, error)

	// AddDocumentsCsvWithContext adds documents from a CSV byte array to the index using the provided context for cancellation.
	AddDocumentsCsvWithContext(ctx context.Context, documents []byte, options *CsvDocumentsQuery) (*TaskInfo, error)

	// AddDocumentsCsvInBatches adds documents from a CSV byte array to the index in batches of specified size.
	AddDocumentsCsvInBatches(documents []byte, batchSize int, options *CsvDocumentsQuery) ([]TaskInfo, error)

	// AddDocumentsCsvInBatchesWithContext adds documents from a CSV byte array to the index in batches of specified size using the provided context for cancellation.
	AddDocumentsCsvInBatchesWithContext(ctx context.Context, documents []byte, batchSize int, options *CsvDocumentsQuery) ([]TaskInfo, error)

	// AddDocumentsCsvFromReaderInBatches adds documents from a CSV reader to the index in batches of specified size.
	AddDocumentsCsvFromReaderInBatches(documents io.Reader, batchSize int, options *CsvDocumentsQuery) ([]TaskInfo, error)

	// AddDocumentsCsvFromReaderInBatchesWithContext adds documents from a CSV reader to the index in batches of specified size using the provided context for cancellation.
	AddDocumentsCsvFromReaderInBatchesWithContext(ctx context.Context, documents io.Reader, batchSize int, options *CsvDocumentsQuery) ([]TaskInfo, error)

	// AddDocumentsCsvFromReader adds documents from a CSV reader to the index.
	AddDocumentsCsvFromReader(documents io.Reader, options *CsvDocumentsQuery) (*TaskInfo, error)

	// AddDocumentsCsvFromReaderWithContext adds documents from a CSV reader to the index using the provided context for cancellation.
	AddDocumentsCsvFromReaderWithContext(ctx context.Context, documents io.Reader, options *CsvDocumentsQuery) (*TaskInfo, error)

	// AddDocumentsNdjson adds documents from a NDJSON byte array to the index.
	AddDocumentsNdjson(documents []byte, primaryKey *string) (*TaskInfo, error)

	// AddDocumentsNdjsonWithContext adds documents from a NDJSON byte array to the index using the provided context for cancellation.
	AddDocumentsNdjsonWithContext(ctx context.Context, documents []byte, primaryKey *string) (*TaskInfo, error)

	// AddDocumentsNdjsonInBatches adds documents from a NDJSON byte array to the index in batches of specified size.
	AddDocumentsNdjsonInBatches(documents []byte, batchSize int, primaryKey *string) ([]TaskInfo, error)

	// AddDocumentsNdjsonInBatchesWithContext adds documents from a NDJSON byte array to the index in batches of specified size using the provided context for cancellation.
	AddDocumentsNdjsonInBatchesWithContext(ctx context.Context, documents []byte, batchSize int, primaryKey *string) ([]TaskInfo, error)

	// AddDocumentsNdjsonFromReader adds documents from a NDJSON reader to the index.
	AddDocumentsNdjsonFromReader(documents io.Reader, primaryKey *string) (*TaskInfo, error)

	// AddDocumentsNdjsonFromReaderWithContext adds documents from a NDJSON reader to the index using the provided context for cancellation.
	AddDocumentsNdjsonFromReaderWithContext(ctx context.Context, documents io.Reader, primaryKey *string) (*TaskInfo, error)

	// AddDocumentsNdjsonFromReaderInBatches adds documents from a NDJSON reader to the index in batches of specified size.
	AddDocumentsNdjsonFromReaderInBatches(documents io.Reader, batchSize int, primaryKey *string) ([]TaskInfo, error)

	// AddDocumentsNdjsonFromReaderInBatchesWithContext adds documents from a NDJSON reader to the index in batches of specified size using the provided context for cancellation.
	AddDocumentsNdjsonFromReaderInBatchesWithContext(ctx context.Context, documents io.Reader, batchSize int, primaryKey *string) ([]TaskInfo, error)

	// UpdateDocuments updates multiple documents in the index.
	UpdateDocuments(documentsPtr interface{}, primaryKey *string) (*TaskInfo, error)

	// UpdateDocumentsWithContext updates multiple documents in the index using the provided context for cancellation.
	UpdateDocumentsWithContext(ctx context.Context, documentsPtr interface{}, primaryKey *string) (*TaskInfo, error)

	// UpdateDocumentsInBatches updates documents in the index in batches of specified size.
	UpdateDocumentsInBatches(documentsPtr interface{}, batchSize int, primaryKey *string) ([]TaskInfo, error)

	// UpdateDocumentsInBatchesWithContext updates documents in the index in batches of specified size using the provided context for cancellation.
	UpdateDocumentsInBatchesWithContext(ctx context.Context, documentsPtr interface{}, batchSize int, primaryKey *string) ([]TaskInfo, error)

	// UpdateDocumentsCsv updates documents in the index from a CSV byte array.
	UpdateDocumentsCsv(documents []byte, options *CsvDocumentsQuery) (*TaskInfo, error)

	// UpdateDocumentsCsvWithContext updates documents in the index from a CSV byte array using the provided context for cancellation.
	UpdateDocumentsCsvWithContext(ctx context.Context, documents []byte, options *CsvDocumentsQuery) (*TaskInfo, error)

	// UpdateDocumentsCsvInBatches updates documents in the index from a CSV byte array in batches of specified size.
	UpdateDocumentsCsvInBatches(documents []byte, batchsize int, options *CsvDocumentsQuery) ([]TaskInfo, error)

	// UpdateDocumentsCsvInBatchesWithContext updates documents in the index from a CSV byte array in batches of specified size using the provided context for cancellation.
	UpdateDocumentsCsvInBatchesWithContext(ctx context.Context, documents []byte, batchsize int, options *CsvDocumentsQuery) ([]TaskInfo, error)

	// UpdateDocumentsNdjson updates documents in the index from a NDJSON byte array.
	UpdateDocumentsNdjson(documents []byte, primaryKey *string) (*TaskInfo, error)

	// UpdateDocumentsNdjsonWithContext updates documents in the index from a NDJSON byte array using the provided context for cancellation.
	UpdateDocumentsNdjsonWithContext(ctx context.Context, documents []byte, primaryKey *string) (*TaskInfo, error)

	// UpdateDocumentsNdjsonInBatches updates documents in the index from a NDJSON byte array in batches of specified size.
	UpdateDocumentsNdjsonInBatches(documents []byte, batchsize int, primaryKey *string) ([]TaskInfo, error)

	// UpdateDocumentsNdjsonInBatchesWithContext updates documents in the index from a NDJSON byte array in batches of specified size using the provided context for cancellation.
	UpdateDocumentsNdjsonInBatchesWithContext(ctx context.Context, documents []byte, batchsize int, primaryKey *string) ([]TaskInfo, error)

	// UpdateDocumentsByFunction update documents by using function
	UpdateDocumentsByFunction(req *UpdateDocumentByFunctionRequest) (*TaskInfo, error)

	// UpdateDocumentsByFunctionWithContext update documents by using function then provided context for cancellation.
	UpdateDocumentsByFunctionWithContext(ctx context.Context, req *UpdateDocumentByFunctionRequest) (*TaskInfo, error)

	// DeleteDocument deletes a single document from the index by identifier.
	DeleteDocument(identifier string) (*TaskInfo, error)

	// DeleteDocumentWithContext deletes a single document from the index by identifier using the provided context for cancellation.
	DeleteDocumentWithContext(ctx context.Context, identifier string) (*TaskInfo, error)

	// DeleteDocuments deletes multiple documents from the index by identifiers.
	DeleteDocuments(identifiers []string) (*TaskInfo, error)

	// DeleteDocumentsWithContext deletes multiple documents from the index by identifiers using the provided context for cancellation.
	DeleteDocumentsWithContext(ctx context.Context, identifiers []string) (*TaskInfo, error)

	// DeleteDocumentsByFilter deletes documents from the index by filter.
	DeleteDocumentsByFilter(filter interface{}) (*TaskInfo, error)

	// DeleteDocumentsByFilterWithContext deletes documents from the index by filter using the provided context for cancellation.
	DeleteDocumentsByFilterWithContext(ctx context.Context, filter interface{}) (*TaskInfo, error)

	// DeleteAllDocuments deletes all documents from the index.
	DeleteAllDocuments() (*TaskInfo, error)

	// DeleteAllDocumentsWithContext deletes all documents from the index using the provided context for cancellation.
	DeleteAllDocumentsWithContext(ctx context.Context) (*TaskInfo, error)
}

type DocumentReader interface {
	// GetDocument retrieves a single document from the index by identifier.
	GetDocument(identifier string, request *DocumentQuery, documentPtr interface{}) error

	// GetDocumentWithContext retrieves a single document from the index by identifier using the provided context for cancellation.
	GetDocumentWithContext(ctx context.Context, identifier string, request *DocumentQuery, documentPtr interface{}) error

	// GetDocuments retrieves multiple documents from the index.
	GetDocuments(param *DocumentsQuery, resp *DocumentsResult) error

	// GetDocumentsWithContext retrieves multiple documents from the index using the provided context for cancellation.
	GetDocumentsWithContext(ctx context.Context, param *DocumentsQuery, resp *DocumentsResult) error
}

type SearchReader interface {
	// Search performs a search query on the index.
	Search(query string, request *SearchRequest) (*SearchResponse, error)

	// SearchWithContext performs a search query on the index using the provided context for cancellation.
	SearchWithContext(ctx context.Context, query string, request *SearchRequest) (*SearchResponse, error)

	// SearchRaw performs a raw search query on the index, returning a JSON response.
	SearchRaw(query string, request *SearchRequest) (*json.RawMessage, error)

	// SearchRawWithContext performs a raw search query on the index using the provided context for cancellation, returning a JSON response.
	SearchRawWithContext(ctx context.Context, query string, request *SearchRequest) (*json.RawMessage, error)

	// FacetSearch performs a facet search query on the index.
	FacetSearch(request *FacetSearchRequest) (*json.RawMessage, error)

	// FacetSearchWithContext performs a facet search query on the index using the provided context for cancellation.
	FacetSearchWithContext(ctx context.Context, request *FacetSearchRequest) (*json.RawMessage, error)

	// SearchSimilarDocuments performs a search for similar documents.
	SearchSimilarDocuments(param *SimilarDocumentQuery, resp *SimilarDocumentResult) error

	// SearchSimilarDocumentsWithContext performs a search for similar documents using the provided context for cancellation.
	SearchSimilarDocumentsWithContext(ctx context.Context, param *SimilarDocumentQuery, resp *SimilarDocumentResult) error
}

type SettingsManager interface {
	SettingsReader

	// UpdateSettings updates the settings of the index.
	UpdateSettings(request *Settings) (*TaskInfo, error)

	// UpdateSettingsWithContext updates the settings of the index using the provided context for cancellation.
	UpdateSettingsWithContext(ctx context.Context, request *Settings) (*TaskInfo, error)

	// ResetSettings resets the settings of the index to default values.
	ResetSettings() (*TaskInfo, error)

	// ResetSettingsWithContext resets the settings of the index to default values using the provided context for cancellation.
	ResetSettingsWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateRankingRules updates the ranking rules of the index.
	UpdateRankingRules(request *[]string) (*TaskInfo, error)

	// UpdateRankingRulesWithContext updates the ranking rules of the index using the provided context for cancellation.
	UpdateRankingRulesWithContext(ctx context.Context, request *[]string) (*TaskInfo, error)

	// ResetRankingRules resets the ranking rules of the index to default values.
	ResetRankingRules() (*TaskInfo, error)

	// ResetRankingRulesWithContext resets the ranking rules of the index to default values using the provided context for cancellation.
	ResetRankingRulesWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateDistinctAttribute updates the distinct attribute of the index.
	UpdateDistinctAttribute(request string) (*TaskInfo, error)

	// UpdateDistinctAttributeWithContext updates the distinct attribute of the index using the provided context for cancellation.
	UpdateDistinctAttributeWithContext(ctx context.Context, request string) (*TaskInfo, error)

	// ResetDistinctAttribute resets the distinct attribute of the index to default value.
	ResetDistinctAttribute() (*TaskInfo, error)

	// ResetDistinctAttributeWithContext resets the distinct attribute of the index to default value using the provided context for cancellation.
	ResetDistinctAttributeWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateSearchableAttributes updates the searchable attributes of the index.
	UpdateSearchableAttributes(request *[]string) (*TaskInfo, error)

	// UpdateSearchableAttributesWithContext updates the searchable attributes of the index using the provided context for cancellation.
	UpdateSearchableAttributesWithContext(ctx context.Context, request *[]string) (*TaskInfo, error)

	// ResetSearchableAttributes resets the searchable attributes of the index to default values.
	ResetSearchableAttributes() (*TaskInfo, error)

	// ResetSearchableAttributesWithContext resets the searchable attributes of the index to default values using the provided context for cancellation.
	ResetSearchableAttributesWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateDisplayedAttributes updates the displayed attributes of the index.
	UpdateDisplayedAttributes(request *[]string) (*TaskInfo, error)

	// UpdateDisplayedAttributesWithContext updates the displayed attributes of the index using the provided context for cancellation.
	UpdateDisplayedAttributesWithContext(ctx context.Context, request *[]string) (*TaskInfo, error)

	// ResetDisplayedAttributes resets the displayed attributes of the index to default values.
	ResetDisplayedAttributes() (*TaskInfo, error)

	// ResetDisplayedAttributesWithContext resets the displayed attributes of the index to default values using the provided context for cancellation.
	ResetDisplayedAttributesWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateStopWords updates the stop words of the index.
	UpdateStopWords(request *[]string) (*TaskInfo, error)

	// UpdateStopWordsWithContext updates the stop words of the index using the provided context for cancellation.
	UpdateStopWordsWithContext(ctx context.Context, request *[]string) (*TaskInfo, error)

	// ResetStopWords resets the stop words of the index to default values.
	ResetStopWords() (*TaskInfo, error)

	// ResetStopWordsWithContext resets the stop words of the index to default values using the provided context for cancellation.
	ResetStopWordsWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateSynonyms updates the synonyms of the index.
	UpdateSynonyms(request *map[string][]string) (*TaskInfo, error)

	// UpdateSynonymsWithContext updates the synonyms of the index using the provided context for cancellation.
	UpdateSynonymsWithContext(ctx context.Context, request *map[string][]string) (*TaskInfo, error)

	// ResetSynonyms resets the synonyms of the index to default values.
	ResetSynonyms() (*TaskInfo, error)

	// ResetSynonymsWithContext resets the synonyms of the index to default values using the provided context for cancellation.
	ResetSynonymsWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateFilterableAttributes updates the filterable attributes of the index.
	UpdateFilterableAttributes(request *[]interface{}) (*TaskInfo, error)

	// UpdateFilterableAttributesWithContext updates the filterable attributes of the index using the provided context for cancellation.
	UpdateFilterableAttributesWithContext(ctx context.Context, request *[]interface{}) (*TaskInfo, error)

	// ResetFilterableAttributes resets the filterable attributes of the index to default values.
	ResetFilterableAttributes() (*TaskInfo, error)

	// ResetFilterableAttributesWithContext resets the filterable attributes of the index to default values using the provided context for cancellation.
	ResetFilterableAttributesWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateSortableAttributes updates the sortable attributes of the index.
	UpdateSortableAttributes(request *[]string) (*TaskInfo, error)

	// UpdateSortableAttributesWithContext updates the sortable attributes of the index using the provided context for cancellation.
	UpdateSortableAttributesWithContext(ctx context.Context, request *[]string) (*TaskInfo, error)

	// ResetSortableAttributes resets the sortable attributes of the index to default values.
	ResetSortableAttributes() (*TaskInfo, error)

	// ResetSortableAttributesWithContext resets the sortable attributes of the index to default values using the provided context for cancellation.
	ResetSortableAttributesWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateTypoTolerance updates the typo tolerance settings of the index.
	UpdateTypoTolerance(request *TypoTolerance) (*TaskInfo, error)

	// UpdateTypoToleranceWithContext updates the typo tolerance settings of the index using the provided context for cancellation.
	UpdateTypoToleranceWithContext(ctx context.Context, request *TypoTolerance) (*TaskInfo, error)

	// ResetTypoTolerance resets the typo tolerance settings of the index to default values.
	ResetTypoTolerance() (*TaskInfo, error)

	// ResetTypoToleranceWithContext resets the typo tolerance settings of the index to default values using the provided context for cancellation.
	ResetTypoToleranceWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdatePagination updates the pagination settings of the index.
	UpdatePagination(request *Pagination) (*TaskInfo, error)

	// UpdatePaginationWithContext updates the pagination settings of the index using the provided context for cancellation.
	UpdatePaginationWithContext(ctx context.Context, request *Pagination) (*TaskInfo, error)

	// ResetPagination resets the pagination settings of the index to default values.
	ResetPagination() (*TaskInfo, error)

	// ResetPaginationWithContext resets the pagination settings of the index to default values using the provided context for cancellation.
	ResetPaginationWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateFaceting updates the faceting settings of the index.
	UpdateFaceting(request *Faceting) (*TaskInfo, error)

	// UpdateFacetingWithContext updates the faceting settings of the index using the provided context for cancellation.
	UpdateFacetingWithContext(ctx context.Context, request *Faceting) (*TaskInfo, error)

	// ResetFaceting resets the faceting settings of the index to default values.
	ResetFaceting() (*TaskInfo, error)

	// ResetFacetingWithContext resets the faceting settings of the index to default values using the provided context for cancellation.
	ResetFacetingWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateEmbedders updates the embedders of the index.
	UpdateEmbedders(request map[string]Embedder) (*TaskInfo, error)

	// UpdateEmbeddersWithContext updates the embedders of the index using the provided context for cancellation.
	UpdateEmbeddersWithContext(ctx context.Context, request map[string]Embedder) (*TaskInfo, error)

	// ResetEmbedders resets the embedders of the index to default values.
	ResetEmbedders() (*TaskInfo, error)

	// ResetEmbeddersWithContext resets the embedders of the index to default values using the provided context for cancellation.
	ResetEmbeddersWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateSearchCutoffMs updates the search cutoff time in milliseconds.
	UpdateSearchCutoffMs(request int64) (*TaskInfo, error)

	// UpdateSearchCutoffMsWithContext updates the search cutoff time in milliseconds using the provided context for cancellation.
	UpdateSearchCutoffMsWithContext(ctx context.Context, request int64) (*TaskInfo, error)

	// ResetSearchCutoffMs resets the search cutoff time in milliseconds to default value.
	ResetSearchCutoffMs() (*TaskInfo, error)

	// ResetSearchCutoffMsWithContext resets the search cutoff time in milliseconds to default value using the provided context for cancellation.
	ResetSearchCutoffMsWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateSeparatorTokens update separator tokens
	// https://www.meilisearch.com/docs/reference/api/settings#update-separator-tokens
	UpdateSeparatorTokens(tokens []string) (*TaskInfo, error)

	// UpdateSeparatorTokensWithContext update separator tokens and support parent context
	// https://www.meilisearch.com/docs/reference/api/settings#update-separator-tokens
	UpdateSeparatorTokensWithContext(ctx context.Context, tokens []string) (*TaskInfo, error)

	// ResetSeparatorTokens reset separator tokens
	// https://www.meilisearch.com/docs/reference/api/settings#reset-separator-tokens
	ResetSeparatorTokens() (*TaskInfo, error)

	// ResetSeparatorTokensWithContext reset separator tokens and support parent context
	// https://www.meilisearch.com/docs/reference/api/settings#reset-separator-tokens
	ResetSeparatorTokensWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateNonSeparatorTokens update non-separator tokens
	// https://www.meilisearch.com/docs/reference/api/settings#update-non-separator-tokens
	UpdateNonSeparatorTokens(tokens []string) (*TaskInfo, error)

	// UpdateNonSeparatorTokensWithContext update non-separator tokens and support parent context
	// https://www.meilisearch.com/docs/reference/api/settings#update-non-separator-tokens
	UpdateNonSeparatorTokensWithContext(ctx context.Context, tokens []string) (*TaskInfo, error)

	// ResetNonSeparatorTokens reset non-separator tokens
	// https://www.meilisearch.com/docs/reference/api/settings#reset-non-separator-tokens
	ResetNonSeparatorTokens() (*TaskInfo, error)

	// ResetNonSeparatorTokensWithContext reset non-separator tokens and support parent context
	// https://www.meilisearch.com/docs/reference/api/settings#reset-non-separator-tokens
	ResetNonSeparatorTokensWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateDictionary update user dictionary
	// https://www.meilisearch.com/docs/reference/api/settings#update-dictionary
	UpdateDictionary(words []string) (*TaskInfo, error)

	// UpdateDictionaryWithContext update user dictionary and support parent context
	// https://www.meilisearch.com/docs/reference/api/settings#update-dictionary
	UpdateDictionaryWithContext(ctx context.Context, words []string) (*TaskInfo, error)

	// ResetDictionary reset user dictionary
	// https://www.meilisearch.com/docs/reference/api/settings#reset-dictionary
	ResetDictionary() (*TaskInfo, error)

	// ResetDictionaryWithContext reset user dictionary and support parent context
	// https://www.meilisearch.com/docs/reference/api/settings#reset-dictionary
	ResetDictionaryWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateProximityPrecision set ProximityPrecision value ByWord or ByAttribute
	// https://www.meilisearch.com/docs/reference/api/settings#update-proximity-precision-settings
	UpdateProximityPrecision(proximityType ProximityPrecisionType) (*TaskInfo, error)

	// UpdateProximityPrecisionWithContext set ProximityPrecision value ByWord or ByAttribute and support parent context
	// https://www.meilisearch.com/docs/reference/api/settings#update-proximity-precision-settings
	UpdateProximityPrecisionWithContext(ctx context.Context, proximityType ProximityPrecisionType) (*TaskInfo, error)

	// ResetProximityPrecision reset ProximityPrecision to default ByWord
	// https://www.meilisearch.com/docs/reference/api/settings#reset-proximity-precision-settings
	ResetProximityPrecision() (*TaskInfo, error)

	// ResetProximityPrecisionWithContext reset ProximityPrecision to default ByWord and support parent context
	// https://www.meilisearch.com/docs/reference/api/settings#reset-proximity-precision-settings
	ResetProximityPrecisionWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateLocalizedAttributes update the localized attributes settings of an index
	// https://www.meilisearch.com/docs/reference/api/settings#update-localized-attribute-settings
	UpdateLocalizedAttributes(request []*LocalizedAttributes) (*TaskInfo, error)

	// UpdateLocalizedAttributesWithContext update the localized attributes settings of an index using the provided context for cancellation
	// https://www.meilisearch.com/docs/reference/api/settings#update-localized-attribute-settings
	UpdateLocalizedAttributesWithContext(ctx context.Context, request []*LocalizedAttributes) (*TaskInfo, error)

	// ResetLocalizedAttributes reset the localized attributes settings
	ResetLocalizedAttributes() (*TaskInfo, error)

	// ResetLocalizedAttributesWithContext reset the localized attributes settings using the provided context for cancellation
	ResetLocalizedAttributesWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdatePrefixSearch updates the prefix search setting of the index.
	UpdatePrefixSearch(request string) (*TaskInfo, error)

	// UpdatePrefixSearchWithContext updates the prefix search setting of the index using the provided context for cancellation.
	UpdatePrefixSearchWithContext(ctx context.Context, request string) (*TaskInfo, error)

	// ResetPrefixSearch resets the prefix search setting of the index to default value.
	ResetPrefixSearch() (*TaskInfo, error)

	// ResetPrefixSearchWithContext resets the prefix search setting of the index to default value using the provided context for cancellation.
	ResetPrefixSearchWithContext(ctx context.Context) (*TaskInfo, error)

	// UpdateFacetSearch updates the facet search setting of the index.
	UpdateFacetSearch(request bool) (*TaskInfo, error)

	// UpdateFacetSearchWithContext updates the facet search setting of the index using the provided context for cancellation.
	UpdateFacetSearchWithContext(ctx context.Context, request bool) (*TaskInfo, error)

	// ResetFacetSearch resets the facet search setting of the index to default value.
	ResetFacetSearch() (*TaskInfo, error)

	// ResetFacetSearchWithContext resets the facet search setting of the index to default value using the provided context for cancellation.
	ResetFacetSearchWithContext(ctx context.Context) (*TaskInfo, error)
}

type SettingsReader interface {
	// GetSettings retrieves the settings of the index.
	GetSettings() (*Settings, error)

	// GetSettingsWithContext retrieves the settings of the index using the provided context for cancellation.
	GetSettingsWithContext(ctx context.Context) (*Settings, error)

	// GetRankingRules retrieves the ranking rules of the index.
	GetRankingRules() (*[]string, error)

	// GetRankingRulesWithContext retrieves the ranking rules of the index using the provided context for cancellation.
	GetRankingRulesWithContext(ctx context.Context) (*[]string, error)

	// GetDistinctAttribute retrieves the distinct attribute of the index.
	GetDistinctAttribute() (*string, error)

	// GetDistinctAttributeWithContext retrieves the distinct attribute of the index using the provided context for cancellation.
	GetDistinctAttributeWithContext(ctx context.Context) (*string, error)

	// GetSearchableAttributes retrieves the searchable attributes of the index.
	GetSearchableAttributes() (*[]string, error)

	// GetSearchableAttributesWithContext retrieves the searchable attributes of the index using the provided context for cancellation.
	GetSearchableAttributesWithContext(ctx context.Context) (*[]string, error)

	// GetDisplayedAttributes retrieves the displayed attributes of the index.
	GetDisplayedAttributes() (*[]string, error)

	// GetDisplayedAttributesWithContext retrieves the displayed attributes of the index using the provided context for cancellation.
	GetDisplayedAttributesWithContext(ctx context.Context) (*[]string, error)

	// GetStopWords retrieves the stop words of the index.
	GetStopWords() (*[]string, error)

	// GetStopWordsWithContext retrieves the stop words of the index using the provided context for cancellation.
	GetStopWordsWithContext(ctx context.Context) (*[]string, error)

	// GetSynonyms retrieves the synonyms of the index.
	GetSynonyms() (*map[string][]string, error)

	// GetSynonymsWithContext retrieves the synonyms of the index using the provided context for cancellation.
	GetSynonymsWithContext(ctx context.Context) (*map[string][]string, error)

	// GetFilterableAttributes retrieves the filterable attributes of the index.
	GetFilterableAttributes() (*[]interface{}, error)

	// GetFilterableAttributesWithContext retrieves the filterable attributes of the index using the provided context for cancellation.
	GetFilterableAttributesWithContext(ctx context.Context) (*[]interface{}, error)

	// GetSortableAttributes retrieves the sortable attributes of the index.
	GetSortableAttributes() (*[]string, error)

	// GetSortableAttributesWithContext retrieves the sortable attributes of the index using the provided context for cancellation.
	GetSortableAttributesWithContext(ctx context.Context) (*[]string, error)

	// GetTypoTolerance retrieves the typo tolerance settings of the index.
	GetTypoTolerance() (*TypoTolerance, error)

	// GetTypoToleranceWithContext retrieves the typo tolerance settings of the index using the provided context for cancellation.
	GetTypoToleranceWithContext(ctx context.Context) (*TypoTolerance, error)

	// GetPagination retrieves the pagination settings of the index.
	GetPagination() (*Pagination, error)

	// GetPaginationWithContext retrieves the pagination settings of the index using the provided context for cancellation.
	GetPaginationWithContext(ctx context.Context) (*Pagination, error)

	// GetFaceting retrieves the faceting settings of the index.
	GetFaceting() (*Faceting, error)

	// GetFacetingWithContext retrieves the faceting settings of the index using the provided context for cancellation.
	GetFacetingWithContext(ctx context.Context) (*Faceting, error)

	// GetEmbedders retrieves the embedders of the index.
	GetEmbedders() (map[string]Embedder, error)

	// GetEmbeddersWithContext retrieves the embedders of the index using the provided context for cancellation.
	GetEmbeddersWithContext(ctx context.Context) (map[string]Embedder, error)

	// GetSearchCutoffMs retrieves the search cutoff time in milliseconds.
	GetSearchCutoffMs() (int64, error)

	// GetSearchCutoffMsWithContext retrieves the search cutoff time in milliseconds using the provided context for cancellation.
	GetSearchCutoffMsWithContext(ctx context.Context) (int64, error)

	// GetSeparatorTokens returns separators tokens
	// https://www.meilisearch.com/docs/reference/api/settings#get-separator-tokens
	GetSeparatorTokens() ([]string, error)

	// GetSeparatorTokensWithContext returns separator tokens and support parent context
	// https://www.meilisearch.com/docs/reference/api/settings#get-separator-tokens
	GetSeparatorTokensWithContext(ctx context.Context) ([]string, error)

	// GetNonSeparatorTokens returns non-separator tokens
	// https://www.meilisearch.com/docs/reference/api/settings#get-non-separator-tokens
	GetNonSeparatorTokens() ([]string, error)

	// GetNonSeparatorTokensWithContext returns non-separator tokens and support parent context
	// https://www.meilisearch.com/docs/reference/api/settings#get-non-separator-tokens
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
	// https://www.meilisearch.com/docs/reference/api/settings#get-dictionary
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
	// https://www.meilisearch.com/docs/reference/api/settings#get-dictionary
	GetDictionaryWithContext(ctx context.Context) ([]string, error)

	// GetProximityPrecision returns ProximityPrecision configuration value
	// https://www.meilisearch.com/docs/reference/api/settings#get-proximity-precision-settings
	GetProximityPrecision() (ProximityPrecisionType, error)

	// GetProximityPrecisionWithContext returns ProximityPrecision configuration value and support parent context
	// https://www.meilisearch.com/docs/reference/api/settings#get-proximity-precision-settings
	GetProximityPrecisionWithContext(ctx context.Context) (ProximityPrecisionType, error)

	// GetLocalizedAttributes get the localized attributes settings of an index
	// https://www.meilisearch.com/docs/reference/api/settings#get-localized-attributes-settings
	GetLocalizedAttributes() ([]*LocalizedAttributes, error)

	// GetLocalizedAttributesWithContext get the localized attributes settings of an index using the provided context for cancellation
	// https://www.meilisearch.com/docs/reference/api/settings#get-localized-attributes-settings
	GetLocalizedAttributesWithContext(ctx context.Context) ([]*LocalizedAttributes, error)

	// GetPrefixSearch retrieves the prefix search setting of the index.
	GetPrefixSearch() (*string, error)

	// GetPrefixSearchWithContext retrieves the prefix search setting of the index using the provided context for cancellation.
	GetPrefixSearchWithContext(ctx context.Context) (*string, error)

	// GetFacetSearch retrieves the facet search setting of the index.
	GetFacetSearch() (bool, error)

	// GetFacetSearchWithContext retrieves the facet search setting of the index using the provided context for cancellation.
	GetFacetSearchWithContext(ctx context.Context) (bool, error)
}
