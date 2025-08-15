package meilisearch

import "context"

type WebhookManager interface {
	WebhookReader

	// AddWebhook add a new webhook to meilisearch.
	AddWebhook(params *AddWebhookRequest) (*Webhook, error)
	// AddWebhookWithContext add a new webhook to meilisearch with a context.
	AddWebhookWithContext(ctx context.Context, params *AddWebhookRequest) (*Webhook, error)

	// UpdateWebhook updates a webhook by uuid.
	// If the webhook has isEditable to false the HTTP call returns an error.
	UpdateWebhook(uuid string, params *UpdateWebhookRequest) (*Webhook, error)
	// UpdateWebhookWithContext updates a webhook by uuid with a context.
	// If the webhook has isEditable to false the HTTP call returns an error.
	UpdateWebhookWithContext(ctx context.Context, uuid string, params *UpdateWebhookRequest) (*Webhook, error)

	// DeleteWebhook deletes a webhook by id.
	// If the webhook has isEditable to false the HTTP call returns an error.
	DeleteWebhook(uuid string) error

	// DeleteWebhookWithContext deletes a webhook by id with a context.
	// If the webhook has isEditable to false the HTTP call returns an error.
	DeleteWebhookWithContext(ctx context.Context, uuid string) error
}

type WebhookReader interface {
	// ListWebhooks lists all the webhooks.
	ListWebhooks() (*WebhookResults, error)
	// ListWebhooksWithContext lists all the webhooks with context.
	ListWebhooksWithContext(ctx context.Context) (*WebhookResults, error)

	// GetWebhook gets a webhook by uuid.
	GetWebhook(uuid string) (*Webhook, error)
	// GetWebhookWithContext gets a webhook by uuid with a context.
	GetWebhookWithContext(ctx context.Context, uuid string) (*Webhook, error)
}
