package meilisearch

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
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
