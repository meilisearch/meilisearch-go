package meilisearch

import (
	"context"
	"fmt"
	"net/http"
)

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
		acceptedStatusCodes: []int{http.StatusNoContent},
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
