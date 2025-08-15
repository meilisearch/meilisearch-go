package meilisearch

import (
	"context"
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
