package meilisearch

import (
	"net/http"
	"strconv"
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

func (c *Client) CreateIndex(config *IndexConfig) (resp *TaskInfo, err error) {
	request := &CreateIndexRequest{
		UID:        config.Uid,
		PrimaryKey: config.PrimaryKey,
	}
	resp = &TaskInfo{}
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

func (c *Client) GetIndexes(param *IndexesQuery) (resp *IndexesResults, err error) {
	resp = &IndexesResults{}
	req := internalRequest{
		endpoint:            "/indexes",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		withQueryParams:     map[string]string{},
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetIndexes",
	}
	if param != nil && param.Limit != 0 {
		req.withQueryParams["limit"] = strconv.FormatInt(param.Limit, 10)
	}
	if param != nil && param.Offset != 0 {
		req.withQueryParams["offset"] = strconv.FormatInt(param.Offset, 10)
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetRawIndexes(param *IndexesQuery) (resp map[string]interface{}, err error) {
	resp = map[string]interface{}{}
	req := internalRequest{
		endpoint:            "/indexes",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		withQueryParams:     map[string]string{},
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetRawIndexes",
	}
	if param != nil && param.Limit != 0 {
		req.withQueryParams["limit"] = strconv.FormatInt(param.Limit, 10)
	}
	if param != nil && param.Offset != 0 {
		req.withQueryParams["offset"] = strconv.FormatInt(param.Offset, 10)
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) DeleteIndex(uid string) (resp *TaskInfo, err error) {
	resp = &TaskInfo{}
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
