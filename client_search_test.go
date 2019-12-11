package meilisearch

import (
	"testing"
)

var search clientSearch

func TestClientSearch_Search(t *testing.T) {
	resp, err := search.Search(SearchRequest{
		Query: "citrons",
		Limit: 10,
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Hits) != 2 {
		t.Fatal("number of hits should be equal to 2")
	}
}

func init() {
	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	resp, err := newClientIndexes(client).Create(CreateIndexRequest{
		Name: "search_tests",
		UID:  "search_tests",
		Schema: Schema{
			"id":   {"identifier", "indexed", "displayed"},
			"name": {"indexed", "indexed", "displayed"},
		},
	})

	if err != nil {
		panic(err)
	}

	documents := newClientDocuments(client, resp.UID)

	u, err := documents.AddOrUpdate(&[]docTest{
		{"0", "J'adore les citrons"},
		{"1", "Les citrons c'est la vie"},
		{"2", "Les ponchos c'est bien !"},
	})
	if err != nil {
		panic(err)
	}

	AwaitAsyncUpdateId(documents, u)

	search = newClientSearch(client, resp.UID)
}
