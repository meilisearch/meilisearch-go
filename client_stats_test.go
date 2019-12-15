package meilisearch

import "testing"

var stats clientStats

func TestClientStats_Get(t *testing.T) {
	if _, err := stats.Get(); err != nil {
		return
	}
}

func TestClientStats_List(t *testing.T) {
	if _, err := stats.List(); err != nil {
		return
	}
}

func init() {
	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	resp, err := newClientIndexes(client).Create(CreateIndexRequest{
		Name: "stats_tests",
		UID:  "stats_tests",
		Schema: Schema{
			"id":   {"identifier", "indexed", "displayed"},
			"name": {"indexed", "indexed", "displayed"},
		},
	})

	if err != nil {
		panic(err)
	}

	stats = newClientStats(client, resp.UID)
}
