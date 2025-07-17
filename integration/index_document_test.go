package integration

import (
	"bytes"
	"crypto/tls"
	"github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_AddOrUpdateDocumentsWithContentEncoding(t *testing.T) {
	tests := []struct {
		Name            string
		ContentEncoding meilisearch.ContentEncoding
		Request         []map[string]interface{}
		Response        struct {
			WantResp *meilisearch.TaskInfo
			DocResp  meilisearch.DocumentsResult
		}
	}{
		{
			Name:            "TestIndexBasicAddDocumentsWithGzip",
			ContentEncoding: meilisearch.GzipEncoding,
			Request: []map[string]interface{}{
				{"ID": "123", "Name": "Pride and Prejudice"},
			},
			Response: struct {
				WantResp *meilisearch.TaskInfo
				DocResp  meilisearch.DocumentsResult
			}{
				WantResp: &meilisearch.TaskInfo{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
				},
				DocResp: meilisearch.DocumentsResult{
					Results: meilisearch.Hits{
						{"ID": toRawMessage("123"), "Name": toRawMessage("Pride and Prejudice")},
					},
					Limit:  3,
					Offset: 0,
					Total:  1,
				},
			},
		},
		{
			Name:            "TestIndexBasicAddDocumentsWithDeflate",
			ContentEncoding: meilisearch.DeflateEncoding,
			Request: []map[string]interface{}{
				{"ID": "123", "Name": "Pride and Prejudice"},
			},
			Response: struct {
				WantResp *meilisearch.TaskInfo
				DocResp  meilisearch.DocumentsResult
			}{
				WantResp: &meilisearch.TaskInfo{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
				},
				DocResp: meilisearch.DocumentsResult{
					Results: meilisearch.Hits{
						{"ID": toRawMessage("123"), "Name": toRawMessage("Pride and Prejudice")},
					},
					Limit:  3,
					Offset: 0,
					Total:  1,
				},
			},
		},
		{
			Name:            "TestIndexBasicAddDocumentsWithBrotli",
			ContentEncoding: meilisearch.BrotliEncoding,
			Request: []map[string]interface{}{
				{"ID": "123", "Name": "Pride and Prejudice"},
			},
			Response: struct {
				WantResp *meilisearch.TaskInfo
				DocResp  meilisearch.DocumentsResult
			}{
				WantResp: &meilisearch.TaskInfo{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
				},
				DocResp: meilisearch.DocumentsResult{
					Results: meilisearch.Hits{
						{"ID": toRawMessage("123"), "Name": toRawMessage("Pride and Prejudice")},
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
			sv := setup(t, "", meilisearch.WithContentEncoding(tt.ContentEncoding, meilisearch.DefaultCompression))
			t.Cleanup(cleanup(sv))

			i := sv.Index("indexUID")

			// Add Documents
			gotResp, err := i.AddDocuments(tt.Request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.Response.WantResp.TaskUID)
			require.Equal(t, gotResp.Status, tt.Response.WantResp.Status)
			require.Equal(t, gotResp.Type, tt.Response.WantResp.Type)
			require.Equal(t, gotResp.IndexUID, "indexUID")
			require.NotZero(t, gotResp.EnqueuedAt)

			testWaitForTask(t, i, gotResp)

			// Get Documents
			var documents meilisearch.DocumentsResult
			err = i.GetDocuments(&meilisearch.DocumentsQuery{Limit: 3}, &documents)
			require.NoError(t, err)
			require.Equal(t, tt.Response.DocResp, documents)

			// Update Documents
			gotResp, err = i.UpdateDocuments(tt.Request)
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

	type args struct {
		UID          string
		client       meilisearch.ServiceManager
		documentsPtr interface{}
	}
	type resp struct {
		wantResp     *meilisearch.TaskInfo
		documentsRes meilisearch.DocumentsResult
	}
	tests := []struct {
		name string
		args args
		resp resp
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
				wantResp: &meilisearch.TaskInfo{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
				},
				documentsRes: meilisearch.DocumentsResult{
					Results: meilisearch.Hits{
						{"ID": toRawMessage("123"), "Name": toRawMessage("Pride and Prejudice")},
					},
					Limit:  3,
					Offset: 0,
					Total:  1,
				},
			},
		},
		{
			name: "TestIndexAddDocumentsWithIntID",
			args: args{
				UID:    "TestIndexBasicAddDocumentsWithIntID",
				client: sv,
				documentsPtr: []map[string]interface{}{
					{"BookID": 123, "Title": "Pride and Prejudice"},
				},
			},
			resp: resp{
				wantResp: &meilisearch.TaskInfo{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
				},
				documentsRes: meilisearch.DocumentsResult{
					Results: meilisearch.Hits{
						{"BookID": toRawMessage(123), "Title": toRawMessage("Pride and Prejudice")},
					},
					Limit:  3,
					Offset: 0,
					Total:  1,
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
			require.Equal(t, tt.resp.wantResp.Status, gotResp.Status)
			require.Equal(t, tt.resp.wantResp.Type, gotResp.Type)
			require.Equal(t, tt.args.UID, gotResp.IndexUID)
			require.NotZero(t, gotResp.EnqueuedAt)

			testWaitForTask(t, i, gotResp)

			var documents meilisearch.DocumentsResult
			err = i.GetDocuments(&meilisearch.DocumentsQuery{Limit: 3}, &documents)
			require.NoError(t, err)
			require.Equal(t, tt.resp.documentsRes, documents)

			gotResp, err = i.UpdateDocuments(tt.args.documentsPtr)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.resp.wantResp.TaskUID)
			require.Equal(t, tt.resp.wantResp.Status, gotResp.Status)
			require.Equal(t, tt.resp.wantResp.Type, gotResp.Type)
			require.Equal(t, tt.args.UID, gotResp.IndexUID)
			require.NotZero(t, gotResp.EnqueuedAt)
		})
	}
}

func TestIndex_AddDocumentsWithPrimaryKey(t *testing.T) {
	sv := setup(t, "")

	type args struct {
		UID          string
		client       meilisearch.ServiceManager
		documentsPtr interface{}
		primaryKey   string
	}
	type resp struct {
		wantResp     *meilisearch.TaskInfo
		documentsRes meilisearch.DocumentsResult
	}
	tests := []struct {
		name string
		args args
		resp resp
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
				wantResp: &meilisearch.TaskInfo{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
				},
				documentsRes: meilisearch.DocumentsResult{
					Results: meilisearch.Hits{
						{"key": toRawMessage("123"), "Name": toRawMessage("Pride and Prejudice")},
					},
					Limit:  3,
					Offset: 0,
					Total:  1,
				},
			},
		},
		{
			name: "TestIndexAddDocumentsWithPrimaryKeyWithIntID",
			args: args{
				UID:    "TestIndexAddDocumentsWithPrimaryKeyWithIntID",
				client: sv,
				documentsPtr: []map[string]interface{}{
					{"key": 123, "Name": "Pride and Prejudice"},
				},
				primaryKey: "key",
			},
			resp: resp{
				wantResp: &meilisearch.TaskInfo{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
				},
				documentsRes: meilisearch.DocumentsResult{
					Results: meilisearch.Hits{
						{"key": toRawMessage(123), "Name": toRawMessage("Pride and Prejudice")},
					},
					Limit:  3,
					Offset: 0,
					Total:  1,
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

			testWaitForTask(t, i, gotResp)

			var documents meilisearch.DocumentsResult
			err = i.GetDocuments(&meilisearch.DocumentsQuery{Limit: 3}, &documents)
			require.NoError(t, err)
			require.Equal(t, tt.resp.documentsRes, documents)
		})
	}
}

func TestIndex_AddOrUpdateDocumentsInBatches(t *testing.T) {
	sv := setup(t, "")

	type argsNoKey struct {
		UID          string
		client       meilisearch.ServiceManager
		documentsPtr interface{}
		batchSize    int
	}

	type argsWithKey struct {
		UID          string
		client       meilisearch.ServiceManager
		documentsPtr interface{}
		batchSize    int
		primaryKey   string
	}

	testsNoKey := []struct {
		name          string
		args          argsNoKey
		wantResp      []meilisearch.TaskInfo
		expectedError meilisearch.Error
	}{
		{
			name: "TestIndexBasicAddDocumentsInBatches",
			args: argsNoKey{
				UID:    "TestIndexBasicAddDocumentsInBatches",
				client: sv,
				documentsPtr: meilisearch.Hits{
					{"ID": toRawMessage("122"), "Name": toRawMessage("Pride and Prejudice")},
					{"ID": toRawMessage("123"), "Name": toRawMessage("Pride and Prejudica")},
					{"ID": toRawMessage("124"), "Name": toRawMessage("Pride and Prejudicb")},
					{"ID": toRawMessage("125"), "Name": toRawMessage("Pride and Prejudicc")},
				},
				batchSize: 2,
			},
			wantResp: []meilisearch.TaskInfo{
				{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
				},
				{
					TaskUID: 1,
					Status:  "enqueued",
					Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
				},
			},
		},
	}

	testsWithKey := []struct {
		name          string
		args          argsWithKey
		wantResp      []meilisearch.TaskInfo
		expectedError meilisearch.Error
	}{
		{
			name: "TestIndexBasicAddDocumentsInBatchesWithKey",
			args: argsWithKey{
				UID:    "TestIndexBasicAddDocumentsInBatchesWithKey",
				client: sv,
				documentsPtr: meilisearch.Hits{
					{"ID": toRawMessage("122"), "Name": toRawMessage("Pride and Prejudice")},
					{"ID": toRawMessage("123"), "Name": toRawMessage("Pride and Prejudica")},
					{"ID": toRawMessage("124"), "Name": toRawMessage("Pride and Prejudicb")},
					{"ID": toRawMessage("125"), "Name": toRawMessage("Pride and Prejudicc")},
				},
				batchSize:  2,
				primaryKey: "ID",
			},
			wantResp: []meilisearch.TaskInfo{
				{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
				},
				{
					TaskUID: 1,
					Status:  "enqueued",
					Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
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

			var documents meilisearch.DocumentsResult
			err = i.GetDocuments(&meilisearch.DocumentsQuery{
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

			var documents meilisearch.DocumentsResult
			err = i.GetDocuments(&meilisearch.DocumentsQuery{
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
		client    meilisearch.ServiceManager
		documents []byte
	}
	type testData struct {
		name     string
		args     args
		wantResp *meilisearch.TaskInfo
	}

	tests := []testData{
		{
			name: "TestBasic",
			args: args{
				UID:       "ndjson",
				client:    sv,
				documents: testNdjsonDocuments,
			},
			wantResp: &meilisearch.TaskInfo{
				TaskUID: 0,
				Status:  "enqueued",
				Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
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
			require.NotEmpty(t, wantDocs, "Parsed NDJSON documents should not be empty")

			var (
				gotResp *meilisearch.TaskInfo
				err     error
			)

			if testReader {
				reader := bytes.NewReader(tt.args.documents)
				gotResp, err = i.AddDocumentsNdjsonFromReader(reader)
			} else {
				gotResp, err = i.AddDocumentsNdjson(tt.args.documents)
			}

			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.wantResp.TaskUID)
			require.Equal(t, tt.wantResp.Status, gotResp.Status)
			require.Equal(t, tt.wantResp.Type, gotResp.Type)
			require.NotZero(t, gotResp.EnqueuedAt)

			testWaitForTask(t, i, gotResp)

			var documents meilisearch.DocumentsResult
			err = i.GetDocuments(&meilisearch.DocumentsQuery{}, &documents)
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
		// ✅ Test both `[]byte` and `io.Reader` methods
		testAddDocumentsNdjson(t, tt, false)
		testAddDocumentsNdjson(t, tt, true)
	}
}

func TestIndex_AddOrUpdateDocumentsCsvInBatches(t *testing.T) {
	sv := setup(t, "")

	type args struct {
		UID       string
		client    meilisearch.ServiceManager
		batchSize int
		documents []byte
	}
	type testData struct {
		name     string
		args     args
		wantResp []meilisearch.TaskInfo
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
			wantResp: []meilisearch.TaskInfo{
				{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
				},
				{
					TaskUID: 1,
					Status:  "enqueued",
					Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
				},
				{
					TaskUID: 2,
					Status:  "enqueued",
					Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
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
				gotResp []meilisearch.TaskInfo
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

			var documents meilisearch.DocumentsResult
			err = i.GetDocuments(&meilisearch.DocumentsQuery{}, &documents)
			require.NoError(t, err)
			require.Equal(t, wantDocs, hitsToStringMaps(documents.Results))

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
		client    meilisearch.ServiceManager
		documents []byte
	}
	type testData struct {
		name     string
		args     args
		wantResp *meilisearch.TaskInfo
	}

	tests := []testData{
		{
			name: "TestIndexBasic",
			args: args{
				UID:       "csv",
				client:    sv,
				documents: testCsvDocuments,
			},
			wantResp: &meilisearch.TaskInfo{
				TaskUID: 0,
				Status:  "enqueued",
				Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
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
				gotResp *meilisearch.TaskInfo
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

			var documents meilisearch.DocumentsResult
			err = i.GetDocuments(&meilisearch.DocumentsQuery{}, &documents)
			require.NoError(t, err)
			require.Equal(t, wantDocs, hitsToStringMaps(documents.Results))
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
		client    meilisearch.ServiceManager
		documents []byte
		options   *meilisearch.CsvDocumentsQuery
	}
	type testData struct {
		name     string
		args     args
		wantResp *meilisearch.TaskInfo
	}

	tests := []testData{
		{
			name: "TestIndexBasicAddDocumentsCsvWithOptions",
			args: args{
				UID:       "csv",
				client:    sv,
				documents: testCsvDocuments,
				options: &meilisearch.CsvDocumentsQuery{
					PrimaryKey:   "id",
					CsvDelimiter: ",",
				},
			},
			wantResp: &meilisearch.TaskInfo{
				TaskUID: 0,
				Status:  "enqueued",
				Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
			},
		},
		{
			name: "TestIndexBasicAddDocumentsCsvWithPrimaryKey",
			args: args{
				UID:       "csv",
				client:    sv,
				documents: testCsvDocuments,
				options: &meilisearch.CsvDocumentsQuery{
					PrimaryKey: "id",
				},
			},
			wantResp: &meilisearch.TaskInfo{
				TaskUID: 0,
				Status:  "enqueued",
				Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
			},
		},
		{
			name: "TestIndexBasicAddDocumentsCsvWithCsvDelimiter",
			args: args{
				UID:       "csv",
				client:    sv,
				documents: testCsvDocuments,
				options: &meilisearch.CsvDocumentsQuery{
					CsvDelimiter: ",",
				},
			},
			wantResp: &meilisearch.TaskInfo{
				TaskUID: 0,
				Status:  "enqueued",
				Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
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
				gotResp *meilisearch.TaskInfo
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

			var documents meilisearch.DocumentsResult
			err = i.GetDocuments(&meilisearch.DocumentsQuery{}, &documents)
			require.NoError(t, err)
			require.Equal(t, wantDocs, hitsToStringMaps(documents.Results))
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
		client    meilisearch.ServiceManager
		batchSize int
		documents []byte
	}
	type testData struct {
		name     string
		args     args
		wantResp []meilisearch.TaskInfo
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
			wantResp: []meilisearch.TaskInfo{
				{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
				},
				{
					TaskUID: 1,
					Status:  "enqueued",
					Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
				},
				{
					TaskUID: 2,
					Status:  "enqueued",
					Type:    meilisearch.TaskTypeDocumentAdditionOrUpdate,
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
				gotResp []meilisearch.TaskInfo
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

			var documents meilisearch.DocumentsResult
			err = i.GetDocuments(&meilisearch.DocumentsQuery{}, &documents)
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
	customSv := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantResp *meilisearch.TaskInfo
	}{
		{
			name: "TestIndexBasicDeleteAllDocuments",
			args: args{
				UID:    "TestIndexBasicDeleteAllDocuments",
				client: sv,
			},
			wantResp: &meilisearch.TaskInfo{
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
			wantResp: &meilisearch.TaskInfo{
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

			var documents meilisearch.DocumentsResult
			err = i.GetDocuments(&meilisearch.DocumentsQuery{Limit: 5}, &documents)
			require.NoError(t, err)
			require.Empty(t, documents.Results)
		})
	}
}

func TestIndex_DeleteOneDocument(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID          string
		PrimaryKey   string
		client       meilisearch.ServiceManager
		identifier   string
		documentsPtr interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantResp *meilisearch.TaskInfo
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
			wantResp: &meilisearch.TaskInfo{
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
			wantResp: &meilisearch.TaskInfo{
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
			wantResp: &meilisearch.TaskInfo{
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
			wantResp: &meilisearch.TaskInfo{
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
			wantResp: &meilisearch.TaskInfo{
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
			wantResp: &meilisearch.TaskInfo{
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
	customSv := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID          string
		client       meilisearch.ServiceManager
		identifier   []string
		documentsPtr []docTest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *meilisearch.TaskInfo
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
			wantResp: &meilisearch.TaskInfo{
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
			wantResp: &meilisearch.TaskInfo{
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
			wantResp: &meilisearch.TaskInfo{
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
			wantResp: &meilisearch.TaskInfo{
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
	customSv := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID            string
		client         meilisearch.ServiceManager
		filterToDelete interface{}
		filterToApply  []string
		documentsPtr   []docTestBooks
	}
	tests := []struct {
		name     string
		args     args
		wantResp *meilisearch.TaskInfo
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
			wantResp: &meilisearch.TaskInfo{
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
			wantResp: &meilisearch.TaskInfo{
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
			wantResp: &meilisearch.TaskInfo{
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
			wantResp: &meilisearch.TaskInfo{
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

			var documents meilisearch.DocumentsResult
			err = i.GetDocuments(&meilisearch.DocumentsQuery{}, &documents)
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

	idx := setupMovieIndex(t, c, "movies")
	t.Cleanup(cleanup(c))

	t.Run("Test Upper Case and Add Sparkles around Movie Titles", func(t *testing.T) {
		task, err := idx.UpdateDocumentsByFunction(&meilisearch.UpdateDocumentByFunctionRequest{
			Filter:   "id > 3000",
			Function: "doc.title = `✨ ${doc.title.to_upper()} ✨`",
		})
		require.NoError(t, err)
		testWaitForTask(t, idx, task)
	})

	t.Run("Test User-defined Context", func(t *testing.T) {
		task, err := idx.UpdateDocumentsByFunction(&meilisearch.UpdateDocumentByFunctionRequest{
			Context: map[string]interface{}{
				"idmax": 50,
			},
			Function: "if doc.id >= context.idmax {\n\t\t    doc = ()\n\t\t  } else {\n\t\t\t  doc.title = `✨ ${doc.title} ✨`\n\t\t\t}",
		})
		require.NoError(t, err)
		testWaitForTask(t, idx, task)
	})
}
