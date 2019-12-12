package meilisearch

import "net/http"

type clientKeys struct {
	client *Client
}

func newClientKeys(client *Client) clientKeys {
	return clientKeys{client: client}
}

func (c clientKeys) Get(key string) (resp *APIKey, err error) {
	resp = &APIKey{}
	req := internalRequest{
		endpoint:            "/keys/" + key,
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Get",
		apiName:             "Keys",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c clientKeys) List() (resp []APIKey, err error) {
	resp = make([]APIKey, 0)
	req := internalRequest{
		endpoint:            "/keys",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "List",
		apiName:             "Keys",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil

}

func (c clientKeys) Create(request CreateApiKeyRequest) (resp *APIKey, err error) {
	resp = &APIKey{}
	req := internalRequest{
		endpoint:            "/keys",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusCreated},
		functionName:        "Create",
		apiName:             "Keys",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil

}

func (c clientKeys) Update(key string, request UpdateApiKeyRequest) (resp *APIKey, err error) {
	resp = &APIKey{}
	req := internalRequest{
		endpoint:            "/keys/" + key,
		method:              http.MethodPut,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Update",
		apiName:             "Keys",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c clientKeys) Delete(key string) (deleted bool, err error) {
	req := internalRequest{
		endpoint:            "/keys/" + key,
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        nil,
		acceptedStatusCodes: []int{http.StatusNoContent},
		functionName:        "Delete",
		apiName:             "Keys",
	}

	if err := c.client.executeRequest(req); err != nil {
		return false, err
	}
	return true, nil
}
