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
}

type clientConfig struct {
	contentEncoding          ContentEncoding
	encodingCompressionLevel EncodingCompressionLevel
	retryOnStatus            map[int]bool
	disableRetry             bool
	maxRetries               uint8
}

type internalRequest struct {
	endpoint        string
	method          string
	contentType     string
	withRequest     interface{}
	withResponse    interface{}
	withQueryParams map[string]string

	acceptedStatusCodes []int

	functionName string
}

func newClient(cli *http.Client, host, apiKey string, cfg clientConfig) *client {
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

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = c.handleStatusCode(req, resp.StatusCode, b, internalError)
	if err != nil {
		return err
	}

	err = c.handleResponse(req, b, internalError)
	if err != nil {
		return err
	}
	return nil
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

	// Create request body
	var body io.Reader = nil
	if req.withRequest != nil {
		if req.method == http.MethodGet || req.method == http.MethodHead {
			return nil, ErrInvalidRequestMethod
		}
		if req.contentType == "" {
			return nil, ErrRequestBodyWithoutContentType
		}

		rawRequest := req.withRequest

		buf := c.bufferPool.Get().(*bytes.Buffer)
		buf.Reset()

		if b, ok := rawRequest.([]byte); ok {
			buf.Write(b)
			body = buf
		} else if reader, ok := rawRequest.(io.Reader); ok {
			// If the request body is an io.Reader then stream it directly
			body = reader
		} else {
			// Otherwise convert it to JSON
			var (
				data []byte
				err  error
			)
			if marshaler, ok := rawRequest.(json.Marshaler); ok {
				data, err = marshaler.MarshalJSON()
				if err != nil {
					return nil, internalError.WithErrCode(ErrCodeMarshalRequest,
						fmt.Errorf("failed to marshal with MarshalJSON: %w", err))
				}
				if data == nil {
					return nil, internalError.WithErrCode(ErrCodeMarshalRequest,
						errors.New("MarshalJSON returned nil data"))
				}
			} else {
				data, err = json.Marshal(rawRequest)
				if err != nil {
					return nil, internalError.WithErrCode(ErrCodeMarshalRequest,
						fmt.Errorf("failed to marshal with json.Marshal: %w", err))
				}
			}
			buf.Write(data)
			body = buf
		}

		if !c.contentEncoding.IsZero() {
			// Get the data from the buffer before encoding
			var bufData []byte
			if buf, ok := body.(*bytes.Buffer); ok {
				bufData = buf.Bytes()
				encodedBuf, err := c.encoder.Encode(bytes.NewReader(bufData))
				if err != nil {
					if buf, ok := body.(*bytes.Buffer); ok {
						c.bufferPool.Put(buf)
					}
					return nil, internalError.WithErrCode(ErrCodeMarshalRequest,
						fmt.Errorf("failed to encode request body: %w", err))
				}
				// Return the original buffer to the pool since we have a new one
				if buf, ok := body.(*bytes.Buffer); ok {
					c.bufferPool.Put(buf)
				}
				body = encodedBuf
			}
		}
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

	if req.withResponse != nil && !c.contentEncoding.IsZero() {
		request.Header.Set("Accept-Encoding", c.contentEncoding.String())
	}

	if req.withRequest != nil && !c.contentEncoding.IsZero() {
		request.Header.Set("Content-Encoding", c.contentEncoding.String())
	}

	request.Header.Set("User-Agent", GetQualifiedVersion())

	resp, err := c.do(request, internalError)
	if err != nil {
		return nil, err
	}

	if body != nil {
		if buf, ok := body.(*bytes.Buffer); ok {
			c.bufferPool.Put(buf)
		}
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
			resp.Body.Close()

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

			var err error
			if resp, ok := req.withResponse.(json.Unmarshaler); ok {
				err = resp.UnmarshalJSON(body)
				req.withResponse = resp
			} else {
				err = json.Unmarshal(body, req.withResponse)
			}
			if err != nil {
				return internalError.WithErrCode(ErrCodeResponseUnmarshalBody, err)
			}
		}
	}
	return nil
}
