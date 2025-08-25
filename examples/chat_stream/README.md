# Chat Streaming Example

This example demonstrates interactive chat streaming using the Meilisearch Go SDK with enterprise chat capabilities:

- **Interactive chat sessions** with streaming responses
- **Workspace management** and automatic workspace selection
- **Real-time streaming** with proper EOF handling
- **Environment-based configuration** for flexible deployment
- **Graceful error handling** and resource cleanup

## Available Examples

### 📝 [Basic Search](../search)
Basic search functionality with index creation, document management, and various search operations.

### 💬 Chat Streaming
Interactive chat experience with streaming responses and knowledge base integration.

## Prerequisites

1. **Meilisearch Enterprise**: Chat functionality requires Meilisearch Enterprise
2. **LLM Integration**: Configure Meilisearch with LLM integration
3. **Chat Workspaces**: Set up appropriate chat workspaces

## Configuration

The example supports configuration via environment variables:

```bash
# Set Meilisearch server URL (optional, defaults to http://localhost:7700)
export MEILI_HOST="http://localhost:7700"

# Set API key (optional, but recommended for production)
export MEILI_API_KEY="your-api-key"
```

## Features Demonstrated

1. **Environment Configuration**: Uses `MEILI_HOST` and `MEILI_API_KEY`
2. **Workspace Discovery**: Automatically lists and selects available workspaces
3. **Interactive Loop**: Simple REPL interface for chat interaction
4. **Streaming Chat**: Real-time streaming responses with proper EOF handling
5. **Resource Management**: Proper client and stream cleanup
6. **Error Handling**: Graceful error handling with informative messages

## Running the Example

```bash
# Set environment variables (optional)
export MEILI_HOST="http://localhost:7700"
export MEILI_API_KEY="your-api-key"

# Run the example
go run ./examples/chat_stream
```

The example will start an interactive chat session where you can ask questions and receive streaming responses. Type 'quit' or 'exit' to end the session.
