package meilisearch

import (
	"net/http"
	"strconv"
	"strings"
)

type clientDocuments struct {
	client   *Client
	indexUID string
}

func newClientDocuments(client *Client, indexUID string) clientDocuments {
	return clientDocuments{client: client, indexUID: indexUID}
}

func (c clientDocuments) Get(identifier string, documentPtr interface{}) error {
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexUID + "/documents/" + identifier,
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

func (c clientDocuments) Delete(identifier string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexUID + "/documents/" + identifier,
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

func (c clientDocuments) Deletes(identifier []string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexUID + "/documents/delete-batch",
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
		endpoint:     "/indexes/" + c.indexUID + "/documents",
		method:       http.MethodGet,
		withRequest:  &request,
		withResponse: documentsPtr,
		withQueryParams: map[string]string{
			"limit":                strconv.FormatInt(request.Limit, 10),
			"offset":               strconv.FormatInt(request.Offset, 10),
			"attributesToRetrieve": strings.Join(request.AttributesToRetrieve, ","),
		},
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "List",
		apiName:             "Documents",
	}

	if err := c.client.executeRequest(req); err != nil {
		return err
	}

	return nil
}

func (c clientDocuments) AddOrReplace(documentsPtr interface{}) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexUID + "/documents",
		method:              http.MethodPost,
		withRequest:         documentsPtr,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "AddOrReplace",
		apiName:             "Documents",
	}

	if err = c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c clientDocuments) AddOrReplaceWithPrimaryKey(documentsPtr interface{}, primaryKey string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexUID + "/documents?primaryKey=" + primaryKey,
		method:              http.MethodPost,
		withRequest:         documentsPtr,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "AddOrReplaceWithPrimaryKey",
		apiName:             "Documents",
	}

	if err = c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c clientDocuments) AddOrUpdate(documentsPtr interface{}) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexUID + "/documents",
		method:              http.MethodPut,
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

func (c clientDocuments) AddOrUpdateWithPrimaryKey(documentsPtr interface{}, primaryKey string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexUID + "/documents?primaryKey=" + primaryKey,
		method:              http.MethodPut,
		withRequest:         documentsPtr,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "AddOrUpdateWithPrimaryKey",
		apiName:             "Documents",
	}

	if err = c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c clientDocuments) DeleteAllDocuments() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexUID + "/documents",
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

func (c clientDocuments) IndexID() string {
	return c.indexUID
}

func (c clientDocuments) Client() *Client {
	return c.client
}
