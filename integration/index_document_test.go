package integration

import (
	"bytes"
	"crypto/tls"
	"testing"

	"github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/require"
)

func Test_GetDocumentsByIDs(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanup(sv))

	_, err := sv.CreateIndex(&meilisearch.IndexConfig{
		Uid: "TestGetDocumentsByIDs",
	})

	require.NoError(t, err)
	request := []map[string]interface{}{
		{"ID": "1", "Name": "Pride and Prejudice 1"},
		{"ID": "2", "Name": "Pride and Prejudice 2"},
		{"ID": "3", "Name": "Pride and Prejudice 3"},
	}
	i := sv.Index("TestGetDocumentsByIDs")

	ts, err := i.AddDocuments(request, nil)
	require.NoError(t, err)

	testWaitForIndexTask(t, i, ts)

	var documents meilisearch.DocumentsResult
	err = sv.Index("TestGetDocumentsByIDs").GetDocuments(&meilisearch.DocumentsQuery{Ids: []string{"1", "2", "3"}}, &documents)
	require.NoError(t, err)

	results := meilisearch.Hits{
		{"ID": toRawMessage("1"), "Name": toRawMessage("Pride and Prejudice 1")},
		{"ID": toRawMessage("2"), "Name": toRawMessage("Pride and Prejudice 2")},
		{"ID": toRawMessage("3"), "Name": toRawMessage("Pride and Prejudice 3")},
	}
	require.Equal(t, results, documents.Results)
}

func Test_GetDocumentsWithQuery(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanup(sv))

	indexUID := "TestGetDocumentsWithQuery"
	_, err := sv.CreateIndex(&meilisearch.IndexConfig{
		Uid: indexUID,
	})
	require.NoError(t, err)

	// Add test documents with sortable fields
	testDocuments := []map[string]interface{}{
		{"id": "1", "title": "Alice in Wonderland", "rating": 4.5, "year": 1865},
		{"id": "2", "title": "Pride and Prejudice", "rating": 4.8, "year": 1813},
		{"id": "3", "title": "The Great Gatsby", "rating": 4.2, "year": 1925},
		{"id": "4", "title": "To Kill a Mockingbird", "rating": 4.9, "year": 1960},
	}

	index := sv.Index(indexUID)
	task, err := index.AddDocuments(testDocuments, nil)
	require.NoError(t, err)
	testWaitForIndexTask(t, index, task)

	// Set sortable attributes for sorting tests
	task, err = index.UpdateSortableAttributes(&[]string{"title", "rating", "year"})
	require.NoError(t, err)
	testWaitForIndexTask(t, index, task)

	tests := []struct {
		name          string
		query         *meilisearch.DocumentsQuery
		expectedCount int
		description   string
	}{
		{
			name:          "Get all documents with nil query",
			query:         nil,
			expectedCount: 4,
			description:   "Should return all documents when query is nil",
		},
		{
			name:          "Get all documents with no query",
			query:         &meilisearch.DocumentsQuery{},
			expectedCount: 4,
			description:   "Should return all documents when no filters applied",
		},
		{
			name: "Get documents with limit",
			query: &meilisearch.DocumentsQuery{
				Limit: 2,
			},
			expectedCount: 2,
			description:   "Should return only 2 documents when limit is set",
		},
		{
			name: "Get documents with offset",
			query: &meilisearch.DocumentsQuery{
				Offset: 2,
			},
			expectedCount: 2,
			description:   "Should return 2 documents when offset is 2",
		},
		{
			name: "Get documents with specific fields",
			query: &meilisearch.DocumentsQuery{
				Fields: []string{"id", "title"},
			},
			expectedCount: 4,
			description:   "Should return all documents but only specified fields",
		},
		{
			name: "Get documents with IDs",
			query: &meilisearch.DocumentsQuery{
				Ids: []string{"1", "3"},
			},
			expectedCount: 2,
			description:   "Should return only documents with specified IDs",
		},
		{
			name: "Get documents with Sort parameter",
			query: &meilisearch.DocumentsQuery{
				Sort: []string{"year:asc"},
			},
			expectedCount: 4,
			description:   "Should return all documents sorted by year ascending",
		},
		{
			name: "Get documents with multiple Sort parameters",
			query: &meilisearch.DocumentsQuery{
				Sort: []string{"rating:desc", "year:asc"},
			},
			expectedCount: 4,
			description:   "Should return all documents sorted by rating desc, then year asc",
		},
		{
			name: "Get documents with combined parameters",
			query: &meilisearch.DocumentsQuery{
				Limit:  3,
				Fields: []string{"id", "title", "rating"},
				Sort:   []string{"rating:desc"},
			},
			expectedCount: 3,
			description:   "Should return 3 documents with specified fields sorted by rating",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var documents meilisearch.DocumentsResult
			err := index.GetDocuments(tt.query, &documents)
			require.NoError(t, err, "GetDocuments should not return an error for: %s", tt.description)
			require.Len(t, documents.Results, tt.expectedCount, "Expected %d documents but got %d for: %s", tt.expectedCount, len(documents.Results), tt.description)

			// Additional validation for specific test cases
			switch tt.name {
			case "Get documents with specific fields":
				// Verify only specified fields are returned
				if len(documents.Results) > 0 {
					doc := documents.Results[0]
					require.Contains(t, doc, "id", "Should contain id field")
					require.Contains(t, doc, "title", "Should contain title field")
					require.NotContains(t, doc, "rating", "Should not contain rating field when not requested")
				}

			case "Get documents with Sort parameter":
				// Verify documents are sorted by year ascending (1813, 1865, 1925, 1960)
				if len(documents.Results) >= 4 {
					expectedOrder := []string{"2", "1", "3", "4"} // IDs in order: 1813, 1865, 1925, 1960
					for i, expectedID := range expectedOrder {
						actualID := documents.Results[i]["id"]
						require.Equal(t, toRawMessage(expectedID), actualID, "Document at position %d should have ID %s (sorted by year ascending)", i, expectedID)
					}
				}

			case "Get documents with multiple Sort parameters":
				// Verify documents are sorted by rating desc, then year asc
				// Expected order: 4.9 (ID:4, 1960), 4.8 (ID:2, 1813), 4.5 (ID:1, 1865), 4.2 (ID:3, 1925)
				if len(documents.Results) >= 4 {
					expectedOrder := []string{"4", "2", "1", "3"} // Sorted by rating desc, then year asc
					for i, expectedID := range expectedOrder {
						actualID := documents.Results[i]["id"]
						require.Equal(t, toRawMessage(expectedID), actualID, "Document at position %d should have ID %s (sorted by rating desc, then year asc)", i, expectedID)
					}
				}

			case "Get documents with IDs":
				// Verify only specified IDs are returned
				for _, doc := range documents.Results {
					id := doc["id"]
					require.Contains(t, []interface{}{toRawMessage("1"), toRawMessage("3")}, id, "Should only return documents with IDs 1 or 3")
				}

			case "Get documents with combined parameters":
				// Verify documents are sorted by rating desc and limited to 3, with only specified fields
				if len(documents.Results) >= 3 {
					expectedOrder := []string{"4", "2", "1"} // Top 3 by rating desc: 4.9, 4.8, 4.5
					for i, expectedID := range expectedOrder {
						doc := documents.Results[i]
						actualID := doc["id"]
						require.Equal(t, toRawMessage(expectedID), actualID, "Document at position %d should have ID %s (sorted by rating desc, limited to 3)", i, expectedID)

						// Verify only specified fields are present
						require.Contains(t, doc, "id", "Should contain id field")
						require.Contains(t, doc, "title", "Should contain title field")
						require.Contains(t, doc, "rating", "Should contain rating field")
						require.NotContains(t, doc, "year", "Should not contain year field when not requested")
					}
				}
			}
		})
	}
}

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
			gotResp, err := i.AddDocuments(tt.Request, nil)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.Response.WantResp.TaskUID)
			require.Equal(t, gotResp.Status, tt.Response.WantResp.Status)
			require.Equal(t, gotResp.Type, tt.Response.WantResp.Type)
			require.Equal(t, gotResp.IndexUID, "indexUID")
			require.NotZero(t, gotResp.EnqueuedAt)

			testWaitForIndexTask(t, i, gotResp)

			// Get Documents
			var documents meilisearch.DocumentsResult
			err = i.GetDocuments(&meilisearch.DocumentsQuery{Limit: 3}, &documents)
			require.NoError(t, err)
			require.Equal(t, tt.Response.DocResp, documents)

			// Update Documents
			gotResp, err = i.UpdateDocuments(tt.Request, nil)
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

			gotResp, err := i.AddDocuments(tt.args.documentsPtr, nil)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.resp.wantResp.TaskUID)
			require.Equal(t, tt.resp.wantResp.Status, gotResp.Status)
			require.Equal(t, tt.resp.wantResp.Type, gotResp.Type)
			require.Equal(t, tt.args.UID, gotResp.IndexUID)
			require.NotZero(t, gotResp.EnqueuedAt)

			testWaitForIndexTask(t, i, gotResp)

			var documents meilisearch.DocumentsResult
			err = i.GetDocuments(&meilisearch.DocumentsQuery{Limit: 3}, &documents)
			require.NoError(t, err)
			require.Equal(t, tt.resp.documentsRes, documents)

			gotResp, err = i.UpdateDocuments(tt.args.documentsPtr, nil)
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

			gotResp, err := i.AddDocuments(tt.args.documentsPtr, &meilisearch.DocumentOptions{PrimaryKey: &tt.args.primaryKey})
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.resp.wantResp.TaskUID)
			require.Equal(t, tt.resp.wantResp.Status, gotResp.Status)
			require.Equal(t, tt.resp.wantResp.Type, gotResp.Type)
			require.Equal(t, tt.args.UID, gotResp.IndexUID)
			require.NotZero(t, gotResp.EnqueuedAt)

			testWaitForIndexTask(t, i, gotResp)

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

			gotResp, err := i.AddDocumentsInBatches(tt.args.documentsPtr, tt.args.batchSize, nil)

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

			gotResp, err = i.UpdateDocumentsInBatches(tt.args.documentsPtr, tt.args.batchSize, nil)
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

			gotResp, err := i.AddDocumentsInBatches(tt.args.documentsPtr, tt.args.batchSize, &meilisearch.DocumentOptions{PrimaryKey: &tt.args.primaryKey})

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
				gotResp, err = i.AddDocumentsNdjsonFromReader(reader, nil)
			} else {
				gotResp, err = i.AddDocumentsNdjson(tt.args.documents, nil)
			}

			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.wantResp.TaskUID)
			require.Equal(t, tt.wantResp.Status, gotResp.Status)
			require.Equal(t, tt.wantResp.Type, gotResp.Type)
			require.NotZero(t, gotResp.EnqueuedAt)

			testWaitForIndexTask(t, i, gotResp)

			var documents meilisearch.DocumentsResult
			err = i.GetDocuments(&meilisearch.DocumentsQuery{}, &documents)
			require.NoError(t, err)
			require.Equal(t, wantDocs, documents.Results)

			if !testReader {
				gotResp, err = i.UpdateDocumentsNdjson(tt.args.documents, nil)
				require.NoError(t, err)
				require.GreaterOrEqual(t, gotResp.TaskUID, tt.wantResp.TaskUID)
				require.Equal(t, tt.wantResp.Status, gotResp.Status)
				require.Equal(t, tt.wantResp.Type, gotResp.Type)
				require.NotZero(t, gotResp.EnqueuedAt)
				testWaitForIndexTask(t, i, gotResp)
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

			testWaitForIndexTask(t, i, gotResp)

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

			testWaitForIndexTask(t, i, gotResp)

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
				gotResp, err = i.AddDocumentsNdjsonFromReaderInBatches(bytes.NewReader(tt.args.documents), tt.args.batchSize, nil)
			} else {
				gotResp, err = i.AddDocumentsNdjsonInBatches(tt.args.documents, tt.args.batchSize, nil)
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
				gotResp, err = i.UpdateDocumentsNdjsonInBatches(tt.args.documents, tt.args.batchSize, nil)
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
			gotResp, err := i.DeleteAllDocuments(&meilisearch.DocumentOptions{})
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.wantResp.TaskUID)
			require.Equal(t, tt.wantResp.Status, gotResp.Status)
			require.Equal(t, tt.wantResp.Type, gotResp.Type)
			require.NotZero(t, gotResp.EnqueuedAt)

			testWaitForIndexTask(t, i, gotResp)

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

			gotAddResp, err := i.AddDocuments(tt.args.documentsPtr, nil)
			require.GreaterOrEqual(t, gotAddResp.TaskUID, tt.wantResp.TaskUID)
			require.NoError(t, err)

			testWaitForIndexTask(t, i, gotAddResp)

			gotResp, err := i.DeleteDocument(tt.args.identifier, nil)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.wantResp.TaskUID)
			require.Equal(t, tt.wantResp.Status, gotResp.Status)
			require.Equal(t, tt.wantResp.Type, gotResp.Type)
			require.NotZero(t, gotResp.EnqueuedAt)

			testWaitForIndexTask(t, i, gotResp)

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

			gotAddResp, err := i.AddDocuments(tt.args.documentsPtr, nil)
			require.NoError(t, err)

			testWaitForIndexTask(t, i, gotAddResp)

			gotResp, err := i.DeleteDocuments(tt.args.identifier, nil)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.wantResp.TaskUID)
			require.Equal(t, tt.wantResp.Status, gotResp.Status)
			require.Equal(t, tt.wantResp.Type, gotResp.Type)
			require.NotZero(t, gotResp.EnqueuedAt)

			testWaitForIndexTask(t, i, gotResp)

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
		filterToApply  []interface{}
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
				filterToApply:  []interface{}{"book_id"},
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
				filterToApply:  []interface{}{"tag"},
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
				filterToApply:  []interface{}{"tag", "year"},
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
				filterToApply:  []interface{}{"book_id", "tag"},
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
		{
			name: "TestIndexDeleteWithAttributeRuleForTagAndYear",
			args: args{
				UID:    "1",
				client: customSv,
				filterToApply: []interface{}{
					meilisearch.AttributeRule{
						AttributePatterns: []string{"tag"},
						Features: meilisearch.AttributeFeatures{
							FacetSearch: false,
							Filter: meilisearch.FilterFeatures{
								Equality:   true,
								Comparison: false,
							},
						},
					},
					meilisearch.AttributeRule{
						AttributePatterns: []string{"year"},
						Features: meilisearch.AttributeFeatures{
							FacetSearch: false,
							Filter: meilisearch.FilterFeatures{
								Equality:   true,
								Comparison: true,
							},
						},
					},
				},
				filterToDelete: []string{"tag = 'Fantasy'", "year > 1900"},
				documentsPtr: []docTestBooks{
					{BookID: 1, Title: "Fantasy Realms", Tag: "Fantasy", Year: 1950},
					{BookID: 1344, Title: "The Hobbit", Tag: "Fantasy", Year: 1937},
				},
			},
			wantResp: &meilisearch.TaskInfo{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    "documentDeletion",
			},
		},
		{
			name: "TestIndexDeleteWithMixedFilterableAttributes",
			args: args{
				UID:    "1",
				client: customSv,
				filterToApply: []interface{}{
					"title",
					meilisearch.AttributeRule{
						AttributePatterns: []string{"year"},
						Features: meilisearch.AttributeFeatures{
							FacetSearch: false,
							Filter: meilisearch.FilterFeatures{
								Equality:   true,
								Comparison: true,
							},
						},
					},
				},
				filterToDelete: []string{"title = 'The Hobbit'", "year > 1930"},
				documentsPtr: []docTestBooks{
					{BookID: 1344, Title: "The Hobbit", Tag: "Fantasy", Year: 1937},
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

			gotAddResp, err := i.AddDocuments(tt.args.documentsPtr, nil)
			require.NoError(t, err)

			testWaitForIndexTask(t, i, gotAddResp)

			if len(tt.args.filterToApply) != 0 {
				gotTask, err := i.UpdateFilterableAttributes(&tt.args.filterToApply)
				require.NoError(t, err)
				testWaitForIndexTask(t, i, gotTask)
			}

			gotResp, err := i.DeleteDocumentsByFilter(tt.args.filterToDelete, nil)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.TaskUID, tt.wantResp.TaskUID)
			require.Equal(t, tt.wantResp.Status, gotResp.Status)
			require.Equal(t, tt.wantResp.Type, gotResp.Type)
			require.NotZero(t, gotResp.EnqueuedAt)

			testWaitForIndexTask(t, i, gotResp)

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
		testWaitForIndexTask(t, idx, task)
	})

	t.Run("Test User-defined Context", func(t *testing.T) {
		task, err := idx.UpdateDocumentsByFunction(&meilisearch.UpdateDocumentByFunctionRequest{
			Context: map[string]interface{}{
				"idmax": 50,
			},
			Function: "if doc.id >= context.idmax {\n\t\t    doc = ()\n\t\t  } else {\n\t\t\t  doc.title = `✨ ${doc.title} ✨`\n\t\t\t}",
		})
		require.NoError(t, err)
		testWaitForIndexTask(t, idx, task)
	})
}
func TestIndex_DocumentOperationsWithCustomMetadata(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanup(sv))

	// Setup a basic index
	indexUID := "TestCustomMetadata"
	_, err := sv.CreateIndex(&meilisearch.IndexConfig{Uid: indexUID, PrimaryKey: "id"})
	require.NoError(t, err)
	i := sv.Index(indexUID)

	filterableAttributes := []interface{}{"id"}

	task, err := i.UpdateFilterableAttributes(&filterableAttributes)
	require.NoError(t, err)
	testWaitForIndexTask(t, i, task)

	// Define common test data
	documents := []map[string]interface{}{
		{"id": "1", "title": "Document 1"},
		{"id": "2", "title": "Document 2"},
	}

	tests := []struct {
		name           string
		action         func(t *testing.T) *meilisearch.TaskInfo
		expectMetadata string
	}{
		{
			name: "AddDocuments with Metadata",
			action: func(t *testing.T) *meilisearch.TaskInfo {
				meta := "meta-add-docs"
				task, err := i.AddDocuments(documents, &meilisearch.DocumentOptions{
					TaskCustomMetadata: meta,
				})
				require.NoError(t, err)
				return task
			},
			expectMetadata: "meta-add-docs",
		},
		{
			name: "UpdateDocuments with Metadata",
			action: func(t *testing.T) *meilisearch.TaskInfo {
				meta := "meta-update-docs"
				updateDocs := []map[string]interface{}{
					{"id": "1", "title": "Updated Document 1"},
				}
				task, err := i.UpdateDocuments(updateDocs, &meilisearch.DocumentOptions{
					TaskCustomMetadata: meta,
				})
				require.NoError(t, err)
				return task
			},
			expectMetadata: "meta-update-docs",
		},
		{
			name: "DeleteDocument (Single) with Metadata",
			action: func(t *testing.T) *meilisearch.TaskInfo {
				meta := "meta-delete-one"
				task, err := i.DeleteDocument("1", &meilisearch.DocumentOptions{
					TaskCustomMetadata: meta,
				})
				require.NoError(t, err)
				return task
			},
			expectMetadata: "meta-delete-one",
		},
		{
			name: "DeleteDocuments (Batch) with Metadata",
			action: func(t *testing.T) *meilisearch.TaskInfo {
				meta := "meta-delete-batch"
				task, err := i.DeleteDocuments([]string{"2"}, &meilisearch.DocumentOptions{
					TaskCustomMetadata: meta,
				})
				require.NoError(t, err)
				return task
			},
			expectMetadata: "meta-delete-batch",
		},
		{
			name: "AddDocumentsNdjson with Metadata",
			action: func(t *testing.T) *meilisearch.TaskInfo {
				meta := "meta-ndjson"
				ndjson := []byte(`{"id": "3", "title": "Ndjson Doc"}`)
				task, err := i.AddDocumentsNdjson(ndjson, &meilisearch.DocumentOptions{
					TaskCustomMetadata: meta,
				})
				require.NoError(t, err)
				return task
			},
			expectMetadata: "meta-ndjson",
		},
		{
			name: "DeleteAllDocuments with Metadata",
			action: func(t *testing.T) *meilisearch.TaskInfo {
				meta := "meta-delete-all"
				task, err := i.DeleteAllDocuments(&meilisearch.DocumentOptions{
					TaskCustomMetadata: meta,
				})
				require.NoError(t, err)
				return task
			},
			expectMetadata: "meta-delete-all",
		},
	}

	// Run Standard Document Option Tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Execute Action
			taskInfo := tt.action(t)

			// 2. Wait for task to be processed (ensures storage)
			testWaitForIndexTask(t, i, taskInfo)

			// 3. Fetch the full task details from the engine
			task, err := sv.GetTask(taskInfo.TaskUID)
			require.NoError(t, err)

			// 4. Verify the metadata matches
			require.Equal(t, tt.expectMetadata, task.CustomMetadata, "CustomMetadata should be persisted and retrievable from the Task object")
		})
	}

	// Special Case: UpdateDocumentsByFunction
	t.Run("UpdateDocumentsByFunction with Metadata", func(t *testing.T) {
		// Ensure we have a doc to update
		setupTask, _ := i.AddDocuments([]map[string]interface{}{{"id": "99", "title": "Function Doc"}}, nil)
		testWaitForIndexTask(t, i, setupTask)

		meta := "meta-function-update"

		// Enable feature
		exp := sv.ExperimentalFeatures()
		exp.SetEditDocumentsByFunction(true)
		_, err := exp.Update()
		require.NoError(t, err)

		// Perform Update
		taskInfo, err := i.UpdateDocumentsByFunction(&meilisearch.UpdateDocumentByFunctionRequest{
			Filter:             "id = 99",
			Function:           "doc.title = \"Updated Function\"",
			TaskCustomMetadata: meta,
		})
		require.NoError(t, err)

		// Wait and Verify
		testWaitForIndexTask(t, i, taskInfo)

		task, err := sv.GetTask(taskInfo.TaskUID)
		require.NoError(t, err)
		require.Equal(t, meta, task.CustomMetadata)
	})

	// Special Case: DeleteDocumentsByFilter
	t.Run("DeleteDocumentsByFilter with Metadata", func(t *testing.T) {
		// Ensure we have a doc to delete
		setupTask, _ := i.AddDocuments([]map[string]interface{}{{"id": "99", "title": "Filter Doc"}}, nil)
		testWaitForIndexTask(t, i, setupTask)

		meta := "meta-delete-filter"
		// Note: "id" was made filterable at the top of the test function
		taskInfo, err := i.DeleteDocumentsByFilter("id = 99", &meilisearch.DocumentOptions{
			TaskCustomMetadata: meta,
		})
		require.NoError(t, err)

		testWaitForIndexTask(t, i, taskInfo)

		task, err := sv.GetTask(taskInfo.TaskUID)
		require.NoError(t, err)
		require.Equal(t, meta, task.CustomMetadata)
	})
}
