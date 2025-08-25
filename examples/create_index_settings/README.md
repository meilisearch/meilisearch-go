# Index Creation and Settings Example

This example demonstrates comprehensive index creation and advanced settings configuration using the Meilisearch Go SDK with article data:

- **Index creation** with primary key specification and validation
- **Advanced settings configuration** for optimal search performance
- **Custom ranking rules** for relevance fine-tuning
- **Comprehensive attribute configuration** for search, filtering, and sorting
- **Typo tolerance settings** with configurable thresholds
- **Stop words and synonyms** for enhanced search quality
- **Task monitoring** with proper completion waiting

## Article Database Structure

The example uses an `Article` struct representing article documents:

```go
type Article struct {
    ID          string   `json:"id"`           // Primary key (unique identifier)
    Title       string   `json:"title"`        // Article title
    Content     string   `json:"content"`      // Full article content
    Author      string   `json:"author"`       // Author name
    Category    string   `json:"category"`     // Article category
    Tags        []string `json:"tags"`         // Search tags
    PublishDate string   `json:"publish_date"` // Publication date (ISO format)
    Rating      float64  `json:"rating"`       // Article rating (1.0-5.0)
    WordCount   int      `json:"word_count"`   // Article length
    Featured    bool     `json:"featured"`     // Featured article flag
}
```

## Sample Article Data

The example includes 3 diverse articles:
- **"Getting Started with Meilisearch"** - Technology, John Doe, Rating: 4.5, 750 words, Featured
- **"Advanced Search Techniques"** - Technology, Jane Smith, Rating: 4.8, 1200 words, Not Featured  
- **"Building Modern Applications"** - Development, Bob Johnson, Rating: 4.2, 950 words, Featured

## Advanced Settings Configuration

### **1. Index Creation**
- Creates "articles" index with "id" as primary key
- Validates index creation with task completion waiting
- Displays success confirmation with detailed logging

### **2. Custom Ranking Rules**
```go
RankingRules: []string{
    "words",        // Number of matching words
    "typo",         // Fewer typos ranked higher
    "proximity",    // Word proximity importance
    "attribute",    // Attribute-based ranking
    "sort",         // Custom sorting rules
    "exactness",    // Exact matches prioritized
    "desc(rating)", // Custom: Higher ratings first
}
```

#### **Ranking Rule Explanation:**
- **words**: Articles matching more query words rank higher
- **typo**: Articles with fewer typos in matches rank higher  
- **proximity**: Words closer together in content rank higher
- **attribute**: Matches in more important attributes (title > content) rank higher
- **sort**: Enables custom sorting capabilities
- **exactness**: Exact word matches rank higher than partial matches
- **desc(rating)**: Custom rule prioritizing higher-rated articles

### **3. Attribute Configuration**

#### **Searchable Attributes (Search Priority Order)**
```go
SearchableAttributes: &[]string{"title", "content", "author", "tags"}
```
- **title**: Highest priority - matches in titles rank highest
- **content**: Main content searchability
- **author**: Author name searching
- **tags**: Tag-based search capabilities

#### **Displayed Attributes (Result Fields)**
```go
DisplayedAttributes: &[]string{"id", "title", "author", "category", "publish_date", "rating", "featured"}
```
- Controls which fields appear in search results
- Excludes sensitive or unnecessary fields (e.g., full content, word_count)

#### **Filterable Attributes (Filter Capabilities)**
```go
FilterableAttributes: []string{"category", "author", "featured", "rating", "publish_date", "word_count"}
```
- **category**: Filter by article category (Technology, Development, etc.)
- **author**: Filter by specific authors
- **featured**: Filter featured vs regular articles
- **rating**: Numeric range filtering (rating > 4.0)
- **publish_date**: Date-based filtering
- **word_count**: Length-based filtering

#### **Sortable Attributes (Sorting Options)**
```go
SortableAttributes: []string{"publish_date", "rating", "word_count"}
```
- **publish_date**: Sort by publication date (newest/oldest first)
- **rating**: Sort by article rating (highest/lowest first) 
- **word_count**: Sort by article length (longest/shortest first)

### **4. Search Enhancement Features**

#### **Stop Words Configuration**
```go
StopWords: &[]string{"the", "a", "an", "and", "or", "but", "in", "on", "at", "to", "for", "of", "with", "by"}
```
- Common words ignored during search to improve relevance
- Reduces noise from frequent, low-meaning words

#### **Synonyms Management**
```go
Synonyms: &map[string][]string{
    "js":     {"javascript"},
    "react":  {"reactjs"},
    "vue":    {"vuejs"},
    "search": {"find", "lookup"},
}
```
- Expands search queries with related terms
- Improves search recall for technical terms and abbreviations

#### **Typo Tolerance Settings**
```go
TypoTolerance: &meilisearch.TypoTolerance{
    Enabled: true,
    MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
        OneTypo:  5,  // Words ≥5 chars allow 1 typo
        TwoTypos: 9,  // Words ≥9 chars allow 2 typos
    },
}
```
- **OneTypo (5)**: Words with 5+ characters can have 1 typo and still match
- **TwoTypos (9)**: Words with 9+ characters can have 2 typos and still match
- Balances search flexibility with accuracy

### **5. Advanced Features**

#### **Distinct Attribute**
```go
DistinctAttribute: stringPtr("author")
```
- Ensures only one result per author in search results
- Prevents author over-representation in results

#### **Proximity Precision**
```go
ProximityPrecision: stringPtr("byWord")
```
- Configures how word proximity is calculated
- "byWord" provides precise word-level proximity scoring

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

# Run the create index settings example
go run ./examples/create_index_settings
```

### **Expected Output**
```
Testing connection to Meilisearch...
✅ Connected to Meilisearch (status: available)

1. Creating index 'articles' with primary key 'id'...
✅ Index created successfully! Task #1 completed

2. Adding sample articles...
✅ Added 3 articles successfully! Task #2 completed

3. Configuring advanced settings...
✅ Settings updated successfully! Task #3 completed

Index and settings configuration completed successfully! 🎉

The example will create an "articles" index, configure comprehensive settings, add sample articles, and demonstrate how the settings improve search results.
