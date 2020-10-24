package meilisearch

import (
	"context"
	"time"
)

// Config configure the Client
type Config struct {

	// Host is the host of your meilisearch database
	// Example: 'http://localhost:7700'
	Host string

	// APIKey is optional
	APIKey string
}

// ClientInterface is interface for all Meilisearch client
type ClientInterface interface {
	WaitForPendingUpdate(ctx context.Context, interval time.Duration, indexID string, updateID *AsyncUpdateID) (UpdateStatus, error)
	DefaultWaitForPendingUpdate(indexUID string, updateID *AsyncUpdateID) (UpdateStatus, error)

	Indexes() APIIndexes
	Version() APIVersion
	Documents(indexID string) APIDocuments
	Search(indexID string) APISearch
	Updates(indexID string) APIUpdates
	Settings(indexID string) APISettings
	Keys() APIKeys
	Stats() APIStats
	Health() APIHealth
}
