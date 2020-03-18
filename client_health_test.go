package meilisearch

import "testing"

func TestClientHealth_Get(t *testing.T) {
	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	if err := client.Health().Get(); err != nil {
		return
	}
}

func TestClientHealth_Set(t *testing.T) {
	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	if err := client.Health().Set(true); err != nil {
		return
	}
}
