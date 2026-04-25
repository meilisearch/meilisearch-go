package meilisearch

import "context"

type SearchRulesManager interface {
	SearchRulesReader

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
	// Get retrieve a single dynamic search rule by its unique identifier.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/dynamic-search-rules/get-a-dynamic-search-rule
	Get(uid string) (*SearchRule, error)

	// GetWithContext retrieve a single dynamic search rule by its unique identifier with a context.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/dynamic-search-rules/get-a-dynamic-search-rule
	GetWithContext(ctx context.Context, uid string) (*SearchRule, error)
}
