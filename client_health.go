package meilisearch

import (
	"net/http"

	"github.com/valyala/fastjson"
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
