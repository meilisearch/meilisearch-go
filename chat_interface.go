package meilisearch

import "context"

type ChatManager interface {
	ChatReader

	// UpdateChatWorkspace updates a chat workspace by its ID.
	UpdateChatWorkspace(ctx context.Context, uid string, settings *ChatWorkspaceSettings) (*ChatWorkspaceSettings, error)

	// ResetChatWorkspace resets a chat workspace by its ID.
	ResetChatWorkspace(ctx context.Context, uid string) (*ChatWorkspaceSettings, error)
}

type ChatReader interface {
	// ChatCompletionStream retrieves a stream of chat completions for a given workspace and query.
	ChatCompletionStream(ctx context.Context, workspace string, query *ChatCompletionQuery) (*Stream[*ChatCompletionStreamChunk], error)

	// ListChatWorkspaces retrieves all chat workspaces.
	ListChatWorkspaces(ctx context.Context, query *ListChatWorkSpaceQuery) (*ListChatWorkspace, error)

	// GetChatWorkspace retrieves a chat workspace by its ID.
	GetChatWorkspace(ctx context.Context, uid string) (*ChatWorkspace, error)

	// GetChatWorkspaceSettings retrieves chat workspace settings by its ID.
	GetChatWorkspaceSettings(ctx context.Context, uid string) (*ChatWorkspaceSettings, error)
}
