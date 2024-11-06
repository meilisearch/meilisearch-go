package meilisearch

import (
	"crypto/tls"
	"github.com/stretchr/testify/require"
	"testing"
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

func TestIndexManagerAndReaders(t *testing.T) {
	c := setup(t, "")
	idx := c.Index("indexUID")
	require.NotNil(t, idx)
	require.NotNil(t, idx.GetIndexReader())
	require.NotNil(t, idx.GetTaskReader())
	require.NotNil(t, idx.GetSettingsManager())
	require.NotNil(t, idx.GetSettingsReader())
	require.NotNil(t, idx.GetSearch())
	require.NotNil(t, idx.GetDocumentManager())
	require.NotNil(t, idx.GetDocumentReader())
}
