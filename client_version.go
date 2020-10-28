package meilisearch

import "net/http"

type clientVersion struct {
	client *Client
}

func newClientVersion(client *Client) clientVersion {
	return clientVersion{client: client}
}

func (c clientVersion) Get() (resp *Version, err error) {
	resp = &Version{}

	req := internalRequest{
		endpoint:            "/version",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Get",
		apiName:             "Version",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}
