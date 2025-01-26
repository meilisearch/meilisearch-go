package meilisearch

import (
	"context"
	"crypto/tls"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

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
		{
			name: "TestTasksWithParams",
			args: args{
				UID:    "indexUID",
				client: sv,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					IndexUIDS: []string{"indexUID"},
					Limit:     10,
					From:      1,
					Statuses:  []TaskStatus{TaskStatusSucceeded},
					Types:     []TaskType{TaskTypeDocumentAdditionOrUpdate},
					Reverse:   true,
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
