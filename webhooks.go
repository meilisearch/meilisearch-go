package meilisearch

import (
	"context"
	"fmt"
	"net/http"
)

func (m *meilisearch) AddWebhook(params *AddWebhookRequest) (*Webhook, error) {
	return m.AddWebhookWithContext(context.Background(), params)
}

func (m *meilisearch) AddWebhookWithContext(ctx context.Context, params *AddWebhookRequest) (*Webhook, error) {
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

func (m *meilisearch) ListWebhooks() (*WebhookResults, error) {
	return m.ListWebhooksWithContext(context.Background())
}

func (m *meilisearch) ListWebhooksWithContext(ctx context.Context) (*WebhookResults, error) {
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

func (m *meilisearch) GetWebhook(uuid string) (*Webhook, error) {
	return m.GetWebhookWithContext(context.Background(), uuid)
}

func (m *meilisearch) GetWebhookWithContext(ctx context.Context, uuid string) (*Webhook, error) {
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

func (m *meilisearch) UpdateWebhook(uuid string, params *UpdateWebhookRequest) (*Webhook, error) {
	return m.UpdateWebhookWithContext(context.Background(), uuid, params)
}

func (m *meilisearch) UpdateWebhookWithContext(ctx context.Context, uuid string, params *UpdateWebhookRequest) (*Webhook, error) {
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

func (m *meilisearch) DeleteWebhook(uuid string) error {
	return m.DeleteWebhookWithContext(context.Background(), uuid)
}

func (m *meilisearch) DeleteWebhookWithContext(ctx context.Context, uuid string) error {
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
