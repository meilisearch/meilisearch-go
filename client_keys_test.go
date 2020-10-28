package meilisearch

import (
	"testing"
)

func TestClientKeys_Get(t *testing.T) {
	if _, err := client.Keys().Get(); err != nil {
		t.Fatal(err)
	}
}
