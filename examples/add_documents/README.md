# Document Management Example

This example demonstrates comprehensive document management using the Meilisearch Go SDK with a User management system:

- **Document addition** - Add single documents and batch operations
- **Document updates** - Update existing documents with new information
- **Document retrieval** - Get documents by ID and bulk document queries
- **Document deletion** - Delete single and multiple documents
- **Index management** - Create index with primary key configuration
- **Environment-based configuration** - Flexible deployment setup

## User Document Structure

The example uses a `User` struct representing user documents:

```go
type User struct {
    ID       int    `json:"id"`        // Primary key
    Name     string `json:"name"`      // User's full name
    Email    string `json:"email"`     // User's email address
    Role     string `json:"role"`      // admin, user, moderator
    Active   bool   `json:"active"`    // Account status
    JoinDate string `json:"join_date"` // Registration date
}
```

## Operations Demonstrated

### 1. **Index Creation**
- Creates "users" index with "id" as primary key
- Handles existing index scenarios gracefully

### 2. **Batch Document Addition**
- Adds multiple User documents in a single operation
- Sample users: Alice (admin), Bob (user), Charlie (moderator)

### 3. **Single Document Addition**
- Adds individual User document (Diana Prince)
- Demonstrates single document workflow

### 4. **Document Updates**
- Updates existing users (Bob becomes admin, Charlie activated)
- Partial updates while preserving other fields

### 5. **Document Retrieval**
- **By ID**: Get specific user by ID (User #2)
- **Bulk Query**: Retrieve multiple documents with pagination
- Displays retrieved document details

### 6. **Single Document Deletion**
- Deletes specific document by ID (User #4)
- Task-based operation with completion waiting

### 7. **Multiple Document Deletion**
- Deletes multiple documents by ID list (Users #1 and #3)
- Batch deletion for efficiency

### 8. **Final Document Count**
- Retrieves remaining documents after operations
- Displays final state of the index

## Document Operations Covered

- `AddDocuments(users)` - Batch add multiple User documents
- `AddDocuments([]User{newUser})` - Add single User document
- `UpdateDocuments(updatedUsers)` - Update existing User records
- `GetDocument("2", nil)` - Retrieve User by ID
- `GetDocuments(&DocumentsQuery{...})` - Bulk retrieve with pagination
- `DeleteDocument("4")` - Delete single User by ID
- `DeleteDocuments([]string{"1", "3"})` - Delete multiple Users by ID list

## Task Management Features

- **Task Waiting**: All operations wait for completion using `waitForTask()`
- **Task Monitoring**: Displays task IDs for tracking
- **Error Handling**: Graceful handling of task failures
- **Timeout Management**: 10-second timeout for task completion

## Configuration

```bash
# Set Meilisearch server URL (defaults to http://localhost:7700)
export MEILI_HOST="http://localhost:7700"

# Set API key (recommended for production)
export MEILI_API_KEY="your-api-key"
```

## Running the Example

```bash
# Set environment variables (optional)
export MEILI_HOST="http://localhost:7700"
export MEILI_API_KEY="your-api-key"

# Run the example
go run ./examples/add_documents
```

The example will create a "users" index, add sample user data, demonstrate various document operations, and show the final state after all operations.

### **Expected Output**
```
Testing connection to Meilisearch...
✅ Connected to Meilisearch

Creating index 'users'...
✅ Created index 'users' (Task ID: 123)

📄 Document Management Examples
===============================

1. Adding documents in batch:
✅ Added 3 documents (Task ID: 124)

2. Adding single document:
✅ Added single document (Task ID: 125)

3. Updating documents:
✅ Updated 2 documents (Task ID: 126)

4. Getting document by ID:
Retrieved document: map[id:2 name:Bob Smith Jr. email:bob.jr@example.com ...]

5. Getting multiple documents:
Retrieved 4 documents:
  1. map[id:1 name:Alice Johnson ...]
  2. map[id:2 name:Bob Smith Jr. ...]
  ...

6. Deleting documents:
✅ Deleted document with ID 4 (Task ID: 127)

7. Deleting multiple documents:
✅ Deleted multiple documents (Task ID: 128)

8. Final document count:
Remaining documents: 2

Document management examples completed successfully! 🎉
