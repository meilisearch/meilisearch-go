package meilisearch

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type Config struct {
	Host string

	// APIKey is optional
	APIKey string
}

type Client struct {
	config     Config
	httpClient http.Client

	// singleton clients which don't need index id
	apiIndexes ApiIndexes
	apiVersion ApiVersion
}

func NewClient(config Config) *Client {
	return NewClientWithCustomHttpClient(config, http.Client{
		Timeout: time.Second,
	})
}

func NewClientWithCustomHttpClient(config Config, client http.Client) *Client {
	c := &Client{
		config:     config,
		httpClient: client,
	}

	c.apiIndexes = newClientIndexes(c)
	c.apiVersion = newClientVersion(c)

	return c
}

func (c *Client) Indexes() ApiIndexes {
	return c.apiIndexes
}

func (c *Client) Version() ApiVersion {
	return c.apiVersion
}

func (c *Client) Documents(indexId string) ApiDocuments {
	return newClientDocuments(c, indexId)
}

func (c *Client) Search(indexId string) ApiSearch {
	return newClientSearch(c, indexId)
}

func (c *Client) Updates(indexId string) ApiUpdates {
	return newClientUpdates(c, indexId)
}

func (c *Client) StopWords(indexId string) ApiStopWords {
	return newClientStopWords(c, indexId)
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

func (c Client) executeRequest(i internalRequest) error {
	meiliErr := &Error{
		Endpoint:           i.endpoint,
		Method:             i.method,
		Function:           i.functionName,
		APIName:            i.apiName,
		RequestToString:    "empty request",
		ResponseToString:   "empty response",
		MeilisearchMessage: "empty meilisearch message",
		StatusCodeExpected: i.acceptedStatusCodes,
	}

	var (
		request *http.Request
		err     error
	)

	if i.withRequest != nil {
		b, err := json.Marshal(i.withRequest)
		if err != nil {
			return meiliErr.WithErrCode(ErrCodeMarshalRequest, err)
		}
		meiliErr.RequestToString = string(b)
		request, err = http.NewRequest(i.method, c.config.Host+i.endpoint, bytes.NewBuffer(b))
	} else {
		request, err = http.NewRequest(i.method, c.config.Host+i.endpoint, nil)
	}

	if err != nil {
		return meiliErr.WithErrCode(ErrCodeRequestCreation, err)
	}

	if c.config.APIKey != "" {
		request.Header.Set("X-Meili-API-Key", c.config.APIKey)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return meiliErr.WithErrCode(ErrCodeRequestExecution, err)
	}

	code := response.StatusCode
	meiliErr.StatusCode = code

	if i.acceptedStatusCodes != nil {
		ok := false
		for _, acceptedCode := range i.acceptedStatusCodes {
			if code == acceptedCode {
				ok = true
				break
			}
		}

		if !ok {
			b, errbody := ioutil.ReadAll(response.Body)
			if errbody == nil {
				meiliErr.ErrorBody(b)
			}

			return meiliErr.WithErrCode(ErrCodeResponseStatusCode)
		}
	}

	if i.withResponse != nil {
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return meiliErr.WithErrCode(ErrCodeResponseReadBody, err)
		}
		meiliErr.ResponseToString = string(b)
		if err := json.Unmarshal(b, i.withResponse); err != nil {
			return meiliErr.WithErrCode(ErrCodeResponseUnmarshalBody, err)
		}
	}

	return nil
}

// AwaitAsyncUpdateId check each 16ms the status of a AsyncUpdateId.
// This method should be avoided.
// TODO: improve this method by returning a channel
func (c Client) AwaitAsyncUpdateId(indexId string, updateId *AsyncUpdateId) UpdateStatus {
	apiUpdates := c.Updates(indexId)
	for {
		update, err := apiUpdates.Get(updateId.UpdateID)
		if err != nil {
			return UpdateStatusUnknown
		}
		if update.Status != UpdateStatusEnqueued {
			return update.Status
		}
		time.Sleep(time.Millisecond * 16)
	}
}

// AwaitAsyncUpdateId check each 25ms the status of a AsyncUpdateId.
// This method should be avoided.
func AwaitAsyncUpdateId(api ApiWithIndexID, updateId *AsyncUpdateId) UpdateStatus {
	return api.Client().AwaitAsyncUpdateId(api.IndexId(), updateId)
}
