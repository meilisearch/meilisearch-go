# Multi-Search Example

This example demonstrates how to perform multiple search queries simultaneously using the Meilisearch Go SDK:

- **Simultaneous queries** across indexes
- **Different search parameters** for each query
- **Batch processing** of search requests
- **Performance optimization** through parallel execution

## Features Demonstrated

1. **Multiple Query Types**: Different search queries with various parameters
2. **Mixed Filters**: Each query can have its own filtering conditions
3. **Different Sorting**: Queries can use different sorting strategies
4. **Batch Results**: Process multiple search results efficiently
5. **Performance**: Execute multiple searches in a single request

## Use Cases

- **Dashboard searches**: Multiple widgets requiring different data
- **Comparison searches**: Side-by-side search results
- **Related content**: Find related items with different criteria
- **Analytics**: Gather multiple search metrics in one call

## Configuration

```bash
# Set Meilisearch server URL (optional, defaults to http://localhost:7700)
export MEILI_HOST="http://localhost:7700"

# Set API key (optional, but recommended for production)
export MEILI_API_KEY="your-api-key"
```

## Running the Example

```bash
# Set environment variables (optional)
export MEILI_HOST="http://localhost:7700"
export MEILI_API_KEY="your-api-key"

# Run the example
go run ./examples/multi_search
```

The example will create a "products" index, add sample product data, and demonstrate multiple simultaneous search queries with different parameters, filters, and sorting options.
