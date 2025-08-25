# Faceted Search Example

This example demonstrates advanced faceted search capabilities using the Meilisearch Go SDK with a comprehensive book catalog:

- **Facet distribution analysis** alongside search results
- **Multi-facet filtering** with complex boolean logic
- **Facet-specific searches** for targeted facet exploration
- **Advanced filter combinations** with category, rating, and date filters
- **Comprehensive book catalog** with rich facetable attributes
- **JSON facet parsing** with proper error handling

## Book Catalog Structure

The example uses a `Book` struct representing book documents:

```go
type Book struct {
    ID          string   `json:"id"`           // Primary key (ISBN or unique ID)
    Title       string   `json:"title"`        // Book title
    Author      string   `json:"author"`       // Primary author
    Genre       string   `json:"genre"`        // Book genre/category
    Language    string   `json:"language"`     // Publication language
    PublishYear int      `json:"publish_year"` // Year of publication
    Rating      float64  `json:"rating"`       // Average rating (1.0-5.0)
    Publisher   string   `json:"publisher"`    // Publishing house
    Pages       int      `json:"pages"`        // Number of pages
    InPrint     bool     `json:"in_print"`     // Current availability
    Series      string   `json:"series"`       // Book series (if applicable)
}
```

## Sample Book Data

The example includes 8 diverse books across multiple genres:
- **"The Hobbit"** - J.R.R. Tolkien, Fantasy, 1937, 4.8 rating, In Print
- **"Dune"** - Frank Herbert, Sci-Fi, 1965, 4.6 rating, In Print
- **"1984"** - George Orwell, Fiction, 1949, 4.7 rating, In Print
- **"Foundation"** - Isaac Asimov, Sci-Fi, 1951, 4.5 rating, Out of Print
- **"Neuromancer"** - William Gibson, Cyberpunk, 1984, 4.4 rating, In Print
- **"The Matrix"** - Simulacra, Sci-Fi, 1999, 4.2 rating, In Print
- **"Blade Runner"** - Philip K. Dick, Cyberpunk, 1968, 4.3 rating, Out of Print
- **"Snow Crash"** - Neal Stephenson, Cyberpunk, 1992, 4.1 rating, In Print

## Faceted Search Operations

### **1. Index Setup and Configuration**
```go
// Configure filterable attributes for faceting
settingsTask, err := index.UpdateSettings(&meilisearch.Settings{
    FilterableAttributes: []string{"genre", "language", "publish_year", "rating", "publisher", "in_print", "series"},
    SortableAttributes:   []string{"rating", "publish_year", "pages"},
})
```

#### **Filterable Attributes (Facet Capabilities)**
- **genre**: Book categories (fantasy, sci-fi, fiction, cyberpunk, etc.)
- **language**: Publication language (English, Spanish, French, etc.)
- **publish_year**: Year of publication for temporal filtering
- **rating**: Numeric rating for quality-based filtering
- **publisher**: Publishing house for brand-based filtering
- **in_print**: Availability status (true/false)
- **series**: Book series for collection-based filtering

### **2. Basic Faceted Search with Distribution**
```go
searchResult, err := client.Index("books").Search("fiction", &meilisearch.SearchRequest{
    Facets: []string{"genre", "language", "publish_year"},
    Limit:  5,
})
```

#### **Facet Distribution Analysis**
```go
func displaySearchResults(query string, result *meilisearch.SearchResponse) {
    fmt.Printf("Search: '%s' - Found %d results\n", query, result.EstimatedTotalHits)
    if len(result.FacetDistribution) > 0 {
        fmt.Println("Facet distribution:")
        var fd map[string]map[string]int
        if err := json.Unmarshal(result.FacetDistribution, &fd); err == nil {
            for facet, distribution := range fd {
                fmt.Printf("  %s: %v\n", facet, distribution)
            }
        }
    }
}
```

**Example Distribution Output:**
```
Facet distribution:
  genre: {"fiction": 3, "fantasy": 1, "sci-fi": 2}
  language: {"english": 5, "spanish": 1}
  publish_year: {"1949": 1, "1965": 1, "1984": 1, "1992": 1, "1999": 1}
```

### **3. Advanced Filtering with Multiple Facets**
```go
searchResult, err = client.Index("books").Search("", &meilisearch.SearchRequest{
    Filter: "genre = fantasy AND publish_year > 2000",
    Facets: []string{"language", "rating", "publisher"},
    Sort:   []string{"rating:desc"},
    Limit:  10,
})
```

#### **Advanced Filter Syntax**
- **Boolean Logic**: `AND`, `OR` operators for complex conditions
- **Comparison Operators**: `>`, `<`, `>=`, `<=`, `=`, `!=`
- **Numeric Ranges**: `rating > 4.0 AND rating < 5.0`
- **String Matching**: `genre = "fantasy"` (exact match)
- **Multiple Conditions**: `genre = fantasy AND publish_year > 2000 AND in_print = true`

### **4. Facet-Specific Search**
```go
facetResult, err := client.Index("books").FacetSearch(&meilisearch.FacetSearchRequest{
    FacetName:  "genre",
    FacetQuery: "sci",
    Q:          "space",
})
```

#### **FacetSearch Parameters**
- **FacetName**: Target facet attribute ("genre", "author", etc.)
- **FacetQuery**: Search within facet values ("sci" matches "sci-fi")
- **Q**: Optional main query to narrow context ("space" for space-related books)
- **Filter**: Additional filters to apply during facet search

#### **Use Cases for Facet Search**
- **Auto-complete**: Suggest genre names as user types "sci" → "sci-fi"
- **Facet Exploration**: Discover available values within a facet
- **Contextual Faceting**: Find facet values relevant to current search
- **Facet Validation**: Check if specific facet values exist

## Facet Distribution Benefits

### **Search Refinement**
- **Category Counts**: "Show 15 sci-fi books, 8 fantasy books"
- **Filter Guidance**: Help users understand available filter options
- **Search Navigation**: Enable drill-down search experiences
- **Result Preview**: Show result distribution before applying filters

### **User Experience Enhancement**
- **Interactive Filtering**: Click facet values to apply filters
- **Search Faceted Navigation**: Build category-based browsing
- **Result Analytics**: Show search result breakdowns
- **Dynamic UI**: Update available filters based on current results

## Advanced Filtering Patterns

### **Complex Boolean Logic**
```go
// Fantasy OR Sci-Fi books published after 1980 with high ratings
Filter: "(genre = fantasy OR genre = sci-fi) AND publish_year > 1980 AND rating > 4.0"

// Available books under 400 pages from specific publishers
Filter: "in_print = true AND pages < 400 AND (publisher = 'Penguin' OR publisher = 'HarperCollins')"

// Recent cyberpunk books with ratings between 4.0 and 5.0
Filter: "genre = cyberpunk AND publish_year > 1990 AND rating >= 4.0 AND rating <= 5.0"
```

### **Numeric Range Filtering**
```go
// Books from the golden age of sci-fi (1938-1946)
Filter: "genre = sci-fi AND publish_year >= 1938 AND publish_year <= 1946"

// Highly rated recent releases
Filter: "publish_year > 2010 AND rating > 4.5"

// Medium-length books for casual reading
Filter: "pages >= 200 AND pages <= 350"
```

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

# Run the faceted search example
go run ./examples/facet_search
```

### **Expected Output**
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
