package meilisearch

import (
	"testing"
)

func TestClientKeys_Get(t *testing.T) {
	var client = NewClient(Config{
		Host:   "http://localhost:7700",
		APIKey: "masterKey",
	})

	if _, err := client.Keys().Get(); err != nil {
		t.Fatal(err)
	}
}
