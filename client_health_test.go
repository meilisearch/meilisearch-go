package meilisearch

import "testing"

var health clientHealth

func TestClientHealth_Get(t *testing.T) {
	if err := health.Get(); err != nil {
		return
	}
}

func TestClientHealth_Set(t *testing.T) {
	if err := health.Set(true); err != nil {
		return
	}
}

func init() {
	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	health = newClientHealth(client)
}
