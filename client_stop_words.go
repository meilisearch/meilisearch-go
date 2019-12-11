package meilisearch

import "net/http"

type clientStopWords struct {
	client  *Client
	indexID string
}

func newClientStopWords(client *Client, indexId string) clientStopWords {
	return clientStopWords{client: client, indexID: indexId}
}

func (c clientStopWords) List() (resp []string, err error) {
	resp = make([]string, 0)
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/stop-words",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: nil,
		functionName:        "List",
		apiName:             "StopWords",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c clientStopWords) Add(words []string) (resp *AsyncUpdateId, err error) {
	resp = &AsyncUpdateId{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/stop-words",
		method:              http.MethodPatch,
		withRequest:         &words,
		withResponse:        &resp,
		acceptedStatusCodes: nil,
		functionName:        "Add",
		apiName:             "StopWords",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c clientStopWords) Deletes(words []string) (resp *AsyncUpdateId, err error) {
	resp = &AsyncUpdateId{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/stop-words",
		method:              http.MethodPost,
		withRequest:         &words,
		withResponse:        &resp,
		acceptedStatusCodes: nil,
		functionName:        "Add",
		apiName:             "StopWords",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c clientStopWords) IndexId() string {
	return c.indexID
}

func (c clientStopWords) Client() *Client {
	return c.client
}
