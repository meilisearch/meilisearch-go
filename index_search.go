package meilisearch

import (
	"encoding/json"
	"errors"
	"net/http"
)

// This constant contains the default values assigned by Meilisearch to the limit in search parameters
//
// Documentation: https://www.meilisearch.com/docs/reference/api/search#search-parameters
const (
	DefaultLimit int64 = 20
)

var ErrNoSearchRequest = errors.New("no search request provided")

func (i Index) SearchRaw(query string, request *SearchRequest) (*json.RawMessage, error) {
	if request == nil {
		return nil, ErrNoSearchRequest
	}

	if request.Limit == 0 {
		request.Limit = DefaultLimit
	}

	searchPostRequestParams := searchPostRequestParams(query, request)

	resp := &json.RawMessage{}

	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/search",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         searchPostRequestParams,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "SearchRaw",
	}

	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (i Index) Search(query string, request *SearchRequest) (*SearchResponse, error) {
	if request == nil {
		return nil, ErrNoSearchRequest
	}

	if request.Limit == 0 {
		request.Limit = DefaultLimit
	}
	if request.IndexUID != "" {
		request.IndexUID = ""
	}

	searchPostRequestParams := searchPostRequestParams(query, request)

	resp := &SearchResponse{}

	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/search",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         searchPostRequestParams,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Search",
	}

	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func searchPostRequestParams(query string, request *SearchRequest) map[string]interface{} {
	params := make(map[string]interface{}, 16)

	if !request.PlaceholderSearch {
		params["q"] = query
	}
	if request.IndexUID != "" {
		params["indexUid"] = request.IndexUID
	}
	if request.Limit != DefaultLimit {
		params["limit"] = request.Limit
	}
	if request.ShowMatchesPosition {
		params["showMatchesPosition"] = request.ShowMatchesPosition
	}
	if request.ShowRankingScore {
		params["showRankingScore"] = request.ShowRankingScore
	}
	if request.Filter != nil {
		params["filter"] = request.Filter
	}
	if request.Offset != 0 {
		params["offset"] = request.Offset
	}
	if request.CropLength != 0 {
		params["cropLength"] = request.CropLength
	}
	if request.HitsPerPage != 0 {
		params["hitsPerPage"] = request.HitsPerPage
	}
	if request.Page != 0 {
		params["page"] = request.Page
	}
	if request.CropMarker != "" {
		params["cropMarker"] = request.CropMarker
	}
	if request.HighlightPreTag != "" {
		params["highlightPreTag"] = request.HighlightPreTag
	}
	if request.HighlightPostTag != "" {
		params["highlightPostTag"] = request.HighlightPostTag
	}
	if request.MatchingStrategy != "" {
		params["matchingStrategy"] = request.MatchingStrategy
	}
	if len(request.AttributesToRetrieve) != 0 {
		params["attributesToRetrieve"] = request.AttributesToRetrieve
	}
	if len(request.AttributesToSearchOn) != 0 {
		params["attributesToSearchOn"] = request.AttributesToSearchOn
	}
	if len(request.AttributesToCrop) != 0 {
		params["attributesToCrop"] = request.AttributesToCrop
	}
	if len(request.AttributesToHighlight) != 0 {
		params["attributesToHighlight"] = request.AttributesToHighlight
	}
	if len(request.Facets) != 0 {
		params["facets"] = request.Facets
	}
	if len(request.Sort) != 0 {
		params["sort"] = request.Sort
	}

	if request.Vector != nil && len(request.Vector) > 0 {
		params["vector"] = request.Vector
	}

	return params
}
