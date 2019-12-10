package meilisearch

import (
	"net/http"
	"strconv"
)

type clientUpdates struct {
	client  *Client
	indexID string
}

func newClientUpdates(client *Client, indexId string) clientUpdates {
	return clientUpdates{client: client, indexID: indexId}
}

func (c clientUpdates) Get(id int64) (resp *Update, err error) {
	resp = &Update{}

	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/updates/" + strconv.FormatInt(id, 10),
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: nil,
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
		endpoint:            "/indexes/" + c.indexID + "/updates",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: nil,
		functionName:        "List",
		apiName:             "Updates",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}
