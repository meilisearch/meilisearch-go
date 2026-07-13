package meilisearch

type UpdateWebhookRequest struct {
	URL     string            `json:"url,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

type Webhook struct {
	UUID       string            `json:"uuid"`
	IsEditable bool              `json:"isEditable"`
	URL        string            `json:"url"`
	Headers    map[string]string `json:"headers"`
}

type WebhookResults struct {
	Result []*Webhook `json:"results"`
}

type AddWebhookRequest struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
}
