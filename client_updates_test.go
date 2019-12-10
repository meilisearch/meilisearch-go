package meilisearch

import (
	"log"
	"testing"
)

var updates clientUpdates

func TestClientUpdates_Get(t *testing.T) {
	resp, err := newClientIndexes(updates.client).UpdateSchema("updates_tests", Schema{
		"id":   {"identifier", "indexed", "displayed"},
		"name": {"indexed", "indexed", "displayed"},
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, err := updates.Get(resp.UpdateID); err != nil {
		t.Fatal(err)
	}
}

func TestClientUpdates_List(t *testing.T) {
	if _, err := updates.List(); err != nil {
		t.Fatal(err)
	}
}

func init() {
	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	resp, err := newClientIndexes(client).Create(CreateIndexRequest{
		Name: "updates_tests",
		UID:  "updates_tests",
		Schema: Schema{
			"id":   {"identifier", "indexed", "displayed"},
			"name": {"indexed", "displayed"},
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	updates = newClientUpdates(client, resp.UID)
}
