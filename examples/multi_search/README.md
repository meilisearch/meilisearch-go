# Multi-Search Example

This example shows how to run multiple searches at once using a single API call.

## What it does

```go
type Product struct {
    ID          int      `json:"id"`
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Category    string   `json:"category"`
    Price       float64  `json:"price"`
    Brand       string   `json:"brand"`
    Tags        []string `json:"tags"`
    InStock     bool     `json:"in_stock"`
}
```

1. Create a "products" index with sample product data
2. Configure filterable attributes for categories, stock status, price, and brand
3. Run 3 searches simultaneously:
   - Search for "laptop" in electronics category
   - Search for "coffee" that's in stock
   - Search for products under $100, sorted by price

## Configuration

```bash
export MEILI_HOST="http://localhost:7700"
export MEILI_API_KEY="your-api-key"
```

## Run it

```bash
go run ./examples/multi_search
```
```

## Best Practices Demonstrated

- **Settings Configuration**: Configure filterable/sortable attributes before searching
- **Task Completion**: Wait for index setup before executing searches
- **Error Handling**: Comprehensive error handling for multi-search operations
- **Resource Management**: Proper client cleanup with `defer client.Close()`
- **Filter Syntax**: Proper filter expression formatting
- **Performance Optimization**: Batch multiple searches for efficiency
