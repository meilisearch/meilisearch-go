package integration

import (
	"os"
	"testing"

	"github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetChatWorkspace(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanupChat(sv))

	resp, err := sv.ExperimentalFeatures().SetChatCompletions(true).Update()
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.ChatCompletions)

	chat := sv.ChatManager()
	require.NotNil(t, chat)

	uid := "test-workspace"

	workspace, err := chat.UpdateChatWorkspace(uid, &meilisearch.ChatWorkspaceSettings{
		Source: meilisearch.OpenaiChatSource,
		ApiKey: "test-api-key",
		Prompts: &meilisearch.ChatWorkspaceSettingsPrompts{
			System: "foobar",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, workspace)

	chatReader := sv.ChatReader()
	require.NotNil(t, chatReader)

	got, err := chatReader.GetChatWorkspace(uid)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, uid, got.UID)
}

func Test_GetChatWorkspaceSettings(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanupChat(sv))

	resp, err := sv.ExperimentalFeatures().SetChatCompletions(true).Update()
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.ChatCompletions)

	chat := sv.ChatManager()
	require.NotNil(t, chat)

	uid := "test-workspace"

	want, err := chat.UpdateChatWorkspace(uid, &meilisearch.ChatWorkspaceSettings{
		Source: meilisearch.OpenaiChatSource,
		ApiKey: "test-api-key",
		Prompts: &meilisearch.ChatWorkspaceSettingsPrompts{
			System: "foobar",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, want)

	got, err := chat.GetChatWorkspaceSettings(uid)
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, want.Source, got.Source)
	assert.Equal(t, want.Prompts, got.Prompts)
}

func Test_ListChatWorkspace(t *testing.T) {
	sv := setup(t, "")

	resp, err := sv.ExperimentalFeatures().SetChatCompletions(true).Update()
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.ChatCompletions)

	chat := sv.ChatManager()
	require.NotNil(t, chat)

	tests := []struct {
		name  string
		query *meilisearch.ListChatWorkSpaceQuery
		reqs  map[string]*meilisearch.ChatWorkspaceSettings
		resp  *meilisearch.ListChatWorkspace
	}{
		{
			name: "TestListChatWorkspaceNormalState",
			query: &meilisearch.ListChatWorkSpaceQuery{
				Limit:  2,
				Offset: 0,
			},
			reqs: map[string]*meilisearch.ChatWorkspaceSettings{
				"test-workspace-1": {
					Source: meilisearch.OpenaiChatSource,
					ApiKey: "test-api-key-1",
					Prompts: &meilisearch.ChatWorkspaceSettingsPrompts{
						System: "system prompt 1",
					},
				},
				"test-workspace-2": {
					Source: meilisearch.OpenaiChatSource,
					ApiKey: "test-api-key-2",
					Prompts: &meilisearch.ChatWorkspaceSettingsPrompts{
						System: "system prompt 2",
					},
				},
			},
			resp: &meilisearch.ListChatWorkspace{
				Results: []*meilisearch.ChatWorkspace{
					{
						UID: "test-workspace-1",
					},
					{
						UID: "test-workspace-2",
					},
				},
				Limit:  2,
				Offset: 0,
				Total:  2,
			},
		},
		{
			name:  "TestListChatWorkspaceWithoutQuery",
			query: nil,
			reqs: map[string]*meilisearch.ChatWorkspaceSettings{
				"test-workspace-1": {
					Source: meilisearch.OpenaiChatSource,
					ApiKey: "test-api-key-1",
					Prompts: &meilisearch.ChatWorkspaceSettingsPrompts{
						System: "system prompt 1",
					},
				},
			},
			resp: &meilisearch.ListChatWorkspace{
				Results: []*meilisearch.ChatWorkspace{
					{
						UID: "test-workspace-1",
					},
				},
				Limit:  20,
				Offset: 0,
				Total:  2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(cleanupChat(sv))

			if tt.reqs != nil {
				for uid, settings := range tt.reqs {
					resp, err := chat.UpdateChatWorkspace(uid, settings)
					require.NoError(t, err)
					require.NotNil(t, resp)
				}
			}

			listResp, err := chat.ListChatWorkspaces(tt.query)
			require.NoError(t, err)
			require.NotNil(t, listResp)

			assert.NotEmpty(t, listResp.Results)
			assert.Equal(t, tt.resp.Limit, listResp.Limit)
			assert.Equal(t, tt.resp.Offset, listResp.Offset)
		})
	}
}

func Test_UpdateChatWorkspaceSettings(t *testing.T) {
	sv := setup(t, "")

	resp, err := sv.ExperimentalFeatures().SetChatCompletions(true).Update()
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.ChatCompletions)

	chat := sv.ChatManager()
	require.NotNil(t, chat)

	tests := []struct {
		name string
		uid  string
		req  *meilisearch.ChatWorkspaceSettings
	}{
		{
			name: "TestUpdateChatWorkspaceSettingsNormalState",
			uid:  "test-workspace",
			req: &meilisearch.ChatWorkspaceSettings{
				Source: meilisearch.OpenaiChatSource,
				ApiKey: "test-api-key",
				Prompts: &meilisearch.ChatWorkspaceSettingsPrompts{
					System: "system prompt",
				},
			},
		},
		{
			name: "TestUpdateChatWorkspaceSettingsFull",
			uid:  "test-workspace",
			req: &meilisearch.ChatWorkspaceSettings{
				Source:       meilisearch.MistralChatSource,
				ApiKey:       "test-api-key",
				OrgId:        "test-org",
				ApiVersion:   "v1",
				BaseUrl:      "https://api.mistral.ai",
				ProjectId:    "test-project",
				DeploymentId: "test-deployment",
				Prompts: &meilisearch.ChatWorkspaceSettingsPrompts{
					System: "system prompt",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(cleanupChat(sv))

			resp, err := chat.UpdateChatWorkspace(tt.uid, tt.req)
			require.NoError(t, err)
			require.NotNil(t, resp)

			got, err := chat.GetChatWorkspaceSettings(tt.uid)
			require.NoError(t, err)
			require.NotNil(t, got)

			assert.Equal(t, resp, got)
			assert.Equal(t, resp.Source, got.Source)
			assert.Equal(t, resp.OrgId, got.OrgId)
			assert.Equal(t, resp.ProjectId, got.ProjectId)
			assert.Equal(t, resp.ApiVersion, got.ApiVersion)
			assert.Equal(t, resp.DeploymentId, got.DeploymentId)
			assert.Equal(t, resp.BaseUrl, got.BaseUrl)
			assert.Equal(t, resp.Prompts, got.Prompts)
		})
	}
}

func Test_ResetChatWorkspaceSettings(t *testing.T) {
	sv := setup(t, "")

	resp, err := sv.ExperimentalFeatures().SetChatCompletions(true).Update()
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.ChatCompletions)

	chat := sv.ChatManager()
	require.NotNil(t, chat)

	uid := "test-workspace"

	want, err := chat.UpdateChatWorkspace(uid, &meilisearch.ChatWorkspaceSettings{
		Source: meilisearch.OpenaiChatSource,
		ApiKey: "test-api-key",
		Prompts: &meilisearch.ChatWorkspaceSettingsPrompts{
			System: "system prompt",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, want)

	got, err := chat.ResetChatWorkspace(uid)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, meilisearch.OpenaiChatSource, got.Source)
}

func Test_ChatCompletionStream(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanupChat(sv))

	resp, err := sv.ExperimentalFeatures().SetChatCompletions(true).Update()
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.ChatCompletions)

	chat := sv.ChatManager()
	require.NotNil(t, chat)

	uid := "test-workspace"
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("OPENAI_API_KEY not set; skipping live chat completion stream integration test")
	}

	_, err = chat.UpdateChatWorkspace(uid, &meilisearch.ChatWorkspaceSettings{
		Source: meilisearch.OpenaiChatSource,
		ApiKey: apiKey,
		Prompts: &meilisearch.ChatWorkspaceSettingsPrompts{
			System: "You are a helpful assistant that answers questions based on the provided context.",
		},
	})
	require.NoError(t, err)

	query := &meilisearch.ChatCompletionQuery{
		Model: "gpt-4o-mini",
		Messages: []*meilisearch.ChatCompletionMessage{
			{
				Role:    meilisearch.ChatRoleUser,
				Content: "Hello, how are you?",
			},
		},
		Stream: true,
	}

	stream, err := chat.ChatCompletionStream(uid, query)
	require.NoError(t, err)
	require.NotNil(t, stream)
	defer func() { _ = stream.Close() }()

	content := ""

	for stream.Next() {
		require.NoError(t, stream.Err())
		chunk := stream.Current()
		require.NotNil(t, chunk)
		assert.NotEmpty(t, chunk.Choices)
		for _, choice := range chunk.Choices {
			require.NotNil(t, choice.Delta)
			if choice.Delta.Content != nil {
				content += *choice.Delta.Content
			}
		}
	}

	assert.NotEmpty(t, content)
	t.Log(content)
}
