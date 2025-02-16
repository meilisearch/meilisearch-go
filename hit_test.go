package meilisearch

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStruct represents a sample struct for decoding tests.
type TestStruct struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func TestHit_Decode(t *testing.T) {
	t.Run("Decode valid Hit into struct", func(t *testing.T) {
		hit := Hit{
			"id":    json.RawMessage(`123`),
			"title": json.RawMessage(`"Test Book"`),
		}

		var result TestStruct
		err := hit.Decode(&result)

		require.NoError(t, err)
		assert.Equal(t, 123, result.ID)
		assert.Equal(t, "Test Book", result.Title)
	})

	t.Run("Decode Hit with invalid JSON", func(t *testing.T) {
		hit := Hit{
			"id":    json.RawMessage(`invalid`), // Invalid JSON
			"title": json.RawMessage(`"Test Book"`),
		}

		var result TestStruct
		err := hit.Decode(&result)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to marshal hit")
	})
}

func TestHits_Decode(t *testing.T) {
	t.Run("Decode multiple Hits into struct slice", func(t *testing.T) {
		hits := Hits{
			{
				"id":    json.RawMessage(`1`),
				"title": json.RawMessage(`"Book One"`),
			},
			{
				"id":    json.RawMessage(`2`),
				"title": json.RawMessage(`"Book Two"`),
			},
		}

		var results []TestStruct
		err := hits.Decode(&results)

		require.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, 1, results[0].ID)
		assert.Equal(t, "Book One", results[0].Title)
		assert.Equal(t, 2, results[1].ID)
		assert.Equal(t, "Book Two", results[1].Title)
	})

	t.Run("Decode empty Hits", func(t *testing.T) {
		hits := Hits{}

		var results []TestStruct
		err := hits.Decode(&results)

		require.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("Decode into non-pointer slice", func(t *testing.T) {
		hits := Hits{
			{
				"id":    json.RawMessage(`1`),
				"title": json.RawMessage(`"Book One"`),
			},
		}

		var results []TestStruct
		err := hits.Decode(results) // Not a pointer

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "v must be a pointer to a slice")
	})

	t.Run("Decode into incorrect type", func(t *testing.T) {
		hits := Hits{
			{
				"id":    json.RawMessage(`1`),
				"title": json.RawMessage(`"Book One"`),
			},
		}

		var results string // Wrong type
		err := hits.Decode(&results)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "v must be a pointer to a slice")
	})
}
