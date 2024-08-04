package meilisearch

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestOptions_WithCustomClient(t *testing.T) {
	meili := setup(t, "", WithCustomClient(http.DefaultClient))
	v, err := meili.Version()
	require.NoError(t, err)
	require.NotZero(t, v.PkgVersion)
}
