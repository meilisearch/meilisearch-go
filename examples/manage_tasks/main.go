package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/meilisearch/meilisearch-go"
)

// Document represents a simple document structure
type Document struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
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

	fmt.Println("\n📋 Task Management Examples")
	fmt.Println("===========================")

	// 1. Create an index to generate tasks
	indexUID := "task_demo"
	fmt.Println("1. Creating index to generate tasks...")
	
	createTask, err := client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        indexUID,
		PrimaryKey: "id",
	})
	if err != nil {
		log.Printf("Index might already exist: %v", err)
	} else {
		fmt.Printf("✅ Index creation task created (Task ID: %d)\n", createTask.TaskUID)
	}

	// 2. Add documents to generate more tasks
	fmt.Println("\n2. Adding documents to generate more tasks...")
	index := client.Index(indexUID)
	
	documents := []Document{
		{ID: 1, Title: "First Document", Content: "This is the first document content"},
		{ID: 2, Title: "Second Document", Content: "This is the second document content"},
		{ID: 3, Title: "Third Document", Content: "This is the third document content"},
	}

	addTask, err := index.AddDocuments(documents)
	if err != nil {
		log.Fatalf("Failed to add documents: %v", err)
	}
	fmt.Printf("✅ Document addition task created (Task ID: %d)\n", addTask.TaskUID)

	// 3. Update settings to generate another task
	fmt.Println("\n3. Updating settings to generate another task...")
	settings := &meilisearch.Settings{
		SearchableAttributes: []string{"title", "content"},
		FilterableAttributes: []string{"id"},
	}

	settingsTask, err := index.UpdateSettings(settings)
	if err != nil {
		log.Fatalf("Failed to update settings: %v", err)
	}
	fmt.Printf("✅ Settings update task created (Task ID: %d)\n", settingsTask.TaskUID)

	// 4. Get specific task information
	fmt.Println("\n4. Getting specific task information...")
	if createTask != nil {
		task, err := client.GetTask(createTask.TaskUID)
		if err != nil {
			log.Printf("Failed to get task: %v", err)
		} else {
			displayTaskInfo("Index Creation", task)
		}
	}

	task, err := client.GetTask(addTask.TaskUID)
	if err != nil {
		log.Printf("Failed to get add task: %v", err)
	} else {
		displayTaskInfo("Document Addition", task)
	}

	// 5. List all tasks with pagination
	fmt.Println("\n5. Listing all tasks...")
	tasks, err := client.GetTasks(&meilisearch.TasksQuery{
		Limit:  10,
		From:   0,
		Statuses: []string{"succeeded", "processing", "enqueued", "failed"},
	})
	if err != nil {
		log.Fatalf("Failed to get tasks: %v", err)
	}

	fmt.Printf("Found %d tasks (showing up to 10):\n", tasks.Total)
	for i, task := range tasks.Results {
		fmt.Printf("  %d. Task #%d - %s - %s (%s)\n", 
			i+1, task.UID, task.Type, task.Status, task.EnqueuedAt.Format("15:04:05"))
	}

	// 6. Filter tasks by type
	fmt.Println("\n6. Filtering tasks by type...")
	documentTasks, err := client.GetTasks(&meilisearch.TasksQuery{
		Types: []string{"documentAdditionOrUpdate"},
		Limit: 5,
	})
	if err != nil {
		log.Printf("Failed to get document tasks: %v", err)
	} else {
		fmt.Printf("Document-related tasks: %d\n", documentTasks.Total)
		for _, task := range documentTasks.Results {
			fmt.Printf("  - Task #%d: %s (%s)\n", task.UID, task.Status, task.Type)
		}
	}

	// 7. Filter tasks by index
	fmt.Println("\n7. Filtering tasks by index...")
	indexTasks, err := client.GetTasks(&meilisearch.TasksQuery{
		IndexUIDS: []string{indexUID},
		Limit:     5,
	})
	if err != nil {
		log.Printf("Failed to get index tasks: %v", err)
	} else {
		fmt.Printf("Tasks for index '%s': %d\n", indexUID, indexTasks.Total)
		for _, task := range indexTasks.Results {
			fmt.Printf("  - Task #%d: %s - %s\n", task.UID, task.Type, task.Status)
		}
	}

	// 8. Wait for specific task completion
	fmt.Println("\n8. Waiting for task completion...")
	fmt.Printf("Waiting for document addition task #%d to complete...\n", addTask.TaskUID)
	
	finalTask, err := client.WaitForTask(addTask.TaskUID, 100*time.Millisecond)
	if err != nil {
		log.Printf("Failed to wait for task: %v", err)
	} else {
		fmt.Printf("✅ Task #%d completed with status: %s\n", finalTask.UID, finalTask.Status)
		if finalTask.Error != nil {
			fmt.Printf("   Error: %v\n", finalTask.Error)
		}
		if finalTask.Duration != nil {
			fmt.Printf("   Duration: %s\n", *finalTask.Duration)
		}
	}

	// 9. Monitor multiple tasks
	fmt.Println("\n9. Monitoring multiple tasks...")
	allTaskUIDs := []int64{addTask.TaskUID, settingsTask.TaskUID}
	if createTask != nil {
		allTaskUIDs = append(allTaskUIDs, createTask.TaskUID)
	}

	for _, taskUID := range allTaskUIDs {
		task, err := client.GetTask(taskUID)
		if err != nil {
			log.Printf("Failed to get task %d: %v", taskUID, err)
			continue
		}
		
		status := "⏳"
		switch task.Status {
		case "succeeded":
			status = "✅"
		case "failed":
			status = "❌"
		case "processing":
			status = "🔄"
		}
		
		fmt.Printf("  %s Task #%d (%s): %s\n", status, task.UID, task.Type, task.Status)
	}

	// 10. Get task statistics
	fmt.Println("\n10. Task statistics...")
	allTasks, err := client.GetTasks(&meilisearch.TasksQuery{
		Limit: 100, // Get more tasks for statistics
	})
	if err != nil {
		log.Printf("Failed to get tasks for statistics: %v", err)
	} else {
		stats := calculateTaskStats(allTasks.Results)
		fmt.Printf("Task Statistics:\n")
		fmt.Printf("  - Total tasks: %d\n", stats.Total)
		fmt.Printf("  - Succeeded: %d\n", stats.Succeeded)
		fmt.Printf("  - Failed: %d\n", stats.Failed)
		fmt.Printf("  - Processing: %d\n", stats.Processing)
		fmt.Printf("  - Enqueued: %d\n", stats.Enqueued)
	}

	fmt.Println("\nTask management examples completed successfully! 🎉")
}

type TaskStats struct {
	Total      int
	Succeeded  int
	Failed     int
	Processing int
	Enqueued   int
}

func calculateTaskStats(tasks []*meilisearch.Task) TaskStats {
	stats := TaskStats{Total: len(tasks)}
	
	for _, task := range tasks {
		switch task.Status {
		case "succeeded":
			stats.Succeeded++
		case "failed":
			stats.Failed++
		case "processing":
			stats.Processing++
		case "enqueued":
			stats.Enqueued++
		}
	}
	
	return stats
}

func displayTaskInfo(name string, task *meilisearch.Task) {
	fmt.Printf("%s Task #%d:\n", name, task.UID)
	fmt.Printf("  - Type: %s\n", task.Type)
	fmt.Printf("  - Status: %s\n", task.Status)
	fmt.Printf("  - Enqueued At: %s\n", task.EnqueuedAt.Format(time.RFC3339))
	
	if task.StartedAt != nil {
		fmt.Printf("  - Started At: %s\n", task.StartedAt.Format(time.RFC3339))
	}
	
	if task.FinishedAt != nil {
		fmt.Printf("  - Finished At: %s\n", task.FinishedAt.Format(time.RFC3339))
	}
	
	if task.Duration != nil {
		fmt.Printf("  - Duration: %s\n", *task.Duration)
	}
	
	if task.Error != nil {
		fmt.Printf("  - Error: %v\n", task.Error)
	}
	fmt.Println()
}

// getenv returns the value of the environment variable named by the key,
// or def if the variable is not present or empty.
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
