package meilisearch

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func sendCsvRecords(documentsCsvFunc func(recs []byte, pk ...string) (resp *TaskInfo, err error), records [][]string, primaryKey ...string) (*TaskInfo, error) {
	b := new(bytes.Buffer)
	w := csv.NewWriter(b)
	w.UseCRLF = true

	err := w.WriteAll(records)
	if err != nil {
		return nil, fmt.Errorf("could not write CSV records: %w", err)
	}

	resp, err := documentsCsvFunc(b.Bytes(), primaryKey...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) saveDocumentsFromReaderInBatches(documents io.Reader, batchSize int, documentsCsvFunc func(recs []byte, pk ...string) (resp *TaskInfo, err error), primaryKey ...string) (resp []TaskInfo, err error) {
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
			resp, err := sendCsvRecords(documentsCsvFunc, records, primaryKey...)
			if err != nil {
				return nil, err
			}
			responses = append(responses, *resp)
			records = nil
		}
	}

	// Send remaining records as the last batch if there is any
	if len(records) > 0 {
		resp, err := sendCsvRecords(documentsCsvFunc, records, primaryKey...)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *resp)
	}

	return responses, nil
}

func (i Index) saveDocumentsInBatches(documentsPtr interface{}, batchSize int, documentFunc func(documentsPtr interface{}, primaryKey ...string) (resp *TaskInfo, err error), primaryKey ...string) (resp []TaskInfo, err error) {
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

		if len(primaryKey) != 0 {
			respID, err := documentFunc(batch, primaryKey[0])
			if err != nil {
				return nil, err
			}

			resp[j] = *respID
		} else {
			respID, err := documentFunc(batch)
			if err != nil {
				return nil, err
			}

			resp[j] = *respID
		}
	}

	return resp, nil
}

func (i Index) GetDocument(identifier string, request *DocumentQuery, documentPtr interface{}) error {
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/documents/" + identifier,
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
	}
	if err := i.client.executeRequest(req); err != nil {
		return err
	}
	return nil
}

func (i Index) GetDocuments(request *DocumentsQuery, resp *DocumentsResult) error {
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/documents",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		withQueryParams:     map[string]string{},
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetDocuments",
	}
	if request != nil {
		if request.Limit != 0 {
			req.withQueryParams["limit"] = strconv.FormatInt(request.Limit, 10)
		}
		if request.Offset != 0 {
			req.withQueryParams["offset"] = strconv.FormatInt(request.Offset, 10)
		}
		if len(request.Fields) != 0 {
			req.withQueryParams["fields"] = strings.Join(request.Fields, ",")
		}
	}
	if err := i.client.executeRequest(req); err != nil {
		return err
	}
	return nil
}

func (i Index) addDocuments(documentsPtr interface{}, contentType string, primaryKey ...string) (resp *TaskInfo, err error) {
	resp = &TaskInfo{}
	endpoint := ""
	if len(primaryKey) == 0 {
		endpoint = "/indexes/" + i.UID + "/documents"
	} else {
		i.PrimaryKey = primaryKey[0] //nolint:golint,staticcheck
		endpoint = "/indexes/" + i.UID + "/documents?primaryKey=" + primaryKey[0]
	}
	req := internalRequest{
		endpoint:            endpoint,
		method:              http.MethodPost,
		contentType:         contentType,
		withRequest:         documentsPtr,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "AddDocuments",
	}
	if err = i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) AddDocuments(documentsPtr interface{}, primaryKey ...string) (resp *TaskInfo, err error) {
	return i.addDocuments(documentsPtr, contentTypeJSON, primaryKey...)
}

func (i Index) AddDocumentsInBatches(documentsPtr interface{}, batchSize int, primaryKey ...string) (resp []TaskInfo, err error) {
	return i.saveDocumentsInBatches(documentsPtr, batchSize, i.AddDocuments, primaryKey...)
}

func (i Index) AddDocumentsCsv(documents []byte, primaryKey ...string) (resp *TaskInfo, err error) {
	// []byte avoids JSON conversion in Client.sendRequest()
	return i.addDocuments(documents, contentTypeCSV, primaryKey...)
}

func (i Index) AddDocumentsCsvFromReader(documents io.Reader, primaryKey ...string) (resp *TaskInfo, err error) {
	// Using io.Reader would avoid JSON conversion in Client.sendRequest(), but
	// read content to memory anyway because of problems with streamed bodies
	data, err := io.ReadAll(documents)
	if err != nil {
		return nil, fmt.Errorf("could not read documents: %w", err)
	}
	return i.addDocuments(data, contentTypeCSV, primaryKey...)
}

func (i Index) AddDocumentsCsvInBatches(documents []byte, batchSize int, primaryKey ...string) (resp []TaskInfo, err error) {
	// Reuse io.Reader implementation
	return i.AddDocumentsCsvFromReaderInBatches(bytes.NewReader(documents), batchSize, primaryKey...)
}

func (i Index) AddDocumentsCsvFromReaderInBatches(documents io.Reader, batchSize int, primaryKey ...string) (resp []TaskInfo, err error) {
	return i.saveDocumentsFromReaderInBatches(documents, batchSize, i.AddDocumentsCsv, primaryKey...)
}

func (i Index) AddDocumentsNdjson(documents []byte, primaryKey ...string) (resp *TaskInfo, err error) {
	// []byte avoids JSON conversion in Client.sendRequest()
	return i.addDocuments([]byte(documents), contentTypeNDJSON, primaryKey...)
}

func (i Index) AddDocumentsNdjsonFromReader(documents io.Reader, primaryKey ...string) (resp *TaskInfo, err error) {
	// Using io.Reader would avoid JSON conversion in Client.sendRequest(), but
	// read content to memory anyway because of problems with streamed bodies
	data, err := io.ReadAll(documents)
	if err != nil {
		return nil, fmt.Errorf("could not read documents: %w", err)
	}
	return i.addDocuments(data, contentTypeNDJSON, primaryKey...)
}

func (i Index) AddDocumentsNdjsonInBatches(documents []byte, batchSize int, primaryKey ...string) (resp []TaskInfo, err error) {
	// Reuse io.Reader implementation
	return i.AddDocumentsNdjsonFromReaderInBatches(bytes.NewReader(documents), batchSize, primaryKey...)
}

func (i Index) AddDocumentsNdjsonFromReaderInBatches(documents io.Reader, batchSize int, primaryKey ...string) (resp []TaskInfo, err error) {
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

		resp, err := i.AddDocumentsNdjson(b.Bytes(), primaryKey...)
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

func (i Index) updateDocuments(documentsPtr interface{}, contentType string, primaryKey ...string) (resp *TaskInfo, err error) {
	resp = &TaskInfo{}
	endpoint := ""
	if len(primaryKey) == 0 {
		endpoint = "/indexes/" + i.UID + "/documents"
	} else {
		i.PrimaryKey = primaryKey[0] //nolint:golint,staticcheck
		endpoint = "/indexes/" + i.UID + "/documents?primaryKey=" + primaryKey[0]
	}
	req := internalRequest{
		endpoint:            endpoint,
		method:              http.MethodPut,
		contentType:         contentType,
		withRequest:         documentsPtr,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateDocuments",
	}
	if err = i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) UpdateDocuments(documentsPtr interface{}, primaryKey ...string) (resp *TaskInfo, err error) {
	return i.updateDocuments(documentsPtr, contentTypeJSON, primaryKey...)
}

func (i Index) UpdateDocumentsInBatches(documentsPtr interface{}, batchSize int, primaryKey ...string) (resp []TaskInfo, err error) {
	return i.saveDocumentsInBatches(documentsPtr, batchSize, i.UpdateDocuments, primaryKey...)
}

func (i Index) UpdateDocumentsCsv(documents []byte, primaryKey ...string) (resp *TaskInfo, err error) {
	return i.updateDocuments(documents, contentTypeCSV, primaryKey...)
}

func (i Index) UpdateDocumentsCsvFromReader(documents io.Reader, primaryKey ...string) (resp *TaskInfo, err error) {
	// Using io.Reader would avoid JSON conversion in Client.sendRequest(), but
	// read content to memory anyway because of problems with streamed bodies
	data, err := io.ReadAll(documents)
	if err != nil {
		return nil, fmt.Errorf("could not read documents: %w", err)
	}
	return i.updateDocuments(data, contentTypeCSV, primaryKey...)
}

func (i Index) UpdateDocumentsCsvInBatches(documents []byte, batchSize int, primaryKey ...string) (resp []TaskInfo, err error) {
	// Reuse io.Reader implementation
	return i.UpdateDocumentsCsvFromReaderInBatches(bytes.NewReader(documents), batchSize, primaryKey...)
}

func (i Index) UpdateDocumentsCsvFromReaderInBatches(documents io.Reader, batchSize int, primaryKey ...string) (resp []TaskInfo, err error) {
	return i.saveDocumentsFromReaderInBatches(documents, batchSize, i.UpdateDocumentsCsv, primaryKey...)
}

func (i Index) UpdateDocumentsNdjson(documents []byte, primaryKey ...string) (resp *TaskInfo, err error) {
	return i.updateDocuments(documents, contentTypeNDJSON, primaryKey...)
}

func (i Index) UpdateDocumentsNdjsonFromReader(documents io.Reader, primaryKey ...string) (resp *TaskInfo, err error) {
	// Using io.Reader would avoid JSON conversion in Client.sendRequest(), but
	// read content to memory anyway because of problems with streamed bodies
	data, err := io.ReadAll(documents)
	if err != nil {
		return nil, fmt.Errorf("could not read documents: %w", err)
	}
	return i.updateDocuments(data, contentTypeNDJSON, primaryKey...)
}

func (i Index) UpdateDocumentsNdjsonInBatches(documents []byte, batchsize int, primaryKey ...string) (resp []TaskInfo, err error) {
	return i.updateDocumentsNdjsonFromReaderInBatches(bytes.NewReader(documents), batchsize, primaryKey...)
}

func (i Index) updateDocumentsNdjsonFromReaderInBatches(documents io.Reader, batchSize int, primaryKey ...string) (resp []TaskInfo, err error) {
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
				return nil, fmt.Errorf("Could not write NDJSON line: %w", err)
			}
			err = b.WriteByte('\n')
			if err != nil {
				return nil, fmt.Errorf("Could not write NDJSON line: %w", err)
			}
		}

		resp, err := i.UpdateDocumentsNdjson(b.Bytes(), primaryKey...)
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
		return nil, fmt.Errorf("Could not read NDJSON: %w", err)
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

func (i Index) DeleteDocument(identifier string) (resp *TaskInfo, err error) {
	resp = &TaskInfo{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/documents/" + identifier,
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "DeleteDocument",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) DeleteDocuments(identifier []string) (resp *TaskInfo, err error) {
	resp = &TaskInfo{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/documents/delete-batch",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         identifier,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "DeleteDocuments",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) DeleteAllDocuments() (resp *TaskInfo, err error) {
	resp = &TaskInfo{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/documents",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "DeleteAllDocuments",
	}
	if err = i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}
