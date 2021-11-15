package meilisearch

import (
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func (i Index) GetDocument(identifier string, documentPtr interface{}) error {
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/documents/" + identifier,
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        documentPtr,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetDocument",
	}
	if err := i.client.executeRequest(req); err != nil {
		return err
	}
	return nil
}

func (i Index) GetDocuments(request *DocumentsRequest, resp interface{}) error {
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/documents",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		withQueryParams:     map[string]string{},
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetDocuments",
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
	if err := i.client.executeRequest(req); err != nil {
		return err
	}
	return nil
}

func (i Index) AddDocuments(documentsPtr interface{}, primaryKey ...string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	endpoint := ""
	if primaryKey == nil {
		endpoint = "/indexes/" + i.UID + "/documents"
	} else {
		i.PrimaryKey = primaryKey[0] //nolint:golint,staticcheck
		endpoint = "/indexes/" + i.UID + "/documents?primaryKey=" + primaryKey[0]
	}
	req := internalRequest{
		endpoint:            endpoint,
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         documentsPtr,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "AddDocuments",
	}
	if err = i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) AddDocumentsInBatches(documentsPtr interface{}, batchSize int, primaryKey ...string) (resp []AsyncUpdateID, err error) {
	arr := reflect.ValueOf(documentsPtr)
	lenDocs := arr.Len()
	numBatches := int(math.Ceil(float64(lenDocs) / float64(batchSize)))
	resp = make([]AsyncUpdateID, numBatches)

	for j := 0; j < numBatches; j++ {
		end := (j + 1) * batchSize
		if end > lenDocs {
			end = lenDocs
		}

		batch := arr.Slice(j*batchSize, end).Interface()

		if primaryKey != nil {
			respID, err := i.AddDocuments(batch, primaryKey[0])
			if err != nil {
				return nil, err
			}

			resp[j] = *respID
		} else {
			respID, err := i.AddDocuments(batch)
			if err != nil {
				return nil, err
			}

			resp[j] = *respID
		}
	}

	return resp, nil
}

func (i Index) UpdateDocuments(documentsPtr interface{}, primaryKey ...string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	endpoint := ""
	if primaryKey == nil {
		endpoint = "/indexes/" + i.UID + "/documents"
	} else {
		i.PrimaryKey = primaryKey[0] //nolint:golint,staticcheck
		endpoint = "/indexes/" + i.UID + "/documents?primaryKey=" + primaryKey[0]
	}
	req := internalRequest{
		endpoint:            endpoint,
		method:              http.MethodPut,
		contentType:         contentTypeJSON,
		withRequest:         documentsPtr,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateDocuments",
	}
	if err = i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) UpdateDocumentsInBatches(documentsPtr interface{}, batchSize int, primaryKey ...string) (resp []AsyncUpdateID, err error) {
	arr := reflect.ValueOf(documentsPtr)
	lenDocs := arr.Len()
	numBatches := int(math.Ceil(float64(lenDocs) / float64(batchSize)))
	resp = make([]AsyncUpdateID, numBatches)

	for j := 0; j < numBatches; j++ {
		end := (j + 1) * batchSize
		if end > lenDocs {
			end = lenDocs
		}

		batch := arr.Slice(j*batchSize, end).Interface()
		if primaryKey != nil {
			respID, err := i.UpdateDocuments(batch, primaryKey[0])
			if err != nil {
				return nil, err
			}

			resp[j] = *respID
		} else {
			respID, err := i.UpdateDocuments(batch)
			if err != nil {
				return nil, err
			}

			resp[j] = *respID
		}
	}

	return resp, nil
}

func (i Index) DeleteDocument(identifier string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/documents/" + identifier,
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "DeleteDocument",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) DeleteDocuments(identifier []string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/documents/delete-batch",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         identifier,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "DeleteDocuments",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) DeleteAllDocuments() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/documents",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "DeleteAllDocuments",
	}
	if err = i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}
