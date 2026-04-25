package meilisearch

import (
	"context"
	"fmt"
	"net/http"
)

func (m *meilisearch) UpdateSearchRule(uid string, params *SearchRulesRequest) (*SearchRule, error) {
	return m.UpdateSearchRuleWithContext(context.Background(), uid, params)
}

func (m *meilisearch) UpdateSearchRuleWithContext(ctx context.Context, uid string, params *SearchRulesRequest) (*SearchRule, error) {
	resp := new(SearchRule)

	req := &internalRequest{
		endpoint:            fmt.Sprintf("/dynamic-search-rules/%s", uid),
		method:              http.MethodPatch,
		contentType:         contentTypeJSON,
		withRequest:         params,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusCreated, http.StatusOK},
		functionName:        "UpdateSearchRule",
	}

	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *meilisearch) ListSearchRules(params *SearchRulesParams) (*SearchRulesResults, error) {
	return m.ListSearchRulesWithContext(context.Background(), params)
}

func (m *meilisearch) ListSearchRulesWithContext(ctx context.Context, params *SearchRulesParams) (*SearchRulesResults, error) {
	resp := new(SearchRulesResults)

	req := &internalRequest{
		endpoint:            "/dynamic-search-rules",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
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

func (m *meilisearch) GetSearchRule(uid string) (*SearchRule, error) {
	return m.GetSearchRuleWithContext(context.Background(), uid)
}

func (m *meilisearch) GetSearchRuleWithContext(ctx context.Context, uid string) (*SearchRule, error) {
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

func (m *meilisearch) DeleteSearchRule(uid string) error {
	return m.DeleteSearchRuleWithContext(context.Background(), uid)
}

func (m *meilisearch) DeleteSearchRuleWithContext(ctx context.Context, uid string) error {
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
