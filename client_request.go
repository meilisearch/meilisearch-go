package meilisearch

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"encoding/json"
)

const (
	contentTypeJSON string = "application/json"
	// These two types are defined in advance of the resolution of this issue
	// https://github.com/meilisearch/meilisearch-go/issues/215. Don't forget to remove the no lint comment when it's done.
	contentTypeNDJSON string = "application/x-dnjson" //nolint
	contentTypeCSV    string = "text/csv"             //nolint
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
	internalError := &Error{
		Endpoint:         req.endpoint,
		Method:           req.method,
		Function:         req.functionName,
		RequestToString:  "empty request",
		ResponseToString: "empty response",
		MeilisearchApiMessage: meilisearchApiMessage{
			Message: "empty meilisearch message",
		},
		StatusCodeExpected: req.acceptedStatusCodes,
	}

	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)
	err := c.sendRequest(&req, internalError, response)
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

func (c *Client) sendRequest(req *internalRequest, internalError *Error, response *fasthttp.Response) error {
	var (
		request *fasthttp.Request

		err error
	)

	// Setup URL
	requestURL, err := url.Parse(c.config.Host + req.endpoint)
	if err != nil {
		return errors.Wrap(err, "unable to parse url")
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

		// A json request is mandatory, so the request interface{} need to be passed as a raw json body.
		rawJSONRequest := req.withRequest
		var data []byte
		var err error
		if raw, ok := rawJSONRequest.(json.Marshaler); ok {
			data, err = raw.MarshalJSON()
		} else {
			data, err = json.Marshal(rawJSONRequest)
		}
		internalError.RequestToString = string(data)
		if err != nil {
			return internalError.WithErrCode(ErrCodeMarshalRequest, err)
		}
		request.SetBody(data)
	}

	// adding request headers
	if req.contentType != "" {
		request.Header.Set("Content-Type", req.contentType)
	}
	if c.config.APIKey != "" {
		request.Header.Set("X-Meili-API-Key", c.config.APIKey)
	}

	// request is sent
	if c.config.Timeout != 0 {
		err = c.httpClient.DoTimeout(request, response, c.config.Timeout)
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

func (c *Client) handleStatusCode(req *internalRequest, response *fasthttp.Response, internalError *Error) error {
	if req.acceptedStatusCodes != nil {

		// A successful status code is required so check if the response status code is in the
		// expected status code list.
		for _, acceptedCode := range req.acceptedStatusCodes {
			if response.StatusCode() == acceptedCode {
				return nil
			}
		}
		// At this point the response status code is a failure.
		rawBody := response.Body()

		internalError.ErrorBody(rawBody)

		if internalError.MeilisearchApiMessage.Code == "" {
			return internalError.WithErrCode(MeilisearchApiErrorWithoutMessage)
		}
		return internalError.WithErrCode(MeilisearchApiError)
	}

	return nil
}

func (c *Client) handleResponse(req *internalRequest, response *fasthttp.Response, internalError *Error) (err error) {
	if req.withResponse != nil {

		// A json response is mandatory, so the response interface{} need to be unmarshal from the response payload.
		rawBody := response.Body()
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
	}
	return nil
}
