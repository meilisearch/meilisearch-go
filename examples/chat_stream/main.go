package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/meilisearch/meilisearch-go"
)

func main() {
	// Initialize the Meilisearch client with environment configuration
	host := getenv("MEILI_HOST", "http://localhost:7700")
	apiKey := os.Getenv("MEILI_API_KEY")
	client := meilisearch.New(host, meilisearch.WithAPIKey(apiKey))
	defer client.Close()

	// Test connection
	fmt.Println("Testing connection to Meilisearch...")
	health, err := client.Health()
	if err != nil {
		log.Fatalf("Failed to connect to Meilisearch: %v", err)
	}
	fmt.Printf("✅ Connected to Meilisearch (status: %s)\n", health.Status)

	// Basic setup - ensure we have some knowledge base
	if err := setupKnowledgeBase(client); err != nil {
		log.Fatalf("Failed to setup knowledge base: %v", err)
	}

	// List available chat workspaces and pick the first one
	workspaceID := "default"
	fmt.Println("\nListing chat workspaces...")
	workspaces, err := client.ListChatWorkspaces(&meilisearch.ListChatWorkSpaceQuery{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		log.Printf("Warning: could not list chat workspaces: %v", err)
	} else {
		fmt.Printf("Found %d chat workspaces\n", len(workspaces.Results))
		for _, ws := range workspaces.Results {
			fmt.Printf("  - Workspace: %s\n", ws.UID)
		}
		if len(workspaces.Results) > 0 && workspaces.Results[0].UID != "" {
			workspaceID = workspaces.Results[0].UID
		}
	}

	fmt.Println("\n--- Interactive Chat Session ---")
	fmt.Println("Type your questions (or 'quit' to exit):")
	
	// Simple interactive loop
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nYou: ")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)
		if userInput == "" || strings.EqualFold(userInput, "exit") || strings.EqualFold(userInput, "quit") {
			fmt.Println("Bye!")
			break
		}
		if err := performChatStream(client, workspaceID, userInput); err != nil {
			log.Printf("Chat stream error: %v", err)
		}
	}
}

func performChatStream(client meilisearch.ServiceManager, workspaceID string, query string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	
	// Create a chat completion query with streaming enabled
	chatQuery := &meilisearch.ChatCompletionQuery{
		Model: "gpt-3.5-turbo",
		Messages: []*meilisearch.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "You are a helpful assistant that answers questions about Meilisearch. Use the provided knowledge base to give accurate answers.",
			},
			{
				Role:    "user",
				Content: query,
			},
		},
		Stream: true,
	}

	// Start streaming chat completion
	stream, err := client.ChatCompletionStreamWithContext(ctx, workspaceID, chatQuery)
	if err != nil {
		return fmt.Errorf("failed to start chat stream: %w", err)
	}
	defer stream.Close()
	
	fmt.Print("Assistant: ")
	
	for {
		hasNext := stream.Next()
		if !hasNext {
			// Check for any error
			if err := stream.Err(); err != nil {
				if !errors.Is(err, io.EOF) {
					return fmt.Errorf("stream error: %w", err)
				}
			}
			break
		}
		
		chunk := stream.Current()
		
		// Print the streaming content from choices
		if chunk != nil && len(chunk.Choices) > 0 {
			if chunk.Choices[0].Delta.Content != nil && *chunk.Choices[0].Delta.Content != "" {
				fmt.Print(*chunk.Choices[0].Delta.Content)
			}
		}
	}
	fmt.Println() // New line after streaming response
	return nil
}

func setupKnowledgeBase(client meilisearch.ServiceManager) error {
	// Basic setup - this is a placeholder for actual knowledge base setup
	// In a real scenario, you would populate indices with relevant documents
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
