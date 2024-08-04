package meilisearch

import (
	"context"
	"encoding/json"
	"net/http"
)

func (i *index) Search(query string, request *SearchRequest) (*SearchResponse, error) {
	return i.SearchWithContext(context.Background(), query, request)
}

func (i *index) SearchWithContext(ctx context.Context, query string, request *SearchRequest) (*SearchResponse, error) {
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

	resp := new(SearchResponse)

	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/search",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Search",
	}

	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (i *index) SearchRaw(query string, request *SearchRequest) (*json.RawMessage, error) {
	return i.SearchRawWithContext(context.Background(), query, request)
}

func (i *index) SearchRawWithContext(ctx context.Context, query string, request *SearchRequest) (*json.RawMessage, error) {
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

	resp := new(json.RawMessage)

	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/search",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "SearchRaw",
	}

	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (i *index) FacetSearch(request *FacetSearchRequest) (*json.RawMessage, error) {
	return i.FacetSearchWithContext(context.Background(), request)
}

func (i *index) FacetSearchWithContext(ctx context.Context, request *FacetSearchRequest) (*json.RawMessage, error) {
	if request == nil {
		return nil, ErrNoFacetSearchRequest
	}

	resp := new(json.RawMessage)

	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/facet-search",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "FacetSearch",
	}

	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (i *index) SearchSimilarDocuments(param *SimilarDocumentQuery, resp *SimilarDocumentResult) error {
	return i.SearchSimilarDocumentsWithContext(context.Background(), param, resp)
}

func (i *index) SearchSimilarDocumentsWithContext(ctx context.Context, param *SimilarDocumentQuery, resp *SimilarDocumentResult) error {
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/similar",
		method:              http.MethodPost,
		withRequest:         param,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "SearchSimilarDocuments",
		contentType:         contentTypeJSON,
	}

	if err := i.client.executeRequest(ctx, req); err != nil {
		return err
	}
	return nil
}
