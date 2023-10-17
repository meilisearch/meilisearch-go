package meilisearch

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIndex_Delete(t *testing.T) {
	type args struct {
		createUid []string
		deleteUid []string
	}
	tests := []struct {
		name   string
		client *Client
		args   args
	}{
		{
			name:   "TestIndexDeleteOneIndex",
			client: defaultClient,
			args: args{
				createUid: []string{"TestIndexDeleteOneIndex"},
				deleteUid: []string{"TestIndexDeleteOneIndex"},
			},
		},
		{
			name:   "TestIndexDeleteOneIndexWithCustomClient",
			client: customClient,
			args: args{
				createUid: []string{"TestIndexDeleteOneIndexWithCustomClient"},
				deleteUid: []string{"TestIndexDeleteOneIndexWithCustomClient"},
			},
		},
		{
			name:   "TestIndexDeleteMultipleIndex",
			client: defaultClient,
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
			client: defaultClient,
			args: args{
				createUid: []string{},
				deleteUid: []string{"TestIndexDeleteNotExistingIndex"},
			},
		},
		{
			name:   "TestIndexDeleteMultipleNotExistingIndex",
			client: defaultClient,
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
				_, err := SetUpEmptyIndex(&IndexConfig{Uid: uid})
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
	type args struct {
		UID    string
		client *Client
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
				client: defaultClient,
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
				client: customClient,
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
			SetUpBasicIndex(tt.args.UID)
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
	type args struct {
		client *Client
		uid    string
	}
	tests := []struct {
		name string
		args args
		want *Index
	}{
		{
			name: "TestBasicNewIndex",
			args: args{
				client: defaultClient,
				uid:    "TestBasicNewIndex",
			},
			want: &Index{
				UID:    "TestBasicNewIndex",
				client: defaultClient,
			},
		},
		{
			name: "TestNewIndexCustomClient",
			args: args{
				client: customClient,
				uid:    "TestNewIndexCustomClient",
			},
			want: &Index{
				UID:    "TestNewIndexCustomClient",
				client: customClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			t.Cleanup(cleanup(c))

			got := newIndex(c, tt.args.uid)
			require.Equal(t, tt.want.UID, got.UID)
			require.Equal(t, tt.want.client, got.client)
			// Timestamps should be empty unless fetched
			require.Zero(t, got.CreatedAt)
			require.Zero(t, got.UpdatedAt)
		})
	}
}

func TestIndex_GetTask(t *testing.T) {
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
			name: "TestIndexBasicGetTask",
			args: args{
				UID:     "TestIndexBasicGetTask",
				client:  defaultClient,
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
				client:  customClient,
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
				client:  defaultClient,
				taskUID: 0,
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
			name: "TestIndexBasicGetTasks",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
		},
		{
			name: "TestIndexGetTasksWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
		},
		{
			name: "TestIndexBasicGetTasksWithFilters",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
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

			_, err = c.WaitForTask(task.TaskUID)
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
	type args struct {
		UID      string
		client   *Client
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
				client:   defaultClient,
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
				client:   customClient,
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
				client:   defaultClient,
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
				client:   defaultClient,
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

			gotTask, err := i.WaitForTask(task.TaskUID, WaitParams{Context: ctx, Interval: tt.args.interval})
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
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantResp *Index
	}{
		{
			name: "TestIndexBasicFetchInfo",
			args: args{
				UID:    "TestIndexBasicFetchInfo",
				client: defaultClient,
			},
			wantResp: &Index{
				UID:        "TestIndexBasicFetchInfo",
				PrimaryKey: "book_id",
			},
		},
		{
			name: "TestIndexFetchInfoWithCustomClient",
			args: args{
				UID:    "TestIndexFetchInfoWithCustomClient",
				client: customClient,
			},
			wantResp: &Index{
				UID:        "TestIndexFetchInfoWithCustomClient",
				PrimaryKey: "book_id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpBasicIndex(tt.args.UID)
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.FetchInfo()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp.UID, gotResp.UID)
			require.Equal(t, tt.wantResp.PrimaryKey, gotResp.PrimaryKey)
			// Make sure that timestamps are also fetched and are updated
			require.NotZero(t, gotResp.CreatedAt)
			require.NotZero(t, gotResp.UpdatedAt)
			require.Equal(t, i.CreatedAt, gotResp.CreatedAt)
			require.Equal(t, i.UpdatedAt, gotResp.UpdatedAt)
		})
	}
}

func TestIndex_FetchPrimaryKey(t *testing.T) {
	type args struct {
		UID    string
		client *Client
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
				client: defaultClient,
			},
			wantPrimaryKey: "book_id",
		},
		{
			name: "TestIndexFetchPrimaryKeyWithCustomClient",
			args: args{
				UID:    "TestIndexFetchPrimaryKeyWithCustomClient",
				client: customClient,
			},
			wantPrimaryKey: "book_id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpBasicIndex(tt.args.UID)
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
	type args struct {
		primaryKey string
		config     IndexConfig
		client     *Client
	}
	tests := []struct {
		name     string
		args     args
		wantResp *Index
	}{
		{
			name: "TestIndexBasicUpdateIndex",
			args: args{
				client: defaultClient,
				config: IndexConfig{
					Uid: "indexUID",
				},
				primaryKey: "book_id",
			},
			wantResp: &Index{
				UID:        "indexUID",
				PrimaryKey: "book_id",
			},
		},
		{
			name: "TestIndexUpdateIndexWithCustomClient",
			args: args{
				client: customClient,
				config: IndexConfig{
					Uid: "indexUID",
				},
				primaryKey: "book_id",
			},
			wantResp: &Index{
				UID:        "indexUID",
				PrimaryKey: "book_id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			t.Cleanup(cleanup(c))

			i, err := SetUpEmptyIndex(&tt.args.config)
			require.NoError(t, err)
			require.Equal(t, tt.args.config.Uid, i.UID)
			// Store original timestamps
			createdAt := i.CreatedAt
			updatedAt := i.UpdatedAt

			gotResp, err := i.UpdateIndex(tt.args.primaryKey)
			require.NoError(t, err)

			_, err = c.WaitForTask(gotResp.TaskUID)
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
