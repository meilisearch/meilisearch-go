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
		contentEncoding: &encodingOpt{
			level: DefaultCompression,
		},
		retryOnStatus: map[int]bool{
			502: true,
			503: true,
			504: true,
		},
		disableRetry: false,
		maxRetries:   3,
	}
)

type meiliOpt struct {
	client          *http.Client
	apiKey          string
	contentEncoding *encodingOpt
	retryOnStatus   map[int]bool
	disableRetry    bool
	maxRetries      uint8
}

type encodingOpt struct {
	encodingType ContentEncoding
	level        EncodingCompressionLevel
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
//
// more: https://www.meilisearch.com/docs/reference/api/keys
func WithAPIKey(key string) Option {
	return func(opt *meiliOpt) {
		opt.apiKey = key
	}
}

// WithContentEncoding support the Content-Encoding header indicates the media type is compressed by a given algorithm.
// compression improves transfer speed and reduces bandwidth consumption by sending and receiving smaller payloads.
// the Accept-Encoding header, instead, indicates the compression algorithm the client understands.
//
// more: https://www.meilisearch.com/docs/reference/api/overview#content-encoding
func WithContentEncoding(encodingType ContentEncoding, level EncodingCompressionLevel) Option {
	return func(opt *meiliOpt) {
		opt.contentEncoding = &encodingOpt{
			encodingType: encodingType,
			level:        level,
		}
	}
}

// WithCustomRetries set retry on specific http error code and max retries (min: 1, max: 255)
func WithCustomRetries(retryOnStatus []int, maxRetries uint8) Option {
	return func(opt *meiliOpt) {
		opt.retryOnStatus = make(map[int]bool)
		for _, status := range retryOnStatus {
			opt.retryOnStatus[status] = true
		}

		if maxRetries == 0 {
			maxRetries = 1
		}

		opt.maxRetries = maxRetries
	}
}

// DisableRetries disable retry logic in client
func DisableRetries() Option {
	return func(opt *meiliOpt) {
		opt.disableRetry = true
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
