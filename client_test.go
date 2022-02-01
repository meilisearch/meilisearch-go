package meilisearch

import (
	"context"
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

func TestClient_GetAllStats(t *testing.T) {
	tests := []struct {
		name   string
		client *Client
	}{
		{
			name:   "TestGetAllStats",
			client: defaultClient,
		},
		{
			name:   "TestGetAllStatsWithCustomClient",
			client: customClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.client.GetAllStats()
			require.NoError(t, err)
			require.NotNil(t, gotResp, "GetAllStats() should not return nil value")
		})
	}
}

func TestClient_GetKeys(t *testing.T) {
	tests := []struct {
		name   string
		client *Client
	}{
		{
			name:   "TestGetKeys",
			client: defaultClient,
		},
		{
			name:   "TestGetKeysWithCustomClient",
			client: defaultClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.client.GetKeys()
			require.NoError(t, err)
			require.NotNil(t, gotResp, "GetKeys() should not return nil value")
		})
	}
}

func TestClient_Health(t *testing.T) {
	tests := []struct {
		name          string
		client        *Client
		wantResp      *Health
		wantErr       bool
		expectedError Error
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
					Name: "meilsearch-client",
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
					Name: "meilsearch-client",
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
		wantResp *Dump
	}{
		{
			name:   "TestCreateDump",
			client: defaultClient,
			wantResp: &Dump{
				Status: "in_progress",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client

			gotResp, err := c.CreateDump()
			require.NoError(t, err)
			if assert.NotNil(t, gotResp, "CreateDump() should not return nil value") {
				require.Equal(t, tt.wantResp.Status, gotResp.Status, "CreateDump() got response status %v, want: %v", gotResp.Status, tt.wantResp.Status)
			}

			// Waiting for CreateDump() to finished
			for {
				gotResp, _ := c.GetDumpStatus(gotResp.UID)
				if gotResp.Status == "done" {
					break
				}
			}
		})
	}
}

func TestClient_GetDumpStatus(t *testing.T) {
	tests := []struct {
		name     string
		client   *Client
		wantResp []DumpStatus
		wantErr  bool
	}{
		{
			name:     "TestGetDumpStatus",
			client:   defaultClient,
			wantResp: []DumpStatus{DumpStatusInProgress, DumpStatusFailed, DumpStatusDone},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client

			dump, err := c.CreateDump()
			require.NoError(t, err, "CreateDump() in TestGetDumpStatus error should be nil")

			gotResp, err := c.GetDumpStatus(dump.UID)
			require.NoError(t, err)
			require.Contains(t, tt.wantResp, gotResp.Status, "GetDumpStatus() got response status %v", gotResp.Status)
			require.NotEqual(t, "failed", gotResp.Status, "GetDumpStatus() response status should not be failed")
		})
	}
}

func TestClient_GetTask(t *testing.T) {
	type args struct {
		UID      string
		client   *Client
		taskID   int64
		document []docTest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestBasicGetTask",
			args: args{
				UID:    "TestBasicGetTask",
				client: defaultClient,
				taskID: 0,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
		},
		{
			name: "TestGetTaskWithCustomClient",
			args: args{
				UID:    "TestGetTaskWithCustomClient",
				client: customClient,
				taskID: 1,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
		},
		{
			name: "TestGetTask",
			args: args{
				UID:    "TestGetTask",
				client: defaultClient,
				taskID: 2,
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

			_, err = c.WaitForTask(task)
			require.NoError(t, err)

			gotResp, err := c.GetTask(task.UID)
			require.NoError(t, err)
			require.NotNil(t, gotResp)
			require.GreaterOrEqual(t, gotResp.UID, tt.args.taskID)
			require.Equal(t, gotResp.IndexUID, tt.args.UID)
			require.Equal(t, gotResp.Status, TaskStatusSucceeded)

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

			_, err = c.WaitForTask(task)
			require.NoError(t, err)

			gotResp, err := i.GetTasks()
			require.NoError(t, err)
			require.NotNil(t, (*gotResp).Results[0].Status)
			require.NotZero(t, (*gotResp).Results[0].UID)
			require.NotNil(t, (*gotResp).Results[0].Type)
		})
	}
}

func TestClient_DefaultWaitForTask(t *testing.T) {
	type args struct {
		UID      string
		client   *Client
		taskID   *Task
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
				UID:      "TestDefaultWaitForTask",
				client:   defaultClient,
				taskID: &Task{
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
				UID:      "TestDefaultWaitForTaskWithCustomClient",
				client:   customClient,
				taskID: &Task{
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

			gotTask, err := c.WaitForTask(task)
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
		taskID   *Task
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
				taskID: &Task{
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
				taskID: &Task{
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
				taskID: &Task{
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
				taskID: &Task{
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

			gotTask, err := c.WaitForTask(task, waitParams{Context: ctx, Interval: tt.args.interval})
			if tt.args.timeout < tt.args.interval {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, gotTask.Status)
			}
		})
	}
}
