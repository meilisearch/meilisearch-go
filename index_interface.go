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
	UpdateIndex(ctx context.Context, params *UpdateIndexRequestParams) (*TaskInfo, error)

	// Delete removes the index identified by the given UID.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/delete-index
	Delete(ctx context.Context, uid string) (bool, error)

	// Compact compacts the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/compact-index
	Compact(ctx context.Context) (*TaskInfo, error)
}

type IndexReader interface {
	// FetchInfo retrieves information about the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/get-index
	FetchInfo(ctx context.Context) (*IndexResult, error)

	// FetchPrimaryKey retrieves the primary key of the index.
	FetchPrimaryKey(ctx context.Context) (*string, error)

	// GetStats retrieves statistical information about the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/stats/get-stats-of-index
	GetStats(ctx context.Context, param *StatsParams) (*StatsIndex, error)
}

type DocumentManager interface {
	DocumentReader

	// AddDocuments adds multiple documents to the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	AddDocuments(ctx context.Context, documentsPtr interface{}, opts *DocumentOptions) (*TaskInfo, error)

	// AddDocumentsInBatches adds documents to the index in batches of specified size.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	AddDocumentsInBatches(ctx context.Context, documentsPtr interface{}, batchSize int, opts *DocumentOptions) ([]TaskInfo, error)

	// AddDocumentsCsv adds documents from a CSV byte array to the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	AddDocumentsCsv(ctx context.Context, documents []byte, options *CsvDocumentsQuery) (*TaskInfo, error)

	// AddDocumentsCsvInBatches adds documents from a CSV byte array to the index in batches of specified size.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	AddDocumentsCsvInBatches(ctx context.Context, documents []byte, batchSize int, options *CsvDocumentsQuery) ([]TaskInfo, error)

	// AddDocumentsCsvFromReaderInBatches adds documents from a CSV reader to the index in batches of specified size.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	AddDocumentsCsvFromReaderInBatches(ctx context.Context, documents io.Reader, batchSize int, options *CsvDocumentsQuery) ([]TaskInfo, error)

	// AddDocumentsCsvFromReader adds documents from a CSV reader to the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	AddDocumentsCsvFromReader(ctx context.Context, documents io.Reader, options *CsvDocumentsQuery) (*TaskInfo, error)

	// AddDocumentsNdjson adds documents from a NDJSON byte array to the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	AddDocumentsNdjson(ctx context.Context, documents []byte, opts *DocumentOptions) (*TaskInfo, error)

	// AddDocumentsNdjsonInBatches adds documents from a NDJSON byte array to the index in batches of specified size.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	AddDocumentsNdjsonInBatches(ctx context.Context, documents []byte, batchSize int, opts *DocumentOptions) ([]TaskInfo, error)

	// AddDocumentsNdjsonFromReader adds documents from a NDJSON reader to the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	AddDocumentsNdjsonFromReader(ctx context.Context, documents io.Reader, opts *DocumentOptions) (*TaskInfo, error)

	// AddDocumentsNdjsonFromReaderInBatches adds documents from a NDJSON reader to the index in batches of specified size.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	AddDocumentsNdjsonFromReaderInBatches(ctx context.Context, documents io.Reader, batchSize int, opts *DocumentOptions) ([]TaskInfo, error)

	// UpdateDocuments updates multiple documents in the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	UpdateDocuments(ctx context.Context, documentsPtr interface{}, opts *DocumentOptions) (*TaskInfo, error)

	// UpdateDocumentsInBatches updates documents in the index in batches of specified size.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	UpdateDocumentsInBatches(ctx context.Context, documentsPtr interface{}, batchSize int, opts *DocumentOptions) ([]TaskInfo, error)

	// UpdateDocumentsCsv updates documents in the index from a CSV byte array.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	UpdateDocumentsCsv(ctx context.Context, documents []byte, options *CsvDocumentsQuery) (*TaskInfo, error)

	// UpdateDocumentsCsvInBatches updates documents in the index from a CSV byte array in batches of specified size.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	UpdateDocumentsCsvInBatches(ctx context.Context, documents []byte, batchsize int, options *CsvDocumentsQuery) ([]TaskInfo, error)

	// UpdateDocumentsNdjson updates documents in the index from a NDJSON byte array.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	UpdateDocumentsNdjson(ctx context.Context, documents []byte, opts *DocumentOptions) (*TaskInfo, error)

	// UpdateDocumentsNdjsonInBatches updates documents in the index from a NDJSON byte array in batches of specified size.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	UpdateDocumentsNdjsonInBatches(ctx context.Context, documents []byte, batchsize int, opts *DocumentOptions) ([]TaskInfo, error)

	// UpdateDocumentsByFunction update documents by using function
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents
	UpdateDocumentsByFunction(ctx context.Context, req *UpdateDocumentByFunctionRequest) (*TaskInfo, error)

	// DeleteDocument deletes a single document from the index by identifier.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/delete-document
	DeleteDocument(ctx context.Context, identifier string, opts *DocumentOptions) (*TaskInfo, error)

	// DeleteDocuments deletes multiple documents from the index by identifiers.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/delete-document
	DeleteDocuments(ctx context.Context, identifiers []string, opts *DocumentOptions) (*TaskInfo, error)

	// DeleteDocumentsByFilter deletes documents from the index by filter.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/delete-documents-by-filter
	DeleteDocumentsByFilter(ctx context.Context, filter interface{}, opts *DocumentOptions) (*TaskInfo, error)

	// DeleteAllDocuments deletes all documents from the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/delete-all-documents
	DeleteAllDocuments(ctx context.Context, opts *DocumentOptions) (*TaskInfo, error)
}

type DocumentReader interface {
	// GetDocument retrieves a single document from the index by identifier.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/get-document
	GetDocument(ctx context.Context, identifier string, request *DocumentQuery, documentPtr interface{}) error

	// GetDocuments retrieves multiple documents from the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/documents/list-documents-with-get
	GetDocuments(ctx context.Context, param *DocumentsQuery, resp *DocumentsResult) error
}

type SearchReader interface {
	// Search performs a search query on the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/search/search-with-get
	Search(ctx context.Context, query string, request *SearchRequest) (*SearchResponse, error)

	// SearchRaw performs a raw search query on the index, returning a JSON response.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/search/search-with-get
	SearchRaw(ctx context.Context, query string, request *SearchRequest) (*json.RawMessage, error)

	// FacetSearch performs a facet search query on the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/facet-search/search-in-facets
	FacetSearch(ctx context.Context, request *FacetSearchRequest) (*json.RawMessage, error)

	// SearchSimilarDocuments performs a search for similar documents.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/similar-documents/get-similar-documents-with-get
	SearchSimilarDocuments(ctx context.Context, param *SimilarDocumentQuery, resp *SimilarDocumentResult) error
}

type SettingsManager interface {
	SettingsReader

	// UpdateSettings updates the settings of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-all-settings
	UpdateSettings(ctx context.Context, request *Settings) (*TaskInfo, error)

	// ResetSettings resets the settings of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-all-settings
	ResetSettings(ctx context.Context) (*TaskInfo, error)

	// UpdateRankingRules updates the ranking rules of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-rankingrules
	UpdateRankingRules(ctx context.Context, request *[]string) (*TaskInfo, error)

	// ResetRankingRules resets the ranking rules of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-rankingrules
	ResetRankingRules(ctx context.Context) (*TaskInfo, error)

	// UpdateDistinctAttribute updates the distinct attribute of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-distinctattribute
	UpdateDistinctAttribute(ctx context.Context, request string) (*TaskInfo, error)

	// ResetDistinctAttribute resets the distinct attribute of the index to default value.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-distinctattribute
	ResetDistinctAttribute(ctx context.Context) (*TaskInfo, error)

	// UpdateSearchableAttributes updates the searchable attributes of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-searchableattributes
	UpdateSearchableAttributes(ctx context.Context, request *[]string) (*TaskInfo, error)

	// ResetSearchableAttributes resets the searchable attributes of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-searchableattributes
	ResetSearchableAttributes(ctx context.Context) (*TaskInfo, error)

	// UpdateDisplayedAttributes updates the displayed attributes of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-displayedattributes
	UpdateDisplayedAttributes(ctx context.Context, request *[]string) (*TaskInfo, error)

	// ResetDisplayedAttributes resets the displayed attributes of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-displayedattributes
	ResetDisplayedAttributes(ctx context.Context) (*TaskInfo, error)

	// UpdateStopWords updates the stop words of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-stopwords
	UpdateStopWords(ctx context.Context, request *[]string) (*TaskInfo, error)

	// ResetStopWords resets the stop words of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-stopwords
	ResetStopWords(ctx context.Context) (*TaskInfo, error)

	// UpdateSynonyms updates the synonyms of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-synonyms
	UpdateSynonyms(ctx context.Context, request *map[string][]string) (*TaskInfo, error)

	// ResetSynonyms resets the synonyms of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-synonyms
	ResetSynonyms(ctx context.Context) (*TaskInfo, error)

	// UpdateFilterableAttributes updates the filterable attributes of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-filterableattributes
	UpdateFilterableAttributes(ctx context.Context, request *[]interface{}) (*TaskInfo, error)

	// ResetFilterableAttributes resets the filterable attributes of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-filterableattributes
	ResetFilterableAttributes(ctx context.Context) (*TaskInfo, error)

	// UpdateSortableAttributes updates the sortable attributes of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-sortableattributes
	UpdateSortableAttributes(ctx context.Context, request *[]string) (*TaskInfo, error)

	// ResetSortableAttributes resets the sortable attributes of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-sortableattributes
	ResetSortableAttributes(ctx context.Context) (*TaskInfo, error)

	// UpdateTypoTolerance updates the typo tolerance settings of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-typotolerance
	UpdateTypoTolerance(ctx context.Context, request *TypoTolerance) (*TaskInfo, error)

	// ResetTypoTolerance resets the typo tolerance settings of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-typotolerance
	ResetTypoTolerance(ctx context.Context) (*TaskInfo, error)

	// UpdatePagination updates the pagination settings of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-pagination
	UpdatePagination(ctx context.Context, request *Pagination) (*TaskInfo, error)

	// ResetPagination resets the pagination settings of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-pagination
	ResetPagination(ctx context.Context) (*TaskInfo, error)

	// UpdateFaceting updates the faceting settings of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-faceting
	UpdateFaceting(ctx context.Context, request *Faceting) (*TaskInfo, error)

	// ResetFaceting resets the faceting settings of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-faceting
	ResetFaceting(ctx context.Context) (*TaskInfo, error)

	// UpdateEmbedders updates the embedders of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-embedders
	UpdateEmbedders(ctx context.Context, request map[string]Embedder) (*TaskInfo, error)

	// ResetEmbedders resets the embedders of the index to default values.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-embedders
	ResetEmbedders(ctx context.Context) (*TaskInfo, error)

	// UpdateSearchCutoffMs updates the search cutoff time in milliseconds.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-searchcutoffms
	UpdateSearchCutoffMs(ctx context.Context, request int64) (*TaskInfo, error)

	// ResetSearchCutoffMs resets the search cutoff time in milliseconds to default value.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-searchcutoffms
	ResetSearchCutoffMs(ctx context.Context) (*TaskInfo, error)

	// UpdateSeparatorTokens update separator tokens
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-separatortokens
	UpdateSeparatorTokens(ctx context.Context, tokens []string) (*TaskInfo, error)

	// ResetSeparatorTokens reset separator tokens
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-separatortokens
	ResetSeparatorTokens(ctx context.Context) (*TaskInfo, error)

	// UpdateNonSeparatorTokens update non-separator tokens
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-nonseparatortokens
	UpdateNonSeparatorTokens(ctx context.Context, tokens []string) (*TaskInfo, error)

	// ResetNonSeparatorTokens reset non-separator tokens
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-nonseparatortokens
	ResetNonSeparatorTokens(ctx context.Context) (*TaskInfo, error)

	// UpdateDictionary update user dictionary
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-dictionary
	UpdateDictionary(ctx context.Context, words []string) (*TaskInfo, error)

	// ResetDictionary reset user dictionary
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-dictionary
	ResetDictionary(ctx context.Context) (*TaskInfo, error)

	// UpdateProximityPrecision set ProximityPrecision value ByWord or ByAttribute
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-proximityprecision
	UpdateProximityPrecision(ctx context.Context, proximityType ProximityPrecisionType) (*TaskInfo, error)

	// ResetProximityPrecision reset ProximityPrecision to default ByWord
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-proximityprecision
	ResetProximityPrecision(ctx context.Context) (*TaskInfo, error)

	// UpdateLocalizedAttributes update the localized attributes settings of an index
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-localizedattributes
	UpdateLocalizedAttributes(ctx context.Context, request []*LocalizedAttributes) (*TaskInfo, error)

	// ResetLocalizedAttributes reset the localized attributes settings
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-localizedattributes
	ResetLocalizedAttributes(ctx context.Context) (*TaskInfo, error)

	// UpdatePrefixSearch updates the prefix search setting of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-prefixsearch
	UpdatePrefixSearch(ctx context.Context, request string) (*TaskInfo, error)

	// ResetPrefixSearch resets the prefix search setting of the index to default value.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-prefixsearch
	ResetPrefixSearch(ctx context.Context) (*TaskInfo, error)

	// UpdateFacetSearch updates the facet search setting of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-facetsearch
	UpdateFacetSearch(ctx context.Context, request bool) (*TaskInfo, error)

	// ResetFacetSearch resets the facet search setting of the index to default value.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-facetsearch
	ResetFacetSearch(ctx context.Context) (*TaskInfo, error)

	// UpdateForeignKeys updates the foreignKeys setting for the index. Set foreignKeys to nil to reset to default.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/update-foreignkeys
	UpdateForeignKeys(ctx context.Context, foreignKeys []ForeignKey) (*TaskInfo, error)

	// ResetForeignKeys resets the foreignKeys setting to its default value.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/reset-foreignkeys
	ResetForeignKeys(ctx context.Context) (*TaskInfo, error)
}

type SettingsReader interface {
	// GetSettings retrieves the settings of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/list-all-settings
	GetSettings(ctx context.Context) (*Settings, error)

	// GetRankingRules retrieves the ranking rules of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-rankingrules
	GetRankingRules(ctx context.Context) (*[]string, error)

	// GetDistinctAttribute retrieves the distinct attribute of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-distinctattribute
	GetDistinctAttribute(ctx context.Context) (*string, error)

	// GetSearchableAttributes retrieves the searchable attributes of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-searchableattributes
	GetSearchableAttributes(ctx context.Context) (*[]string, error)

	// GetDisplayedAttributes retrieves the displayed attributes of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-displayedattributes
	GetDisplayedAttributes(ctx context.Context) (*[]string, error)

	// GetStopWords retrieves the stop words of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-stopwords
	GetStopWords(ctx context.Context) (*[]string, error)

	// GetSynonyms retrieves the synonyms of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-synonyms
	GetSynonyms(ctx context.Context) (*map[string][]string, error)

	// GetFilterableAttributes retrieves the filterable attributes of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-filterableattributes
	GetFilterableAttributes(ctx context.Context) (*[]interface{}, error)

	// GetSortableAttributes retrieves the sortable attributes of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-sortableattributes
	GetSortableAttributes(ctx context.Context) (*[]string, error)

	// GetTypoTolerance retrieves the typo tolerance settings of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-typotolerance
	GetTypoTolerance(ctx context.Context) (*TypoTolerance, error)

	// GetPagination retrieves the pagination settings of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-pagination
	GetPagination(ctx context.Context) (*Pagination, error)

	// GetFaceting retrieves the faceting settings of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-faceting
	GetFaceting(ctx context.Context) (*Faceting, error)

	// GetEmbedders retrieves the embedders of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-embedders
	GetEmbedders(ctx context.Context) (map[string]Embedder, error)

	// GetSearchCutoffMs retrieves the search cutoff time in milliseconds.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-searchcutoffms
	GetSearchCutoffMs(ctx context.Context) (int64, error)

	// GetSeparatorTokens returns separators tokens
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-separatortokens
	GetSeparatorTokens(ctx context.Context) ([]string, error)

	// GetNonSeparatorTokens returns non-separator tokens
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-nonseparatortokens
	GetNonSeparatorTokens(ctx context.Context) ([]string, error)

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
	GetDictionary(ctx context.Context) ([]string, error)

	// GetProximityPrecision returns ProximityPrecision configuration value
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-proximityprecision
	GetProximityPrecision(ctx context.Context) (ProximityPrecisionType, error)

	// GetLocalizedAttributes get the localized attributes settings of an index
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-localizedattributes
	GetLocalizedAttributes(ctx context.Context) ([]*LocalizedAttributes, error)

	// GetPrefixSearch retrieves the prefix search setting of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-prefixsearch
	GetPrefixSearch(ctx context.Context) (*string, error)

	// GetFacetSearch retrieves the facet search setting of the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-facetsearch
	GetFacetSearch(ctx context.Context) (bool, error)

	// GetForeignKeys returns the current value of the foreignKeys setting for the index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/settings/get-foreignkeys
	GetForeignKeys(ctx context.Context) ([]ForeignKey, error)
}
