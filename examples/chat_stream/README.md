# Chat Streaming Example

This example shows how to use Meilisearch's chat features with streaming responses.

## What it does

- Check server health
- List available chat workspaces
- Start an interactive chat session
- Stream responses in real-time
- Type 'quit' to exit

## Prerequisites

 **Meilisearch Enterprise** with chat features enabled

## Configuration

```bash
export MEILI_HOST="http://localhost:7700"
export MEILI_API_KEY="your-enterprise-api-key"
```

## Run it

```bash
go run ./examples/chat_stream
```
```

The example will start an interactive chat session where you can ask questions and receive streaming responses. Type 'quit' or 'exit' to end the session.
