package meilisearch

import "context"

type SearchRulesManager interface {
	SearchRulesReader
	// UpdateSearchRule update a dynamic search rule or create a new one if it doesn't exist.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/dynamic-search-rules/update-a-dynamic-search-rule-or-create-a-new-one-if-it-doesnt-exist#body-priority-one-of-0
	UpdateSearchRule(uid string, params *SearchRulesRequest) (*SearchRule, error)

	// UpdateSearchRuleWithContext update a dynamic search rule or create a new one if it doesn't exist with a context.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/dynamic-search-rules/update-a-dynamic-search-rule-or-create-a-new-one-if-it-doesnt-exist#body-priority-one-of-0
	UpdateSearchRuleWithContext(ctx context.Context, uid string, params *SearchRulesRequest) (*SearchRule, error)

	// DeleteSearchRule deletes a dynamic search rule by its unique identifier.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/dynamic-search-rules/delete-a-dynamic-search-rule
	DeleteSearchRule(uid string) error

	// DeleteSearchRuleWithContext deletes a dynamic search rule by its unique identifier with a context.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/dynamic-search-rules/delete-a-dynamic-search-rule
	DeleteSearchRuleWithContext(ctx context.Context, uid string) error
}

type SearchRulesReader interface {
	// ListSearchRules returns all dynamic search rules configured on the instance.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/dynamic-search-rules/list-dynamic-search-rules
	ListSearchRules(params *SearchRulesParams) (*SearchRulesResults, error)

	// ListSearchRulesWithContext returns all dynamic search rules configured on the instance with a context.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/dynamic-search-rules/list-dynamic-search-rules
	ListSearchRulesWithContext(ctx context.Context, params *SearchRulesParams) (*SearchRulesResults, error)

	// GetSearchRule retrieve a single dynamic search rule by its unique identifier.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/dynamic-search-rules/get-a-dynamic-search-rule
	GetSearchRule(uid string) (*SearchRule, error)

	// GetSearchRuleWithContext retrieve a single dynamic search rule by its unique identifier with a context.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/dynamic-search-rules/get-a-dynamic-search-rule
	GetSearchRuleWithContext(ctx context.Context, uid string) (*SearchRule, error)
}
