# Basic Search Example

This example demonstrates comprehensive search functionality using the Meilisearch Go SDK with a movie database:

- **Index creation** with primary key specification
- **Settings configuration** for optimal search performance
- **Document indexing** with movie data
- **Simple text search** with basic queries
- **Advanced search** with filters, facets, and highlighting
- **Facet distribution** parsing and display

## Movie Database Structure

The example uses a `Movie` struct representing movie documents:

```go
type Movie struct {
    ID     string   `json:"id"`     // Primary key
    Title  string   `json:"title"`  // Movie title
    Year   int      `json:"year"`   // Release year
    Rating float64  `json:"rating"` // IMDb rating
    Genres []string `json:"genres"` // Movie genres
}
```

## Sample Movie Data

The example includes 5 classic movies:
- **The Dark Knight** (2008) - Action, Crime, Drama - Rating: 9.0
- **Inception** (2010) - Action, Sci-Fi, Thriller - Rating: 8.8  
- **The Godfather** (1972) - Crime, Drama - Rating: 9.2
- **Pulp Fiction** (1994) - Crime, Drama - Rating: 8.9
- **Fight Club** (1999) - Drama - Rating: 8.8

## Search Operations Demonstrated

### **1. Index Setup**
- Creates "movies" index with "id" as primary key
- Configures filterable attributes: `["genres", "year", "rating"]`
- Configures sortable attributes: `["year", "rating"]`
- Waits for settings to be applied before indexing

### **2. Document Indexing**
- Adds 5 movie documents to the index
- Uses `AddDocuments()` for batch operation
- Waits for indexing task completion with timeout
- Displays success confirmation with task ID

### **3. Simple Search**
- **Query**: "action" 
- **Results**: Movies containing "action" in any field
- **Limit**: 5 results maximum
- **Display**: Shows movie title, year, and rating

### **4. Advanced Search with Filters**
- **Query**: "drama"
- **Filter**: `year > 1990` (movies after 1990)
- **Facets**: `["genres", "year"]` for distribution analysis
- **Highlighting**: Highlights matches in title and overview
- **Limit**: 10 results maximum

### **5. Facet Distribution Analysis**
- Parses `FacetDistribution` JSON response
- Shows genre and year distribution statistics
- Handles JSON parsing errors gracefully
- Displays facet counts for search results

## Search Features Covered

### **SearchRequest Parameters**
- `Filter`: Filter expression (`"year > 1990"`)
- `Facets`: Facet fields for distribution (`["genres", "year"]`)
- `AttributesToHighlight`: Fields to highlight (`["title", "overview"]`)
- `Limit`: Maximum number of results (5, 10)

### **Response Processing**
- `searchResult.Hits`: Array of matching documents
- `searchResult.FacetDistribution`: JSON raw message for facet data
- Error handling for search operations
- Type conversion for numeric fields (`hit["year"].(float64)`)

## Configuration

```bash
# Set Meilisearch server URL (defaults to http://localhost:7700)
export MEILI_HOST="http://localhost:7700"

# Set API key (optional, recommended for production)
export MEILI_API_KEY="your-api-key"
```

## Running the Example

```bash
# Set environment variables (optional)
export MEILI_HOST="http://localhost:7700"
export MEILI_API_KEY="your-api-key"

# Run the search example
go run ./examples/search
```

### **Expected Output**
```
Testing connection to Meilisearch...
✅ Connected to Meilisearch

Adding 5 movies to the index...
✅ Indexed documents successfully! (Task ID: 123)

1. Simple search for 'action':
Found 2 results
  1. The Dark Knight (2008) - Rating: 9.0
  2. Inception (2010) - Rating: 8.8

2. Advanced search with filters and facets:
Found 3 drama movies after 1990
  1. Pulp Fiction (1994) - Rating: 8.9
  2. Fight Club (1999) - Rating: 8.8

Facet distribution:
  genres: map[crime:1 drama:3]
  year: map[1994:1 1999:1]

Search example completed successfully!
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
