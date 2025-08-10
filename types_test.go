package meilisearch

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type sampleStructure struct {
	ImportantString string `json:"important_string"`
}

func Test_GolangJSONEncoder(t *testing.T) {
	t.Parallel()

	var (
		ss = &sampleStructure{
			ImportantString: "Hello World",
		}
		importantString             = `{"important_string":"Hello World"}`
		jsonEncoder     JSONMarshal = json.Marshal
	)

	raw, err := jsonEncoder(ss)
	require.NoError(t, err)

	require.Equal(t, string(raw), importantString)
}

func Test_DefaultJSONEncoder(t *testing.T) {
	t.Parallel()

	var (
		ss = &sampleStructure{
			ImportantString: "Hello World",
		}
		importantString             = `{"important_string":"Hello World"}`
		jsonEncoder     JSONMarshal = json.Marshal
	)

	raw, err := jsonEncoder(ss)
	require.NoError(t, err)

	require.Equal(t, string(raw), importantString)
}

func Test_DefaultJSONDecoder(t *testing.T) {
	t.Parallel()

	var (
		ss              sampleStructure
		importantString               = []byte(`{"important_string":"Hello World"}`)
		jsonDecoder     JSONUnmarshal = json.Unmarshal
	)

	err := jsonDecoder(importantString, &ss)
	require.NoError(t, err)
	require.Equal(t, "Hello World", ss.ImportantString)
}

func TestSearchRequest_validate(t *testing.T) {
	t.Parallel()

	t.Run("Hybrid is nil", func(t *testing.T) {
		sr := &SearchRequest{Hybrid: nil}
		sr.validate()
		// Should not panic or set anything
		require.Nil(t, sr.Hybrid)
	})

	t.Run("Hybrid non-nil, Embedder empty", func(t *testing.T) {
		sr := &SearchRequest{Hybrid: &SearchRequestHybrid{Embedder: ""}}
		sr.validate()
		require.NotNil(t, sr.Hybrid)
		require.Equal(t, "default", sr.Hybrid.Embedder)
	})

	t.Run("Hybrid non-nil, Embedder set", func(t *testing.T) {
		sr := &SearchRequest{Hybrid: &SearchRequestHybrid{Embedder: "custom"}}
		sr.validate()
		require.NotNil(t, sr.Hybrid)
		require.Equal(t, "custom", sr.Hybrid.Embedder)
	})
}

func TestDocumentsQuery_JSONSerialization(t *testing.T) {
	t.Parallel()

	t.Run("DocumentsQuery with Sort field", func(t *testing.T) {
		dq := &DocumentsQuery{
			Limit:  10,
			Offset: 0,
			Fields: []string{"title", "id"},
			Sort:   []string{"title:asc", "id:desc"},
		}

		expected := `{"limit":10,"fields":["title","id"],"sort":["title:asc","id:desc"]}`

		raw, err := json.Marshal(dq)
		require.NoError(t, err)
		require.JSONEq(t, expected, string(raw))
	})

	t.Run("DocumentsQuery with all fields", func(t *testing.T) {
		dq := &DocumentsQuery{
			Offset:          5,
			Limit:           20,
			Fields:          []string{"title", "author"},
			Filter:          "rating > 4",
			RetrieveVectors: true,
			Ids:             []string{"1", "2", "3"},
			Sort:            []string{"rating:desc"},
		}

		expected := `{"offset":5,"limit":20,"fields":["title","author"],"filter":"rating > 4","retrieveVectors":true,"ids":["1","2","3"],"sort":["rating:desc"]}`

		raw, err := json.Marshal(dq)
		require.NoError(t, err)
		require.JSONEq(t, expected, string(raw))
	})
}
