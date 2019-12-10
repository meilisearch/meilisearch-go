package meilisearch

import (
	"testing"
)

var version clientVersion

func TestClientVersion_Get(t *testing.T) {
	if _, err := version.Get(); err != nil {
		t.Fatal(err)
	}
}

func init() {
	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	version = newClientVersion(client)
}
