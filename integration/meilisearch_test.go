package integration

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/meilisearch/meilisearch-go"
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
	customSv := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	tests := []struct {
		name   string
		client meilisearch.ServiceManager
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
		sv            meilisearch.ServiceManager
		expectedError bool
	}{
		{
			name:          "TestTimeoutError",
			sv:            sv,
			expectedError: true,
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
	customSv := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	tests := []struct {
		name   string
		client meilisearch.ServiceManager
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
	customSv := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	t.Cleanup(cleanup(sv, customSv))

	tests := []struct {
		name   string
		client meilisearch.ServiceManager
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
	customSv := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		client  meilisearch.ServiceManager
		request *meilisearch.KeysQuery
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
				request: &meilisearch.KeysQuery{},
			},
		},
		{
			name: "TestGetKeysWithLimit",
			args: args{
				client: sv,
				request: &meilisearch.KeysQuery{
					Limit: 1,
				},
			},
		},
		{
			name: "TestGetKeysWithOffset",
			args: args{
				client: sv,
				request: &meilisearch.KeysQuery{
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
	customSv := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	tests := []struct {
		name   string
		client meilisearch.ServiceManager
		Key    meilisearch.Key
	}{
		{
			name:   "TestCreateBasicKey",
			client: sv,
			Key: meilisearch.Key{
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
		},
		{
			name:   "TestCreateKeyWithCustomClient",
			client: customSv,
			Key: meilisearch.Key{
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
		},
		{
			name:   "TestCreateKeyWithExpirationAt",
			client: sv,
			Key: meilisearch.Key{
				Actions:   []string{"*"},
				Indexes:   []string{"*"},
				ExpiresAt: time.Now().Add(time.Hour * 10),
			},
		},
		{
			name:   "TestCreateKeyWithDescription",
			client: sv,
			Key: meilisearch.Key{
				Name:        "TestCreateKeyWithDescription",
				Description: "TestCreateKeyWithDescription",
				Actions:     []string{"*"},
				Indexes:     []string{"*"},
			},
		},
		{
			name:   "TestCreateKeyWithActions",
			client: sv,
			Key: meilisearch.Key{
				Name:        "TestCreateKeyWithActions",
				Description: "TestCreateKeyWithActions",
				Actions:     []string{"documents.add", "documents.delete"},
				Indexes:     []string{"*"},
			},
		},
		{
			name:   "TestCreateKeyWithIndexes",
			client: sv,
			Key: meilisearch.Key{
				Name:        "TestCreateKeyWithIndexes",
				Description: "TestCreateKeyWithIndexes",
				Actions:     []string{"*"},
				Indexes:     []string{"movies", "games"},
			},
		},
		{
			name:   "TestCreateKeyWithWildcardedAction",
			client: sv,
			Key: meilisearch.Key{
				Name:        "TestCreateKeyWithWildcardedAction",
				Description: "TestCreateKeyWithWildcardedAction",
				Actions:     []string{"documents.*"},
				Indexes:     []string{"movies", "games"},
			},
		},
		{
			name:   "TestCreateKeyWithUID",
			client: sv,
			Key: meilisearch.Key{
				Name:    "TestCreateKeyWithUID",
				UID:     "9aec34f4-e44c-4917-86c2-9c9403abb3b6",
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
		},
		{
			name:   "TestCreateKeyWithAllOptions",
			client: sv,
			Key: meilisearch.Key{
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

			gotResp, err := tt.client.CreateKey(&tt.Key)
			require.NoError(t, err)

			gotKey, err := tt.client.GetKey(gotResp.Key)
			require.NoError(t, err)
			require.Equal(t, tt.Key.Name, gotKey.Name)
			require.Equal(t, tt.Key.Description, gotKey.Description)
			if tt.Key.UID != "" {
				require.Equal(t, tt.Key.UID, gotKey.UID)
			}
			require.Equal(t, tt.Key.Actions, gotKey.Actions)
			require.Equal(t, tt.Key.Indexes, gotKey.Indexes)
			if !tt.Key.ExpiresAt.IsZero() {
				require.Equal(t, tt.Key.ExpiresAt.Format(Format), gotKey.ExpiresAt.Format(Format))
			}
		})
	}
}

func Test_UpdateKey(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	tests := []struct {
		name        string
		client      meilisearch.ServiceManager
		keyToCreate meilisearch.Key
		keyToUpdate meilisearch.Key
	}{
		{
			name:   "TestUpdateKeyWithDescription",
			client: sv,
			keyToCreate: meilisearch.Key{
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
			keyToUpdate: meilisearch.Key{
				Description: "TestUpdateKeyWithDescription",
			},
		},
		{
			name:   "TestUpdateKeyWithCustomClientWithDescription",
			client: customSv,
			keyToCreate: meilisearch.Key{
				Actions: []string{"*"},
				Indexes: []string{"TestUpdateKeyWithCustomClientWithDescription"},
			},
			keyToUpdate: meilisearch.Key{
				Description: "TestUpdateKeyWithCustomClientWithDescription",
			},
		},
		{
			name:   "TestUpdateKeyWithName",
			client: sv,
			keyToCreate: meilisearch.Key{
				Actions: []string{"*"},
				Indexes: []string{"TestUpdateKeyWithName"},
			},
			keyToUpdate: meilisearch.Key{
				Name: "TestUpdateKeyWithName",
			},
		},
		{
			name:   "TestUpdateKeyWithNameAndAction",
			client: sv,
			keyToCreate: meilisearch.Key{
				Actions: []string{"search"},
				Indexes: []string{"*"},
			},
			keyToUpdate: meilisearch.Key{
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
	customSv := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	tests := []struct {
		name   string
		client meilisearch.ServiceManager
		Key    meilisearch.Key
	}{
		{
			name:   "TestDeleteBasicKey",
			client: sv,
			Key: meilisearch.Key{
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
		},
		{
			name:   "TestDeleteKeyWithCustomClient",
			client: customSv,
			Key: meilisearch.Key{
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
		},
		{
			name:   "TestDeleteKeyWithExpirationAt",
			client: sv,
			Key: meilisearch.Key{
				Actions:   []string{"*"},
				Indexes:   []string{"*"},
				ExpiresAt: time.Now().Add(time.Hour * 10),
			},
		},
		{
			name:   "TestDeleteKeyWithDescription",
			client: sv,
			Key: meilisearch.Key{
				Description: "TestDeleteKeyWithDescription",
				Actions:     []string{"*"},
				Indexes:     []string{"*"},
			},
		},
		{
			name:   "TestDeleteKeyWithActions",
			client: sv,
			Key: meilisearch.Key{
				Description: "TestDeleteKeyWithActions",
				Actions:     []string{"documents.add", "documents.delete"},
				Indexes:     []string{"*"},
			},
		},
		{
			name:   "TestDeleteKeyWithIndexes",
			client: sv,
			Key: meilisearch.Key{
				Description: "TestDeleteKeyWithIndexes",
				Actions:     []string{"*"},
				Indexes:     []string{"movies", "games"},
			},
		},
		{
			name:   "TestDeleteKeyWithAllOptions",
			client: sv,
			Key: meilisearch.Key{
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

			gotKey, err := c.CreateKey(&tt.Key)
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
	customSv := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	badSv := setup(t, "http://wrongurl:1234")

	tests := []struct {
		name     string
		client   meilisearch.ServiceManager
		wantResp *meilisearch.Health
		wantErr  bool
	}{
		{
			name:   "TestHealth",
			client: sv,
			wantResp: &meilisearch.Health{
				Status: "available",
			},
			wantErr: false,
		},
		{
			name:   "TestHealthWithCustomClient",
			client: customSv,
			wantResp: &meilisearch.Health{
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
				require.Equal(t, tt.wantResp, gotResp, "meilisearch.Health() got response %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func Test_IsHealthy(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	badSv := setup(t, "http://wrongurl:1234")

	tests := []struct {
		name   string
		client meilisearch.ServiceManager
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
		client   meilisearch.ServiceManager
		wantResp *meilisearch.Task
	}{
		{
			name:   "TestCreateDump",
			client: sv,
			wantResp: &meilisearch.Task{
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
			require.Equal(t, meilisearch.TaskStatusSucceeded, taskInfo.Status)
			require.NotEmpty(t, taskInfo.Details.DumpUid)
		})
	}
}

func Test_GetTask(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID      string
		client   meilisearch.ServiceManager
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

			taskInfo, err := i.AddDocuments(tt.args.document, nil)
			require.NoError(t, err)

			_, err = c.WaitForTask(taskInfo.TaskUID, 0)
			require.NoError(t, err)

			gotResp, err := c.GetTask(taskInfo.TaskUID)
			require.NoError(t, err)
			require.NotNil(t, gotResp)
			require.NotNil(t, gotResp.Details)
			require.GreaterOrEqual(t, gotResp.UID, tt.args.taskUID)
			require.Equal(t, tt.args.UID, gotResp.IndexUID)
			require.Equal(t, meilisearch.TaskStatusSucceeded, gotResp.Status)
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
	customSv := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID      string
		client   meilisearch.ServiceManager
		document []docTest
		query    *meilisearch.TasksQuery
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
				query: &meilisearch.TasksQuery{
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
				query: &meilisearch.TasksQuery{
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
				query: &meilisearch.TasksQuery{
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
				query: &meilisearch.TasksQuery{
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
				query: &meilisearch.TasksQuery{
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
				query: &meilisearch.TasksQuery{
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
				query: &meilisearch.TasksQuery{
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

			taskInfo, err := i.AddDocuments(tt.args.document, nil)
			require.NoError(t, err)

			_, err = c.WaitForTask(taskInfo.TaskUID, 0)
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
	customSv := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID             string
		client          meilisearch.ServiceManager
		document        []docTest
		query           *meilisearch.TasksQuery
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
				query: &meilisearch.TasksQuery{
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
				query: &meilisearch.TasksQuery{
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
				query: &meilisearch.TasksQuery{
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
				query: &meilisearch.TasksQuery{
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
				query: &meilisearch.TasksQuery{
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
				query: &meilisearch.TasksQuery{
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
				query: &meilisearch.TasksQuery{
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
				query: &meilisearch.TasksQuery{
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
				query: &meilisearch.TasksQuery{
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
				query: &meilisearch.TasksQuery{
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

			taskInfo, err := i.AddDocuments(tt.args.document, nil)
			require.NoError(t, err)

			_, err = c.WaitForTask(taskInfo.TaskUID, 0)
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
	brokenSv := setup(t, "", meilisearch.WithAPIKey("wrong"))

	type args struct {
		UID             string
		client          meilisearch.ServiceManager
		document        []docTest
		query           *meilisearch.TasksQuery
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

			_, err = i.AddDocuments(tt.args.document, nil)
			require.Error(t, err)

			_, err = c.GetTasks(tt.args.query)
			require.Error(t, err)

			_, err = c.GetStats()
			require.Error(t, err)

			_, err = c.CreateKey(&meilisearch.Key{
				Name: "Wrong",
			})
			require.Error(t, err)

			_, err = c.GetKey("Wrong")
			require.Error(t, err)

			_, err = c.UpdateKey("Wrong", &meilisearch.Key{
				Name: "Wrong",
			})
			require.Error(t, err)

			_, err = c.CreateDump()
			require.Error(t, err)

			_, err = c.GetTask(1)
			require.Error(t, err)

			_, err = c.DeleteTasks(nil)
			require.Error(t, err)

			_, err = c.SwapIndexes([]*meilisearch.SwapIndexesParams{
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
		client meilisearch.ServiceManager
		query  *meilisearch.CancelTasksQuery
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
				query: &meilisearch.CancelTasksQuery{
					Statuses: []meilisearch.TaskStatus{meilisearch.TaskStatusSucceeded},
				},
			},
			want: "?statuses=succeeded",
		},
		{
			name: "TestCancelTasksWithIndexUIDFilter",
			args: args{
				UID:    "indexUID",
				client: sv,
				query: &meilisearch.CancelTasksQuery{
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
				query: &meilisearch.CancelTasksQuery{
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
				query: &meilisearch.CancelTasksQuery{
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
				query: &meilisearch.CancelTasksQuery{
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
				query: &meilisearch.CancelTasksQuery{
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
				query: &meilisearch.CancelTasksQuery{
					Statuses:        []meilisearch.TaskStatus{meilisearch.TaskStatusEnqueued},
					Types:           []meilisearch.TaskType{meilisearch.TaskTypeDocumentAdditionOrUpdate},
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
				require.Equal(t, meilisearch.TaskTypeTaskCancelation, gotResp.Type)
				require.Equal(t, tt.want, gotTask.Details.OriginalFilter)
			}
		})
	}
}

func Test_DeleteTasks(t *testing.T) {
	sv := setup(t, "")

	type args struct {
		UID    string
		client meilisearch.ServiceManager
		query  *meilisearch.DeleteTasksQuery
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
				query: &meilisearch.DeleteTasksQuery{
					Statuses: []meilisearch.TaskStatus{meilisearch.TaskStatusEnqueued},
				},
			},
			want: "?statuses=enqueued",
		},
		{
			name: "TestDeleteTasksWithUidFilter",
			args: args{
				UID:    "indexUID",
				client: sv,
				query: &meilisearch.DeleteTasksQuery{
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
				query: &meilisearch.DeleteTasksQuery{
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
				query: &meilisearch.DeleteTasksQuery{
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
				query: &meilisearch.DeleteTasksQuery{
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
				query: &meilisearch.DeleteTasksQuery{
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
				query: &meilisearch.DeleteTasksQuery{
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
				query: &meilisearch.DeleteTasksQuery{
					Statuses:        []meilisearch.TaskStatus{meilisearch.TaskStatusEnqueued},
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
			require.Equal(t, meilisearch.TaskTypeTaskDeletion, gotResp.Type)
			require.NotNil(t, gotTask.Details.OriginalFilter)
			require.Equal(t, tt.want, gotTask.Details.OriginalFilter)
		})
	}
}

func Test_SwapIndexes(t *testing.T) {
	sv := setup(t, "")

	type args struct {
		UID    string
		client meilisearch.ServiceManager
		query  []*meilisearch.SwapIndexesParams
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
				query: []*meilisearch.SwapIndexesParams{
					{Indexes: []string{"IndexA", "IndexB"}},
				},
			},
		},
		{
			name: "TestSwapIndexesWithMultipleIndexes",
			args: args{
				UID:    "indexUID",
				client: sv,
				query: []*meilisearch.SwapIndexesParams{
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
					taskInfo, err := c.CreateIndex(&meilisearch.IndexConfig{
						Uid: idx,
					})
					require.NoError(t, err)
					_, err = c.WaitForTask(taskInfo.TaskUID, 0)
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
			require.Equal(t, gotTask.Status, meilisearch.TaskStatusSucceeded)
		})
	}
}

func Test_DefaultWaitForTask(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID      string
		client   meilisearch.ServiceManager
		taskUID  *meilisearch.Task
		document []docTest
	}
	tests := []struct {
		name string
		args args
		want meilisearch.TaskStatus
	}{
		{
			name: "TestDefaultWaitForTask",
			args: args{
				UID:    "TestDefaultWaitForTask",
				client: sv,
				taskUID: &meilisearch.Task{
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
				taskUID: &meilisearch.Task{
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

			taskInfo, err := c.Index(tt.args.UID).AddDocuments(tt.args.document, nil)
			require.NoError(t, err)

			gotTask, err := c.WaitForTask(taskInfo.TaskUID, 0)
			require.NoError(t, err)
			require.Equal(t, tt.want, gotTask.Status)
		})
	}
}

func Test_WaitForTaskWithContext(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID      string
		client   meilisearch.ServiceManager
		interval time.Duration
		timeout  time.Duration
		taskUID  *meilisearch.Task
		document []docTest
	}
	tests := []struct {
		name string
		args args
		want meilisearch.TaskStatus
	}{
		{
			name: "TestWaitForTask50",
			args: args{
				UID:      "TestWaitForTask50",
				client:   sv,
				interval: time.Millisecond * 50,
				timeout:  time.Second * 5,
				taskUID: &meilisearch.Task{
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
				taskUID: &meilisearch.Task{
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
				taskUID: &meilisearch.Task{
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
				taskUID: &meilisearch.Task{
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

			taskInfo, err := c.Index(tt.args.UID).AddDocuments(tt.args.document, nil)
			require.NoError(t, err)

			ctx, cancelFunc := context.WithTimeout(context.Background(), tt.args.timeout)
			defer cancelFunc()

			gotTask, err := c.WaitForTaskWithContext(ctx, taskInfo.TaskUID, 0)
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

			_, _ = sv.Index("foo").Search("bar", &meilisearch.SearchRequest{})
			time.Sleep(5 * time.Second)
			_, err := sv.Index("foo").Search("bar", &meilisearch.SearchRequest{})
			var e *meilisearch.Error
			if errors.As(err, &e) && e.ErrCode == meilisearch.MeilisearchCommunicationError {
				require.NoErrorf(t, e, "unexpected meilisearch.Error")
			}
		}()
	}
	g.Wait()
}

func Test_GenerateTenantToken(t *testing.T) {
	sv := setup(t, "")
	privateSv := setup(t, "", meilisearch.WithAPIKey(getPrivateKey(sv)))

	type args struct {
		IndexUIDS   string
		client      meilisearch.ServiceManager
		APIKeyUID   string
		searchRules map[string]interface{}
		options     *TenantTokenOptions
		filter      []interface{}
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
				options: &meilisearch.TenantTokenOptions{
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
				options: &meilisearch.TenantTokenOptions{
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
				options: &meilisearch.TenantTokenOptions{
					APIKey:    getPrivateKey(sv),
					ExpiresAt: time.Now().Add(time.Hour * 10),
				},
				filter: nil,
			},
			wantErr:    false,
			wantFilter: false,
		},
		{
			name: "TestGenerateTenantTokenWithMultipleFilters",
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
				filter: []interface{}{
					"year",
					map[string]interface{}{
						"attributePatterns": []interface{}{"book_id"},
						"features": map[string]interface{}{
							"facetSearch": false,
							"filter": map[string]interface{}{
								"equality":   false,
								"comparison": true,
							},
						},
					},
				},
			},
			wantErr:    false,
			wantFilter: true,
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
				filter: []interface{}{
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
				filter: []interface{}{
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
				client:    setup(t, "", meilisearch.WithAPIKey("")),
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
				options: &meilisearch.TenantTokenOptions{
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
					require.NoError(t, err, "UpdateFilterableAttributes() in TestGenerateTenantToken meilisearch.Error should be nil")
					testWaitForTask(t, c.Index(tt.args.IndexUIDS), gotTask)
				} else {
					_, err := setUpEmptyIndex(sv, &meilisearch.IndexConfig{Uid: tt.args.IndexUIDS})
					require.NoError(t, err, "CreateIndex() in TestGenerateTenantToken meilisearch.Error should be nil")
				}

				client := setup(t, "", meilisearch.WithAPIKey(token))

				_, err = client.Index(tt.args.IndexUIDS).Search("", &meilisearch.SearchRequest{})

				require.NoError(t, err)
			}
		})
	}
}

func TestClient_MultiSearch(t *testing.T) {
	sv := setup(t, "")

	feat, err := sv.ExperimentalFeatures().SetNetwork(true).Update()
	require.NoError(t, err)
	require.NotNil(t, feat)
	require.True(t, feat.Network)

	type args struct {
		client  meilisearch.ServiceManager
		queries *meilisearch.MultiSearchRequest
		UIDS    []string
	}
	tests := []struct {
		name    string
		args    args
		want    *meilisearch.MultiSearchResponse
		setup   string
		wantErr bool
	}{
		{
			name: "TestClientMultiSearchOneIndex",
			args: args{
				client: sv,
				queries: &meilisearch.MultiSearchRequest{
					Queries: []*meilisearch.SearchRequest{
						{
							IndexUID: "TestClientMultiSearchOneIndex",
							Query:    "wonder",
						},
					},
				},
				UIDS: []string{"TestClientMultiSearchOneIndex"},
			},
			want: &meilisearch.MultiSearchResponse{
				Results: []meilisearch.SearchResponse{
					{
						Hits: meilisearch.Hits{
							{"book_id": toRawMessage(1), "title": toRawMessage("Alice In Wonderland")},
						},
						EstimatedTotalHits: 1,
						Offset:             0,
						Limit:              20,
						Query:              "wonder",
						IndexUID:           "TestClientMultiSearchOneIndex",
					},
				},
			},
			setup: "books",
		},
		{
			name: "TestClientMultiSearchOnTwoIndexes",
			args: args{
				client: sv,
				queries: &meilisearch.MultiSearchRequest{
					Queries: []*meilisearch.SearchRequest{
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
			want: &meilisearch.MultiSearchResponse{
				Results: []meilisearch.SearchResponse{
					{
						Hits: meilisearch.Hits{
							{"book_id": toRawMessage(1), "title": toRawMessage("Alice In Wonderland")},
						},
						EstimatedTotalHits: 1,
						Offset:             0,
						Limit:              20,
						Query:              "wonder",
						IndexUID:           "TestClientMultiSearchOnTwoIndexes1",
					},
					{
						Hits: meilisearch.Hits{
							{"book_id": toRawMessage(456), "title": toRawMessage("Le Petit Prince")},
							{"book_id": toRawMessage(4), "title": toRawMessage("Harry Potter and the Half-Blood Prince")},
						},
						EstimatedTotalHits: 2,
						Offset:             0,
						Limit:              20,
						Query:              "prince",
						IndexUID:           "TestClientMultiSearchOnTwoIndexes2",
					},
				},
			},
			setup: "books",
		},
		{
			name: "TestClientMultiSearchWithFederation",
			args: args{
				client: sv,
				queries: &meilisearch.MultiSearchRequest{
					Queries: []*meilisearch.SearchRequest{
						{
							IndexUID: "TestClientMultiSearchOnTwoIndexes1",
							Query:    "wonder",
						},
						{
							IndexUID: "TestClientMultiSearchOnTwoIndexes2",
							Query:    "prince",
						},
					},
					Federation: &meilisearch.MultiSearchFederation{},
				},
				UIDS: []string{"TestClientMultiSearchOnTwoIndexes1", "TestClientMultiSearchOnTwoIndexes2"},
			},
			want: &meilisearch.MultiSearchResponse{
				Results: nil,
				Hits: meilisearch.Hits{
					{
						"_federation": toRawMessage(map[string]interface{}{
							"indexUid": "TestClientMultiSearchOnTwoIndexes2", "queriesPosition": 1.0, "weightedRankingScore": 0.8787878787878788,
						}),
						"book_id": toRawMessage(456), "title": toRawMessage("Le Petit Prince"),
					},
					{
						"_federation": toRawMessage(map[string]interface{}{
							"indexUid": "TestClientMultiSearchOnTwoIndexes1", "queriesPosition": 0.0, "weightedRankingScore": 0.8712121212121212,
						}),
						"book_id": toRawMessage(1), "title": toRawMessage("Alice In Wonderland"),
					},
					{
						"_federation": toRawMessage(map[string]interface{}{
							"indexUid": "TestClientMultiSearchOnTwoIndexes2", "queriesPosition": 1.0, "weightedRankingScore": 0.8333333333333334,
						}),
						"book_id": toRawMessage(4), "title": toRawMessage("Harry Potter and the Half-Blood Prince"),
					},
				},
				ProcessingTimeMs:   0,
				Offset:             0,
				Limit:              20,
				EstimatedTotalHits: 3,
				SemanticHitCount:   0,
			},
			setup: "books",
		},
		{
			name: "TestClientMultiSearchWithFederationFacetsByIndex",
			args: args{
				client: sv,
				queries: &meilisearch.MultiSearchRequest{
					Federation: &meilisearch.MultiSearchFederation{
						FacetsByIndex: map[string][]string{
							"movies": {
								"title",
								"id",
							},
							"comics": {
								"title",
							},
						},
					},
					Queries: []*meilisearch.SearchRequest{
						{
							IndexUID: "movies",
							Query:    "Batman",
						}, {
							IndexUID: "comics",
							Query:    "Batman",
						},
					},
				},
				UIDS: []string{"movies", "comics"},
			},
			want: &meilisearch.MultiSearchResponse{
				Results: nil,
				Hits: meilisearch.Hits{

					{
						"id":           toRawMessage(31),
						"title":        toRawMessage("Batman"),
						"genres":       toRawMessage([]string{"Action", "Thriller"}),
						"overview":     toRawMessage("Follow the adventures of the Dark Knight as he battles crime in Gotham City."),
						"cover":        toRawMessage("https://example.com/comics/batman.jpg"),
						"release_date": toRawMessage(1625097600),
						"_federation": toRawMessage(map[string]interface{}{
							"indexUid":             "comics",
							"queriesPosition":      1,
							"weightedRankingScore": toRawMessage(1.0),
						}),
					},
					{
						"id":           toRawMessage(87),
						"title":        toRawMessage("The Batman"),
						"genres":       toRawMessage([]string{"Action", "Thriller"}),
						"overview":     toRawMessage("When a sadistic serial killer begins murdering key political figures in Gotham, the Batman is forced to investigate the city's hidden corruption and question his family's involvement."),
						"poster":       toRawMessage("https://example.com/comics/batman.jpg"),
						"release_date": toRawMessage(1625097600),
						"_federation": toRawMessage(map[string]interface{}{
							"indexUid":             "movies",
							"queriesPosition":      0,
							"weightedRankingScore": 0.9242424242424242,
						}),
					},
				},
				Offset:             0,
				Limit:              20,
				EstimatedTotalHits: 2,
				SemanticHitCount:   0,
			},
			setup: "movies",
		},
		{
			name: "TestClientMultiSearchWithFederationFacetsByIndexWithMergeFacets",
			args: args{
				client: sv,
				queries: &meilisearch.MultiSearchRequest{
					Federation: &meilisearch.MultiSearchFederation{
						FacetsByIndex: map[string][]string{
							"movies": {
								"title",
								"id",
							},
							"comics": {
								"title",
							},
						},
						MergeFacets: &meilisearch.MultiSearchFederationMergeFacets{
							MaxValuesPerFacet: 10,
						},
					},
					Queries: []*meilisearch.SearchRequest{
						{
							IndexUID: "movies",
							Query:    "Batman",
						}, {
							IndexUID: "comics",
							Query:    "Batman",
						},
					},
				},
				UIDS: []string{"movies", "comics"},
			},
			want: &meilisearch.MultiSearchResponse{
				Results: nil,
				Hits: meilisearch.Hits{

					{
						"id":           toRawMessage(31),
						"title":        toRawMessage("Batman"),
						"genres":       toRawMessage([]string{"Action", "Thriller"}),
						"overview":     toRawMessage("Follow the adventures of the Dark Knight as he battles crime in Gotham City."),
						"cover":        toRawMessage("https://example.com/comics/batman.jpg"),
						"release_date": toRawMessage(1625097600),
						"_federation": toRawMessage(map[string]interface{}{
							"indexUid":             "comics",
							"queriesPosition":      1,
							"weightedRankingScore": toRawMessage(1.0),
						}),
					},
					{
						"id":           toRawMessage(87),
						"title":        toRawMessage("The Batman"),
						"genres":       toRawMessage([]string{"Action", "Thriller"}),
						"overview":     toRawMessage("When a sadistic serial killer begins murdering key political figures in Gotham, the Batman is forced to investigate the city's hidden corruption and question his family's involvement."),
						"poster":       toRawMessage("https://example.com/comics/batman.jpg"),
						"release_date": toRawMessage(1625097600),
						"_federation": toRawMessage(map[string]interface{}{
							"indexUid":             "movies",
							"queriesPosition":      0,
							"weightedRankingScore": 0.9242424242424242,
						}),
					},
				},
				Offset:             0,
				Limit:              20,
				EstimatedTotalHits: 2,
				SemanticHitCount:   0,
			},
			setup: "movies",
		},
		{
			name: "TestClientMultiSearchNoIndex",
			args: args{
				client: sv,
				queries: &meilisearch.MultiSearchRequest{
					Queries: []*meilisearch.SearchRequest{
						{
							Query: "",
						},
					},
				},
				UIDS: []string{"TestClientMultiSearchNoIndex"},
			},
			setup:   "books",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, UID := range tt.args.UIDS {
				if tt.setup == "books" {
					setUpBasicIndex(sv, UID)
				} else {
					if UID == "movies" {
						setupMovieIndex(t, sv, UID)
					}

					if UID == "comics" {
						setupComicIndex(t, sv, UID)
					}
				}
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

				// Compare results while ignoring ProcessingTimeMs
				require.Equal(t, tt.want.Results, got.Results)
				require.Equal(t, tt.want.EstimatedTotalHits, got.EstimatedTotalHits)
				require.Equal(t, tt.want.SemanticHitCount, got.SemanticHitCount)
				require.Equal(t, tt.want.Offset, got.Offset)
				require.Equal(t, tt.want.Limit, got.Limit)

				for i := range tt.want.Hits {
					require.Equal(t, len(tt.want.Hits), len(got.Hits))

					var (
						wants map[string]interface{}
						gots  map[string]interface{}
					)

					err := tt.want.Hits[i].Decode(&wants)
					require.NoError(t, err)

					err = got.Hits[i].Decode(&gots)
					require.NoError(t, err)

					for k, v := range wants {
						gotVal := gots[k]
						switch wantFloat := v.(type) {
						case float64:
							require.InEpsilon(t, wantFloat, gotVal.(float64), 1e-9)
						default:
							require.Equal(t, v, gotVal)
						}
					}
				}
			}
		})
	}
}

func Test_CreateIndex(t *testing.T) {
	tests := []struct {
		Name       string
		Encoding   meilisearch.ContentEncoding
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
			Name:     "Create index without primary meilisearch.Key",
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
			Encoding: meilisearch.GzipEncoding,
		},
		{
			Name:     "Create index with content encoding brotli",
			IndexUID: "foobar",
			Encoding: meilisearch.GzipEncoding,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			c := setup(t, "")
			if !tt.Encoding.IsZero() {
				c = setup(t, "", meilisearch.WithContentEncoding(tt.Encoding, meilisearch.DefaultCompression))
			}

			t.Cleanup(cleanup(c))

			info, err := c.CreateIndex(&meilisearch.IndexConfig{
				Uid:        tt.IndexUID,
				PrimaryKey: tt.PrimaryKey,
			})

			if tt.WantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, info)
				taskInfo, err := c.WaitForTask(info.TaskUID, 0)
				require.NoError(t, err)
				require.Equal(t, taskInfo.Status, meilisearch.TaskStatusSucceeded)
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
		ContentEncoding meilisearch.ContentEncoding
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
			ContentEncoding: meilisearch.DeflateEncoding,
			WantErr:         false,
		},
		{
			Name:            "Basic get list of indexes with encoding",
			Indexes:         []string{"foo", "bar"},
			ContentEncoding: meilisearch.BrotliEncoding,
			WantErr:         false,
		},
		{
			Name:            "Get Empty list",
			ContentEncoding: meilisearch.BrotliEncoding,
			WantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			c := setup(t, "")
			if !tt.ContentEncoding.IsZero() {
				c = setup(t, "", meilisearch.WithContentEncoding(tt.ContentEncoding, meilisearch.DefaultCompression))
			}

			t.Cleanup(cleanup(c))

			for _, idx := range tt.Indexes {
				info, err := c.CreateIndex(&meilisearch.IndexConfig{
					Uid:        idx,
					PrimaryKey: "id", // Adding a default primary meilisearch.Key
				})
				if tt.WantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					require.NotNil(t, info)
					taskInfo, err := c.WaitForTask(info.TaskUID, 0)
					require.NoError(t, err)
					require.Equal(t, taskInfo.Status, meilisearch.TaskStatusSucceeded)
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
		ContentEncoding meilisearch.ContentEncoding
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
			Name:    "Got meilisearch.Error on delete index",
			WantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			c := setup(t, "")
			if !tt.ContentEncoding.IsZero() {
				c = setup(t, "", meilisearch.WithContentEncoding(tt.ContentEncoding, meilisearch.DefaultCompression))
			}

			t.Cleanup(cleanup(c))

			if len(tt.IndexUID) != 0 {
				info, err := c.CreateIndex(&meilisearch.IndexConfig{
					Uid: tt.IndexUID,
				})
				if tt.WantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					taskInfo, err := c.WaitForTask(info.TaskUID, 0)
					require.NoError(t, err)
					require.Equal(t, taskInfo.Status, meilisearch.TaskStatusSucceeded)
				}
			}

			info, err := c.DeleteIndex(tt.IndexUID)
			if tt.WantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, info)
				taskInfo, err := c.WaitForTask(info.TaskUID, 0)
				require.NoError(t, err)
				require.Equal(t, taskInfo.Status, meilisearch.TaskStatusSucceeded)
			}

		})
	}
}

func Test_CreateSnapshot(t *testing.T) {
	c := setup(t, "")
	taskInfo, err := c.CreateSnapshot()
	require.NoError(t, err)
	testWaitForTask(t, c.Index("indexUID"), taskInfo)
}

func TestGetServiceManagerAndReaders(t *testing.T) {
	c := setup(t, "")
	require.NotNil(t, c.ServiceReader())
	require.NotNil(t, c.TaskManager())
	require.NotNil(t, c.TaskReader())
	require.NotNil(t, c.KeyManager())
	require.NotNil(t, c.KeyReader())
}

func TestGetBatch(t *testing.T) {
	c := setup(t, "")
	indexUID := "indexUID"

	info, err := c.CreateIndex(&meilisearch.IndexConfig{
		Uid: indexUID,
	})

	require.NoError(t, err)
	taskInfo, err := c.WaitForTask(info.TaskUID, 0)
	require.NoError(t, err)
	require.Equal(t, taskInfo.Status, meilisearch.TaskStatusSucceeded)

	info, err = c.DeleteIndex(indexUID)
	require.NoError(t, err)
	taskInfo, err = c.WaitForTask(info.TaskUID, 0)
	require.NoError(t, err)
	require.Equal(t, taskInfo.Status, meilisearch.TaskStatusSucceeded)

	batches, err := c.GetBatches(nil)
	require.NoError(t, err)
	require.NotEmpty(t, batches.Results)

	for _, bt := range batches.Results {
		batch, err := c.GetBatch(bt.UID)
		require.NoError(t, err)

		require.NotZero(t, batch.StartedAt)
		if batch.Progress != nil {
			require.GreaterOrEqual(t, batch.Progress.Percentage, 0.0)
		}
		if batch.Stats != nil {
			require.GreaterOrEqual(t, batch.Stats.TotalNbTasks, 0)
		}
	}
}

func TestGetBatches(t *testing.T) {
	c := setup(t, "")
	indexUID := "indexUID"

	info, err := c.CreateIndex(&meilisearch.IndexConfig{
		Uid: indexUID,
	})

	require.NoError(t, err)
	taskInfo, err := c.WaitForTask(info.TaskUID, 0)
	require.NoError(t, err)
	require.Equal(t, taskInfo.Status, meilisearch.TaskStatusSucceeded)

	info, err = c.DeleteIndex(indexUID)
	require.NoError(t, err)
	taskInfo, err = c.WaitForTask(info.TaskUID, 0)
	require.NoError(t, err)
	require.Equal(t, taskInfo.Status, meilisearch.TaskStatusSucceeded)

	tests := []struct {
		name   string
		params *meilisearch.BatchesQuery
		limit  int
	}{
		{
			name:   "TestGetBatchesWithNoFilter",
			params: nil,
			limit:  -1, // No limit
		},
		{
			name:   "TestGetBatchesWithLimit",
			params: &meilisearch.BatchesQuery{Limit: 1},
			limit:  1,
		},
		{
			name: "TestGetBatchesWithSpecificTypes",
			params: &meilisearch.BatchesQuery{Types: []string{
				string(meilisearch.TaskTypeIndexCreation),
				string(meilisearch.TaskTypeIndexDeletion),
			}},
			limit: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batches, err := c.GetBatches(tt.params)
			require.NoError(t, err)
			require.NotEmpty(t, batches.Results)

			if tt.limit == -1 {
				require.Greater(t, len(batches.Results), 1)
			} else {
				require.LessOrEqual(t, len(batches.Results), tt.limit)
			}

			batch := batches.Results[0]
			require.NotZero(t, batch.StartedAt)
			require.NotEmpty(t, batch.BatchStrategy)
			require.NotNil(t, batch.Stats)
			if tt.params != nil && tt.params.Limit > 0 {
				require.LessOrEqual(t, int64(len(batches.Results)), tt.params.Limit)
			}
		})
	}

}
