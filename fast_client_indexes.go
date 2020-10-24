package meilisearch

import (
	"net/http"
)

type fastClientIndexes struct {
	client *FastHTTPClient
}

func newFastClientIndexes(client *FastHTTPClient) fastClientIndexes {
	return fastClientIndexes{client: client}
}

func (c fastClientIndexes) Get(uid string) (resp *Index, err error) {
	resp = &Index{}
	req := internalRawRequest{
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

func (c fastClientIndexes) List() (resp []Index, err error) {
	resp = []Index{}

	req := internalRawRequest{
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

func (c fastClientIndexes) Create(request CreateIndexRequest) (resp *CreateIndexResponse, err error) {
	resp = &CreateIndexResponse{}
	req := internalRawRequest{
		endpoint:            "/indexes",
		method:              http.MethodPost,
		withRequest:         request,
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

func (c fastClientIndexes) UpdateName(uid string, name string) (resp *Index, err error) {
	resp = &Index{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + uid,
		method:              http.MethodPut,
		withRequest:         &Name{Name: name},
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

func (c fastClientIndexes) UpdatePrimaryKey(uid string, primaryKey string) (resp *Index, err error) {
	resp = &Index{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + uid,
		method:              http.MethodPut,
		withRequest:         &PrimaryKey{PrimaryKey: primaryKey},
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

func (c fastClientIndexes) Delete(uid string) (ok bool, err error) {
	req := internalRawRequest{
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
