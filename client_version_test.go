package meilisearch

import (
	"testing"
)

func TestClientVersion_Get(t *testing.T) {
	if _, err := client.Version().Get(); err != nil {
		t.Fatal(err)
	}
}
