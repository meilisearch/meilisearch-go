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

	return c
}

func (c *Client) Indexes() ApiIndexes {
	return c.apiIndexes
}

func (c *Client) Documents(indexId string) ApiDocuments {
	return newClientDocuments(c, indexId)
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
