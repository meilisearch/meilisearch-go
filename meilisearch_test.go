package meilisearch

import (
	"context"
	"crypto/tls"
	"errors"
	"math"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Version(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	tests := []struct {
		name   string
		client ServiceManager
	}{
		{
			name:   "TestVersion",
			client: sv,
		},
		{
			name:   "TestVersionWithCustomClient",
			client: customSv,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.client.Version()
			require.NoError(t, err)
			require.NotNil(t, gotResp, "Version() should not return nil value")
		})
	}
}

func TestClient_TimeoutError(t *testing.T) {
	sv := setup(t, "")

	t.Cleanup(cleanup(sv))

	tests := []struct {
		name          string
		sv            ServiceManager
		expectedError Error
	}{
		{
			name: "TestTimeoutError",
			sv:   sv,
			expectedError: Error{
				MeilisearchApiError: meilisearchApiError{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			time.Sleep(2 * time.Second)
			gotResp, err := tt.sv.VersionWithContext(ctx)
			require.Error(t, err)
			require.Nil(t, gotResp)
		})
	}
}

func Test_GetStats(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	tests := []struct {
		name   string
		client ServiceManager
	}{
		{
			name:   "TestGetStats",
			client: sv,
		},
		{
			name:   "TestGetStatsWithCustomClient",
			client: customSv,
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

func Test_GetKey(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	t.Cleanup(cleanup(sv, customSv))

	tests := []struct {
		name   string
		client ServiceManager
	}{
		{
			name:   "TestGetKey",
			client: sv,
		},
		{
			name:   "TestGetKeyWithCustomClient",
			client: customSv,
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

func Test_GetKeys(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		client  ServiceManager
		request *KeysQuery
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestBasicGetKeys",
			args: args{
				client:  sv,
				request: nil,
			},
		},
		{
			name: "TestGetKeysWithCustomClient",
			args: args{
				client:  customSv,
				request: nil,
			},
		},
		{
			name: "TestGetKeysWithEmptyParam",
			args: args{
				client:  sv,
				request: &KeysQuery{},
			},
		},
		{
			name: "TestGetKeysWithLimit",
			args: args{
				client: sv,
				request: &KeysQuery{
					Limit: 1,
				},
			},
		},
		{
			name: "TestGetKeysWithOffset",
			args: args{
				client: sv,
				request: &KeysQuery{
					Limit:  2,
					Offset: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.args.client.GetKeys(tt.args.request)

			require.NoError(t, err)
			require.NotNil(t, gotResp, "GetKeys() should not return nil value")
			switch {
			case tt.args.request != nil && tt.args.request.Limit != 0 && tt.args.request.Offset == 0:
				require.Equal(t, tt.args.request.Limit, int64(len(gotResp.Results)))
			case tt.args.request != nil && tt.args.request.Limit == 2 && tt.args.request.Offset == 1:
				require.GreaterOrEqual(t, len(gotResp.Results), int(tt.args.request.Limit-tt.args.request.Offset))
			default:
				require.GreaterOrEqual(t, len(gotResp.Results), 2)
			}
		})
	}
}

func Test_CreateKey(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	tests := []struct {
		name   string
		client ServiceManager
		key    Key
	}{
		{
			name:   "TestCreateBasicKey",
			client: sv,
			key: Key{
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
		},
		{
			name:   "TestCreateKeyWithCustomClient",
			client: customSv,
			key: Key{
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
		},
		{
			name:   "TestCreateKeyWithExpirationAt",
			client: sv,
			key: Key{
				Actions:   []string{"*"},
				Indexes:   []string{"*"},
				ExpiresAt: time.Now().Add(time.Hour * 10),
			},
		},
		{
			name:   "TestCreateKeyWithDescription",
			client: sv,
			key: Key{
				Name:        "TestCreateKeyWithDescription",
				Description: "TestCreateKeyWithDescription",
				Actions:     []string{"*"},
				Indexes:     []string{"*"},
			},
		},
		{
			name:   "TestCreateKeyWithActions",
			client: sv,
			key: Key{
				Name:        "TestCreateKeyWithActions",
				Description: "TestCreateKeyWithActions",
				Actions:     []string{"documents.add", "documents.delete"},
				Indexes:     []string{"*"},
			},
		},
		{
			name:   "TestCreateKeyWithIndexes",
			client: sv,
			key: Key{
				Name:        "TestCreateKeyWithIndexes",
				Description: "TestCreateKeyWithIndexes",
				Actions:     []string{"*"},
				Indexes:     []string{"movies", "games"},
			},
		},
		{
			name:   "TestCreateKeyWithWildcardedAction",
			client: sv,
			key: Key{
				Name:        "TestCreateKeyWithWildcardedAction",
				Description: "TestCreateKeyWithWildcardedAction",
				Actions:     []string{"documents.*"},
				Indexes:     []string{"movies", "games"},
			},
		},
		{
			name:   "TestCreateKeyWithUID",
			client: sv,
			key: Key{
				Name:    "TestCreateKeyWithUID",
				UID:     "9aec34f4-e44c-4917-86c2-9c9403abb3b6",
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
		},
		{
			name:   "TestCreateKeyWithAllOptions",
			client: sv,
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
			t.Cleanup(cleanup(tt.client))

			gotResp, err := tt.client.CreateKey(&tt.key)
			require.NoError(t, err)

			gotKey, err := tt.client.GetKey(gotResp.Key)
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

func Test_UpdateKey(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	tests := []struct {
		name        string
		client      ServiceManager
		keyToCreate Key
		keyToUpdate Key
	}{
		{
			name:   "TestUpdateKeyWithDescription",
			client: sv,
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
			client: customSv,
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
			client: sv,
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
			client: sv,
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

func Test_DeleteKey(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	tests := []struct {
		name   string
		client ServiceManager
		key    Key
	}{
		{
			name:   "TestDeleteBasicKey",
			client: sv,
			key: Key{
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
		},
		{
			name:   "TestDeleteKeyWithCustomClient",
			client: customSv,
			key: Key{
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
		},
		{
			name:   "TestDeleteKeyWithExpirationAt",
			client: sv,
			key: Key{
				Actions:   []string{"*"},
				Indexes:   []string{"*"},
				ExpiresAt: time.Now().Add(time.Hour * 10),
			},
		},
		{
			name:   "TestDeleteKeyWithDescription",
			client: sv,
			key: Key{
				Description: "TestDeleteKeyWithDescription",
				Actions:     []string{"*"},
				Indexes:     []string{"*"},
			},
		},
		{
			name:   "TestDeleteKeyWithActions",
			client: sv,
			key: Key{
				Description: "TestDeleteKeyWithActions",
				Actions:     []string{"documents.add", "documents.delete"},
				Indexes:     []string{"*"},
			},
		},
		{
			name:   "TestDeleteKeyWithIndexes",
			client: sv,
			key: Key{
				Description: "TestDeleteKeyWithIndexes",
				Actions:     []string{"*"},
				Indexes:     []string{"movies", "games"},
			},
		},
		{
			name:   "TestDeleteKeyWithAllOptions",
			client: sv,
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

func Test_Health(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	badSv := setup(t, "http://wrongurl:1234")

	tests := []struct {
		name     string
		client   ServiceManager
		wantResp *Health
		wantErr  bool
	}{
		{
			name:   "TestHealth",
			client: sv,
			wantResp: &Health{
				Status: "available",
			},
			wantErr: false,
		},
		{
			name:   "TestHealthWithCustomClient",
			client: customSv,
			wantResp: &Health{
				Status: "available",
			},
			wantErr: false,
		},
		{
			name:    "TestHealthWithBadUrl",
			client:  badSv,
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

func Test_IsHealthy(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	badSv := setup(t, "http://wrongurl:1234")

	tests := []struct {
		name   string
		client ServiceManager
		want   bool
	}{
		{
			name:   "TestIsHealthy",
			client: sv,
			want:   true,
		},
		{
			name:   "TestIsHealthyWithCustomClient",
			client: customSv,
			want:   true,
		},
		{
			name:   "TestIsHealthyWIthBadUrl",
			client: badSv,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.client.IsHealthy()
			require.Equal(t, tt.want, got, "IsHealthy() got response %v, want %v", got, tt.want)
		})
	}
}

func Test_CreateDump(t *testing.T) {
	sv := setup(t, "")

	tests := []struct {
		name     string
		client   ServiceManager
		wantResp *Task
	}{
		{
			name:   "TestCreateDump",
			client: sv,
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

			taskInfo, err := c.WaitForTask(task.TaskUID, 0)

			require.NoError(t, err)
			require.NotNil(t, taskInfo)
			require.NotNil(t, taskInfo.Details)
			require.Equal(t, TaskStatusSucceeded, taskInfo.Status)
			require.NotEmpty(t, taskInfo.Details.DumpUid)
		})
	}
}

func Test_GetTask(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID      string
		client   ServiceManager
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
				client:  sv,
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
				client:  customSv,
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
				client:  sv,
				taskUID: 2,
				document: []docTest{
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
		},
	}

	t.Cleanup(cleanup(sv, customSv))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			task, err := i.AddDocuments(tt.args.document)
			require.NoError(t, err)

			_, err = c.WaitForTask(task.TaskUID, 0)
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

func Test_GetTasks(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID      string
		client   ServiceManager
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
				client: sv,
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
				client: customSv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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

			_, err = c.WaitForTask(task.TaskUID, 0)
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

func Test_GetTasksUsingClient(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID             string
		client          ServiceManager
		document        []docTest
		query           *TasksQuery
		expectedResults int
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestBasicGetTasks",
			args: args{
				UID:    "indexUID",
				client: sv,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query:           nil,
				expectedResults: 1,
			},
		},
		{
			name: "TestGetTasksWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customSv,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query:           nil,
				expectedResults: 1,
			},
		},
		{
			name: "TestGetTasksWithLimit",
			args: args{
				UID:    "indexUID",
				client: sv,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					Limit: 1,
				},
				expectedResults: 1,
			},
		},
		{
			name: "TestGetTasksWithLimit",
			args: args{
				UID:    "indexUID",
				client: sv,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					Limit: 1,
				},
				expectedResults: 1,
			},
		},
		{
			name: "TestGetTasksWithFrom",
			args: args{
				UID:    "indexUID",
				client: sv,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					From: 0,
				},
				expectedResults: 1,
			},
		},
		{
			name: "TestGetTasksWithFrom_1",
			args: args{
				UID:    "indexUID",
				client: sv,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					From: 1,
				},
				expectedResults: 0,
			},
		},
		{
			name: "TestGetTasksWithParameters",
			args: args{
				UID:    "indexUID",
				client: sv,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					Limit:     1,
					From:      0,
					IndexUIDS: []string{"indexUID"},
				},
				expectedResults: 1,
			},
		},
		{
			name: "TestGetTasksWithDateFilter",
			args: args{
				UID:    "indexUID",
				client: sv,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					Limit:            1,
					BeforeEnqueuedAt: time.Now(),
				},
				expectedResults: 1,
			},
		},

		{
			name: "TestGetTasksWithBeforeStartedAt",
			args: args{
				UID:    "indexUID",
				client: sv,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					Limit:           1,
					BeforeStartedAt: time.Now(),
				},
				expectedResults: 1,
			},
		},
		{
			name: "TestGetTasksWithAfterStartedAt",
			args: args{
				UID:    "indexUID",
				client: sv,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					Limit:          1,
					AfterStartedAt: time.Now().Add(-time.Hour),
				},
				expectedResults: 0,
			},
		},
		{
			name: "TestGetTasksWithBeforeFinishedAt",
			args: args{
				UID:    "indexUID",
				client: sv,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					Limit:            1,
					BeforeFinishedAt: time.Now().Add(time.Hour),
				},
				expectedResults: 1,
			},
		},
		{
			name: "TestGetTasksWithAfterFinishedAt",
			args: args{
				UID:    "indexUID",
				client: sv,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					Limit:           1,
					AfterFinishedAt: time.Now().Add(-time.Hour),
				},
				expectedResults: 0,
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

			_, err = c.WaitForTask(task.TaskUID, 0)
			require.NoError(t, err)

			gotResp, err := c.GetTasks(tt.args.query)
			require.NoError(t, err)
			require.NotNil(t, gotResp)
			// require.Equal(t, tt.args.expectedResults, len((*gotResp).Results))

			if tt.args.expectedResults > 0 {
				require.NotNil(t, (*gotResp).Results[0].Status)
				require.NotZero(t, (*gotResp).Results[0].UID)
				require.NotNil(t, (*gotResp).Results[0].Type)
			}
			if tt.args.query != nil {
				if tt.args.query.Limit != 0 {
					require.Equal(t, tt.args.query.Limit, (*gotResp).Limit)
				} else {
					require.Equal(t, int64(20), (*gotResp).Limit)
				}
				if tt.args.query.From != 0 && tt.args.expectedResults > 0 {
					require.Equal(t, tt.args.query.From, (*gotResp).From)
				}
			}
		})
	}
}

func Test_GetTasksUsingClientAllFailures(t *testing.T) {
	brokenSv := setup(t, "", WithAPIKey("wrong"))

	type args struct {
		UID             string
		client          ServiceManager
		document        []docTest
		query           *TasksQuery
		expectedResults int
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestBasicGetTasks",
			args: args{
				UID:    "indexUID",
				client: brokenSv,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query:           nil,
				expectedResults: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			t.Cleanup(cleanup(c))
			i := c.Index("NOT_EXISTS")

			_, err := c.DeleteIndex("NOT_EXISTS")
			require.Error(t, err)

			_, err = c.WaitForTask(math.MaxInt32, 0)
			require.Error(t, err)

			_, err = i.AddDocuments(tt.args.document)
			require.Error(t, err)

			_, err = c.GetTasks(tt.args.query)
			require.Error(t, err)

			_, err = c.GetStats()
			require.Error(t, err)

			_, err = c.CreateKey(&Key{
				Name: "Wrong",
			})
			require.Error(t, err)

			_, err = c.GetKey("Wrong")
			require.Error(t, err)

			_, err = c.UpdateKey("Wrong", &Key{
				Name: "Wrong",
			})
			require.Error(t, err)

			_, err = c.CreateDump()
			require.Error(t, err)

			_, err = c.GetTask(1)
			require.Error(t, err)

			_, err = c.DeleteTasks(nil)
			require.Error(t, err)

			_, err = c.SwapIndexes([]*SwapIndexesParams{
				{Indexes: []string{"Wrong", "Worse"}},
			})
			require.Error(t, err)
		})
	}
}

func Test_CancelTasks(t *testing.T) {
	sv := setup(t, "")

	type args struct {
		UID    string
		client ServiceManager
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
				client: sv,
				query:  nil,
			},
			want: "",
		},
		{
			name: "TestCancelTasksWithStatutes",
			args: args{
				UID:    "indexUID",
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
			} else {
				require.NoError(t, err)

				_, err = c.WaitForTask(gotResp.TaskUID, 0)
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

func Test_DeleteTasks(t *testing.T) {
	sv := setup(t, "")

	type args struct {
		UID    string
		client ServiceManager
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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

			_, err = c.WaitForTask(gotResp.TaskUID, 0)
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

func Test_SwapIndexes(t *testing.T) {
	sv := setup(t, "")

	type args struct {
		UID    string
		client ServiceManager
		query  []*SwapIndexesParams
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestBasicSwapIndexes",
			args: args{
				UID:    "indexUID",
				client: sv,
				query: []*SwapIndexesParams{
					{Indexes: []string{"IndexA", "IndexB"}},
				},
			},
		},
		{
			name: "TestSwapIndexesWithMultipleIndexes",
			args: args{
				UID:    "indexUID",
				client: sv,
				query: []*SwapIndexesParams{
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

			for _, params := range tt.args.query {
				for _, idx := range params.Indexes {
					task, err := c.CreateIndex(&IndexConfig{
						Uid: idx,
					})
					require.NoError(t, err)
					_, err = c.WaitForTask(task.TaskUID, 0)
					require.NoError(t, err)
				}
			}

			gotResp, err := c.SwapIndexes(tt.args.query)
			require.NoError(t, err)

			_, err = c.WaitForTask(gotResp.TaskUID, 0)
			require.NoError(t, err)

			gotTask, err := c.GetTask(gotResp.TaskUID)
			require.NoError(t, err)

			require.NotNil(t, gotResp.Status)
			require.NotNil(t, gotResp.Type)
			require.NotNil(t, gotResp.TaskUID)
			require.NotNil(t, gotResp.EnqueuedAt)
			require.Equal(t, gotTask.Status, TaskStatusSucceeded)
		})
	}
}

func Test_DefaultWaitForTask(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID      string
		client   ServiceManager
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
				client: sv,
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
				client: customSv,
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

			gotTask, err := c.WaitForTask(task.TaskUID, 0)
			require.NoError(t, err)
			require.Equal(t, tt.want, gotTask.Status)
		})
	}
}

func Test_WaitForTaskWithContext(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID      string
		client   ServiceManager
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
				client:   sv,
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
				client:   customSv,
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
				client:   sv,
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
				client:   sv,
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

			gotTask, err := c.WaitForTaskWithContext(ctx, task.TaskUID, 0)
			if tt.args.timeout < tt.args.interval {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, gotTask.Status)
			}
		})
	}
}

func Test_ConnectionCloseByServer(t *testing.T) {
	sv := setup(t, "")

	// Simulate 10 clients sending requests.
	g := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		g.Add(1)
		go func() {
			defer g.Done()

			_, _ = sv.Index("foo").Search("bar", &SearchRequest{})
			time.Sleep(5 * time.Second)
			_, err := sv.Index("foo").Search("bar", &SearchRequest{})
			var e *Error
			if errors.As(err, &e) && e.ErrCode == MeilisearchCommunicationError {
				require.NoErrorf(t, e, "unexpected error")
			}
		}()
	}
	g.Wait()
}

func Test_GenerateTenantToken(t *testing.T) {
	sv := setup(t, "")
	privateSv := setup(t, "", WithAPIKey(getPrivateKey(sv)))

	type args struct {
		IndexUIDS   string
		client      ServiceManager
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
				client:    privateSv,
				APIKeyUID: getPrivateUIDKey(sv),
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
				client:    sv,
				APIKeyUID: getPrivateUIDKey(sv),
				searchRules: map[string]interface{}{
					"*": map[string]string{},
				},
				options: &TenantTokenOptions{
					APIKey: getPrivateKey(sv),
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
				client:    privateSv,
				APIKeyUID: getPrivateUIDKey(sv),
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
				client:    sv,
				APIKeyUID: getPrivateUIDKey(sv),
				searchRules: map[string]interface{}{
					"*": map[string]string{},
				},
				options: &TenantTokenOptions{
					APIKey:    getPrivateKey(sv),
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
				client:    privateSv,
				APIKeyUID: getPrivateUIDKey(sv),
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
				client:    privateSv,
				APIKeyUID: getPrivateUIDKey(sv),
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
				client:      privateSv,
				APIKeyUID:   getPrivateUIDKey(sv),
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
				client:    setup(t, "", WithAPIKey("")),
				APIKeyUID: getPrivateUIDKey(sv),
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
				client:    sv,
				APIKeyUID: getPrivateUIDKey(sv),
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
				client:    sv,
				APIKeyUID: getPrivateUIDKey(sv) + "1234",
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
				client:    sv,
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
					_, err := setUpEmptyIndex(sv, &IndexConfig{Uid: tt.args.IndexUIDS})
					require.NoError(t, err, "CreateIndex() in TestGenerateTenantToken error should be nil")
				}

				client := setup(t, "", WithAPIKey(token))

				_, err = client.Index(tt.args.IndexUIDS).Search("", &SearchRequest{})

				require.NoError(t, err)
			}
		})
	}
}

func TestClient_MultiSearch(t *testing.T) {
	sv := setup(t, "")

	type args struct {
		client  ServiceManager
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
				client: sv,
				queries: &MultiSearchRequest{
					Queries: []*SearchRequest{
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
				client: sv,
				queries: &MultiSearchRequest{
					Queries: []*SearchRequest{
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
			name: "TestClientMultiSearchWithFederation",
			args: args{
				client: sv,
				queries: &MultiSearchRequest{
					Queries: []*SearchRequest{
						{
							IndexUID: "TestClientMultiSearchOnTwoIndexes1",
							Query:    "wonder",
						},
						{
							IndexUID: "TestClientMultiSearchOnTwoIndexes2",
							Query:    "prince",
						},
					},
					Federation: &MultiSearchFederation{},
				},
				UIDS: []string{"TestClientMultiSearchOnTwoIndexes1", "TestClientMultiSearchOnTwoIndexes2"},
			},
			want: &MultiSearchResponse{
				Results: nil,
				Hits: []interface{}{
					map[string]interface{}{"_federation": map[string]interface{}{"indexUid": "TestClientMultiSearchOnTwoIndexes2", "queriesPosition": 1.0, "weightedRankingScore": 0.8787878787878788}, "book_id": 456.0, "title": "Le Petit Prince"},
					map[string]interface{}{"_federation": map[string]interface{}{"indexUid": "TestClientMultiSearchOnTwoIndexes1", "queriesPosition": 0.0, "weightedRankingScore": 0.8712121212121212}, "book_id": 1.0, "title": "Alice In Wonderland"},
					map[string]interface{}{"_federation": map[string]interface{}{"indexUid": "TestClientMultiSearchOnTwoIndexes2", "queriesPosition": 1.0, "weightedRankingScore": 0.8333333333333334}, "book_id": 4.0, "title": "Harry Potter and the Half-Blood Prince"}},
				ProcessingTimeMs:   0,
				Offset:             0,
				Limit:              20,
				EstimatedTotalHits: 3,
				SemanticHitCount:   0,
			},
		},
		{
			name: "TestClientMultiSearchNoIndex",
			args: args{
				client: sv,
				queries: &MultiSearchRequest{
					Queries: []*SearchRequest{
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
				setUpBasicIndex(sv, UID)
			}
			c := tt.args.client
			t.Cleanup(cleanup(c))

			got, err := c.MultiSearch(tt.args.queries)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				got.ProcessingTimeMs = 0 // Can vary.
				require.Equal(t, got, tt.want)
			}
		})
	}
}

func Test_CreateIndex(t *testing.T) {
	tests := []struct {
		Name       string
		Encoding   ContentEncoding
		IndexUID   string
		PrimaryKey string
		WantErr    bool
	}{
		{
			Name:       "Basic create index",
			IndexUID:   "foobar",
			PrimaryKey: "id",
			WantErr:    false,
		},
		{
			Name:     "Create index without primary key",
			IndexUID: "foobar",
			WantErr:  false,
		},
		{
			Name:    "Empty index UID",
			WantErr: true,
		},
		{
			Name:     "Create index with content encoding gzip",
			IndexUID: "foobar",
			Encoding: GzipEncoding,
		},
		{
			Name:     "Create index with content encoding brotli",
			IndexUID: "foobar",
			Encoding: GzipEncoding,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			c := setup(t, "")
			if !tt.Encoding.IsZero() {
				c = setup(t, "", WithContentEncoding(tt.Encoding, DefaultCompression))
			}

			t.Cleanup(cleanup(c))

			info, err := c.CreateIndex(&IndexConfig{
				Uid:        tt.IndexUID,
				PrimaryKey: tt.PrimaryKey,
			})

			if tt.WantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, info)
				task, err := c.WaitForTask(info.TaskUID, 0)
				require.NoError(t, err)
				require.Equal(t, task.Status, TaskStatusSucceeded)
				got, err := c.GetIndex(tt.IndexUID)
				require.NoError(t, err)
				require.Equal(t, got.UID, tt.IndexUID)
				require.Equal(t, got.PrimaryKey, tt.PrimaryKey)
			}
		})
	}
}

func Test_ListIndex(t *testing.T) {
	tests := []struct {
		Name            string
		Indexes         []string
		ContentEncoding ContentEncoding
		WantErr         bool
	}{
		{
			Name:    "Basic get list of indexes",
			Indexes: []string{"foo", "bar"},
			WantErr: false,
		},
		{
			Name:            "Basic get list of indexes with deflate encoding",
			Indexes:         []string{"foo", "bar"},
			ContentEncoding: DeflateEncoding,
			WantErr:         false,
		},
		{
			Name:            "Basic get list of indexes with encoding",
			Indexes:         []string{"foo", "bar"},
			ContentEncoding: BrotliEncoding,
			WantErr:         false,
		},
		{
			Name:            "Get Empty list",
			ContentEncoding: BrotliEncoding,
			WantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			c := setup(t, "")
			if !tt.ContentEncoding.IsZero() {
				c = setup(t, "", WithContentEncoding(tt.ContentEncoding, DefaultCompression))
			}

			t.Cleanup(cleanup(c))

			for _, idx := range tt.Indexes {
				info, err := c.CreateIndex(&IndexConfig{
					Uid: idx,
				})
				if tt.WantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					require.NotNil(t, info)
					task, err := c.WaitForTask(info.TaskUID, 0)
					require.NoError(t, err)
					require.Equal(t, task.Status, TaskStatusSucceeded)
				}
			}

			got, err := c.ListIndexes(nil)
			if tt.WantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, len(got.Results), len(tt.Indexes))
			}
		})
	}
}

func Test_DeleteIndex(t *testing.T) {
	tests := []struct {
		Name            string
		IndexUID        string
		ContentEncoding ContentEncoding
		WantErr         bool
	}{
		{
			Name:     "Delete an index",
			IndexUID: "foobar",
			WantErr:  false,
		},
		{
			Name:     "Delete an index with encoding",
			IndexUID: "foobar",
			WantErr:  false,
		},
		{
			Name:    "Got Error on delete index",
			WantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			c := setup(t, "")
			if !tt.ContentEncoding.IsZero() {
				c = setup(t, "", WithContentEncoding(tt.ContentEncoding, DefaultCompression))
			}

			t.Cleanup(cleanup(c))

			if len(tt.IndexUID) != 0 {
				info, err := c.CreateIndex(&IndexConfig{
					Uid: tt.IndexUID,
				})
				if tt.WantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					task, err := c.WaitForTask(info.TaskUID, 0)
					require.NoError(t, err)
					require.Equal(t, task.Status, TaskStatusSucceeded)
				}
			}

			info, err := c.DeleteIndex(tt.IndexUID)
			if tt.WantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, info)
				task, err := c.WaitForTask(info.TaskUID, 0)
				require.NoError(t, err)
				require.Equal(t, task.Status, TaskStatusSucceeded)
			}

		})
	}
}

func Test_CreateSnapshot(t *testing.T) {
	c := setup(t, "")
	task, err := c.CreateSnapshot()
	require.NoError(t, err)
	testWaitForTask(t, c.Index("indexUID"), task)
}

func TestGetServiceManagerAndReaders(t *testing.T) {
	c := setup(t, "")
	require.NotNil(t, c.ServiceReader())
	require.NotNil(t, c.TaskManager())
	require.NotNil(t, c.TaskReader())
	require.NotNil(t, c.KeyManager())
	require.NotNil(t, c.KeyReader())
}
