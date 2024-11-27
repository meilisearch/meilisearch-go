package meilisearch

import (
	"bytes"
	"crypto/tls"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_AddOrUpdateDocumentsWithContentEncoding(t *testing.T) {
	tests := []struct {
		Name            string
		ContentEncoding ContentEncoding
		Request         interface{}
		Response        struct {
			WantResp *TaskInfo
			DocResp  DocumentsResult
		}
	}{
		{
			Name:            "TestIndexBasicAddDocumentsWithGzip",
			ContentEncoding: GzipEncoding,
			Request: []map[string]interface{}{
				{"ID": "123", "Name": "Pride and Prejudice"},
			},
			Response: struct {
				WantResp *TaskInfo
				DocResp  DocumentsResult
			}{WantResp: &TaskInfo{
				TaskUID: 0,
				Status:  "enqueued",
				Type:    TaskTypeDocumentAdditionOrUpdate,
			},
				DocResp: DocumentsResult{
					Results: []map[string]interface{}{
						{"ID": "123", "Name": "Pride and Prejudice"},
					},
					Limit:  3,
					Offset: 0,
					Total:  1,
				},
			},
		},
		{
			Name:            "TestIndexBasicAddDocumentsWithDeflate",
			ContentEncoding: DeflateEncoding,
			Request: []map[string]interface{}{
				{"ID": "123", "Name": "Pride and Prejudice"},
			},
			Response: struct {
				WantResp *TaskInfo
				DocResp  DocumentsResult
			}{WantResp: &TaskInfo{
				TaskUID: 0,
				Status:  "enqueued",
				Type:    TaskTypeDocumentAdditionOrUpdate,
			},
				DocResp: DocumentsResult{
					Results: []map[string]interface{}{
						{"ID": "123", "Name": "Pride and Prejudice"},
					},
					Limit:  3,
					Offset: 0,
					Total:  1,
				},
			},
		},
		{
			Name:            "TestIndexBasicAddDocumentsWithBrotli",
			ContentEncoding: BrotliEncoding,
			Request: []map[string]interface{}{
				{"ID": "123", "Name": "Pride and Prejudice"},
			},
			Response: struct {
				WantResp *TaskInfo
				DocResp  DocumentsResult
			}{WantResp: &TaskInfo{
				TaskUID: 0,
				Status:  "enqueued",
				Type:    TaskTypeDocumentAdditionOrUpdate,
			},
				DocResp: DocumentsResult{
					Results: []map[string]interface{}{
						{"ID": "123", "Name": "Pride and Prejudice"},
					},
					Limit:  3,
					Offset: 0,
					Total:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			sv := setup(t, "", WithContentEncoding(tt.ContentEncoding, DefaultCompression))
			t.Cleanup(cleanup(sv))

			i := sv.Index("indexUID")
			gotResp, err := i.AddDocuments(&tt.Request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.Response.WantResp.TaskUID)
			require.Equal(t, gotResp.Status, tt.Response.WantResp.Status)
			require.Equal(t, gotResp.Type, tt.Response.WantResp.Type)
			require.Equal(t, gotResp.IndexUID, "indexUID")
			require.NotZero(t, gotResp.EnqueuedAt)
			require.NoError(t, err)

			testWaitForTask(t, i, gotResp)
			var documents DocumentsResult
			err = i.GetDocuments(&DocumentsQuery{
				Limit: 3,
			}, &documents)
			require.NoError(t, err)
			require.Equal(t, tt.Response.DocResp, documents)

			gotResp, err = i.UpdateDocuments(&tt.Request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.Response.WantResp.TaskUID)
			require.Equal(t, gotResp.Status, tt.Response.WantResp.Status)
			require.Equal(t, gotResp.Type, tt.Response.WantResp.Type)
			require.NotZero(t, gotResp.EnqueuedAt)
		})
	}
}

func TestIndex_AddOrUpdateDocuments(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID          string
		client       ServiceManager
		documentsPtr interface{}
	}
	type resp struct {
		wantResp     *TaskInfo
		documentsRes DocumentsResult
	}
	tests := []struct {
		name          string
		args          args
		resp          resp
		expectedError Error
	}{
		{
			name: "TestIndexBasicAddDocuments",
			args: args{
				UID:    "TestIndexBasicAddDocuments",
				client: sv,
				documentsPtr: []map[string]interface{}{
					{"ID": "123", "Name": "Pride and Prejudice"},
				},
			},
			resp: resp{
				wantResp: &TaskInfo{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    TaskTypeDocumentAdditionOrUpdate,
				},
				documentsRes: DocumentsResult{
					Results: []map[string]interface{}{
						{"ID": "123", "Name": "Pride and Prejudice"},
					},
					Limit:  3,
					Offset: 0,
					Total:  1,
				},
			},
		},
		{
			name: "TestIndexAddDocumentsWithCustomClient",
			args: args{
				UID:    "TestIndexAddDocumentsWithCustomClient",
				client: customSv,
				documentsPtr: []map[string]interface{}{
					{"ID": "123", "Name": "Pride and Prejudice"},
				},
			},
			resp: resp{
				wantResp: &TaskInfo{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    TaskTypeDocumentAdditionOrUpdate,
				},
				documentsRes: DocumentsResult{
					Results: []map[string]interface{}{
						{"ID": "123", "Name": "Pride and Prejudice"},
					},
					Limit:  3,
					Offset: 0,
					Total:  1,
				},
			},
		},
		{
			name: "TestIndexMultipleAddDocuments",
			args: args{
				UID:    "TestIndexMultipleAddDocuments",
				client: sv,
				documentsPtr: []map[string]interface{}{
					{"ID": "1", "Name": "Alice In Wonderland"},
					{"ID": "123", "Name": "Pride and Prejudice"},
					{"ID": "456", "Name": "Le Petit Prince"},
				},
			},
			resp: resp{
				wantResp: &TaskInfo{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    TaskTypeDocumentAdditionOrUpdate,
				},
				documentsRes: DocumentsResult{
					Results: []map[string]interface{}{
						{"ID": "1", "Name": "Alice In Wonderland"},
						{"ID": "123", "Name": "Pride and Prejudice"},
						{"ID": "456", "Name": "Le Petit Prince"},
					},
					Limit:  3,
					Offset: 0,
					Total:  3,
				},
			},
		},
		{
			name: "TestIndexBasicAddDocumentsWithIntID",
			args: args{
				UID:    "TestIndexBasicAddDocumentsWithIntID",
				client: sv,
				documentsPtr: []map[string]interface{}{
					{"BookID": float64(123), "Title": "Pride and Prejudice"},
				},
			},
			resp: resp{
				wantResp: &TaskInfo{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    TaskTypeDocumentAdditionOrUpdate,
				},
				documentsRes: DocumentsResult{
					Results: []map[string]interface{}{
						{"BookID": float64(123), "Title": "Pride and Prejudice"},
					},
					Limit:  3,
					Offset: 0,
					Total:  1,
				},
			},
		},
		{
			name: "TestIndexAddDocumentsWithIntIDWithCustomClient",
			args: args{
				UID:    "TestIndexAddDocumentsWithIntIDWithCustomClient",
				client: customSv,
				documentsPtr: []map[string]interface{}{
					{"BookID": float64(123), "Title": "Pride and Prejudice"},
				},
			},
			resp: resp{
				wantResp: &TaskInfo{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    TaskTypeDocumentAdditionOrUpdate,
				},
				documentsRes: DocumentsResult{
					Results: []map[string]interface{}{
						{"BookID": float64(123), "Title": "Pride and Prejudice"},
					},
					Limit:  3,
					Offset: 0,
					Total:  1,
				},
			},
		},
		{
			name: "TestIndexMultipleAddDocumentsWithIntID",
			args: args{
				UID:    "TestIndexMultipleAddDocumentsWithIntID",
				client: sv,
				documentsPtr: []map[string]interface{}{
					{"BookID": float64(1), "Title": "Alice In Wonderland"},
					{"BookID": float64(123), "Title": "Pride and Prejudice"},
					{"BookID": float64(456), "Title": "Le Petit Prince", "Tag": "Conte"},
				},
			},
			resp: resp{
				wantResp: &TaskInfo{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    TaskTypeDocumentAdditionOrUpdate,
				},
				documentsRes: DocumentsResult{
					Results: []map[string]interface{}{
						{"BookID": float64(1), "Title": "Alice In Wonderland"},
						{"BookID": float64(123), "Title": "Pride and Prejudice"},
						{"BookID": float64(456), "Title": "Le Petit Prince", "Tag": "Conte"},
					},
					Limit:  3,
					Offset: 0,
					Total:  3,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.AddDocuments(tt.args.documentsPtr)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.resp.wantResp.TaskUID)
			require.Equal(t, gotResp.Status, tt.resp.wantResp.Status)
			require.Equal(t, gotResp.Type, tt.resp.wantResp.Type)
			require.Equal(t, gotResp.IndexUID, tt.args.UID)
			require.NotZero(t, gotResp.EnqueuedAt)
			require.NoError(t, err)

			testWaitForTask(t, i, gotResp)
			var documents DocumentsResult
			err = i.GetDocuments(&DocumentsQuery{
				Limit: 3,
			}, &documents)
			require.NoError(t, err)
			require.Equal(t, tt.resp.documentsRes, documents)

			gotResp, err = i.UpdateDocuments(tt.args.documentsPtr)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.resp.wantResp.TaskUID)
			require.Equal(t, gotResp.Status, tt.resp.wantResp.Status)
			require.Equal(t, gotResp.Type, tt.resp.wantResp.Type)
			require.Equal(t, gotResp.IndexUID, tt.args.UID)
			require.NotZero(t, gotResp.EnqueuedAt)
			require.NoError(t, err)
		})
	}
}

func TestIndex_AddDocumentsWithPrimaryKey(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID          string
		client       ServiceManager
		documentsPtr interface{}
		primaryKey   string
	}
	type resp struct {
		wantResp     *TaskInfo
		documentsRes DocumentsResult
	}
	tests := []struct {
		name          string
		args          args
		resp          resp
		expectedError Error
	}{
		{
			name: "TestIndexBasicAddDocumentsWithPrimaryKey",
			args: args{
				UID:    "TestIndexBasicAddDocumentsWithPrimaryKey",
				client: sv,
				documentsPtr: []map[string]interface{}{
					{"key": "123", "Name": "Pride and Prejudice"},
				},
				primaryKey: "key",
			},
			resp: resp{
				wantResp: &TaskInfo{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    TaskTypeDocumentAdditionOrUpdate,
				},
				documentsRes: DocumentsResult{
					Results: []map[string]interface{}{
						{"key": "123", "Name": "Pride and Prejudice"},
					},
					Limit:  3,
					Offset: 0,
					Total:  1,
				},
			},
		},
		{
			name: "TestIndexAddDocumentsWithPrimaryKeyWithCustomClient",
			args: args{
				UID:    "TestIndexAddDocumentsWithPrimaryKeyWithCustomClient",
				client: customSv,
				documentsPtr: []map[string]interface{}{
					{"key": "123", "Name": "Pride and Prejudice"},
				},
				primaryKey: "key",
			},
			resp: resp{
				wantResp: &TaskInfo{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    TaskTypeDocumentAdditionOrUpdate,
				},
				documentsRes: DocumentsResult{
					Results: []map[string]interface{}{
						{"key": "123", "Name": "Pride and Prejudice"},
					},
					Limit:  3,
					Offset: 0,
					Total:  1,
				},
			},
		},
		{
			name: "TestIndexMultipleAddDocumentsWithPrimaryKey",
			args: args{
				UID:    "TestIndexMultipleAddDocumentsWithPrimaryKey",
				client: sv,
				documentsPtr: []map[string]interface{}{
					{"key": "1", "Name": "Alice In Wonderland"},
					{"key": "123", "Name": "Pride and Prejudice"},
					{"key": "456", "Name": "Le Petit Prince"},
				},
				primaryKey: "key",
			},
			resp: resp{
				wantResp: &TaskInfo{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    TaskTypeDocumentAdditionOrUpdate,
				},
				documentsRes: DocumentsResult{
					Results: []map[string]interface{}{
						{"key": "1", "Name": "Alice In Wonderland"},
						{"key": "123", "Name": "Pride and Prejudice"},
						{"key": "456", "Name": "Le Petit Prince"},
					},
					Limit:  3,
					Offset: 0,
					Total:  3,
				},
			},
		},
		{
			name: "TestIndexAddDocumentsWithPrimaryKeyWithIntID",
			args: args{
				UID:    "TestIndexAddDocumentsWithPrimaryKeyWithIntID",
				client: sv,
				documentsPtr: []map[string]interface{}{
					{"key": float64(123), "Name": "Pride and Prejudice"},
				},
				primaryKey: "key",
			},
			resp: resp{
				wantResp: &TaskInfo{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    TaskTypeDocumentAdditionOrUpdate,
				},
				documentsRes: DocumentsResult{
					Results: []map[string]interface{}{
						{"key": float64(123), "Name": "Pride and Prejudice"},
					},
					Limit:  3,
					Offset: 0,
					Total:  1,
				},
			},
		},
		{
			name: "TestIndexMultipleAddDocumentsWithPrimaryKeyWithIntID",
			args: args{
				UID:    "TestIndexMultipleAddDocumentsWithPrimaryKeyWithIntID",
				client: sv,
				documentsPtr: []map[string]interface{}{
					{"key": float64(1), "Name": "Alice In Wonderland"},
					{"key": float64(123), "Name": "Pride and Prejudice"},
					{"key": float64(456), "Name": "Le Petit Prince"},
				},
				primaryKey: "key",
			},
			resp: resp{
				wantResp: &TaskInfo{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    TaskTypeDocumentAdditionOrUpdate,
				},
				documentsRes: DocumentsResult{
					Results: []map[string]interface{}{
						{"key": float64(1), "Name": "Alice In Wonderland"},
						{"key": float64(123), "Name": "Pride and Prejudice"},
						{"key": float64(456), "Name": "Le Petit Prince"},
					},
					Limit:  3,
					Offset: 0,
					Total:  3,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.AddDocuments(tt.args.documentsPtr, tt.args.primaryKey)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.resp.wantResp.TaskUID)
			require.Equal(t, tt.resp.wantResp.Status, gotResp.Status)
			require.Equal(t, tt.resp.wantResp.Type, gotResp.Type)
			require.Equal(t, tt.args.UID, gotResp.IndexUID)
			require.NotZero(t, gotResp.EnqueuedAt)
			require.NoError(t, err)

			testWaitForTask(t, i, gotResp)

			var documents DocumentsResult
			err = i.GetDocuments(&DocumentsQuery{Limit: 3}, &documents)
			require.NoError(t, err)
			require.Equal(t, tt.resp.documentsRes, documents)
		})
	}
}

func TestIndex_AddOrUpdateDocumentsInBatches(t *testing.T) {
	sv := setup(t, "")

	type argsNoKey struct {
		UID          string
		client       ServiceManager
		documentsPtr interface{}
		batchSize    int
	}

	type argsWithKey struct {
		UID          string
		client       ServiceManager
		documentsPtr interface{}
		batchSize    int
		primaryKey   string
	}

	testsNoKey := []struct {
		name          string
		args          argsNoKey
		wantResp      []TaskInfo
		expectedError Error
	}{
		{
			name: "TestIndexBasicAddDocumentsInBatches",
			args: argsNoKey{
				UID:    "TestIndexBasicAddDocumentsInBatches",
				client: sv,
				documentsPtr: []map[string]interface{}{
					{"ID": "122", "Name": "Pride and Prejudice"},
					{"ID": "123", "Name": "Pride and Prejudica"},
					{"ID": "124", "Name": "Pride and Prejudicb"},
					{"ID": "125", "Name": "Pride and Prejudicc"},
				},
				batchSize: 2,
			},
			wantResp: []TaskInfo{
				{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    TaskTypeDocumentAdditionOrUpdate,
				},
				{
					TaskUID: 1,
					Status:  "enqueued",
					Type:    TaskTypeDocumentAdditionOrUpdate,
				},
			},
		},
	}

	testsWithKey := []struct {
		name          string
		args          argsWithKey
		wantResp      []TaskInfo
		expectedError Error
	}{
		{
			name: "TestIndexBasicAddDocumentsInBatchesWithKey",
			args: argsWithKey{
				UID:    "TestIndexBasicAddDocumentsInBatchesWithKey",
				client: sv,
				documentsPtr: []map[string]interface{}{
					{"ID": "122", "Name": "Pride and Prejudice"},
					{"ID": "123", "Name": "Pride and Prejudica"},
					{"ID": "124", "Name": "Pride and Prejudicb"},
					{"ID": "125", "Name": "Pride and Prejudicc"},
				},
				batchSize:  2,
				primaryKey: "ID",
			},
			wantResp: []TaskInfo{
				{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    TaskTypeDocumentAdditionOrUpdate,
				},
				{
					TaskUID: 1,
					Status:  "enqueued",
					Type:    TaskTypeDocumentAdditionOrUpdate,
				},
			},
		},
	}

	for _, tt := range testsNoKey {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.AddDocumentsInBatches(tt.args.documentsPtr, tt.args.batchSize)

			require.NoError(t, err)
			for i := 0; i < 2; i++ {
				require.GreaterOrEqual(t, gotResp[i].TaskUID, tt.wantResp[i].TaskUID)
				require.Equal(t, gotResp[i].Status, tt.wantResp[i].Status)
				require.Equal(t, gotResp[i].Type, tt.wantResp[i].Type)
				require.Equal(t, gotResp[i].IndexUID, tt.args.UID)
				require.NotZero(t, gotResp[i].EnqueuedAt)
			}

			testWaitForBatchTask(t, i, gotResp)

			var documents DocumentsResult
			err = i.GetDocuments(&DocumentsQuery{
				Limit: 4,
			}, &documents)

			require.NoError(t, err)
			require.Equal(t, tt.args.documentsPtr, documents.Results)

			gotResp, err = i.UpdateDocumentsInBatches(tt.args.documentsPtr, tt.args.batchSize)
			require.NoError(t, err)
			for i := 0; i < 2; i++ {
				require.GreaterOrEqual(t, gotResp[i].TaskUID, tt.wantResp[i].TaskUID)
				require.Equal(t, gotResp[i].Status, tt.wantResp[i].Status)
				require.Equal(t, gotResp[i].Type, tt.wantResp[i].Type)
				require.Equal(t, gotResp[i].IndexUID, tt.args.UID)
				require.NotZero(t, gotResp[i].EnqueuedAt)
			}

			testWaitForBatchTask(t, i, gotResp)
		})
	}

	for _, tt := range testsWithKey {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.AddDocumentsInBatches(tt.args.documentsPtr, tt.args.batchSize, tt.args.primaryKey)

			require.NoError(t, err)
			for i := 0; i < 2; i++ {
				require.GreaterOrEqual(t, gotResp[i].TaskUID, tt.wantResp[i].TaskUID)
				require.Equal(t, gotResp[i].Status, tt.wantResp[i].Status)
				require.Equal(t, gotResp[i].Type, tt.wantResp[i].Type)
				require.Equal(t, gotResp[i].IndexUID, tt.args.UID)
				require.NotZero(t, gotResp[i].EnqueuedAt)
			}

			testWaitForBatchTask(t, i, gotResp)

			var documents DocumentsResult
			err = i.GetDocuments(&DocumentsQuery{
				Limit: 4,
			}, &documents)

			require.NoError(t, err)
			require.Equal(t, tt.args.documentsPtr, documents.Results)
		})
	}
}

func TestIndex_AddOrUpdateDocumentsNdjson(t *testing.T) {
	sv := setup(t, "")

	type args struct {
		UID       string
		client    ServiceManager
		documents []byte
	}
	type testData struct {
		name     string
		args     args
		wantResp *TaskInfo
	}

	tests := []testData{
		{
			name: "TestIndexBasic",
			args: args{
				UID:       "ndjson",
				client:    sv,
				documents: testNdjsonDocuments,
			},
			wantResp: &TaskInfo{
				TaskUID: 0,
				Status:  "enqueued",
				Type:    TaskTypeDocumentAdditionOrUpdate,
			},
		},
	}

	testAddDocumentsNdjson := func(t *testing.T, tt testData, testReader bool) {
		name := tt.name + "AddDocumentsNdjson"
		if testReader {
			name += "FromReader"
		}

		uid := tt.args.UID
		if testReader {
			uid += "-reader"
		} else {
			uid += "-string"
		}

		t.Run(name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(uid)
			t.Cleanup(cleanup(c))

			wantDocs := testParseNdjsonDocuments(t, bytes.NewReader(tt.args.documents))

			var (
				gotResp *TaskInfo
				err     error
			)

			if testReader {
				gotResp, err = i.AddDocumentsNdjsonFromReader(bytes.NewReader(tt.args.documents))
			} else {
				gotResp, err = i.AddDocumentsNdjson(tt.args.documents)
			}

			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.wantResp.TaskUID)
			require.Equal(t, tt.wantResp.Status, gotResp.Status)
			require.Equal(t, tt.wantResp.Type, gotResp.Type)
			require.NotZero(t, gotResp.EnqueuedAt)

			testWaitForTask(t, i, gotResp)

			var documents DocumentsResult
			err = i.GetDocuments(&DocumentsQuery{}, &documents)
			require.NoError(t, err)
			require.Equal(t, wantDocs, documents.Results)

			if !testReader {
				gotResp, err = i.UpdateDocumentsNdjson(tt.args.documents)
				require.NoError(t, err)
				require.GreaterOrEqual(t, gotResp.TaskUID, tt.wantResp.TaskUID)
				require.Equal(t, tt.wantResp.Status, gotResp.Status)
				require.Equal(t, tt.wantResp.Type, gotResp.Type)
				require.NotZero(t, gotResp.EnqueuedAt)
				testWaitForTask(t, i, gotResp)
			}
		})
	}

	for _, tt := range tests {
		// Test both the string and io.Reader receiving versions
		testAddDocumentsNdjson(t, tt, false)
		testAddDocumentsNdjson(t, tt, true)
	}
}

func TestIndex_AddOrUpdateDocumentsCsvInBatches(t *testing.T) {
	sv := setup(t, "")

	type args struct {
		UID       string
		client    ServiceManager
		batchSize int
		documents []byte
	}
	type testData struct {
		name     string
		args     args
		wantResp []TaskInfo
	}

	tests := []testData{
		{
			name: "TestIndexBasic",
			args: args{
				UID:       "csvbatch",
				client:    sv,
				batchSize: 2,
				documents: testCsvDocuments,
			},
			wantResp: []TaskInfo{
				{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    TaskTypeDocumentAdditionOrUpdate,
				},
				{
					TaskUID: 1,
					Status:  "enqueued",
					Type:    TaskTypeDocumentAdditionOrUpdate,
				},
				{
					TaskUID: 2,
					Status:  "enqueued",
					Type:    TaskTypeDocumentAdditionOrUpdate,
				},
			},
		},
	}

	testAddDocumentsCsvInBatches := func(t *testing.T, tt testData, testReader bool) {
		name := tt.name + "AddDocumentsCsv"
		if testReader {
			name += "FromReader"
		}
		name += "InBatches"

		uid := tt.args.UID
		if testReader {
			uid += "-reader"
		} else {
			uid += "-string"
		}

		t.Run(name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(uid)
			t.Cleanup(cleanup(c))

			wantDocs := testParseCsvDocuments(t, bytes.NewReader(tt.args.documents))

			var (
				gotResp []TaskInfo
				err     error
			)

			if testReader {
				gotResp, err = i.AddDocumentsCsvFromReaderInBatches(bytes.NewReader(tt.args.documents), tt.args.batchSize, nil)
			} else {
				gotResp, err = i.AddDocumentsCsvInBatches(tt.args.documents, tt.args.batchSize, nil)
			}

			require.NoError(t, err)
			for i := 0; i < 2; i++ {
				require.GreaterOrEqual(t, gotResp[i].TaskUID, tt.wantResp[i].TaskUID)
				require.Equal(t, gotResp[i].Status, tt.wantResp[i].Status)
				require.Equal(t, gotResp[i].Type, tt.wantResp[i].Type)
				require.NotZero(t, gotResp[i].EnqueuedAt)
			}

			testWaitForBatchTask(t, i, gotResp)

			var documents DocumentsResult
			err = i.GetDocuments(&DocumentsQuery{}, &documents)
			require.NoError(t, err)
			require.Equal(t, wantDocs, documents.Results)

			if !testReader {
				gotResp, err = i.UpdateDocumentsCsvInBatches(tt.args.documents, tt.args.batchSize, nil)
				require.NoError(t, err)
				for i := 0; i < 2; i++ {
					require.GreaterOrEqual(t, gotResp[i].TaskUID, tt.wantResp[i].TaskUID)
					require.Equal(t, gotResp[i].Status, tt.wantResp[i].Status)
					require.Equal(t, gotResp[i].Type, tt.wantResp[i].Type)
					require.NotZero(t, gotResp[i].EnqueuedAt)
				}
			}

		})
	}

	for _, tt := range tests {
		// Test both the string and io.Reader receiving versions
		testAddDocumentsCsvInBatches(t, tt, false)
		testAddDocumentsCsvInBatches(t, tt, true)
	}
}

func TestIndex_AddDocumentsCsv(t *testing.T) {
	sv := setup(t, "")

	type args struct {
		UID       string
		client    ServiceManager
		documents []byte
	}
	type testData struct {
		name     string
		args     args
		wantResp *TaskInfo
	}

	tests := []testData{
		{
			name: "TestIndexBasic",
			args: args{
				UID:       "csv",
				client:    sv,
				documents: testCsvDocuments,
			},
			wantResp: &TaskInfo{
				TaskUID: 0,
				Status:  "enqueued",
				Type:    TaskTypeDocumentAdditionOrUpdate,
			},
		},
	}

	testAddDocumentsCsv := func(t *testing.T, tt testData, testReader bool) {
		name := tt.name + "AddDocumentsCsv"
		if testReader {
			name += "FromReader"
		}

		uid := tt.args.UID
		if testReader {
			uid += "-reader"
		} else {
			uid += "-string"
		}

		t.Run(name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(uid)
			t.Cleanup(cleanup(c))

			wantDocs := testParseCsvDocuments(t, bytes.NewReader(tt.args.documents))

			var (
				gotResp *TaskInfo
				err     error
			)

			if testReader {
				gotResp, err = i.AddDocumentsCsvFromReader(bytes.NewReader(tt.args.documents), nil)
			} else {
				gotResp, err = i.AddDocumentsCsv(tt.args.documents, nil)
			}

			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.wantResp.TaskUID)
			require.Equal(t, tt.wantResp.Status, gotResp.Status)
			require.Equal(t, tt.wantResp.Type, gotResp.Type)
			require.NotZero(t, gotResp.EnqueuedAt)

			testWaitForTask(t, i, gotResp)

			var documents DocumentsResult
			err = i.GetDocuments(&DocumentsQuery{}, &documents)
			require.NoError(t, err)
			require.Equal(t, wantDocs, documents.Results)
		})
	}

	for _, tt := range tests {
		// Test both the string and io.Reader receiving versions
		testAddDocumentsCsv(t, tt, false)
		testAddDocumentsCsv(t, tt, true)
	}
}

func TestIndex_AddDocumentsCsvWithOptions(t *testing.T) {
	sv := setup(t, "")

	type args struct {
		UID       string
		client    ServiceManager
		documents []byte
		options   *CsvDocumentsQuery
	}
	type testData struct {
		name     string
		args     args
		wantResp *TaskInfo
	}

	tests := []testData{
		{
			name: "TestIndexBasicAddDocumentsCsvWithOptions",
			args: args{
				UID:       "csv",
				client:    sv,
				documents: testCsvDocuments,
				options: &CsvDocumentsQuery{
					PrimaryKey:   "id",
					CsvDelimiter: ",",
				},
			},
			wantResp: &TaskInfo{
				TaskUID: 0,
				Status:  "enqueued",
				Type:    TaskTypeDocumentAdditionOrUpdate,
			},
		},
		{
			name: "TestIndexBasicAddDocumentsCsvWithPrimaryKey",
			args: args{
				UID:       "csv",
				client:    sv,
				documents: testCsvDocuments,
				options: &CsvDocumentsQuery{
					PrimaryKey: "id",
				},
			},
			wantResp: &TaskInfo{
				TaskUID: 0,
				Status:  "enqueued",
				Type:    TaskTypeDocumentAdditionOrUpdate,
			},
		},
		{
			name: "TestIndexBasicAddDocumentsCsvWithCsvDelimiter",
			args: args{
				UID:       "csv",
				client:    sv,
				documents: testCsvDocuments,
				options: &CsvDocumentsQuery{
					CsvDelimiter: ",",
				},
			},
			wantResp: &TaskInfo{
				TaskUID: 0,
				Status:  "enqueued",
				Type:    TaskTypeDocumentAdditionOrUpdate,
			},
		},
	}

	testAddDocumentsCsv := func(t *testing.T, tt testData, testReader bool) {
		name := tt.name + "AddDocumentsCsv"
		if testReader {
			name += "FromReader"
		}

		uid := tt.args.UID
		if testReader {
			uid += "-reader"
		} else {
			uid += "-string"
		}

		t.Run(name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(uid)
			t.Cleanup(cleanup(c))

			wantDocs := testParseCsvDocuments(t, bytes.NewReader(tt.args.documents))

			var (
				gotResp *TaskInfo
				err     error
			)

			if testReader {
				gotResp, err = i.AddDocumentsCsvFromReader(bytes.NewReader(tt.args.documents), tt.args.options)
			} else {
				gotResp, err = i.AddDocumentsCsv(tt.args.documents, tt.args.options)
			}

			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.wantResp.TaskUID)
			require.Equal(t, tt.wantResp.Status, gotResp.Status)
			require.Equal(t, tt.wantResp.Type, gotResp.Type)
			require.NotZero(t, gotResp.EnqueuedAt)

			testWaitForTask(t, i, gotResp)

			var documents DocumentsResult
			err = i.GetDocuments(&DocumentsQuery{}, &documents)
			require.NoError(t, err)
			require.Equal(t, wantDocs, documents.Results)
		})
	}

	for _, tt := range tests {
		// Test both the string and io.Reader receiving versions
		testAddDocumentsCsv(t, tt, false)
		testAddDocumentsCsv(t, tt, true)
	}
}

func TestIndex_AddOrUpdateDocumentsNdjsonInBatches(t *testing.T) {
	sv := setup(t, "")

	type args struct {
		UID       string
		client    ServiceManager
		batchSize int
		documents []byte
	}
	type testData struct {
		name     string
		args     args
		wantResp []TaskInfo
	}

	tests := []testData{
		{
			name: "TestIndexBasic",
			args: args{
				UID:       "ndjsonbatch",
				client:    sv,
				batchSize: 2,
				documents: testNdjsonDocuments,
			},
			wantResp: []TaskInfo{
				{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    TaskTypeDocumentAdditionOrUpdate,
				},
				{
					TaskUID: 1,
					Status:  "enqueued",
					Type:    TaskTypeDocumentAdditionOrUpdate,
				},
				{
					TaskUID: 2,
					Status:  "enqueued",
					Type:    TaskTypeDocumentAdditionOrUpdate,
				},
			},
		},
	}

	testAddDocumentsNdjsonInBatches := func(t *testing.T, tt testData, testReader bool) {
		name := tt.name + "AddDocumentsNdjson"
		if testReader {
			name += "FromReader"
		}
		name += "InBatches"

		uid := tt.args.UID
		if testReader {
			uid += "-reader"
		} else {
			uid += "-string"
		}

		t.Run(name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(uid)
			t.Cleanup(cleanup(c))

			wantDocs := testParseNdjsonDocuments(t, bytes.NewReader(tt.args.documents))

			var (
				gotResp []TaskInfo
				err     error
			)

			if testReader {
				gotResp, err = i.AddDocumentsNdjsonFromReaderInBatches(bytes.NewReader(tt.args.documents), tt.args.batchSize)
			} else {
				gotResp, err = i.AddDocumentsNdjsonInBatches(tt.args.documents, tt.args.batchSize)
			}

			require.NoError(t, err)
			for i := 0; i < 2; i++ {
				require.GreaterOrEqual(t, gotResp[i].TaskUID, tt.wantResp[i].TaskUID)
				require.Equal(t, gotResp[i].Status, tt.wantResp[i].Status)
				require.Equal(t, gotResp[i].Type, tt.wantResp[i].Type)
				require.NotZero(t, gotResp[i].EnqueuedAt)
			}

			testWaitForBatchTask(t, i, gotResp)

			var documents DocumentsResult
			err = i.GetDocuments(&DocumentsQuery{}, &documents)
			require.NoError(t, err)
			require.Equal(t, wantDocs, documents.Results)

			if !testReader {
				gotResp, err = i.UpdateDocumentsNdjsonInBatches(tt.args.documents, tt.args.batchSize)
				require.NoError(t, err)
				for i := 0; i < 2; i++ {
					require.GreaterOrEqual(t, gotResp[i].TaskUID, tt.wantResp[i].TaskUID)
					require.Equal(t, gotResp[i].Status, tt.wantResp[i].Status)
					require.Equal(t, gotResp[i].Type, tt.wantResp[i].Type)
					require.NotZero(t, gotResp[i].EnqueuedAt)
				}
				testWaitForBatchTask(t, i, gotResp)
			}
		})
	}

	for _, tt := range tests {
		// Test both the string and io.Reader receiving versions
		testAddDocumentsNdjsonInBatches(t, tt, false)
		testAddDocumentsNdjsonInBatches(t, tt, true)
	}
}

func TestIndex_DeleteAllDocuments(t *testing.T) {
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
		wantResp *TaskInfo
	}{
		{
			name: "TestIndexBasicDeleteAllDocuments",
			args: args{
				UID:    "TestIndexBasicDeleteAllDocuments",
				client: sv,
			},
			wantResp: &TaskInfo{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    "documentDeletion",
			},
		},
		{
			name: "TestIndexDeleteAllDocumentsWithCustomClient",
			args: args{
				UID:    "TestIndexDeleteAllDocumentsWithCustomClient",
				client: customSv,
			},
			wantResp: &TaskInfo{
				TaskUID: 2,
				Status:  "enqueued",
				Type:    "documentDeletion",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			setUpBasicIndex(tt.args.client, tt.args.UID)
			gotResp, err := i.DeleteAllDocuments()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.wantResp.TaskUID)
			require.Equal(t, tt.wantResp.Status, gotResp.Status)
			require.Equal(t, tt.wantResp.Type, gotResp.Type)
			require.NotZero(t, gotResp.EnqueuedAt)

			testWaitForTask(t, i, gotResp)

			var documents DocumentsResult
			err = i.GetDocuments(&DocumentsQuery{Limit: 5}, &documents)
			require.NoError(t, err)
			require.Empty(t, documents.Results)
		})
	}
}

func TestIndex_DeleteOneDocument(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID          string
		PrimaryKey   string
		client       ServiceManager
		identifier   string
		documentsPtr interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantResp *TaskInfo
	}{
		{
			name: "TestIndexBasicDeleteOneDocument",
			args: args{
				UID:        "1",
				client:     sv,
				identifier: "123",
				documentsPtr: []map[string]interface{}{
					{"ID": "123", "Name": "Pride and Prejudice"},
				},
			},
			wantResp: &TaskInfo{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    "documentDeletion",
			},
		},
		{
			name: "TestIndexDeleteOneDocumentWithCustomClient",
			args: args{
				UID:        "2",
				client:     customSv,
				identifier: "123",
				documentsPtr: []map[string]interface{}{
					{"ID": "123", "Name": "Pride and Prejudice"},
				},
			},
			wantResp: &TaskInfo{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    "documentDeletion",
			},
		},
		{
			name: "TestIndexDeleteOneDocumentinMultiple",
			args: args{
				UID:        "3",
				client:     sv,
				identifier: "456",
				documentsPtr: []map[string]interface{}{
					{"ID": "123", "Name": "Pride and Prejudice"},
					{"ID": "456", "Name": "Le Petit Prince"},
					{"ID": "1", "Name": "Alice In Wonderland"},
				},
			},
			wantResp: &TaskInfo{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    "documentDeletion",
			},
		},
		{
			name: "TestIndexBasicDeleteOneDocumentWithIntID",
			args: args{
				UID:        "4",
				client:     sv,
				identifier: "123",
				documentsPtr: []map[string]interface{}{
					{"BookID": 123, "Title": "Pride and Prejudice"},
				},
			},
			wantResp: &TaskInfo{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    "documentDeletion",
			},
		},
		{
			name: "TestIndexDeleteOneDocumentWithIntIDWithCustomClient",
			args: args{
				UID:        "5",
				client:     customSv,
				identifier: "123",
				documentsPtr: []map[string]interface{}{
					{"BookID": 123, "Title": "Pride and Prejudice"},
				},
			},
			wantResp: &TaskInfo{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    "documentDeletion",
			},
		},
		{
			name: "TestIndexDeleteOneDocumentWithIntIDinMultiple",
			args: args{
				UID:        "6",
				client:     sv,
				identifier: "456",
				documentsPtr: []map[string]interface{}{
					{"BookID": 123, "Title": "Pride and Prejudice"},
					{"BookID": 456, "Title": "Le Petit Prince"},
					{"BookID": 1, "Title": "Alice In Wonderland"},
				},
			},
			wantResp: &TaskInfo{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    "documentDeletion",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotAddResp, err := i.AddDocuments(tt.args.documentsPtr)
			require.GreaterOrEqual(t, gotAddResp.TaskUID, tt.wantResp.TaskUID)
			require.NoError(t, err)

			testWaitForTask(t, i, gotAddResp)

			gotResp, err := i.DeleteDocument(tt.args.identifier)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.wantResp.TaskUID)
			require.Equal(t, tt.wantResp.Status, gotResp.Status)
			require.Equal(t, tt.wantResp.Type, gotResp.Type)
			require.NotZero(t, gotResp.EnqueuedAt)

			testWaitForTask(t, i, gotResp)

			var document []map[string]interface{}
			err = i.GetDocument(tt.args.identifier, nil, &document)
			require.Error(t, err)
			require.Empty(t, document)
		})
	}
}

func TestIndex_DeleteDocuments(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID          string
		client       ServiceManager
		identifier   []string
		documentsPtr []docTest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *TaskInfo
	}{
		{
			name: "TestIndexBasicDeleteDocuments",
			args: args{
				UID:        "1",
				client:     sv,
				identifier: []string{"123"},
				documentsPtr: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
			wantResp: &TaskInfo{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    "documentDeletion",
			},
		},
		{
			name: "TestIndexDeleteDocumentsWithCustomClient",
			args: args{
				UID:        "2",
				client:     customSv,
				identifier: []string{"123"},
				documentsPtr: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
			wantResp: &TaskInfo{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    "documentDeletion",
			},
		},
		{
			name: "TestIndexDeleteOneDocumentOnMultiple",
			args: args{
				UID:        "3",
				client:     sv,
				identifier: []string{"123"},
				documentsPtr: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			wantResp: &TaskInfo{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    "documentDeletion",
			},
		},
		{
			name: "TestIndexDeleteMultipleDocuments",
			args: args{
				UID:        "4",
				client:     sv,
				identifier: []string{"123", "456", "1"},
				documentsPtr: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			wantResp: &TaskInfo{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    "documentDeletion",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotAddResp, err := i.AddDocuments(tt.args.documentsPtr)
			require.NoError(t, err)

			testWaitForTask(t, i, gotAddResp)

			gotResp, err := i.DeleteDocuments(tt.args.identifier)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.wantResp.TaskUID)
			require.Equal(t, tt.wantResp.Status, gotResp.Status)
			require.Equal(t, tt.wantResp.Type, gotResp.Type)
			require.NotZero(t, gotResp.EnqueuedAt)

			testWaitForTask(t, i, gotResp)

			var document docTest
			for _, identifier := range tt.args.identifier {
				err = i.GetDocument(identifier, nil, &document)
				require.Error(t, err)
				require.Empty(t, document)
			}
		})
	}
}

func TestIndex_DeleteDocumentsByFilter(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID            string
		client         ServiceManager
		filterToDelete interface{}
		filterToApply  []string
		documentsPtr   []docTestBooks
	}
	tests := []struct {
		name     string
		args     args
		wantResp *TaskInfo
	}{
		{
			name: "TestIndexDeleteDocumentsByFilterString",
			args: args{
				UID:            "1",
				client:         sv,
				filterToApply:  []string{"book_id"},
				filterToDelete: "book_id = 123",
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "Pride and Prejudice", Tag: "Romance", Year: 1813},
				},
			},
			wantResp: &TaskInfo{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    "documentDeletion",
			},
		},
		{
			name: "TestIndexDeleteMultipleDocumentsByFilterArrayOfString",
			args: args{
				UID:            "1",
				client:         customSv,
				filterToApply:  []string{"tag"},
				filterToDelete: []string{"tag = 'Epic fantasy'"},
				documentsPtr: []docTestBooks{
					{BookID: 1344, Title: "The Hobbit", Tag: "Epic fantasy", Year: 1937},
					{BookID: 4, Title: "Harry Potter and the Half-Blood Prince", Tag: "Epic fantasy", Year: 2005},
					{BookID: 42, Title: "The Hitchhiker's Guide to the Galaxy", Tag: "Epic fantasy", Year: 1978},
				},
			},
			wantResp: &TaskInfo{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    "documentDeletion",
			},
		},
		{
			name: "TestIndexDeleteMultipleDocumentsAndMultipleFiltersWithArrayOfString",
			args: args{
				UID:            "1",
				client:         customSv,
				filterToApply:  []string{"tag", "year"},
				filterToDelete: []string{"tag = 'Epic fantasy'", "year > 1936"},
				documentsPtr: []docTestBooks{
					{BookID: 1344, Title: "The Hobbit", Tag: "Epic fantasy", Year: 1937},
					{BookID: 4, Title: "Harry Potter and the Half-Blood Prince", Tag: "Epic fantasy", Year: 2005},
					{BookID: 42, Title: "The Hitchhiker's Guide to the Galaxy", Tag: "Epic fantasy", Year: 1978},
				},
			},
			wantResp: &TaskInfo{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    "documentDeletion",
			},
		},
		{
			name: "TestIndexDeleteMultipleDocumentsAndMultipleFiltersWithInterface",
			args: args{
				UID:            "1",
				client:         customSv,
				filterToApply:  []string{"book_id", "tag"},
				filterToDelete: []interface{}{[]string{"tag = 'Epic fantasy'", "book_id = 123"}},
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "Pride and Prejudice", Tag: "Romance", Year: 1813},
					{BookID: 1344, Title: "The Hobbit", Tag: "Epic fantasy", Year: 1937},
					{BookID: 4, Title: "Harry Potter and the Half-Blood Prince", Tag: "Epic fantasy", Year: 2005},
					{BookID: 42, Title: "The Hitchhiker's Guide to the Galaxy", Tag: "Epic fantasy", Year: 1978},
				},
			},
			wantResp: &TaskInfo{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    "documentDeletion",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotAddResp, err := i.AddDocuments(tt.args.documentsPtr)
			require.NoError(t, err)

			testWaitForTask(t, i, gotAddResp)

			if tt.args.filterToApply != nil && len(tt.args.filterToApply) != 0 {
				gotTask, err := i.UpdateFilterableAttributes(&tt.args.filterToApply)
				require.NoError(t, err)
				testWaitForTask(t, i, gotTask)
			}

			gotResp, err := i.DeleteDocumentsByFilter(tt.args.filterToDelete)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.wantResp.TaskUID)
			require.Equal(t, tt.wantResp.Status, gotResp.Status)
			require.Equal(t, tt.wantResp.Type, gotResp.Type)
			require.NotZero(t, gotResp.EnqueuedAt)

			testWaitForTask(t, i, gotResp)

			var documents DocumentsResult
			err = i.GetDocuments(&DocumentsQuery{}, &documents)
			require.NoError(t, err)
			require.Zero(t, len(documents.Results))
		})
	}
}

func TestIndex_UpdateDocumentsByFunction(t *testing.T) {
	c := setup(t, "")

	exp := c.ExperimentalFeatures()
	exp.SetEditDocumentsByFunction(true)
	res, err := exp.Update()
	require.NoError(t, err)
	require.True(t, res.EditDocumentsByFunction)

	idx := setupMovieIndex(t, c)
	t.Cleanup(cleanup(c))

	t.Run("Test Upper Case and Add Sparkles around Movie Titles", func(t *testing.T) {
		task, err := idx.UpdateDocumentsByFunction(&UpdateDocumentByFunctionRequest{
			Filter:   "id > 3000",
			Function: "doc.title = `✨ ${doc.title.to_upper()} ✨`",
		})
		require.NoError(t, err)
		testWaitForTask(t, idx, task)
	})

	t.Run("Test User-defined Context", func(t *testing.T) {
		task, err := idx.UpdateDocumentsByFunction(&UpdateDocumentByFunctionRequest{
			Context: map[string]interface{}{
				"idmax": 50,
			},
			Function: "if doc.id >= context.idmax {\n\t\t    doc = ()\n\t\t  } else {\n\t\t\t  doc.title = `✨ ${doc.title} ✨`\n\t\t\t}",
		})
		require.NoError(t, err)
		testWaitForTask(t, idx, task)
	})
}
