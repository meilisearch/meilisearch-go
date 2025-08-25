package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/meilisearch/meilisearch-go"
)

// Movie represents a document structure for our search index
type Movie struct {
	ID       int      `json:"id"`
	Title    string   `json:"title"`
	Overview string   `json:"overview"`
	Genres   []string `json:"genres"`
	Rating   float64  `json:"rating"`
	Year     int      `json:"year"`
}

func main() {
	// Initialize the Meilisearch client
	// Replace with your Meilisearch server URL and API key
	client := meilisearch.New("http://localhost:7700", meilisearch.WithAPIKey("your-api-key"))

	// Test connection to Meilisearch
	fmt.Println("Testing connection to Meilisearch...")
	health, err := client.Health()
	if err != nil {
		log.Fatalf("Failed to connect to Meilisearch: %v", err)
	}
	fmt.Printf("Meilisearch is %s\n", health.Status)

	// Create or get an index
	indexUID := "movies"
	index := client.Index(indexUID)

	// Create the index with a primary key
	fmt.Printf("Creating index '%s'...\n", indexUID)
	task, err := client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        indexUID,
		PrimaryKey: "id",
	})
	if err != nil {
		log.Printf("Index might already exist: %v", err)
	} else {
		// Wait for the index creation task to complete
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		_, err = client.WaitForTaskWithContext(ctx, task.TaskUID, 100*time.Millisecond)
		if err != nil {
			log.Printf("Warning: Index creation task didn't complete: %v", err)
		} else {
			fmt.Println("Index created successfully!")
		}
	}

	// Prepare sample movie data
	movies := []Movie{
		{
			ID:       1,
			Title:    "The Shawshank Redemption",
			Overview: "Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.",
			Genres:   []string{"Drama"},
			Rating:   9.3,
			Year:     1994,
		},
		{
			ID:       2,
			Title:    "The Godfather",
			Overview: "The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.",
			Genres:   []string{"Crime", "Drama"},
			Rating:   9.2,
			Year:     1972,
		},
		{
			ID:       3,
			Title:    "Pulp Fiction",
			Overview: "The lives of two mob hitmen, a boxer, a gangster and his wife intertwine in four tales of violence and redemption.",
			Genres:   []string{"Crime", "Drama"},
			Rating:   8.9,
			Year:     1994,
		},
		{
			ID:       4,
			Title:    "The Dark Knight",
			Overview: "When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests.",
			Genres:   []string{"Action", "Crime", "Drama"},
			Rating:   9.0,
			Year:     2008,
		},
	}

	// Add documents to the index
	fmt.Println("Adding documents to the index...")
	addTask, err := index.AddDocuments(movies)
	if err != nil {
		log.Fatalf("Failed to add documents: %v", err)
	}

	// Wait for the document addition task to complete
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	finalTask, err := client.WaitForTaskWithContext(ctx, addTask.TaskUID, 100*time.Millisecond)
	if err != nil {
		log.Fatalf("Failed to wait for document addition: %v", err)
	}
	fmt.Printf("Documents added successfully! Task status: %s\n", finalTask.Status)

	// Perform various searches
	fmt.Println("\n--- Search Examples ---")

	// Basic search
	searchResult, err := index.Search("godfather", &meilisearch.SearchRequest{})
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}
	fmt.Printf("Search for 'godfather' found %d results:\n", searchResult.EstimatedTotalHits)
	for _, hit := range searchResult.Hits {
		fmt.Printf("  - %v\n", hit)
	}

	// Search with filters and facets
	searchResult, err = index.Search("drama", &meilisearch.SearchRequest{
		Filter:              []string{"year > 1990"},
		FacetsDistribution:  []string{"genres", "year"},
		AttributesToHighlight: []string{"title", "overview"},
		Limit:               10,
	})
	if err != nil {
		log.Fatalf("Advanced search failed: %v", err)
	}
	fmt.Printf("\nSearch for 'drama' with year > 1990 found %d results:\n", searchResult.EstimatedTotalHits)
	for _, hit := range searchResult.Hits {
		fmt.Printf("  - %v\n", hit)
	}
	
	if searchResult.FacetsDistribution != nil {
		fmt.Println("\nFacets distribution:")
		for facet, distribution := range *searchResult.FacetsDistribution {
			fmt.Printf("  %s: %v\n", facet, distribution)
		}
	}

	fmt.Println("\nSearch example completed successfully!")
}
