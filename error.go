package meilisearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// ErrCode are all possible errors found during requests
type ErrCode int

const (
	// ErrCodeUnknown default error code, undefined
	ErrCodeUnknown ErrCode = 0
	// ErrCodeMarshalRequest impossible to serialize request body
	ErrCodeMarshalRequest ErrCode = iota + 1
	// ErrCodeResponseUnmarshalBody impossible deserialize the response body
	ErrCodeResponseUnmarshalBody
	// APIError send by the meilisearch api
	APIError
	// APIErrorWithoutMessage APIError send by the meilisearch api
	APIErrorWithoutMessage
	TimeoutError
	// CommunicationError impossible execute a request
	CommunicationError
	// MaxRetriesExceeded used max retries and exceeded
	MaxRetriesExceeded
)

const (
	rawStringCtx                           = `(path "${method} ${endpoint}" with method "${function}")`
	rawStringMarshalRequest                = `unable to marshal body from request: '${request}'`
	rawStringResponseUnmarshalBody         = `unable to unmarshal body from response: '${response}' status code: ${statusCode}`
	rawStringAPIError                      = `unaccepted status code found: ${statusCode} expected: ${statusCodeExpected}, APIError Message: ${message}, Code: ${code}, Type: ${type}, Link: ${link}`
	rawStringAPIErrorWithoutMessage        = `unaccepted status code found: ${statusCode} expected: ${statusCodeExpected}, APIError Message: ${message}`
	rawStringMeilisearchTimeoutError       = `MeilisearchTimeoutError`
	rawStringMeilisearchCommunicationError = `MeilisearchCommunicationError unable to execute request`
	rawStringMeilisearchMaxRetriesExceeded = "failed to request and max retries exceeded"
)

func (e ErrCode) rawMessage() string {
	switch e {
	case ErrCodeMarshalRequest:
		return rawStringMarshalRequest + " " + rawStringCtx
	case ErrCodeResponseUnmarshalBody:
		return rawStringResponseUnmarshalBody + " " + rawStringCtx
	case APIError:
		return rawStringAPIError + " " + rawStringCtx
	case APIErrorWithoutMessage:
		return rawStringAPIErrorWithoutMessage + " " + rawStringCtx
	case TimeoutError:
		return rawStringMeilisearchTimeoutError + " " + rawStringCtx
	case CommunicationError:
		return rawStringMeilisearchCommunicationError + " " + rawStringCtx
	case MaxRetriesExceeded:
		return rawStringMeilisearchMaxRetriesExceeded + " " + rawStringCtx
	default:
		return rawStringCtx
	}
}

// APIErrCode represents Meilisearch API error codes returned by the Meilisearch server.
type APIErrCode string

const (
	// APIErrCodeAPIKeyAlreadyExists A key with this uid already exists.
	APIErrCodeAPIKeyAlreadyExists APIErrCode = "api_key_already_exists"
	// APIErrCodeAPIKeyNotFound The requested API key could not be found.
	APIErrCodeAPIKeyNotFound APIErrCode = "api_key_not_found"
	// APIErrCodeBadRequest The request is invalid, check the error message for more information.
	APIErrCodeBadRequest APIErrCode = "bad_request"
	// APIErrCodeBatchNotFound The requested batch does not exist. Please ensure that you are using the correct uid.
	APIErrCodeBatchNotFound APIErrCode = "batch_not_found"
	// APIErrCodeDatabaseSizeLimitReached The requested database has reached its maximum size.
	APIErrCodeDatabaseSizeLimitReached APIErrCode = "database_size_limit_reached"
	// APIErrCodeDocumentFieldsLimitReached A document exceeds the maximum limit of 65,536 attributes.
	APIErrCodeDocumentFieldsLimitReached APIErrCode = "document_fields_limit_reached"
	// APIErrCodeDocumentNotFound The requested document can't be retrieved. Either it doesn't exist, or the database was left in an inconsistent state.
	APIErrCodeDocumentNotFound APIErrCode = "document_not_found"
	// APIErrCodeDumpProcessFailed An error occurred during the dump creation process. The task was aborted.
	APIErrCodeDumpProcessFailed APIErrCode = "dump_process_failed"
	// APIErrCodeFacetSearchDisabled The /facet-search route has been queried while the facetSearch index setting is set to false.
	APIErrCodeFacetSearchDisabled APIErrCode = "facet_search_disabled"
	// APIErrCodeFeatureNotEnabled You have tried using an experimental feature without activating it.
	APIErrCodeFeatureNotEnabled APIErrCode = "feature_not_enabled"
	// APIErrCodeImmutableAPIKeyActions The actions field of an API key cannot be modified.
	APIErrCodeImmutableAPIKeyActions APIErrCode = "immutable_api_key_actions"
	// APIErrCodeImmutableAPIKeyCreatedAt The createdAt field of an API key cannot be modified.
	APIErrCodeImmutableAPIKeyCreatedAt APIErrCode = "immutable_api_key_created_at"
	// APIErrCodeImmutableAPIKeyExpiresAt The expiresAt field of an API key cannot be modified.
	APIErrCodeImmutableAPIKeyExpiresAt APIErrCode = "immutable_api_key_expires_at"
	// APIErrCodeImmutableAPIKeyIndexes The indexes field of an API key cannot be modified.
	APIErrCodeImmutableAPIKeyIndexes APIErrCode = "immutable_api_key_indexes"
	// APIErrCodeImmutableAPIKeyKey The key field of an API key cannot be modified.
	APIErrCodeImmutableAPIKeyKey APIErrCode = "immutable_api_key_key"
	// APIErrCodeImmutableAPIKeyUID The uid field of an API key cannot be modified.
	APIErrCodeImmutableAPIKeyUID APIErrCode = "immutable_api_key_uid"
	// APIErrCodeImmutableAPIKeyUpdatedAt The updatedAt field of an API key cannot be modified.
	APIErrCodeImmutableAPIKeyUpdatedAt APIErrCode = "immutable_api_key_updated_at"
	// APIErrCodeImmutableIndexUID The uid field of an index cannot be modified.
	APIErrCodeImmutableIndexUID APIErrCode = "immutable_index_uid"
	// APIErrCodeImmutableIndexUpdatedAt The updatedAt field of an index cannot be modified.
	APIErrCodeImmutableIndexUpdatedAt APIErrCode = "immutable_index_updated_at"
	// APIErrCodeImmutableWebhook You tried to modify a reserved webhook. Reserved webhooks are configured by Meilisearch Cloud and have isEditable set to false. Webhooks created with an instance option are also immutable.
	APIErrCodeImmutableWebhook APIErrCode = "immutable_webhook"
	// APIErrCodeImmutableWebhookUUID You tried to manually set a webhook uuid. Meilisearch automatically generates uuid for webhooks.
	APIErrCodeImmutableWebhookUUID APIErrCode = "immutable_webhook_uuid"
	// APIErrCodeImmutableWebhookIsEditable You tried to manually set a webhook's isEditable field. Meilisearch automatically sets isEditable for all webhooks. Only reserved webhooks have isEditable set to false.
	APIErrCodeImmutableWebhookIsEditable APIErrCode = "immutable_webhook_is_editable"
	// APIErrCodeIndexAlreadyExists An index with this uid already exists, check out our guide on index creation.
	APIErrCodeIndexAlreadyExists APIErrCode = "index_already_exists"
	// APIErrCodeIndexCreationFailed An error occurred while trying to create an index, check out our guide on index creation.
	APIErrCodeIndexCreationFailed APIErrCode = "index_creation_failed"
	// APIErrCodeIndexNotFound An index with this uid was not found, check out our guide on index creation.
	APIErrCodeIndexNotFound APIErrCode = "index_not_found"
	// APIErrCodeIndexPrimaryKeyAlreadyExists The requested index already has a primary key that cannot be changed.
	APIErrCodeIndexPrimaryKeyAlreadyExists APIErrCode = "index_primary_key_already_exists"
	// APIErrCodeIndexPrimaryKeyMultipleCandidatesFound Primary key inference failed because the received documents contain multiple fields ending with id. Use the update index endpoint to manually set a primary key.
	APIErrCodeIndexPrimaryKeyMultipleCandidatesFound APIErrCode = "index_primary_key_multiple_candidates_found"
	// APIErrCodeInternal Meilisearch experienced an internal error. Check the error message, and open an issue if necessary.
	APIErrCodeInternal APIErrCode = "internal"
	// APIErrCodeInvalidAPIKey The requested resources are protected with an API key. The provided API key is invalid. Read more about it in our security tutorial.
	APIErrCodeInvalidAPIKey APIErrCode = "invalid_api_key"
	// APIErrCodeInvalidAPIKeyActions The actions field for the provided API key resource is invalid. It should be an array of strings representing action names.
	APIErrCodeInvalidAPIKeyActions APIErrCode = "invalid_api_key_actions"
	// APIErrCodeInvalidAPIKeyDescription The description field for the provided API key resource is invalid. It should either be a string or set to null.
	APIErrCodeInvalidAPIKeyDescription APIErrCode = "invalid_api_key_description"
	// APIErrCodeInvalidAPIKeyExpiresAt The expiresAt field for the provided API key resource is invalid. It should either show a future date or datetime in the RFC 3339 format or be set to null.
	APIErrCodeInvalidAPIKeyExpiresAt APIErrCode = "invalid_api_key_expires_at"
	// APIErrCodeInvalidAPIKeyIndexes The indexes field for the provided API key resource is invalid. It should be an array of strings representing index names.
	APIErrCodeInvalidAPIKeyIndexes APIErrCode = "invalid_api_key_indexes"
	// APIErrCodeInvalidAPIKeyLimit The limit parameter is invalid. It should be an integer.
	APIErrCodeInvalidAPIKeyLimit APIErrCode = "invalid_api_key_limit"
	// APIErrCodeInvalidAPIKeyName The given name is invalid. It should either be a string or set to null.
	APIErrCodeInvalidAPIKeyName APIErrCode = "invalid_api_key_name"
	// APIErrCodeInvalidAPIKeyOffset The offset parameter is invalid. It should be an integer.
	APIErrCodeInvalidAPIKeyOffset APIErrCode = "invalid_api_key_offset"
	// APIErrCodeInvalidAPIKeyUID The given uid is invalid. The uid must follow the uuid v4 format.
	APIErrCodeInvalidAPIKeyUID APIErrCode = "invalid_api_key_uid"
	// APIErrCodeInvalidSearchAttributesToSearchOn The value passed to attributesToSearchOn is invalid. attributesToSearchOn accepts an array of strings indicating document attributes. Attributes given to attributesToSearchOn must be present in the searchableAttributes list.
	APIErrCodeInvalidSearchAttributesToSearchOn APIErrCode = "invalid_search_attributes_to_search_on"
	// APIErrCodeInvalidSearchMedia The value passed to media is not a valid JSON object.
	APIErrCodeInvalidSearchMedia APIErrCode = "invalid_search_media"
	// APIErrCodeInvalidSearchMediaAndVector The search query contains non-null values for both media and vector. These two parameters are mutually exclusive, since media generates vector embeddings via the embedder configured in hybrid.
	APIErrCodeInvalidSearchMediaAndVector APIErrCode = "invalid_search_media_and_vector"
	// APIErrCodeInvalidFilter The provided filter expression is invalid. This may happen if the filter syntax is malformed, uses an unsupported operator, or references an attribute not listed in filterableAttributes.
	APIErrCodeInvalidFilter APIErrCode = "invalid_filter"
	// APIErrCodeInvalidContentType The Content-Type header is not supported by Meilisearch. Currently, Meilisearch only supports JSON, CSV, and NDJSON.
	APIErrCodeInvalidContentType APIErrCode = "invalid_content_type"
	// APIErrCodeInvalidDocumentCSVDelimiter The csvDelimiter parameter is invalid. It should either be a string or a single ASCII character.
	APIErrCodeInvalidDocumentCSVDelimiter APIErrCode = "invalid_document_csv_delimiter"
	// APIErrCodeInvalidDocumentID The provided document identifier does not meet the format requirements. A document identifier must be of type integer or string, composed only of alphanumeric characters (a-z A-Z 0-9), hyphens (-), and underscores (_).
	APIErrCodeInvalidDocumentID APIErrCode = "invalid_document_id"
	// APIErrCodeInvalidDocumentFields The fields parameter is invalid. It should be a string.
	APIErrCodeInvalidDocumentFields APIErrCode = "invalid_document_fields"
	// APIErrCodeInvalidDocumentFilter The filter parameter is invalid or the attribute used for filtering is not filterable.
	APIErrCodeInvalidDocumentFilter APIErrCode = "invalid_document_filter"
	// APIErrCodeInvalidDocumentLimit The limit parameter is invalid. It should be an integer.
	APIErrCodeInvalidDocumentLimit APIErrCode = "invalid_document_limit"
	// APIErrCodeInvalidDocumentOffset The offset parameter is invalid. It should be an integer.
	APIErrCodeInvalidDocumentOffset APIErrCode = "invalid_document_offset"
	// APIErrCodeInvalidDocumentSort The sort parameter is invalid or the attribute used for sorting is not sortable.
	APIErrCodeInvalidDocumentSort APIErrCode = "invalid_document_sort"
	// APIErrCodeInvalidDocumentGeoField The provided _geo field of one or more documents is invalid. Meilisearch expects _geo to be an object with lat and lng fields.
	APIErrCodeInvalidDocumentGeoField APIErrCode = "invalid_document_geo_field"
	// APIErrCodeInvalidDocumentGeoJSONField The geojson field in one or more documents is invalid or doesn't match the GeoJSON specification.
	APIErrCodeInvalidDocumentGeoJSONField APIErrCode = "invalid_document_geojson_field"
	// APIErrCodeInvalidExportURL The export target instance URL is invalid or could not be reached.
	APIErrCodeInvalidExportURL APIErrCode = "invalid_export_url"
	// APIErrCodeInvalidExportAPIKey The supplied security key does not have the required permissions to access the target instance.
	APIErrCodeInvalidExportAPIKey APIErrCode = "invalid_export_api_key"
	// APIErrCodeInvalidExportPayloadSize The provided payload size is invalid. The payload size must be a string indicating the maximum payload size in a human-readable format.
	APIErrCodeInvalidExportPayloadSize APIErrCode = "invalid_export_payload_size"
	// APIErrCodeInvalidExportIndexesPatterns The provided index pattern is invalid. The index pattern must be an alphanumeric string, optionally including a wildcard.
	APIErrCodeInvalidExportIndexesPatterns APIErrCode = "invalid_export_indexes_patterns"
	// APIErrCodeInvalidExportIndexFilter The provided index export filter is not a valid filter expression.
	APIErrCodeInvalidExportIndexFilter APIErrCode = "invalid_export_index_filter"
	// APIErrCodeInvalidFacetSearchFacetName The attribute used for the facetName field is either not a string or not defined in the filterableAttributes list.
	APIErrCodeInvalidFacetSearchFacetName APIErrCode = "invalid_facet_search_facet_name"
	// APIErrCodeInvalidFacetSearchFacetQuery The provided value for facetQuery is invalid. It should either be a string or null.
	APIErrCodeInvalidFacetSearchFacetQuery APIErrCode = "invalid_facet_search_facet_query"
	// APIErrCodeInvalidIndexLimit The limit parameter is invalid. It should be an integer.
	APIErrCodeInvalidIndexLimit APIErrCode = "invalid_index_limit"
	// APIErrCodeInvalidIndexOffset The offset parameter is invalid. It should be an integer.
	APIErrCodeInvalidIndexOffset APIErrCode = "invalid_index_offset"
	// APIErrCodeInvalidIndexUID There is an error in the provided index format, check out our guide on index creation.
	APIErrCodeInvalidIndexUID APIErrCode = "invalid_index_uid"
	// APIErrCodeInvalidIndexPrimaryKey The primaryKey field is invalid. It should either be a string or set to null.
	APIErrCodeInvalidIndexPrimaryKey APIErrCode = "invalid_index_primary_key"
	// APIErrCodeInvalidMultiSearchQueryFederated A multi-search query includes federationOptions but the top-level federation object is null or missing.
	APIErrCodeInvalidMultiSearchQueryFederated APIErrCode = "invalid_multi_search_query_federated"
	// APIErrCodeInvalidMultiSearchQueryPagination A multi-search query contains page, hitsPerPage, limit or offset, but the top-level federation object is not null.
	APIErrCodeInvalidMultiSearchQueryPagination APIErrCode = "invalid_multi_search_query_pagination"
	// APIErrCodeInvalidMultiSearchQueryPosition federationOptions.queryPosition is not a positive integer.
	APIErrCodeInvalidMultiSearchQueryPosition APIErrCode = "invalid_multi_search_query_position"
	// APIErrCodeInvalidMultiSearchWeight A multi-search query contains a negative value for federated.weight.
	APIErrCodeInvalidMultiSearchWeight APIErrCode = "invalid_multi_search_weight"
	// APIErrCodeInvalidMultiSearchQueriesRankingRules Two or more queries in a multi-search request have incompatible results.
	APIErrCodeInvalidMultiSearchQueriesRankingRules APIErrCode = "invalid_multi_search_queries_ranking_rules"
	// APIErrCodeInvalidMultiSearchFacets federation.facetsByIndex.<INDEX_NAME> contains a value that is not in the filterable attributes list.
	APIErrCodeInvalidMultiSearchFacets APIErrCode = "invalid_multi_search_facets"
	// APIErrCodeInvalidMultiSearchSortFacetValuesBy federation.mergeFacets.sortFacetValuesBy is not a string or doesn't have one of the allowed values.
	APIErrCodeInvalidMultiSearchSortFacetValuesBy APIErrCode = "invalid_multi_search_sort_facet_values_by"
	// APIErrCodeInvalidMultiSearchQueryFacets A query in the queries array contains facets when federation is present and non-null.
	APIErrCodeInvalidMultiSearchQueryFacets APIErrCode = "invalid_multi_search_query_facets"
	// APIErrCodeInvalidMultiSearchMergeFacets federation.mergeFacets is not an object or contains unexpected fields.
	APIErrCodeInvalidMultiSearchMergeFacets APIErrCode = "invalid_multi_search_merge_facets"
	// APIErrCodeInvalidMultiSearchMaxValuesPerFacet federation.mergeFacets.maxValuesPerFacet is not a positive integer.
	APIErrCodeInvalidMultiSearchMaxValuesPerFacet APIErrCode = "invalid_multi_search_max_values_per_facet"
	// APIErrCodeInvalidMultiSearchFacetOrder Two or more indexes have a different faceting.sortFacetValuesBy for the same requested facet.
	APIErrCodeInvalidMultiSearchFacetOrder APIErrCode = "invalid_multi_search_facet_order"
	// APIErrCodeInvalidMultiSearchFacetsByIndex facetsByIndex is not an object or contains unknown fields.
	APIErrCodeInvalidMultiSearchFacetsByIndex APIErrCode = "invalid_multi_search_facets_by_index"
	// APIErrCodeInvalidMultiSearchRemote federationOptions.remote is not network.self and is not a key in network.remotes.
	APIErrCodeInvalidMultiSearchRemote APIErrCode = "invalid_multi_search_remote"
	// APIErrCodeInvalidNetworkSelf The network object contains a self that is not a string or null.
	APIErrCodeInvalidNetworkSelf APIErrCode = "invalid_network_self"
	// APIErrCodeInvalidNetworkRemotes The network object contains a remotes that is not an object or null.
	APIErrCodeInvalidNetworkRemotes APIErrCode = "invalid_network_remotes"
	// APIErrCodeInvalidNetworkURL One of the remotes in the network object contains a url that is not a string.
	APIErrCodeInvalidNetworkURL APIErrCode = "invalid_network_url"
	// APIErrCodeInvalidNetworkSearchAPIKey One of the remotes in the network object contains a searchAPIKey that is not a string or null.
	APIErrCodeInvalidNetworkSearchAPIKey APIErrCode = "invalid_network_search_api_key"
	// APIErrCodeInvalidSearchAttributesToCrop The attributesToCrop parameter is invalid. It should be an array of strings, a string, or set to null.
	APIErrCodeInvalidSearchAttributesToCrop APIErrCode = "invalid_search_attributes_to_crop"
	// APIErrCodeInvalidSearchAttributesToHighlight The attributesToHighlight parameter is invalid. It should be an array of strings, a string, or set to null.
	APIErrCodeInvalidSearchAttributesToHighlight APIErrCode = "invalid_search_attributes_to_highlight"
	// APIErrCodeInvalidSearchAttributesToRetrieve The attributesToRetrieve parameter is invalid. It should be an array of strings, a string, or set to null.
	APIErrCodeInvalidSearchAttributesToRetrieve APIErrCode = "invalid_search_attributes_to_retrieve"
	// APIErrCodeInvalidSearchCropLength The cropLength parameter is invalid. It should be an integer.
	APIErrCodeInvalidSearchCropLength APIErrCode = "invalid_search_crop_length"
	// APIErrCodeInvalidSearchCropMarker The cropMarker parameter is invalid. It should be a string or set to null.
	APIErrCodeInvalidSearchCropMarker APIErrCode = "invalid_search_crop_marker"
	// APIErrCodeInvalidSearchEmbedder embedder is invalid. It should be a string corresponding to the name of a configured embedder.
	APIErrCodeInvalidSearchEmbedder APIErrCode = "invalid_search_embedder"
	// APIErrCodeInvalidSearchFacets The facets parameter is invalid or the attribute used for faceting is not filterable.
	APIErrCodeInvalidSearchFacets APIErrCode = "invalid_search_facets"
	// APIErrCodeInvalidSearchFilter The syntax for the filter parameter is invalid or the attribute used for filtering is not filterable.
	APIErrCodeInvalidSearchFilter APIErrCode = "invalid_search_filter"
	// APIErrCodeInvalidSearchHighlightPostTag The highlightPostTag parameter is invalid. It should be a string.
	APIErrCodeInvalidSearchHighlightPostTag APIErrCode = "invalid_search_highlight_post_tag"
	// APIErrCodeInvalidSearchHighlightPreTag The highlightPreTag parameter is invalid. It should be a string.
	APIErrCodeInvalidSearchHighlightPreTag APIErrCode = "invalid_search_highlight_pre_tag"
	// APIErrCodeInvalidSearchHitsPerPage The hitsPerPage parameter is invalid. It should be an integer.
	APIErrCodeInvalidSearchHitsPerPage APIErrCode = "invalid_search_hits_per_page"
	// APIErrCodeInvalidSearchHybridQuery The hybrid parameter is neither null nor an object, or it is an object with unknown keys.
	APIErrCodeInvalidSearchHybridQuery APIErrCode = "invalid_search_hybrid_query"
	// APIErrCodeInvalidSearchLimit The limit parameter is invalid. It should be an integer.
	APIErrCodeInvalidSearchLimit APIErrCode = "invalid_search_limit"
	// APIErrCodeInvalidSearchLocales The locales parameter is invalid.
	APIErrCodeInvalidSearchLocales APIErrCode = "invalid_search_locales"
	// APIErrCodeInvalidSettingsEmbedder The embedders index setting value is invalid.
	APIErrCodeInvalidSettingsEmbedder APIErrCode = "invalid_settings_embedder"
	// APIErrCodeInvalidSettingsFacetSearch The facetSearch index setting value is invalid.
	APIErrCodeInvalidSettingsFacetSearch APIErrCode = "invalid_settings_facet_search"
	// APIErrCodeInvalidSettingsLocalizedAttributes The localizedAttributes index setting value is invalid.
	APIErrCodeInvalidSettingsLocalizedAttributes APIErrCode = "invalid_settings_localized_attributes"
	// APIErrCodeInvalidSearchMatchingStrategy The matchingStrategy parameter is invalid. It should either be set to last or all.
	APIErrCodeInvalidSearchMatchingStrategy APIErrCode = "invalid_search_matching_strategy"
	// APIErrCodeInvalidSearchOffset The offset parameter is invalid. It should be an integer.
	APIErrCodeInvalidSearchOffset APIErrCode = "invalid_search_offset"
	// APIErrCodeInvalidSettingsPrefixSearch The prefixSearch index setting value is invalid.
	APIErrCodeInvalidSettingsPrefixSearch APIErrCode = "invalid_settings_prefix_search"
	// APIErrCodeInvalidSearchPage The page parameter is invalid. It should be an integer.
	APIErrCodeInvalidSearchPage APIErrCode = "invalid_search_page"
	// APIErrCodeInvalidSearchQ The q parameter is invalid. It should be a string or set to null.
	APIErrCodeInvalidSearchQ APIErrCode = "invalid_search_q"
	// APIErrCodeInvalidSearchRankingScoreThreshold The rankingScoreThreshold in a search or multi-search request is not a number between 0.0 and 1.0.
	APIErrCodeInvalidSearchRankingScoreThreshold APIErrCode = "invalid_search_ranking_score_threshold"
	// APIErrCodeInvalidSearchShowMatchesPosition The showMatchesPosition parameter is invalid. It should either be a boolean or set to null.
	APIErrCodeInvalidSearchShowMatchesPosition APIErrCode = "invalid_search_show_matches_position"
	// APIErrCodeInvalidSearchSort The sort parameter is invalid or the attribute used for sorting is not sortable.
	APIErrCodeInvalidSearchSort APIErrCode = "invalid_search_sort"
	// APIErrCodeInvalidSettingsDisplayedAttributes The value of displayed attributes is invalid. It should be an empty array, an array of strings, or set to null.
	APIErrCodeInvalidSettingsDisplayedAttributes APIErrCode = "invalid_settings_displayed_attributes"
	// APIErrCodeInvalidSettingsDistinctAttribute The value of distinct attributes is invalid. It should be a string or set to null.
	APIErrCodeInvalidSettingsDistinctAttribute APIErrCode = "invalid_settings_distinct_attribute"
	// APIErrCodeInvalidSettingsFacetingSortFacetValuesBy The value provided for the sortFacetValuesBy object is incorrect. The accepted values are alpha or count.
	APIErrCodeInvalidSettingsFacetingSortFacetValuesBy APIErrCode = "invalid_settings_faceting_sort_facet_values_by"
	// APIErrCodeInvalidSettingsFacetingMaxValuesPerFacet The value for the maxValuesPerFacet field is invalid. It should either be an integer or set to null.
	APIErrCodeInvalidSettingsFacetingMaxValuesPerFacet APIErrCode = "invalid_settings_faceting_max_values_per_facet"
	// APIErrCodeInvalidSettingsFilterableAttributes The value of filterable attributes is invalid. It should be an empty array, an array of strings, or set to null.
	APIErrCodeInvalidSettingsFilterableAttributes APIErrCode = "invalid_settings_filterable_attributes"
	// APIErrCodeInvalidSettingsPagination The value for the maxTotalHits field is invalid. It should either be an integer or set to null.
	APIErrCodeInvalidSettingsPagination APIErrCode = "invalid_settings_pagination"
	// APIErrCodeInvalidSettingsRankingRules The settings payload has an invalid ranking rules format.
	APIErrCodeInvalidSettingsRankingRules APIErrCode = "invalid_settings_ranking_rules"
	// APIErrCodeInvalidSettingsSearchableAttributes The value of searchable attributes is invalid. It should be an empty array, an array of strings or set to null.
	APIErrCodeInvalidSettingsSearchableAttributes APIErrCode = "invalid_settings_searchable_attributes"
	// APIErrCodeInvalidSettingsSearchCutoffMS The specified value for searchCutoffMs is invalid. It should be an integer indicating the cutoff in milliseconds.
	APIErrCodeInvalidSettingsSearchCutoffMS APIErrCode = "invalid_settings_search_cutoff_ms"
	// APIErrCodeInvalidSettingsSortableAttributes The value of sortable attributes is invalid. It should be an empty array, an array of strings or set to null.
	APIErrCodeInvalidSettingsSortableAttributes APIErrCode = "invalid_settings_sortable_attributes"
	// APIErrCodeInvalidSettingsStopWords The value of stop words is invalid. It should be an empty array, an array of strings or set to null.
	APIErrCodeInvalidSettingsStopWords APIErrCode = "invalid_settings_stop_words"
	// APIErrCodeInvalidSettingsSynonyms The value of the synonyms is invalid. It should either be an object or set to null.
	APIErrCodeInvalidSettingsSynonyms APIErrCode = "invalid_settings_synonyms"
	// APIErrCodeInvalidSettingsTypoTolerance Typo tolerance field is invalid.
	APIErrCodeInvalidSettingsTypoTolerance APIErrCode = "invalid_settings_typo_tolerance"
	// APIErrCodeInvalidSimilarID The provided target document identifier is invalid.
	APIErrCodeInvalidSimilarID APIErrCode = "invalid_similar_id"
	// APIErrCodeNotFoundSimilarID Meilisearch could not find the target document.
	APIErrCodeNotFoundSimilarID APIErrCode = "not_found_similar_id"
	// APIErrCodeInvalidSimilarAttributesToRetrieve attributesToRetrieve is invalid. It should be an array of strings, a string, or set to null.
	APIErrCodeInvalidSimilarAttributesToRetrieve APIErrCode = "invalid_similar_attributes_to_retrieve"
	// APIErrCodeInvalidSimilarEmbedder embedder is invalid. It should be a string corresponding to the name of a configured embedder.
	APIErrCodeInvalidSimilarEmbedder APIErrCode = "invalid_similar_embedder"
	// APIErrCodeInvalidSimilarFilter filter is invalid or contains a filter expression with a missing or invalid operator.
	APIErrCodeInvalidSimilarFilter APIErrCode = "invalid_similar_filter"
	// APIErrCodeInvalidSimilarLimit limit is invalid. It should be an integer.
	APIErrCodeInvalidSimilarLimit APIErrCode = "invalid_similar_limit"
	// APIErrCodeInvalidSimilarOffset offset is invalid. It should be an integer.
	APIErrCodeInvalidSimilarOffset APIErrCode = "invalid_similar_offset"
	// APIErrCodeInvalidSimilarShowRankingScore ranking_score is invalid. It should be a boolean.
	APIErrCodeInvalidSimilarShowRankingScore APIErrCode = "invalid_similar_score"
	// APIErrCodeInvalidSimilarShowRankingScoreDetails ranking_score_details is invalid. It should be a boolean.
	APIErrCodeInvalidSimilarShowRankingScoreDetails APIErrCode = "invalid_similar_score_details"
	// APIErrCodeInvalidSimilarRankingScoreThreshold The rankingScoreThreshold in a similar documents request is not a number between 0.0 and 1.0.
	APIErrCodeInvalidSimilarRankingScoreThreshold APIErrCode = "invalid_similar_ranking_score_threshold"
	// APIErrCodeInvalidState The database is in an invalid state. Deleting the database and re-indexing should solve the problem.
	APIErrCodeInvalidState APIErrCode = "invalid_state"
	// APIErrCodeInvalidStoreFile The data.ms folder is in an invalid state. Your b file is corrupted or the data.ms folder has been replaced by a file.
	APIErrCodeInvalidStoreFile APIErrCode = "invalid_store_file"
	// APIErrCodeInvalidSwapDuplicateIndexFound The indexes used in the indexes array for a swap index request have been declared multiple times. You must declare each index only once.
	APIErrCodeInvalidSwapDuplicateIndexFound APIErrCode = "invalid_swap_duplicate_index_found"
	// APIErrCodeInvalidSwapIndexes The payload doesn't contain exactly two index uids for a swap operation or contains an invalid index name.
	APIErrCodeInvalidSwapIndexes APIErrCode = "invalid_swap_indexes"
	// APIErrCodeInvalidTaskAfterEnqueuedAt The afterEnqueuedAt query parameter is invalid.
	APIErrCodeInvalidTaskAfterEnqueuedAt APIErrCode = "invalid_task_after_enqueued_at"
	// APIErrCodeInvalidTaskAfterFinishedAt The afterFinishedAt query parameter is invalid.
	APIErrCodeInvalidTaskAfterFinishedAt APIErrCode = "invalid_task_after_finished_at"
	// APIErrCodeInvalidTaskAfterStartedAt The afterStartedAt query parameter is invalid.
	APIErrCodeInvalidTaskAfterStartedAt APIErrCode = "invalid_task_after_started_at"
	// APIErrCodeInvalidTaskBeforeEnqueuedAt The beforeEnqueuedAt query parameter is invalid.
	APIErrCodeInvalidTaskBeforeEnqueuedAt APIErrCode = "invalid_task_before_enqueued_at"
	// APIErrCodeInvalidTaskBeforeFinishedAt The beforeFinishedAt query parameter is invalid.
	APIErrCodeInvalidTaskBeforeFinishedAt APIErrCode = "invalid_task_before_finished_at"
	// APIErrCodeInvalidTaskBeforeStartedAt The beforeStartedAt query parameter is invalid.
	APIErrCodeInvalidTaskBeforeStartedAt APIErrCode = "invalid_task_before_started_at"
	// APIErrCodeInvalidTaskCanceledBy The canceledBy query parameter is invalid. It should be an integer. Multiple uids should be separated by commas.
	APIErrCodeInvalidTaskCanceledBy APIErrCode = "invalid_task_canceled_by"
	// APIErrCodeInvalidTaskIndexUIDs The indexUids query parameter contains an invalid index uid.
	APIErrCodeInvalidTaskIndexUIDs APIErrCode = "invalid_task_index_uids"
	// APIErrCodeInvalidTaskLimit The limit parameter is invalid. It must be an integer.
	APIErrCodeInvalidTaskLimit APIErrCode = "invalid_task_limit"
	// APIErrCodeInvalidTaskStatuses The requested task status is invalid. Please use one of the possible values.
	APIErrCodeInvalidTaskStatuses APIErrCode = "invalid_task_statuses"
	// APIErrCodeInvalidTaskTypes The requested task type is invalid. Please use one of the possible values.
	APIErrCodeInvalidTaskTypes APIErrCode = "invalid_task_types"
	// APIErrCodeInvalidTaskUIDs The uids query parameter is invalid.
	APIErrCodeInvalidTaskUIDs APIErrCode = "invalid_task_uids"
	// APIErrCodeInvalidWebhooks The create webhook request did not contain a valid JSON payload. Meilisearch also returns this error when you try to create more than 20 webhooks.
	APIErrCodeInvalidWebhooks APIErrCode = "invalid_webhooks"
	// APIErrCodeInvalidWebhookURL The provided webhook URL isn’t a valid JSON string, is null, is missing, or its value cannot be parsed as a valid URL.
	APIErrCodeInvalidWebhookURL APIErrCode = "invalid_webhook_url"
	// APIErrCodeInvalidWebhookHeaders The provided webhook headers field is not a JSON object or not a valid HTTP header. Meilisearch also returns this error if you set more than 200 header fields for a single webhook.
	APIErrCodeInvalidWebhookHeaders APIErrCode = "invalid_webhook_headers"
	// APIErrCodeInvalidWebhookUUID The provided webhook uuid is not a valid uuid v4 value.
	APIErrCodeInvalidWebhookUUID APIErrCode = "invalid_webhook_uuid"
	// APIErrCodeIOError This error generally occurs when the host system has no space left on the device or when the database doesn't have read or write access.
	APIErrCodeIOError APIErrCode = "io_error"
	// APIErrCodeIndexPrimaryKeyNoCandidateFound Primary key inference failed as the received documents do not contain any fields ending with id. Manually designate the primary key, or add some field ending with id to your documents.
	APIErrCodeIndexPrimaryKeyNoCandidateFound APIErrCode = "index_primary_key_no_candidate_found"
	// APIErrCodeMalformedPayload The Content-Type header does not match the request body payload format or the format is invalid.
	APIErrCodeMalformedPayload APIErrCode = "malformed_payload"
	// APIErrCodeMissingAPIKeyActions The actions field is missing from payload.
	APIErrCodeMissingAPIKeyActions APIErrCode = "missing_api_key_actions"
	// APIErrCodeMissingAPIKeyExpiresAt The expiresAt field is missing from payload.
	APIErrCodeMissingAPIKeyExpiresAt APIErrCode = "missing_api_key_expires_at"
	// APIErrCodeMissingAPIKeyIndexes The indexes field is missing from payload.
	APIErrCodeMissingAPIKeyIndexes APIErrCode = "missing_api_key_indexes"
	// APIErrCodeMissingAuthorizationHeader This error happens if the requested resources are protected with an API key that was not provided in the request header.
	APIErrCodeMissingAuthorizationHeader APIErrCode = "missing_authorization_header"
	// APIErrCodeMissingContentType The payload does not contain a Content-Type header. Currently, Meilisearch only supports JSON, CSV, and NDJSON.
	APIErrCodeMissingContentType APIErrCode = "missing_content_type"
	// APIErrCodeMissingDocumentFilter This payload is missing the filter field.
	APIErrCodeMissingDocumentFilter APIErrCode = "missing_document_filter"
	// APIErrCodeMissingDocumentID A document does not contain any value for the required primary key, and is thus invalid. Check documents in the current addition for the invalid ones.
	APIErrCodeMissingDocumentID APIErrCode = "missing_document_id"
	// APIErrCodeMissingIndexUID The payload is missing the uid field.
	APIErrCodeMissingIndexUID APIErrCode = "missing_index_uid"
	// APIErrCodeMissingFacetSearchFacetName The facetName parameter is required.
	APIErrCodeMissingFacetSearchFacetName APIErrCode = "missing_facet_search_facet_name"
	// APIErrCodeMissingMasterKey You need to set a master key before you can access the /keys route.
	APIErrCodeMissingMasterKey APIErrCode = "missing_master_key"
	// APIErrCodeMissingNetworkURL One of the remotes in the network object does not contain the url field.
	APIErrCodeMissingNetworkURL APIErrCode = "missing_network_url"
	// APIErrCodeMissingPayload The Content-Type header was specified, but no request body was sent to the server or the request body is empty.
	APIErrCodeMissingPayload APIErrCode = "missing_payload"
	// APIErrCodeMissingSwapIndexes The index swap payload is missing the indexes object.
	APIErrCodeMissingSwapIndexes APIErrCode = "missing_swap_indexes"
	// APIErrCodeMissingTaskFilters The cancel tasks and delete tasks endpoints require one of the available query parameters.
	APIErrCodeMissingTaskFilters APIErrCode = "missing_task_filters"
	// APIErrCodeNoSpaceLeftOnDevice The host system partition reaches its maximum capacity and can no longer accept writes or the tasks queue reaches its limit.
	APIErrCodeNoSpaceLeftOnDevice APIErrCode = "no_space_left_on_device"
	// APIErrCodeNotFound The requested resources could not be found.
	APIErrCodeNotFound APIErrCode = "not_found"
	// APIErrCodePayloadTooLarge The payload sent to the server was too large.
	APIErrCodePayloadTooLarge APIErrCode = "payload_too_large"
	// APIErrCodeTaskNotFound The requested task does not exist. Please ensure that you are using the correct uid.
	APIErrCodeTaskNotFound APIErrCode = "task_not_found"
	// APIErrCodeTooManyOpenFiles Indexing a large batch of documents can result in Meilisearch opening too many file descriptors.
	APIErrCodeTooManyOpenFiles APIErrCode = "too_many_open_files"
	// APIErrCodeTooManySearchRequests You have reached the limit of concurrent search requests.
	APIErrCodeTooManySearchRequests APIErrCode = "too_many_search_requests"
	// APIErrCodeUnretrievableDocument The document exists in store, but there was an error retrieving it. This probably comes from an inconsistent state in the database.
	APIErrCodeUnretrievableDocument APIErrCode = "unretrievable_document"
	// APIErrCodeVectorEmbeddingError Error while generating embeddings.
	APIErrCodeVectorEmbeddingError APIErrCode = "vector_embedding_error"
	// APIErrCodeRemoteBadResponse The remote instance answered with a response that this instance could not use as a federated search response.
	APIErrCodeRemoteBadResponse APIErrCode = "remote_bad_response"
	// APIErrCodeRemoteBadRequest The remote instance answered with 400 BAD REQUEST.
	APIErrCodeRemoteBadRequest APIErrCode = "remote_bad_request"
	// APIErrCodeRemoteCouldNotSendRequest There was an error while sending the remote federated search request.
	APIErrCodeRemoteCouldNotSendRequest APIErrCode = "remote_could_not_send_request"
	// APIErrCodeRemoteInvalidAPIKey The remote instance answered with 403 FORBIDDEN or 401 UNAUTHORIZED to this instance’s request.
	APIErrCodeRemoteInvalidAPIKey APIErrCode = "remote_invalid_api_key"
	// APIErrCodeRemoteRemoteError The remote instance answered with 500 INTERNAL ERROR.
	APIErrCodeRemoteRemoteError APIErrCode = "remote_remote_error"
	// APIErrCodeRemoteTimeout The proxy did not answer in the allocated time.
	APIErrCodeRemoteTimeout APIErrCode = "remote_timeout"
	// APIErrCodeWebhookNotFound The provided webhook uuid does not correspond to any configured webhooks in the instance.
	APIErrCodeWebhookNotFound APIErrCode = "webhook_not_found"
)

// APIErrorDetails represents the details of an error returned by Meilisearch API.
type APIErrorDetails struct {
	Message string     `json:"message"`
	Code    APIErrCode `json:"code"`
	Type    string     `json:"type"`
	Link    string     `json:"link"`
}

type meilisearchApiError = APIErrorDetails

// Error is the internal error structure that all exposed method use.
// So ALL errors returned by this library can be cast to this struct (as a pointer)
type Error struct {
	// Endpoint is the path of the request (host is not in)
	Endpoint string

	// Method is the HTTP verb of the request
	Method string

	// Function name used
	Function string

	// RequestToString is the raw request into string ('empty request' if not present)
	RequestToString string

	// RequestToString is the raw request into string ('empty response' if not present)
	ResponseToString string

	// Error info from meilisearch api
	// Message is the raw request into string ('empty meilisearch message' if not present)
	APIError APIErrorDetails

	// StatusCode of the request
	StatusCode int

	// StatusCode expected by the endpoint to be considered as a success
	StatusCodeExpected []int

	rawMessage string

	// OriginError is the origin error that produce the current Error. It can be nil in case of a bad status code.
	OriginError error

	// ErrCode is the internal error code that represent the different step when executing a request that can produce
	// an error.
	ErrCode ErrCode

	encoder
}

// Error return a well human formatted message.
func (e *Error) Error() string {
	message := namedSprintf(e.rawMessage, map[string]interface{}{
		"endpoint":           e.Endpoint,
		"method":             e.Method,
		"function":           e.Function,
		"request":            e.RequestToString,
		"response":           e.ResponseToString,
		"statusCodeExpected": e.StatusCodeExpected,
		"statusCode":         e.StatusCode,
		"message":            e.APIError.Message,
		"code":               e.APIError.Code,
		"type":               e.APIError.Type,
		"link":               e.APIError.Link,
	})
	if e.OriginError != nil {
		return fmt.Sprintf("%s: %s", message, e.OriginError.Error())
	}

	return message
}

// Unwrap returns the underlying error if one exists.
func (e *Error) Unwrap() error {
	return e.OriginError
}

// HasCode returns true if the Meilisearch API error code matches the given code.
func (e *Error) HasCode(code APIErrCode) bool {
	return e.APIError.Code == code
}

// WithErrCode add an error code to an error
func (e *Error) WithErrCode(err ErrCode, errs ...error) *Error {
	if errs != nil {
		e.OriginError = errs[0]
	}

	e.rawMessage = err.rawMessage()
	e.ErrCode = err
	return e
}

// ErrorBody add a body to an error
func (e *Error) ErrorBody(body []byte) {
	msg := meilisearchApiError{}

	if e.encoder != nil {
		err := e.Decode(body, &msg)
		if err == nil {
			e.APIError.Message = msg.Message
			e.APIError.Code = msg.Code
			e.APIError.Type = msg.Type
			e.APIError.Link = msg.Link
		}
		return
	}

	e.ResponseToString = string(body)
	err := json.Unmarshal(body, &msg)
	if err == nil {
		e.APIError.Message = msg.Message
		e.APIError.Code = msg.Code
		e.APIError.Type = msg.Type
		e.APIError.Link = msg.Link
	}
}

// VersionErrorHintMessage a hint to the error message if it may come from a version incompatibility with meilisearch
func VersionErrorHintMessage(err error, req *internalRequest) error {
	return fmt.Errorf("%w. Hint: It might not be working because you're not up to date with the "+
		"Meilisearch version that %s call requires", err, req.functionName)
}

func namedSprintf(format string, params map[string]interface{}) string {
	for key, val := range params {
		format = strings.ReplaceAll(format, "${"+key+"}", fmt.Sprintf("%v", val))
	}
	return format
}

// General errors
var (
	ErrInvalidRequestMethod          = errors.New("request body is not expected for GET and HEAD requests")
	ErrRequestBodyWithoutContentType = errors.New("request body without Content-Type is not allowed")
	ErrNoSearchRequest               = errors.New("no search request provided")
	ErrNoFacetSearchRequest          = errors.New("no search facet request provided")
	ErrConnectingFailed              = errors.New("meilisearch is not connected")
	ErrMeilisearchNotAvailable       = errors.New("meilisearch service is not available")
)
