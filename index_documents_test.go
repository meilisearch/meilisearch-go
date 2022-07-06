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
		wantResp     *Task
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
				wantResp: &Task{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
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
				wantResp: &Task{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
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
				wantResp: &Task{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
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
				wantResp: &Task{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
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
				wantResp: &Task{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
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
				wantResp: &Task{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
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
		wantResp     *Task
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
				wantResp: &Task{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
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
				wantResp: &Task{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
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
				wantResp: &Task{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
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
				wantResp: &Task{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
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
				wantResp: &Task{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
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
		wantResp      []Task
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
			wantResp: []Task{
				{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
				},
				{
					TaskUID: 1,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
				},
			},
		},
	}

	testsWithKey := []struct {
		name          string
		args          argsWithKey
		wantResp      []Task
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
			wantResp: []Task{
				{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
				},
				{
					TaskUID: 1,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
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
		wantResp *Task
	}

	tests := []testData{
		{
			name: "TestIndexBasic",
			args: args{
				UID:       "csv",
				client:    defaultClient,
				documents: testCsvDocuments,
			},
			wantResp: &Task{
				TaskUID: 0,
				Status:  "enqueued",
				Type:    "documentAdditionOrUpdate",
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
				gotResp *Task
				err     error
			)

			if testReader {
				gotResp, err = i.AddDocumentsCsvFromReader(bytes.NewReader(tt.args.documents))
			} else {
				gotResp, err = i.AddDocumentsCsv(tt.args.documents)
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
		wantResp []Task
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
			wantResp: []Task{
				{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
				},
				{
					TaskUID: 1,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
				},
				{
					TaskUID: 2,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
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
				gotResp []Task
				err     error
			)

			if testReader {
				gotResp, err = i.AddDocumentsCsvFromReaderInBatches(bytes.NewReader(tt.args.documents), tt.args.batchSize)
			} else {
				gotResp, err = i.AddDocumentsCsvInBatches(tt.args.documents, tt.args.batchSize)
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
		wantResp *Task
	}

	tests := []testData{
		{
			name: "TestIndexBasic",
			args: args{
				UID:       "ndjson",
				client:    defaultClient,
				documents: testNdjsonDocuments,
			},
			wantResp: &Task{
				TaskUID: 0,
				Status:  "enqueued",
				Type:    "documentAdditionOrUpdate",
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
				gotResp *Task
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
		wantResp []Task
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
			wantResp: []Task{
				{
					TaskUID: 0,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
				},
				{
					TaskUID: 1,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
				},
				{
					TaskUID: 2,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
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
				gotResp []Task
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
		wantResp *Task
	}{
		{
			name: "TestIndexBasicDeleteAllDocuments",
			args: args{
				UID:    "TestIndexBasicDeleteAllDocuments",
				client: defaultClient,
			},
			wantResp: &Task{
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
			wantResp: &Task{
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
		wantResp *Task
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
			wantResp: &Task{
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
			wantResp: &Task{
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
			wantResp: &Task{
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
			wantResp: &Task{
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
			wantResp: &Task{
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
			wantResp: &Task{
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
		wantResp *Task
	}{
		{
			name: "TestIndexBasicDeleteDocument",
			args: args{
				UID:        "1",
				client:     defaultClient,
				identifier: []string{"123"},
				documentsPtr: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
			wantResp: &Task{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    "documentDeletion",
			},
		},
		{
			name: "TestIndexDeleteDocumentWithCustomClient",
			args: args{
				UID:        "2",
				client:     customClient,
				identifier: []string{"123"},
				documentsPtr: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
			wantResp: &Task{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    "documentDeletion",
			},
		},
		{
			name: "TestIndexBasicDeleteDocument",
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
			wantResp: &Task{
				TaskUID: 1,
				Status:  "enqueued",
				Type:    "documentDeletion",
			},
		},
		{
			name: "TestIndexBasicDeleteDocument",
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
			wantResp: &Task{
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
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestIndexBasicGetDocuments",
			args: args{
				UID:     "TestIndexBasicGetDocuments",
				client:  defaultClient,
				request: nil,
				resp:    &DocumentsResult{},
			},
		},
		{
			name: "TestIndexGetDocumentsWithCustomClient",
			args: args{
				UID:     "TestIndexGetDocumentsWithCustomClient",
				client:  customClient,
				request: nil,
				resp:    &DocumentsResult{},
			},
		},
		{
			name: "TestIndexGetDocumentsWithEmptyStruct",
			args: args{
				UID:     "TestIndexGetDocumentsWithNoExistingDocument",
				client:  defaultClient,
				request: &DocumentsQuery{},
				resp:    &DocumentsResult{},
			},
		},
		{
			name: "TestIndexGetDocumentsWithLimit",
			args: args{
				UID:    "TestIndexGetDocumentsWithNoExistingDocument",
				client: defaultClient,
				request: &DocumentsQuery{
					Limit: 3,
				},
				resp: &DocumentsResult{},
			},
		},
		{
			name: "TestIndexGetDocumentsWithLimit",
			args: args{
				UID:    "TestIndexGetDocumentsWithNoExistingDocument",
				client: defaultClient,
				request: &DocumentsQuery{
					Fields: []string{"title"},
				},
				resp: &DocumentsResult{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))
			SetUpBasicIndex(tt.args.UID)

			err := i.GetDocuments(tt.args.request, tt.args.resp)
			require.NoError(t, err)
			if tt.args.request != nil && tt.args.request.Limit != 0 {
				require.Equal(t, tt.args.request.Limit, int64(len(tt.args.resp.Results)))
			} else {
				require.Equal(t, 6, len(tt.args.resp.Results))
			}
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
				Type:    "documentAdditionOrUpdate",
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
				Type:    "documentAdditionOrUpdate",
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
				Type:    "documentAdditionOrUpdate",
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
				Type:    "documentAdditionOrUpdate",
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
				Type:    "documentAdditionOrUpdate",
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
				Type:    "documentAdditionOrUpdate",
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
				Type:    "documentAdditionOrUpdate",
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
				Type:    "documentAdditionOrUpdate",
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
				Type:    "documentAdditionOrUpdate",
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
				Type:    "documentAdditionOrUpdate",
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
		want []Task
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
			want: []Task{
				{
					TaskUID: 1,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
				},
				{
					TaskUID: 2,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
				},
			},
		},
	}

	testsWithKey := []struct {
		name string
		args argsWithKey
		want []Task
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
			want: []Task{
				{
					TaskUID: 1,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
				},
				{
					TaskUID: 2,
					Status:  "enqueued",
					Type:    "documentAdditionOrUpdate",
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
