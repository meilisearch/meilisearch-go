package meilisearch

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"net/url"
	"time"
)

// FastHTTPClient is a structure that give you the power for interacting with an high-level api with meilisearch.
type FastHTTPClient struct {
	config     Config
	httpClient *fasthttp.Client

	// singleton clients which don't need index id
	apiIndexes APIIndexes
	apiKeys    APIKeys
	apiStats   APIStats
	apiHealth  APIHealth
	apiVersion APIVersion
}

// Indexes return an APIIndexes client.
func (c *FastHTTPClient) Indexes() APIIndexes {
	return c.apiIndexes
}

// Version return an APIVersion client.
func (c *FastHTTPClient) Version() APIVersion {
	return c.apiVersion
}

// Documents return an APIDocuments client.
func (c *FastHTTPClient) Documents(indexID string) APIDocuments {
	return newFastClientDocuments(c, indexID)
}

// Search return an APISearch client.
func (c *FastHTTPClient) Search(indexID string) APISearch {
	return newFastClientSearch(c, indexID)
}

// Updates return an APIUpdates client.
func (c *FastHTTPClient) Updates(indexID string) APIUpdates {
	return newFastClientUpdates(c, indexID)
}

// Settings return an APISettings client.
func (c *FastHTTPClient) Settings(indexID string) APISettings {
	return newFastClientSettings(c, indexID)
}

// Keys return an APIKeys client.
func (c *FastHTTPClient) Keys() APIKeys {
	return c.apiKeys
}

// Stats return an APIStats client.
func (c *FastHTTPClient) Stats() APIStats {
	return c.apiStats
}

// Health return an APIHealth client.
func (c *FastHTTPClient) Health() APIHealth {
	return c.apiHealth
}

// NewFastHTTPCustomClient creates Meilisearch with custom fasthttp.Client
func NewFastHTTPCustomClient(config Config, client *fasthttp.Client) ClientInterface {
	c := &FastHTTPClient{
		config:     config,
		httpClient: client,
	}

	c.apiIndexes = newFastClientIndexes(c)
	c.apiKeys = newFastClientKeys(c)
	c.apiHealth = newFastClientHealth(c)
	c.apiStats = newFastClientStats(c)
	c.apiVersion = newFastClientVersion(c)

	return c
}

// NewClient creates Meilisearch with default fasthttp.Client
func NewClient(config Config) ClientInterface {
	client := &fasthttp.Client{
		Name: "meilsearch-client",
	}

	c := &FastHTTPClient{
		config:     config,
		httpClient: client,
	}

	c.apiIndexes = newFastClientIndexes(c)
	c.apiKeys = newFastClientKeys(c)
	c.apiHealth = newFastClientHealth(c)
	c.apiStats = newFastClientStats(c)
	c.apiVersion = newFastClientVersion(c)

	return c
}

type internalRawRequest struct {
	endpoint string
	method   string

	withRequest     interface{}
	withResponse    interface{}
	withQueryParams map[string]string

	acceptedStatusCodes []int

	functionName string
	apiName      string
}

func (c *FastHTTPClient) executeRequest(req internalRawRequest) error {
	internalError := &Error{
		Endpoint:           req.endpoint,
		Method:             req.method,
		Function:           req.functionName,
		APIName:            req.apiName,
		RequestToString:    "empty request",
		ResponseToString:   "empty response",
		MeilisearchMessage: "empty meilisearch message",
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

func (c *FastHTTPClient) sendRequest(req *internalRawRequest, internalError *Error, response *fasthttp.Response) error {
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
	request.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" {
		request.Header.Set("X-Meili-API-Key", c.config.APIKey)
	}

	// request is sent
	err = c.httpClient.Do(request, response)

	// request execution fail
	if err != nil {
		return internalError.WithErrCode(ErrCodeRequestExecution, err)
	}

	return nil
}

func (c *FastHTTPClient) handleStatusCode(req *internalRawRequest, response *fasthttp.Response, internalError *Error) error {
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

		return internalError.WithErrCode(ErrCodeResponseStatusCode)
	}

	return nil
}

func (c *FastHTTPClient) handleResponse(req *internalRawRequest, response *fasthttp.Response, internalError *Error) (err error) {
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

// DefaultWaitForPendingUpdate checks each 50ms the status of a WaitForPendingUpdate.
// This is a default implementation of WaitForPendingUpdate.
func (c FastHTTPClient) DefaultWaitForPendingUpdate(indexUID string, updateID *AsyncUpdateID) (UpdateStatus, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	return c.WaitForPendingUpdate(ctx, time.Millisecond*50, indexUID, updateID)
}

// WaitForPendingUpdate waits for the end of an update.
// The function will check by regular interval provided in parameter interval
// the UpdateStatus. If it is not UpdateStatusEnqueued or the ctx cancelled
// we return the UpdateStatus.
func (c FastHTTPClient) WaitForPendingUpdate(
	ctx context.Context,
	interval time.Duration,
	indexID string,
	updateID *AsyncUpdateID) (UpdateStatus, error) {

	apiUpdates := c.Updates(indexID)
	for {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		update, err := apiUpdates.Get(updateID.UpdateID)
		if err != nil {
			return UpdateStatusUnknown, nil
		}
		if update.Status != UpdateStatusEnqueued {
			return update.Status, nil
		}
		time.Sleep(interval)
	}
}
