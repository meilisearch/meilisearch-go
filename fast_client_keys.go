package meilisearch

import "net/http"

type fastClientKeys struct {
	client *FastHTTPClient
}

func newFastClientKeys(client *FastHTTPClient) fastClientKeys {
	return fastClientKeys{client: client}
}

func (c fastClientKeys) Get() (resp *Keys, err error) {
	resp = &Keys{}
	req := internalRawRequest{
		endpoint:            "/keys",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Get",
		apiName:             "Keys",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}
