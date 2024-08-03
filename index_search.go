package meilisearch

import (
	"encoding/json"
	"errors"
	"net/http"
)

// This constant contains the default values assigned by Meilisearch to the limit in search parameters
//
// Documentation: https://www.meilisearch.com/docs/reference/api/search#search-parameters
const (
	DefaultLimit int64 = 20
)

var ErrNoSearchRequest = errors.New("no search request provided")

func (i Index) SearchRaw(query string, request *SearchRequest) (*json.RawMessage, error) {
	if request == nil {
		return nil, ErrNoSearchRequest
	}

	if query != "" {
		request.Query = query
	}

	if request.IndexUID != "" {
		request.IndexUID = ""
	}

	request.validate()

	resp := &json.RawMessage{}

	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/search",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "SearchRaw",
	}

	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (i Index) Search(query string, request *SearchRequest) (*SearchResponse, error) {
	if request == nil {
		return nil, ErrNoSearchRequest
	}

	if query != "" {
		request.Query = query
	}

	if request.IndexUID != "" {
		request.IndexUID = ""
	}

	request.validate()

	resp := &SearchResponse{}

	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/search",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Search",
	}

	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (i Index) SearchSimilarDocuments(param *SimilarDocumentQuery, resp *SimilarDocumentResult) error {
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/similar",
		method:              http.MethodPost,
		withRequest:         param,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "SearchSimilarDocuments",
		contentType:         contentTypeJSON,
	}

	if err := i.client.executeRequest(req); err != nil {
		return err
	}
	return nil
}
