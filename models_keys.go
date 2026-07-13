package meilisearch

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Key allow the user to connect to the meilisearch instance
//
// Documentation: https://www.meilisearch.com/docs/learn/security/master_api_keys#protecting-a-meilisearch-instance
type Key struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Key         string    `json:"key,omitempty"`
	UID         string    `json:"uid,omitempty"`
	Actions     []string  `json:"actions,omitempty"`
	Indexes     []string  `json:"indexes,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

// KeyParsed this structure is used to send the exact ISO-8601 time format managed by meilisearch
type KeyParsed struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	UID         string   `json:"uid,omitempty"`
	Actions     []string `json:"actions,omitempty"`
	Indexes     []string `json:"indexes,omitempty"`
	ExpiresAt   *string  `json:"expiresAt"`
}

// KeyUpdate this structure is used to update a Key
type KeyUpdate struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// KeysResults return of multiple keys is wrap in a KeysResults
type KeysResults struct {
	Results []Key `json:"results"`
	Offset  int64 `json:"offset"`
	Limit   int64 `json:"limit"`
	Total   int64 `json:"total"`
}

type KeysQuery struct {
	Limit  int64
	Offset int64
}

// TenantTokenOptions information to create a tenant token
//
// ExpiresAt is a time.Time when the key will expire.
// Note that if an ExpiresAt value is included it should be in UTC time.
// ApiKey is the API key parent of the token.
type TenantTokenOptions struct {
	APIKey    string
	ExpiresAt time.Time
}

// TenantTokenClaims custom Claims structure to create a Tenant Token
type TenantTokenClaims struct {
	APIKeyUID   string      `json:"apiKeyUid"`
	SearchRules interface{} `json:"searchRules"`
	jwt.RegisteredClaims
}
