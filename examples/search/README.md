- Meilisearch server running (default: `http://localhost:7700`)
- Valid API key (if authentication is enabled)

## Configuration

The example supports configuration via environment variables:

```bash
# Set Meilisearch server URL (optional, defaults to http://localhost:7700)
export MEILI_HOST="http://localhost:7700"

# Set API key (optional, but recommended for production)
export MEILI_API_KEY="your-api-key"
```

## Running the Example

1. **Start the Meilisearch server:**
   ```bash
   ./meilisearch
   ```

2. **Set environment variables (optional):**
   ```bash
   export MEILI_HOST="http://localhost:7700"
   export MEILI_API_KEY="your-api-key"  # If authentication is enabled
   ```

3. **Run the example:**
   ```bash
   go run ./examples/search
