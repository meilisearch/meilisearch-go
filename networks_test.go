package meilisearch

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateNetwork_ResponseType(t *testing.T) {
	t.Parallel()

	// An update is enqueued as a task whenever the network has a leader once updated, even if the request does not set the leader.
	const (
		taskBody    = `{"taskUid":1,"indexUid":null,"status":"enqueued","type":"networkTopologyChange","enqueuedAt":"2026-11-08T15:32:27.357211593Z"}`
		networkBody = `{"self":"ms-00","remotes":{},"shards":{},"leader":null,"version":"01a1bbca-0c21-7251-89fa-29c0a945af83"}`
	)

	tests := []struct {
		name        string
		request     *UpdateNetworkRequest
		statusCode  int
		body        string
		wantTask    bool
		wantTaskUID int64
	}{
		{
			name:        "leader set by the request",
			request:     &UpdateNetworkRequest{Leader: String("ms-00")},
			statusCode:  http.StatusAccepted,
			body:        taskBody,
			wantTask:    true,
			wantTaskUID: 1,
		},
		{
			// The first task of an instance has uid 0, so the response type cannot be detected from the uid.
			name:       "first task of the instance",
			request:    &UpdateNetworkRequest{Leader: String("ms-00")},
			statusCode: http.StatusAccepted,
			body:       `{"taskUid":0,"indexUid":null,"status":"enqueued","type":"networkTopologyChange","enqueuedAt":"2026-11-11T12:42:09.559137459Z"}`,
			wantTask:   true,
		},
		{
			name:        "self updated while a leader is configured",
			request:     &UpdateNetworkRequest{Self: String("ms-00")},
			statusCode:  http.StatusAccepted,
			body:        taskBody,
			wantTask:    true,
			wantTaskUID: 1,
		},
		{
			name:       "self updated without leader",
			request:    &UpdateNetworkRequest{Self: String("ms-00")},
			statusCode: http.StatusOK,
			body:       networkBody,
		},
		{
			name:       "leader removed",
			request:    &UpdateNetworkRequest{Leader: Null[string]()},
			statusCode: http.StatusOK,
			body:       networkBody,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodPatch, r.Method)
				require.Equal(t, "/network", r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			got, err := New(srv.URL).UpdateNetwork(context.Background(), tt.request)
			require.NoError(t, err)

			if tt.wantTask {
				task, ok := got.(*Task)
				require.True(t, ok, "got %T, want *Task", got)
				require.Equal(t, tt.wantTaskUID, task.TaskUID)
				require.Equal(t, TaskStatusEnqueued, task.Status)
				require.Equal(t, TaskTypeNetworkTopologyChange, task.Type)
			} else {
				network, ok := got.(*Network)
				require.True(t, ok, "got %T, want *Network", got)
				require.Equal(t, "ms-00", network.Self)
				require.Equal(t, "01a1bbca-0c21-7251-89fa-29c0a945af83", network.Version)
			}
		})
	}
}

func TestUpdateNetwork_InvalidResponseBody(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"taskUid":"not-a-number"}`))
	}))
	defer srv.Close()

	got, err := New(srv.URL).UpdateNetwork(context.Background(), &UpdateNetworkRequest{Self: String("ms-00")})
	require.Nil(t, got)

	var meiliErr *Error
	require.True(t, errors.As(err, &meiliErr))
	require.Equal(t, ErrCodeResponseUnmarshalBody, meiliErr.ErrCode)
	require.Equal(t, "UpdateNetwork", meiliErr.Function)
}
