package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/meilisearch/meilisearch-go"
)

// User represents a user document
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Active   bool   `json:"active"`
	JoinDate string `json:"join_date"`
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

	// Create index for document operations
	indexUID := "users"
	if err := createIndex(client, indexUID); err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}

	index := client.Index(indexUID)

	fmt.Println("\n📄 Document Management Examples")
	fmt.Println("===============================")

	// 1. Add documents (batch)
	fmt.Println("1. Adding documents in batch:")
	users := []User{
		{ID: 1, Name: "Alice Johnson", Email: "alice@example.com", Role: "admin", Active: true, JoinDate: "2023-01-15"},
		{ID: 2, Name: "Bob Smith", Email: "bob@example.com", Role: "user", Active: true, JoinDate: "2023-02-20"},
		{ID: 3, Name: "Charlie Brown", Email: "charlie@example.com", Role: "moderator", Active: false, JoinDate: "2023-03-10"},
	}

	task, err := index.AddDocuments(users, nil)
	if err != nil {
		log.Fatalf("Failed to add documents: %v", err)
	}
	
	if err := waitForTask(client, task.TaskUID); err != nil {
		log.Fatalf("Failed to wait for add task: %v", err)
	}
	fmt.Printf("✅ Added %d documents (Task ID: %d)\n", len(users), task.TaskUID)

	// 2. Add single document
	fmt.Println("\n2. Adding single document:")
	newUser := User{
		ID:       4,
		Name:     "Diana Prince",
		Email:    "diana@example.com",
		Role:     "user",
		Active:   true,
		JoinDate: "2023-04-05",
	}

	task, err = index.AddDocuments([]User{newUser}, nil)
	if err != nil {
		log.Fatalf("Failed to add single document: %v", err)
	}
	
	if err := waitForTask(client, task.TaskUID); err != nil {
		log.Fatalf("Failed to wait for add task: %v", err)
	}
	fmt.Printf("✅ Added single document (Task ID: %d)\n", task.TaskUID)

	// 3. Update documents
	fmt.Println("\n3. Updating documents:")
	updatedUsers := []User{
		{ID: 2, Name: "Bob Smith Jr.", Email: "bob.jr@example.com", Role: "admin", Active: true, JoinDate: "2023-02-20"},
		{ID: 3, Name: "Charlie Brown", Email: "charlie@example.com", Role: "moderator", Active: true, JoinDate: "2023-03-10"}, // Activating user
	}

	task, err = index.UpdateDocuments(updatedUsers, nil)
	if err != nil {
		log.Fatalf("Failed to update documents: %v", err)
	}
	
	if err := waitForTask(client, task.TaskUID); err != nil {
		log.Fatalf("Failed to wait for update task: %v", err)
	}
	fmt.Printf("✅ Updated %d documents (Task ID: %d)\n", len(updatedUsers), task.TaskUID)

	// 4. Get document by ID
	fmt.Println("\n4. Getting document by ID:")
	var doc User
	err = index.GetDocument("2", nil, &doc)
	if err != nil {
		log.Fatalf("Failed to get document: %v", err)
	}
	fmt.Printf("Retrieved document: %v\n", doc)

	// 5. Get multiple documents
	fmt.Println("\n5. Getting multiple documents:")
	var docs meilisearch.DocumentsResult
	err = index.GetDocuments(&meilisearch.DocumentsQuery{
		Limit:  10,
		Offset: 0,
	}, &docs)
	if err != nil {
		log.Fatalf("Failed to get documents: %v", err)
	}
	fmt.Printf("Retrieved %d documents:\n", len(docs.Results))
	for i, doc := range docs.Results {
		fmt.Printf("  %d. %v\n", i+1, doc)
	}

	// 6. Delete documents
	fmt.Println("\n6. Deleting documents:")
	task, err = index.DeleteDocument("4")
	if err != nil {
		log.Fatalf("Failed to delete document: %v", err)
	}
	
	if err := waitForTask(client, task.TaskUID); err != nil {
		log.Fatalf("Failed to wait for delete task: %v", err)
	}
	fmt.Printf("✅ Deleted document with ID 4 (Task ID: %d)\n", task.TaskUID)

	// 7. Delete multiple documents
	fmt.Println("\n7. Deleting multiple documents:")
	task, err = index.DeleteDocuments([]string{"1", "3"})
	if err != nil {
		log.Fatalf("Failed to delete documents: %v", err)
	}
	
	if err := waitForTask(client, task.TaskUID); err != nil {
		log.Fatalf("Failed to wait for delete task: %v", err)
	}
	fmt.Printf("✅ Deleted multiple documents (Task ID: %d)\n", task.TaskUID)

	// 8. Final document count
	fmt.Println("\n8. Final document count:")
	var finalDocs meilisearch.DocumentsResult
	err = index.GetDocuments(&meilisearch.DocumentsQuery{Limit: 100}, &finalDocs)
	if err != nil {
		log.Fatalf("Failed to get final documents: %v", err)
	}
	fmt.Printf("Remaining documents: %d\n", len(finalDocs.Results))

	fmt.Println("\nDocument management examples completed successfully! 🎉")
}

func createIndex(client meilisearch.ServiceManager, indexUID string) error {
	fmt.Printf("Creating index '%s'...\n", indexUID)
	
	task, err := client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        indexUID,
		PrimaryKey: "id",
	})
	if err != nil {
		log.Printf("Index might already exist: %v", err)
		return nil // Continue if index exists
	}
	
	return waitForTask(client, task.TaskUID)
}

func waitForTask(client meilisearch.ServiceManager, taskUID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	_, err := client.WaitForTaskWithContext(ctx, taskUID, 100*time.Millisecond)
	return err
}

// getenv returns the value of the environment variable named by the key,
// or def if the variable is not present or empty.
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
