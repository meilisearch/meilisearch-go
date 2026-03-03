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

	ServiceReader() ServiceReader

	TaskManager() TaskManager
	TaskReader() TaskReader

	KeyManager() KeyManager
	KeyReader() KeyReader

	ChatManager() ChatManager
	ChatReader() ChatReader

	WebhookManager() WebhookManager
	WebhookReader() WebhookReader

	// CreateIndex creates a new index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/create-index
	CreateIndex(config *IndexConfig) (*TaskInfo, error)

	// CreateIndexWithContext creates a new index with a context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/create-index
	CreateIndexWithContext(ctx context.Context, config *IndexConfig) (*TaskInfo, error)

	// DeleteIndex deletes a specific index.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/delete-index
	DeleteIndex(uid string) (*TaskInfo, error)

	// DeleteIndexWithContext deletes a specific index with a context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/delete-index
	DeleteIndexWithContext(ctx context.Context, uid string) (*TaskInfo, error)

	// SwapIndexes swaps two existing indexes if rename is false; use rename: true if the second index does not exist.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/swap-indexes
	SwapIndexes(param []*SwapIndexesParams) (*TaskInfo, error)

	// SwapIndexesWithContext swaps two existing indexes with a context if rename is false; use rename: true if the second index does not exist.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/swap-indexes
	SwapIndexesWithContext(ctx context.Context, param []*SwapIndexesParams) (*TaskInfo, error)

	// GenerateTenantToken generates a tenant token for multi-tenancy.
	GenerateTenantToken(apiKeyUID string, searchRules map[string]interface{}, options *TenantTokenOptions) (string, error)

	// CreateDump creates a database dump.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/backups/create-dump
	CreateDump() (*TaskInfo, error)

	// CreateDumpWithContext creates a database dump with a context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/backups/create-dump
	CreateDumpWithContext(ctx context.Context) (*TaskInfo, error)

	// CreateSnapshot create database snapshot from meilisearch
	//
	// docs: https://www.meilisearch.com/docs/reference/api/backups/create-snapshot
	CreateSnapshot() (*TaskInfo, error)

	// CreateSnapshotWithContext create database snapshot from meilisearch and support parent context
	//
	// docs: https://www.meilisearch.com/docs/reference/api/backups/create-snapshot
	CreateSnapshotWithContext(ctx context.Context) (*TaskInfo, error)

	// ExperimentalFeatures returns the experimental features manager.
	ExperimentalFeatures() *ExperimentalFeatures

	// Export transfers data from your origin instance to a remote target instance.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/export/export-to-a-remote-meilisearch
	Export(params *ExportParams) (*TaskInfo, error)

	// ExportWithContext transfers data from your origin instance to a remote target instance with a context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/export/export-to-a-remote-meilisearch
	ExportWithContext(ctx context.Context, params *ExportParams) (*TaskInfo, error)

	// Experimental: UpdateNetwork updates the network object.
	//
	// 	- If leader is set to a value then UpdateNetwork will return a *[Task] object.
	// 	- If leader is not set or explicitly set to null it will return a *[Network] object.
	//	- Updates are partial; only the provided fields are updated.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/experimental-features/configure-network-topology
	UpdateNetwork(params *UpdateNetworkRequest) (any, error)

	// Experimental: UpdateNetworkWithContext updates the network object with a context.
	// 	- If leader is set to a value then UpdateNetwork will return a *[Task] object.
	// 	- If leader is not set or explicitly set to null it will return a *[Network] object.
	// 	- Updates are partial; only the provided fields are updated.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/experimental-features/configure-network-topology
	UpdateNetworkWithContext(ctx context.Context, params *UpdateNetworkRequest) (any, error)

	// Close closes the connection to the Meilisearch server.
	Close()
}

type ServiceReader interface {
	// Index retrieves an IndexManager for a specific index.
	Index(uid string) IndexManager

	// GetIndex fetches the details of a specific index and returns a *[IndexResult] object.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/get-index
	GetIndex(indexID string) (*IndexResult, error)

	// GetIndexWithContext fetches the details of a specific index and returns a *[IndexResult] object, with a context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/get-index
	GetIndexWithContext(ctx context.Context, indexID string) (*IndexResult, error)

	// GetRawIndex fetches the raw JSON representation of a specific index and returns it as a map
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/get-index
	GetRawIndex(uid string) (map[string]interface{}, error)

	// GetRawIndexWithContext fetches the raw JSON representation of a specific index and returns it as a map, with a context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/get-index
	GetRawIndexWithContext(ctx context.Context, uid string) (map[string]interface{}, error)

	// ListIndexes lists all indexes.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/list-indexes
	ListIndexes(param *IndexesQuery) (*IndexesResults, error)

	// ListIndexesWithContext lists all indexes with a context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/list-indexes
	ListIndexesWithContext(ctx context.Context, param *IndexesQuery) (*IndexesResults, error)

	// GetRawIndexes fetches the raw JSON representation of all indexes.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/list-indexes
	GetRawIndexes(param *IndexesQuery) (map[string]interface{}, error)

	// GetRawIndexesWithContext fetches the raw JSON representation of all indexes with a context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/indexes/list-indexes
	GetRawIndexesWithContext(ctx context.Context, param *IndexesQuery) (map[string]interface{}, error)

	// MultiSearch performs a multi-index search.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/search/perform-a-multi-search
	MultiSearch(queries *MultiSearchRequest) (*MultiSearchResponse, error)

	// MultiSearchWithContext performs a multi-index search with a context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/search/perform-a-multi-search
	MultiSearchWithContext(ctx context.Context, queries *MultiSearchRequest) (*MultiSearchResponse, error)

	// GetStats fetches global stats.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/stats/get-stats-of-all-indexes
	GetStats() (*Stats, error)

	// GetStatsWithContext fetches global stats with a context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/stats/get-stats-of-all-indexes
	GetStatsWithContext(ctx context.Context) (*Stats, error)

	// Version fetches the version of the Meilisearch server.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/version/get-version
	Version() (*Version, error)

	// VersionWithContext fetches the version of the Meilisearch server with a context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/version/get-version
	VersionWithContext(ctx context.Context) (*Version, error)

	// Health checks the health of the Meilisearch server.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/health/get-health
	Health() (*Health, error)

	// HealthWithContext checks the health of the Meilisearch server with a context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/health/get-health
	HealthWithContext(ctx context.Context) (*Health, error)

	// IsHealthy checks if the Meilisearch server is healthy.
	IsHealthy() bool

	// GetBatches allows you to monitor how Meilisearch is grouping and processing asynchronous operations.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/async-task-management/list-batches
	GetBatches(param *BatchesQuery) (*BatchesResults, error)

	// GetBatchesWithContext allows you to monitor how Meilisearch is grouping and processing asynchronous operations with a context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/async-task-management/list-batches
	GetBatchesWithContext(ctx context.Context, param *BatchesQuery) (*BatchesResults, error)

	// GetBatch retrieves a specific batch by its UID.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/async-task-management/get-batch
	GetBatch(batchUID int) (*Batch, error)

	// GetBatchWithContext retrieves a specific batch by its UID with a context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/async-task-management/get-batch
	GetBatchWithContext(ctx context.Context, batchUID int) (*Batch, error)

	// Experimental: GetNetwork gets the current value of the instance’s network object.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/experimental-features/get-network-topology#get-network-topology
	GetNetwork() (*Network, error)

	// Experimental: GetNetworkWithContext gets the current value of the instance’s network object with a context.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/experimental-features/get-network-topology#get-network-topology
	GetNetworkWithContext(ctx context.Context) (*Network, error)
}

type KeyManager interface {
	KeyReader

	// CreateKey creates a new API key and returns the details of the created [Key]
	//
	// docs: https://www.meilisearch.com/docs/reference/api/keys/create-api-key
	CreateKey(request *Key) (*Key, error)

	// CreateKeyWithContext creates a new API key and returns the details of the created [Key], with a context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/keys/create-api-key
	CreateKeyWithContext(ctx context.Context, request *Key) (*Key, error)

	// UpdateKey updates a specific API key.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/keys/update-api-key
	UpdateKey(keyOrUID string, request *Key) (*Key, error)

	// UpdateKeyWithContext updates a specific API key with a context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/keys/update-api-key
	UpdateKeyWithContext(ctx context.Context, keyOrUID string, request *Key) (*Key, error)

	// DeleteKey deletes a specific API key.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/keys/delete-api-key
	DeleteKey(keyOrUID string) (bool, error)

	// DeleteKeyWithContext deletes a specific API key with a context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/keys/delete-api-key
	DeleteKeyWithContext(ctx context.Context, keyOrUID string) (bool, error)
}

type KeyReader interface {
	// GetKey fetches the details of a specific API key.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/keys/get-api-key
	GetKey(identifier string) (*Key, error)

	// GetKeyWithContext fetches the details of a specific API key with a context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/keys/get-api-key
	GetKeyWithContext(ctx context.Context, identifier string) (*Key, error)

	// GetKeys lists all API keys.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/keys/list-api-keys
	GetKeys(param *KeysQuery) (*KeysResults, error)

	// GetKeysWithContext lists all API keys with a context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/keys/list-api-keys
	GetKeysWithContext(ctx context.Context, param *KeysQuery) (*KeysResults, error)
}

type TaskManager interface {
	TaskReader

	// CancelTasks cancels specific tasks.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/async-task-management/cancel-tasks
	CancelTasks(param *CancelTasksQuery) (*TaskInfo, error)

	// CancelTasksWithContext cancels specific tasks with a context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/async-task-management/cancel-tasks
	CancelTasksWithContext(ctx context.Context, param *CancelTasksQuery) (*TaskInfo, error)

	// DeleteTasks deletes specific tasks.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/async-task-management/delete-tasks
	DeleteTasks(param *DeleteTasksQuery) (*TaskInfo, error)

	// DeleteTasksWithContext deletes specific tasks with a context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/async-task-management/delete-tasks
	DeleteTasksWithContext(ctx context.Context, param *DeleteTasksQuery) (*TaskInfo, error)
}

type TaskReader interface {
	// GetTask retrieves a task by its UID.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/async-task-management/get-task
	GetTask(taskUID int64) (*Task, error)

	// GetTaskWithContext retrieves a task by its UID using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/async-task-management/get-task
	GetTaskWithContext(ctx context.Context, taskUID int64) (*Task, error)

	// GetTasks retrieves multiple tasks based on query parameters.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/async-task-management/list-tasks
	GetTasks(param *TasksQuery) (*TaskResult, error)

	// GetTasksWithContext retrieves multiple tasks based on query parameters using the provided context for cancellation.
	//
	// docs: https://www.meilisearch.com/docs/reference/api/async-task-management/list-tasks
	GetTasksWithContext(ctx context.Context, param *TasksQuery) (*TaskResult, error)

	// WaitForTask waits for a task to complete by its UID with the given interval.
	WaitForTask(taskUID int64, interval time.Duration) (*Task, error)

	// WaitForTaskWithContext waits for a task to complete by its UID with the given interval using the provided context for cancellation.
	WaitForTaskWithContext(ctx context.Context, taskUID int64, interval time.Duration) (*Task, error)
}
