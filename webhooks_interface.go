package meilisearch

import "context"

type WebhookManager interface {
	WebhookReader

	// AddWebhook add a new webhook to meilisearch.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/webhooks/create-webhook
	AddWebhook(ctx context.Context, params *AddWebhookRequest) (*Webhook, error)

	// UpdateWebhook modifies a previously existing webhook.
	// If the webhook has isEditable to false the HTTP call returns an error.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/webhooks/update-webhook
	UpdateWebhook(ctx context.Context, uuid string, params *UpdateWebhookRequest) (*Webhook, error)

	// DeleteWebhook deletes an existing webhook. Will also fail when the webhook doesn’t exist.
	// If the webhook has isEditable to false the HTTP call returns an error.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/webhooks/delete-webhook
	DeleteWebhook(ctx context.Context, uuid string) error
}

type WebhookReader interface {
	// ListWebhooks lists all the webhooks.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/webhooks/list-webhooks
	ListWebhooks(ctx context.Context) (*WebhookResults, error)

	// GetWebhook gets a webhook by uuid.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/webhooks/get-webhook
	GetWebhook(ctx context.Context, uuid string) (*Webhook, error)
}
