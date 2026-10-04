package meilisearch

import (
	"context"
	"fmt"
	"net/http"
)

func (m *meilisearch) AddWebhook(ctx context.Context, params *AddWebhookRequest) (*Webhook, error) {
	resp := new(Webhook)
	req := &internalRequest{
		endpoint:            "/webhooks",
		method:              http.MethodPost,
		withRequest:         params,
		contentType:         contentTypeJSON,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusCreated, http.StatusOK},
		functionName:        "AddWebhook",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) ListWebhooks(ctx context.Context) (*WebhookResults, error) {
	resp := new(WebhookResults)
	req := &internalRequest{
		endpoint:            "/webhooks",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "ListWebhooks",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) GetWebhook(ctx context.Context, uuid string) (*Webhook, error) {
	resp := new(Webhook)
	req := &internalRequest{
		endpoint:            fmt.Sprintf("/webhooks/%s", uuid),
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetWebhook",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) UpdateWebhook(ctx context.Context, uuid string, params *UpdateWebhookRequest) (*Webhook, error) {
	resp := new(Webhook)
	req := &internalRequest{
		endpoint:            fmt.Sprintf("/webhooks/%s", uuid),
		method:              http.MethodPatch,
		withRequest:         params,
		contentType:         contentTypeJSON,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "UpdateWebhook",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) DeleteWebhook(ctx context.Context, uuid string) error {
	req := &internalRequest{
		endpoint:            fmt.Sprintf("/webhooks/%s", uuid),
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        nil,
		acceptedStatusCodes: []int{http.StatusNoContent},
		functionName:        "DeleteWebhook",
	}
	return m.client.executeRequest(ctx, req)
}
