package meilisearch

import "context"

type WebhookManager interface {
	WebhookReader

	// AddWebhook add a new webhook to meilisearch.
	AddWebhook(params *AddWebhookQuery) (*Webhook, error)
	// AddWebhookWithContext add a new webhook to meilisearch with a context.
	AddWebhookWithContext(ctx context.Context, params *AddWebhookQuery) (*Webhook, error)
}

type WebhookReader interface{}
