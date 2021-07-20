package meilisearch

import (
	"net/http"
)

// This constant contains the default values assigned by meilisearch to the limit in search parameters
//
// Documentation: https://docs.meilisearch.com/reference/features/search_parameters.html
const (
	DefaultLimit int64 = 20
)

func (i Index) Search(query string, request *SearchRequest) (*SearchResponse, error) {
	resp := &SearchResponse{}

	searchPostRequestParams := map[string]interface{}{}

	if request.Limit == 0 {
		request.Limit = DefaultLimit
	}

	if !request.PlaceholderSearch {
		searchPostRequestParams["q"] = query
	}
	if request.Limit != DefaultLimit {
		searchPostRequestParams["limit"] = request.Limit
	}
	if request.Matches {
		searchPostRequestParams["matches"] = request.Matches
	}
	if request.Filter != nil {
		searchPostRequestParams["filter"] = request.Filter
	}
	if request.Offset != 0 {
		searchPostRequestParams["offset"] = request.Offset
	}
	if request.CropLength != 0 {
		searchPostRequestParams["cropLength"] = request.CropLength
	}
	if len(request.AttributesToRetrieve) != 0 {
		searchPostRequestParams["attributesToRetrieve"] = request.AttributesToRetrieve
	}
	if len(request.AttributesToCrop) != 0 {
		searchPostRequestParams["attributesToCrop"] = request.AttributesToCrop
	}
	if len(request.AttributesToHighlight) != 0 {
		searchPostRequestParams["attributesToHighlight"] = request.AttributesToHighlight
	}
	if len(request.FacetsDistribution) != 0 {
		searchPostRequestParams["facetsDistribution"] = request.FacetsDistribution
	}

	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/search",
		method:              http.MethodPost,
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
