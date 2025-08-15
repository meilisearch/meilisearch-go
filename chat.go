package meilisearch

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

func (m *meilisearch) ChatCompletionStream(workspace string, query *ChatCompletionQuery) (*Stream[*ChatCompletionStreamChunk], error) {
	return m.ChatCompletionStreamWithContext(context.Background(), workspace, query)
}

func (m *meilisearch) ChatCompletionStreamWithContext(ctx context.Context, workspace string, query *ChatCompletionQuery) (*Stream[*ChatCompletionStreamChunk], error) {
	if query == nil {
		return nil, fmt.Errorf("query cannot be nil")
	}

	if !query.Stream {
		query.Stream = true
	}

	req := &internalRequest{
		endpoint:            fmt.Sprintf("/chats/%s/chat/completions", workspace),
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         query,
		acceptedStatusCodes: []int{http.StatusOK},
	}

	internalError := &Error{
		Endpoint:         req.endpoint,
		Method:           req.method,
		Function:         req.functionName,
		RequestToString:  "empty request",
		ResponseToString: "empty response",
		MeilisearchApiError: meilisearchApiError{
			Message: "empty meilisearch message",
		},
		StatusCodeExpected: req.acceptedStatusCodes,
		encoder:            m.client.encoder,
	}

	resp, err := m.client.sendRequest(ctx, req, internalError)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	dec := NewDecoder(resp)
	if dec == nil {
		// Defensive: if something went wrong constructing the decoder,
		// close the response to avoid leaking the connection.
		_ = resp.Body.Close()
		return nil, fmt.Errorf("failed to create stream decoder: nil decoder")
	}

	return NewStream[*ChatCompletionStreamChunk](dec, m.client.jsonUnmarshal), nil
}

func (m *meilisearch) GetChatWorkspace(id string) (*ChatWorkspace, error) {
	return m.GetChatWorkspaceWithContext(context.Background(), id)
}

func (m *meilisearch) GetChatWorkspaceWithContext(ctx context.Context, uid string) (*ChatWorkspace, error) {
	resp := new(ChatWorkspace)
	req := &internalRequest{
		endpoint:            "/chats/" + uid,
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetChatWorkspace",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) ListChatWorkspaces(query *ListChatWorkSpaceQuery) (*ListChatWorkspace, error) {
	return m.ListChatWorkspacesWithContext(context.Background(), query)
}

func (m *meilisearch) ListChatWorkspacesWithContext(ctx context.Context, query *ListChatWorkSpaceQuery) (*ListChatWorkspace, error) {
	resp := new(ListChatWorkspace)
	req := &internalRequest{
		endpoint:            "/chats",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		withQueryParams:     make(map[string]string),
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "ListChatWorkspaces",
	}

	if query != nil {
		req.withQueryParams["limit"] = strconv.FormatInt(query.Limit, 10)
		req.withQueryParams["offset"] = strconv.FormatInt(query.Offset, 10)
	}

	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) GetChatWorkspaceSettings(id string) (*ChatWorkspaceSettings, error) {
	return m.GetChatWorkspaceSettingsWithContext(context.Background(), id)
}

func (m *meilisearch) GetChatWorkspaceSettingsWithContext(ctx context.Context, uid string) (*ChatWorkspaceSettings, error) {
	resp := new(ChatWorkspaceSettings)
	req := &internalRequest{
		endpoint:            "/chats/" + uid + "/settings",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetChatWorkspaceSettings",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) UpdateChatWorkspace(id string, chatWorkspace *ChatWorkspaceSettings) (*ChatWorkspaceSettings, error) {
	return m.UpdateChatWorkspaceWithContext(context.Background(), id, chatWorkspace)
}

func (m *meilisearch) UpdateChatWorkspaceWithContext(ctx context.Context, uid string, settings *ChatWorkspaceSettings) (*ChatWorkspaceSettings, error) {
	resp := new(ChatWorkspaceSettings)
	req := &internalRequest{
		endpoint:            "/chats/" + uid + "/settings",
		method:              http.MethodPatch,
		withRequest:         settings,
		withResponse:        resp,
		contentType:         contentTypeJSON,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "UpdateChatWorkspace",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) ResetChatWorkspace(id string) (*ChatWorkspaceSettings, error) {
	return m.ResetChatWorkspaceWithContext(context.Background(), id)
}

func (m *meilisearch) ResetChatWorkspaceWithContext(ctx context.Context, uid string) (*ChatWorkspaceSettings, error) {
	resp := new(ChatWorkspaceSettings)
	req := &internalRequest{
		endpoint:            "/chats/" + uid + "/settings",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "ResetChatWorkspace",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}
