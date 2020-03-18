package meilisearch

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// Config configure the Client
type Config struct {

	// Host is the host of your meilisearch database
	// Example: 'http://localhost:7700'
	Host string

	// APIKey is optional
	APIKey string
}

// Client is a structure that give you the power for interacting with an high-level api with meilisearch.
type Client struct {
	config     Config
	httpClient http.Client

	// singleton clients which don't need index id
	apiIndexes APIIndexes
	apiVersion APIVersion
	apiKeys    APIKeys
	apiStats   APIStats
	apiHealth  APIHealth
}

// NewClient create a Client from a Config structure.
func NewClient(config Config) *Client {
	return NewClientWithCustomHTTPClient(config, http.Client{
		Timeout: time.Second,
	})
}

// NewClientWithCustomHTTPClient create a Client from a Config structure and a http.Client which you can customize.
func NewClientWithCustomHTTPClient(config Config, client http.Client) *Client {
	c := &Client{
		config:     config,
		httpClient: client,
	}

	c.apiIndexes = newClientIndexes(c)
	c.apiVersion = newClientVersion(c)
	c.apiKeys = newClientKeys(c)
	c.apiHealth = newClientHealth(c)
	c.apiStats = newClientStats(c)

	return c
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

type internalRequest struct {
	endpoint string
	method   string

	withRequest  interface{}
	withResponse interface{}

	acceptedStatusCodes []int

	functionName string
	apiName      string
}

func (c Client) executeRequest(req internalRequest) error {
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

	response, err := c.sendRequest(&req, internalError)
	if err != nil {
		return err
	}

	internalError.StatusCode = response.StatusCode

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

func (c Client) sendRequest(req *internalRequest, internalError *Error) (*http.Response, error) {
	var (
		request *http.Request
		err     error
	)

	if req.withRequest != nil {

		// A json request is mandatory, so the request interface{} need to be passed as a raw json body.
		rawJSONRequest, errJSONMarshalling := json.Marshal(req.withRequest)
		if errJSONMarshalling != nil {
			return nil, internalError.WithErrCode(ErrCodeMarshalRequest, errJSONMarshalling)
		}

		internalError.RequestToString = string(rawJSONRequest)

		request, err = http.NewRequest(req.method, c.config.Host+req.endpoint, bytes.NewBuffer(rawJSONRequest))
	} else {
		request, err = http.NewRequest(req.method, c.config.Host+req.endpoint, nil)
	}

	if err != nil {
		return nil, internalError.WithErrCode(ErrCodeRequestCreation, err)
	}

	// adding apikey to the request
	if c.config.APIKey != "" {
		request.Header.Set("X-Meili-API-Key", c.config.APIKey)
	}

	// request is sent
	response, err := c.httpClient.Do(request)

	// request execution fail
	if err != nil {
		return nil, internalError.WithErrCode(ErrCodeRequestExecution, err)
	}

	return response, nil
}

func (c Client) handleStatusCode(req *internalRequest, response *http.Response, internalError *Error) error {
	if req.acceptedStatusCodes != nil {

		// A successful status code is required so check if the response status code is in the
		// expected status code list.
		for _, acceptedCode := range req.acceptedStatusCodes {
			if response.StatusCode == acceptedCode {
				return nil
			}
		}

		// At this point the response status code is a failure.
		rawBody, err := ioutil.ReadAll(response.Body)
		if err == nil {
			internalError.ErrorBody(rawBody)
		} else {
			return internalError.WithErrCode(ErrCodeResponseStatusCode, err)
		}

		return internalError.WithErrCode(ErrCodeResponseStatusCode)
	}

	return nil
}

func (c Client) handleResponse(req *internalRequest, response *http.Response, internalError *Error) error {
	if req.withResponse != nil {

		// A json response is mandatory, so the response interface{} need to be unmarshal from the response payload.
		rawBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return internalError.WithErrCode(ErrCodeResponseReadBody, err)
		}

		internalError.ResponseToString = string(rawBody)

		if err := json.Unmarshal(rawBody, req.withResponse); err != nil {
			return internalError.WithErrCode(ErrCodeResponseUnmarshalBody, err)
		}
	}
	return nil
}

func contains(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// AwaitAsyncUpdateID check each 16ms the status of a AsyncUpdateID.
// This method should be avoided.
// TODO: improve this method by returning a channel
func (c Client) AwaitAsyncUpdateID(indexID string, updateID *AsyncUpdateID) UpdateStatus {
	apiUpdates := c.Updates(indexID)
	for {
		update, err := apiUpdates.Get(updateID.UpdateID)
		if err != nil {
			return UpdateStatusUnknown
		}
		if update.Status != UpdateStatusEnqueued {
			return update.Status
		}
		time.Sleep(time.Millisecond * 16)
	}
}

// AwaitAsyncUpdateID check each 25ms the status of a AsyncUpdateID.
// This method should be avoided.
func AwaitAsyncUpdateID(api APIWithIndexID, updateID *AsyncUpdateID) UpdateStatus {
	return api.Client().AwaitAsyncUpdateID(api.IndexID(), updateID)
}
