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

func (i *index) AddDocuments(ctx context.Context, documentsPtr interface{}, opts *DocumentOptions) (*TaskInfo, error) {
	return i.addDocuments(ctx, documentsPtr, contentTypeJSON, transformDocumentOptionsToMap(opts))
}

func (i *index) AddDocumentsInBatches(ctx context.Context, documentsPtr interface{}, batchSize int, opts *DocumentOptions) ([]TaskInfo, error) {
	return i.saveDocumentsInBatches(ctx, documentsPtr, batchSize, i.AddDocuments, opts)
}

func (i *index) AddDocumentsCsv(ctx context.Context, documents []byte, options *CsvDocumentsQuery) (*TaskInfo, error) {
	// []byte avoids JSON conversion in Client.sendRequest()
	return i.addDocuments(ctx, documents, contentTypeCSV, transformCsvDocumentsQueryToMap(options))
}

func (i *index) AddDocumentsCsvInBatches(ctx context.Context, documents []byte, batchSize int, options *CsvDocumentsQuery) ([]TaskInfo, error) {
	// Reuse io.Reader implementation
	return i.AddDocumentsCsvFromReaderInBatches(ctx, bytes.NewReader(documents), batchSize, options)
}

func (i *index) AddDocumentsCsvFromReaderInBatches(ctx context.Context, documents io.Reader, batchSize int, options *CsvDocumentsQuery) (resp []TaskInfo, err error) {
	return i.saveDocumentsFromReaderInBatches(ctx, documents, batchSize, i.AddDocumentsCsv, options)
}

func (i *index) AddDocumentsCsvFromReader(ctx context.Context, documents io.Reader, options *CsvDocumentsQuery) (resp *TaskInfo, err error) {
	return i.addDocumentsFromReader(ctx, documents, contentTypeCSV, transformCsvDocumentsQueryToMap(options))
}

func (i *index) AddDocumentsNdjson(ctx context.Context, documents []byte, opts *DocumentOptions) (*TaskInfo, error) {
	// []byte avoids JSON conversion in Client.sendRequest()
	return i.addDocumentsFromReader(ctx, bytes.NewReader(documents), contentTypeNDJSON, transformDocumentOptionsToMap(opts))
}

func (i *index) AddDocumentsNdjsonInBatches(ctx context.Context, documents []byte, batchSize int, opts *DocumentOptions) ([]TaskInfo, error) {
	// Reuse io.Reader implementation
	return i.AddDocumentsNdjsonFromReaderInBatches(ctx, bytes.NewReader(documents), batchSize, opts)
}

func (i *index) AddDocumentsNdjsonFromReaderInBatches(ctx context.Context, documents io.Reader, batchSize int, opts *DocumentOptions) (resp []TaskInfo, err error) {
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

		resp, err := i.AddDocumentsNdjson(ctx, b.Bytes(), opts)
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

func (i *index) AddDocumentsNdjsonFromReader(ctx context.Context, documents io.Reader, opts *DocumentOptions) (resp *TaskInfo, err error) {
	// Using io.Reader would avoid JSON conversion in Client.sendRequest(), but
	// read content to memory anyway because of problems with streamed bodies
	data, err := io.ReadAll(documents)
	if err != nil {
		return nil, fmt.Errorf("could not read documents: %w", err)
	}
	return i.addDocuments(ctx, data, contentTypeNDJSON, transformDocumentOptionsToMap(opts))
}

func (i *index) UpdateDocuments(ctx context.Context, documentsPtr interface{}, opts *DocumentOptions) (*TaskInfo, error) {
	return i.updateDocuments(ctx, documentsPtr, contentTypeJSON, transformDocumentOptionsToMap(opts))
}

func (i *index) UpdateDocumentsInBatches(ctx context.Context, documentsPtr interface{}, batchSize int, opts *DocumentOptions) ([]TaskInfo, error) {
	return i.saveDocumentsInBatches(ctx, documentsPtr, batchSize, i.UpdateDocuments, opts)
}

func (i *index) UpdateDocumentsCsv(ctx context.Context, documents []byte, options *CsvDocumentsQuery) (*TaskInfo, error) {
	return i.updateDocuments(ctx, documents, contentTypeCSV, transformCsvDocumentsQueryToMap(options))
}

func (i *index) UpdateDocumentsCsvInBatches(ctx context.Context, documents []byte, batchSize int, options *CsvDocumentsQuery) ([]TaskInfo, error) {
	// Reuse io.Reader implementation
	return i.updateDocumentsCsvFromReaderInBatches(ctx, bytes.NewReader(documents), batchSize, options)
}

func (i *index) UpdateDocumentsNdjson(ctx context.Context, documents []byte, opts *DocumentOptions) (*TaskInfo, error) {
	return i.updateDocuments(ctx, documents, contentTypeNDJSON, transformDocumentOptionsToMap(opts))
}

func (i *index) UpdateDocumentsNdjsonInBatches(ctx context.Context, documents []byte, batchSize int, opts *DocumentOptions) ([]TaskInfo, error) {
	return i.updateDocumentsNdjsonFromReaderInBatches(ctx, bytes.NewReader(documents), batchSize, opts)
}

func (i *index) UpdateDocumentsByFunction(ctx context.Context, req *UpdateDocumentByFunctionRequest) (*TaskInfo, error) {
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
	if req.TaskCustomMetadata != "" {
		r.withQueryParams = map[string]string{
			"customMetadata": req.TaskCustomMetadata,
		}
	}
	if err := i.client.executeRequest(ctx, r); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetDocument(ctx context.Context, identifier string, request *DocumentQuery, documentPtr interface{}) error {
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

func (i *index) GetDocuments(ctx context.Context, param *DocumentsQuery, resp *DocumentsResult) error {
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

func (i *index) DeleteDocument(ctx context.Context, identifier string, opts *DocumentOptions) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/documents/" + identifier,
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		withQueryParams:     transformDocumentOptionsToMap(opts),
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "DeleteDocument",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) DeleteDocuments(ctx context.Context, identifiers []string, opts *DocumentOptions) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/documents/delete-batch",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         identifiers,
		withResponse:        resp,
		withQueryParams:     transformDocumentOptionsToMap(opts),
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "DeleteDocuments",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) DeleteDocumentsByFilter(ctx context.Context, filter interface{}, opts *DocumentOptions) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:    "/indexes/" + i.uid + "/documents/delete",
		method:      http.MethodPost,
		contentType: contentTypeJSON,
		withRequest: map[string]interface{}{
			"filter": filter,
		},
		withResponse:        resp,
		withQueryParams:     transformDocumentOptionsToMap(opts),
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "DeleteDocumentsByFilter",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, VersionErrorHintMessage(err, req)
	}
	return resp, nil
}

func (i *index) DeleteAllDocuments(ctx context.Context, opts *DocumentOptions) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/documents",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		withQueryParams:     transformDocumentOptionsToMap(opts),
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

func (i *index) saveDocumentsInBatches(ctx context.Context, documentsPtr interface{}, batchSize int, documentFunc func(ctx context.Context, documentsPtr interface{}, opts *DocumentOptions) (resp *TaskInfo, err error), opts *DocumentOptions) (resp []TaskInfo, err error) {
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

		respID, err := documentFunc(ctx, batch, opts)
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
	return i.saveDocumentsFromReaderInBatches(ctx, documents, batchSize, i.UpdateDocumentsCsv, options)
}

func (i *index) updateDocumentsNdjsonFromReaderInBatches(ctx context.Context, documents io.Reader, batchSize int, opts *DocumentOptions) (resp []TaskInfo, err error) {
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

		resp, err := i.UpdateDocumentsNdjson(ctx, b.Bytes(), opts)
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
