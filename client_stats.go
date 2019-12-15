package meilisearch

import "net/http"

type clientStats struct {
	client *Client
}

func newClientStats(client *Client) clientStats {
	return clientStats{client: client}
}

func (c clientStats) Get(indexId string) (resp *Stats, err error) {
	resp = &Stats{}
	req := internalRequest{
		endpoint:            "/stats/" + indexId,
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Get",
		apiName:             "Stats",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c clientStats) List() (resp []Stats, err error) {
	resp = []Stats{}
	req := internalRequest{
		endpoint:            "/stats",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Get",
		apiName:             "Stats",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}
