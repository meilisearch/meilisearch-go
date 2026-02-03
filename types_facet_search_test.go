package meilisearch

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFacetSearchRequest_FilterType(t *testing.T) {
	// Trying to encode a list of strings as filter
	// This simulates what a user might want to do: Filter: []string{"genre = horror", "year > 2000"}
	
	// Current implementation:
	// Filter string `json:"filter,omitempty"`
	
	// If I want to pass an array, I can't assign it to Filter string.
	// This confirms the limitation statically.
	
	req := FacetSearchRequest{
		FacetName: "genres",
		Filter:    []string{"genre = horror", "year > 2000"},
	}
	
	bytes, err := json.Marshal(req)
	require.NoError(t, err)

	// We only check if the filter part is correct by unmarshaling to a map
	var resultMap map[string]interface{}
	err = json.Unmarshal(bytes, &resultMap)
	require.NoError(t, err)
	
	filter, ok := resultMap["filter"]
	require.True(t, ok)
	
	// Check if it's a slice
	filterSlice, ok := filter.([]interface{})
	require.True(t, ok, "Filter should be a slice")
	require.Equal(t, 2, len(filterSlice))
	require.Equal(t, "genre = horror", filterSlice[0])
	require.Equal(t, "year > 2000", filterSlice[1])
}

func TestFacetSearchRequest_FilterString(t *testing.T) {
	req := FacetSearchRequest{
		FacetName: "genres",
		Filter:    "genre = horror",
	}

	bytes, err := json.Marshal(req)
	require.NoError(t, err)
	
	var resultMap map[string]interface{}
	err = json.Unmarshal(bytes, &resultMap)
	require.NoError(t, err)
	
	filter, ok := resultMap["filter"]
	require.True(t, ok)
	require.Equal(t, "genre = horror", filter)
}
