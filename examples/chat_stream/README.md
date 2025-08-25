# Meilisearch Go SDK Examples

This directory contains runnable examples demonstrating how to use the Meilisearch Go SDK for common and advanced operations.

## Available Examples

### 📝 [Basic Search](./search)
Demonstrates core search functionality including:
- Client initialization and connection testing
- Index creation and management
- Document addition and indexing
- Various search operations (basic, filtered, faceted)
- Task management and waiting for operations

### 💬 [Chat Streaming](./chat_stream)
Shows streaming capabilities and conversational search:
- Chat workspace management
- Real-time streaming responses
- Interactive chat sessions
- Knowledge base integration
- Context-aware responses

## Running Examples

Each example is self-contained and can be run with:

```bash
go run ./examples/<example-name>
```

For example:
```bash
go run ./examples/search
go run ./examples/chat_stream
```

## Prerequisites

- Go 1.20 or higher
- Meilisearch server running (default: http://localhost:7700)
- Valid API key (if authentication is enabled)

## Configuration

Before running the examples, make sure to:
1. Update the Meilisearch server URL in each example
2. Set the appropriate API key
3. Ensure your Meilisearch instance has the required features enabled

## Code Style

All examples follow idiomatic Go style and include:
- Comprehensive error handling
- Clear documentation and comments
- Realistic use cases
- Best practices for production usage

## Contributing

When adding new examples:
- Follow the established directory structure
- Include a detailed README for each example
- Ensure examples are self-contained and runnable
- Add comprehensive comments explaining each step
- Handle errors appropriately
- Test examples against a real Meilisearch instance
