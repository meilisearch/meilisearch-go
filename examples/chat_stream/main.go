package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/meilisearch/meilisearch-go"
)

// Document represents a knowledge base document for chat context
type Document struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Topic   string `json:"topic"`
}

func main() {
	// Initialize the Meilisearch client
	// Replace with your Meilisearch server URL and API key
	client := meilisearch.New("http://localhost:7700", meilisearch.WithAPIKey("your-api-key"))

	// Test connection
	fmt.Println("Testing connection to Meilisearch...")
	health, err := client.Health()
	if err != nil {
		log.Fatalf("Failed to connect to Meilisearch: %v", err)
	}
	fmt.Printf("Meilisearch is %s\n", health.Status)

	// Setup knowledge base for chat context
	if err := setupKnowledgeBase(client); err != nil {
		log.Fatalf("Failed to setup knowledge base: %v", err)
	}

	// List available chat workspaces
	fmt.Println("\nListing chat workspaces...")
	workspaces, err := client.ListChatWorkspaces(&meilisearch.ListChatWorkSpaceQuery{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		log.Printf("Warning: Could not list chat workspaces: %v", err)
	} else {
		fmt.Printf("Found %d chat workspaces\n", len(workspaces.Results))
		for _, workspace := range workspaces.Results {
			fmt.Printf("  - Workspace: %s\n", workspace.UID)
		}
	}

	// Start interactive chat session
	fmt.Println("\n--- Interactive Chat Session ---")
	fmt.Println("Type your questions (or 'quit' to exit):")
	
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("\nYou: ")
		if !scanner.Scan() {
			break
		}
		
		userInput := strings.TrimSpace(scanner.Text())
		if userInput == "" {
			continue
		}
		if strings.ToLower(userInput) == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		// Demonstrate chat streaming
		if err := performChatStream(client, userInput); err != nil {
			log.Printf("Chat stream error: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error: %v", err)
	}
}

func setupKnowledgeBase(client meilisearch.ServiceManager) error {
	fmt.Println("Setting up knowledge base...")
	
	// Create knowledge base index
	indexUID := "knowledge_base"
	index := client.Index(indexUID)

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

	// Add sample knowledge base documents
	documents := []Document{
		{
			ID:      1,
			Title:   "Getting Started with Meilisearch",
			Content: "Meilisearch is a powerful, fast, open-source, easy to use and deploy search engine. It provides instant search capabilities with typo tolerance and filtering.",
			Topic:   "basics",
		},
		{
			ID:      2,
			Title:   "Search Features",
			Content: "Meilisearch offers faceted search, geo-search, full-text search with typo tolerance, synonyms, and custom ranking rules.",
			Topic:   "features",
		},
		{
			ID:      3,
			Title:   "API Integration",
			Content: "The Meilisearch API is RESTful and supports multiple SDKs including Go, JavaScript, Python, and more for easy integration.",
			Topic:   "integration",
		},
	}

	addTask, err := index.AddDocuments(documents)
	if err != nil {
		return fmt.Errorf("failed to add documents: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	_, err = client.WaitForTaskWithContext(ctx, addTask.TaskUID, 100*time.Millisecond)
	if err != nil {
		return fmt.Errorf("failed to wait for document addition: %w", err)
	}

	fmt.Println("Knowledge base setup completed!")
	return nil
}

func performChatStream(client meilisearch.ServiceManager, query string) error {
	// Create a chat completion query with streaming enabled
	chatQuery := &meilisearch.ChatCompletionQuery{
		Model: "gpt-3.5-turbo",
		Messages: []meilisearch.ChatMessage{
			{
				Role:    "system",
				Content: "You are a helpful assistant that answers questions about Meilisearch. Use the provided knowledge base to give accurate answers.",
			},
			{
				Role:    "user",
				Content: query,
			},
		},
		Stream:      true,
		MaxTokens:   200,
		Temperature: 0.7,
	}

	// Note: This example assumes you have a chat workspace configured
	// In a real implementation, you would need to have Meilisearch configured
	// with chat capabilities and appropriate workspace setup
	workspaceID := "default" // Replace with your actual workspace ID

	fmt.Print("Assistant: ")
	
	// Start streaming chat completion
	stream, err := client.ChatCompletionStream(workspaceID, chatQuery)
	if err != nil {
		return fmt.Errorf("failed to start chat stream: %w", err)
	}
	defer stream.Close()

	// Process streaming responses
	for {
		chunk, err := stream.Next()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("stream error: %w", err)
		}

		if chunk != nil && len(chunk.Choices) > 0 {
			delta := chunk.Choices[0].Delta
			if delta.Content != "" {
				fmt.Print(delta.Content)
			}
		}
	}
	
	fmt.Println() // New line after streaming response
	return nil
}
