package meilisearch

import (
	"testing"
)

func TestClientStats_Get(t *testing.T) {
	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	if _, err := client.Stats().Get("stats_tests"); err != nil {
		return
	}
}

func TestClientStats_List(t *testing.T) {
	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	if _, err := client.Stats().List(); err != nil {
		return
	}
}
