# Task Management Example

This example demonstrates comprehensive task management and monitoring using the Meilisearch Go SDK with real-time task tracking:

- **Task creation** through various Meilisearch operations
- **Real-time task monitoring** with status tracking and completion waiting
- **Advanced task filtering** by type, status, index, and date ranges
- **Task statistics calculation** and performance monitoring
- **Asynchronous operation handling** with proper timeout management
- **Error handling** for failed tasks with detailed error reporting

## Task Management Operations

### **1. Index Creation Task**
- Creates "task_demo" index with primary key "id"
- **Task Type**: `indexCreation`
- **Expected Status**: `succeeded` (quick operation)
- Demonstrates basic async operation with task generation

### **2. Document Addition Task**
```go
documents := []map[string]interface{}{
    {"id": 1, "title": "Task Management Guide", "category": "tutorial"},
    {"id": 2, "title": "Advanced Meilisearch", "category": "guide"},
    {"id": 3, "title": "Performance Optimization", "category": "advanced"},
}
```
- **Task Type**: `documentAdditionOrUpdate`
- **Expected Status**: `succeeded` (with processing time)
- Shows document indexing with task monitoring

### **3. Settings Update Task**
```go
settings := &meilisearch.Settings{
    FilterableAttributes: []string{"category"},
    SortableAttributes:   []string{"title"},
}
```
- **Task Type**: `settingsUpdate`
- **Expected Status**: `succeeded` (configuration change)
- Demonstrates settings modification with task tracking

## Task Information Display

### **Individual Task Details**
For each task, the example displays:
```go
func displayTaskInfo(taskName string, task *meilisearch.Task) {
    fmt.Printf("Task #%d:\n", task.UID)
    fmt.Printf("  - Type: %s\n", task.Type)
    fmt.Printf("  - Status: %s\n", task.Status)

    if !task.StartedAt.IsZero() {
        fmt.Printf("  - Started At: %s\n", task.StartedAt.Format(time.RFC3339))
    }
    if !task.FinishedAt.IsZero() {
        fmt.Printf("  - Finished At: %s\n", task.FinishedAt.Format(time.RFC3339))
    }
    if task.Duration != "" {
        fmt.Printf("  - Duration: %s\n", task.Duration)
    }
    if task.Status == meilisearch.TaskStatus("failed") {
        fmt.Printf("  - Error: %v\n", task.Error)
    }
}
```

### **Task Status Types**
- **enqueued**: Task queued and waiting for processing
- **processing**: Task currently being executed by Meilisearch  
- **succeeded**: Task completed successfully without errors
- **failed**: Task failed with error details available

### **Task Type Categories**
- **indexCreation**: Index creation and initialization
- **documentAdditionOrUpdate**: Document indexing and updates
- **settingsUpdate**: Index settings and configuration changes
- **documentDeletion**: Document removal operations
- **indexDeletion**: Index deletion operations
- **taskCancelation**: Task cancellation operations
- **taskDeletion**: Task cleanup operations

## Advanced Task Filtering

### **4. Task Listing with Status Filtering**
```go
tasks, err := client.GetTasks(&meilisearch.TasksQuery{
    Limit: 10,
    From:  0,
    Statuses: []meilisearch.TaskStatus{
        meilisearch.TaskStatus("succeeded"),
        meilisearch.TaskStatus("processing"),
        meilisearch.TaskStatus("enqueued"),
        meilisearch.TaskStatus("failed"),
    },
})
```
- **Limit**: Maximum number of tasks to retrieve (pagination)
- **From**: Starting task UID for pagination  
- **Statuses**: Filter by multiple task statuses simultaneously
- **Use Case**: Monitor tasks across different completion states

### **5. Task Type Filtering**
```go
documentTasks, err := client.GetTasks(&meilisearch.TasksQuery{
    Types: []meilisearch.TaskType{
        meilisearch.TaskType("documentAdditionOrUpdate")
    },
    Limit: 5,
})
```
- **Types**: Filter by specific task types
- **Use Case**: Monitor document-related operations specifically
- **Benefits**: Focus on relevant operation categories

## Task Completion Monitoring

### **6. WaitForTask Implementation**
```go
finalTask, err := client.WaitForTask(addTask.TaskUID, 100*time.Millisecond)
if err != nil {
    log.Printf("Failed to wait for task: %v", err)
} else {
    fmt.Printf("✅ Task #%d completed with status: %s\n", finalTask.UID, finalTask.Status)
    if finalTask.Status == meilisearch.TaskStatus("failed") {
        fmt.Printf("   Error: %v\n", finalTask.Error)
    }
    if finalTask.Duration != "" {
        fmt.Printf("   Duration: %s\n", finalTask.Duration)
    }
}
```

#### **WaitForTask Parameters**
- **TaskUID**: Unique identifier of task to monitor
- **Interval**: Polling interval (100ms for responsive monitoring)
- **Returns**: Final task state with completion details
- **Error Handling**: Comprehensive error reporting for failures

## Task Statistics and Analysis

### **Performance Metrics**
- **Task Duration**: Execution time for completed tasks
- **Success Rate**: Percentage of succeeded vs failed tasks
- **Task Distribution**: Count by type and status
- **Processing Patterns**: Peak times and operation frequencies

### **Error Analysis**
- **Failed Task Details**: Complete error information
- **Failure Patterns**: Common failure types and causes
- **Recovery Strategies**: Retry logic and error handling

## Configuration

```bash
# Set Meilisearch server URL (optional, defaults to http://localhost:7700)
export MEILI_HOST="http://localhost:7700"

# Set API key (optional, but recommended for production)
export MEILI_API_KEY="your-api-key"
```

**⚠️ Important**: Task management operations generally require an **Admin API key** (not a Search key). Ensure your API key has appropriate permissions for:
- Task listing and filtering
- Task status monitoring  
- Task cancellation (if needed)

## Running the Example

```bash
# Set environment variables (optional)
export MEILI_HOST="http://localhost:7700"
export MEILI_API_KEY="your-admin-api-key"

# Run the task management example
go run ./examples/manage_tasks
```

### **Expected Output**
```
Testing connection to Meilisearch...
✅ Connected to Meilisearch (status: available)

1. Creating index for task demonstration...
✅ Task #1 (Index Creation) completed with status: succeeded
   Duration: 2ms

2. Adding documents to generate tasks...
✅ Task #2 (Document Addition) completed with status: succeeded  
   Duration: 15ms

3. Updating settings...
✅ Task #3 (Settings Update) completed with status: succeeded
   Duration: 8ms

4. Retrieving individual task information...
Task #2:
  - Type: documentAdditionOrUpdate
  - Status: succeeded
  - Started At: 2024-01-15T10:30:45Z
  - Finished At: 2024-01-15T10:30:45Z
  - Duration: 15ms

5. Listing all tasks...
Found 3 tasks total

6. Filtering tasks by type...
Found 1 document-related tasks

Task management example completed successfully! 🎉

## Best Practices Shown

- Always wait for critical tasks to complete
- Use appropriate timeouts for task waiting
- Filter tasks efficiently for monitoring
- Handle task errors gracefully
- Monitor task performance with statistics
