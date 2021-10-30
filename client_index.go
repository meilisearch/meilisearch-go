package meilisearch

import "net/http"

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

func (c *Client) CreateIndex(config *IndexConfig) (resp *Index, err error) {
	request := &CreateIndexRequest{
		UID:        config.Uid,
		PrimaryKey: config.PrimaryKey,
	}
	resp = newIndex(c, config.Uid)
	req := internalRequest{
		endpoint:            "/indexes",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusCreated},
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

func (c *Client) GetOrCreateIndex(config *IndexConfig) (resp *Index, err error) {
	resp, err = c.GetIndex(config.Uid)
	if err == nil {
		return resp, err
	}
	return c.CreateIndex(config)
}

func (c *Client) DeleteIndex(uid string) (ok bool, err error) {
	req := internalRequest{
		endpoint:            "/indexes/" + uid,
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        nil,
		acceptedStatusCodes: []int{http.StatusNoContent},
		functionName:        "DeleteIndex",
	}
	// err is not nil if status code is not 204 StatusNoContent
	if err := c.executeRequest(req); err != nil {
		return false, err
	}
	return true, nil
}

func (c *Client) DeleteIndexIfExists(uid string) (ok bool, err error) {
	req := internalRequest{
		endpoint:            "/indexes/" + uid,
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        nil,
		acceptedStatusCodes: []int{http.StatusNoContent},
		functionName:        "DeleteIndex",
	}
	// err is not nil if status code is not 204 StatusNoContent
	if err := c.executeRequest(req); err != nil {
		if err.(*Error).MeilisearchApiMessage.Code != "index_not_found" {
			return false, err
		}
		return false, nil
	}
	return true, nil
}
