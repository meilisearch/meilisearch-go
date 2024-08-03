package meilisearch

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTypes_UnmarshalJSON(t *testing.T) {
	var raw RawType
	data := []byte(`"some data"`)

	err := json.Unmarshal(data, &raw)
	require.NoError(t, err)

	expected := RawType(data)
	require.Equal(t, expected, raw)
}

func TestTypes_MarshalJSON(t *testing.T) {
	raw := RawType(`"some data"`)

	data, err := json.Marshal(raw)
	require.NoError(t, err)

	expected := []byte(`"some data"`)
	require.Equal(t, data, expected)
}

func TestTypes_ValidateSearchRequest(t *testing.T) {
	req := &SearchRequest{
		Limit: 0,
		Hybrid: &SearchRequestHybrid{
			Embedder: "",
		},
	}

	req.validate()

	assert.Equal(t, req.Limit, DefaultLimit)
	assert.Equal(t, req.Hybrid.Embedder, "default")
}
