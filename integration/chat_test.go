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
		MinVersion:         tls.VersionTLS12,
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

			workspaceUID := "test-workspace-" + tt.name
			apiKey := "test-api-key-placeholder"

			updateResp, err := tt.client.UpdateWorkspaceSettings(
				workspaceUID,
				&meilisearch.ChatWorkspaceSettings{
					Source:       "azureOpenAi",
					OrgID:        "your-azure-org-id",
					APIVersion:   "your-api-version",
					DeploymentID: "your-deployment-id",
					APIKey:       apiKey,
					BaseURL:      "https://your-resource.openai.azure.com",
					Prompts: &meilisearch.Prompts{
						System: "You are a helpful customer support assistant.",
					},
				},
			)
			require.NoError(t, err)

			getResp, err := tt.client.GetWorkspaceSettings(workspaceUID)
			require.NoError(t, err)

			require.Equal(t, updateResp, getResp)
		})
	}
}
