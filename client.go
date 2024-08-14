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
)

type client struct {
	client     *http.Client
	host       string
	apiKey     string
	bufferPool *sync.Pool
}

type internalRequest struct {
	endpoint    string
	method      string
	contentType string

	withRequest     interface{}
	withResponse    interface{}
	withQueryParams map[string]string

	acceptedStatusCodes []int

	functionName string
}

func newClient(cli *http.Client, host, apiKey string) *client {
	return &client{
		client: cli,
		host:   host,
		apiKey: apiKey,
		bufferPool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
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
		if b, ok := rawRequest.([]byte); ok {
			// If the request body is already a []byte then use it directly
			buf := c.bufferPool.Get().(*bytes.Buffer)
			buf.Reset()
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
					return nil, internalError.WithErrCode(ErrCodeMarshalRequest, fmt.Errorf("failed to marshal with MarshalJSON: %w", err))
				}
				if data == nil {
					return nil, internalError.WithErrCode(ErrCodeMarshalRequest, errors.New("MarshalJSON returned nil data"))
				}
			} else {
				data, err = json.Marshal(rawRequest)
				if err != nil {
					return nil, internalError.WithErrCode(ErrCodeMarshalRequest, fmt.Errorf("failed to marshal with json.Marshal: %w", err))
				}
			}
			buf := c.bufferPool.Get().(*bytes.Buffer)
			buf.Reset()
			buf.Write(data)
			body = buf
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

	request.Header.Set("User-Agent", GetQualifiedVersion())

	resp, err := c.client.Do(request)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, internalError.WithErrCode(MeilisearchTimeoutError, err)
		}
		return nil, internalError.WithErrCode(MeilisearchCommunicationError, err)
	}

	if body != nil {
		if buf, ok := body.(*bytes.Buffer); ok {
			c.bufferPool.Put(buf)
		}
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

		internalError.ResponseToString = string(body)

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
	return nil
}
