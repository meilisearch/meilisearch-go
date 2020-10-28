package meilisearch

import (
	"net/http"
)

type clientSearch struct {
	client   *Client
	indexUID string
}

func newClientSearch(client *Client, indexUID string) clientSearch {
	return clientSearch{client: client, indexUID: indexUID}
}

func (c clientSearch) Search(request SearchRequest) (*SearchResponse, error) {

	resp := &SearchResponse{}

	searchPostRequestParams := map[string]interface{}{}

	if request.Limit == 0 {
		request.Limit = 20
	}

	if !request.PlaceholderSearch {
		searchPostRequestParams["q"] = request.Query
	}
	if request.Filters != "" {
		searchPostRequestParams["filters"] = request.Filters
	}
	if request.Offset != 0 {
		searchPostRequestParams["offset"] = request.Offset
	}
	if request.Limit != 20 {
		searchPostRequestParams["limit"] = request.Limit
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
	if request.Matches {
		searchPostRequestParams["matches"] = request.Matches
	}
	if len(request.FacetsDistribution) != 0 {
		searchPostRequestParams["facetsDistribution"] = request.FacetsDistribution
	}
	if request.FacetFilters != nil {
		searchPostRequestParams["facetFilters"] = request.FacetFilters
	}

	req := internalRequest{
		endpoint:            "/indexes/" + c.indexUID + "/search",
		method:              http.MethodPost,
		withRequest:         searchPostRequestParams,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Search",
		apiName:             "Search",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c clientSearch) IndexID() string {
	return c.indexUID
}

func (c clientSearch) Client() ClientInterface {
	return c.client
}
