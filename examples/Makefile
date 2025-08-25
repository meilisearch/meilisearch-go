.PHONY: help search chat_stream all clean

# Default target
help:
	@echo "Available commands:"
	@echo "  make search      - Run the basic search example"
	@echo "  make chat_stream - Run the chat streaming example"
	@echo "  make all         - Run all examples"
	@echo "  make clean       - Clean up any temporary files"
	@echo ""
	@echo "Prerequisites:"
	@echo "  - Meilisearch server running on http://localhost:7700"
	@echo "  - Go 1.20 or higher installed"

# Run the basic search example
search:
	@echo "Running basic search example..."
	go run ./search/main.go

# Run the chat streaming example
chat_stream:
	@echo "Running chat streaming example..."
	go run ./chat_stream/main.go

# Run all examples
all: search chat_stream

# Clean up temporary files
clean:
	@echo "Cleaning up..."
	go clean ./...
