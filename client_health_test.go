package meilisearch

import "testing"

func TestClientHealth_Get(t *testing.T) {
	if err := client.Health().Get(); err != nil {
		t.Fatal(err)
	}
}
