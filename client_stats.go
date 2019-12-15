package meilisearch

import "net/http"

type clientStats struct {
	client  *Client
	indexID string
}

func newClientStats(client *Client, indexId string) clientStats {
	return clientStats{client: client, indexID: indexId}
}

func (c clientStats) Get() (resp *Stats, err error) {
	resp = &Stats{}
	req := internalRequest{
		endpoint:            "/stats/" + c.indexID,
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

func (c clientStats) IndexId() string {
	return c.indexID
}

func (c clientStats) Client() *Client {
	return c.client
}
