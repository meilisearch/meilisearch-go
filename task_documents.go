package meilisearch

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/andybalholm/brotli"
)

const (
	taskDocumentsContentType      = "application/x-ndjson"
	taskDocumentsMaxScanTokenSize = 64 << 20
)

func (m *meilisearch) GetTaskDocuments(taskUID int64, dst interface{}) error {
	return m.GetTaskDocumentsWithContext(context.Background(), taskUID, dst)
}

func (m *meilisearch) GetTaskDocumentsWithContext(ctx context.Context, taskUID int64, dst interface{}) error {
	sliceValue, sliceElemType, err := validateTaskDocumentsDestination(dst)
	if err != nil {
		return err
	}

	req := &internalRequest{
		endpoint:             "/tasks/" + strconv.FormatInt(taskUID, 10) + "/documents",
		method:               http.MethodGet,
		withRequest:          nil,
		withResponse:         nil,
		withQueryParams:      nil,
		withResponseEncoding: true,
		acceptedStatusCodes:  []int{http.StatusOK},
		functionName:         "GetTaskDocuments",
	}
	internalError := &Error{
		Endpoint:         req.endpoint,
		Method:           req.method,
		Function:         req.functionName,
		RequestToString:  "empty request",
		ResponseToString: "empty response",
		MeilisearchApiError: meilisearchApiError{
			Message: "empty meilisearch message",
		},
		StatusCodeExpected: req.acceptedStatusCodes,
	}

	resp, err := m.client.sendRequest(ctx, req, internalError)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	internalError.StatusCode = resp.StatusCode
	if resp.StatusCode != http.StatusOK {
		if responseUsesClientEncoding(resp, m.client.contentEncoding) {
			internalError.encoder = m.client.encoder
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return m.client.handleStatusCode(req, resp.StatusCode, body, internalError)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(contentType)), taskDocumentsContentType) {
		return fmt.Errorf("GetTaskDocuments: expected Content-Type to start with %q, got %q", taskDocumentsContentType, contentType)
	}

	body, closeBody, err := taskDocumentsResponseBody(resp)
	if err != nil {
		return err
	}
	if closeBody {
		defer func() {
			_ = body.Close()
		}()
	}

	result := sliceValue
	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 0, 64<<10), taskDocumentsMaxScanTokenSize)
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}

		elemPtr := reflect.New(sliceElemType)
		if err := m.client.jsonUnmarshal(line, elemPtr.Interface()); err != nil {
			return fmt.Errorf("GetTaskDocuments: failed to unmarshal NDJSON line: %w", err)
		}
		result = reflect.Append(result, elemPtr.Elem())
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("GetTaskDocuments: failed to read NDJSON response: %w", err)
	}

	sliceValue.Set(result)
	return nil
}

func validateTaskDocumentsDestination(dst interface{}) (reflect.Value, reflect.Type, error) {
	if dst == nil {
		return reflect.Value{}, nil, fmt.Errorf("GetTaskDocuments: dst must be a non-nil pointer to a slice")
	}

	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr || dstValue.IsNil() {
		return reflect.Value{}, nil, fmt.Errorf("GetTaskDocuments: dst must be a non-nil pointer to a slice")
	}

	sliceValue := dstValue.Elem()
	if sliceValue.Kind() != reflect.Slice {
		return reflect.Value{}, nil, fmt.Errorf("GetTaskDocuments: dst must point to a slice, got %s", sliceValue.Kind())
	}

	return sliceValue, sliceValue.Type().Elem(), nil
}

func taskDocumentsResponseBody(resp *http.Response) (io.ReadCloser, bool, error) {
	contentEncoding := strings.ToLower(strings.TrimSpace(resp.Header.Get("Content-Encoding")))
	switch ContentEncoding(contentEncoding) {
	case "":
		return resp.Body, false, nil
	case GzipEncoding:
		body, err := gzip.NewReader(resp.Body)
		return body, true, err
	case DeflateEncoding:
		body, err := zlib.NewReader(resp.Body)
		return body, true, err
	case BrotliEncoding:
		return io.NopCloser(brotli.NewReader(resp.Body)), true, nil
	default:
		return nil, false, fmt.Errorf("GetTaskDocuments: unsupported Content-Encoding %q", contentEncoding)
	}
}

func responseUsesClientEncoding(resp *http.Response, contentEncoding ContentEncoding) bool {
	if contentEncoding.IsZero() {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(resp.Header.Get("Content-Encoding")), contentEncoding.String())
}
