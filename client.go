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
	"time"
)

type client struct {
	client          *http.Client
	host            string
	apiKey          string
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
		client:        cli,
		host:          host,
		apiKey:        apiKey,
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
		APIError: meilisearchApiError{
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
	sliceValue, _, err := validateNDJSONDestination(req.functionName, req.withResponse)
	if err != nil {
		return err
	}

	if err := c.handleStreamingStatusCode(req, resp, internalError); err != nil {
		return err
	}

	if err := c.handleContentType(req, resp, internalError); err != nil {
		return err
	}

	// The Meilisearch API returns concatenated JSON values (akin to NDJSON
	// but the server does not stream them). Read every value through the
	// response decoder up front and assemble a single JSON array, then
	// unmarshal once into dst so the SDK does not expose streaming
	// semantics to its callers.
	dec, err := c.responseDecoder(resp)
	if err != nil {
		return fmt.Errorf("%s: failed to create response decoder: %w", req.functionName, err)
	}
	defer func() {
		_ = dec.Close()
	}()

	values := make([]json.RawMessage, 0)
	totalBytes := 2

	for {
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("%s: failed to decode NDJSON: %w", req.functionName, err)
		}
		values = append(values, raw)
		totalBytes += len(raw) + 1
	}

	if len(values) == 0 {
		sliceValue.Set(sliceValue.Slice(0, 0))
		return nil
	}

	var buf bytes.Buffer
	buf.Grow(totalBytes)
	buf.WriteByte('[')
	for i, v := range values {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.Write(v)
	}
	buf.WriteByte(']')

	if err := c.jsonUnmarshal(buf.Bytes(), sliceValue.Addr().Interface()); err != nil {
		return fmt.Errorf("%s: failed to unmarshal NDJSON response: %w", req.functionName, err)
	}
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

	var bodyBytes []byte
	if body != nil {
		bodyBytes, err = io.ReadAll(body)
		// Close the original body (handles sync.Pool internally if you implement a custom closer, or let GC handle it)
		_ = body.Close()
		if err != nil {
			return nil, fmt.Errorf("unable to read body: %w", err)
		}
		body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	request, err := http.NewRequestWithContext(ctx, req.method, apiURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	if bodyBytes != nil {
		request.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(bodyBytes)), nil
		}

		request.ContentLength = int64(len(bodyBytes))
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
		if retriesCount > 0 && req.GetBody != nil {
			newBody, bodyErr := req.GetBody()
			if bodyErr != nil {
				return nil, internalError.WithErrCode(CommunicationError,
					fmt.Errorf("failed to rewind body on retry: %w", bodyErr))
			}
			req.Body = newBody
		}

		resp, err = c.client.Do(req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return nil, internalError.WithErrCode(TimeoutError, err)
			}
			if errors.Is(err, context.Canceled) {
				return nil, internalError.WithErrCode(TimeoutError, err)
			}
			return nil, internalError.WithErrCode(CommunicationError, err)
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
				return nil, internalError.WithErrCode(TimeoutError, err)
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
		return nil, internalError.WithErrCode(MaxRetriesExceeded, nil)
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

		if internalError.APIError.Code == "" {
			return internalError.WithErrCode(APIErrorWithoutMessage)
		}
		return internalError.WithErrCode(APIError)
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
	var body io.ReadCloser

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
			data, err := c.jsonMarshal(req.withRequest)
			if err != nil {
				return nil, internalError.WithErrCode(ErrCodeMarshalRequest,
					fmt.Errorf("failed to marshal request with json.Marshal: %w", err))
			}
			body = io.NopCloser(bytes.NewReader(data))
		}

		if !c.contentEncoding.IsZero() {
			compressedBody, err := c.encoder.Encode(body)
			_ = body.Close()

			if err != nil {
				return nil, internalError.WithErrCode(ErrCodeMarshalRequest,
					fmt.Errorf("failed to encode request body: %w", err))
			}
			return compressedBody, nil
		}
	}

	return body, nil
}
