package meilisearch

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestHit_Decode_Errors(t *testing.T) {
	t.Run("Decode with nil pointer", func(t *testing.T) {
		hit := Hit{
			"id":    json.RawMessage(`123`),
			"title": json.RawMessage(`"Test Book"`),
		}

		err := hit.Decode(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "vPtr must be a non-nil pointer")
	})

	t.Run("Decode with non-pointer value", func(t *testing.T) {
		hit := Hit{
			"id":    json.RawMessage(`123`),
			"title": json.RawMessage(`"Test Book"`),
		}

		var notPtr TestStruct
		err := hit.Decode(notPtr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "vPtr must be a non-nil pointer")
	})
}

func TestHits_Len(t *testing.T) {
	t.Run("Returns correct length", func(t *testing.T) {
		hits := Hits{
			{"id": json.RawMessage(`1`)},
			{"id": json.RawMessage(`2`)},
		}
		assert.Equal(t, 2, hits.Len())
	})
}

func TestHits_Decode_HitDecodeError(t *testing.T) {
	t.Run("One of hits causes decoding error", func(t *testing.T) {
		hits := Hits{
			{
				"id":    json.RawMessage(`1`),
				"title": json.RawMessage(`"Book One"`),
			},
			{
				"id":    json.RawMessage(`invalid`),
				"title": json.RawMessage(`"Book Two"`),
			},
		}

		var results []TestStruct
		err := hits.Decode(&results)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode hit at index 1")
	})
}
