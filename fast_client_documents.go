package meilisearch

import (
	"net/http"
	"strconv"
	"strings"
)

type fastClientDocuments struct {
	client   *FastHTTPClient
	indexUID string
}

func newFastClientDocuments(client *FastHTTPClient, indexUID string) fastClientDocuments {
	return fastClientDocuments{client: client, indexUID: indexUID}
}

func (c fastClientDocuments) Get(identifier string, documentPtr interface{}) error {
	req := internalRawRequest{
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

func (c fastClientDocuments) Delete(identifier string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
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

func (c fastClientDocuments) Deletes(identifier []string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/documents/delete-batch",
		method:              http.MethodPost,
		withRequest:         identifier,
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

func (c fastClientDocuments) List(request ListDocumentsRequest, response interface{}) error {
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/documents",
		method:              http.MethodGet,
		withRequest:         request,
		withResponse:        response,
		withQueryParams:     map[string]string{},
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "List",
		apiName:             "Documents",
	}

	if request.Limit != 0 {
		req.withQueryParams["limit"] = strconv.FormatInt(request.Limit, 10)
	}
	if request.Offset != 0 {
		req.withQueryParams["offset"] = strconv.FormatInt(request.Offset, 10)
	}
	if len(request.AttributesToRetrieve) != 0 {
		req.withQueryParams["attributesToRetrieve"] = strings.Join(request.AttributesToRetrieve, ",")
	}

	if err := c.client.executeRequest(req); err != nil {
		return err
	}

	return nil
}

func (c fastClientDocuments) AddOrReplace(documentsPtr interface{}) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
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

func (c fastClientDocuments) AddOrReplaceWithPrimaryKey(documentsPtr interface{}, primaryKey string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
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

func (c fastClientDocuments) AddOrUpdate(documentsPtr interface{}) (*AsyncUpdateID, error) {
	var err error
	resp := &AsyncUpdateID{}
	req := internalRawRequest{
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

func (c fastClientDocuments) AddOrUpdateWithPrimaryKey(documentsPtr interface{}, primaryKey string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
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

func (c fastClientDocuments) DeleteAllDocuments() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
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

func (c fastClientDocuments) IndexID() string {
	return c.indexUID
}

func (c fastClientDocuments) Client() ClientInterface {
	return c.client
}
