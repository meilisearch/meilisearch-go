package meilisearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"
)

type client struct {
	client          *http.Client
	host            string
	apiKey          string
	bufferPool      *sync.Pool
	encoder         encoder
	contentEncoding ContentEncoding
	retryOnStatus   map[int]bool
	disableRetry    bool
	maxRetries      uint8
	retryBackoff    func(attempt uint8) time.Duration

	jsonMarshal   JSONMarshal
	jsonUnmarshal JSONUnmarshal
}

type clientConfig struct {
	contentEncoding          ContentEncoding
	encodingCompressionLevel EncodingCompressionLevel
	retryOnStatus            map[int]bool
	disableRetry             bool
	maxRetries               uint8
	jsonMarshal              JSONMarshal
	jsonUnmarshal            JSONUnmarshal
}

type internalRequest struct {
	endpoint             string
	method               string
	contentType          string
	withRequest          interface{}
	withResponse         interface{}
	withQueryParams      map[string]string
	withResponseEncoding bool

	acceptedStatusCodes []int
	acceptedContentType string

	functionName string
}

func newClient(cli *http.Client, host, apiKey string, cfg *clientConfig) *client {
	c := &client{
		client: cli,
		host:   host,
		apiKey: apiKey,
		bufferPool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		disableRetry:  cfg.disableRetry,
		maxRetries:    cfg.maxRetries,
		retryOnStatus: cfg.retryOnStatus,
		jsonMarshal:   cfg.jsonMarshal,
		jsonUnmarshal: cfg.jsonUnmarshal,
	}

	if c.retryOnStatus == nil {
		c.retryOnStatus = map[int]bool{
			502: true,
			503: true,
			504: true,
		}
	}

	if !c.disableRetry && c.retryBackoff == nil {
		c.retryBackoff = func(attempt uint8) time.Duration {
			return time.Second * time.Duration(attempt)
		}
	}

	if !cfg.contentEncoding.IsZero() {
		c.contentEncoding = cfg.contentEncoding
		c.encoder = newEncoding(cfg.contentEncoding, cfg.encodingCompressionLevel)
	}

	return c
}

func (c *client) executeRequest(ctx context.Context, req *internalRequest) error {
	if req.acceptedContentType == contentTypeNDJSON && req.withResponse != nil {
		if _, _, err := validateNDJSONDestination(req.functionName, req.withResponse); err != nil {
			return err
		}
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
		encoder:            c.encoder,
	}

	resp, err := c.sendRequest(ctx, req, internalError)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	internalError.StatusCode = resp.StatusCode

	if req.acceptedContentType == contentTypeNDJSON && req.withResponse != nil {
		return c.handleNDJSONResponse(req, resp, internalError)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = c.handleStatusCode(req, resp.StatusCode, b, internalError)
	if err != nil {
		return err
	}

	err = c.handleContentType(req, resp, internalError)
	if err != nil {
		return err
	}

	err = c.handleResponse(req, b, internalError)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) handleNDJSONResponse(req *internalRequest, resp *http.Response, internalError *Error) error {
	sliceValue, sliceElemType, err := validateNDJSONDestination(req.functionName, req.withResponse)
	if err != nil {
		return err
	}

	if err := c.handleStreamingStatusCode(req, resp, internalError); err != nil {
		return err
	}

	if err := c.handleContentType(req, resp, internalError); err != nil {
		return err
	}

	result := sliceValue.Slice(0, 0)
	dec, err := c.responseDecoder(resp)
	if err != nil {
		return fmt.Errorf("%s: failed to create response decoder: %w", req.functionName, err)
	}
	defer func() {
		_ = dec.Close()
	}()

	for {
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("%s: failed to decode NDJSON: %w", req.functionName, err)
		}

		elemPtr := reflect.New(sliceElemType)
		if err := c.jsonUnmarshal(raw, elemPtr.Interface()); err != nil {
			return fmt.Errorf("%s: failed to unmarshal NDJSON document: %w", req.functionName, err)
		}
		result = reflect.Append(result, elemPtr.Elem())
	}

	sliceValue.Set(result)
	return nil
}

func (c *client) handleStreamingStatusCode(req *internalRequest, resp *http.Response, internalError *Error) error {
	if req.acceptedStatusCodes == nil {
		return nil
	}

	for _, acceptedCode := range req.acceptedStatusCodes {
		if resp.StatusCode == acceptedCode {
			return nil
		}
	}

	if responseUsesClientEncoding(resp, c.contentEncoding) {
		internalError.encoder = c.encoder
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return c.handleStatusCode(req, resp.StatusCode, body, internalError)
}

func (c *client) handleContentType(req *internalRequest, resp *http.Response, internalError *Error) error {
	if req.acceptedContentType == "" {
		return nil
	}

	if err := validateContentType(req.functionName, req.acceptedContentType, resp.Header.Get("Content-Type")); err != nil {
		return internalError.WithErrCode(ErrCodeResponseUnmarshalBody, err)
	}
	return nil
}

func validateNDJSONDestination(functionName string, dst interface{}) (reflect.Value, reflect.Type, error) {
	if dst == nil {
		return reflect.Value{}, nil, fmt.Errorf("%s: dst must be a non-nil pointer to a slice", functionName)
	}

	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr || dstValue.IsNil() {
		return reflect.Value{}, nil, fmt.Errorf("%s: dst must be a non-nil pointer to a slice", functionName)
	}

	sliceValue := dstValue.Elem()
	if sliceValue.Kind() != reflect.Slice {
		return reflect.Value{}, nil, fmt.Errorf("%s: dst must point to a slice, got %s", functionName, sliceValue.Kind())
	}

	return sliceValue, sliceValue.Type().Elem(), nil
}

func validateContentType(functionName, expectedContentType, contentType string) error {
	normalizedExpected := strings.ToLower(strings.TrimSpace(expectedContentType))
	normalized := strings.ToLower(strings.TrimSpace(contentType))
	if strings.HasPrefix(normalized, normalizedExpected) {
		return nil
	}
	return fmt.Errorf("%s: unexpected Content-Type %q, expected %q", functionName, contentType, expectedContentType)
}

func (c *client) responseDecoder(resp *http.Response) (streamDecoder, error) {
	contentEncoding := strings.ToLower(strings.TrimSpace(resp.Header.Get("Content-Encoding")))
	ce := ContentEncoding(contentEncoding)
	switch ce {
	case "":
		return &jsonStreamDecoder{Decoder: json.NewDecoder(resp.Body)}, nil
	case GzipEncoding, DeflateEncoding, BrotliEncoding:
		encoder := c.encoder
		if encoder == nil || ce != c.contentEncoding {
			encoder = newEncoding(ce, DefaultCompression)
		}
		return encoder.Decoder(resp.Body)
	default:
		return nil, fmt.Errorf("unsupported Content-Encoding %q", contentEncoding)
	}
}

func responseUsesClientEncoding(resp *http.Response, contentEncoding ContentEncoding) bool {
	if contentEncoding.IsZero() {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(resp.Header.Get("Content-Encoding")), contentEncoding.String())
}

func (c *client) sendRequest(
	ctx context.Context,
	req *internalRequest,
	internalError *Error,
) (*http.Response, error) {

	apiURL, err := url.Parse(c.host + req.endpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to parse url: %w", err)
	}

	if req.withQueryParams != nil {
		query := apiURL.Query()
		for key, value := range req.withQueryParams {
			query.Set(key, value)
		}

		apiURL.RawQuery = query.Encode()
	}

	body, err := c.buildBody(req, internalError)
	if err != nil {
		return nil, err
	}

	// Create the HTTP request
	request, err := http.NewRequestWithContext(ctx, req.method, apiURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	// adding request headers
	if req.contentType != "" {
		request.Header.Set("Content-Type", req.contentType)
	}
	if c.apiKey != "" {
		request.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	if (req.withResponse != nil || req.withResponseEncoding) && !c.contentEncoding.IsZero() {
		request.Header.Set("Accept-Encoding", c.contentEncoding.String())
	}

	if req.withRequest != nil && !c.contentEncoding.IsZero() {
		request.Header.Set("Content-Encoding", c.contentEncoding.String())
	}

	request.Header.Set("User-Agent", GetQualifiedVersion())

	resp, err := c.do(request, internalError)
	if err != nil {
		if rc, ok := body.(io.Closer); ok {
			_ = rc.Close()
		}
		return nil, err
	}

	if rc, ok := body.(io.Closer); ok {
		_ = rc.Close()
	}

	return resp, nil
}

func (c *client) do(req *http.Request, internalError *Error) (resp *http.Response, err error) {
	retriesCount := uint8(0)

	for {
		resp, err = c.client.Do(req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return nil, internalError.WithErrCode(MeilisearchTimeoutError, err)
			}
			return nil, internalError.WithErrCode(MeilisearchCommunicationError, err)
		}

		// Exit if retries are disabled
		if c.disableRetry {
			break
		}

		// Check if response status is retryable and we haven't exceeded max retries
		if c.retryOnStatus[resp.StatusCode] && retriesCount < c.maxRetries {
			retriesCount++

			// Close response body to prevent memory leaks
			_ = resp.Body.Close()

			// Handle backoff with context cancellation support
			backoff := c.retryBackoff(retriesCount)
			timer := time.NewTimer(backoff)

			select {
			case <-req.Context().Done():
				err := req.Context().Err()
				timer.Stop()
				return nil, internalError.WithErrCode(MeilisearchTimeoutError, err)
			case <-timer.C:
				// Retry after backoff
				timer.Stop()
			}

			continue
		}

		break
	}

	// Return error if retries exceeded the maximum limit
	if !c.disableRetry && retriesCount >= c.maxRetries {
		return nil, internalError.WithErrCode(MeilisearchMaxRetriesExceeded, nil)
	}

	return resp, nil
}

func (c *client) handleStatusCode(req *internalRequest, statusCode int, body []byte, internalError *Error) error {
	if req.acceptedStatusCodes != nil {

		// A successful status code is required so check if the response status code is in the
		// expected status code list.
		for _, acceptedCode := range req.acceptedStatusCodes {
			if statusCode == acceptedCode {
				return nil
			}
		}

		internalError.ErrorBody(body)

		if internalError.MeilisearchApiError.Code == "" {
			return internalError.WithErrCode(MeilisearchApiErrorWithoutMessage)
		}
		return internalError.WithErrCode(MeilisearchApiError)
	}

	return nil
}

func (c *client) handleResponse(req *internalRequest, body []byte, internalError *Error) (err error) {
	if req.withResponse != nil {
		if !c.contentEncoding.IsZero() {
			if err := c.encoder.Decode(body, req.withResponse); err != nil {
				return internalError.WithErrCode(ErrCodeResponseUnmarshalBody, err)
			}
		} else {
			internalError.ResponseToString = string(body)
			if internalError.ResponseToString == nullBody {
				req.withResponse = nil
				return nil
			}

			if err := c.jsonUnmarshal(body, req.withResponse); err != nil {
				return internalError.WithErrCode(ErrCodeResponseUnmarshalBody, err)
			}
		}
	}
	return nil
}

func (c *client) buildBody(req *internalRequest, internalError *Error) (io.ReadCloser, error) {
	var body io.Reader
	var buf *bytes.Buffer
	var bufFromPool bool

	if req.withRequest != nil {
		if req.method == http.MethodGet || req.method == http.MethodHead {
			return nil, ErrInvalidRequestMethod
		}
		if req.contentType == "" {
			return nil, ErrRequestBodyWithoutContentType
		}

		switch v := req.withRequest.(type) {
		case io.ReadCloser:
			body = v
		case io.Reader:
			body = io.NopCloser(v)
		case []byte:
			body = io.NopCloser(bytes.NewReader(v))
		default:
			buf = c.bufferPool.Get().(*bytes.Buffer)
			buf.Reset()
			bufFromPool = true

			data, err := c.jsonMarshal(req.withRequest)
			if err != nil {
				c.bufferPool.Put(buf)
				return nil, internalError.WithErrCode(ErrCodeMarshalRequest,
					fmt.Errorf("failed to marshal request with json.Marshal: %w", err))
			}
			buf.Write(data)
			body = buf
		}

		if !c.contentEncoding.IsZero() {
			compressedBody, err := c.encoder.Encode(body)
			if bufFromPool {
				c.bufferPool.Put(buf)
			}
			if err != nil {
				return nil, internalError.WithErrCode(ErrCodeMarshalRequest,
					fmt.Errorf("failed to encode request body: %w", err))
			}
			return compressedBody, nil
		}
	}

	// If not compressed, make sure body is io.ReadCloser
	switch b := body.(type) {
	case nil:
		return nil, nil
	case io.ReadCloser:
		return b, nil
	default:
		return io.NopCloser(b), nil
	}
}
