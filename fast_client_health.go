package meilisearch

import (
	"github.com/valyala/fastjson"
	"net/http"
)

type fastClientHealth struct {
	client *FastHTTPClient
	arp    *fastjson.ArenaPool
}

func newFastClientHealth(client *FastHTTPClient) fastClientHealth {
	return fastClientHealth{client: client}
}

func (c fastClientHealth) Get() error {
	req := internalRawRequest{
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

func (c fastClientHealth) Update(health bool) error {

	req := internalRawRequest{
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
