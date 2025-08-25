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
		log.Fatalf("Failed to setup knowledge base: %v", err)
	}

	// List available chat workspaces and pick the first one
	var workspaceID string = "default"
	fmt.Println("\nListing chat workspaces...")
	workspaces, err := client.ListChatWorkspaces(&meilisearch.ListChatWorkSpaceQuery{
		Limit:  10,
		fmt.Printf("Found %d chat workspaces\n", len(workspaces.Results))
		for _, workspace := range workspaces.Results {
			fmt.Printf("  - Workspace: %s\n", workspace.UID)
			if workspaceID == "default" && workspace.UID != "" {
				workspaceID = workspace.UID // Use first available workspace
			}
		}
	}

		}

		// Demonstrate chat streaming
		if err := performChatStream(client, workspaceID, userInput); err != nil {
			log.Printf("Chat stream error: %v", err)
		}
	}
	return nil
}

func performChatStream(client meilisearch.ServiceManager, workspaceID string, query string) error {
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

	// workspaceID provided by caller (picked from the listed workspaces)

	fmt.Print("Assistant: ")
	
	for {
		chunk, err := stream.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("stream error: %w", err)
	fmt.Println() // New line after streaming response
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
