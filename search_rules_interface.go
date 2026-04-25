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
}
