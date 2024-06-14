package meilisearch

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestClient_Transport(t *testing.T) {
	called := false
	client := NewClient(ClientConfig{
		Host:   getenv("MEILISEARCH_URL", "http://localhost:7700"),
		APIKey: masterKey,
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			called = true
			return http.DefaultTransport.RoundTrip(req)
		}),
	})
	_, err := client.GetVersion()
	require.NoError(t, err)
	require.True(t, called)
}
