package meilisearch

import (
	"net/http"
)

func (c *Client) Index(uid string) *Index {
	return newIndex(c, uid)
}

func (c *Client) GetIndex(uid string) (resp *Index, err error) {
	return newIndex(c, uid).FetchInfo()
}

func (c *Client) GetRawIndex(uid string) (resp map[string]interface{}, err error) {
	resp = map[string]interface{}{}
	req := internalRequest{
		endpoint:            "/indexes/" + uid,
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetRawIndex",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) CreateIndex(config *IndexConfig) (resp *Task, err error) {
	request := &CreateIndexRequest{
		UID:        config.Uid,
		PrimaryKey: config.PrimaryKey,
	}
	resp = &Task{}
	req := internalRequest{
		endpoint:            "/indexes",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "CreateIndex",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetAllIndexes() (resp []*Index, err error) {
	resp = []*Index{}
	req := internalRequest{
		endpoint:            "/indexes",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetAllIndexes",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetAllRawIndexes() (resp []map[string]interface{}, err error) {
	resp = []map[string]interface{}{}
	req := internalRequest{
		endpoint:            "/indexes",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetAllRawIndexes",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) DeleteIndex(uid string) (resp *Task, err error) {
	resp = &Task{}
	req := internalRequest{
		endpoint:            "/indexes/" + uid,
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "DeleteIndex",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}
