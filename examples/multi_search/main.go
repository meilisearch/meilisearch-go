package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/meilisearch/meilisearch-go"
)

// Product represents a product document
type Product struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Price       float64  `json:"price"`
	Brand       string   `json:"brand"`
	Tags        []string `json:"tags"`
	InStock     bool     `json:"in_stock"`
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

	// Setup the products index
	if err := setupProductsIndex(client); err != nil {
		log.Fatalf("Failed to setup products index: %v", err)
	}

	// Demonstrate multi-search capabilities
	fmt.Println("\n🔍 Multi-Search Examples")
	fmt.Println("========================")

	// Perform multi-search with different queries
	multiSearchRequest := &meilisearch.MultiSearchRequest{
		Queries: []*meilisearch.SearchRequest{
			{
				IndexUID: "products",
				Query:    "laptop",
				Limit:    5,
				Filter:   []string{"category = electronics"},
			},
			{
				IndexUID: "products", 
				Query:    "coffee",
				Limit:    3,
				Filter:   []string{"in_stock = true"},
			},
			{
				IndexUID: "products",
				Query:    "",
				Limit:    10,
				Filter:   []string{"price < 100"},
				Sort:     []string{"price:asc"},
			},
		},
	}

	// Execute multi-search
	results, err := client.MultiSearch(multiSearchRequest)
	if err != nil {
		log.Fatalf("Multi-search failed: %v", err)
	}

	// Display results for each query
	fmt.Printf("📊 Executed %d search queries simultaneously:\n\n", len(results.Results))
	
	for i, result := range results.Results {
		query := multiSearchRequest.Queries[i]
		fmt.Printf("Query %d: '%s' in %s\n", i+1, query.Query, query.IndexUID)
		fmt.Printf("Filter: %v\n", query.Filter)
		fmt.Printf("Found %d results:\n", result.EstimatedTotalHits)
		
		for j, hit := range result.Hits {
			fmt.Printf("  %d. %v\n", j+1, hit)
		}
		fmt.Println()
	}

	fmt.Println("Multi-search example completed successfully! 🎉")
}

func setupProductsIndex(client meilisearch.ServiceManager) error {
	fmt.Println("Setting up products index...")
	
	indexUID := "products"
	
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

	// Add sample products
	products := []Product{
		{ID: 1, Name: "Gaming Laptop", Description: "High-performance laptop for gaming", Category: "electronics", Price: 1299.99, Brand: "TechBrand", Tags: []string{"gaming", "laptop", "computer"}, InStock: true},
		{ID: 2, Name: "Coffee Maker", Description: "Automatic drip coffee maker", Category: "appliances", Price: 89.99, Brand: "BrewMaster", Tags: []string{"coffee", "kitchen", "appliance"}, InStock: true},
		{ID: 3, Name: "Wireless Mouse", Description: "Ergonomic wireless mouse", Category: "electronics", Price: 29.99, Brand: "TechBrand", Tags: []string{"mouse", "wireless", "computer"}, InStock: true},
		{ID: 4, Name: "Coffee Beans", Description: "Premium arabica coffee beans", Category: "food", Price: 24.99, Brand: "CoffeeCorp", Tags: []string{"coffee", "beans", "premium"}, InStock: false},
		{ID: 5, Name: "Bluetooth Headphones", Description: "Noise-canceling headphones", Category: "electronics", Price: 199.99, Brand: "AudioTech", Tags: []string{"headphones", "bluetooth", "audio"}, InStock: true},
	}

	index := client.Index(indexUID)
	addTask, err := index.AddDocuments(products)
	if err != nil {
		return fmt.Errorf("failed to add documents: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	_, err = client.WaitForTaskWithContext(ctx, addTask.TaskUID, 100*time.Millisecond)
	if err != nil {
		return fmt.Errorf("failed to wait for document addition: %w", err)
	}

	fmt.Println("✅ Products index setup completed!")
	return nil
}

// getenv returns the value of the environment variable named by the key,
// or def if the variable is not present or empty.
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
