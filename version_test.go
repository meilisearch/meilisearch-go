package meilisearch

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
