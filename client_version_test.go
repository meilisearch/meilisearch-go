package meilisearch

import (
	"testing"
)

func TestClientVersion_Get(t *testing.T) {
	var client = NewClient(Config{
		Host:   "http://localhost:7700",
		APIKey: "masterKey",
	})

	if _, err := client.Version().Get(); err != nil {
		t.Fatal(err)
	}
}
