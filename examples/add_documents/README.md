# Document Management Example

This example shows how to manage documents in Meilisearch. You'll learn to add, update, retrieve, and delete documents.

## What it does

```go
type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    Active   bool   `json:"active"`
    JoinDate string `json:"join_date"`
}
```

1. Create a "users" index
2. Add multiple documents at once
3. Add single documents
4. Update existing documents
5. Get documents by ID
6. Get multiple documents with pagination
7. Delete single document
8. Delete multiple documents

## Configuration

```bash
# Set Meilisearch server URL (defaults to http://localhost:7700)
export MEILI_HOST="http://localhost:7700"

# Set API key (recommended for production)
export MEILI_API_KEY="your-api-key"
```

## Run it

```bash
go run ./examples/add_documents
```
