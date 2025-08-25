# Document Management Example

This example demonstrates comprehensive document management using the Meilisearch Go SDK:

- **Adding documents** (single and batch operations)
- **Updating documents** with new data
- **Retrieving documents** by ID and in batches
- **Deleting documents** (single and multiple)
- **Task monitoring** for all operations

## Features Demonstrated

1. **Batch Operations**: Add multiple documents efficiently
2. **Single Document Operations**: Add, update, delete individual documents
3. **Document Retrieval**: Get documents by ID or fetch in batches
4. **Update Operations**: Modify existing documents
5. **Delete Operations**: Remove documents by ID or in batches
6. **Task Management**: Monitor async operations completion
7. **Error Handling**: Robust error handling for all operations

## Operations Covered

- `AddDocuments()` - Add documents to index
- `UpdateDocuments()` - Update existing documents
- `GetDocument()` - Retrieve single document by ID
- `GetDocuments()` - Retrieve multiple documents with pagination
- `DeleteDocument()` - Delete single document
- `DeleteDocuments()` - Delete multiple documents
- `WaitForTask()` - Monitor task completion

## Configuration

```bash
# Set Meilisearch server URL (optional, defaults to http://localhost:7700)
export MEILI_HOST="http://localhost:7700"

# Set API key (optional, but recommended for production)
export MEILI_API_KEY="your-api-key"
```

## Running the Example

```bash
go run ./examples/add_documents
```

The example will create a "users" index and demonstrate all document management operations with detailed logging of each step and task completion.
