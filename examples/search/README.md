# Basic Search Example

This example demonstrates the core functionality of the Meilisearch Go SDK:

- **Client initialization** with connection testing
- **Index creation** with primary key configuration
- **Document management** (adding documents to an index)
- **Search operations** with various parameters and filters
- **Task management** (waiting for async operations to complete)

## Features Demonstrated

1. **Basic Setup**: Initialize client and test connection
2. **Index Management**: Create an index with proper configuration
3. **Document Operations**: Add structured documents (movies) to the index
4. **Search Capabilities**:
   - Simple text search
   - Filtered search with conditions
   - Faceted search with distribution
   - Highlighted results
   - Pagination controls

## Prerequisites

- Go 1.20 or higher
- Meilisearch server running (default: `http://localhost:7700`)
- Valid API key (if authentication is enabled)

## Running the Example

1. **Start Meilisearch server:**
   ```bash
   # Using Docker
   docker run -it --rm -p 7700:7700 getmeili/meilisearch:latest
   
   # Or download and run directly
   ./meilisearch
   ```

2. **Update configuration:**
   Edit `main.go` to match your Meilisearch server URL and API key.

3. **Run the example:**
   ```bash
   go run ./examples/search
   ```

The example will create a "movies" index, add sample movie documents, and demonstrate various search operations with detailed output.
