import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/meilisearch/meilisearch-go"

func main() {
	// Initialize the Meilisearch client with environment configuration
	host := getenv("MEILI_HOST", "http://localhost:7700")
	apiKey := os.Getenv("MEILI_API_KEY")
	client := meilisearch.New(host, meilisearch.WithAPIKey(apiKey))
	defer client.Close()

	// Test connection to Meilisearch
	fmt.Println("Testing connection to Meilisearch...")
	}
	fmt.Printf("Index '%s' is ready!\n", indexUID)

	// Configure filterable and facet attributes
	fmt.Println("Configuring filterable/faceted attributes...")
	settingsTask, err := index.UpdateSettings(&meilisearch.Settings{
		FilterableAttributes: &[]string{"year", "genres"},
	})
	if err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = client.WaitForTaskWithContext(ctx, settingsTask.TaskUID, 100*time.Millisecond)
	}
	if err != nil {
		log.Fatalf("Failed to apply settings: %v", err)
	}

	// Prepare sample movie data
	movies := []Movie{
		{ID: "1", Title: "The Dark Knight", Year: 2008, Rating: 9.0, Genres: []string{"action", "crime", "drama"}},
		{ID: "2", Title: "Inception", Year: 2010, Rating: 8.8, Genres: []string{"action", "sci-fi", "thriller"}},
		{ID: "3", Title: "The Godfather", Year: 1972, Rating: 9.2, Genres: []string{"crime", "drama"}},
		{ID: "4", Title: "Pulp Fiction", Year: 1994, Rating: 8.9, Genres: []string{"crime", "drama"}},
		{ID: "5", Title: "Fight Club", Year: 1999, Rating: 8.8, Genres: []string{"drama"}},
	}

	fmt.Printf("Adding %d movies to the index...\n", len(movies))
	task, err := index.AddDocuments(movies)
	if err != nil {
		log.Fatalf("Failed to add documents: %v", err)
	}

	// Wait for the task to complete
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = client.WaitForTaskWithContext(ctx, task.TaskUID, 100*time.Millisecond)
	if err != nil {
		log.Fatalf("Failed to index documents: %v", err)
	}
	fmt.Printf("✅ Indexed documents successfully! (Task ID: %d)\n", task.TaskUID)

	// Simple search
	fmt.Println("\n1. Simple search for 'action':")
	searchResult, err := index.Search("action", &meilisearch.SearchRequest{
		Limit: 5,
	})
	if err != nil {
		log.Fatalf("Failed to search: %v", err)
	}

	fmt.Printf("Found %d results\n", len(searchResult.Hits))
	for i, hit := range searchResult.Hits {
	fmt.Printf("Found %d drama movies after 1990\n", len(searchResult.Hits))
	for i, hit := range searchResult.Hits {
		fmt.Printf("  %d. %s (%d) - Rating: %.1f\n", i+1, hit["title"], int(hit["year"].(float64)), hit["rating"])
	}

	// Search with filters and facets
	fmt.Println("\n2. Advanced search with filters and facets:")
	searchResult, err = index.Search("drama", &meilisearch.SearchRequest{
		Filter:                "year > 1990",
		Facets:                []string{"genres", "year"},
		AttributesToHighlight: []string{"title", "overview"},
		Limit:                 10,
	})
	if err != nil {
		log.Fatalf("Failed to search with filters: %v", err)
	fmt.Printf("Found %d drama movies after 1990\n", len(searchResult.Hits))
	for i, hit := range searchResult.Hits {
		fmt.Printf("  %d. %s (%d) - Rating: %.1f\n", i+1, hit["title"], int(hit["year"].(float64)), hit["rating"])
	}
	}

	if len(searchResult.FacetDistribution) > 0 {
		fmt.Println("\nFacet distribution:")
		var facets map[string]map[string]float64
		if err := json.Unmarshal(searchResult.FacetDistribution, &facets); err != nil {
			log.Printf("Failed to parse facetDistribution: %v", err)
		} else {
			for facet, distribution := range facets {
				fmt.Printf("  %s: %v\n", facet, distribution)
			}
		}
	}

	fmt.Println("\nSearch example completed successfully!")
}

type Movie struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Year   int      `json:"year"`
	Rating float64  `json:"rating"`
	Genres []string `json:"genres"`
}

// getenv returns the value of the environment variable named by the key,
// or def if the variable is not present or empty.
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
