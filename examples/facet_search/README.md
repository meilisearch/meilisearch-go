# Faceted Search Example

This example demonstrates advanced faceted search capabilities using the Meilisearch Go SDK:

- **Facet distribution** for search results
- **Filterable attributes** configuration
- **Facet-specific searches** with targeted queries
- **Combined filtering** with multiple facets

## Features Demonstrated

1. **Basic Faceting**: Get facet distributions alongside search results
2. **Advanced Filtering**: Combine multiple facet filters in searches
3. **Facet Search**: Search within specific facet values
4. **Index Configuration**: Set up filterable and sortable attributes
5. **Real-world Data**: Books dataset with multiple facetable dimensions

## Faceted Attributes

The example configures these facetable attributes:
- **genre**: Book categories (fiction, fantasy, sci-fi, etc.)
- **language**: Publication language
- **publish_year**: Year of publication
- **rating**: Book rating
- **publisher**: Publishing house
- **in_print**: Availability status

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
go run ./examples/facet_search
```

The example will create a "books" index, configure facetable attributes, add sample book data, and demonstrate various faceted search scenarios including basic faceting, advanced filtering, and facet-specific searches.
