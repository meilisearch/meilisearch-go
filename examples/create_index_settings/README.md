# Index Creation and Settings Example

This example demonstrates comprehensive index creation and advanced settings configuration using the Meilisearch Go SDK:

- **Index creation** with primary key specification
- **Advanced settings** configuration for optimal search
- **Ranking rules** customization
- **Attribute configuration** for search, filtering, and sorting
- **Typo tolerance** and **synonym** management

## Features Demonstrated

1. **Basic Index Creation**: Create index with primary key
2. **Ranking Rules**: Custom search relevance configuration
3. **Attribute Management**:
   - Searchable attributes (search priority)
   - Displayed attributes (result fields)
   - Filterable attributes (filter capabilities)
   - Sortable attributes (sorting options)
4. **Search Enhancement**:
   - Stop words configuration
   - Synonyms for better matching
   - Typo tolerance settings
5. **Performance Settings**:
   - Search cutoff time
   - Pagination limits
   - Faceting configuration
6. **Index Information**: Retrieve index metadata

## Settings Configured

### Ranking Rules
- **words**: Matches more words first
- **typo**: Fewer typos first
- **proximity**: Words closer together first
- **attribute**: Matches in more important attributes first
- **sort**: Custom sorting rules
- **exactness**: Exact matches first
- **desc(rating)**: Custom rule for rating-based ranking

### Attributes
- **Searchable**: title, content, author, tags
- **Filterable**: category, author, featured, rating, publish_date, word_count
- **Sortable**: publish_date, rating, word_count
- **Displayed**: id, title, author, category, publish_date, rating, featured

### Search Enhancement
- **Stop Words**: Common words ignored in search
- **Synonyms**: Related terms for better matching
- **Typo Tolerance**: Configurable typo handling
- **Distinct Results**: Avoid duplicates by author

## Configuration

```bash
# Set Meilisearch server URL (optional, defaults to http://localhost:7700)
export MEILI_HOST="http://localhost:7700"

# Set API key (optional, but recommended for production)
export MEILI_API_KEY="your-api-key"
```

## Running the Example

```bash
go run ./examples/create_index_settings
```

The example will create an "articles" index, configure comprehensive settings, add sample articles, and demonstrate how the settings improve search results.
