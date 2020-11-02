package meilisearch

import (
	"context"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"net/url"
	"time"

	"encoding/json"
)

// Config configure the Client
type Config struct {

	// Host is the host of your meilisearch database
	// Example: 'http://localhost:7700'
	Host string

	// APIKey is optional
	APIKey string
}

// ClientInterface is interface for all Meilisearch client
type ClientInterface interface {
	WaitForPendingUpdate(ctx context.Context, interval time.Duration, indexID string, updateID *AsyncUpdateID) (UpdateStatus, error)
	DefaultWaitForPendingUpdate(indexUID string, updateID *AsyncUpdateID) (UpdateStatus, error)

	Indexes() APIIndexes
	Version() APIVersion
	Documents(indexID string) APIDocuments
	Search(indexID string) APISearch
	Updates(indexID string) APIUpdates
	Settings(indexID string) APISettings
	Keys() APIKeys
	Stats() APIStats
	Health() APIHealth
}

// Client is a structure that give you the power for interacting with an high-level api with meilisearch.
type Client struct {
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
func (c *Client) Indexes() APIIndexes {
	return c.apiIndexes
}

// Version return an APIVersion client.
func (c *Client) Version() APIVersion {
	return c.apiVersion
}

// Documents return an APIDocuments client.
func (c *Client) Documents(indexID string) APIDocuments {
	return newClientDocuments(c, indexID)
}

// Search return an APISearch client.
func (c *Client) Search(indexID string) APISearch {
	return newClientSearch(c, indexID)
}

// Updates return an APIUpdates client.
func (c *Client) Updates(indexID string) APIUpdates {
	return newClientUpdates(c, indexID)
}

// Settings return an APISettings client.
func (c *Client) Settings(indexID string) APISettings {
	return newClientSettings(c, indexID)
}

// Keys return an APIKeys client.
func (c *Client) Keys() APIKeys {
	return c.apiKeys
}

// Stats return an APIStats client.
func (c *Client) Stats() APIStats {
	return c.apiStats
}

// Health return an APIHealth client.
func (c *Client) Health() APIHealth {
	return c.apiHealth
}

// NewFastHTTPCustomClient creates Meilisearch with custom fasthttp.Client
func NewFastHTTPCustomClient(config Config, client *fasthttp.Client) ClientInterface {
	c := &Client{
		config:     config,
		httpClient: client,
	}

	c.apiIndexes = newClientIndexes(c)
	c.apiKeys = newClientKeys(c)
	c.apiHealth = newClientHealth(c)
	c.apiStats = newClientStats(c)
	c.apiVersion = newClientVersion(c)

	return c
}

// NewClient creates Meilisearch with default fasthttp.Client
func NewClient(config Config) ClientInterface {
	client := &fasthttp.Client{
		Name: "meilsearch-client",
	}

	c := &Client{
		config:     config,
		httpClient: client,
	}

	c.apiIndexes = newClientIndexes(c)
	c.apiKeys = newClientKeys(c)
	c.apiHealth = newClientHealth(c)
	c.apiStats = newClientStats(c)
	c.apiVersion = newClientVersion(c)

	return c
}

type internalRequest struct {
	endpoint string
	method   string

	withRequest     interface{}
	withResponse    interface{}
	withQueryParams map[string]string

	acceptedStatusCodes []int

	functionName string
	apiName      string
}

func (c *Client) executeRequest(req internalRequest) error {
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

		return internalError.WithErrCode(ErrCodeResponseStatusCode)
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

// DefaultWaitForPendingUpdate checks each 50ms the status of a WaitForPendingUpdate.
// This is a default implementation of WaitForPendingUpdate.
func (c Client) DefaultWaitForPendingUpdate(indexUID string, updateID *AsyncUpdateID) (UpdateStatus, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	return c.WaitForPendingUpdate(ctx, time.Millisecond*50, indexUID, updateID)
}

// WaitForPendingUpdate waits for the end of an update.
// The function will check by regular interval provided in parameter interval
// the UpdateStatus. If it is not UpdateStatusEnqueued or the ctx cancelled
// we return the UpdateStatus.
func (c Client) WaitForPendingUpdate(
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
