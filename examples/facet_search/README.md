# Facet Search Example

This example shows how to search with facets (filters) using a book catalog.

## What it does

Uses a Book struct with these fields:

```go
type Book struct {
    ID          string  `json:"id"`
    Title       string  `json:"title"`
    Author      string  `json:"author"`
    Genre       string  `json:"genre"`
    Language    string  `json:"language"`
    PublishYear int     `json:"publish_year"`
    Rating      float64 `json:"rating"`
    Publisher   string  `json:"publisher"`
    Pages       int     `json:"pages"`
    InPrint     bool    `json:"in_print"`
    Series      string  `json:"series"`
}
```

1. Create a "books" index with 8 sample books
2. Configure filterable attributes (genre, language, rating, etc.)
3. Search with facets to get result counts by category
4. Use advanced filters like "genre = fantasy AND rating > 4.0"
5. Search within specific facets (like finding "sci" in genres)

## Configuration

```bash
export MEILI_HOST="http://localhost:7700"
export MEILI_API_KEY="your-api-key"
```

## Run it

```bash
go run ./examples/facet_search
```
```
Testing connection to Meilisearch...
✅ Connected to Meilisearch (status: available)

📚 Setting up books index with facetable attributes...
✅ Added 8 books to the index

1. Basic faceted search with distribution:
Search: 'fiction' - Found 4 results
Facet distribution:
  genre: {"fiction": 1, "sci-fi": 2, "cyberpunk": 1}
  language: {"english": 4}
  publish_year: {"1949": 1, "1965": 1, "1984": 1, "1992": 1}

2. Faceted search with filters:
Search: '' with filter 'genre = fantasy AND publish_year > 2000'
Found 0 results (no fantasy books after 2000 in dataset)
Refining to: 'genre = sci-fi AND rating > 4.0'
Found 3 results with high-rated sci-fi books

3. Facet-specific search:
Facet search for 'sci' in genre facet with query 'space':
Found 2 matching facet values: ["sci-fi"]

Faceted search example completed successfully! 🎉

The example will create a "books" index, configure facetable attributes, add sample book data, and demonstrate various faceted search scenarios including basic faceting, advanced filtering, and facet-specific searches.
