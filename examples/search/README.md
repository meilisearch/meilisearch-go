# Search Example

This example shows how to search documents in Meilisearch with basic and advanced search features.

## What it does

```go
type Movie struct {
    ID     string   `json:"id"`
    Title  string   `json:"title"`
    Year   int      `json:"year"`
    Rating float64  `json:"rating"`
    Genres []string `json:"genres"`
}
```

1. Create a "movies" index
2. Configure search settings (filterable and sortable attributes)
3. Index sample movie data
4. Simple search for "action" movies
5. Advanced search with filters and facets for "drama" movies

## Configuration

```bash
export MEILI_HOST="http://localhost:7700"
export MEILI_API_KEY="your-api-key"
```

## Run it

```bash
go run ./examples/search
```
```

## Best Practices Demonstrated

- **Settings Configuration**: Configure filterable/sortable attributes before indexing
- **Task Completion**: Always wait for indexing and settings tasks
- **Error Handling**: Comprehensive error handling throughout operations
- **Resource Management**: Proper client cleanup with `defer client.Close()`
- **Facet Processing**: Safe JSON unmarshalling with error handling
- **Environment Configuration**: Flexible host and API key setup

## Advanced Usage

The example demonstrates production-ready patterns:

- **Timeout Management**: 10-second timeout for indexing operations
- **Settings First**: Apply index settings before adding documents
- **Batch Operations**: Efficient bulk document indexing
- **Type Safety**: Proper type conversion for search results
- **JSON Handling**: Safe parsing of facet distribution data

## Troubleshooting

- Ensure Meilisearch server is running on the configured host
- Verify API key permissions if using authentication
- Check that filterable attributes are configured before filtering
- Confirm documents are indexed before searching
