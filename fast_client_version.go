package meilisearch

import "net/http"

type fastClientVersion struct {
	client *FastHTTPClient
}

func newFastClientVersion(client *FastHTTPClient) fastClientVersion {
	return fastClientVersion{client: client}
}

func (c fastClientVersion) Get() (resp *Version, err error) {
	resp = &Version{}

	req := internalRawRequest{
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
