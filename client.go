package meilisearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Config struct {
	Host string

	// APIKey is optional
	APIKey string
}

type Client struct {
	config     Config
	httpClient http.Client
}

func NewClient(config Config) *Client {
	return &Client{
		config: config,
		httpClient: http.Client{
			Timeout: time.Second,
		},
	}
}

func NewClientWithCustomHttpClient(config Config, client http.Client) *Client {
	return &Client{
		config:     config,
		httpClient: client,
	}
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
	errContext := fmt.Sprintf(`Endpoint="%s %s" Function=%s ApiName=%s`, i.method, i.endpoint, i.functionName, i.apiName)

	var (
		request *http.Request
		err     error
	)

	if i.withRequest != nil {
		b, err := json.Marshal(i.withRequest)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to marshal body from request (%s)", errContext))
		}
		request, err = http.NewRequest(i.method, path.Join(c.config.Host, i.endpoint), bytes.NewBuffer(b))
	} else {
		request, err = http.NewRequest(i.method, path.Join(c.config.Host, i.endpoint), nil)
	}

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to create new request (%s)", errContext))
	}

	if c.config.APIKey != "" {
		request.Header.Set("X-Meili-API-Key", c.config.APIKey)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to execute request (%s)", errContext))
	}

	code := response.StatusCode
	if i.acceptedStatusCodes != nil {
		ok := false
		for _, acceptedCode := range i.acceptedStatusCodes {
			if code == acceptedCode {
				ok = true
				break
			}
		}

		if !ok {
			return fmt.Errorf("status code received is not a status code of success: %v (%s)", i.acceptedStatusCodes, errContext)
		}
	}

	if i.withResponse != nil {
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to read body from response (%s)", errContext))
		}
		if err := json.Unmarshal(b, i.withResponse); err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to unmarshal body from response (%s)", errContext))
		}
	}

	return nil
}

func IsStatusCodeErr(err error) bool {
	return strings.HasPrefix(err.Error(), "status code received is not a status code of success")
}

func IsRequestMarshalErr(err error) bool {
	return strings.HasPrefix(err.Error(), "unable to marshal body from request")
}

func IsRequestCreationErr(err error) bool {
	return strings.HasPrefix(err.Error(), "unable to create new request")
}

func IsRequestExecutionErr(err error) bool {
	return strings.HasPrefix(err.Error(), "unable to execute request")
}

func IsResponseBodyErr(err error) bool {
	return strings.HasPrefix(err.Error(), "unable to read body from response")
}

func IsResponseUnmarshalErr(err error) bool {
	return strings.HasPrefix(err.Error(), "unable to unmarshal body from response")
}
