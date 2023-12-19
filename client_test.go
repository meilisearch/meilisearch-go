package meilisearch

import (
	"context"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func TestClient_Version(t *testing.T) {
	tests := []struct {
		name   string
		client *Client
	}{
		{
			name:   "TestVersion",
			client: defaultClient,
		},
		{
			name:   "TestVersionWithCustomClient",
			client: customClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.client.GetVersion()
			require.NoError(t, err)
			require.NotNil(t, gotResp, "Version() should not return nil value")
		})
	}
}

func TestClient_TimeoutError(t *testing.T) {
	tests := []struct {
		name          string
		client        *Client
		expectedError Error
	}{
		{
			name:   "TestTimeoutError",
			client: timeoutClient,
			expectedError: Error{
				MeilisearchApiError: meilisearchApiError{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.client.GetVersion()
			require.Error(t, err)
			require.Nil(t, gotResp)
			require.Equal(t, tt.expectedError.MeilisearchApiError.Code,
				err.(*Error).MeilisearchApiError.Code)
		})
	}
}

func TestClient_GetStats(t *testing.T) {
	tests := []struct {
		name   string
		client *Client
	}{
		{
			name:   "TestGetStats",
			client: defaultClient,
		},
		{
			name:   "TestGetStatsWithCustomClient",
			client: customClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.client.GetStats()
			require.NoError(t, err)
			require.NotNil(t, gotResp, "GetStats() should not return nil value")
		})
	}
}

func TestClient_GetKey(t *testing.T) {
	tests := []struct {
		name   string
		client *Client
	}{
		{
			name:   "TestGetKey",
			client: defaultClient,
		},
		{
			name:   "TestGetKeyWithCustomClient",
			client: customClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.client.GetKeys(nil)
			require.NoError(t, err)

			gotKey, err := tt.client.GetKey(gotResp.Results[0].Key)
			require.NoError(t, err)
			require.NotNil(t, gotKey.ExpiresAt)
			require.NotNil(t, gotKey.CreatedAt)
			require.NotNil(t, gotKey.UpdatedAt)
		})
	}
}

func TestClient_GetKeys(t *testing.T) {
	type args struct {
		client  *Client
		request *KeysQuery
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestBasicGetKeys",
			args: args{
				client:  defaultClient,
				request: nil,
			},
		},
		{
			name: "TestGetKeysWithCustomClient",
			args: args{
				client:  customClient,
				request: nil,
			},
		},
		{
			name: "TestGetKeysWithEmptyParam",
			args: args{
				client:  defaultClient,
				request: &KeysQuery{},
			},
		},
		{
			name: "TestGetKeysWithLimit",
			args: args{
				client: defaultClient,
				request: &KeysQuery{
					Limit: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.args.client.GetKeys(tt.args.request)

			require.NoError(t, err)
			require.NotNil(t, gotResp, "GetKeys() should not return nil value")
			if tt.args.request != nil && tt.args.request.Limit != 0 {
				require.Equal(t, tt.args.request.Limit, int64(len(gotResp.Results)))
			} else {
				require.GreaterOrEqual(t, len(gotResp.Results), 2)
			}
		})
	}
}

func TestClient_CreateKey(t *testing.T) {
	tests := []struct {
		name   string
		client *Client
		key    Key
	}{
		{
			name:   "TestCreateBasicKey",
			client: defaultClient,
			key: Key{
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
		},
		{
			name:   "TestCreateKeyWithCustomClient",
			client: customClient,
			key: Key{
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
		},
		{
			name:   "TestCreateKeyWithExpirationAt",
			client: defaultClient,
			key: Key{
				Actions:   []string{"*"},
				Indexes:   []string{"*"},
				ExpiresAt: time.Now().Add(time.Hour * 10),
			},
		},
		{
			name:   "TestCreateKeyWithDescription",
			client: defaultClient,
			key: Key{
				Name:        "TestCreateKeyWithDescription",
				Description: "TestCreateKeyWithDescription",
				Actions:     []string{"*"},
				Indexes:     []string{"*"},
			},
		},
		{
			name:   "TestCreateKeyWithActions",
			client: defaultClient,
			key: Key{
				Name:        "TestCreateKeyWithActions",
				Description: "TestCreateKeyWithActions",
				Actions:     []string{"documents.add", "documents.delete"},
				Indexes:     []string{"*"},
			},
		},
		{
			name:   "TestCreateKeyWithIndexes",
			client: defaultClient,
			key: Key{
				Name:        "TestCreateKeyWithIndexes",
				Description: "TestCreateKeyWithIndexes",
				Actions:     []string{"*"},
				Indexes:     []string{"movies", "games"},
			},
		},
		{
			name:   "TestCreateKeyWithWildcardedAction",
			client: defaultClient,
			key: Key{
				Name:        "TestCreateKeyWithWildcardedAction",
				Description: "TestCreateKeyWithWildcardedAction",
				Actions:     []string{"documents.*"},
				Indexes:     []string{"movies", "games"},
			},
		},
		{
			name:   "TestCreateKeyWithUID",
			client: defaultClient,
			key: Key{
				Name:    "TestCreateKeyWithUID",
				UID:     "9aec34f4-e44c-4917-86c2-9c9403abb3b6",
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
		},
		{
			name:   "TestCreateKeyWithAllOptions",
			client: defaultClient,
			key: Key{
				Name:        "TestCreateKeyWithAllOptions",
				Description: "TestCreateKeyWithAllOptions",
				UID:         "9aec34f4-e44c-4917-86c2-9c9403abb3b6",
				Actions:     []string{"documents.add", "documents.delete"},
				Indexes:     []string{"movies", "games"},
				ExpiresAt:   time.Now().Add(time.Hour * 10),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const Format = "2006-01-02T15:04:05"
			c := tt.client
			t.Cleanup(cleanup(c))

			gotResp, err := c.CreateKey(&tt.key)
			require.NoError(t, err)

			gotKey, err := c.GetKey(gotResp.Key)
			require.NoError(t, err)
			require.Equal(t, tt.key.Name, gotKey.Name)
			require.Equal(t, tt.key.Description, gotKey.Description)
			if tt.key.UID != "" {
				require.Equal(t, tt.key.UID, gotKey.UID)
			}
			require.Equal(t, tt.key.Actions, gotKey.Actions)
			require.Equal(t, tt.key.Indexes, gotKey.Indexes)
			if !tt.key.ExpiresAt.IsZero() {
				require.Equal(t, tt.key.ExpiresAt.Format(Format), gotKey.ExpiresAt.Format(Format))
			}
		})
	}
}

func TestClient_UpdateKey(t *testing.T) {
	tests := []struct {
		name        string
		client      *Client
		keyToCreate Key
		keyToUpdate Key
	}{
		{
			name:   "TestUpdateKeyWithDescription",
			client: defaultClient,
			keyToCreate: Key{
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
			keyToUpdate: Key{
				Description: "TestUpdateKeyWithDescription",
			},
		},
		{
			name:   "TestUpdateKeyWithCustomClientWithDescription",
			client: customClient,
			keyToCreate: Key{
				Actions: []string{"*"},
				Indexes: []string{"TestUpdateKeyWithCustomClientWithDescription"},
			},
			keyToUpdate: Key{
				Description: "TestUpdateKeyWithCustomClientWithDescription",
			},
		},
		{
			name:   "TestUpdateKeyWithName",
			client: defaultClient,
			keyToCreate: Key{
				Actions: []string{"*"},
				Indexes: []string{"TestUpdateKeyWithName"},
			},
			keyToUpdate: Key{
				Name: "TestUpdateKeyWithName",
			},
		},
		{
			name:   "TestUpdateKeyWithNameAndAction",
			client: defaultClient,
			keyToCreate: Key{
				Actions: []string{"search"},
				Indexes: []string{"*"},
			},
			keyToUpdate: Key{
				Name: "TestUpdateKeyWithName",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const Format = "2006-01-02T15:04:05"
			c := tt.client
			t.Cleanup(cleanup(c))

			gotResp, err := c.CreateKey(&tt.keyToCreate)
			require.NoError(t, err)

			if tt.keyToCreate.Description != "" {
				require.Equal(t, tt.keyToCreate.Description, gotResp.Description)
			}
			if len(tt.keyToCreate.Actions) != 0 {
				require.Equal(t, tt.keyToCreate.Actions, gotResp.Actions)
			}
			if len(tt.keyToCreate.Indexes) != 0 {
				require.Equal(t, tt.keyToCreate.Indexes, gotResp.Indexes)
			}
			if !tt.keyToCreate.ExpiresAt.IsZero() {
				require.Equal(t, tt.keyToCreate.ExpiresAt.Format(Format), gotResp.ExpiresAt.Format(Format))
			}

			gotKey, err := c.UpdateKey(gotResp.Key, &tt.keyToUpdate)
			require.NoError(t, err)

			if tt.keyToUpdate.Description != "" {
				require.Equal(t, tt.keyToUpdate.Description, gotKey.Description)
			}
			if len(tt.keyToUpdate.Actions) != 0 {
				require.Equal(t, tt.keyToUpdate.Actions, gotKey.Actions)
			}
			if len(tt.keyToUpdate.Indexes) != 0 {
				require.Equal(t, tt.keyToUpdate.Indexes, gotKey.Indexes)
			}
			if tt.keyToUpdate.Description != "" {
				require.Equal(t, tt.keyToUpdate.Name, gotKey.Name)
			}
		})
	}
}

func TestClient_DeleteKey(t *testing.T) {
	tests := []struct {
		name   string
		client *Client
		key    Key
	}{
		{
			name:   "TestDeleteBasicKey",
			client: defaultClient,
			key: Key{
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
		},
		{
			name:   "TestDeleteKeyWithCustomClient",
			client: customClient,
			key: Key{
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
		},
		{
			name:   "TestDeleteKeyWithExpirationAt",
			client: defaultClient,
			key: Key{
				Actions:   []string{"*"},
				Indexes:   []string{"*"},
				ExpiresAt: time.Now().Add(time.Hour * 10),
			},
		},
		{
			name:   "TestDeleteKeyWithDescription",
			client: defaultClient,
			key: Key{
				Description: "TestDeleteKeyWithDescription",
				Actions:     []string{"*"},
				Indexes:     []string{"*"},
			},
		},
		{
			name:   "TestDeleteKeyWithActions",
			client: defaultClient,
			key: Key{
				Description: "TestDeleteKeyWithActions",
				Actions:     []string{"documents.add", "documents.delete"},
				Indexes:     []string{"*"},
			},
		},
		{
			name:   "TestDeleteKeyWithIndexes",
			client: defaultClient,
			key: Key{
				Description: "TestDeleteKeyWithIndexes",
				Actions:     []string{"*"},
				Indexes:     []string{"movies", "games"},
			},
		},
		{
			name:   "TestDeleteKeyWithAllOptions",
			client: defaultClient,
			key: Key{
				Description: "TestDeleteKeyWithAllOptions",
				Actions:     []string{"documents.add", "documents.delete"},
				Indexes:     []string{"movies", "games"},
				ExpiresAt:   time.Now().Add(time.Hour * 10),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client

			gotKey, err := c.CreateKey(&tt.key)
			require.NoError(t, err)

			gotResp, err := c.DeleteKey(gotKey.Key)
			require.NoError(t, err)
			require.True(t, gotResp)

			gotResp, err = c.DeleteKey(gotKey.Key)
			require.Error(t, err)
			require.False(t, gotResp)
		})
	}
}

func TestClient_Health(t *testing.T) {
	tests := []struct {
		name     string
		client   *Client
		wantResp *Health
		wantErr  bool
	}{
		{
			name:   "TestHealth",
			client: defaultClient,
			wantResp: &Health{
				Status: "available",
			},
			wantErr: false,
		},
		{
			name:   "TestHealthWithCustomClient",
			client: customClient,
			wantResp: &Health{
				Status: "available",
			},
			wantErr: false,
		},
		{
			name: "TestHealthWithBadUrl",
			client: &Client{
				config: ClientConfig{
					Host:   "http://wrongurl:1234",
					APIKey: masterKey,
				},
				httpClient: &fasthttp.Client{
					Name: "meilisearch-client",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.client.Health()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantResp, gotResp, "Health() got response %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestClient_IsHealthy(t *testing.T) {
	tests := []struct {
		name   string
		client *Client
		want   bool
	}{
		{
			name:   "TestIsHealthy",
			client: defaultClient,
			want:   true,
		},
		{
			name:   "TestIsHealthyWithCustomClient",
			client: customClient,
			want:   true,
		},
		{
			name: "TestIsHealthyWIthBadUrl",
			client: &Client{
				config: ClientConfig{
					Host:   "http://wrongurl:1234",
					APIKey: masterKey,
				},
				httpClient: &fasthttp.Client{
					Name: "meilisearch-client",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.client.IsHealthy()
			require.Equal(t, tt.want, got, "IsHealthy() got response %v, want %v", got, tt.want)
		})
	}
}

func TestClient_CreateDump(t *testing.T) {
	tests := []struct {
		name     string
		client   *Client
		wantResp *Task
	}{
		{
			name:   "TestCreateDump",
			client: defaultClient,
			wantResp: &Task{
				Status: "enqueued",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client

			task, err := c.CreateDump()
			require.NoError(t, err)
			if assert.NotNil(t, task, "CreateDump() should not return nil value") {
				require.Equal(t, tt.wantResp.Status, task.Status, "CreateDump() got response status %v, want: %v", task.Status, tt.wantResp.Status)
			}

			taskInfo, err := c.WaitForTask(task.TaskUID)

			require.NoError(t, err)
			require.NotNil(t, taskInfo)
			require.NotNil(t, taskInfo.Details)
			require.Equal(t, TaskStatusSucceeded, taskInfo.Status)
			require.NotEmpty(t, taskInfo.Details.DumpUid)
		})
	}
}

func TestClient_GetTask(t *testing.T) {
	type args struct {
		UID      string
		client   *Client
		taskUID  int64
		document []docTest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestBasicGetTask",
			args: args{
				UID:     "TestBasicGetTask",
				client:  defaultClient,
				taskUID: 0,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
		},
		{
			name: "TestGetTaskWithCustomClient",
			args: args{
				UID:     "TestGetTaskWithCustomClient",
				client:  customClient,
				taskUID: 1,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
		},
		{
			name: "TestGetTask",
			args: args{
				UID:     "TestGetTask",
				client:  defaultClient,
				taskUID: 2,
				document: []docTest{
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
		},
	}

	t.Cleanup(cleanup(defaultClient))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			task, err := i.AddDocuments(tt.args.document)
			require.NoError(t, err)

			_, err = c.WaitForTask(task.TaskUID)
			require.NoError(t, err)

			gotResp, err := c.GetTask(task.TaskUID)
			require.NoError(t, err)
			require.NotNil(t, gotResp)
			require.NotNil(t, gotResp.Details)
			require.GreaterOrEqual(t, gotResp.UID, tt.args.taskUID)
			require.Equal(t, tt.args.UID, gotResp.IndexUID)
			require.Equal(t, TaskStatusSucceeded, gotResp.Status)
			require.Equal(t, int64(len(tt.args.document)), gotResp.Details.ReceivedDocuments)
			require.Equal(t, int64(len(tt.args.document)), gotResp.Details.IndexedDocuments)

			// Make sure that timestamps are also retrieved
			require.NotZero(t, gotResp.EnqueuedAt)
			require.NotZero(t, gotResp.StartedAt)
			require.NotZero(t, gotResp.FinishedAt)
		})
	}
}

func TestClient_GetTasks(t *testing.T) {
	type args struct {
		UID      string
		client   *Client
		document []docTest
		query    *TasksQuery
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestBasicGetTasks",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: nil,
			},
		},
		{
			name: "TestGetTasksWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: nil,
			},
		},
		{
			name: "TestGetTasksWithLimit",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					Limit: 1,
				},
			},
		},
		{
			name: "TestGetTasksWithLimit",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					Limit: 1,
				},
			},
		},
		{
			name: "TestGetTasksWithFrom",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					From: 0,
				},
			},
		},
		{
			name: "TestGetTasksWithParameters",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					Limit:     1,
					From:      0,
					IndexUIDS: []string{"indexUID"},
				},
			},
		},
		{
			name: "TestGetTasksWithUidFilter",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					Limit: 1,
					UIDS:  []int64{1},
				},
			},
		},
		{
			name: "TestGetTasksWithDateFilter",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					Limit:            1,
					BeforeEnqueuedAt: time.Now(),
				},
			},
		},
		{
			name: "TestGetTasksWithCanceledByFilter",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					Limit:      1,
					CanceledBy: []int64{1},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			task, err := i.AddDocuments(tt.args.document)
			require.NoError(t, err)

			_, err = c.WaitForTask(task.TaskUID)
			require.NoError(t, err)

			gotResp, err := i.GetTasks(tt.args.query)
			require.NoError(t, err)
			require.NotNil(t, (*gotResp).Results[0].Status)
			require.NotZero(t, (*gotResp).Results[0].UID)
			require.NotNil(t, (*gotResp).Results[0].Type)
			if tt.args.query != nil {
				if tt.args.query.Limit != 0 {
					require.Equal(t, tt.args.query.Limit, (*gotResp).Limit)
				} else {
					require.Equal(t, int64(20), (*gotResp).Limit)
				}
				if tt.args.query.From != 0 {
					require.Equal(t, tt.args.query.From, (*gotResp).From)
				}
			}
		})
	}
}

func TestClient_CancelTasks(t *testing.T) {
	type args struct {
		UID    string
		client *Client
		query  *CancelTasksQuery
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestCancelTasksWithNoFilters",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  nil,
			},
			want: "",
		},
		{
			name: "TestCancelTasksWithStatutes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query: &CancelTasksQuery{
					Statuses: []TaskStatus{TaskStatusSucceeded},
				},
			},
			want: "?statuses=succeeded",
		},
		{
			name: "TestCancelTasksWithIndexUIDFilter",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query: &CancelTasksQuery{
					IndexUIDS: []string{"0"},
				},
			},
			want: "?indexUids=0",
		},
		{
			name: "TestCancelTasksWithMultipleIndexUIDsFilter",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query: &CancelTasksQuery{
					IndexUIDS: []string{"0", "1"},
				},
			},
			want: "?indexUids=0%2C1",
		},
		{
			name: "TestCancelTasksWithUidFilter",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query: &CancelTasksQuery{
					UIDS: []int64{0},
				},
			},
			want: "?uids=0",
		},
		{
			name: "TestCancelTasksWithMultipleUidsFilter",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query: &CancelTasksQuery{
					UIDS: []int64{0, 1},
				},
			},
			want: "?uids=0%2C1",
		},
		{
			name: "TestCancelTasksWithDateFilter",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query: &CancelTasksQuery{
					BeforeEnqueuedAt: time.Now(),
				},
			},
			want: strings.NewReplacer(":", "%3A").Replace("?beforeEnqueuedAt=" + time.Now().Format("2006-01-02T15:04:05Z")),
		},
		{
			name: "TestCancelTasksWithParameters",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query: &CancelTasksQuery{
					Statuses:        []TaskStatus{TaskStatusEnqueued},
					Types:           []TaskType{TaskTypeDocumentAdditionOrUpdate},
					IndexUIDS:       []string{"indexUID"},
					UIDS:            []int64{1},
					AfterEnqueuedAt: time.Now(),
				},
			},
			want: "?afterEnqueuedAt=" + strings.NewReplacer(":", "%3A").Replace(time.Now().Format("2006-01-02T15:04:05Z")) + "&indexUids=indexUID&statuses=enqueued&types=documentAdditionOrUpdate&uids=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			t.Cleanup(cleanup(c))

			gotResp, err := c.CancelTasks(tt.args.query)
			if tt.args.query == nil {
				require.Error(t, err)
				require.Equal(t, "missing_task_filters",
					err.(*Error).MeilisearchApiError.Code)
			} else {
				require.NoError(t, err)

				_, err = c.WaitForTask(gotResp.TaskUID)
				require.NoError(t, err)

				gotTask, err := c.GetTask(gotResp.TaskUID)
				require.NoError(t, err)

				require.NotNil(t, gotResp.Status)
				require.NotNil(t, gotResp.Type)
				require.NotNil(t, gotResp.TaskUID)
				require.NotNil(t, gotResp.EnqueuedAt)
				require.Equal(t, "", gotResp.IndexUID)
				require.Equal(t, TaskTypeTaskCancelation, gotResp.Type)
				require.Equal(t, tt.want, gotTask.Details.OriginalFilter)
			}
		})
	}
}

func TestClient_DeleteTasks(t *testing.T) {
	type args struct {
		UID    string
		client *Client
		query  *DeleteTasksQuery
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestBasicDeleteTasks",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query: &DeleteTasksQuery{
					Statuses: []TaskStatus{TaskStatusEnqueued},
				},
			},
			want: "?statuses=enqueued",
		},
		{
			name: "TestDeleteTasksWithUidFilter",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query: &DeleteTasksQuery{
					UIDS: []int64{1},
				},
			},
			want: "?uids=1",
		},
		{
			name: "TestDeleteTasksWithMultipleUidsFilter",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query: &DeleteTasksQuery{
					UIDS: []int64{0, 1},
				},
			},
			want: "?uids=0%2C1",
		},
		{
			name: "TestDeleteTasksWithIndexUIDFilter",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query: &DeleteTasksQuery{
					IndexUIDS: []string{"0"},
				},
			},
			want: "?indexUids=0",
		},
		{
			name: "TestDeleteTasksWithMultipleIndexUIDsFilter",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query: &DeleteTasksQuery{
					IndexUIDS: []string{"0", "1"},
				},
			},
			want: "?indexUids=0%2C1",
		},
		{
			name: "TestDeleteTasksWithDateFilter",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query: &DeleteTasksQuery{
					BeforeEnqueuedAt: time.Now(),
				},
			},
			want: strings.NewReplacer(":", "%3A").Replace("?beforeEnqueuedAt=" + time.Now().Format("2006-01-02T15:04:05Z")),
		},
		{
			name: "TestDeleteTasksWithCanceledByFilter",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query: &DeleteTasksQuery{
					CanceledBy: []int64{1},
				},
			},
			want: "?canceledBy=1",
		},
		{
			name: "TestDeleteTasksWithParameters",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query: &DeleteTasksQuery{
					Statuses:        []TaskStatus{TaskStatusEnqueued},
					IndexUIDS:       []string{"indexUID"},
					UIDS:            []int64{1},
					AfterEnqueuedAt: time.Now(),
				},
			},
			want: "?afterEnqueuedAt=" + strings.NewReplacer(":", "%3A").Replace(time.Now().Format("2006-01-02T15:04:05Z")) + "&indexUids=indexUID&statuses=enqueued&uids=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			t.Cleanup(cleanup(c))

			gotResp, err := c.DeleteTasks(tt.args.query)
			require.NoError(t, err)

			_, err = c.WaitForTask(gotResp.TaskUID)
			require.NoError(t, err)

			gotTask, err := c.GetTask(gotResp.TaskUID)
			require.NoError(t, err)

			require.NotNil(t, gotResp.Status)
			require.NotNil(t, gotResp.Type)
			require.NotNil(t, gotResp.TaskUID)
			require.NotNil(t, gotResp.EnqueuedAt)
			require.Equal(t, "", gotResp.IndexUID)
			require.Equal(t, TaskTypeTaskDeletion, gotResp.Type)
			require.NotNil(t, gotTask.Details.OriginalFilter)
			require.Equal(t, tt.want, gotTask.Details.OriginalFilter)
		})
	}
}

func TestClient_SwapIndexes(t *testing.T) {
	type args struct {
		UID    string
		client *Client
		query  []SwapIndexesParams
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestBasicSwapIndexes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query: []SwapIndexesParams{
					{Indexes: []string{"IndexA", "IndexB"}},
				},
			},
		},
		{
			name: "TestSwapIndexesWithMultipleIndexes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query: []SwapIndexesParams{
					{Indexes: []string{"IndexA", "IndexB"}},
					{Indexes: []string{"Index1", "Index2"}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			t.Cleanup(cleanup(c))

			gotResp, err := c.SwapIndexes(tt.args.query)
			require.NoError(t, err)

			_, err = c.WaitForTask(gotResp.TaskUID)
			require.NoError(t, err)

			gotTask, err := c.GetTask(gotResp.TaskUID)
			require.NoError(t, err)

			require.NotNil(t, gotResp.Status)
			require.NotNil(t, gotResp.Type)
			require.NotNil(t, gotResp.TaskUID)
			require.NotNil(t, gotResp.EnqueuedAt)
			require.Equal(t, "", gotResp.IndexUID)
			require.Equal(t, TaskTypeIndexSwap, gotResp.Type)
			require.Equal(t, tt.args.query, gotTask.Details.Swaps)
		})
	}
}

func TestClient_DefaultWaitForTask(t *testing.T) {
	type args struct {
		UID      string
		client   *Client
		taskUID  *Task
		document []docTest
	}
	tests := []struct {
		name string
		args args
		want TaskStatus
	}{
		{
			name: "TestDefaultWaitForTask",
			args: args{
				UID:    "TestDefaultWaitForTask",
				client: defaultClient,
				taskUID: &Task{
					UID: 0,
				},
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			want: "succeeded",
		},
		{
			name: "TestDefaultWaitForTaskWithCustomClient",
			args: args{
				UID:    "TestDefaultWaitForTaskWithCustomClient",
				client: customClient,
				taskUID: &Task{
					UID: 0,
				},
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			want: "succeeded",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			t.Cleanup(cleanup(c))

			task, err := c.Index(tt.args.UID).AddDocuments(tt.args.document)
			require.NoError(t, err)

			gotTask, err := c.WaitForTask(task.TaskUID)
			require.NoError(t, err)
			require.Equal(t, tt.want, gotTask.Status)
		})
	}
}

func TestClient_WaitForTaskWithContext(t *testing.T) {
	type args struct {
		UID      string
		client   *Client
		interval time.Duration
		timeout  time.Duration
		taskUID  *Task
		document []docTest
	}
	tests := []struct {
		name string
		args args
		want TaskStatus
	}{
		{
			name: "TestWaitForTask50",
			args: args{
				UID:      "TestWaitForTask50",
				client:   defaultClient,
				interval: time.Millisecond * 50,
				timeout:  time.Second * 5,
				taskUID: &Task{
					UID: 0,
				},
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			want: "succeeded",
		},
		{
			name: "TestWaitForTask50WithCustomClient",
			args: args{
				UID:      "TestWaitForTask50WithCustomClient",
				client:   customClient,
				interval: time.Millisecond * 50,
				timeout:  time.Second * 5,
				taskUID: &Task{
					UID: 0,
				},
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			want: "succeeded",
		},
		{
			name: "TestWaitForTask10",
			args: args{
				UID:      "TestWaitForTask10",
				client:   defaultClient,
				interval: time.Millisecond * 10,
				timeout:  time.Second * 5,
				taskUID: &Task{
					UID: 1,
				},
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			want: "succeeded",
		},
		{
			name: "TestWaitForTaskWithTimeout",
			args: args{
				UID:      "TestWaitForTaskWithTimeout",
				client:   defaultClient,
				interval: time.Millisecond * 50,
				timeout:  time.Millisecond * 10,
				taskUID: &Task{
					UID: 1,
				},
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			want: "succeeded",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			t.Cleanup(cleanup(c))

			task, err := c.Index(tt.args.UID).AddDocuments(tt.args.document)
			require.NoError(t, err)

			ctx, cancelFunc := context.WithTimeout(context.Background(), tt.args.timeout)
			defer cancelFunc()

			gotTask, err := c.WaitForTask(task.TaskUID, WaitParams{Context: ctx, Interval: tt.args.interval})
			if tt.args.timeout < tt.args.interval {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, gotTask.Status)
			}
		})
	}
}

func TestClient_ConnectionCloseByServer(t *testing.T) {
	meili := NewClient(ClientConfig{Host: getenv("MEILISEARCH_URL", "http://localhost:7700")})

	// Simulate 10 clients sending requests.
	g := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		g.Add(1)
		go func() {
			defer g.Done()

			_, _ = meili.Index("foo").Search("bar", &SearchRequest{})
			time.Sleep(5 * time.Second)
			_, err := meili.Index("foo").Search("bar", &SearchRequest{})
			if e, ok := err.(*Error); ok && e.ErrCode == MeilisearchCommunicationError {
				require.NoErrorf(t, e, "unexpected error")
			}
		}()
	}
	g.Wait()
}

func TestClient_GenerateTenantToken(t *testing.T) {
	type args struct {
		IndexUIDS   string
		client      *Client
		APIKeyUID   string
		searchRules map[string]interface{}
		options     *TenantTokenOptions
		filter      []string
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantFilter bool
	}{
		{
			name: "TestDefaultGenerateTenantToken",
			args: args{
				IndexUIDS: "TestDefaultGenerateTenantToken",
				client:    privateClient,
				APIKeyUID: GetPrivateUIDKey(),
				searchRules: map[string]interface{}{
					"*": map[string]string{},
				},
				options: nil,
				filter:  nil,
			},
			wantErr:    false,
			wantFilter: false,
		},
		{
			name: "TestGenerateTenantTokenWithApiKey",
			args: args{
				IndexUIDS: "TestGenerateTenantTokenWithApiKey",
				client:    defaultClient,
				APIKeyUID: GetPrivateUIDKey(),
				searchRules: map[string]interface{}{
					"*": map[string]string{},
				},
				options: &TenantTokenOptions{
					APIKey: GetPrivateKey(),
				},
				filter: nil,
			},
			wantErr:    false,
			wantFilter: false,
		},
		{
			name: "TestGenerateTenantTokenWithOnlyExpiresAt",
			args: args{
				IndexUIDS: "TestGenerateTenantTokenWithOnlyExpiresAt",
				client:    privateClient,
				APIKeyUID: GetPrivateUIDKey(),
				searchRules: map[string]interface{}{
					"*": map[string]string{},
				},
				options: &TenantTokenOptions{
					ExpiresAt: time.Now().Add(time.Hour * 10),
				},
				filter: nil,
			},
			wantErr:    false,
			wantFilter: false,
		},
		{
			name: "TestGenerateTenantTokenWithApiKeyAndExpiresAt",
			args: args{
				IndexUIDS: "TestGenerateTenantTokenWithApiKeyAndExpiresAt",
				client:    defaultClient,
				APIKeyUID: GetPrivateUIDKey(),
				searchRules: map[string]interface{}{
					"*": map[string]string{},
				},
				options: &TenantTokenOptions{
					APIKey:    GetPrivateKey(),
					ExpiresAt: time.Now().Add(time.Hour * 10),
				},
				filter: nil,
			},
			wantErr:    false,
			wantFilter: false,
		},
		{
			name: "TestGenerateTenantTokenWithFilters",
			args: args{
				IndexUIDS: "indexUID",
				client:    privateClient,
				APIKeyUID: GetPrivateUIDKey(),
				searchRules: map[string]interface{}{
					"*": map[string]string{
						"filter": "book_id > 1000",
					},
				},
				options: nil,
				filter: []string{
					"book_id",
				},
			},
			wantErr:    false,
			wantFilter: true,
		},
		{
			name: "TestGenerateTenantTokenWithFilterOnOneINdex",
			args: args{
				IndexUIDS: "indexUID",
				client:    privateClient,
				APIKeyUID: GetPrivateUIDKey(),
				searchRules: map[string]interface{}{
					"indexUID": map[string]string{
						"filter": "year > 2000",
					},
				},
				options: nil,
				filter: []string{
					"year",
				},
			},
			wantErr:    false,
			wantFilter: true,
		},
		{
			name: "TestGenerateTenantTokenWithoutSearchRules",
			args: args{
				IndexUIDS:   "TestGenerateTenantTokenWithoutSearchRules",
				client:      privateClient,
				APIKeyUID:   GetPrivateUIDKey(),
				searchRules: nil,
				options:     nil,
				filter:      nil,
			},
			wantErr:    true,
			wantFilter: false,
		},
		{
			name: "TestGenerateTenantTokenWithoutApiKey",
			args: args{
				IndexUIDS: "TestGenerateTenantTokenWithoutApiKey",
				client: NewClient(ClientConfig{
					Host:   getenv("MEILISEARCH_URL", "http://localhost:7700"),
					APIKey: "",
				}),
				APIKeyUID: GetPrivateUIDKey(),
				searchRules: map[string]interface{}{
					"*": map[string]string{},
				},
				options: nil,
				filter:  nil,
			},
			wantErr:    true,
			wantFilter: false,
		},
		{
			name: "TestGenerateTenantTokenWithBadExpiresAt",
			args: args{
				IndexUIDS: "TestGenerateTenantTokenWithBadExpiresAt",
				client:    defaultClient,
				APIKeyUID: GetPrivateUIDKey(),
				searchRules: map[string]interface{}{
					"*": map[string]string{},
				},
				options: &TenantTokenOptions{
					ExpiresAt: time.Now().Add(-time.Hour * 10),
				},
				filter: nil,
			},
			wantErr:    true,
			wantFilter: false,
		},
		{
			name: "TestGenerateTenantTokenWithBadAPIKeyUID",
			args: args{
				IndexUIDS: "TestGenerateTenantTokenWithBadAPIKeyUID",
				client:    defaultClient,
				APIKeyUID: GetPrivateUIDKey() + "1234",
				searchRules: map[string]interface{}{
					"*": map[string]string{},
				},
				options: nil,
				filter:  nil,
			},
			wantErr:    true,
			wantFilter: false,
		},
		{
			name: "TestGenerateTenantTokenWithEmptyAPIKeyUID",
			args: args{
				IndexUIDS: "TestGenerateTenantTokenWithEmptyAPIKeyUID",
				client:    defaultClient,
				APIKeyUID: "",
				searchRules: map[string]interface{}{
					"*": map[string]string{},
				},
				options: nil,
				filter:  nil,
			},
			wantErr:    true,
			wantFilter: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			t.Cleanup(cleanup(c))

			token, err := c.GenerateTenantToken(tt.args.APIKeyUID, tt.args.searchRules, tt.args.options)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				if tt.wantFilter {
					gotTask, err := c.Index(tt.args.IndexUIDS).UpdateFilterableAttributes(&tt.args.filter)
					require.NoError(t, err, "UpdateFilterableAttributes() in TestGenerateTenantToken error should be nil")
					testWaitForTask(t, c.Index(tt.args.IndexUIDS), gotTask)
				} else {
					_, err := SetUpEmptyIndex(&IndexConfig{Uid: tt.args.IndexUIDS})
					require.NoError(t, err, "CreateIndex() in TestGenerateTenantToken error should be nil")
				}

				client := NewClient(ClientConfig{
					Host:   getenv("MEILISEARCH_URL", "http://localhost:7700"),
					APIKey: token,
				})

				_, err = client.Index(tt.args.IndexUIDS).Search("", &SearchRequest{})

				require.NoError(t, err)
			}
		})
	}
}

func TestClient_MultiSearch(t *testing.T) {
	type args struct {
		client  *Client
		queries *MultiSearchRequest
		UIDS    []string
	}
	tests := []struct {
		name    string
		args    args
		want    *MultiSearchResponse
		wantErr bool
	}{
		{
			name: "TestClientMultiSearchOneIndex",
			args: args{
				client: defaultClient,
				queries: &MultiSearchRequest{
					[]SearchRequest{
						{
							IndexUID: "TestClientMultiSearchOneIndex",
							Query:    "wonder",
						},
					},
				},
				UIDS: []string{"TestClientMultiSearchOneIndex"},
			},
			want: &MultiSearchResponse{
				Results: []SearchResponse{
					{
						Hits: []interface{}{
							map[string]interface{}{
								"book_id": float64(1),
								"title":   "Alice In Wonderland",
							},
						},
						EstimatedTotalHits: 1,
						Offset:             0,
						Limit:              20,
						Query:              "wonder",
						IndexUID:           "TestClientMultiSearchOneIndex",
					},
				},
			},
		},
		{
			name: "TestClientMultiSearchOnTwoIndexes",
			args: args{
				client: defaultClient,
				queries: &MultiSearchRequest{
					[]SearchRequest{
						{
							IndexUID: "TestClientMultiSearchOnTwoIndexes1",
							Query:    "wonder",
						},
						{
							IndexUID: "TestClientMultiSearchOnTwoIndexes2",
							Query:    "prince",
						},
					},
				},
				UIDS: []string{"TestClientMultiSearchOnTwoIndexes1", "TestClientMultiSearchOnTwoIndexes2"},
			},
			want: &MultiSearchResponse{
				Results: []SearchResponse{
					{
						Hits: []interface{}{
							map[string]interface{}{
								"book_id": float64(1),
								"title":   "Alice In Wonderland",
							},
						},
						EstimatedTotalHits: 1,
						Offset:             0,
						Limit:              20,
						Query:              "wonder",
						IndexUID:           "TestClientMultiSearchOnTwoIndexes1",
					},
					{
						Hits: []interface{}{
							map[string]interface{}{
								"book_id": float64(456),
								"title":   "Le Petit Prince",
							},
							map[string]interface{}{
								"book_id": float64(4),
								"title":   "Harry Potter and the Half-Blood Prince",
							},
						},
						EstimatedTotalHits: 2,
						Offset:             0,
						Limit:              20,
						Query:              "prince",
						IndexUID:           "TestClientMultiSearchOnTwoIndexes2",
					},
				},
			},
		},
		{
			name: "TestClientMultiSearchNoIndex",
			args: args{
				client: defaultClient,
				queries: &MultiSearchRequest{
					[]SearchRequest{
						{
							Query: "",
						},
					},
				},
				UIDS: []string{"TestClientMultiSearchNoIndex"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, UID := range tt.args.UIDS {
				SetUpBasicIndex(UID)
			}
			c := tt.args.client
			t.Cleanup(cleanup(c))

			got, err := c.MultiSearch(tt.args.queries)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				for i := 0; i < len(tt.want.Results); i++ {
					if !reflect.DeepEqual(got.Results[i].Hits, tt.want.Results[i].Hits) {
						t.Errorf("Client.MultiSearch() = %v, want %v", got.Results[i].Hits, tt.want.Results[i].Hits)
					}
					require.Equal(t, tt.want.Results[i].EstimatedTotalHits, got.Results[i].EstimatedTotalHits)
					require.Equal(t, tt.want.Results[i].Offset, got.Results[i].Offset)
					require.Equal(t, tt.want.Results[i].Limit, got.Results[i].Limit)
					require.Equal(t, tt.want.Results[i].Query, got.Results[i].Query)
					require.Equal(t, tt.want.Results[i].IndexUID, got.Results[i].IndexUID)
				}
			}
		})
	}
}
