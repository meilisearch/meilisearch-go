package meilisearch

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

var (
	defaultMeiliOpt = &meiliOpt{
		client: &http.Client{
			Transport: baseTransport(),
		},
	}
)

type meiliOpt struct {
	client *http.Client
	apiKey string
}

type Option func(*meiliOpt)

// WithCustomClient set custom http.Client
func WithCustomClient(client *http.Client) Option {
	return func(opt *meiliOpt) {
		opt.client = client
	}
}

// WithCustomClientWithTLS client support tls configuration
func WithCustomClientWithTLS(tlsConfig *tls.Config) Option {
	return func(opt *meiliOpt) {
		trans := baseTransport()
		trans.TLSClientConfig = tlsConfig
		opt.client = &http.Client{Transport: trans}
	}
}

// WithAPIKey is API key or master key.
// more: https://www.meilisearch.com/docs/reference/api/keys
func WithAPIKey(key string) Option {
	return func(opt *meiliOpt) {
		opt.apiKey = key
	}
}

func baseTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
