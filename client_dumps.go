package meilisearch

import (
	"net/http"
)

type clientDumps struct {
	client *Client
}

func newClientDumps(client *Client) clientDumps {
	return clientDumps{client: client}
}

func (c clientDumps) Create() (resp *Dump, err error) {
	resp = &Dump{}
	req := internalRequest{
		endpoint:            "/dumps",
		method:              http.MethodPost,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "Create",
		apiName:             "Dumps",
	}
	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c clientDumps) GetStatus(dumpUID string) (resp *Dump, err error) {
	resp = &Dump{}
	req := internalRequest{
		endpoint:            "/dumps/" + dumpUID + "/status",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetStatus",
		apiName:             "Dumps",
	}
	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}
