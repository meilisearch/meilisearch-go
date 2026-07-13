package meilisearch

import "time"

type ChatWorkspace struct {
	UID string `json:"uid"`
}

type ChatWorkspaceSettings struct {
	Source       ChatSource                    `json:"source"`
	OrgId        string                        `json:"orgId"`
	ProjectId    string                        `json:"projectId"`
	ApiVersion   string                        `json:"apiVersion"`
	DeploymentId string                        `json:"deploymentId"`
	BaseUrl      string                        `json:"baseUrl"`
	ApiKey       string                        `json:"apiKey,omitempty"`
	Prompts      *ChatWorkspaceSettingsPrompts `json:"prompts"`
}

type ChatWorkspaceSettingsPrompts struct {
	System              string `json:"system"`
	SearchDescription   string `json:"searchDescription"`
	SearchQParam        string `json:"searchQParam"`
	SearchFilterParam   string `json:"searchFilterParam"`
	SearchIndexUidParam string `json:"searchIndexUidParam"`
}

type ListChatWorkspace struct {
	Results []*ChatWorkspace `json:"results"`
	Offset  int64            `json:"offset"`
	Limit   int64            `json:"limit"`
	Total   int64            `json:"total"`
}

type ListChatWorkSpaceQuery struct {
	Limit  int64
	Offset int64
}

type ChatCompletionQuery struct {
	Model    string                   `json:"model"`
	Messages []*ChatCompletionMessage `json:"messages"`
	Stream   bool                     `json:"stream"`
}

type ChatCompletionMessage struct {
	Role    ChatRole `json:"role"`
	Content string   `json:"content"`
}

type ChatCompletionStreamChunk struct {
	ID                string                  `json:"id"`
	Object            *string                 `json:"object,omitempty"`
	Created           Timestampz              `json:"created,omitempty"`
	Model             string                  `json:"model,omitempty"`
	Choices           []*ChatCompletionChoice `json:"choices"`
	ServiceTier       *string                 `json:"service_tier,omitempty"`
	SystemFingerprint *string                 `json:"system_fingerprint,omitempty"`
	Usage             any                     `json:"usage,omitempty"`
}

type ChatCompletionChoice struct {
	Index        int64                      `json:"index"`
	Delta        *ChatCompletionChoiceDelta `json:"delta"`
	FinishReason *string                    `json:"finish_reason,omitempty"`
	Logprobs     any                        `json:"logprobs"`
}

type ChatCompletionChoiceDelta struct {
	Content      *string   `json:"content,omitempty"`
	Role         *ChatRole `json:"role,omitempty"`
	Refusal      *string   `json:"refusal,omitempty"`
	FunctionCall *string   `json:"function_call,omitempty"`
	ToolCalls    *string   `json:"tool_calls,omitempty"`
}

type Timestampz int64

func (t Timestampz) String() string {
	return time.Unix(int64(t), 0).UTC().Format(time.RFC3339)
}

func (t Timestampz) ToTime() time.Time {
	return time.Unix(int64(t), 0).UTC()
}
