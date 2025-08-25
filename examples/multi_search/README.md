# Multi-Search Example

This example demonstrates simultaneous multi-index searching using the Meilisearch Go SDK with a product catalog:

- **Multiple concurrent searches** with different queries and filters
- **Single API call** for multiple search operations
- **Product index setup** with comprehensive configuration
- **Filter combinations** with category, stock status, and price ranges
- **Sorting operations** with price-based ordering
- **Performance optimization** through batched search requests

## Product Catalog Structure

The example uses a `Product` struct representing product documents:

```go
type Product struct {
    ID          int      `json:"id"`          // Primary key
    Name        string   `json:"name"`        // Product name
    Description string   `json:"description"` // Product description
    Category    string   `json:"category"`    // Product category
    Price       float64  `json:"price"`       // Product price
    Brand       string   `json:"brand"`       // Brand name
    Tags        []string `json:"tags"`        // Search tags
    InStock     bool     `json:"in_stock"`    // Stock status
}
```

## Sample Product Data

The example includes 6 diverse products:
- **Gaming Laptop** - TechBrand, $1299.99, Electronics, In Stock
- **Organic Coffee** - CoffeeCorp, $15.99, Food, Out of Stock  
- **Wireless Mouse** - TechBrand, $29.99, Electronics, In Stock
- **Coffee Beans** - CoffeeCorp, $24.99, Food, Out of Stock
- **Bluetooth Headphones** - AudioTech, $199.99, Electronics, In Stock
- **Desk Chair** - FurniturePlus, $249.99, Furniture, In Stock

## Multi-Search Operations

### **1. Index Setup**
- Creates "products" index with "id" as primary key
- **Filterable attributes**: `["category", "in_stock", "price", "brand"]`
- **Sortable attributes**: `["price"]`
- Waits for settings configuration before indexing

### **2. Concurrent Search Queries**

The example executes **3 simultaneous searches** in a single API call:

#### **Query 1: Electronics Laptop Search**
```go
{
    IndexUID: "products",
    Query:    "laptop",
    Limit:    5,
    Filter:   []string{`category = "electronics"`},
}
```
- **Purpose**: Find laptop products in electronics category
- **Expected**: Gaming Laptop results

#### **Query 2: Available Coffee Search**  
```go
{
    IndexUID: "products",
    Query:    "coffee", 
    Limit:    3,
    Filter:   []string{"in_stock = true"},
}
```
- **Purpose**: Find coffee products currently in stock
- **Expected**: Available coffee items only

#### **Query 3: Budget Products with Sorting**
```go
{
    IndexUID: "products",
    Query:    "",           // Empty query = all products
    Limit:    10,
    Filter:   []string{"price < 100"},
    Sort:     []string{"price:asc"},
}
```
- **Purpose**: Find affordable products under $100, sorted by price
- **Expected**: Wireless Mouse ($29.99), etc., sorted by price ascending

## Multi-Search Benefits

### **Performance Advantages**
- **Single API call**: Reduces network round trips
- **Concurrent execution**: Searches run simultaneously on server
- **Reduced latency**: No sequential search delays
- **Resource efficiency**: Optimal server resource utilization

### **Use Cases Demonstrated**
- **Category-based filtering**: Electronics vs Food vs Furniture
- **Stock status filtering**: Available vs Out of Stock
- **Price range filtering**: Budget-conscious shopping
- **Sorting integration**: Price-based product ordering
- **Mixed query types**: Text search vs filtered browsing

## Results Processing

### **MultiSearchResponse Structure**
```go
results, err := client.MultiSearch(multiSearchRequest)
// results.Results[0] = Laptop search results
// results.Results[1] = Coffee search results  
// results.Results[2] = Budget products results
```

### **Individual Result Display**
For each search result:
- **Query identification**: Shows search query and filters used
- **Hit count**: Number of matching products found
- **Product details**: ID, Name, Price, Brand, Stock status
- **Performance metrics**: Search execution insights

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

# Run the multi-search example
go run ./examples/multi_search
```

### **Expected Output**
```
Testing connection to Meilisearch...
✅ Connected to Meilisearch

📊 Executed 3 search queries simultaneously:

Query 1: "laptop" with filter 'category = "electronics"'
Found 1 results:
  1. Gaming Laptop ($1299.99) - TechBrand [In Stock: true]

Query 2: "coffee" with filter 'in_stock = true'  
Found 0 results:
  (No coffee products currently in stock)

Query 3: "" with filter 'price < 100' sorted by price:asc
Found 1 results:
  1. Wireless Mouse ($29.99) - TechBrand [In Stock: true]

Multi-search example completed successfully! 🎉
```

## Best Practices Demonstrated

- **Settings Configuration**: Configure filterable/sortable attributes before searching
- **Task Completion**: Wait for index setup before executing searches
- **Error Handling**: Comprehensive error handling for multi-search operations
- **Resource Management**: Proper client cleanup with `defer client.Close()`
- **Filter Syntax**: Proper filter expression formatting
- **Performance Optimization**: Batch multiple searches for efficiency
