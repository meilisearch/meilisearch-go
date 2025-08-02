package meilisearch

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOptions_WithCustomProxy(t *testing.T) {
	proxy := func(*http.Request) (*url.URL, error) { return nil, nil }
	meili := New("localhost:7700", WithCustomProxy(proxy))
	require.NotNil(t, meili)

	m, ok := meili.(*meilisearch)
	require.True(t, ok)
	transport, ok := m.client.client.Transport.(*http.Transport)
	require.True(t, ok)

	require.Equal(t, reflect.ValueOf(proxy).Pointer(), reflect.ValueOf(transport.Proxy).Pointer())
}

func TestOptions_WithCustomDialContext(t *testing.T) {
	dial := func(ctx context.Context, network, addr string) (net.Conn, error) { return nil, nil }
	meili := New("localhost:7700", WithCustomDialContext(dial))
	require.NotNil(t, meili)

	m, ok := meili.(*meilisearch)
	require.True(t, ok)
	transport, ok := m.client.client.Transport.(*http.Transport)
	require.True(t, ok)
	require.Equal(t, reflect.ValueOf(dial).Pointer(), reflect.ValueOf(transport.DialContext).Pointer())
}

func TestOptions_WithCustomMaxIdleConns(t *testing.T) {
	meili := New("localhost:7700", WithCustomMaxIdleConns(50))
	require.NotNil(t, meili)

	m, ok := meili.(*meilisearch)
	require.True(t, ok)
	transport, ok := m.client.client.Transport.(*http.Transport)
	require.True(t, ok)
	require.Equal(t, 50, transport.MaxIdleConns)
}

func TestOptions_WithCustomMaxIdleConnsPerHost(t *testing.T) {
	meili := New("localhost:7700", WithCustomMaxIdleConnsPerHost(50))
	require.NotNil(t, meili)

	m, ok := meili.(*meilisearch)
	require.True(t, ok)
	transport, ok := m.client.client.Transport.(*http.Transport)
	require.True(t, ok)
	require.Equal(t, 50, transport.MaxIdleConnsPerHost)
}

func TestOptions_WithCustomIdleConnTimeout(t *testing.T) {
	timeout := time.Second * 5
	meili := New("localhost:7700", WithCustomIdleConnTimeout(timeout))
	require.NotNil(t, meili)

	m, ok := meili.(*meilisearch)
	require.True(t, ok)
	transport, ok := m.client.client.Transport.(*http.Transport)
	require.True(t, ok)
	require.Equal(t, timeout, transport.IdleConnTimeout)
}

func TestOptions_WithCustomTLSHandshakeTimeout(t *testing.T) {
	timeout := time.Second * 5
	meili := New("localhost:7700", WithCustomTLSHandshakeTimeout(timeout))
	require.NotNil(t, meili)

	m, ok := meili.(*meilisearch)
	require.True(t, ok)
	transport, ok := m.client.client.Transport.(*http.Transport)
	require.True(t, ok)
	require.Equal(t, timeout, transport.TLSHandshakeTimeout)
}

func TestOptions_WithCustomExpectContinueTimeout(t *testing.T) {
	timeout := time.Second * 5
	meili := New("localhost:7700", WithCustomExpectContinueTimeout(timeout))
	require.NotNil(t, meili)

	m, ok := meili.(*meilisearch)
	require.True(t, ok)
	transport, ok := m.client.client.Transport.(*http.Transport)
	require.True(t, ok)
	require.Equal(t, timeout, transport.ExpectContinueTimeout)
}

func TestOptions_WithCustomClient(t *testing.T) {
	meili := New("localhost:7700", WithCustomClient(http.DefaultClient))
	require.NotNil(t, meili)

	m, ok := meili.(*meilisearch)
	require.True(t, ok)

	require.Equal(t, m.client.client, http.DefaultClient)
}

func TestOptions_WithCustomClientWithTLS(t *testing.T) {
	tl := new(tls.Config)
	meili := New("localhost:7700", WithCustomClientWithTLS(tl))
	require.NotNil(t, meili)

	m, ok := meili.(*meilisearch)
	require.True(t, ok)

	tr, ok := m.client.client.Transport.(*http.Transport)
	require.True(t, ok)

	require.Equal(t, tr.TLSClientConfig, tl)
}

func TestOptions_WithAPIKey(t *testing.T) {
	meili := New("localhost:7700", WithAPIKey("foobar"))
	require.NotNil(t, meili)

	m, ok := meili.(*meilisearch)
	require.True(t, ok)

	require.Equal(t, m.client.apiKey, "foobar")
}

func TestOptions_WithContentEncoding(t *testing.T) {
	meili := New("localhost:7700", WithContentEncoding(GzipEncoding, DefaultCompression))
	require.NotNil(t, meili)

	m, ok := meili.(*meilisearch)
	require.True(t, ok)

	require.Equal(t, m.client.contentEncoding, GzipEncoding)
	require.NotNil(t, m.client.encoder)
}

func TestOptions_WithCustomRetries(t *testing.T) {
	meili := New("localhost:7700", WithCustomRetries([]int{http.StatusInternalServerError}, 10))
	require.NotNil(t, meili)

	m, ok := meili.(*meilisearch)
	require.True(t, ok)
	require.True(t, m.client.retryOnStatus[http.StatusInternalServerError])
	require.Equal(t, m.client.maxRetries, uint8(10))

	meili = New("localhost:7700", WithCustomRetries([]int{http.StatusInternalServerError}, 0))
	require.NotNil(t, meili)

	m, ok = meili.(*meilisearch)
	require.True(t, ok)
	require.True(t, m.client.retryOnStatus[http.StatusInternalServerError])
	require.Equal(t, m.client.maxRetries, uint8(1))
}

func TestOptions_DisableRetries(t *testing.T) {
	meili := New("localhost:7700", DisableRetries())
	require.NotNil(t, meili)

	m, ok := meili.(*meilisearch)
	require.True(t, ok)
	require.Equal(t, m.client.disableRetry, true)
}

func TestOptions_WithCustomJsonMarshalAndUnmarshaler(t *testing.T) {
	meili := New("localhost:7700", WithCustomJsonMarshaler(json.Marshal),
		WithCustomJsonUnmarshaler(json.Unmarshal))
	require.NotNil(t, meili)

	m, ok := meili.(*meilisearch)
	require.True(t, ok)

	require.NotNil(t, m.client.jsonMarshal)
	require.NotNil(t, m.client.jsonUnmarshal)
}
