package meilisearch

import "net/http"

type clientKeys struct {
	client *Client
}

func newClientKeys(client *Client) clientKeys {
	return clientKeys{client: client}
}

func (c clientKeys) Get() (resp *Keys, err error) {
	resp = &Keys{}
	req := internalRequest{
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
