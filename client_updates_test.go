package meilisearch

import (
	"testing"
)

func TestClientUpdates_List(t *testing.T) {
	var indexUID = "TestClientUpdates_List"

	var client = NewClient(Config{
		Host:   "http://localhost:7700",
		APIKey: "masterKey",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, err := client.Updates(indexUID).List(); err != nil {
		t.Fatal(err)
	}
}
