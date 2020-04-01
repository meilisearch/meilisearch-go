package meilisearch

import (
	"testing"
)

func TestClientSearch_Search(t *testing.T) {
	var indexUID = "TestClientSearch_Search"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.
		Documents(indexUID).
		AddOrUpdate([]docTest{
			{"0", "J'adore les citrons"},
			{"1", "Les citrons c'est la vie"},
			{"2", "Les ponchos c'est bien !"},
		})

	if err != nil {
		t.Fatal(err)
	}

	client.defaultWaitForPendingUpdate(indexUID, updateIDRes)

	resp, err := client.Search(indexUID).Search(SearchRequest{
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
