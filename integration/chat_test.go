package integration

import (
	"crypto/tls"
	"testing"

	"github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/require"
)

func Test_ChatWorkspaceSettings(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	tests := []struct {
		name   string
		client meilisearch.ServiceManager
	}{
		{
			name:   "TestChatWorkspaceSettings",
			client: sv,
		},
		{
			name:   "TestChatWorkspaceSettingsWithCustomClient",
			client: customSv,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			updateResp, err := tt.client.UpdateWorkspaceSettings(
				"workspace_uid",
				&meilisearch.ChatWorkspaceSettings{
					Source:       "azureOpenAi",
					OrgID:        "your-azure-org-id",
					APIVersion:   "your-api-version",
					DeploymentID: "your-deployment-id",
					APIKey:       "OPEN_AI_API_KEY",
					BaseURL:      "https://your-resource.openai.azure.com",
					Prompts: &meilisearch.Prompts{
						System: "You are a helpful customer support assistant.",
					},
				},
			)
			require.NoError(t, err)

			getResp, err := tt.client.GetWorkspaceSettings("workspace_uid")
			require.NoError(t, err)

			require.Equal(t, updateResp, getResp)
		})
	}
}
