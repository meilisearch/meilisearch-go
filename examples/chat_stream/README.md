## Available Examples

### 📝 [Basic Search](../search)
Basic search functionality with index creation, document management, and various search operations.

### 💬 Chat Streaming
Interactive chat experience with streaming responses and knowledge base integration.

2. **Enterprise Setup**: Configure Meilisearch with LLM integration
3. **Workspace Configuration**: Set up appropriate chat workspaces

## Configuration

The example supports configuration via environment variables:

```bash
# Set Meilisearch server URL (optional, defaults to http://localhost:7700)
export MEILI_HOST="http://localhost:7700"

# Set API key (optional, but recommended for production)
export MEILI_API_KEY="your-api-key"
```

## Configuration

## Running the Example

```bash
# Set environment variables (optional)
export MEILI_HOST="http://localhost:7700"
export MEILI_API_KEY="your-api-key"

# Run the example
go run ./examples/chat_stream
```