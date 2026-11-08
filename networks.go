package meilisearch

import (
	"context"
	"net/http"
)

func (m *meilisearch) UpdateNetwork(ctx context.Context, params *UpdateNetworkRequest) (any, error) {
	resp := new(struct {
		Network
		Task
	})
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

	// Whether Meilisearch enqueues a task or applies the update immediately depends on the network
	// state (a leader is configured after the update), not only on the request. Network and Task
	// share no JSON field, so the body fills only one of them: only an enqueued task has a status.
	if resp.Status != "" {
		return &resp.Task, nil
	}
	return &resp.Network, nil
}

func (m *meilisearch) GetNetwork(ctx context.Context) (*Network, error) {
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
