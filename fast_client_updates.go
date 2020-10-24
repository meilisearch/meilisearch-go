package meilisearch

import (
	"net/http"
	"strconv"
)

type fastClientUpdates struct {
	client   *FastHTTPClient
	indexUID string
}

func newFastClientUpdates(client *FastHTTPClient, indexUID string) fastClientUpdates {
	return fastClientUpdates{client: client, indexUID: indexUID}
}

func (c fastClientUpdates) Get(id int64) (resp *Update, err error) {
	resp = &Update{}

	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/updates/" + strconv.FormatInt(id, 10),
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Get",
		apiName:             "Updates",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c fastClientUpdates) List() (resp []Update, err error) {
	resp = []Update{}

	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/updates",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "List",
		apiName:             "Updates",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c fastClientUpdates) IndexID() string {
	return c.indexUID
}

func (c fastClientUpdates) Client() ClientInterface {
	return c.client
}
