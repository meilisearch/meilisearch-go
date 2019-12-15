package meilisearch

import "net/http"

type clientHealth struct {
	client *Client
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
		acceptedStatusCodes: []int{http.StatusNoContent},
		functionName:        "Get",
		apiName:             "Health",
	}

	return c.client.executeRequest(req)
}

func (c clientHealth) Set(health bool) error {
	req := internalRequest{
		endpoint: "/health",
		method:   http.MethodPut,
		withRequest: map[string]bool{
			"health": health,
		},
		withResponse:        nil,
		acceptedStatusCodes: []int{http.StatusNoContent},
		functionName:        "Set",
		apiName:             "Health",
	}

	return c.client.executeRequest(req)
}
