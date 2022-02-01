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
				require.Equal(t, gotResp.Type, "indexCreation")
				require.Equal(t, gotResp.Status, TaskStatusEnqueued)
				// Make sure that timestamps are also retrieved
				require.NotZero(t, gotResp.EnqueuedAt)

				_, err := c.WaitForTask(gotResp)
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
					require.Equal(t, gotResp.Type, "indexDeletion")
					// Make sure that timestamps are also retrieved
					require.NotZero(t, gotResp.EnqueuedAt)
				}
			}
		})
	}
}

func TestClient_GetAllIndexes(t *testing.T) {
	type args struct {
		uid []string
	}
	tests := []struct {
		name     string
		client   *Client
		args     args
		wantResp []Index
	}{
		{
			name:   "TestGetAllIndexesOnNoIndexes",
			client: defaultClient,
			args: args{
				uid: []string{},
			},
			wantResp: []Index{},
		},
		{
			name:   "TestBasicGetAllIndexes",
			client: defaultClient,
			args: args{
				uid: []string{"TestBasicGetAllIndexes"},
			},
			wantResp: []Index{
				{
					UID: "TestBasicGetAllIndexes",
				},
			},
		},
		{
			name:   "TestGetAllIndexesWithCustomClient",
			client: customClient,
			args: args{
				uid: []string{"TestGetAllIndexesWithCustomClient"},
			},
			wantResp: []Index{
				{
					UID: "TestGetAllIndexesWithCustomClient",
				},
			},
		},
		{
			name:   "TestGetAllIndexesOnMultipleIndex",
			client: defaultClient,
			args: args{
				uid: []string{
					"TestGetAllIndexesOnMultipleIndex_1",
					"TestGetAllIndexesOnMultipleIndex_2",
					"TestGetAllIndexesOnMultipleIndex_3",
				},
			},
			wantResp: []Index{
				{
					UID: "TestGetAllIndexesOnMultipleIndex_1",
				},
				{
					UID: "TestGetAllIndexesOnMultipleIndex_2",
				},
				{
					UID: "TestGetAllIndexesOnMultipleIndex_3",
				},
			},
		},
		{
			name:   "TestGetAllIndexesOnMultipleIndexWithPrimaryKey",
			client: defaultClient,
			args: args{
				uid: []string{
					"TestGetAllIndexesOnMultipleIndexWithPrimaryKey_1",
					"TestGetAllIndexesOnMultipleIndexWithPrimaryKey_2",
					"TestGetAllIndexesOnMultipleIndexWithPrimaryKey_3",
				},
			},
			wantResp: []Index{
				{
					UID:        "TestGetAllIndexesOnMultipleIndexWithPrimaryKey_1",
					PrimaryKey: "PrimaryKey1",
				},
				{
					UID:        "TestGetAllIndexesOnMultipleIndexWithPrimaryKey_2",
					PrimaryKey: "PrimaryKey2",
				},
				{
					UID:        "TestGetAllIndexesOnMultipleIndexWithPrimaryKey_3",
					PrimaryKey: "PrimaryKey3",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client
			t.Cleanup(cleanup(c))

			for _, uid := range tt.args.uid {
				_, err := SetUpEmptyIndex(&IndexConfig{Uid: uid})
				require.NoError(t, err, "CreateIndex() in TestGetAllIndexes error should be nil")
			}
			gotResp, err := c.GetAllIndexes()
			require.NoError(t, err)
			require.Equal(t, len(tt.wantResp), len(gotResp))
		})
	}
}

func TestClient_GetAllRawIndexes(t *testing.T) {
	type args struct {
		uid []string
	}
	tests := []struct {
		name     string
		client   *Client
		args     args
		wantResp []map[string]interface{}
	}{
		{
			name:   "TestGetAllRawIndexesOnNoIndexes",
			client: defaultClient,
			args: args{
				uid: []string{},
			},
			wantResp: []map[string]interface{}{},
		},
		{
			name:   "TestBasicGetAllRawIndexes",
			client: defaultClient,
			args: args{
				uid: []string{"TestBasicGetAllRawIndexes"},
			},
			wantResp: []map[string]interface{}{
				{
					"uid": "TestBasicGetAllRawIndexes",
				},
			},
		},
		{
			name:   "TestGetAllRawIndexesWithCustomClient",
			client: customClient,
			args: args{
				uid: []string{"TestGetAllRawIndexesWithCustomClient"},
			},
			wantResp: []map[string]interface{}{
				{
					"uid": "TestGetAllRawIndexesWithCustomClient",
				},
			},
		},
		{
			name:   "TestGetAllRawIndexesOnMultipleIndex",
			client: defaultClient,
			args: args{
				uid: []string{
					"TestGetAllRawIndexesOnMultipleIndex_1",
					"TestGetAllRawIndexesOnMultipleIndex_2",
					"TestGetAllRawIndexesOnMultipleIndex_3",
				},
			},
			wantResp: []map[string]interface{}{
				{
					"uid": "TestGetAllRawIndexesOnMultipleIndex_1",
				},
				{
					"uid": "TestGetAllRawIndexesOnMultipleIndex_2",
				},
				{
					"uid": "TestGetAllRawIndexesOnMultipleIndex_3",
				},
			},
		},
		{
			name:   "TestGetAllRawIndexesOnMultipleIndexWithPrimaryKey",
			client: defaultClient,
			args: args{
				uid: []string{
					"TestGetAllRawIndexesOnMultipleIndexWithPrimaryKey_1",
					"TestGetAllRawIndexesOnMultipleIndexWithPrimaryKey_2",
					"TestGetAllRawIndexesOnMultipleIndexWithPrimaryKey_3",
				},
			},
			wantResp: []map[string]interface{}{
				{
					"uid":        "TestGetAllRawIndexesOnMultipleIndexWithPrimaryKey_1",
					"primaryKey": "PrimaryKey1",
				},
				{
					"uid":        "TestGetAllRawIndexesOnMultipleIndexWithPrimaryKey_2",
					"primaryKey": "PrimaryKey2",
				},
				{
					"uid":        "TestGetAllRawIndexesOnMultipleIndexWithPrimaryKey_3",
					"primaryKey": "PrimaryKey3",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client
			t.Cleanup(cleanup(c))

			for _, uid := range tt.args.uid {
				_, err := SetUpEmptyIndex(&IndexConfig{Uid: uid})
				require.NoError(t, err, "CreateIndex() in TestGetAllRawIndexes error should be nil")
			}
			gotResp, err := c.GetAllRawIndexes()
			require.NoError(t, err)
			require.Equal(t, len(tt.wantResp), len(gotResp))
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
