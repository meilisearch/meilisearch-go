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
	CreateIndex(config *IndexConfig) (*TaskInfo, error)

	// CreateIndexWithContext creates a new index with a context for cancellation.
	CreateIndexWithContext(ctx context.Context, config *IndexConfig) (*TaskInfo, error)

	// DeleteIndex deletes a specific index.
	DeleteIndex(uid string) (*TaskInfo, error)

	// DeleteIndexWithContext deletes a specific index with a context for cancellation.
	DeleteIndexWithContext(ctx context.Context, uid string) (*TaskInfo, error)

	// SwapIndexes swaps two existing indexes if rename is false; use rename: true if the second index does not exist.
	SwapIndexes(param []*SwapIndexesParams) (*TaskInfo, error)

	// SwapIndexesWithContext swaps two existing indexes with a context if rename is false; use rename: true if the second index does not exist.
	SwapIndexesWithContext(ctx context.Context, param []*SwapIndexesParams) (*TaskInfo, error)

	// GenerateTenantToken generates a tenant token for multi-tenancy.
	GenerateTenantToken(apiKeyUID string, searchRules map[string]interface{}, options *TenantTokenOptions) (string, error)

	// CreateDump creates a database dump.
	CreateDump() (*TaskInfo, error)

	// CreateDumpWithContext creates a database dump with a context for cancellation.
	CreateDumpWithContext(ctx context.Context) (*TaskInfo, error)

	// CreateSnapshot create database snapshot from meilisearch
	CreateSnapshot() (*TaskInfo, error)

	// CreateSnapshotWithContext create database snapshot from meilisearch and support parent context
	CreateSnapshotWithContext(ctx context.Context) (*TaskInfo, error)

	// ExperimentalFeatures returns the experimental features manager.
	ExperimentalFeatures() *ExperimentalFeatures

	// Export transfers data from your origin instance to a remote target instance.
	Export(params *ExportParams) (*TaskInfo, error)

	// ExportWithContext transfers data from your origin instance to a remote target instance with a context for cancellation.
	ExportWithContext(ctx context.Context, params *ExportParams) (*TaskInfo, error)

	// UpdateNetwork updates the network object.
	// Updates are partial; only the provided fields are updated.
	UpdateNetwork(params *Network) (*Network, error)

	// UpdateNetworkWithContext updates the network object with a context.
	// Updates are partial; only the provided fields are updated.
	UpdateNetworkWithContext(ctx context.Context, params *Network) (*Network, error)

	// Close closes the connection to the Meilisearch server.
	Close()
}

type ServiceReader interface {
	// Index retrieves an IndexManager for a specific index.
	Index(uid string) IndexManager

	// GetIndex fetches the details of a specific index.
	GetIndex(indexID string) (*IndexResult, error)

	// GetIndexWithContext fetches the details of a specific index with a context for cancellation.
	GetIndexWithContext(ctx context.Context, indexID string) (*IndexResult, error)

	// GetRawIndex fetches the raw JSON representation of a specific index.
	GetRawIndex(uid string) (map[string]interface{}, error)

	// GetRawIndexWithContext fetches the raw JSON representation of a specific index with a context for cancellation.
	GetRawIndexWithContext(ctx context.Context, uid string) (map[string]interface{}, error)

	// ListIndexes lists all indexes.
	ListIndexes(param *IndexesQuery) (*IndexesResults, error)

	// ListIndexesWithContext lists all indexes with a context for cancellation.
	ListIndexesWithContext(ctx context.Context, param *IndexesQuery) (*IndexesResults, error)

	// GetRawIndexes fetches the raw JSON representation of all indexes.
	GetRawIndexes(param *IndexesQuery) (map[string]interface{}, error)

	// GetRawIndexesWithContext fetches the raw JSON representation of all indexes with a context for cancellation.
	GetRawIndexesWithContext(ctx context.Context, param *IndexesQuery) (map[string]interface{}, error)

	// MultiSearch performs a multi-index search.
	MultiSearch(queries *MultiSearchRequest) (*MultiSearchResponse, error)

	// MultiSearchWithContext performs a multi-index search with a context for cancellation.
	MultiSearchWithContext(ctx context.Context, queries *MultiSearchRequest) (*MultiSearchResponse, error)

	// GetStats fetches global stats.
	GetStats() (*Stats, error)

	// GetStatsWithContext fetches global stats with a context for cancellation.
	GetStatsWithContext(ctx context.Context) (*Stats, error)

	// Version fetches the version of the Meilisearch server.
	Version() (*Version, error)

	// VersionWithContext fetches the version of the Meilisearch server with a context for cancellation.
	VersionWithContext(ctx context.Context) (*Version, error)

	// Health checks the health of the Meilisearch server.
	Health() (*Health, error)

	// HealthWithContext checks the health of the Meilisearch server with a context for cancellation.
	HealthWithContext(ctx context.Context) (*Health, error)

	// IsHealthy checks if the Meilisearch server is healthy.
	IsHealthy() bool

	// GetBatches allows you to monitor how Meilisearch is grouping and processing asynchronous operations.
	GetBatches(param *BatchesQuery) (*BatchesResults, error)

	// GetBatchesWithContext allows you to monitor how Meilisearch is grouping and processing asynchronous operations with a context for cancellation.
	GetBatchesWithContext(ctx context.Context, param *BatchesQuery) (*BatchesResults, error)

	// GetBatch retrieves a specific batch by its UID.
	GetBatch(batchUID int) (*Batch, error)

	// GetBatchWithContext retrieves a specific batch by its UID with a context for cancellation.
	GetBatchWithContext(ctx context.Context, batchUID int) (*Batch, error)

	// GetNetwork gets the current value of the instance’s network object.
	GetNetwork() (*Network, error)

	// GetNetworkWithContext gets the current value of the instance’s network object with a context.
	GetNetworkWithContext(ctx context.Context) (*Network, error)
}

type KeyManager interface {
	KeyReader

	// CreateKey creates a new API key.
	CreateKey(request *Key) (*Key, error)

	// CreateKeyWithContext creates a new API key with a context for cancellation.
	CreateKeyWithContext(ctx context.Context, request *Key) (*Key, error)

	// UpdateKey updates a specific API key.
	UpdateKey(keyOrUID string, request *Key) (*Key, error)

	// UpdateKeyWithContext updates a specific API key with a context for cancellation.
	UpdateKeyWithContext(ctx context.Context, keyOrUID string, request *Key) (*Key, error)

	// DeleteKey deletes a specific API key.
	DeleteKey(keyOrUID string) (bool, error)

	// DeleteKeyWithContext deletes a specific API key with a context for cancellation.
	DeleteKeyWithContext(ctx context.Context, keyOrUID string) (bool, error)
}

type KeyReader interface {
	// GetKey fetches the details of a specific API key.
	GetKey(identifier string) (*Key, error)

	// GetKeyWithContext fetches the details of a specific API key with a context for cancellation.
	GetKeyWithContext(ctx context.Context, identifier string) (*Key, error)

	// GetKeys lists all API keys.
	GetKeys(param *KeysQuery) (*KeysResults, error)

	// GetKeysWithContext lists all API keys with a context for cancellation.
	GetKeysWithContext(ctx context.Context, param *KeysQuery) (*KeysResults, error)
}

type TaskManager interface {
	TaskReader

	// CancelTasks cancels specific tasks.
	CancelTasks(param *CancelTasksQuery) (*TaskInfo, error)

	// CancelTasksWithContext cancels specific tasks with a context for cancellation.
	CancelTasksWithContext(ctx context.Context, param *CancelTasksQuery) (*TaskInfo, error)

	// DeleteTasks deletes specific tasks.
	DeleteTasks(param *DeleteTasksQuery) (*TaskInfo, error)

	// DeleteTasksWithContext deletes specific tasks with a context for cancellation.
	DeleteTasksWithContext(ctx context.Context, param *DeleteTasksQuery) (*TaskInfo, error)
}

type TaskReader interface {
	// GetTask retrieves a task by its UID.
	GetTask(taskUID int64) (*Task, error)

	// GetTaskWithContext retrieves a task by its UID using the provided context for cancellation.
	GetTaskWithContext(ctx context.Context, taskUID int64) (*Task, error)

	// GetTasks retrieves multiple tasks based on query parameters.
	GetTasks(param *TasksQuery) (*TaskResult, error)

	// GetTasksWithContext retrieves multiple tasks based on query parameters using the provided context for cancellation.
	GetTasksWithContext(ctx context.Context, param *TasksQuery) (*TaskResult, error)

	// WaitForTask waits for a task to complete by its UID with the given interval.
	WaitForTask(taskUID int64, interval time.Duration) (*Task, error)

	// WaitForTaskWithContext waits for a task to complete by its UID with the given interval using the provided context for cancellation.
	WaitForTaskWithContext(ctx context.Context, taskUID int64, interval time.Duration) (*Task, error)
}
