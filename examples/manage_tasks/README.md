# Task Management Example

This example shows how to monitor and manage Meilisearch tasks.

## What it does

1. Create an index (generates a task)
2. Add documents (generates a task)  
3. Update settings (generates a task)
4. Wait for tasks to complete
5. Get individual task details
6. List all tasks with filters
7. Filter tasks by type (document operations only)

## Task information shown

For each task you see:
- Task ID and type
- Status (enqueued, processing, succeeded, failed)
- Start and finish times
- Duration
- Error details if it failed

## Configuration

```bash
export MEILI_HOST="http://localhost:7700"
export MEILI_API_KEY="your-admin-api-key"
```

Note: You need an admin API key (not a search key) to manage tasks.

## Run it

```bash
go run ./examples/manage_tasks
```
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
