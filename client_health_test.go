package meilisearch

import "testing"

func TestClientHealth_Get(t *testing.T) {
	var client = NewClient(Config{
		Host:   "http://localhost:7700",
		APIKey: "masterKey",
	})

	if err := client.Health().Get(); err != nil {
		t.Fatal(err)
	}
}

func TestClientHealth_Set(t *testing.T) {
	var client = NewClient(Config{
		Host:   "http://localhost:7700",
		APIKey: "masterKey",
	})

	if err := client.Health().Update(true); err != nil {
		t.Fatal(err)
	}
}
