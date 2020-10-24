package meilisearch

import "net/http"

type fastClientStats struct {
	client *FastHTTPClient
}

func newFastClientStats(client *FastHTTPClient) fastClientStats {
	return fastClientStats{client: client}
}

func (c fastClientStats) Get(indexUID string) (resp *StatsIndex, err error) {
	resp = &StatsIndex{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + indexUID + "/stats",
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

func (c fastClientStats) GetAll() (resp *Stats, err error) {
	resp = &Stats{}
	req := internalRawRequest{
		endpoint:            "/stats",
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
