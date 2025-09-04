package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/meilisearch/meilisearch-go"
)

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

	fmt.Println("\n🏗️  Index Creation and Settings Examples")
	fmt.Println("========================================")

	// 1. Create basic index
	fmt.Println("1. Creating basic index:")
	indexUID := "articles"
	task, err := client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        indexUID,
		PrimaryKey: "id",
	})
	if err != nil {
		log.Printf("Index might already exist: %v", err)
	} else {
		if err := waitForTask(client, task.TaskUID); err != nil {
			log.Fatalf("Failed to create index: %v", err)
		}
		fmt.Printf("✅ Created index '%s' (Task ID: %d)\n", indexUID, task.TaskUID)
	}

	index := client.Index(indexUID)

	// 2. Configure comprehensive settings
	fmt.Println("\n2. Configuring index settings:")
	settings := &meilisearch.Settings{
		// Searchable attributes - fields that are searched
		SearchableAttributes: []string{"title", "content", "author", "tags"},
		
		// Displayed attributes - fields returned in search results
		DisplayedAttributes: []string{"id", "title", "author", "publish_date", "category", "summary"},
		
		// Filterable attributes - fields that can be used in filters
		FilterableAttributes: []string{"category", "author", "publish_date", "status", "featured"},
		
		// Sortable attributes - fields that can be used for sorting
		SortableAttributes: []string{"publish_date", "title", "author"},
		
		// Ranking rules - control result relevance
		RankingRules: []string{
			"words",
			"typo", 
			"proximity",
			"attribute",
			"sort",
			"exactness",
		},
		
		// Stop words - words ignored during search
		StopWords: []string{"the", "a", "an", "and", "or", "but", "in", "on", "at", "to", "for", "of", "with", "by"},
		
		// Synonyms - alternative words
		Synonyms: map[string][]string{
			"programming": {"coding", "development"},
			"javascript":  {"js", "ecmascript"},
			"golang":      {"go"},
		},
		
		// Distinct attribute - deduplicate results
		DistinctAttribute: stringPtr("title"),
		
		// Typo tolerance settings
		TypoTolerance: &meilisearch.TypoTolerance{
			Enabled: true,
			MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
				OneTypo:  5,
				TwoTypos: 9,
			},
		},
		
		// Pagination settings
		Pagination: &meilisearch.Pagination{
			MaxTotalHits: 1000,
		},
	}

	settingsTask, err := index.UpdateSettings(settings)
	if err != nil {
		log.Fatalf("Failed to update settings: %v", err)
	}
	
	if err := waitForTask(client, settingsTask.TaskUID); err != nil {
		log.Fatalf("Failed to wait for settings update: %v", err)
	}
	fmt.Printf("✅ Updated index settings (Task ID: %d)\n", settingsTask.TaskUID)

	// 3. Retrieve and display current settings
	fmt.Println("\n3. Current index settings:")
	currentSettings, err := index.GetSettings()
	if err != nil {
		log.Fatalf("Failed to get settings: %v", err)
	}
	
	fmt.Printf("Searchable attributes: %v\n", currentSettings.SearchableAttributes)
	fmt.Printf("Filterable attributes: %v\n", currentSettings.FilterableAttributes)
	fmt.Printf("Sortable attributes: %v\n", currentSettings.SortableAttributes)

	fmt.Println("\nIndex creation and settings examples completed successfully! 🎉")
}

// waitForTask waits for a task to complete
func waitForTask(client meilisearch.ServiceManager, taskUID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := client.WaitForTaskWithContext(ctx, taskUID, 100*time.Millisecond)
	return err
}

// stringPtr returns a pointer to a string
func stringPtr(s string) *string {
	return &s
}

// getenv returns the value of the environment variable named by the key,
// or def if the variable is not present or empty.
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}