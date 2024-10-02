package meilisearch

import (
	"crypto/tls"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestOptions_WithCustomClient(t *testing.T) {
	meili := setup(t, "", WithCustomClient(http.DefaultClient))
	require.NotNil(t, meili)

	m, ok := meili.(*meilisearch)
	require.True(t, ok)

	require.Equal(t, m.client.client, http.DefaultClient)
}

func TestOptions_WithCustomClientWithTLS(t *testing.T) {
	tl := new(tls.Config)
	meili := setup(t, "", WithCustomClientWithTLS(tl))
	require.NotNil(t, meili)

	m, ok := meili.(*meilisearch)
	require.True(t, ok)

	tr, ok := m.client.client.Transport.(*http.Transport)
	require.True(t, ok)

	require.Equal(t, tr.TLSClientConfig, tl)
}

func TestOptions_WithAPIKey(t *testing.T) {
	meili := setup(t, "", WithAPIKey("foobar"))
	require.NotNil(t, meili)

	m, ok := meili.(*meilisearch)
	require.True(t, ok)

	require.Equal(t, m.client.apiKey, "foobar")
}

func TestOptions_WithContentEncoding(t *testing.T) {
	meili := setup(t, "", WithContentEncoding(GzipEncoding, DefaultCompression))
	require.NotNil(t, meili)

	m, ok := meili.(*meilisearch)
	require.True(t, ok)

	require.Equal(t, m.client.contentEncoding, GzipEncoding)
	require.NotNil(t, m.client.encoder)
}

func TestOptions_WithCustomRetries(t *testing.T) {
	meili := setup(t, "", WithCustomRetries([]int{http.StatusInternalServerError}, 10))
	require.NotNil(t, meili)

	m, ok := meili.(*meilisearch)
	require.True(t, ok)
	require.True(t, m.client.retryOnStatus[http.StatusInternalServerError])
	require.Equal(t, m.client.maxRetries, uint8(10))

	meili = setup(t, "", WithCustomRetries([]int{http.StatusInternalServerError}, 0))
	require.NotNil(t, meili)

	m, ok = meili.(*meilisearch)
	require.True(t, ok)
	require.True(t, m.client.retryOnStatus[http.StatusInternalServerError])
	require.Equal(t, m.client.maxRetries, uint8(1))
}

func TestOptions_DisableRetries(t *testing.T) {
	meili := setup(t, "", DisableRetries())
	require.NotNil(t, meili)

	m, ok := meili.(*meilisearch)
	require.True(t, ok)
	require.Equal(t, m.client.disableRetry, true)
}
