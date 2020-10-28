package meilisearch

import (
	"net/http"
	"strconv"
)

type clientUpdates struct {
	client   *Client
	indexUID string
}

func newClientUpdates(client *Client, indexUID string) clientUpdates {
	return clientUpdates{client: client, indexUID: indexUID}
}

func (c clientUpdates) Get(id int64) (resp *Update, err error) {
	resp = &Update{}

	req := internalRequest{
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

func (c clientUpdates) List() (resp []Update, err error) {
	resp = []Update{}

	req := internalRequest{
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

func (c clientUpdates) IndexID() string {
	return c.indexUID
}

func (c clientUpdates) Client() ClientInterface {
	return c.client
}
