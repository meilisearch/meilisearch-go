package meilisearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateIndex(t *testing.T) {
	type args struct {
		config IndexConfig
	}
	tests := []struct {
		name          string
		client        *Client
		args          args
		wantResp      *Index
		wantErr       bool
		expectedError Error
	}{
		{
			name:   "TestBasicCreateIndex",
			client: defaultClient,
			args: args{
				config: IndexConfig{
					Uid: "TestBasicCreateIndex",
				},
			},
			wantResp: &Index{
				UID: "TestBasicCreateIndex",
			},
			wantErr: false,
		},
		{
			name:   "TestCreateIndexWithCustomClient",
			client: customClient,
			args: args{
				config: IndexConfig{
					Uid: "TestCreateIndexWithCustomClient",
				},
			},
			wantResp: &Index{
				UID: "TestCreateIndexWithCustomClient",
			},
			wantErr: false,
		},
		{
			name:   "TestCreateIndexWithPrimaryKey",
			client: defaultClient,
			args: args{
				config: IndexConfig{
					Uid:        "TestCreateIndexWithPrimaryKey",
					PrimaryKey: "PrimaryKey",
				},
			},
			wantResp: &Index{
				UID:        "TestCreateIndexWithPrimaryKey",
				PrimaryKey: "PrimaryKey",
			},
			wantErr: false,
		},
		{
			name:   "TestCreateIndexInvalidUid",
			client: defaultClient,
			args: args{
				config: IndexConfig{
					Uid: "TestCreateIndexInvalidUid*",
				},
			},
			wantErr: true,
			expectedError: Error{
				MeilisearchApiError: meilisearchApiError{
					Code: "invalid_index_uid",
				},
			},
		},
		{
			name:   "TestCreateIndexTimeout",
			client: timeoutClient,
			args: args{
				config: IndexConfig{
					Uid: "TestCreateIndexTimeout",
				},
			},
			wantErr: true,
			expectedError: Error{
				MeilisearchApiError: meilisearchApiError{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client
			t.Cleanup(cleanup(c))

			gotResp, err := c.CreateIndex(&tt.args.config)

			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, tt.expectedError.MeilisearchApiError.Code,
					err.(*Error).MeilisearchApiError.Code)
			} else {
				require.NoError(t, err)
				require.Equal(t, gotResp.Type, TaskTypeIndexCreation)
				require.Equal(t, gotResp.Status, TaskStatusEnqueued)
				// Make sure that timestamps are also retrieved
				require.NotZero(t, gotResp.EnqueuedAt)

				_, err := c.WaitForTask(gotResp.TaskUID)
				require.NoError(t, err)

				index, err := c.GetIndex(tt.args.config.Uid)

				require.NoError(t, err)
				if assert.NotNil(t, index) {
					require.Equal(t, tt.wantResp.UID, gotResp.IndexUID)
					require.Equal(t, tt.wantResp.UID, index.UID)
					require.Equal(t, tt.wantResp.PrimaryKey, index.PrimaryKey)
				}
			}
		})
	}
}

func TestClient_DeleteIndex(t *testing.T) {
	type args struct {
		createUid []string
		deleteUid []string
	}
	tests := []struct {
		name          string
		client        *Client
		args          args
		wantErr       bool
		expectedError []Error
	}{
		{
			name:   "TestBasicDeleteIndex",
			client: defaultClient,
			args: args{
				createUid: []string{"TestBasicDeleteIndex"},
				deleteUid: []string{"TestBasicDeleteIndex"},
			},
			wantErr: false,
		},
		{
			name:   "TestDeleteIndexWithCustomClient",
			client: customClient,
			args: args{
				createUid: []string{"TestDeleteIndexWithCustomClient"},
				deleteUid: []string{"TestDeleteIndexWithCustomClient"},
			},
			wantErr: false,
		},
		{
			name:   "TestMultipleDeleteIndex",
			client: defaultClient,
			args: args{
				createUid: []string{
					"TestMultipleDeleteIndex_2",
					"TestMultipleDeleteIndex_3",
					"TestMultipleDeleteIndex_4",
					"TestMultipleDeleteIndex_5",
				},
				deleteUid: []string{
					"TestMultipleDeleteIndex_2",
					"TestMultipleDeleteIndex_3",
					"TestMultipleDeleteIndex_4",
					"TestMultipleDeleteIndex_5",
				},
			},
			wantErr: false,
		},
		{
			name:   "TestNotExistingDeleteIndex",
			client: defaultClient,
			args: args{
				deleteUid: []string{"TestNotExistingDeleteIndex"},
			},
			wantErr: false,
		},
		{
			name:   "TestMultipleNotExistingDeleteIndex",
			client: defaultClient,
			args: args{
				deleteUid: []string{
					"TestMultipleNotExistingDeleteIndex_2",
					"TestMultipleNotExistingDeleteIndex_3",
					"TestMultipleNotExistingDeleteIndex_4",
					"TestMultipleNotExistingDeleteIndex_5",
				},
			},
			wantErr: false,
		},
		{
			name:   "TestDeleteIndexTimeout",
			client: timeoutClient,
			args: args{
				deleteUid: []string{"TestDeleteIndexTimeout"},
			},
			wantErr: true,
			expectedError: []Error{
				{
					MeilisearchApiError: meilisearchApiError{},
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
				require.NoError(t, err, "CreateIndex() in TestDeleteIndex error should be nil")
			}
			for k := range tt.args.deleteUid {
				gotResp, err := c.DeleteIndex(tt.args.deleteUid[k])
				if tt.wantErr {
					require.Error(t, err)
					require.Equal(t, tt.expectedError[k].MeilisearchApiError.Code,
						err.(*Error).MeilisearchApiError.Code)
				} else {
					require.NoError(t, err)
					require.Equal(t, gotResp.Type, TaskTypeIndexDeletion)
					// Make sure that timestamps are also retrieved
					require.NotZero(t, gotResp.EnqueuedAt)
				}
			}
		})
	}
}

func TestClient_GetIndexes(t *testing.T) {
	type args struct {
		uid     []string
		request *IndexesQuery
	}
	tests := []struct {
		name     string
		client   *Client
		args     args
		wantResp *IndexesResults
	}{
		{
			name:   "TestGetIndexesOnNoIndexes",
			client: defaultClient,
			args: args{
				uid:     []string{},
				request: nil,
			},
			wantResp: &IndexesResults{
				Offset: 0,
				Limit:  20,
				Total:  0,
			},
		},
		{
			name:   "TestBasicGetIndexes",
			client: defaultClient,
			args: args{
				uid:     []string{"TestBasicGetIndexes"},
				request: nil,
			},
			wantResp: &IndexesResults{
				Results: []Index{
					{
						UID: "TestBasicGetIndexes",
					},
				},
				Offset: 0,
				Limit:  20,
				Total:  1,
			},
		},
		{
			name:   "TestGetIndexesWithCustomClient",
			client: customClient,
			args: args{
				uid:     []string{"TestGetIndexesWithCustomClient"},
				request: nil,
			},
			wantResp: &IndexesResults{
				Results: []Index{
					{
						UID: "TestGetIndexesWithCustomClient",
					},
				},
				Offset: 0,
				Limit:  20,
				Total:  1,
			},
		},
		{
			name:   "TestGetIndexesOnMultipleIndex",
			client: defaultClient,
			args: args{
				uid: []string{
					"TestGetIndexesOnMultipleIndex_1",
					"TestGetIndexesOnMultipleIndex_2",
					"TestGetIndexesOnMultipleIndex_3",
				},
				request: nil,
			},
			wantResp: &IndexesResults{
				Results: []Index{
					{
						UID: "TestGetIndexesOnMultipleIndex_1",
					},
					{
						UID: "TestGetIndexesOnMultipleIndex_2",
					},
					{
						UID: "TestGetIndexesOnMultipleIndex_3",
					},
				},
				Offset: 0,
				Limit:  20,
				Total:  3,
			},
		},
		{
			name:   "TestGetIndexesOnMultipleIndexWithPrimaryKey",
			client: defaultClient,
			args: args{
				uid: []string{
					"TestGetIndexesOnMultipleIndexWithPrimaryKey_1",
					"TestGetIndexesOnMultipleIndexWithPrimaryKey_2",
					"TestGetIndexesOnMultipleIndexWithPrimaryKey_3",
				},
				request: nil,
			},
			wantResp: &IndexesResults{
				Results: []Index{
					{
						UID:        "TestGetIndexesOnMultipleIndexWithPrimaryKey_1",
						PrimaryKey: "PrimaryKey1",
					},
					{
						UID:        "TestGetIndexesOnMultipleIndexWithPrimaryKey_2",
						PrimaryKey: "PrimaryKey2",
					},
					{
						UID:        "TestGetIndexesOnMultipleIndexWithPrimaryKey_3",
						PrimaryKey: "PrimaryKey3",
					},
				},
				Offset: 0,
				Limit:  20,
				Total:  3,
			},
		},
		{
			name:   "TestGetIndexesWithLimit",
			client: defaultClient,
			args: args{
				uid: []string{
					"TestGetIndexesWithLimit_1",
					"TestGetIndexesWithLimit_2",
					"TestGetIndexesWithLimit_3",
				},
				request: &IndexesQuery{
					Limit: 1,
				},
			},
			wantResp: &IndexesResults{
				Results: []Index{
					{
						UID: "TestGetIndexesWithLimit_1",
					},
				},
				Offset: 0,
				Limit:  1,
				Total:  3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client
			t.Cleanup(cleanup(c))

			for _, uid := range tt.args.uid {
				_, err := SetUpEmptyIndex(&IndexConfig{Uid: uid})
				require.NoError(t, err, "CreateIndex() in TestGetIndexes error should be nil")
			}
			gotResp, err := c.GetIndexes(tt.args.request)
			require.NoError(t, err)
			require.Equal(t, len(tt.wantResp.Results), len(gotResp.Results))
			for i := range gotResp.Results {
				require.Equal(t, tt.wantResp.Results[i].UID, gotResp.Results[i].UID)
			}
			require.Equal(t, tt.wantResp.Limit, gotResp.Limit)
			require.Equal(t, tt.wantResp.Offset, gotResp.Offset)
			require.Equal(t, tt.wantResp.Total, gotResp.Total)
		})
	}
}

func TestClient_GetRawIndexes(t *testing.T) {
	type args struct {
		uid     []string
		request *IndexesQuery
	}
	tests := []struct {
		name     string
		client   *Client
		args     args
		wantResp map[string]interface{}
	}{
		{
			name:   "TestGetRawIndexesOnNoIndexes",
			client: defaultClient,
			args: args{
				uid:     []string{},
				request: nil,
			},
			wantResp: map[string]interface{}{
				"results": []map[string]string{},
				"offset":  float64(0),
				"limit":   float64(20),
				"total":   float64(0),
			},
		},
		{
			name:   "TestBasicGetRawIndexes",
			client: defaultClient,
			args: args{
				uid:     []string{"TestBasicGetRawIndexes"},
				request: nil,
			},
			wantResp: map[string]interface{}{
				"results": []map[string]string{
					{
						"uid": "TestBasicGetRawIndexes",
					},
				},
				"offset": float64(0),
				"limit":  float64(20),
				"total":  float64(1),
			},
		},
		{
			name:   "TestGetRawIndexesWithCustomClient",
			client: customClient,
			args: args{
				uid:     []string{"TestGetRawIndexesWithCustomClient"},
				request: nil,
			},
			wantResp: map[string]interface{}{
				"results": []map[string]string{
					{
						"uid": "TestGetRawIndexesWithCustomClient",
					},
				},
				"offset": float64(0),
				"limit":  float64(20),
				"total":  float64(1),
			},
		},
		{
			name:   "TestGetRawIndexesOnMultipleIndex",
			client: defaultClient,
			args: args{
				uid: []string{
					"TestGetRawIndexesOnMultipleIndex_1",
					"TestGetRawIndexesOnMultipleIndex_2",
					"TestGetRawIndexesOnMultipleIndex_3",
				},
				request: nil,
			},
			wantResp: map[string]interface{}{
				"results": []map[string]string{
					{
						"uid": "TestGetRawIndexesOnMultipleIndex_1",
					},
					{
						"uid": "TestGetRawIndexesOnMultipleIndex_2",
					},
					{
						"uid": "TestGetRawIndexesOnMultipleIndex_3",
					},
				},
				"offset": float64(0),
				"limit":  float64(20),
				"total":  float64(3),
			},
		},
		{
			name:   "TestGetRawIndexesOnMultipleIndexWithPrimaryKey",
			client: defaultClient,
			args: args{
				uid: []string{
					"TestGetRawIndexesOnMultipleIndexWithPrimaryKey_1",
					"TestGetRawIndexesOnMultipleIndexWithPrimaryKey_2",
					"TestGetRawIndexesOnMultipleIndexWithPrimaryKey_3",
				},
				request: nil,
			},
			wantResp: map[string]interface{}{
				"results": []map[string]string{
					{
						"uid":        "TestGetRawIndexesOnMultipleIndex_1",
						"primaryKey": "PrimaryKey1",
					},
					{
						"uid":        "TestGetRawIndexesOnMultipleIndex_2",
						"primaryKey": "PrimaryKey2",
					},
					{
						"uid":        "TestGetRawIndexesOnMultipleIndex_3",
						"primaryKey": "PrimaryKey3",
					},
				},
				"offset": float64(0),
				"limit":  float64(20),
				"total":  float64(3),
			},
		},
		{
			name:   "TestGetRawIndexesWithLimit",
			client: defaultClient,
			args: args{
				uid: []string{
					"TestGetRawIndexesWithLimit_1",
					"TestGetRawIndexesWithLimit_2",
					"TestGetRawIndexesWithLimit_3",
				},
				request: &IndexesQuery{
					Limit: 1,
				},
			},
			wantResp: map[string]interface{}{
				"results": []map[string]interface{}{
					{
						"uid": "TestGetIndexesWithLimit_1",
					},
				},
				"lenResults": 1,
				"offset":     float64(0),
				"limit":      float64(1),
				"total":      float64(3),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client
			t.Cleanup(cleanup(c))

			for _, uid := range tt.args.uid {
				_, err := SetUpEmptyIndex(&IndexConfig{Uid: uid})
				require.NoError(t, err, "CreateIndex() in TestGetRawIndexes error should be nil")
			}
			gotResp, err := c.GetRawIndexes(tt.args.request)

			require.NoError(t, err)
			require.Equal(t, tt.wantResp["limit"], gotResp["limit"])
			require.Equal(t, tt.wantResp["offset"], gotResp["offset"])
			require.Equal(t, tt.wantResp["total"], gotResp["total"])
		})
	}
}

func TestClient_GetIndex(t *testing.T) {
	type args struct {
		config IndexConfig
		uid    string
	}
	tests := []struct {
		name     string
		client   *Client
		args     args
		wantResp *Index
		wantCmp  bool
	}{
		{
			name:   "TestBasicGetIndex",
			client: defaultClient,
			args: args{
				config: IndexConfig{
					Uid: "TestBasicGetIndex",
				},
				uid: "TestBasicGetIndex",
			},
			wantResp: &Index{
				UID: "TestBasicGetIndex",
			},
			wantCmp: false,
		},
		{
			name:   "TestGetIndexWithCustomClient",
			client: customClient,
			args: args{
				config: IndexConfig{
					Uid: "TestGetIndexWithCustomClient",
				},
				uid: "TestGetIndexWithCustomClient",
			},
			wantResp: &Index{
				UID: "TestGetIndexWithCustomClient",
			},
			wantCmp: false,
		},
		{
			name:   "TestGetIndexWithPrimaryKey",
			client: defaultClient,
			args: args{
				config: IndexConfig{
					Uid:        "TestGetIndexWithPrimaryKey",
					PrimaryKey: "PrimaryKey",
				},
				uid: "TestGetIndexWithPrimaryKey",
			},
			wantResp: &Index{
				UID:        "TestGetIndexWithPrimaryKey",
				PrimaryKey: "PrimaryKey",
			},
			wantCmp: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client
			t.Cleanup(cleanup(c))

			gotCreatedResp, err := SetUpEmptyIndex(&tt.args.config)
			if tt.args.config.Uid != "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			gotResp, err := c.GetIndex(tt.args.uid)
			if err != nil {
				t.Errorf("GetIndex() error = %v", err)
				return
			} else {
				require.NoError(t, err)
				require.Equal(t, gotCreatedResp.UID, gotResp.UID)
				require.Equal(t, tt.wantResp.UID, gotResp.UID)
				require.Equal(t, tt.args.config.Uid, gotResp.UID)
				require.Equal(t, tt.wantResp.PrimaryKey, gotResp.PrimaryKey)
				// Make sure that timestamps are also retrieved
				require.NotZero(t, gotResp.CreatedAt)
				require.NotZero(t, gotResp.UpdatedAt)
			}
		})
	}
}

func TestClient_GetRawIndex(t *testing.T) {
	type args struct {
		config IndexConfig
		uid    string
	}
	tests := []struct {
		name     string
		client   *Client
		args     args
		wantResp map[string]interface{}
	}{
		{
			name:   "TestBasicGetRawIndex",
			client: defaultClient,
			args: args{
				config: IndexConfig{
					Uid: "TestBasicGetRawIndex",
				},
				uid: "TestBasicGetRawIndex",
			},
			wantResp: map[string]interface{}{
				"uid": string("TestBasicGetRawIndex"),
			},
		},
		{
			name:   "TestGetRawIndexWithCustomClient",
			client: customClient,
			args: args{
				config: IndexConfig{
					Uid: "TestGetRawIndexWithCustomClient",
				},
				uid: "TestGetRawIndexWithCustomClient",
			},
			wantResp: map[string]interface{}{
				"uid": string("TestGetRawIndexWithCustomClient"),
			},
		},
		{
			name:   "TestGetRawIndexWithPrimaryKey",
			client: defaultClient,
			args: args{
				config: IndexConfig{
					Uid:        "TestGetRawIndexWithPrimaryKey",
					PrimaryKey: "PrimaryKey",
				},
				uid: "TestGetRawIndexWithPrimaryKey",
			},
			wantResp: map[string]interface{}{
				"uid":        string("TestGetRawIndexWithPrimaryKey"),
				"primaryKey": "PrimaryKey",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client
			t.Cleanup(cleanup(c))

			_, err := SetUpEmptyIndex(&tt.args.config)
			require.NoError(t, err)

			gotResp, err := c.GetRawIndex(tt.args.uid)
			if err != nil {
				t.Errorf("GetRawIndex() error = %v", err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantResp["uid"], gotResp["uid"])
			require.Equal(t, tt.wantResp["primaryKey"], gotResp["primaryKey"])
		})
	}
}

func TestClient_Index(t *testing.T) {
	type args struct {
		uid string
	}
	tests := []struct {
		name   string
		client *Client
		args   args
		want   Index
	}{
		{
			name:   "TestBasicIndex",
			client: defaultClient,
			args: args{
				uid: "TestBasicIndex",
			},
			want: Index{
				UID: "TestBasicIndex",
			},
		},
		{
			name:   "TestIndexWithCustomClient",
			client: customClient,
			args: args{
				uid: "TestIndexWithCustomClient",
			},
			want: Index{
				UID: "TestIndexWithCustomClient",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.client.Index(tt.args.uid)
			require.NotNil(t, got)
			require.Equal(t, tt.want.UID, got.UID)
			// Timestamps should be empty unless fetched
			require.Zero(t, got.CreatedAt)
			require.Zero(t, got.UpdatedAt)
		})
	}
}
