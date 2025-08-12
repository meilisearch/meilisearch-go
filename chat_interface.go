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
	// ListChatWorkspaces retrieves all chat workspaces.
	ListChatWorkspaces(query *ListChatWorkSpaceQuery) (*ListChatWorkspace, error)
	// ListChatWorkspacesWithContext retrieves all chat workspaces with a context.
	ListChatWorkspacesWithContext(ctx context.Context, query *ListChatWorkSpaceQuery) (*ListChatWorkspace, error)

	// GetChatWorkspace retrieves a chat workspace by its ID.
	GetChatWorkspace(uid string) (*ChatWorkspace, error)
	// GetChatWorkSpaceWithContext retrieves a chat workspace by its ID with a context.
	GetChatWorkSpaceWithContext(ctx context.Context, uid string) (*ChatWorkspace, error)

	// GetChatWorkspaceSettings retrieves chat workspace settings by its ID.
	GetChatWorkspaceSettings(uid string) (*ChatWorkspaceSettings, error)
	// GetChatWorkspaceSettingsWithContext retrieves chat workspace settings by its ID with a context.
	GetChatWorkspaceSettingsWithContext(ctx context.Context, uid string) (*ChatWorkspaceSettings, error)
}
