package integration

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/require"
)

func TestIndex_GetTask(t *testing.T) {
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

			task, err := i.AddDocuments(tt.args.document, nil)
			require.NoError(t, err)

			_, err = c.WaitForTask(task.TaskUID, 0)
			require.NoError(t, err)

			gotResp, err := i.GetTask(task.TaskUID)
			require.NoError(t, err)
			require.NotNil(t, gotResp)
			require.GreaterOrEqual(t, gotResp.UID, tt.args.taskUID)
			require.Equal(t, gotResp.IndexUID, tt.args.UID)
			require.Equal(t, gotResp.Status, meilisearch.TaskStatusSucceeded)

			// Make sure that timestamps are also retrieved
			require.NotZero(t, gotResp.EnqueuedAt)
			require.NotZero(t, gotResp.StartedAt)
			require.NotZero(t, gotResp.FinishedAt)
		})
	}
}

func TestIndex_GetTasks(t *testing.T) {
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
				query: &meilisearch.TasksQuery{
					Statuses: []meilisearch.TaskStatus{meilisearch.TaskStatusSucceeded},
					Types:    []meilisearch.TaskType{meilisearch.TaskTypeDocumentAdditionOrUpdate},
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
				query: &meilisearch.TasksQuery{
					IndexUIDS: []string{"indexUID"},
					Limit:     10,
					From:      0,
					Statuses:  []meilisearch.TaskStatus{meilisearch.TaskStatusSucceeded},
					Types:     []meilisearch.TaskType{meilisearch.TaskTypeDocumentAdditionOrUpdate},
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

			task, err := i.AddDocuments(tt.args.document, nil)
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
	customSv := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID      string
		client   meilisearch.ServiceManager
		interval time.Duration
		timeout  time.Duration
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

			task, err := i.AddDocuments(tt.args.document, nil)
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

func TestGetTaskDocuments(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanup(sv))

	// The /tasks/{task_id}/documents route is gated by an experimental
	// feature in Meilisearch v1.13. Enable it before issuing the request.
	_, err := sv.ExperimentalFeatures().SetGetTaskDocumentsRoute(true).Update()
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = sv.ExperimentalFeatures().SetGetTaskDocumentsRoute(false).Update()
	})

	uid := "TestGetTaskDocuments"
	i := sv.Index(uid)

	// The endpoint only returns the document payload while the task is
	// still in `enqueued` or `processing` state. Submit a payload large
	// enough that the task does not finish processing before we issue
	// the GetTaskDocuments call.
	const documentCount = 5000
	documents := make([]docTest, 0, documentCount)
	for n := 0; n < documentCount; n++ {
		documents = append(documents, docTest{
			ID:   strconv.Itoa(n + 1),
			Name: fmt.Sprintf("doc-%d", n+1),
		})
	}

	task, err := i.AddDocuments(documents, nil)
	require.NoError(t, err)

	var docs []docTest
	err = sv.GetTaskDocuments(task.TaskUID, &docs)
	require.NoError(t, err)
	require.NotEmpty(t, docs)

	// Drain the task so cleanup is deterministic.
	_, err = sv.WaitForTask(task.TaskUID, 0)
	require.NoError(t, err)
}
