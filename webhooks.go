package meilisearch

import (
	"context"
	"fmt"
	"net/http"
)

func (m *meilisearch) AddWebhook(params *AddWebhookQuery) (*Webhook, error) {
	return m.AddWebhookWithContext(context.Background(), params)
}

func (m *meilisearch) AddWebhookWithContext(ctx context.Context, params *AddWebhookQuery) (*Webhook, error) {
	resp := new(Webhook)
	req := &internalRequest{
		endpoint:            "/webhooks",
		method:              http.MethodPost,
		withRequest:         params,
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
