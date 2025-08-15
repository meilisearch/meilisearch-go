package meilisearch

import (
	"encoding/json"
	"testing"
	"time"

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

func TestTimestampz_MarshalJSON(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    Timestampz
		expected string
	}{
		{"epoch", Timestampz(0), `"1970-01-01T00:00:00Z"`},
		{"now", Timestampz(time.Now().UTC().Unix()), ""}, // will check RFC3339
		{"negative", Timestampz(-1), `"1969-12-31T23:59:59Z"`},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b, err := c.input.MarshalJSON()
			require.NoError(t, err)
			if c.expected != "" {
				require.Equal(t, c.expected, string(b))
			} else {
				// Check RFC3339 format
				var s string
				require.NoError(t, json.Unmarshal(b, &s))
				_, err := time.Parse(time.RFC3339, s)
				require.NoError(t, err)
			}
		})
	}
}

func TestTimestampz_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    string
		wantUnix int64
		wantErr  bool
	}{
		{"epoch", `"1970-01-01T00:00:00Z"`, 0, false},
		{"negative", `"1969-12-31T23:59:59Z"`, -1, false},
		{"invalid format", `"not-a-date"`, 0, true},
		{"not a string", `12345`, 0, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ts Timestampz
			err := ts.UnmarshalJSON([]byte(c.input))
			if c.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, Timestampz(c.wantUnix), ts)
			}
		})
	}
}

func TestTimestampz_String(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    Timestampz
		expected string
	}{
		{0, "1970-01-01T00:00:00Z"},
		{-1, "1969-12-31T23:59:59Z"},
	}

	for _, c := range cases {
		require.Equal(t, c.expected, c.input.String())
	}
}

func TestTimestampz_ToTime(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    Timestampz
		expected time.Time
	}{
		{0, time.Unix(0, 0).UTC()},
		{-1, time.Unix(-1, 0).UTC()},
	}

	for _, c := range cases {
		require.Equal(t, c.expected, c.input.ToTime())
	}
}
