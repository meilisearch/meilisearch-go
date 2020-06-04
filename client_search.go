package meilisearch

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

	values := url.Values{}

	if request.Limit == 0 {
		request.Limit = 20
	}

	values.Add("q", request.Query)
	if request.Filters != "" {
		values.Add("filters", request.Filters)
	}
	if request.Offset != 0 {
		values.Add("offset", strconv.FormatInt(request.Offset, 10))
	}
	if request.Limit != 20 {
		values.Add("limit", strconv.FormatInt(request.Limit, 10))
	}
	if request.CropLength != 0 {
		values.Add("cropLength", strconv.FormatInt(request.CropLength, 10))
	}
	if len(request.AttributesToRetrieve) != 0 {
		values.Add("attributesToRetrieve", strings.Join(request.AttributesToRetrieve, ","))
	}
	if len(request.AttributesToCrop) != 0 {
		values.Add("attributesToCrop", strings.Join(request.AttributesToCrop, ","))
	}
	if len(request.AttributesToHighlight) != 0 {
		values.Add("attributesToHighlight", strings.Join(request.AttributesToHighlight, ","))
	}
	if request.Matches {
		values.Add("matches", strconv.FormatBool(request.Matches))
	}

	req := internalRequest{
		endpoint:            "/indexes/" + c.indexUID + "/search?" + values.Encode(),
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
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

func (c clientSearch) Client() *Client {
	return c.client
}
