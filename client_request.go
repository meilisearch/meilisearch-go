package meilisearch

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/valyala/fasthttp"

	"encoding/json"
)

const (
	contentTypeJSON   string = "application/json"
	contentTypeNDJSON string = "application/x-ndjson"
	contentTypeCSV    string = "text/csv"
)

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

func (c *Client) executeRequest(req internalRequest) error {
	return c.requestExecutor.executeRequest(req, c)
}

type requestExecutor interface {
	executeRequest(req internalRequest, client *Client) error
}

// fasthttpRequestExecutor is a requestExecutor implementation using fasthttp
type fasthttpRequestExecutor struct {
	httpClient *fasthttp.Client
}

var _ requestExecutor = &fasthttpRequestExecutor{}

// executeRequest implements requestExecutor interface
func (c *fasthttpRequestExecutor) executeRequest(req internalRequest, client *Client) error {
	internalError := &Error{
		Endpoint:         req.endpoint,
		Method:           req.method,
		Function:         req.functionName,
		RequestToString:  "empty request",
		ResponseToString: "empty response",
		MeilisearchApiError: meilisearchApiError{
			Message: "empty Meilisearch message",
		},
		StatusCodeExpected: req.acceptedStatusCodes,
	}

	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)
	err := c.sendRequest(&req, internalError, response, client)
	if err != nil {
		return err
	}
	internalError.StatusCode = response.StatusCode()

	err = c.handleStatusCode(&req, response, internalError)
	if err != nil {
		return err
	}

	err = c.handleResponse(&req, response, internalError)
	if err != nil {
		return err
	}
	return nil
}

func (c *fasthttpRequestExecutor) sendRequest(req *internalRequest, internalError *Error, response *fasthttp.Response, client *Client) error {
	var (
		request *fasthttp.Request

		err error
	)

	// Setup URL
	requestURL, err := url.Parse(client.config.Host + req.endpoint)
	if err != nil {
		return fmt.Errorf("unable to parse url: %w", err)
	}

	// Build query parameters
	if req.withQueryParams != nil {
		query := requestURL.Query()
		for key, value := range req.withQueryParams {
			query.Set(key, value)
		}

		requestURL.RawQuery = query.Encode()
	}

	request = fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)

	request.SetRequestURI(requestURL.String())
	request.Header.SetMethod(req.method)

	if req.withRequest != nil {
		if req.method == http.MethodGet || req.method == http.MethodHead {
			return fmt.Errorf("sendRequest: request body is not expected for GET and HEAD requests")
		}
		if req.contentType == "" {
			return fmt.Errorf("sendRequest: request body without Content-Type is not allowed")
		}

		rawRequest := req.withRequest
		if b, ok := rawRequest.([]byte); ok {
			// If the request body is already a []byte then use it directly
			request.SetBody(b)
		} else if reader, ok := rawRequest.(io.Reader); ok {
			// If the request body is an io.Reader then stream it directly until io.EOF
			// NOTE: Avoid using this, due to problems with streamed request bodies
			request.SetBodyStream(reader, -1)
		} else {
			// Otherwise convert it to JSON
			var (
				data []byte
				err  error
			)
			if marshaler, ok := rawRequest.(json.Marshaler); ok {
				data, err = marshaler.MarshalJSON()
			} else {
				data, err = json.Marshal(rawRequest)
			}
			internalError.RequestToString = string(data)
			if err != nil {
				return internalError.WithErrCode(ErrCodeMarshalRequest, err)
			}
			request.SetBody(data)
		}
	}

	// adding request headers
	if req.contentType != "" {
		request.Header.Set("Content-Type", req.contentType)
	}
	if client.config.APIKey != "" {
		request.Header.Set("Authorization", "Bearer "+client.config.APIKey)
	}

	request.Header.Set("User-Agent", GetQualifiedVersion())

	// request is sent
	if client.config.Timeout != 0 {
		err = c.httpClient.DoTimeout(request, response, client.config.Timeout)
	} else {
		err = c.httpClient.Do(request, response)
	}

	// request execution timeout
	if err == fasthttp.ErrTimeout {
		return internalError.WithErrCode(MeilisearchTimeoutError, err)
	}
	// request execution fail
	if err != nil {
		return internalError.WithErrCode(MeilisearchCommunicationError, err)
	}

	return nil
}

func (c *fasthttpRequestExecutor) handleStatusCode(req *internalRequest, response *fasthttp.Response, internalError *Error) error {
	if req.acceptedStatusCodes != nil {
		return handleStatusCode(req, response.StatusCode(), response.Body(), internalError)
	}
	return nil
}

func handleStatusCode(req *internalRequest, statusCode int, rawBody []byte, internalError *Error) error {
	// A successful status code is required so check if the response status code is in the
	// expected status code list.
	for _, acceptedCode := range req.acceptedStatusCodes {
		if statusCode == acceptedCode {
			return nil
		}
	}
	// At this point the response status code is a failure.
	internalError.ErrorBody(rawBody)

	if internalError.MeilisearchApiError.Code == "" {
		return internalError.WithErrCode(MeilisearchApiErrorWithoutMessage)
	}
	return internalError.WithErrCode(MeilisearchApiError)
}

func (c *fasthttpRequestExecutor) handleResponse(req *internalRequest, response *fasthttp.Response, internalError *Error) (err error) {
	if req.withResponse != nil {
		return handleResponse(req, response.Body(), internalError)
	}
	return nil
}

func handleResponse(req *internalRequest, rawBody []byte, internalError *Error) error {
	// A json response is mandatory, so the response interface{} need to be unmarshal from the response payload.
	internalError.ResponseToString = string(rawBody)

	var err error
	if resp, ok := req.withResponse.(json.Unmarshaler); ok {
		err = resp.UnmarshalJSON(rawBody)
		req.withResponse = resp
	} else {
		err = json.Unmarshal(rawBody, req.withResponse)
	}
	if err != nil {
		return internalError.WithErrCode(ErrCodeResponseUnmarshalBody, err)
	}
	return nil
}

type nethttpRequestExecutor struct {
	httpClient *http.Client
}

var _ requestExecutor = &nethttpRequestExecutor{}

func (c *nethttpRequestExecutor) executeRequest(req internalRequest, client *Client) error {
	internalError := &Error{
		Endpoint:         req.endpoint,
		Method:           req.method,
		Function:         req.functionName,
		RequestToString:  "empty request",
		ResponseToString: "empty response",
		MeilisearchApiError: meilisearchApiError{
			Message: "empty Meilisearch message",
		},
		StatusCodeExpected: req.acceptedStatusCodes,
	}

	response, err := c.sendRequest(&req, internalError, client)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	internalError.StatusCode = response.StatusCode

	rawBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = c.handleStatusCode(&req, response.StatusCode, rawBody, internalError)
	if err != nil {
		return err
	}

	err = c.handleResponse(&req, rawBody, internalError)
	if err != nil {
		return err
	}
	return nil
}

func (c *nethttpRequestExecutor) sendRequest(req *internalRequest, internalError *Error, client *Client) (*http.Response, error) {
	var (
		request *http.Request
		err     error
	)

	// Setup URL
	requestURL, err := url.Parse(client.config.Host + req.endpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to parse url: %w", err)
	}

	// Build query parameters
	if req.withQueryParams != nil {
		query := requestURL.Query()
		for key, value := range req.withQueryParams {
			query.Set(key, value)
		}
		requestURL.RawQuery = query.Encode()
	}

	// Create request
	request, err = http.NewRequest(req.method, requestURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	// Set request body
	if req.withRequest != nil {
		if req.method == http.MethodGet || req.method == http.MethodHead {
			return nil, fmt.Errorf("sendRequest: request body is not expected for GET and HEAD requests")
		}

		if req.contentType == "" {
			return nil, fmt.Errorf("sendRequest: request body without Content-Type is not allowed")
		}

		rawRequest := req.withRequest
		if b, ok := rawRequest.([]byte); ok {
			// If the request body is already a []byte then use it directly
			request.Body = io.NopCloser(bytes.NewReader(b))
		} else if readCloser, ok := rawRequest.(io.ReadCloser); ok {
			// If the request body is an io.ReadCloser then use it directly
			request.Body = readCloser
		} else if reader, ok := rawRequest.(io.Reader); ok {
			// If the request body is an io.Reader then use it directly
			request.Body = io.NopCloser(reader)
		} else {
			// Otherwise convert it to JSON
			var (
				data []byte
				err  error
			)
			if marshaler, ok := rawRequest.(json.Marshaler); ok {
				data, err = marshaler.MarshalJSON()
			} else {
				data, err = json.Marshal(rawRequest)
			}
			internalError.RequestToString = string(data)
			if err != nil {
				return nil, internalError.WithErrCode(ErrCodeMarshalRequest, err)
			}
			request.Body = io.NopCloser(bytes.NewReader(data))
		}
	}

	// Set request headers
	if req.contentType != "" {
		request.Header.Set("Content-Type", req.contentType)
	}
	if client.config.APIKey != "" {
		request.Header.Set("Authorization", "Bearer "+client.config.APIKey)
	}
	request.Header.Set("User-Agent", GetQualifiedVersion())

	// Send request
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, internalError.WithErrCode(MeilisearchCommunicationError, err)
	}

	return response, nil
}

func (c *nethttpRequestExecutor) handleStatusCode(req *internalRequest, statusCode int, rawBody []byte, internalError *Error) error {
	if req.acceptedStatusCodes != nil {
		return handleStatusCode(req, statusCode, rawBody, internalError)
	}
	return nil
}

func (c *nethttpRequestExecutor) handleResponse(req *internalRequest, rawBody []byte, internalError *Error) error {
	if req.withResponse != nil {
		return handleResponse(req, rawBody, internalError)
	}
	return nil
}
