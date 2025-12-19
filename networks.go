package meilisearch

import (
	"context"
	"net/http"
)

func (m *meilisearch) UpdateNetwork(params *UpdateNetworkRequest) (any, error) {
	return m.UpdateNetworkWithContext(context.Background(), params)
}

func (m *meilisearch) UpdateNetworkWithContext(ctx context.Context, params *UpdateNetworkRequest) (any, error) {
	var resp any
	if params.Leader.Valid() {
		resp = new(Task)
	} else {
		resp = new(Network)
	}
	req := &internalRequest{
		endpoint:            "/network",
		method:              http.MethodPatch,
		contentType:         contentTypeJSON,
		withRequest:         params,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK, http.StatusAccepted},
		functionName:        "UpdateNetwork",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) GetNetwork() (*Network, error) {
	return m.GetNetworkWithContext(context.Background())
}

func (m *meilisearch) GetNetworkWithContext(ctx context.Context) (*Network, error) {
	resp := new(Network)
	req := &internalRequest{
		endpoint:            "/network",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetNetwork",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}
