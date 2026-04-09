package meilisearch

import "context"

type WebhookManager interface {
	WebhookReader

	// AddWebhook add a new webhook to meilisearch.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/webhooks/create-webhook
	AddWebhook(params *AddWebhookRequest) (*Webhook, error)

	// AddWebhookWithContext add a new webhook to meilisearch with a context.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/webhooks/create-webhook
	AddWebhookWithContext(ctx context.Context, params *AddWebhookRequest) (*Webhook, error)

	// UpdateWebhook modifies a previously existing webhook.
	// If the webhook has isEditable to false the HTTP call returns an error.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/webhooks/update-webhook
	UpdateWebhook(uuid string, params *UpdateWebhookRequest) (*Webhook, error)

	// UpdateWebhookWithContext modifies a previously existing webhook with a context.
	// If the webhook has isEditable to false the HTTP call returns an error.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/webhooks/update-webhook
	UpdateWebhookWithContext(ctx context.Context, uuid string, params *UpdateWebhookRequest) (*Webhook, error)

	// DeleteWebhook deletes an existing webhook. Will also fail when the webhook doesn’t exist.
	// If the webhook has isEditable to false the HTTP call returns an error.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/webhooks/delete-webhook
	DeleteWebhook(uuid string) error

	// DeleteWebhookWithContext deletes an existing webhook with a context. Will also fail when the webhook doesn’t exist.
	// If the webhook has isEditable to false the HTTP call returns an error.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/webhooks/delete-webhook
	DeleteWebhookWithContext(ctx context.Context, uuid string) error
}

type WebhookReader interface {
	// ListWebhooks lists all the webhooks.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/webhooks/list-webhooks
	ListWebhooks() (*WebhookResults, error)

	// ListWebhooksWithContext lists all the webhooks with context.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/webhooks/list-webhooks
	ListWebhooksWithContext(ctx context.Context) (*WebhookResults, error)

	// GetWebhook gets a webhook by uuid.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/webhooks/get-webhook
	GetWebhook(uuid string) (*Webhook, error)

	// GetWebhookWithContext gets a webhook by uuid with a context.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/webhooks/get-webhook
	GetWebhookWithContext(ctx context.Context, uuid string) (*Webhook, error)
}
