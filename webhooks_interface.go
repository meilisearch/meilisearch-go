package meilisearch

import "context"

type WebhookManager interface {
	WebhookReader

	// AddWebhook add a new webhook to meilisearch.
	AddWebhook(params *AddWebhookQuery) (*Webhook, error)
	// AddWebhookWithContext add a new webhook to meilisearch with a context.
	AddWebhookWithContext(ctx context.Context, params *AddWebhookQuery) (*Webhook, error)
}

type WebhookReader interface{
	// ListWebhooks lists all the webhooks.
	ListWebhooks() (*WebhookResults, error)
	// ListWebhooksWithContext lists all the webhooks with context.
	ListWebhooksWithContext(ctx context.Context) (*WebhookResults, error)
	// GetWebhook gets a webhook by uuid.
	GetWebhook(uuid string) (*Webhook, error)
	// GetWebhookWithContext gets a webhook by uuid with a context.
	GetWebhookWithContext(ctx context.Context, uuid string) (*Webhook, error)
}
