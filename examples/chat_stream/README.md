# Chat Streaming Example

This example demonstrates interactive chat streaming using the Meilisearch Go SDK with enterprise chat capabilities:

- **Interactive chat sessions** with streaming responses  
- **Workspace discovery** and automatic workspace selection
- **Real-time streaming** with proper EOF handling and content parsing
- **Environment-based configuration** for flexible deployment
- **Graceful error handling** and comprehensive resource cleanup

## Interactive Chat Features

### **1. Connection Testing**
- Tests connection to Meilisearch server with health check
- Displays server health status for verification

### **2. Knowledge Base Setup**
- Initializes basic knowledge base (placeholder for actual implementation)
- Sets up foundation for chat functionality

### **3. Workspace Discovery**
- **Lists available chat workspaces** using `ListChatWorkspaces()`
- **Auto-selects workspace**: Uses first available workspace or defaults to "default"
- **Displays workspace information**: Shows all discovered workspace UIDs

### **4. Interactive Chat Session**
- **REPL Interface**: Read-Eval-Print Loop for continuous interaction
- **User Input**: Accepts natural language questions
- **Exit Commands**: Type 'quit' or 'exit' to end session
- **Streaming Responses**: Real-time streaming of AI assistant responses

## Chat Streaming Implementation

### **ChatCompletionQuery Configuration**
```go
chatQuery := &meilisearch.ChatCompletionQuery{
    Model: "gpt-3.5-turbo",
    Messages: []*meilisearch.ChatCompletionMessage{
        {
            Role:    "system",
            Content: "You are a helpful assistant that answers questions about Meilisearch...",
        },
        {
            Role:    "user", 
            Content: query,
        },
    },
    Stream: true,
}
```

### **Stream Processing**
- **Context with timeout**: 60-second timeout for chat operations
- **EOF handling**: Proper `errors.Is(err, io.EOF)` detection
- **Content parsing**: Extracts content from `chunk.Choices[0].Delta.Content`
- **Pointer handling**: Safely dereferences content pointers
- **Real-time display**: Prints streaming content as it arrives

## Prerequisites

1. **Meilisearch Enterprise**: Chat functionality requires Meilisearch Enterprise
2. **LLM Integration**: Configure Meilisearch with LLM integration  
3. **Chat Workspaces**: Set up appropriate chat workspaces
4. **API Key**: Recommended for production deployments

## Configuration

```bash
# Set Meilisearch server URL (defaults to http://localhost:7700)
export MEILI_HOST="http://localhost:7700"

# Set API key (required for Enterprise features)
export MEILI_API_KEY="your-enterprise-api-key"
```

## Running the Example

```bash
# Set environment variables
export MEILI_HOST="http://localhost:7700"
export MEILI_API_KEY="your-enterprise-api-key"

# Run the interactive chat
go run ./examples/chat_stream
```

### **Sample Session**
```
Testing connection to Meilisearch...
✅ Connected to Meilisearch (status: available)

Listing chat workspaces...
Found 2 chat workspaces
  - Workspace: knowledge-base
  - Workspace: support-docs

--- Interactive Chat Session ---
Type your questions (or 'quit' to exit):

You: What is Meilisearch?
Assistant: Meilisearch is a powerful, fast, and hyper-relevant search engine...

You: How do I create an index?
Assistant: To create an index in Meilisearch, you can use the CreateIndex method...

You: quit
Bye!
```

## Implementation Details

### **Stream Management**
- **Proper cleanup**: `defer stream.Close()` ensures resource cleanup
- **Context cancellation**: `defer cancel()` prevents resource leaks
- **Error handling**: Distinguishes between EOF and actual errors

### **Content Processing**
- **Null checks**: Validates `chunk != nil && len(chunk.Choices) > 0`
- **Content validation**: Checks for non-empty content before printing
- **Pointer handling**: Safely dereferences `*chunk.Choices[0].Delta.Content`

### **Interactive Features**
- **Health Check**: Verifies server availability with `client.Health()`
- **Workspace Listing**: Discovers available workspaces automatically
- **REPL Loop**: Continuous interaction until 'quit' or 'exit'
- **Graceful Exit**: Clean shutdown with proper resource cleanup

```bash
# Expected Output
Testing connection to Meilisearch...
✅ Connected to Meilisearch (status: available)

Listing chat workspaces...
Found 1 chat workspaces
  - Workspace: default

--- Interactive Chat Session ---
Type your questions (or 'quit' to exit):

You: [Your question here]
Assistant: [Streaming response appears here in real-time]

You: quit
Bye!
```

## Best Practices Demonstrated

- **Environment Configuration**: Flexible host and API key setup
- **Resource Management**: Proper client and stream cleanup
- **Timeout Handling**: Appropriate timeouts for chat operations
- **Error Handling**: Graceful error handling with informative messages
- **Interactive Design**: User-friendly REPL interface
- **Production Ready**: Enterprise-grade chat implementation

## Advanced Configuration

```bash
# Optional advanced configuration
export MEILI_HOST="http://localhost:7700"
export MEILI_API_KEY="your-enterprise-api-key"

go run ./examples/chat_stream
```

The example will start an interactive chat session where you can ask questions and receive streaming responses. Type 'quit' or 'exit' to end the session.
