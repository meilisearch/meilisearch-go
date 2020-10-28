package meilisearch

import (
	"github.com/valyala/fastjson"
	"net/http"
)

type clientHealth struct {
	client *Client
	arp    *fastjson.ArenaPool
}

func newClientHealth(client *Client) clientHealth {
	return clientHealth{client: client}
}

func (c clientHealth) Get() error {
	req := internalRequest{
		endpoint:            "/health",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        nil,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Get",
		apiName:             "Health",
	}

	return c.client.executeRequest(req)
}

func (c clientHealth) Update(health bool) error {

	req := internalRequest{
		endpoint:            "/health",
		method:              http.MethodPut,
		withRequest:         Health{Health: health},
		withResponse:        nil,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Set",
		apiName:             "Health",
	}

	return c.client.executeRequest(req)
}
