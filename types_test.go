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

func TestUpdateNetwork_MarshalJSON(t *testing.T) {
	t.Parallel()

	type R = Opt[UpdateRemote]

	tests := []struct {
		name     string
		in       UpdateNetworkRequest
		wantJSON string
	}{
		{
			name:     "omit all when both fields are zero value",
			in:       UpdateNetworkRequest{},
			wantJSON: `{}`,
		},
		{
			name: "self set, remotes omitted",
			in: UpdateNetworkRequest{
				Self: String("primary-node"),
			},
			wantJSON: `{"self":"primary-node"}`,
		},
		{
			name: "self null, remotes omitted",
			in: UpdateNetworkRequest{
				Self: Null[string](),
			},
			wantJSON: `{"self":null}`,
		},
		{
			name: "remotes set (one valid remote, one null), self omitted",
			in: UpdateNetworkRequest{
				Remotes: NewOpt(map[string]R{
					"east": NewOpt(UpdateRemote{
						URL: String("https://east.example.com"),
					}),
					"west": Null[UpdateRemote](),
				}),
			},
			wantJSON: `{
				"remotes": {
					"east": { "url": "https://east.example.com" },
					"west": null
				}
			}`,
		},
		{
			name: "self set and remotes set",
			in: UpdateNetworkRequest{
				Self: String("primary"),
				Remotes: NewOpt(map[string]R{
					"a": NewOpt(UpdateRemote{URL: String("https://a.example.com"), SearchAPIKey: String("sek_a")}),
					"b": NewOpt(UpdateRemote{URL: Null[string]()}),
				}),
			},
			wantJSON: `{
				"self":"primary",
				"remotes": {
					"a": {"url":"https://a.example.com","searchApiKey":"sek_a"},
					"b": {"url":null}
				}
			}`,
		},
		{
			name: "remotes explicitly null",
			in: UpdateNetworkRequest{
				Self:    String("primary"),
				Remotes: Null[map[string]R](),
			},
			wantJSON: `{"self":"primary","remotes":null}`,
		},
		{
			name: "version set",
			in: UpdateNetworkRequest{
				Version: String("uuid"),
			},
			wantJSON: `{"version": "uuid"}`,
		},
		{
			name: "leader explicitly null",
			in: UpdateNetworkRequest{
				Version: Null[string](),
			},
			wantJSON: `{"version": null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.in.MarshalJSON()
			require.NoError(t, err)
			require.JSONEq(t, tt.wantJSON, string(got))
		})
	}
}

func TestRemote_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		in       UpdateRemote
		wantJSON string
	}{
		{
			name:     "omit all when both fields are zero value",
			in:       UpdateRemote{},
			wantJSON: `{}`,
		},
		{
			name:     "url set, searchApiKey omitted",
			in:       UpdateRemote{URL: String("https://east.example.com")},
			wantJSON: `{"url":"https://east.example.com"}`,
		},
		{
			name:     "url null, searchApiKey omitted",
			in:       UpdateRemote{URL: Null[string]()},
			wantJSON: `{"url":null}`,
		},
		{
			name:     "url set, searchApiKey null",
			in:       UpdateRemote{URL: String("https://east.example.com"), SearchAPIKey: Null[string]()},
			wantJSON: `{"url":"https://east.example.com","searchApiKey":null}`,
		},
		{
			name:     "both set",
			in:       UpdateRemote{URL: String("https://east.example.com"), SearchAPIKey: String("sek_abc")},
			wantJSON: `{"url":"https://east.example.com","searchApiKey":"sek_abc"}`,
		},
		{
			name:     "writeApiKey set",
			in:       UpdateRemote{WriteAPIKey: String("TEST-API-KEY")},
			wantJSON: `{"writeApiKey": "TEST-API-KEY"}`,
		},
		{
			name:     "writeApiKey null",
			in:       UpdateRemote{WriteAPIKey: Null[string]()},
			wantJSON: `{"writeApiKey": null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.in.MarshalJSON()
			require.NoError(t, err)
			require.JSONEq(t, tt.wantJSON, string(got))
		})
	}
}
