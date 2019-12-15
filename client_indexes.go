package meilisearch

import "net/http"

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

func (c clientIndexes) Update(uid string, name string) (resp *Index, err error) {
	resp = &Index{}
	req := internalRequest{
		endpoint: "/indexes/" + uid,
		method:   http.MethodPut,
		withRequest: &map[string]string{
			"name": name,
		},
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Update",
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

func (c clientIndexes) GetRawSchema(uid string) (resp *RawSchema, err error) {
	resp = &RawSchema{}
	req := internalRequest{
		endpoint:            "/indexes/" + uid + "/schema?raw=true",
		method:              http.MethodGet,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetRawSchema",
		apiName:             "Indexes",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c clientIndexes) GetSchema(uid string) (resp *Schema, err error) {
	resp = &Schema{}
	req := internalRequest{
		endpoint:            "/indexes/" + uid + "/schema",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetSchema",
		apiName:             "Indexes",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c clientIndexes) UpdateSchema(uid string, schema Schema) (resp *AsyncUpdateId, err error) {
	resp = &AsyncUpdateId{}
	req := internalRequest{
		endpoint:            "/indexes/" + uid + "/schema",
		method:              http.MethodPut,
		withRequest:         &schema,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateSchema",
		apiName:             "Indexes",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c clientIndexes) UpdateWithRawSchema(uid string, schema RawSchema) (resp *AsyncUpdateId, err error) {
	resp = &AsyncUpdateId{}
	req := internalRequest{
		endpoint:            "/indexes/" + uid + "/schema?raw=true",
		method:              http.MethodPut,
		withRequest:         &schema,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateWithRawSchema",
		apiName:             "Indexes",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}
