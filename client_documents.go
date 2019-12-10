package meilisearch

import "net/http"

type clientDocuments struct {
	client  *Client
	indexId string
}

func newClientDocuments(client *Client, indexId string) *clientDocuments {
	return &clientDocuments{client: client, indexId: indexId}
}

func (c clientDocuments) Get(identifier string, documentPtr interface{}) error {
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexId + "/documents/" + identifier,
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

func (c clientDocuments) Delete(identifier string) (resp *UpdateIdResponse, err error) {
	resp = &UpdateIdResponse{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexId + "/documents/" + identifier,
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: nil,
		functionName:        "Delete",
		apiName:             "Documents",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c clientDocuments) Deletes(identifier []string) (resp *UpdateIdResponse, err error) {
	resp = &UpdateIdResponse{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexId + "/documents/delete",
		method:              http.MethodPost,
		withRequest:         &identifier,
		withResponse:        resp,
		acceptedStatusCodes: nil,
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
		endpoint:            "/indexes/" + c.indexId + "/documents",
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

func (c clientDocuments) AddOrUpdate(documentsPtr interface{}) (resp *UpdateIdResponse, err error) {
	resp = &UpdateIdResponse{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexId + "/documents",
		method:              http.MethodPost,
		withRequest:         documentsPtr,
		withResponse:        resp,
		acceptedStatusCodes: nil,
		functionName:        "AddOrUpdate",
		apiName:             "Documents",
	}

	if err = c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c clientDocuments) ClearAllDocuments() (resp *UpdateIdResponse, err error) {
	resp = &UpdateIdResponse{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexId + "/documents",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: nil,
		functionName:        "ClearAllDocuments",
		apiName:             "Documents",
	}

	if err = c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}
