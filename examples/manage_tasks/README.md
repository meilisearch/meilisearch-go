# Task Management Example

This example demonstrates comprehensive task management using the Meilisearch Go SDK:

- **Task creation** through various operations
- **Task monitoring** and status tracking
- **Task filtering** by type, status, and index
- **Task statistics** and performance monitoring
- **Async operation handling** with proper waiting

## Features Demonstrated

1. **Task Creation**: Generate tasks through various operations
2. **Task Information**: Get detailed information about specific tasks
3. **Task Listing**: List all tasks with pagination
4. **Task Filtering**: Filter tasks by:
   - Task type (documentAdditionOrUpdate, indexCreation, etc.)
   - Status (succeeded, failed, processing, enqueued)
   - Index UID
   - Date ranges
5. **Task Monitoring**: Wait for task completion
6. **Task Statistics**: Calculate task performance metrics
7. **Error Handling**: Handle task failures and errors

## Task Types Covered

- **indexCreation**: Index creation operations
- **documentAdditionOrUpdate**: Document add/update operations
- **settingsUpdate**: Settings configuration changes
- **documentDeletion**: Document deletion operations
- **indexDeletion**: Index deletion operations

## Task Statuses

- **enqueued**: Task is waiting to be processed
- **processing**: Task is currently being executed
- **succeeded**: Task completed successfully
- **failed**: Task failed with an error

## Operations Demonstrated

- `CreateIndex()` - Generate index creation tasks
- `AddDocuments()` - Generate document addition tasks
- `UpdateSettings()` - Generate settings update tasks
- `GetTask()` - Retrieve specific task information
- `GetTasks()` - List tasks with filtering and pagination
- `WaitForTask()` - Wait for task completion
- Task filtering by type, status, and index
- Task statistics calculation

## Configuration

```bash
# Set Meilisearch server URL (optional, defaults to http://localhost:7700)
export MEILI_HOST="http://localhost:7700"

# Set API key (optional, but recommended for production)
export MEILI_API_KEY="your-api-key"
```

## Running the Example

```bash
go run ./examples/manage_tasks
```

The example will create various tasks, demonstrate different ways to monitor and filter them, and show how to properly handle asynchronous operations in Meilisearch.

## Best Practices Shown

- Always wait for critical tasks to complete
- Use appropriate timeouts for task waiting
- Filter tasks efficiently for monitoring
- Handle task errors gracefully
- Monitor task performance with statistics
