package meilisearch

import (
	"context"
	"net/http"
)

func (m *meilisearch) UpdateNetwork(params *Network) (*Network, error) {
	return m.UpdateNetworkWithContext(context.Background(), params)
}

func (m *meilisearch) UpdateNetworkWithContext(ctx context.Context, params *Network) (*Network, error) {
	resp := new(Network)
	req := &internalRequest{
		endpoint:            "/network",
		method:              http.MethodPatch,
		contentType:         contentTypeJSON,
		withRequest:         params,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
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
