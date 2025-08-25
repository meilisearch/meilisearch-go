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
		{

	// Search with filters and facets
	searchResult, err = index.Search("drama", &meilisearch.SearchRequest{
		Filter:                "year > 1990",
		Facets:                []string{"genres", "year"},
		AttributesToHighlight: []string{"title", "overview"},
		Limit:                 10,
	})
	if err != nil {
		log.Fatalf("Failed to search with filters: %v", err)
		fmt.Printf("  %d. %s (%d) - Rating: %.1f\n", i+1, hit["title"], int(hit["year"].(float64)), hit["rating"])
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

// getenv returns the value of the environment variable named by the key,
// or def if the variable is not present or empty.
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
