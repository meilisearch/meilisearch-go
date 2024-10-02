package meilisearch

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLocateString(t *testing.T) {
	require.Equal(t, ENG.String(), "eng")
}
