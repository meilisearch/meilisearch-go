package meilisearch

import (
	"encoding/json"
	"errors"
	"net/http"
)

var ErrNoFacetSearchRequest = errors.New("no search facet request provided")

func (i Index) FacetSearch(request *FacetSearchRequest) (*json.RawMessage, error) {
	if request == nil {
		return nil, ErrNoFacetSearchRequest
	}

	searchPostRequestParams := FacetSearchPostRequestParams(request)

	resp := &json.RawMessage{}

	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/facet-search",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         searchPostRequestParams,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "FacetSearch",
	}

	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func FacetSearchPostRequestParams(request *FacetSearchRequest) map[string]interface{} {
	params := make(map[string]interface{}, 22)

	if request.Q != "" {
		params["q"] = request.Q
	}
	if request.FacetName != "" {
		params["facetName"] = request.FacetName
	}
	if request.FacetQuery != "" {
		params["facetQuery"] = request.FacetQuery
	}
	if request.Filter != "" {
		params["filter"] = request.Filter
	}
	if request.MatchingStrategy != "" {
		params["matchingStrategy"] = request.MatchingStrategy
	}
	if len(request.AttributesToSearchOn) != 0 {
		params["attributesToSearchOn"] = request.AttributesToSearchOn
	}

	return params
}
