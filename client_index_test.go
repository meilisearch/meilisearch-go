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
					Uid: "TestBasicCreateIndex",
				},
			},
			wantResp: &Index{
				UID: "TestBasicCreateIndex",
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
				MeilisearchApiMessage: meilisearchApiMessage{
					ErrorCode: "invalid_index_uid",
				},
			},
		},
		{
			name:   "TestCreateIndexAlreadyExist",
			client: defaultClient,
			args: args{
				config: IndexConfig{
					Uid: "indexUID",
				},
			},
			wantErr: true,
			expectedError: Error{
				MeilisearchApiMessage: meilisearchApiMessage{
					ErrorCode: "index_already_exists",
				},
			},
		},
		{
			name:   "TestCreateIndexTimeout",
			client: timeoutClient,
			args: args{
				config: IndexConfig{
					Uid: "indexUID",
				},
			},
			wantErr: true,
			expectedError: Error{
				MeilisearchApiMessage: meilisearchApiMessage{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client
			t.Cleanup(cleanup(c))
			SetUpBasicIndex()

			gotResp, err := c.CreateIndex(&tt.args.config)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, tt.expectedError.MeilisearchApiMessage.ErrorCode,
					err.(*Error).MeilisearchApiMessage.ErrorCode)
			} else {
				require.NoError(t, err)
				if assert.NotNil(t, gotResp) {
					require.Equal(t, tt.wantResp.UID, gotResp.UID)
					require.Equal(t, tt.wantResp.PrimaryKey, gotResp.PrimaryKey)
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
				createUid: []string{"1"},
				deleteUid: []string{"1"},
			},
			wantErr: false,
		},
		{
			name:   "TestDeleteIndexWithCustomClient",
			client: customClient,
			args: args{
				createUid: []string{"1"},
				deleteUid: []string{"1"},
			},
			wantErr: false,
		},
		{
			name:   "TestMultipleDeleteIndex",
			client: defaultClient,
			args: args{
				createUid: []string{"2", "3", "4", "5"},
				deleteUid: []string{"2", "3", "4", "5"},
			},
			wantErr: false,
		},
		{
			name:   "TestNotExistingDeleteIndex",
			client: defaultClient,
			args: args{
				deleteUid: []string{"1"},
			},
			wantErr: true,
			expectedError: []Error{
				{
					MeilisearchApiMessage: meilisearchApiMessage{
						ErrorCode: "index_not_found",
					},
				},
			},
		},
		{
			name:   "TestMultipleNotExistingDeleteIndex",
			client: defaultClient,
			args: args{
				deleteUid: []string{"2", "3", "4", "5"},
			},
			wantErr: true,
			expectedError: []Error{
				{
					MeilisearchApiMessage: meilisearchApiMessage{
						ErrorCode: "index_not_found",
					},
				},
				{
					MeilisearchApiMessage: meilisearchApiMessage{
						ErrorCode: "index_not_found",
					},
				},
				{
					MeilisearchApiMessage: meilisearchApiMessage{
						ErrorCode: "index_not_found",
					},
				},
				{
					MeilisearchApiMessage: meilisearchApiMessage{
						ErrorCode: "index_not_found",
					},
				},
			},
		},
		{
			name:   "TestDeleteIndexTimeout",
			client: timeoutClient,
			args: args{
				deleteUid: []string{"1"},
			},
			wantErr: true,
			expectedError: []Error{
				{
					MeilisearchApiMessage: meilisearchApiMessage{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client
			t.Cleanup(cleanup(c))

			for _, uid := range tt.args.createUid {
				_, err := c.CreateIndex(&IndexConfig{Uid: uid})
				require.NoError(t, err, "CreateIndex() in TestDeleteIndex error should be nil")
			}
			for k := range tt.args.deleteUid {
				gotOk, err := c.DeleteIndex(tt.args.deleteUid[k])
				if tt.wantErr {
					require.Error(t, err)
					require.Equal(t, tt.expectedError[k].MeilisearchApiMessage.ErrorCode,
						err.(*Error).MeilisearchApiMessage.ErrorCode)
				} else {
					require.NoError(t, err)
					require.True(t, gotOk)
				}
			}
		})
	}
}

func TestClient_DeleteIndexIfExists(t *testing.T) {
	type args struct {
		createUid []string
		deleteUid []string
	}
	tests := []struct {
		name          string
		client        *Client
		args          args
		wantOk        bool
		wantErr       bool
		expectedError []Error
	}{
		{
			name:   "TestBasicDeleteIndexIfExists",
			client: defaultClient,
			args: args{
				createUid: []string{"1"},
				deleteUid: []string{"1"},
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name:   "TestDeleteIndexIfExistsWithCustomClient",
			client: customClient,
			args: args{
				createUid: []string{"1"},
				deleteUid: []string{"1"},
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name:   "TestMultipleDeleteIndexIfExists",
			client: defaultClient,
			args: args{
				createUid: []string{"2", "3", "4", "5"},
				deleteUid: []string{"2", "3", "4", "5"},
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name:   "TestNotExistingDeleteIndexIfExists",
			client: defaultClient,
			args: args{
				deleteUid: []string{"1"},
			},
			wantOk:  false,
			wantErr: false,
		},
		{
			name:   "TestMultipleNotExistingDeleteIndexIfExists",
			client: defaultClient,
			args: args{
				deleteUid: []string{"2", "3", "4", "5"},
			},
			wantOk:  false,
			wantErr: false,
		},
		{
			name:   "TestDeleteIndexIfExistsTimeout",
			client: timeoutClient,
			args: args{
				deleteUid: []string{"1"},
			},
			wantOk:  false,
			wantErr: true,
			expectedError: []Error{
				{
					MeilisearchApiMessage: meilisearchApiMessage{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client
			t.Cleanup(cleanup(c))

			for _, uid := range tt.args.createUid {
				_, err := c.CreateIndex(&IndexConfig{Uid: uid})
				require.NoError(t, err, "CreateIndex() in TestDeleteIndexIfExists error should be nil")
			}
			for k := range tt.args.deleteUid {
				gotOk, err := c.DeleteIndexIfExists(tt.args.deleteUid[k])
				if tt.wantErr {
					require.Error(t, err)
					require.Equal(t, tt.expectedError[k].MeilisearchApiMessage.ErrorCode,
						err.(*Error).MeilisearchApiMessage.ErrorCode)
				} else {
					require.NoError(t, err)
					if tt.wantOk {
						require.True(t, gotOk)
					} else {
						require.False(t, gotOk)
					}
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
			name:   "TestGelAllIndexesOnNoIndexes",
			client: defaultClient,
			args: args{
				uid: []string{},
			},
			wantResp: []Index{},
		},
		{
			name:   "TestBasicGelAllIndexes",
			client: defaultClient,
			args: args{
				uid: []string{"1"},
			},
			wantResp: []Index{
				{
					UID: "1",
				},
			},
		},
		{
			name:   "TestGelAllIndexesWithCustomClient",
			client: customClient,
			args: args{
				uid: []string{"1"},
			},
			wantResp: []Index{
				{
					UID: "1",
				},
			},
		},
		{
			name:   "TestGelAllIndexesOnMultipleIndex",
			client: defaultClient,
			args: args{
				uid: []string{"1", "2", "3"},
			},
			wantResp: []Index{
				{
					UID: "1",
				},
				{
					UID: "2",
				},
				{
					UID: "3",
				},
			},
		},
		{
			name:   "TestGelAllIndexesOnMultipleIndexWithPrimaryKey",
			client: defaultClient,
			args: args{
				uid: []string{"1", "2", "3"},
			},
			wantResp: []Index{
				{
					UID:        "1",
					PrimaryKey: "PrimaryKey1",
				},
				{
					UID:        "2",
					PrimaryKey: "PrimaryKey2",
				},
				{
					UID:        "3",
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
				_, err := c.CreateIndex(&IndexConfig{Uid: uid})
				require.NoError(t, err, "CreateIndex() in TestGetAllIndexes error should be nil")
			}
			gotResp, err := c.GetAllIndexes()
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
		wantErr  bool
		wantCmp  bool
	}{
		{
			name:   "TestBasicGetIndex",
			client: defaultClient,
			args: args{
				config: IndexConfig{
					Uid: "1",
				},
				uid: "1",
			},
			wantResp: &Index{
				UID: "1",
			},
			wantErr: false,
			wantCmp: false,
		},
		{
			name:   "TestGetIndexWithCustomClient",
			client: customClient,
			args: args{
				config: IndexConfig{
					Uid: "1",
				},
				uid: "1",
			},
			wantResp: &Index{
				UID: "1",
			},
			wantErr: false,
			wantCmp: false,
		},
		{
			name:   "TestGetIndexWithPrimaryKey",
			client: defaultClient,
			args: args{
				config: IndexConfig{
					Uid:        "1",
					PrimaryKey: "PrimaryKey",
				},
				uid: "1",
			},
			wantResp: &Index{
				UID:        "1",
				PrimaryKey: "PrimaryKey",
			},
			wantErr: false,
			wantCmp: false,
		},
		{
			name:   "TestGetIndexOnNotExistingIndex",
			client: defaultClient,
			args: args{
				config: IndexConfig{},
				uid:    "1",
			},
			wantResp: nil,
			wantErr:  true,
			wantCmp:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client
			t.Cleanup(cleanup(c))

			gotCreatedResp, err := c.CreateIndex(&tt.args.config)
			if tt.args.config.Uid != "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			gotResp, err := c.GetIndex(tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantResp.UID, gotResp.UID)
				require.Equal(t, gotCreatedResp.UID, gotResp.UID)
				require.Equal(t, tt.wantResp.PrimaryKey, gotResp.PrimaryKey)
				require.Equal(t, gotCreatedResp.PrimaryKey, gotResp.PrimaryKey)
			}
		})
	}
}

func TestClient_GetOrCreateIndex(t *testing.T) {
	type args struct {
		config IndexConfig
	}
	tests := []struct {
		name     string
		client   *Client
		args     args
		wantResp *Index
	}{
		{
			name:   "TestBasicGetOrCreateIndex",
			client: defaultClient,
			args: args{
				config: IndexConfig{
					Uid: "TestBasicGetOrCreateIndex",
				},
			},
			wantResp: &Index{
				UID: "TestBasicGetOrCreateIndex",
			},
		},
		{
			name:   "TestGetOrCreateIndexWithCustomClient",
			client: customClient,
			args: args{
				config: IndexConfig{
					Uid: "TestBasicGetOrCreateIndex",
				},
			},
			wantResp: &Index{
				UID: "TestBasicGetOrCreateIndex",
			},
		},
		{
			name:   "TestGetOrCreateIndexWithPrimaryKey",
			client: defaultClient,
			args: args{
				config: IndexConfig{
					Uid:        "TestGetOrCreateIndexWithPrimaryKey",
					PrimaryKey: "PrimaryKey",
				},
			},
			wantResp: &Index{
				UID:        "TestGetOrCreateIndexWithPrimaryKey",
				PrimaryKey: "PrimaryKey",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client
			t.Cleanup(cleanup(c))

			gotResp, err := c.GetOrCreateIndex(&tt.args.config)
			require.NoError(t, err)
			if assert.NotNil(t, gotResp) {
				require.Equal(t, tt.wantResp.UID, gotResp.UID)
				require.Equal(t, tt.wantResp.PrimaryKey, gotResp.PrimaryKey)
			}
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
				uid: "1",
			},
			want: Index{
				UID: "1",
			},
		},
		{
			name:   "TestIndexWithCustomClient",
			client: customClient,
			args: args{
				uid: "1",
			},
			want: Index{
				UID: "1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.client.Index(tt.args.uid)
			require.NotNil(t, got)
			require.Equal(t, tt.want.UID, got.UID)
		})
	}
}
