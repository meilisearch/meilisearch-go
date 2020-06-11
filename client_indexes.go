package meilisearch

import (
	"net/http"
)

type clientIndexes struct {
	client *Client
}

func newClientIndexes(client *Client) clientIndexes {
	return clientIndexes{client: client}
}

func (c clientIndexes) Get(uid string) (resp *Index, err error) {
	resp = &Index{}
	req := internalRequest{
		endpoint:            "/indexes/" + uid,
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Get",
		apiName:             "Indexes",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c clientIndexes) List() (resp []Index, err error) {
	resp = []Index{}
	req := internalRequest{
		endpoint:            "/indexes",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "List",
		apiName:             "Indexes",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c clientIndexes) Create(request CreateIndexRequest) (resp *CreateIndexResponse, err error) {
	resp = &CreateIndexResponse{}
	req := internalRequest{
		endpoint:            "/indexes",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusCreated},
		functionName:        "Create",
		apiName:             "Indexes",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c clientIndexes) UpdateName(uid string, name string) (resp *Index, err error) {
	resp = &Index{}
	req := internalRequest{
		endpoint: "/indexes/" + uid,
		method:   http.MethodPut,
		withRequest: &map[string]string{
			"name": name,
		},
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "UpdateName",
		apiName:             "Indexes",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c clientIndexes) UpdatePrimaryKey(uid string, primaryKey string) (resp *Index, err error) {
	resp = &Index{}
	req := internalRequest{
		endpoint: "/indexes/" + uid,
		method:   http.MethodPut,
		withRequest: &map[string]string{
			"primaryKey": primaryKey,
		},
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "UpdatePrimaryKey",
		apiName:             "Indexes",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c clientIndexes) Delete(uid string) (ok bool, err error) {
	req := internalRequest{
		endpoint:            "/indexes/" + uid,
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        nil,
		acceptedStatusCodes: []int{http.StatusNoContent},
		functionName:        "Delete",
		apiName:             "Indexes",
	}

	// err is not nil if status code is not 204 StatusNoContent
	if err := c.client.executeRequest(req); err != nil {
		return false, err
	}

	return true, nil
}

func (c clientIndexes) DeleteAllIndexes() (ok bool, err error) {
	list, err := c.List()

	if err != nil {
		return false, err
	}

	for _, index := range list {
		c.Delete(index.UID)
	}

	return true, nil
}
