package meilisearch

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

// ClientConfig configure the Client
type ClientConfig struct {

	// Host is the host of your Meilisearch database
	// Example: 'http://localhost:7700'
	Host string

	// APIKey is optional
	APIKey string

	// Timeout is optional
	Timeout time.Duration
}

type waitParams struct {
	Context  context.Context
	Interval time.Duration
}

// ClientInterface is interface for all Meilisearch client
type ClientInterface interface {
	Index(uid string) *Index
	GetIndex(indexID string) (resp *Index, err error)
	GetRawIndex(uid string) (resp map[string]interface{}, err error)
	GetAllIndexes() (resp []*Index, err error)
	GetAllRawIndexes() (resp []map[string]interface{}, err error)
	CreateIndex(config *IndexConfig) (resp *Task, err error)
	DeleteIndex(uid string) (resp *Task, err error)
	CreateKey(request *Key) (resp *Key, err error)
	GetKey(identifier string) (resp *Key, err error)
	GetKeys() (resp *ResultKey, err error)
	UpdateKey(identifier string, request *Key) (resp *Key, err error)
	DeleteKey(identifier string) (resp bool, err error)
	GetAllStats() (resp *Stats, err error)
	CreateDump() (resp *Dump, err error)
	GetDumpStatus(dumpUID string) (resp *Dump, err error)
	Version() (*Version, error)
	GetVersion() (resp *Version, err error)
	Health() (*Health, error)
	IsHealthy() bool
	GetTask(taskID int64) (resp *Task, err error)
	GetTasks() (resp *ResultTask, err error)
	WaitForTask(task *Task, options ...waitParams) (*Task, error)
}

var _ ClientInterface = &Client{}

// NewFastHTTPCustomClient creates Meilisearch with custom fasthttp.Client
func NewFastHTTPCustomClient(config ClientConfig, client *fasthttp.Client) *Client {
	c := &Client{
		config:     config,
		httpClient: client,
	}
	return c
}

// NewClient creates Meilisearch with default fasthttp.Client
func NewClient(config ClientConfig) *Client {
	client := &fasthttp.Client{
		Name: "meilsearch-client",
	}
	c := &Client{
		config:     config,
		httpClient: client,
	}
	return c
}

func (c *Client) Version() (resp *Version, err error) {
	resp = &Version{}
	req := internalRequest{
		endpoint:            "/version",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Version",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetVersion() (resp *Version, err error) {
	return c.Version()
}

func (c *Client) GetAllStats() (resp *Stats, err error) {
	resp = &Stats{}
	req := internalRequest{
		endpoint:            "/stats",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetAllStats",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) CreateKey(request *Key) (resp *Key, err error) {
	parsedRequest := convertKeyToParsedKey(*request)
	resp = &Key{}
	req := internalRequest{
		endpoint:            "/keys",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         &parsedRequest,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusCreated},
		functionName:        "CreateKey",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetKey(identifier string) (resp *Key, err error) {
	resp = &Key{}
	req := internalRequest{
		endpoint:            "/keys/" + identifier,
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetKey",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetKeys() (resp *ResultKey, err error) {
	resp = &ResultKey{}
	req := internalRequest{
		endpoint:            "/keys",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetKeys",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) UpdateKey(identifier string, request *Key) (resp *Key, err error) {
	parsedRequest := convertKeyToParsedKey(*request)
	resp = &Key{}
	req := internalRequest{
		endpoint:            "/keys/" + identifier,
		method:              http.MethodPatch,
		contentType:         contentTypeJSON,
		withRequest:         &parsedRequest,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "UpdateKey",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) DeleteKey(identifier string) (resp bool, err error) {
	req := internalRequest{
		endpoint:            "/keys/" + identifier,
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        nil,
		acceptedStatusCodes: []int{http.StatusNoContent},
		functionName:        "DeleteKey",
	}
	if err := c.executeRequest(req); err != nil {
		return false, err
	}
	return true, nil
}

func (c *Client) Health() (resp *Health, err error) {
	resp = &Health{}
	req := internalRequest{
		endpoint:            "/health",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Health",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) IsHealthy() bool {
	if _, err := c.Health(); err != nil {
		return false
	}
	return true
}

func (c *Client) CreateDump() (resp *Dump, err error) {
	resp = &Dump{}
	req := internalRequest{
		endpoint:            "/dumps",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "CreateDump",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetDumpStatus(dumpUID string) (resp *Dump, err error) {
	resp = &Dump{}
	req := internalRequest{
		endpoint:            "/dumps/" + dumpUID + "/status",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetDumpStatus",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetTask(taskID int64) (resp *Task, err error) {
	resp = &Task{}
	req := internalRequest{
		endpoint:            "/tasks/" + strconv.FormatInt(taskID, 10),
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetTask",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetTasks() (resp *ResultTask, err error) {
	resp = &ResultTask{}
	req := internalRequest{
		endpoint:            "/tasks",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetTasks",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

// WaitForTask waits for a task to be processed.
// The function will check by regular interval provided in parameter interval
// the TaskStatus.
// If no ctx and interval are provided WaitForTask will check each 50ms the
// status of a task.
func (c *Client) WaitForTask(task *Task, options ...waitParams) (*Task, error) {
	if options == nil {
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
		defer cancelFunc()
		options = []waitParams{
			{
				Context:  ctx,
				Interval: time.Millisecond * 50,
			},
		}
	}
	for {
		if err := options[0].Context.Err(); err != nil {
			return nil, err
		}
		getTask, err := c.GetTask(task.UID)
		if err != nil {
			return nil, err
		}
		if getTask.Status != TaskStatusEnqueued && getTask.Status != TaskStatusProcessing {
			return getTask, nil
		}
		time.Sleep(options[0].Interval)
	}
}

// This function allows the user to create a Key with an ExpiredAt in time.Time
// and transform the Key structure into a KeyParsed structure to send the time format
// managed by Meilisearch
func convertKeyToParsedKey(key Key) (resp KeyParsed) {
	resp = KeyParsed{Description: key.Description, Actions: key.Actions, Indexes: key.Indexes}

	// Convert time.Time to *string to feat the exact ISO-8601
	// format of Meilisearch
	if !key.ExpiresAt.IsZero() {
		const Format = "2006-01-02T15:04:05"
		timeParsedToString := key.ExpiresAt.Format(Format)
		resp.ExpiresAt = &timeParsedToString
	}
	return resp
}
