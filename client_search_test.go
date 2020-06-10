package meilisearch

import (
	"fmt"
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
		AddOrUpdate([]docTestBooks{
			{Book_id: 123, Title: "Pride and Prejudice", Tag: "Nice book"},
			{Book_id: 456, Title: "Le Petit Prince", Tag: "Nice book"},
			{Book_id: 1, Title: "Alice In Wonderland", Tag: "Nice book"},
			{Book_id: 1344, Title: "The Hobbit", Tag: "Nice book"},
			{Book_id: 4, Title: "Harry Potter and the Half-Blood Prince", Tag: "Interesting book"},
			{Book_id: 42, Title: "The Hitchhiker's Guide to the Galaxy", Tag: "Interesting book"},
			{Book_id: 24, Title: "You are a princess", Tag: "Interesting book"},
		})

	if err != nil {
		t.Fatal(err)
	}

	client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)

	// Test basic search

	resp, err := client.Search(indexUID).Search(SearchRequest{
		Query: "prince",
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Hits) != 3 {
		fmt.Println(resp)
		t.Fatal("number of hits should be equal to 3")
	}

	// Test basic search with limit

	resp, err = client.Search(indexUID).Search(SearchRequest{
		Query: "prince",
		Limit: 1,
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Hits) != 1 {
		fmt.Println(resp)
		t.Fatal("number of hits should be equal to 1")
	}
	title := resp.Hits[0].(map[string]interface{})["title"]
	if title != "Le Petit Prince" {
		fmt.Println(resp)
		t.Fatal("Should have found: Le Petit Prince")
	}

	// Test basic search with offset

	resp, err = client.Search(indexUID).Search(SearchRequest{
		Query:  "prince",
		Offset: 1,
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Hits) != 2 {
		fmt.Println(resp)
		t.Fatal("number of hits should be equal to 2")
	}

}
