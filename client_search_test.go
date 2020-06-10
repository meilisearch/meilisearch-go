package meilisearch

import (
	"fmt"
	"reflect"
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

	// Test basic search with attributesToRetrieve

	resp, err = client.Search(indexUID).Search(SearchRequest{
		Query:                "prince",
		AttributesToRetrieve: []string{"book_id", "title"},
	})

	if err != nil {
		t.Fatal(err)
	}

	if resp.Hits[0].(map[string]interface{})["title"] == nil {
		fmt.Println(resp)
		t.Fatal("attributesToRetrieve: Couldn't retrieve field in response")
	}
	if resp.Hits[0].(map[string]interface{})["tag"] != nil {
		fmt.Println(resp)
		t.Fatal("attributesToRetrieve: Retrieve unrequested field in response")
	}

	// Test basic search with attributesToCrop

	resp, err = client.Search(indexUID).Search(SearchRequest{
		Query:            "to",
		AttributesToCrop: []string{"title"},
		CropLength:       7,
	})

	if err != nil {
		t.Fatal(err)
	}

	if resp.Hits[0].(map[string]interface{})["title"] == nil {
		fmt.Println(resp)
		t.Fatal("attributesToCrop: Couldn't retrieve field in response")
	}
	formatted := resp.Hits[0].(map[string]interface{})["_formatted"]
	if formatted.(map[string]interface{})["title"] != "Guide to the" {
		fmt.Println(resp)
		t.Fatal("attributesToCrop: CropLength didn't work as expected")
	}

	// Test basic search with filters

	resp, err = client.Search(indexUID).Search(SearchRequest{
		Query:   "and",
		Filters: "tag = \"Nice book\"",
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Hits) != 1 {
		fmt.Println(resp)
		t.Fatal("filters: Unable to filter properly")
	}
	if resp.Hits[0].(map[string]interface{})["title"] != "Pride and Prejudice" {
		fmt.Println(resp)
		t.Fatal("filters: Unable to filter properly")
	}

	// Test basic search with matches

	resp, err = client.Search(indexUID).Search(SearchRequest{
		Query:   "and",
		Matches: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	if resp.Hits[0].(map[string]interface{})["_matchesInfo"] == nil {
		fmt.Println(resp)
		t.Fatal("matches: Mathes info not found")
	}

	// Test basic search with facetsDistribution

	r2, err := client.Settings(indexUID).UpdateAttributesForFaceting([]string{"tag"})

	if err != nil {
		t.Fatal(err)
	}

	client.DefaultWaitForPendingUpdate(indexUID, r2)

	resp, err = client.Search(indexUID).Search(SearchRequest{
		Query:              "prince",
		FacetsDistribution: []string{"*"},
	})

	if err != nil {
		t.Fatal(err)
	}

	tagCount := resp.FacetsDistribution.(map[string]interface{})["tag"]

	if len(tagCount.(map[string]interface{})) != 2 {
		fmt.Println(tagCount.(map[string]interface{}))
		t.Fatal("facetsDistribution: Wrong count of facet options")
	}

	if tagCount.(map[string]interface{})["interesting book"] != float64(2) {
		fmt.Println(reflect.TypeOf(tagCount.(map[string]interface{})["interesting book"]))
		t.Fatal("facetsDistribution: Wrong count on facetDistribution")
	}

	r2, _ = client.Settings(indexUID).ResetAttributesForFaceting()
	client.DefaultWaitForPendingUpdate(indexUID, r2)

	// Test basic search with facetFilters

	r2, err = client.Settings(indexUID).UpdateAttributesForFaceting([]string{"tag", "title"})

	if err != nil {
		t.Fatal(err)
	}

	client.DefaultWaitForPendingUpdate(indexUID, r2)

	resp, err = client.Search(indexUID).Search(SearchRequest{
		Query:        "prince",
		FacetFilters: []string{"tag:interesting book"},
	})
	if err != nil {
		fmt.Println("Error:", err)
	}

	if len(resp.Hits) != 2 {
		fmt.Println(resp)
		t.Fatal("facetsFilters: Error on single attribute facet search")
	}

	resp, err = client.Search(indexUID).Search(SearchRequest{
		Query:        "prince",
		FacetFilters: []string{"tag:interesting book", "tag:nice book"},
	})
	if err != nil {
		fmt.Println("Error:", err)
	}

	if len(resp.Hits) != 0 {
		fmt.Println(resp)
		t.Fatal("facetsFilters: Error on 'AND' in attribute facet search")
	}

	resp, err = client.Search(indexUID).Search(SearchRequest{
		Query:        "prince",
		FacetFilters: [][]string{{"tag:interesting book", "tag:nice book"}},
	})
	if err != nil {
		fmt.Println("Error:", err)
	}

	if len(resp.Hits) != 3 {
		fmt.Println(resp)
		t.Fatal("facetsFilters: Error on 'OR' in attribute facet search")
	}

}
