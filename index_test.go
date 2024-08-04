package meilisearch

import (
	"context"
	"crypto/tls"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestIndex_Delete(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		createUid []string
		deleteUid []string
	}
	tests := []struct {
		name   string
		client ServiceManager
		args   args
	}{
		{
			name:   "TestIndexDeleteOneIndex",
			client: sv,
			args: args{
				createUid: []string{"TestIndexDeleteOneIndex"},
				deleteUid: []string{"TestIndexDeleteOneIndex"},
			},
		},
		{
			name:   "TestIndexDeleteOneIndexWithCustomClient",
			client: customSv,
			args: args{
				createUid: []string{"TestIndexDeleteOneIndexWithCustomClient"},
				deleteUid: []string{"TestIndexDeleteOneIndexWithCustomClient"},
			},
		},
		{
			name:   "TestIndexDeleteMultipleIndex",
			client: sv,
			args: args{
				createUid: []string{
					"TestIndexDeleteMultipleIndex_1",
					"TestIndexDeleteMultipleIndex_2",
					"TestIndexDeleteMultipleIndex_3",
					"TestIndexDeleteMultipleIndex_4",
					"TestIndexDeleteMultipleIndex_5",
				},
				deleteUid: []string{
					"TestIndexDeleteMultipleIndex_1",
					"TestIndexDeleteMultipleIndex_2",
					"TestIndexDeleteMultipleIndex_3",
					"TestIndexDeleteMultipleIndex_4",
					"TestIndexDeleteMultipleIndex_5",
				},
			},
		},
		{
			name:   "TestIndexDeleteNotExistingIndex",
			client: sv,
			args: args{
				createUid: []string{},
				deleteUid: []string{"TestIndexDeleteNotExistingIndex"},
			},
		},
		{
			name:   "TestIndexDeleteMultipleNotExistingIndex",
			client: sv,
			args: args{
				createUid: []string{},
				deleteUid: []string{
					"TestIndexDeleteMultipleNotExistingIndex_1",
					"TestIndexDeleteMultipleNotExistingIndex_2",
					"TestIndexDeleteMultipleNotExistingIndex_3",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client
			t.Cleanup(cleanup(c))

			for _, uid := range tt.args.createUid {
				_, err := setUpEmptyIndex(sv, &IndexConfig{Uid: uid})
				require.NoError(t, err, "CreateIndex() in DeleteTest error should be nil")
			}
			for k := range tt.args.deleteUid {
				i := c.Index(tt.args.deleteUid[k])
				gotResp, err := i.Delete(tt.args.deleteUid[k])
				require.True(t, gotResp)
				require.NoError(t, err)
			}
		})
	}
}

func TestIndex_GetStats(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantResp *StatsIndex
	}{
		{
			name: "TestIndexBasicGetStats",
			args: args{
				UID:    "TestIndexBasicGetStats",
				client: sv,
			},
			wantResp: &StatsIndex{
				NumberOfDocuments: 6,
				IsIndexing:        false,
				FieldDistribution: map[string]int64{"book_id": 6, "title": 6},
			},
		},
		{
			name: "TestIndexGetStatsWithCustomClient",
			args: args{
				UID:    "TestIndexGetStatsWithCustomClient",
				client: customSv,
			},
			wantResp: &StatsIndex{
				NumberOfDocuments: 6,
				IsIndexing:        false,
				FieldDistribution: map[string]int64{"book_id": 6, "title": 6},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpBasicIndex(sv, tt.args.UID)
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetStats()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
		})
	}
}

func Test_newIndex(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		client ServiceManager
		uid    string
	}
	tests := []struct {
		name string
		args args
		want IndexManager
	}{
		{
			name: "TestBasicNewIndex",
			args: args{
				client: sv,
				uid:    "TestBasicNewIndex",
			},
			want: sv.Index("TestBasicNewIndex"),
		},
		{
			name: "TestNewIndexCustomClient",
			args: args{
				client: sv,
				uid:    "TestNewIndexCustomClient",
			},
			want: customSv.Index("TestNewIndexCustomClient"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			t.Cleanup(cleanup(c))

			gotIdx := c.Index(tt.args.uid)

			task, err := c.CreateIndex(&IndexConfig{Uid: tt.args.uid})
			require.NoError(t, err)

			testWaitForTask(t, gotIdx, task)

			gotIdxResult, err := gotIdx.FetchInfo()
			require.NoError(t, err)

			wantIdxResult, err := tt.want.FetchInfo()
			require.NoError(t, err)

			require.Equal(t, gotIdxResult.UID, wantIdxResult.UID)
			// Timestamps should be empty unless fetched
			require.NotZero(t, gotIdxResult.CreatedAt)
			require.NotZero(t, gotIdxResult.UpdatedAt)
		})
	}
}

func TestIndex_GetTask(t *testing.T) {
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
			name: "TestIndexBasicGetTask",
			args: args{
				UID:     "TestIndexBasicGetTask",
				client:  sv,
				taskUID: 0,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
		},
		{
			name: "TestIndexGetTaskWithCustomClient",
			args: args{
				UID:     "TestIndexGetTaskWithCustomClient",
				client:  customSv,
				taskUID: 0,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
		},
		{
			name: "TestIndexGetTask",
			args: args{
				UID:     "TestIndexGetTask",
				client:  sv,
				taskUID: 0,
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

			gotResp, err := i.GetTask(task.TaskUID)
			require.NoError(t, err)
			require.NotNil(t, gotResp)
			require.GreaterOrEqual(t, gotResp.UID, tt.args.taskUID)
			require.Equal(t, gotResp.IndexUID, tt.args.UID)
			require.Equal(t, gotResp.Status, TaskStatusSucceeded)

			// Make sure that timestamps are also retrieved
			require.NotZero(t, gotResp.EnqueuedAt)
			require.NotZero(t, gotResp.StartedAt)
			require.NotZero(t, gotResp.FinishedAt)
		})
	}
}

func TestIndex_GetTasks(t *testing.T) {
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
			name: "TestIndexBasicGetTasks",
			args: args{
				UID:    "indexUID",
				client: sv,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
		},
		{
			name: "TestIndexGetTasksWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customSv,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
		},
		{
			name: "TestIndexBasicGetTasksWithFilters",
			args: args{
				UID:    "indexUID",
				client: sv,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					Statuses: []TaskStatus{TaskStatusSucceeded},
					Types:    []TaskType{TaskTypeDocumentAdditionOrUpdate},
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

			gotResp, err := i.GetTasks(nil)
			require.NoError(t, err)
			require.NotNil(t, (*gotResp).Results[0].Status)
			require.NotZero(t, (*gotResp).Results[0].UID)
			require.NotNil(t, (*gotResp).Results[0].Type)
		})
	}
}

func TestIndex_WaitForTask(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID      string
		client   ServiceManager
		interval time.Duration
		timeout  time.Duration
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
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			task, err := i.AddDocuments(tt.args.document)
			require.NoError(t, err)

			ctx, cancelFunc := context.WithTimeout(context.Background(), tt.args.timeout)
			defer cancelFunc()

			gotTask, err := i.WaitForTaskWithContext(ctx, task.TaskUID, 0)
			if tt.args.timeout < tt.args.interval {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, gotTask.Status)
			}
		})
	}
}

func TestIndex_FetchInfo(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))
	broken := setup(t, "", WithAPIKey("wrong"))

	type args struct {
		UID    string
		client ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantResp *IndexResult
	}{
		{
			name: "TestIndexBasicFetchInfo",
			args: args{
				UID:    "TestIndexBasicFetchInfo",
				client: sv,
			},
			wantResp: &IndexResult{
				UID:        "TestIndexBasicFetchInfo",
				PrimaryKey: "book_id",
			},
		},
		{
			name: "TestIndexFetchInfoWithCustomClient",
			args: args{
				UID:    "TestIndexFetchInfoWithCustomClient",
				client: customSv,
			},
			wantResp: &IndexResult{
				UID:        "TestIndexFetchInfoWithCustomClient",
				PrimaryKey: "book_id",
			},
		},
		{
			name: "TestIndexFetchInfoWithBrokenClient",
			args: args{
				UID:    "TestIndexFetchInfoWithCustomClient",
				client: broken,
			},
			wantResp: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpBasicIndex(sv, tt.args.UID)
			c := tt.args.client
			t.Cleanup(cleanup(c))

			i := c.Index(tt.args.UID)

			gotResp, err := i.FetchInfo()

			if tt.wantResp == nil {
				require.Error(t, err)
				require.Nil(t, gotResp)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantResp.UID, gotResp.UID)
				require.Equal(t, tt.wantResp.PrimaryKey, gotResp.PrimaryKey)
				// Make sure that timestamps are also fetched and are updated
				require.NotZero(t, gotResp.CreatedAt)
				require.NotZero(t, gotResp.UpdatedAt)
			}

		})
	}
}

func TestIndex_FetchPrimaryKey(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
	}
	tests := []struct {
		name           string
		args           args
		wantPrimaryKey string
	}{
		{
			name: "TestIndexBasicFetchPrimaryKey",
			args: args{
				UID:    "TestIndexBasicFetchPrimaryKey",
				client: sv,
			},
			wantPrimaryKey: "book_id",
		},
		{
			name: "TestIndexFetchPrimaryKeyWithCustomClient",
			args: args{
				UID:    "TestIndexFetchPrimaryKeyWithCustomClient",
				client: customSv,
			},
			wantPrimaryKey: "book_id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpBasicIndex(tt.args.client, tt.args.UID)
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotPrimaryKey, err := i.FetchPrimaryKey()
			require.NoError(t, err)
			require.Equal(t, &tt.wantPrimaryKey, gotPrimaryKey)
		})
	}
}

func TestIndex_UpdateIndex(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		primaryKey string
		config     IndexConfig
		client     ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantResp *IndexResult
	}{
		{
			name: "TestIndexBasicUpdateIndex",
			args: args{
				client: sv,
				config: IndexConfig{
					Uid: "indexUID",
				},
				primaryKey: "book_id",
			},
			wantResp: &IndexResult{
				UID:        "indexUID",
				PrimaryKey: "book_id",
			},
		},
		{
			name: "TestIndexUpdateIndexWithCustomClient",
			args: args{
				client: customSv,
				config: IndexConfig{
					Uid: "indexUID",
				},
				primaryKey: "book_id",
			},
			wantResp: &IndexResult{
				UID:        "indexUID",
				PrimaryKey: "book_id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			t.Cleanup(cleanup(c))

			i, err := setUpEmptyIndex(tt.args.client, &tt.args.config)
			require.NoError(t, err)
			require.Equal(t, tt.args.config.Uid, i.UID)
			// Store original timestamps
			createdAt := i.CreatedAt
			updatedAt := i.UpdatedAt

			gotResp, err := i.UpdateIndex(tt.args.primaryKey)
			require.NoError(t, err)

			_, err = c.WaitForTask(gotResp.TaskUID, 0)
			require.NoError(t, err)

			require.NoError(t, err)
			require.Equal(t, tt.wantResp.UID, gotResp.IndexUID)

			gotIndex, err := c.GetIndex(tt.args.config.Uid)
			require.NoError(t, err)
			require.Equal(t, tt.wantResp.PrimaryKey, gotIndex.PrimaryKey)
			// Make sure that timestamps were correctly updated as well
			require.Equal(t, createdAt, gotIndex.CreatedAt)
			require.NotEqual(t, updatedAt, gotIndex.UpdatedAt)
		})
	}
}
