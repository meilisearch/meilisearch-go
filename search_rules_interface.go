package meilisearch

import "context"

type SearchRulesManager interface {
	SearchRulesReader
	// Update update a dynamic search rule or create a new one if it doesn't exist.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/dynamic-search-rules/update-a-dynamic-search-rule-or-create-a-new-one-if-it-doesnt-exist#body-priority-one-of-0
	Update(uid string, params *SearchRulesRequest) (*SearchRule, error)

	// UpdateWithContext update a dynamic search rule or create a new one if it doesn't exist with a context.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/dynamic-search-rules/update-a-dynamic-search-rule-or-create-a-new-one-if-it-doesnt-exist#body-priority-one-of-0
	UpdateWithContext(ctx context.Context, uid string, params *SearchRulesRequest) (*SearchRule, error)

	// Delete deletes a dynamic search rule by its unique identifier.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/dynamic-search-rules/delete-a-dynamic-search-rule
	Delete(uid string) error

	// DeleteWithContext deletes a dynamic search rule by its unique identifier with a context.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/dynamic-search-rules/delete-a-dynamic-search-rule
	DeleteWithContext(ctx context.Context, uid string) error
}

type SearchRulesReader interface {
	// List returns all dynamic search rules configured on the instance.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/dynamic-search-rules/list-dynamic-search-rules
	List(params *SearchRulesParams) (*SearchRulesResults, error)

	// ListWithContext returns all dynamic search rules configured on the instance with a context.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/dynamic-search-rules/list-dynamic-search-rules
	ListWithContext(ctx context.Context, params *SearchRulesParams) (*SearchRulesResults, error)

	// Get retrieve a single dynamic search rule by its unique identifier.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/dynamic-search-rules/get-a-dynamic-search-rule
	Get(uid string) (*SearchRule, error)

	// GetWithContext retrieve a single dynamic search rule by its unique identifier with a context.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/dynamic-search-rules/get-a-dynamic-search-rule
	GetWithContext(ctx context.Context, uid string) (*SearchRule, error)
}
