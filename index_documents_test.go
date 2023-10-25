package meilisearch

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndex_AddDocuments(t *testing.T) {
	type args struct {
		UID          string
		client       *Client
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
				client: defaultClient,
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
				client: customClient,
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
				client: defaultClient,
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
				client: defaultClient,
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
				client: customClient,
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
				client: defaultClient,
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
		})
	}
}

func TestIndex_AddDocumentsWithPrimaryKey(t *testing.T) {
	type args struct {
		UID          string
		client       *Client
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
				client: defaultClient,
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
				client: customClient,
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
				client: defaultClient,
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
				client: defaultClient,
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
				client: defaultClient,
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

func TestIndex_AddDocumentsInBatches(t *testing.T) {
	type argsNoKey struct {
		UID          string
		client       *Client
		documentsPtr interface{}
		batchSize    int
	}

	type argsWithKey struct {
		UID          string
		client       *Client
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
				client: defaultClient,
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
				client: defaultClient,
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

func testParseCsvDocuments(t *testing.T, documents io.Reader) []map[string]interface{} {
	var (
		docs   []map[string]interface{}
		header []string
	)
	r := csv.NewReader(documents)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		if header == nil {
			header = record
			continue
		}
		doc := make(map[string]interface{})
		for i, key := range header {
			doc[key] = record[i]
		}
		docs = append(docs, doc)
	}
	return docs
}

var testCsvDocuments = []byte(`id,name
1,Alice In Wonderland
2,Pride and Prejudice
3,Le Petit Prince
4,The Great Gatsby
5,Don Quixote
`)

func TestIndex_AddDocumentsCsv(t *testing.T) {
	type args struct {
		UID       string
		client    *Client
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
				client:    defaultClient,
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
	type args struct {
		UID       string
		client    *Client
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
				client:    defaultClient,
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
				client:    defaultClient,
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
				client:    defaultClient,
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

func TestIndex_AddDocumentsCsvInBatches(t *testing.T) {
	type args struct {
		UID       string
		client    *Client
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
				client:    defaultClient,
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
		})
	}

	for _, tt := range tests {
		// Test both the string and io.Reader receiving versions
		testAddDocumentsCsvInBatches(t, tt, false)
		testAddDocumentsCsvInBatches(t, tt, true)
	}
}

func testParseNdjsonDocuments(t *testing.T, documents io.Reader) []map[string]interface{} {
	var docs []map[string]interface{}
	scanner := bufio.NewScanner(documents)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		doc := make(map[string]interface{})
		err := json.Unmarshal([]byte(line), &doc)
		require.NoError(t, err)
		docs = append(docs, doc)
	}
	require.NoError(t, scanner.Err())
	return docs
}

var testNdjsonDocuments = []byte(`{"id": 1, "name": "Alice In Wonderland"}
{"id": 2, "name": "Pride and Prejudice"}
{"id": 3, "name": "Le Petit Prince"}
{"id": 4, "name": "The Great Gatsby"}
{"id": 5, "name": "Don Quixote"}
`)

func TestIndex_AddDocumentsNdjson(t *testing.T) {
	type args struct {
		UID       string
		client    *Client
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
				client:    defaultClient,
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
		})
	}

	for _, tt := range tests {
		// Test both the string and io.Reader receiving versions
		testAddDocumentsNdjson(t, tt, false)
		testAddDocumentsNdjson(t, tt, true)
	}
}

func TestIndex_AddDocumentsNdjsonInBatches(t *testing.T) {
	type args struct {
		UID       string
		client    *Client
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
				client:    defaultClient,
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
		})
	}

	for _, tt := range tests {
		// Test both the string and io.Reader receiving versions
		testAddDocumentsNdjsonInBatches(t, tt, false)
		testAddDocumentsNdjsonInBatches(t, tt, true)
	}
}

func TestIndex_DeleteAllDocuments(t *testing.T) {
	type args struct {
		UID    string
		client *Client
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
				client: defaultClient,
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
				client: customClient,
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

			SetUpBasicIndex(tt.args.UID)
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
	type args struct {
		UID          string
		PrimaryKey   string
		client       *Client
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
				client:     defaultClient,
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
				client:     customClient,
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
				client:     defaultClient,
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
				client:     defaultClient,
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
				client:     customClient,
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
				client:     defaultClient,
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
	type args struct {
		UID          string
		client       *Client
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
				client:     defaultClient,
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
				client:     customClient,
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
				client:     defaultClient,
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
				client:     defaultClient,
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
	type args struct {
		UID            string
		client         *Client
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
				client:         defaultClient,
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
				client:         customClient,
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
				client:         customClient,
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
				client:         customClient,
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

func TestIndex_GetDocument(t *testing.T) {
	type args struct {
		UID         string
		client      *Client
		identifier  string
		request     *DocumentQuery
		documentPtr *docTestBooks
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestIndexBasicGetDocument",
			args: args{
				UID:         "TestIndexBasicGetDocument",
				client:      defaultClient,
				identifier:  "123",
				request:     nil,
				documentPtr: &docTestBooks{},
			},
			wantErr: false,
		},
		{
			name: "TestIndexGetDocumentWithCustomClient",
			args: args{
				UID:         "TestIndexGetDocumentWithCustomClient",
				client:      customClient,
				identifier:  "123",
				request:     nil,
				documentPtr: &docTestBooks{},
			},
			wantErr: false,
		},
		{
			name: "TestIndexGetDocumentWithNoExistingDocument",
			args: args{
				UID:         "TestIndexGetDocumentWithNoExistingDocument",
				client:      defaultClient,
				identifier:  "125",
				request:     nil,
				documentPtr: &docTestBooks{},
			},
			wantErr: true,
		},
		{
			name: "TestIndexGetDocumentWithEmptyParameters",
			args: args{
				UID:         "TestIndexGetDocumentWithEmptyParameters",
				client:      defaultClient,
				identifier:  "125",
				request:     &DocumentQuery{},
				documentPtr: &docTestBooks{},
			},
			wantErr: true,
		},
		{
			name: "TestIndexGetDocumentWithParametersFields",
			args: args{
				UID:        "TestIndexGetDocumentWithParametersFields",
				client:     defaultClient,
				identifier: "125",
				request: &DocumentQuery{
					Fields: []string{"book_id", "title"},
				},
				documentPtr: &docTestBooks{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))
			SetUpBasicIndex(tt.args.UID)

			require.Empty(t, tt.args.documentPtr)
			err := i.GetDocument(tt.args.identifier, tt.args.request, tt.args.documentPtr)
			if tt.wantErr {
				require.Error(t, err)
				require.Empty(t, tt.args.documentPtr)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, tt.args.documentPtr)
				require.Equal(t, strconv.Itoa(tt.args.documentPtr.BookID), tt.args.identifier)
			}
		})
	}
}

func TestIndex_GetDocuments(t *testing.T) {
	type args struct {
		UID     string
		client  *Client
		request *DocumentsQuery
		resp    *DocumentsResult
		filter  []string
	}
	tests := []struct {
		name   string
		args   args
		result int64
	}{
		{
			name: "TestIndexBasicGetDocuments",
			args: args{
				client:  defaultClient,
				request: nil,
				resp:    &DocumentsResult{},
			},
			result: 20,
		},
		{
			name: "TestIndexGetDocumentsWithCustomClient",
			args: args{
				client:  customClient,
				request: nil,
				resp:    &DocumentsResult{},
			},
			result: 20,
		},
		{
			name: "TestIndexGetDocumentsWithEmptyStruct",
			args: args{
				client:  defaultClient,
				request: &DocumentsQuery{},
				resp:    &DocumentsResult{},
			},
			result: 20,
		},
		{
			name: "TestIndexGetDocumentsWithLimit",
			args: args{
				client: defaultClient,
				request: &DocumentsQuery{
					Limit: 3,
				},
				resp: &DocumentsResult{},
			},
			result: 3,
		},
		{
			name: "TestIndexGetDocumentsWithFields",
			args: args{
				client: defaultClient,
				request: &DocumentsQuery{
					Fields: []string{"title"},
				},
				resp: &DocumentsResult{},
			},
			result: 20,
		},
		{
			name: "TestIndexGetDocumentsWithFilterAsString",
			args: args{
				client: defaultClient,
				request: &DocumentsQuery{
					Filter: "book_id = 123",
				},
				resp: &DocumentsResult{},
				filter: []string{
					"book_id",
				},
			},
			result: 1,
		},
		{
			name: "TestIndexGetDocumentsWithFilterAsArray",
			args: args{
				client: defaultClient,
				request: &DocumentsQuery{
					Filter: []string{"tag = Tragedy"},
				},
				resp: &DocumentsResult{},
				filter: []string{
					"tag",
				},
			},
			result: 3,
		},
		{
			name: "TestIndexGetDocumentsWithMultipleFilterWithArrayOfString",
			args: args{
				client: defaultClient,
				request: &DocumentsQuery{
					Filter: []string{"tag = Tragedy", "book_id = 742"},
				},
				resp: &DocumentsResult{},
				filter: []string{
					"tag",
					"book_id",
				},
			},
			result: 1,
		},
		{
			name: "TestIndexGetDocumentsWithMultipleFilterWithInterface",
			args: args{
				client: defaultClient,
				request: &DocumentsQuery{
					Filter: []interface{}{[]string{"tag = Tragedy", "book_id = 123"}},
				},
				resp: &DocumentsResult{},
				filter: []string{
					"tag",
					"book_id",
				},
			},
			result: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index("indexUID")
			t.Cleanup(cleanup(c))
			SetUpIndexForFaceting()

			if tt.args.request != nil && tt.args.request.Filter != nil {
				gotTask, err := i.UpdateFilterableAttributes(&tt.args.filter)
				require.NoError(t, err)
				testWaitForTask(t, i, gotTask)
			}

			err := i.GetDocuments(tt.args.request, tt.args.resp)
			require.NoError(t, err)
			if tt.args.request != nil && tt.args.request.Limit != 0 {
				require.Equal(t, tt.args.request.Limit, int64(len(tt.args.resp.Results)))
			}
			require.Equal(t, tt.result, int64(len(tt.args.resp.Results)))
		})
	}
}

func TestIndex_UpdateDocuments(t *testing.T) {
	type args struct {
		UID          string
		client       *Client
		documentsPtr []docTestBooks
	}
	tests := []struct {
		name string
		args args
		want *Task
	}{
		{
			name: "TestIndexBasicUpdateDocument",
			args: args{
				UID:    "TestIndexBasicUpdateDocument",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
				},
			},
			want: &Task{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    TaskTypeDocumentAdditionOrUpdate,
			},
		},
		{
			name: "TestIndexUpdateDocumentWithCustomClient",
			args: args{
				UID:    "TestIndexUpdateDocumentWithCustomClient",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
				},
			},
			want: &Task{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    TaskTypeDocumentAdditionOrUpdate,
			},
		},
		{
			name: "TestIndexUpdateDocumentOnMultipleDocuments",
			args: args{
				UID:    "TestIndexUpdateDocumentOnMultipleDocuments",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
					{BookID: 1344, Title: "Harry Potter and the Half-Blood Prince"},
					{BookID: 4, Title: "The Hobbit"},
					{BookID: 42, Title: "The Great Gatsby"},
				},
			},
			want: &Task{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    TaskTypeDocumentAdditionOrUpdate,
			},
		},
		{
			name: "TestIndexUpdateDocumentWithNoExistingDocument",
			args: args{
				UID:    "TestIndexUpdateDocumentWithNoExistingDocument",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 237, Title: "One Hundred Years of Solitude"},
				},
			},
			want: &Task{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    TaskTypeDocumentAdditionOrUpdate,
			},
		},
		{
			name: "TestIndexUpdateDocumentWithNoExistingMultipleDocuments",
			args: args{
				UID:    "TestIndexUpdateDocumentWithNoExistingMultipleDocuments",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 246, Title: "One Hundred Years of Solitude"},
					{BookID: 834, Title: "To Kill a Mockingbird"},
					{BookID: 44, Title: "Don Quixote"},
					{BookID: 594, Title: "The Great Gatsby"},
				},
			},
			want: &Task{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    TaskTypeDocumentAdditionOrUpdate,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))
			SetUpBasicIndex(tt.args.UID)

			got, err := i.UpdateDocuments(tt.args.documentsPtr)
			require.NoError(t, err)
			require.GreaterOrEqual(t, got.TaskUID, tt.want.TaskUID)
			require.Equal(t, got.Status, tt.want.Status)
			require.Equal(t, got.Type, tt.want.Type)
			require.NotZero(t, got.EnqueuedAt)

			testWaitForTask(t, i, got)

			var document docTestBooks
			for _, identifier := range tt.args.documentsPtr {
				err = i.GetDocument(strconv.Itoa(identifier.BookID), nil, &document)
				require.NoError(t, err)
				require.Equal(t, identifier.BookID, document.BookID)
				require.Equal(t, identifier.Title, document.Title)
			}
		})
	}
}

func TestIndex_UpdateDocumentsWithPrimaryKey(t *testing.T) {
	type args struct {
		UID          string
		client       *Client
		documentsPtr []docTestBooks
		primaryKey   string
	}
	tests := []struct {
		name string
		args args
		want *Task
	}{
		{
			name: "TestIndexBasicUpdateDocumentsWithPrimaryKey",
			args: args{
				UID:    "TestIndexBasicUpdateDocumentsWithPrimaryKey",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
				},
				primaryKey: "book_id",
			},
			want: &Task{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    TaskTypeDocumentAdditionOrUpdate,
			},
		},
		{
			name: "TestIndexUpdateDocumentsWithPrimaryKeyWithCustomClient",
			args: args{
				UID:    "TestIndexUpdateDocumentsWithPrimaryKeyWithCustomClient",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
				},
				primaryKey: "book_id",
			},
			want: &Task{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    TaskTypeDocumentAdditionOrUpdate,
			},
		},
		{
			name: "TestIndexUpdateDocumentsWithPrimaryKeyOnMultipleDocuments",
			args: args{
				UID:    "TestIndexUpdateDocumentsWithPrimaryKeyOnMultipleDocuments",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
					{BookID: 1344, Title: "Harry Potter and the Half-Blood Prince"},
					{BookID: 4, Title: "The Hobbit"},
					{BookID: 42, Title: "The Great Gatsby"},
				},
				primaryKey: "book_id",
			},
			want: &Task{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    TaskTypeDocumentAdditionOrUpdate,
			},
		},
		{
			name: "TestIndexUpdateDocumentsWithPrimaryKeyWithNoExistingDocument",
			args: args{
				UID:    "TestIndexUpdateDocumentsWithPrimaryKeyWithNoExistingDocument",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 237, Title: "One Hundred Years of Solitude"},
				},
				primaryKey: "book_id",
			},
			want: &Task{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    TaskTypeDocumentAdditionOrUpdate,
			},
		},
		{
			name: "TestIndexUpdateDocumentsWithPrimaryKeyWithNoExistingMultipleDocuments",
			args: args{
				UID:    "TestIndexUpdateDocumentsWithPrimaryKeyWithNoExistingMultipleDocuments",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 246, Title: "One Hundred Years of Solitude"},
					{BookID: 834, Title: "To Kill a Mockingbird"},
					{BookID: 44, Title: "Don Quixote"},
					{BookID: 594, Title: "The Great Gatsby"},
				},
				primaryKey: "book_id",
			},
			want: &Task{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    TaskTypeDocumentAdditionOrUpdate,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))
			SetUpBasicIndex(tt.args.UID)

			got, err := i.UpdateDocuments(tt.args.documentsPtr, tt.args.primaryKey)
			require.NoError(t, err)
			require.GreaterOrEqual(t, got.TaskUID, tt.want.TaskUID)
			require.Equal(t, got.Status, tt.want.Status)
			require.Equal(t, got.Type, tt.want.Type)
			require.NotZero(t, got.EnqueuedAt)

			testWaitForTask(t, i, got)

			var document docTestBooks
			for _, identifier := range tt.args.documentsPtr {
				err = i.GetDocument(strconv.Itoa(identifier.BookID), nil, &document)
				require.NoError(t, err)
				require.Equal(t, identifier.BookID, document.BookID)
				require.Equal(t, identifier.Title, document.Title)
			}
		})
	}
}

func TestIndex_UpdateDocumentsInBatches(t *testing.T) {
	type argsNoKey struct {
		UID          string
		client       *Client
		documentsPtr []docTestBooks
		batchSize    int
	}

	type argsWithKey struct {
		UID          string
		client       *Client
		documentsPtr []docTestBooks
		batchSize    int
		primaryKey   string
	}

	testsNoKey := []struct {
		name string
		args argsNoKey
		want []TaskInfo
	}{
		{
			name: "TestIndexBatchUpdateDocuments",
			args: argsNoKey{
				UID:    "TestIndexBatchUpdateDocuments",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
					{BookID: 124, Title: "One Hundred Years of Solitude 2"},
				},
				batchSize: 1,
			},
			want: []TaskInfo{
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

	testsWithKey := []struct {
		name string
		args argsWithKey
		want []TaskInfo
	}{
		{
			name: "TestIndexBatchUpdateDocuments",
			args: argsWithKey{
				UID:    "TestIndexBatchUpdateDocuments",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
					{BookID: 124, Title: "One Hundred Years of Solitude 2"},
				},
				batchSize:  1,
				primaryKey: "book_id",
			},
			want: []TaskInfo{
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

	for _, tt := range testsNoKey {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))
			SetUpBasicIndex(tt.args.UID)

			got, err := i.UpdateDocumentsInBatches(tt.args.documentsPtr, tt.args.batchSize)
			require.NoError(t, err)
			for i := 0; i < 2; i++ {
				require.GreaterOrEqual(t, got[i].TaskUID, tt.want[i].TaskUID)
				require.Equal(t, got[i].Status, tt.want[i].Status)
				require.Equal(t, got[i].Type, tt.want[i].Type)
				require.NotZero(t, got[i].EnqueuedAt)
			}

			testWaitForBatchTask(t, i, got)

			var document docTestBooks
			for _, identifier := range tt.args.documentsPtr {
				err = i.GetDocument(strconv.Itoa(identifier.BookID), nil, &document)
				require.NoError(t, err)
				require.Equal(t, identifier.BookID, document.BookID)
				require.Equal(t, identifier.Title, document.Title)
			}
		})
	}

	for _, tt := range testsWithKey {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))
			SetUpBasicIndex(tt.args.UID)

			got, err := i.UpdateDocumentsInBatches(tt.args.documentsPtr, tt.args.batchSize, tt.args.primaryKey)
			require.NoError(t, err)
			for i := 0; i < 2; i++ {
				require.GreaterOrEqual(t, got[i].TaskUID, tt.want[i].TaskUID)
				require.Equal(t, got[i].Status, tt.want[i].Status)
				require.Equal(t, got[i].Type, tt.want[i].Type)
				require.NotZero(t, got[i].EnqueuedAt)
			}

			testWaitForBatchTask(t, i, got)

			var document docTestBooks
			for _, identifier := range tt.args.documentsPtr {
				err = i.GetDocument(strconv.Itoa(identifier.BookID), nil, &document)
				require.NoError(t, err)
				require.Equal(t, identifier.BookID, document.BookID)
				require.Equal(t, identifier.Title, document.Title)
			}
		})
	}
}

func TestIndex_UpdateDocumentsCsv(t *testing.T) {
	type args struct {
		UID       string
		client    *Client
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
				client:    defaultClient,
				documents: testCsvDocuments,
			},
			wantResp: &TaskInfo{
				TaskUID: 0,
				Status:  "enqueued",
				Type:    TaskTypeDocumentAdditionOrUpdate,
			},
		},
	}

	testUpdateDocumentsCsv := func(t *testing.T, tt testData, testReader bool) {
		name := tt.name + "UpdateDocumentsCsv"
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
				gotResp, err = i.UpdateDocumentsCsvFromReader(bytes.NewReader(tt.args.documents), nil)
			} else {
				gotResp, err = i.UpdateDocumentsCsv(tt.args.documents, nil)
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
		testUpdateDocumentsCsv(t, tt, false)
		testUpdateDocumentsCsv(t, tt, true)
	}
}

func TestIndex_UpdateDocumentsCsvWithOptions(t *testing.T) {
	type args struct {
		UID       string
		client    *Client
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
			name: "TestIndexBasicUpdateDocumentsCsvWithOptions",
			args: args{
				UID:       "csv",
				client:    defaultClient,
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
			name: "TestIndexBasicUpdateDocumentsCsvWithPrimaryKey",
			args: args{
				UID:       "csv",
				client:    defaultClient,
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
			name: "TestIndexBasicUpdateDocumentsCsvWithCsvDelimiter",
			args: args{
				UID:       "csv",
				client:    defaultClient,
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

	testUpdateDocumentsCsv := func(t *testing.T, tt testData, testReader bool) {
		name := tt.name + "UpdateDocumentsCsv"
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
				gotResp, err = i.UpdateDocumentsCsvFromReader(bytes.NewReader(tt.args.documents), tt.args.options)
			} else {
				gotResp, err = i.UpdateDocumentsCsv(tt.args.documents, tt.args.options)
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
		testUpdateDocumentsCsv(t, tt, false)
		testUpdateDocumentsCsv(t, tt, true)
	}
}

func TestIndex_UpdateDocumentsCsvInBatches(t *testing.T) {
	type args struct {
		UID       string
		client    *Client
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
				client:    defaultClient,
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

	testUpdateDocumentsCsvInBatches := func(t *testing.T, tt testData, testReader bool) {
		name := tt.name + "UpdateDocumentsCsv"
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
				gotResp, err = i.UpdateDocumentsCsvFromReaderInBatches(bytes.NewReader(tt.args.documents), tt.args.batchSize, nil)
			} else {
				gotResp, err = i.UpdateDocumentsCsvInBatches(tt.args.documents, tt.args.batchSize, nil)
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
		})
	}

	for _, tt := range tests {
		// Test both the string and io.Reader receiving versions
		testUpdateDocumentsCsvInBatches(t, tt, false)
		testUpdateDocumentsCsvInBatches(t, tt, true)
	}
}

func TestIndex_UpdateDocumentsNdjson(t *testing.T) {
	type args struct {
		UID       string
		client    *Client
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
				client:    defaultClient,
				documents: testNdjsonDocuments,
			},
			wantResp: &TaskInfo{
				TaskUID: 0,
				Status:  "enqueued",
				Type:    TaskTypeDocumentAdditionOrUpdate,
			},
		},
	}

	testUpdateDocumentsNdjson := func(t *testing.T, tt testData, testReader bool) {
		name := tt.name + "UpdateDocumentsNdjson"
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
				gotResp, err = i.UpdateDocumentsNdjsonFromReader(bytes.NewReader(tt.args.documents))
			} else {
				gotResp, err = i.UpdateDocumentsNdjson(tt.args.documents)
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
		testUpdateDocumentsNdjson(t, tt, false)
		testUpdateDocumentsNdjson(t, tt, true)
	}
}

func TestIndex_UpdateDocumentsNdjsonInBatches(t *testing.T) {
	type args struct {
		UID       string
		client    *Client
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
				client:    defaultClient,
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

	testUpdateDocumentsNdjsonInBatches := func(t *testing.T, tt testData, testReader bool) {
		name := tt.name + "UpdateDocumentsNdjson"
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
				gotResp, err = i.updateDocumentsNdjsonFromReaderInBatches(bytes.NewReader(tt.args.documents), tt.args.batchSize)
			} else {
				gotResp, err = i.UpdateDocumentsNdjsonInBatches(tt.args.documents, tt.args.batchSize)
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
		})
	}

	for _, tt := range tests {
		// Test both the string and io.Reader receiving versions
		testUpdateDocumentsNdjsonInBatches(t, tt, false)
		testUpdateDocumentsNdjsonInBatches(t, tt, true)
	}
}

func Test_transformStringVariadicToMap(t *testing.T) {
	type args struct {
		primaryKey []string
	}
	tests := []struct {
		name        string
		args        args
		wantOptions map[string]string
	}{
		{
			name: "TestCreateOptionsInterface",
			args: args{
				[]string{
					"id",
				},
			},
			wantOptions: map[string]string{
				"primaryKey": "id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOptions := transformStringVariadicToMap(tt.args.primaryKey...)
			require.Equal(t, tt.wantOptions, gotOptions)
		})
	}
}

func Test_generateQueryForOptions(t *testing.T) {
	type args struct {
		options map[string]string
	}
	tests := []struct {
		name         string
		args         args
		wantUrlQuery string
	}{
		{
			name: "TestGenerateQueryForOptions",
			args: args{
				options: map[string]string{
					"primaryKey":   "id",
					"csvDelimiter": ",",
				},
			},
			wantUrlQuery: "csvDelimiter=%2C&primaryKey=id",
		},
		{
			name: "TestGenerateQueryForPrimaryKey",
			args: args{
				options: map[string]string{
					"primaryKey": "id",
				},
			},
			wantUrlQuery: "primaryKey=id",
		},
		{
			name: "TestGenerateQueryForCsvDelimiter",
			args: args{
				options: map[string]string{
					"csvDelimiter": ",",
				},
			},
			wantUrlQuery: "csvDelimiter=%2C",
		},
		{
			name: "TestGenerateQueryWithNull",
			args: args{
				options: nil,
			},
			wantUrlQuery: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUrlQuery := generateQueryForOptions(tt.args.options)
			require.Equal(t, tt.wantUrlQuery, gotUrlQuery)
		})
	}
}

func Test_transformCsvDocumentsQueryToMap(t *testing.T) {
	type args struct {
		options *CsvDocumentsQuery
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "TestTransformCsvDocumentsQueryToMap",
			args: args{
				options: &CsvDocumentsQuery{
					PrimaryKey:   "id",
					CsvDelimiter: ",",
				},
			},
			want: map[string]string{
				"primaryKey":   "id",
				"csvDelimiter": ",",
			},
		},
		{
			name: "TestTransformCsvDocumentsQueryToMapWithPrimaryKey",
			args: args{
				options: &CsvDocumentsQuery{
					PrimaryKey: "id",
				},
			},
			want: map[string]string{
				"primaryKey": "id",
			},
		},
		{
			name: "TestTransformCsvDocumentsQueryToMapEmpty",
			args: args{
				options: &CsvDocumentsQuery{},
			},
			want: map[string]string{},
		},
		{
			name: "TestTransformCsvDocumentsQueryToMapNull",
			args: args{
				options: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transformCsvDocumentsQueryToMap(tt.args.options)
			require.Equal(t, tt.want, got)
		})
	}
}
