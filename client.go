package meilisearch

import (
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

// ClientConfig configure the Client
type ClientConfig struct {

	// Host is the host of your meilisearch database
	// Example: 'http://localhost:7700'
	Host string

	// APIKey is optional
	APIKey string

	// Timeout is optional
	Timeout time.Duration
}

// ClientInterface is interface for all Meilisearch client
type ClientInterface interface {
	Index(uid string) *Index
	GetIndex(indexID string) (resp *Index, err error)
	GetRawIndex(uid string) (resp map[string]interface{}, err error)
	GetAllIndexes() (resp []*Index, err error)
	GetAllRawIndexes() (resp []map[string]interface{}, err error)
	CreateIndex(config *IndexConfig) (resp *Index, err error)
	DeleteIndex(uid string) (bool, error)
	DeleteIndexIfExists(uid string) (bool, error)
	GetKeys() (resp *Keys, err error)
	GetAllStats() (resp *Stats, err error)
	CreateDump() (resp *Dump, err error)
	GetDumpStatus(dumpUID string) (resp *Dump, err error)
	Version() (*Version, error)
	GetVersion() (resp *Version, err error)
	Health() (*Health, error)
	IsHealthy() bool
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

func (c *Client) GetKeys() (resp *Keys, err error) {
	resp = &Keys{}
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
