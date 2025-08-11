package meilisearch

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"net/http"
	"reflect"
	"strings"
)

func (i *index) AddDocuments(documentsPtr interface{}, primaryKey *string) (*TaskInfo, error) {
	return i.AddDocumentsWithContext(context.Background(), documentsPtr, primaryKey)
}

func (i *index) AddDocumentsWithContext(ctx context.Context, documentsPtr interface{}, primaryKey *string) (*TaskInfo, error) {
	return i.addDocuments(ctx, documentsPtr, contentTypeJSON, transformStringToMap(primaryKey))
}

func (i *index) AddDocumentsInBatches(documentsPtr interface{}, batchSize int, primaryKey *string) ([]TaskInfo, error) {
	return i.AddDocumentsInBatchesWithContext(context.Background(), documentsPtr, batchSize, primaryKey)
}

func (i *index) AddDocumentsInBatchesWithContext(ctx context.Context, documentsPtr interface{}, batchSize int, primaryKey *string) ([]TaskInfo, error) {
	return i.saveDocumentsInBatches(ctx, documentsPtr, batchSize, i.AddDocumentsWithContext, primaryKey)
}

func (i *index) AddDocumentsCsv(documents []byte, options *CsvDocumentsQuery) (*TaskInfo, error) {
	return i.AddDocumentsCsvWithContext(context.Background(), documents, options)
}

func (i *index) AddDocumentsCsvWithContext(ctx context.Context, documents []byte, options *CsvDocumentsQuery) (*TaskInfo, error) {
	// []byte avoids JSON conversion in Client.sendRequest()
	return i.addDocuments(ctx, documents, contentTypeCSV, transformCsvDocumentsQueryToMap(options))
}

func (i *index) AddDocumentsCsvInBatches(documents []byte, batchSize int, options *CsvDocumentsQuery) ([]TaskInfo, error) {
	return i.AddDocumentsCsvInBatchesWithContext(context.Background(), documents, batchSize, options)
}

func (i *index) AddDocumentsCsvInBatchesWithContext(ctx context.Context, documents []byte, batchSize int, options *CsvDocumentsQuery) ([]TaskInfo, error) {
	// Reuse io.Reader implementation
	return i.AddDocumentsCsvFromReaderInBatchesWithContext(ctx, bytes.NewReader(documents), batchSize, options)
}

func (i *index) AddDocumentsCsvFromReaderInBatches(documents io.Reader, batchSize int, options *CsvDocumentsQuery) (resp []TaskInfo, err error) {
	return i.AddDocumentsCsvFromReaderInBatchesWithContext(context.Background(), documents, batchSize, options)
}

func (i *index) AddDocumentsCsvFromReaderInBatchesWithContext(ctx context.Context, documents io.Reader, batchSize int, options *CsvDocumentsQuery) (resp []TaskInfo, err error) {
	return i.saveDocumentsFromReaderInBatches(ctx, documents, batchSize, i.AddDocumentsCsvWithContext, options)
}

func (i *index) AddDocumentsCsvFromReader(documents io.Reader, options *CsvDocumentsQuery) (resp *TaskInfo, err error) {
	return i.AddDocumentsCsvFromReaderWithContext(context.Background(), documents, options)
}

func (i *index) AddDocumentsCsvFromReaderWithContext(ctx context.Context, documents io.Reader, options *CsvDocumentsQuery) (resp *TaskInfo, err error) {
	return i.addDocumentsFromReader(ctx, documents, contentTypeCSV, transformCsvDocumentsQueryToMap(options))
}

func (i *index) AddDocumentsNdjson(documents []byte, primaryKey *string) (*TaskInfo, error) {
	return i.AddDocumentsNdjsonWithContext(context.Background(), documents, primaryKey)
}

func (i *index) AddDocumentsNdjsonWithContext(ctx context.Context, documents []byte, primaryKey *string) (*TaskInfo, error) {
	// []byte avoids JSON conversion in Client.sendRequest()
	return i.addDocumentsFromReader(ctx, bytes.NewReader(documents), contentTypeNDJSON, transformStringToMap(primaryKey))
}

func (i *index) AddDocumentsNdjsonInBatches(documents []byte, batchSize int, primaryKey *string) ([]TaskInfo, error) {
	return i.AddDocumentsNdjsonInBatchesWithContext(context.Background(), documents, batchSize, primaryKey)
}

func (i *index) AddDocumentsNdjsonInBatchesWithContext(ctx context.Context, documents []byte, batchSize int, primaryKey *string) ([]TaskInfo, error) {
	// Reuse io.Reader implementation
	return i.AddDocumentsNdjsonFromReaderInBatchesWithContext(ctx, bytes.NewReader(documents), batchSize, primaryKey)
}

func (i *index) AddDocumentsNdjsonFromReaderInBatches(documents io.Reader, batchSize int, primaryKey *string) (resp []TaskInfo, err error) {
	return i.AddDocumentsNdjsonFromReaderInBatchesWithContext(context.Background(), documents, batchSize, primaryKey)
}

func (i *index) AddDocumentsNdjsonFromReaderInBatchesWithContext(ctx context.Context, documents io.Reader, batchSize int, primaryKey *string) (resp []TaskInfo, err error) {
	// NDJSON files supposed to contain a valid JSON document in each line, so
	// it's safe to split by lines.
	// Lines are read and sent continuously to avoid reading all content into
	// memory. However, this means that only part of the documents might be
	// added successfully.

	sendNdjsonLines := func(lines []string) (*TaskInfo, error) {
		b := new(bytes.Buffer)
		for _, line := range lines {
			_, err := b.WriteString(line)
			if err != nil {
				return nil, fmt.Errorf("could not write NDJSON line: %w", err)
			}
			err = b.WriteByte('\n')
			if err != nil {
				return nil, fmt.Errorf("could not write NDJSON line: %w", err)
			}
		}

		resp, err := i.AddDocumentsNdjsonWithContext(ctx, b.Bytes(), primaryKey)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}

	var (
		responses []TaskInfo
		lines     []string
	)

	scanner := bufio.NewScanner(documents)
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines (NDJSON might not allow this, but just to be sure)
		if line == "" {
			continue
		}

		lines = append(lines, line)
		// After reaching batchSize send NDJSON lines
		if len(lines) == batchSize {
			resp, err := sendNdjsonLines(lines)
			if err != nil {
				return nil, err
			}
			responses = append(responses, *resp)
			lines = nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("could not read NDJSON: %w", err)
	}

	// Send remaining records as the last batch if there is any
	if len(lines) > 0 {
		resp, err := sendNdjsonLines(lines)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *resp)
	}

	return responses, nil
}

func (i *index) AddDocumentsNdjsonFromReader(documents io.Reader, primaryKey *string) (resp *TaskInfo, err error) {
	return i.AddDocumentsNdjsonFromReaderWithContext(context.Background(), documents, primaryKey)
}

func (i *index) AddDocumentsNdjsonFromReaderWithContext(ctx context.Context, documents io.Reader, primaryKey *string) (resp *TaskInfo, err error) {
	// Using io.Reader would avoid JSON conversion in Client.sendRequest(), but
	// read content to memory anyway because of problems with streamed bodies
	data, err := io.ReadAll(documents)
	if err != nil {
		return nil, fmt.Errorf("could not read documents: %w", err)
	}
	return i.addDocuments(ctx, data, contentTypeNDJSON, transformStringToMap(primaryKey))
}

func (i *index) UpdateDocuments(documentsPtr interface{}, primaryKey *string) (*TaskInfo, error) {
	return i.UpdateDocumentsWithContext(context.Background(), documentsPtr, primaryKey)
}

func (i *index) UpdateDocumentsWithContext(ctx context.Context, documentsPtr interface{}, primaryKey *string) (*TaskInfo, error) {
	return i.updateDocuments(ctx, documentsPtr, contentTypeJSON, transformStringToMap(primaryKey))
}

func (i *index) UpdateDocumentsInBatches(documentsPtr interface{}, batchSize int, primaryKey *string) ([]TaskInfo, error) {
	return i.UpdateDocumentsInBatchesWithContext(context.Background(), documentsPtr, batchSize, primaryKey)
}

func (i *index) UpdateDocumentsInBatchesWithContext(ctx context.Context, documentsPtr interface{}, batchSize int, primaryKey *string) ([]TaskInfo, error) {
	return i.saveDocumentsInBatches(ctx, documentsPtr, batchSize, i.UpdateDocumentsWithContext, primaryKey)
}

func (i *index) UpdateDocumentsCsv(documents []byte, options *CsvDocumentsQuery) (*TaskInfo, error) {
	return i.UpdateDocumentsCsvWithContext(context.Background(), documents, options)
}

func (i *index) UpdateDocumentsCsvWithContext(ctx context.Context, documents []byte, options *CsvDocumentsQuery) (*TaskInfo, error) {
	return i.updateDocuments(ctx, documents, contentTypeCSV, transformCsvDocumentsQueryToMap(options))
}

func (i *index) UpdateDocumentsCsvInBatches(documents []byte, batchSize int, options *CsvDocumentsQuery) ([]TaskInfo, error) {
	return i.UpdateDocumentsCsvInBatchesWithContext(context.Background(), documents, batchSize, options)
}

func (i *index) UpdateDocumentsCsvInBatchesWithContext(ctx context.Context, documents []byte, batchSize int, options *CsvDocumentsQuery) ([]TaskInfo, error) {
	// Reuse io.Reader implementation
	return i.updateDocumentsCsvFromReaderInBatches(ctx, bytes.NewReader(documents), batchSize, options)
}

func (i *index) UpdateDocumentsNdjson(documents []byte, primaryKey *string) (*TaskInfo, error) {
	return i.UpdateDocumentsNdjsonWithContext(context.Background(), documents, primaryKey)
}

func (i *index) UpdateDocumentsNdjsonWithContext(ctx context.Context, documents []byte, primaryKey *string) (*TaskInfo, error) {
	return i.updateDocuments(ctx, documents, contentTypeNDJSON, transformStringToMap(primaryKey))
}

func (i *index) UpdateDocumentsNdjsonInBatches(documents []byte, batchSize int, primaryKey *string) ([]TaskInfo, error) {
	return i.UpdateDocumentsNdjsonInBatchesWithContext(context.Background(), documents, batchSize, primaryKey)
}

func (i *index) UpdateDocumentsNdjsonInBatchesWithContext(ctx context.Context, documents []byte, batchSize int, primaryKey *string) ([]TaskInfo, error) {
	return i.updateDocumentsNdjsonFromReaderInBatches(ctx, bytes.NewReader(documents), batchSize, primaryKey)
}

func (i *index) UpdateDocumentsByFunction(req *UpdateDocumentByFunctionRequest) (*TaskInfo, error) {
	return i.UpdateDocumentsByFunctionWithContext(context.Background(), req)
}

func (i *index) UpdateDocumentsByFunctionWithContext(ctx context.Context, req *UpdateDocumentByFunctionRequest) (*TaskInfo, error) {
	resp := new(TaskInfo)
	r := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/documents/edit",
		method:              http.MethodPost,
		withRequest:         req,
		withResponse:        resp,
		contentType:         contentTypeJSON,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateDocumentsByFunction",
	}
	if err := i.client.executeRequest(ctx, r); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetDocument(identifier string, request *DocumentQuery, documentPtr interface{}) error {
	return i.GetDocumentWithContext(context.Background(), identifier, request, documentPtr)
}

func (i *index) GetDocumentWithContext(ctx context.Context, identifier string, request *DocumentQuery, documentPtr interface{}) error {
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/documents/" + identifier,
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        documentPtr,
		withQueryParams:     map[string]string{},
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetDocument",
	}
	if request != nil {
		if len(request.Fields) != 0 {
			req.withQueryParams["fields"] = strings.Join(request.Fields, ",")
		}
		if request.RetrieveVectors {
			req.withQueryParams["retrieveVectors"] = "true"
		}
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return err
	}
	return nil
}

func (i *index) GetDocuments(param *DocumentsQuery, resp *DocumentsResult) error {
	return i.GetDocumentsWithContext(context.Background(), param, resp)
}

func (i *index) GetDocumentsWithContext(ctx context.Context, param *DocumentsQuery, resp *DocumentsResult) error {
	if param == nil {
		param = &DocumentsQuery{}
	}
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/documents/fetch",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         param,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetDocuments",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return VersionErrorHintMessage(err, req)
	}
	return nil
}

func (i *index) DeleteDocument(identifier string) (*TaskInfo, error) {
	return i.DeleteDocumentWithContext(context.Background(), identifier)
}

func (i *index) DeleteDocumentWithContext(ctx context.Context, identifier string) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/documents/" + identifier,
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "DeleteDocument",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) DeleteDocuments(identifiers []string) (*TaskInfo, error) {
	return i.DeleteDocumentsWithContext(context.Background(), identifiers)
}

func (i *index) DeleteDocumentsWithContext(ctx context.Context, identifiers []string) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/documents/delete-batch",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         identifiers,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "DeleteDocuments",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) DeleteDocumentsByFilter(filter interface{}) (*TaskInfo, error) {
	return i.DeleteDocumentsByFilterWithContext(context.Background(), filter)
}

func (i *index) DeleteDocumentsByFilterWithContext(ctx context.Context, filter interface{}) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:    "/indexes/" + i.uid + "/documents/delete",
		method:      http.MethodPost,
		contentType: contentTypeJSON,
		withRequest: map[string]interface{}{
			"filter": filter,
		},
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "DeleteDocumentsByFilter",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, VersionErrorHintMessage(err, req)
	}
	return resp, nil
}

func (i *index) DeleteAllDocuments() (*TaskInfo, error) {
	return i.DeleteAllDocumentsWithContext(context.Background())
}

func (i *index) DeleteAllDocumentsWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/documents",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "DeleteAllDocuments",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) addDocuments(ctx context.Context, documents interface{}, contentType string, options map[string]string) (*TaskInfo, error) {
	resp := new(TaskInfo)
	endpoint := "/indexes/" + i.uid + "/documents"
	if len(options) > 0 {
		for key, val := range options {
			if key == "primaryKey" {
				i.primaryKey = val
			}
		}
		endpoint += "?" + generateQueryForOptions(options)
	}
	req := &internalRequest{
		endpoint:            endpoint,
		method:              http.MethodPost,
		contentType:         contentType,
		withRequest:         documents,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "AddDocuments",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) addDocumentsFromReader(ctx context.Context, r io.Reader, contentType string, options map[string]string) (*TaskInfo, error) {
	resp := new(TaskInfo)
	endpoint := "/indexes/" + i.uid + "/documents"
	if len(options) > 0 {
		for key, val := range options {
			if key == "primaryKey" {
				i.primaryKey = val
			}
		}
		endpoint += "?" + generateQueryForOptions(options)
	}
	req := &internalRequest{
		endpoint:            endpoint,
		method:              http.MethodPost,
		contentType:         contentType,
		withRequest:         r,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "AddDocuments",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) saveDocumentsFromReaderInBatches(ctx context.Context, documents io.Reader, batchSize int, documentsCsvFunc func(ctx context.Context, recs []byte, op *CsvDocumentsQuery) (resp *TaskInfo, err error), options *CsvDocumentsQuery) (resp []TaskInfo, err error) {
	// Because of the possibility of multiline fields it's not safe to split
	// into batches by lines, we'll have to parse the file and reassemble it
	// into smaller parts. RFC 4180 compliant input with a header row is
	// expected.
	// Records are read and sent continuously to avoid reading all content
	// into memory. However, this means that only part of the documents might
	// be added successfully.

	var (
		responses []TaskInfo
		header    []string
		records   [][]string
	)

	r := csv.NewReader(documents)
	for {
		// Read CSV record (empty lines and comments are already skipped by csv.Reader)
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("could not read CSV record: %w", err)
		}

		// Store first record as header
		if header == nil {
			header = record
			continue
		}

		// Add header record to every batch
		if len(records) == 0 {
			records = append(records, header)
		}

		records = append(records, record)

		// After reaching batchSize (not counting the header record) assemble a CSV file and send records
		if len(records) == batchSize+1 {
			resp, err := sendCsvRecords(ctx, documentsCsvFunc, records, options)
			if err != nil {
				return nil, err
			}
			responses = append(responses, *resp)
			records = nil
		}
	}

	// Send remaining records as the last batch if there is any
	if len(records) > 0 {
		resp, err := sendCsvRecords(ctx, documentsCsvFunc, records, options)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *resp)
	}

	return responses, nil
}

func (i *index) saveDocumentsInBatches(ctx context.Context, documentsPtr interface{}, batchSize int, documentFunc func(ctx context.Context, documentsPtr interface{}, primaryKey *string) (resp *TaskInfo, err error), primaryKey *string) (resp []TaskInfo, err error) {
	arr := reflect.ValueOf(documentsPtr)
	lenDocs := arr.Len()
	numBatches := int(math.Ceil(float64(lenDocs) / float64(batchSize)))
	resp = make([]TaskInfo, numBatches)

	for j := 0; j < numBatches; j++ {
		end := (j + 1) * batchSize
		if end > lenDocs {
			end = lenDocs
		}

		batch := arr.Slice(j*batchSize, end).Interface()

		respID, err := documentFunc(ctx, batch, primaryKey)
		if err != nil {
			return nil, err
		}

		resp[j] = *respID

	}

	return resp, nil
}

func (i *index) updateDocuments(ctx context.Context, documentsPtr interface{}, contentType string, options map[string]string) (resp *TaskInfo, err error) {
	resp = &TaskInfo{}
	endpoint := ""
	if options == nil {
		endpoint = "/indexes/" + i.uid + "/documents"
	} else {
		for key, val := range options {
			if key == "primaryKey" {
				i.primaryKey = val
			}
		}
		endpoint = "/indexes/" + i.uid + "/documents?" + generateQueryForOptions(options)
	}
	req := &internalRequest{
		endpoint:            endpoint,
		method:              http.MethodPut,
		contentType:         contentType,
		withRequest:         documentsPtr,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateDocuments",
	}
	if err = i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) updateDocumentsCsvFromReaderInBatches(ctx context.Context, documents io.Reader, batchSize int, options *CsvDocumentsQuery) (resp []TaskInfo, err error) {
	return i.saveDocumentsFromReaderInBatches(ctx, documents, batchSize, i.UpdateDocumentsCsvWithContext, options)
}

func (i *index) updateDocumentsNdjsonFromReaderInBatches(ctx context.Context, documents io.Reader, batchSize int, primaryKey *string) (resp []TaskInfo, err error) {
	// NDJSON files supposed to contain a valid JSON document in each line, so
	// it's safe to split by lines.
	// Lines are read and sent continuously to avoid reading all content into
	// memory. However, this means that only part of the documents might be
	// added successfully.

	sendNdjsonLines := func(lines []string) (*TaskInfo, error) {
		b := new(bytes.Buffer)
		for _, line := range lines {
			_, err := b.WriteString(line)
			if err != nil {
				return nil, fmt.Errorf("could not write NDJSON line: %w", err)
			}
			err = b.WriteByte('\n')
			if err != nil {
				return nil, fmt.Errorf("could not write NDJSON line: %w", err)
			}
		}

		resp, err := i.UpdateDocumentsNdjsonWithContext(ctx, b.Bytes(), primaryKey)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}

	var (
		responses []TaskInfo
		lines     []string
	)

	scanner := bufio.NewScanner(documents)
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines (NDJSON might not allow this, but just to be sure)
		if line == "" {
			continue
		}

		lines = append(lines, line)
		// After reaching batchSize send NDJSON lines
		if len(lines) == batchSize {
			resp, err := sendNdjsonLines(lines)
			if err != nil {
				return nil, err
			}
			responses = append(responses, *resp)
			lines = nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("could not read NDJSON: %w", err)
	}

	// Send remaining records as the last batch if there is any
	if len(lines) > 0 {
		resp, err := sendNdjsonLines(lines)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *resp)
	}

	return responses, nil
}
