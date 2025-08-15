package meilisearch

import "context"

type ChatManager interface {
	ChatReader

	// UpdateChatWorkspace updates a chat workspace by its ID.
	UpdateChatWorkspace(uid string, settings *ChatWorkspaceSettings) (*ChatWorkspaceSettings, error)
	// UpdateChatWorkspaceWithContext updates a chat workspace by its ID with a context.
	UpdateChatWorkspaceWithContext(ctx context.Context, uid string, settings *ChatWorkspaceSettings) (*ChatWorkspaceSettings, error)

	// ResetChatWorkspace resets a chat workspace by its ID.
	ResetChatWorkspace(uid string) (*ChatWorkspaceSettings, error)
	// ResetChatWorkspaceWithContext resets a chat workspace by its ID with a context.
	ResetChatWorkspaceWithContext(ctx context.Context, uid string) (*ChatWorkspaceSettings, error)
}

type ChatReader interface {
	// ChatCompletionStream retrieves a stream of chat completions for a given workspace and query.
	ChatCompletionStream(workspace string, query *ChatCompletionQuery) (*Stream[*ChatCompletionStreamChunk], error)
	// ChatCompletionStreamWithContext retrieves a stream of chat completions for a given workspace and query with a context.
	ChatCompletionStreamWithContext(ctx context.Context, workspace string, query *ChatCompletionQuery) (*Stream[*ChatCompletionStreamChunk], error)

	// ListChatWorkspaces retrieves all chat workspaces.
	ListChatWorkspaces(query *ListChatWorkSpaceQuery) (*ListChatWorkspace, error)
	// ListChatWorkspacesWithContext retrieves all chat workspaces with a context.
	ListChatWorkspacesWithContext(ctx context.Context, query *ListChatWorkSpaceQuery) (*ListChatWorkspace, error)

	// GetChatWorkspace retrieves a chat workspace by its ID.
	GetChatWorkspace(uid string) (*ChatWorkspace, error)
	// GetChatWorkspaceWithContext retrieves a chat workspace by its ID with a context.
	GetChatWorkspaceWithContext(ctx context.Context, uid string) (*ChatWorkspace, error)

	// GetChatWorkspaceSettings retrieves chat workspace settings by its ID.
	GetChatWorkspaceSettings(uid string) (*ChatWorkspaceSettings, error)
	// GetChatWorkspaceSettingsWithContext retrieves chat workspace settings by its ID with a context.
	GetChatWorkspaceSettingsWithContext(ctx context.Context, uid string) (*ChatWorkspaceSettings, error)
}
