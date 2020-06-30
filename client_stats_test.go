package meilisearch

import (
	"testing"
)

func TestClientStats_Get(t *testing.T) {
	var client = NewClient(Config{
		Host:   "http://localhost:7700",
		APIKey: "masterKey",
	})

	if _, err := client.Stats().GetAll(); err != nil {
		t.Fatal(err)
	}

	var indexUID = "TestClientStats_Get"

	if _, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := client.Stats().Get(indexUID); err != nil {
		t.Fatal(err)
	}
}
