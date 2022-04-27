package meilisearch

import (
	"testing"
	"fmt"
	"regexp"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestVersion_GetQualifiedVersion(t *testing.T) {
	version := GetQualifiedVersion()

	require.NotNil(t, version)
	require.Equal(t, version, fmt.Sprintf("Meilisearch Go (v%s)", VERSION))
}

func TestVersion_qualifiedVersionFormat(t *testing.T) {
	version := getQualifiedVersion("2.2.5")

	require.NotNil(t, version)
	require.Equal(t, version, "Meilisearch Go (v2.2.5)")
}

func TestVersion_constVERSIONFormat(t *testing.T) {
	match, _ := regexp.MatchString("[0-9]+.[0-9]+.[0-9]+", VERSION)

	assert.True(t, match)
}
