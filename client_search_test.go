package meilisearch

import (
	"fmt"
	"testing"
)

func TestClientSearch_Search(t *testing.T) {
	var indexUID = "TestClientSearch_Search"

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	booksTest := []docTestBooks{
		{BookID: 123, Title: "Pride and Prejudice", Tag: "Nice book"},
		{BookID: 456, Title: "Le Petit Prince", Tag: "Nice book"},
		{BookID: 1, Title: "Alice In Wonderland", Tag: "Nice book"},
		{BookID: 1344, Title: "The Hobbit", Tag: "Nice book"},
		{BookID: 4, Title: "Harry Potter and the Half-Blood Prince", Tag: "Interesting book"},
		{BookID: 42, Title: "The Hitchhiker's Guide to the Galaxy", Tag: "Interesting book"},
		{BookID: 24, Title: "You are a princess", Tag: "Interesting book"},
	}

	updateIDRes, err := client.
		Documents(indexUID).
		AddOrUpdate(booksTest)

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
		t.Fatal("Basic search: number of hits should be equal to 3")
	}
	title := resp.Hits[0].(map[string]interface{})["title"]
	if title != booksTest[1].Title {
		fmt.Println(resp)
		t.Fatalf("Basic search: should have found %s\n", booksTest[1].Title)
	}
	if resp.NbHits != 3 {
		fmt.Println(resp)
		t.Fatalf("Basic search: wrong number of hits, should have 3, got %d\n", resp.NbHits)
	}

	// Test basic empty search

	resp, err = client.Search(indexUID).Search(SearchRequest{
		Query: "",
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Hits) != len(booksTest) {
		fmt.Println(resp)
		t.Fatal("Basic placeholder search with an empty string: should return placeholder results")
	}

	// Test basic placeholder search

	resp, err = client.Search(indexUID).Search(SearchRequest{
		PlaceholderSearch: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Hits) != len(booksTest) {
		fmt.Println(resp)
		t.Fatal("Basic placeholder search with no Query: should return placeholder results")
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
		t.Fatal("Search offset: number of hits should be equal to 1")
	}
	title = resp.Hits[0].(map[string]interface{})["title"]
	if title != booksTest[1].Title {
		fmt.Println(resp)
		t.Fatalf("Basic search: should have found %s\n", booksTest[1].Title)
	}

	// Test basic placeholder search with limit

	resp, err = client.Search(indexUID).Search(SearchRequest{
		PlaceholderSearch: true,
		Limit:             3,
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Hits) != 3 {
		fmt.Println(resp)
		t.Fatal("Basic placeholder search with limit: should return 3 results")
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
	retrievedTitles := []string{
		fmt.Sprint(resp.Hits[0].(map[string]interface{})["title"]),
		fmt.Sprint(resp.Hits[1].(map[string]interface{})["title"]),
	}
	expectedTitles := []string{
		booksTest[4].Title,
		booksTest[6].Title,
	}

	for title := range expectedTitles {
		found := false
		for retrievedTitle := range retrievedTitles {
			if title == retrievedTitle {
				found = true
				break
			}
		}
		if !found {
			fmt.Println(resp)
			t.Fatal("Search offset: should have found 'Harry Potter and the Half-Blood Prince'")
		}
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

	if tagCount.(map[string]interface{})["Interesting book"] != float64(2) {
		fmt.Println(tagCount.(map[string]interface{})["Interesting book"])
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
