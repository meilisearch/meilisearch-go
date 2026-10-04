package meilisearch

import (
	"context"
	"time"
)

type ServiceManager interface {
	ServiceReader
	KeyManager
	TaskManager
	ChatManager
	ChatReader
	WebhookManager
	SearchRulesManager

	ServiceReader() ServiceReader

	TaskManager() TaskManager
	TaskReader() TaskReader

	KeyManager() KeyManager
	KeyReader() KeyReader

	SearchRulesManager() SearchRulesManager
	SearchRulesReader() SearchRulesReader

	ChatManager() ChatManager
	ChatReader() ChatReader

	WebhookManager() WebhookManager
	WebhookReader() WebhookReader

	// CreateIndex creates a new index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/create-index
	CreateIndex(ctx context.Context, config *IndexConfig) (*TaskInfo, error)

	// DeleteIndex deletes a specific index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/delete-index
	DeleteIndex(ctx context.Context, uid string) (*TaskInfo, error)

	// SwapIndexes swaps two existing indexes if rename is false; use rename: true if the second index does not exist.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/swap-indexes
	SwapIndexes(ctx context.Context, param []*SwapIndexesParams) (*TaskInfo, error)

	// GenerateTenantToken generates a tenant token for multi-tenancy.
	GenerateTenantToken(apiKeyUID string, searchRules map[string]interface{}, options *TenantTokenOptions) (string, error)

	// CreateDump creates a database dump.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/backups/create-dump
	CreateDump(ctx context.Context) (*TaskInfo, error)

	// CreateSnapshot create database snapshot from meilisearch
	//
	// docs: https://www.meilisearch.com/docs/reference/api/backups/create-snapshot
	CreateSnapshot(ctx context.Context) (*TaskInfo, error)

	// ExperimentalFeatures returns the experimental features manager.
	ExperimentalFeatures() *ExperimentalFeatures

	// Export transfers data from your origin instance to a remote target instance.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/export/export-to-a-remote-meilisearch
	Export(ctx context.Context, params *ExportParams) (*TaskInfo, error)

	// Experimental: UpdateNetwork updates the network object.
	//
	// 	- If leader is set to a value then UpdateNetwork will return a *[Task] object.
	// 	- If leader is not set or explicitly set to null it will return a *[Network] object.
	//	- Updates are partial; only the provided fields are updated.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/experimental-features/configure-network-topology
	UpdateNetwork(ctx context.Context, params *UpdateNetworkRequest) (any, error)

	// RenderTemplate renders any template or fragment on any input
	//
	// docs: https://www.meilisearch.com/docs/reference/api/template/render-documents-with-post
	RenderTemplate(ctx context.Context, params *RenderTemplateParams) (*RenderTemplateResponse, error)

	// Close closes the connection to the Meilisearch server.
	Close()
}

type ServiceReader interface {
	// Index retrieves an IndexManager for a specific index.
	Index(uid string) IndexManager

	// GetIndex fetches the details of a specific index and returns a *[IndexResult] object.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/get-index
	GetIndex(ctx context.Context, indexID string) (*IndexResult, error)

	// GetRawIndex fetches the raw JSON representation of a specific index and returns it as a map
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/get-index
	GetRawIndex(ctx context.Context, uid string) (map[string]interface{}, error)

	// ListIndexes lists all indexes.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/list-indexes
	ListIndexes(ctx context.Context, param *IndexesQuery) (*IndexesResults, error)

	// GetRawIndexes fetches the raw JSON representation of all indexes.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/list-indexes
	GetRawIndexes(ctx context.Context, param *IndexesQuery) (map[string]interface{}, error)

	// MultiSearch performs a multi-index search.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/search/perform-a-multi-search
	MultiSearch(ctx context.Context, queries *MultiSearchRequest) (*MultiSearchResponse, error)

	// GetStats fetches global stats.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/stats/get-stats-of-all-indexes
	GetStats(ctx context.Context, param *StatsParams) (*Stats, error)

	// Version fetches the version of the Meilisearch server.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/version/get-version
	Version(ctx context.Context) (*Version, error)

	// Health checks the health of the Meilisearch server.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/health/get-health
	Health(ctx context.Context) (*Health, error)

	// IsHealthy checks if the Meilisearch server is healthy.
	IsHealthy() bool

	// GetBatches allows you to monitor how Meilisearch is grouping and processing asynchronous operations.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/async-task-management/list-batches
	GetBatches(ctx context.Context, param *BatchesQuery) (*BatchesResults, error)

	// GetBatch retrieves a specific batch by its UID.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/async-task-management/get-batch
	GetBatch(ctx context.Context, batchUID int) (*Batch, error)

	// Experimental: GetNetwork gets the current value of the instance’s network object.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/experimental-features/get-network-topology#get-network-topology
	GetNetwork(ctx context.Context) (*Network, error)
}

type KeyManager interface {
	KeyReader

	// CreateKey creates a new API key and returns the details of the created [Key]
	//
	// docs: https://www.meilisearch.com/docs/reference/api/keys/create-api-key
	CreateKey(ctx context.Context, request *Key) (*Key, error)

	// UpdateKey updates a specific API key.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/keys/update-api-key
	UpdateKey(ctx context.Context, keyOrUID string, request *Key) (*Key, error)

	// DeleteKey deletes a specific API key.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/keys/delete-api-key
	DeleteKey(ctx context.Context, keyOrUID string) (bool, error)
}

type KeyReader interface {
	// GetKey fetches the details of a specific API key.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/keys/get-api-key
	GetKey(ctx context.Context, identifier string) (*Key, error)

	// GetKeys lists all API keys.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/keys/list-api-keys
	GetKeys(ctx context.Context, param *KeysQuery) (*KeysResults, error)
}

type TaskManager interface {
	TaskReader

	// CancelTasks cancels specific tasks.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/async-task-management/cancel-tasks
	CancelTasks(ctx context.Context, param *CancelTasksQuery) (*TaskInfo, error)

	// DeleteTasks deletes specific tasks.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/async-task-management/delete-tasks
	DeleteTasks(ctx context.Context, param *DeleteTasksQuery) (*TaskInfo, error)

	// GetTaskDocuments retrieves the documents associated with a task (added, updated, or deleted).
	// dst must be a non-nil pointer to a slice that receives the streamed NDJSON documents.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/async-task-management/get-tasks-documents
	GetTaskDocuments(ctx context.Context, taskUID int64, dst interface{}) error
}

type TaskReader interface {
	// GetTask retrieves a task by its UID.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/async-task-management/get-task
	GetTask(ctx context.Context, taskUID int64) (*Task, error)

	// GetTasks retrieves multiple tasks based on query parameters.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/async-task-management/list-tasks
	GetTasks(ctx context.Context, param *TasksQuery) (*TaskResult, error)

	// WaitForTask waits for a task to complete by its UID with the given interval.
	WaitForTask(ctx context.Context, taskUID int64, interval time.Duration) (*Task, error)
}
