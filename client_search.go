package meilisearch

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type clientSearch struct {
	client  *Client
	indexID string
}

func newClientSearch(client *Client, indexId string) clientSearch {
	return clientSearch{client: client, indexID: indexId}
}

func (c clientSearch) Search(request SearchRequest) (*SearchResponse, error) {

	resp := &SearchResponse{}

	values := url.Values{}

	values.Add("q", request.Query)
	values.Add("filters", request.Filters)
	values.Add("offset", strconv.FormatInt(request.Offset, 10))
	values.Add("limit", strconv.FormatInt(request.Limit, 10))
	values.Add("timeoutMs", strconv.FormatInt(request.TimeoutMs, 10))
	values.Add("cropLength", strconv.FormatInt(request.CropLength, 10))
	values.Add("attributesToRetrieve", strings.Join(request.AttributesToRetrieve, ","))
	values.Add("attributesToSearchIn", strings.Join(request.AttributesToSearchIn, ","))
	values.Add("attributesToCrop", strings.Join(request.AttributesToCrop, ","))
	values.Add("attributesToHighlight", strings.Join(request.AttributesToHighlight, ","))
	values.Add("matches", strconv.FormatBool(request.Matches))

	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/search?" + values.Encode(),
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: nil,
		functionName:        "Search",
		apiName:             "Search",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}
