package meilisearch

import "context"

type SearchRulesManager interface {
	SearchRulesReader
	// UpdateSearchRule update a dynamic search rule or create a new one if it doesn't exist.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/search-rules/create-or-update-a-search-rule
	UpdateSearchRule(ctx context.Context, uid string, params *SearchRulesRequest) (*Task, error)

	// DeleteSearchRule deletes a dynamic search rule by its unique identifier.
	// If uid is nil, all dynamic search rules will be deleted.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/dynamic-search-rules/delete-a-dynamic-search-rule
	DeleteSearchRule(ctx context.Context, uid *string) (*Task, error)
}

type SearchRulesReader interface {
	// ListSearchRules returns all dynamic search rules configured on the instance.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/dynamic-search-rules/list-dynamic-search-rules
	ListSearchRules(ctx context.Context, params *SearchRulesParams) (*SearchRulesResults, error)

	// GetSearchRule retrieve a single dynamic search rule by its unique identifier.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/dynamic-search-rules/get-a-dynamic-search-rule
	GetSearchRule(ctx context.Context, uid string) (*SearchRule, error)
}
