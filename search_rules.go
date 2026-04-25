package meilisearch

import (
	"context"
	"fmt"
	"net/http"
)

func (m *meilisearch) Update(uid string, params *SearchRulesRequest) (*SearchRule, error) {
	return m.UpdateWithContext(context.Background(), uid, params)
}

func (m *meilisearch) UpdateWithContext(ctx context.Context, uid string, params *SearchRulesRequest) (*SearchRule, error) {
	resp := new(SearchRule)

	req := &internalRequest{
		endpoint:            fmt.Sprintf("/dynamic-search-rules/%s", uid),
		method:              http.MethodPatch,
		withRequest:         params,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "UpdateSearchRule",
	}

	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *meilisearch) List(params *SearchRulesParams) (*SearchRulesResults, error) {
	return m.ListWithContext(context.Background(), params)
}

func (m *meilisearch) ListWithContext(ctx context.Context, params *SearchRulesParams) (*SearchRulesResults, error) {
	resp := new(SearchRulesResults)

	req := &internalRequest{
		endpoint:            "/dynamic-search-rules",
		method:              http.MethodPost,
		withRequest:         params,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "ListSearchRules",
	}

	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *meilisearch) Get(uid string) (*SearchRule, error) {
	return m.GetWithContext(context.Background(), uid)
}

func (m *meilisearch) GetWithContext(ctx context.Context, uid string) (*SearchRule, error) {
	resp := new(SearchRule)

	req := &internalRequest{
		endpoint:            fmt.Sprintf("/dynamic-search-rules/%s", uid),
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetSearchRule",
	}

	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *meilisearch) Delete(uid string) error {
	return m.DeleteWithContext(context.Background(), uid)
}

func (m *meilisearch) DeleteWithContext(ctx context.Context, uid string) error {
	req := &internalRequest{
		endpoint:            fmt.Sprintf("/dynamic-search-rules/%s", uid),
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        nil,
		acceptedStatusCodes: []int{http.StatusNoContent},
		functionName:        "DeleteSearchRule",
	}
	return m.client.executeRequest(ctx, req)
}
