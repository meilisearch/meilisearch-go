package meilisearch

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"time"
)

type meiliOpt struct {
	client          *http.Client
	transport       *http.Transport
	apiKey          string
	contentEncoding *encodingOpt
	retryOnStatus   map[int]bool
	disableRetry    bool
	maxRetries      uint8
	jsonMarshaler   JSONMarshal
	jsonUnmarshaler JSONUnmarshal
}

type encodingOpt struct {
	encodingType ContentEncoding
	level        EncodingCompressionLevel
}

type Option func(*meiliOpt)

func _defaultOpts() *meiliOpt {
	transport := baseTransport()
	return &meiliOpt{
		client: &http.Client{
			Transport: transport,
		},
		transport: transport,
		contentEncoding: &encodingOpt{
			level: DefaultCompression,
		},
		retryOnStatus: map[int]bool{
			502: true,
			503: true,
			504: true,
		},
		disableRetry:    false,
		maxRetries:      3,
		jsonMarshaler:   json.Marshal,
		jsonUnmarshaler: json.Unmarshal,
	}
}

// WithCustomClient set custom http.Client
func WithCustomClient(client *http.Client) Option {
	return func(opt *meiliOpt) {
		opt.client = client
	}
}

// WithCustomProxy sets a custom proxy function
func WithCustomProxy(proxy func(*http.Request) (*url.URL, error)) Option {
	return func(opt *meiliOpt) {
		if proxy != nil {
			opt.transport.Proxy = proxy
		}
	}
}

// WithCustomDialContext sets a custom dial context
func WithCustomDialContext(dial func(ctx context.Context, network, addr string) (net.Conn, error)) Option {
	return func(opt *meiliOpt) {
		if dial != nil {
			opt.transport.DialContext = dial
		}
	}
}

// WithCustomMaxIdleConns sets the max number of idle connections
func WithCustomMaxIdleConns(maxIdleConns int) Option {
	return func(opt *meiliOpt) {
		opt.transport.MaxIdleConns = maxIdleConns
	}
}

// WithCustomMaxIdleConnsPerHost sets max idle connections per host
func WithCustomMaxIdleConnsPerHost(maxIdleConnsPerHost int) Option {
	return func(opt *meiliOpt) {
		opt.transport.MaxIdleConnsPerHost = maxIdleConnsPerHost
	}
}

// WithCustomIdleConnTimeout sets the idle connection timeout
func WithCustomIdleConnTimeout(timeout time.Duration) Option {
	return func(opt *meiliOpt) {
		opt.transport.IdleConnTimeout = timeout
	}
}

// WithCustomTLSHandshakeTimeout sets the TLS handshake timeout
func WithCustomTLSHandshakeTimeout(timeout time.Duration) Option {
	return func(opt *meiliOpt) {
		opt.transport.TLSHandshakeTimeout = timeout
	}
}

// WithCustomExpectContinueTimeout sets the Expect-Continue timeout
func WithCustomExpectContinueTimeout(timeout time.Duration) Option {
	return func(opt *meiliOpt) {
		opt.transport.ExpectContinueTimeout = timeout
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

// WithCustomJsonMarshaler set custom marshal from external packages instead encoding/json.
// we use encoding/json as default json library due to stability and producibility. However,
// the standard library is a bit slow compared to 3rd party libraries. If you're not happy with
// the performance of encoding/json.
//
// supported package: goccy/go-json, bytedance/sonic, segmentio/encoding, minio/simdjson-go, wI2L/jettison, mailru/easyjson.
//
// default is encoding/json
func WithCustomJsonMarshaler(marshal JSONMarshal) Option {
	return func(opt *meiliOpt) {
		opt.jsonMarshaler = marshal
	}
}

// WithCustomJsonUnmarshaler set custom unmarshal from external packages instead encoding/json.
// we use encoding/json as default json library due to stability and producibility. However,
// the standard library is a bit slow compared to 3rd party libraries. If you're not happy with
// the performance of encoding/json.
//
// supported package: goccy/go-json, bytedance/sonic, segmentio/encoding, minio/simdjson-go, wI2L/jettison, mailru/easyjson.
//
// default is encoding/json
func WithCustomJsonUnmarshaler(unmarshal JSONUnmarshal) Option {
	return func(opt *meiliOpt) {
		opt.jsonUnmarshaler = unmarshal
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
