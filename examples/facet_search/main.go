package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/meilisearch/meilisearch-go"
)

// Book represents a book document with facetable attributes
type Book struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Genre       string   `json:"genre"`
	Language    string   `json:"language"`
	PublishYear int      `json:"publish_year"`
	Rating      float64  `json:"rating"`
	Pages       int      `json:"pages"`
	Publisher   string   `json:"publisher"`
	Tags        []string `json:"tags"`
	InPrint     bool     `json:"in_print"`
}

func main() {
	// Initialize the Meilisearch client
	host := getenv("MEILI_HOST", "http://localhost:7700")
	apiKey := os.Getenv("MEILI_API_KEY")
	client := meilisearch.New(host, meilisearch.WithAPIKey(apiKey))
	defer client.Close()

	// Test connection
	fmt.Println("Testing connection to Meilisearch...")
	if !client.IsHealthy() {
		log.Fatal("Meilisearch is not available")
	}
	fmt.Println("✅ Connected to Meilisearch")

	// Setup the books index with facetable attributes
	if err := setupBooksIndex(client); err != nil {
		log.Fatalf("Failed to setup books index: %v", err)
	}

	// Demonstrate facet search capabilities
	fmt.Println("\n🔍 Faceted Search Examples")
	fmt.Println("==========================")

	// Basic faceted search
	fmt.Println("1. Basic faceted search with distribution:")
	searchResult, err := client.Index("books").Search("fiction", &meilisearch.SearchRequest{
		Facets: []string{"genre", "language", "publish_year"},
		Limit:  5,
	})
	if err != nil {
		log.Fatalf("Faceted search failed: %v", err)
	}

	displaySearchResults("fiction", searchResult)

	// Advanced faceted search with filters
	fmt.Println("\n2. Faceted search with filters:")
	searchResult, err = client.Index("books").Search("", &meilisearch.SearchRequest{
		Filter: "genre = fantasy AND publish_year > 2000",
		Facets: []string{"language", "rating", "publisher"},
		Sort:   []string{"rating:desc"},
		Limit:  10,
	})
	if err != nil {
		log.Fatalf("Advanced faceted search failed: %v", err)
	}

	displaySearchResults("fantasy books after 2000", searchResult)

	// Facet search with specific facet query
	fmt.Println("\n3. Facet-specific search:")
	facetRaw, err := client.Index("books").FacetSearch(&meilisearch.FacetSearchRequest{
		FacetName:  "genre",
		FacetQuery: "sci",
		Q:          "space",
	})
	if err != nil {
		log.Fatalf("Facet search failed: %v", err)
	}

	// Unmarshal the raw JSON response into FacetSearchResponse
	var facetResult meilisearch.FacetSearchResponse
	err = json.Unmarshal(*facetRaw, &facetResult)
	if err != nil {
		log.Fatalf("Failed to unmarshal facet search response: %v", err)
	}

	fmt.Printf("Facet search for 'sci' in genre facet with query 'space':\n")

	results := make([]Book, 0)
	if err := facetResult.FacetHits.DecodeInto(&results); err != nil {
		log.Fatalf("Facet search failed: %v", err)
	}

	for _, result := range results {
		fmt.Printf(" - %s by %s (%d)\n", result.Title, result.Author, result.PublishYear)
	}

	fmt.Println("\nFaceted search examples completed successfully! 🎉")
}

func setupBooksIndex(client meilisearch.ServiceManager) error {
	fmt.Println("Setting up books index with facetable attributes...")

	indexUID := "books"

	// Create index
	task, err := client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        indexUID,
		PrimaryKey: "id",
	})
	if err != nil {
		log.Printf("Index might already exist: %v", err)
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = client.WaitForTaskWithContext(ctx, task.TaskUID, 100*time.Millisecond)
		if err != nil {
			return fmt.Errorf("index creation failed: %w", err)
		}
	}

	// Configure facetable attributes
	index := client.Index(indexUID)
	settings := &meilisearch.Settings{
		FilterableAttributes: []string{"genre", "language", "publish_year", "rating", "publisher", "in_print"},
		SortableAttributes:   []string{"rating", "publish_year", "pages"},
	}

	settingsTask, err := index.UpdateSettings(settings)
	if err != nil {
		return fmt.Errorf("failed to update settings: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = client.WaitForTaskWithContext(ctx, settingsTask.TaskUID, 100*time.Millisecond)
	if err != nil {
		return fmt.Errorf("failed to wait for settings update: %w", err)
	}

	// Add sample books
	books := []Book{
		{ID: 1, Title: "Dune", Author: "Frank Herbert", Genre: "science-fiction", Language: "English", PublishYear: 1965, Rating: 4.5, Pages: 688, Publisher: "Ace Books", Tags: []string{"space", "politics"}, InPrint: true},
		{ID: 2, Title: "The Hobbit", Author: "J.R.R. Tolkien", Genre: "fantasy", Language: "English", PublishYear: 1937, Rating: 4.7, Pages: 310, Publisher: "George Allen & Unwin", Tags: []string{"adventure", "magic"}, InPrint: true},
		{ID: 3, Title: "1984", Author: "George Orwell", Genre: "dystopian", Language: "English", PublishYear: 1949, Rating: 4.6, Pages: 328, Publisher: "Secker & Warburg", Tags: []string{"politics", "surveillance"}, InPrint: true},
		{ID: 4, Title: "Foundation", Author: "Isaac Asimov", Genre: "science-fiction", Language: "English", PublishYear: 1951, Rating: 4.3, Pages: 244, Publisher: "Gnome Press", Tags: []string{"space", "mathematics"}, InPrint: true},
		{ID: 5, Title: "Harry Potter", Author: "J.K. Rowling", Genre: "fantasy", Language: "English", PublishYear: 1997, Rating: 4.8, Pages: 309, Publisher: "Bloomsbury", Tags: []string{"magic", "school"}, InPrint: true},
	}

	addTask, err := index.AddDocuments(books, nil)
	if err != nil {
		return fmt.Errorf("failed to add documents: %w", err)
	}

	_, err = client.WaitForTaskWithContext(ctx, addTask.TaskUID, 100*time.Millisecond)
	if err != nil {
		return fmt.Errorf("failed to wait for document addition: %w", err)
	}

	fmt.Println("✅ Books index with facetable attributes setup completed!")
	return nil
}

func displaySearchResults(query string, result *meilisearch.SearchResponse) {
	fmt.Printf("Search: '%s' - Found %d results\n", query, result.EstimatedTotalHits)
	if len(result.FacetDistribution) > 0 {
		fmt.Println("Facet distribution:")
		var fd map[string]map[string]int
		if err := json.Unmarshal(result.FacetDistribution, &fd); err == nil {
			for facet, distribution := range fd {
				fmt.Printf("  %s: %v\n", facet, distribution)
			}
		} else {
			fmt.Printf("  (failed to parse facet distribution: %v)\n", err)
		}
	}
}

// getenv returns the value of the environment variable named by the key,
// or def if the variable is not present or empty.
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
