package meilisearch

import "net/http"

type clientDocuments struct {
	client  *Client
	indexID string
}

func newClientDocuments(client *Client, indexId string) clientDocuments {
	return clientDocuments{client: client, indexID: indexId}
}

func (c clientDocuments) Get(identifier string, documentPtr interface{}) error {
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/documents/" + identifier,
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        documentPtr,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Get",
		apiName:             "Documents",
	}

	if err := c.client.executeRequest(req); err != nil {
		return err
	}

	return nil
}

func (c clientDocuments) Delete(identifier string) (resp *AsyncUpdateId, err error) {
	resp = &AsyncUpdateId{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/documents/" + identifier,
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "Delete",
		apiName:             "Documents",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c clientDocuments) Deletes(identifier []string) (resp *AsyncUpdateId, err error) {
	resp = &AsyncUpdateId{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/documents/delete",
		method:              http.MethodPost,
		withRequest:         &identifier,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "Deletes",
		apiName:             "Documents",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c clientDocuments) List(request ListDocumentsRequest, documentsPtr interface{}) error {
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/documents",
		method:              http.MethodGet,
		withRequest:         &request,
		withResponse:        documentsPtr,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "List",
		apiName:             "Documents",
	}

	if err := c.client.executeRequest(req); err != nil {
		return err
	}

	return nil
}

func (c clientDocuments) AddOrUpdate(documentsPtr interface{}) (resp *AsyncUpdateId, err error) {
	resp = &AsyncUpdateId{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/documents",
		method:              http.MethodPost,
		withRequest:         documentsPtr,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "AddOrUpdate",
		apiName:             "Documents",
	}

	if err = c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c clientDocuments) DeleteAllDocuments() (resp *AsyncUpdateId, err error) {
	resp = &AsyncUpdateId{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/documents",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "DeleteAllDocuments",
		apiName:             "Documents",
	}

	if err = c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c clientDocuments) IndexId() string {
	return c.indexID
}

func (c clientDocuments) Client() *Client {
	return c.client
}
