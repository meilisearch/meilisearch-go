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
			name: "leader set",
			in: UpdateNetworkRequest{
				Leader: String("uuid"),
			},
			wantJSON: `{"leader": "uuid"}`,
		},
		{
			name: "leader explicitly null",
			in: UpdateNetworkRequest{
				Leader: Null[string](),
			},
			wantJSON: `{"leader": null}`,
		},
		{
			name: "version set",
			in: UpdateNetworkRequest{
				Version: String("uuid"),
			},
			wantJSON: `{"version": "uuid"}`,
		},
		{
			name: "version explicitly null",
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

func TestUpdateRemoteShard_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		in       UpdateRemoteShard
		wantJSON string
	}{
		{
			name:     "omit all when all fields are zero value",
			in:       UpdateRemoteShard{},
			wantJSON: `{}`,
		},
		{
			name:     "remotes set",
			in:       UpdateRemoteShard{Remotes: NewOpt([]string{"remote1", "remote2"})},
			wantJSON: `{"remotes":["remote1","remote2"]}`,
		},
		{
			name:     "remotes null",
			in:       UpdateRemoteShard{Remotes: Null[[]string]()},
			wantJSON: `{"remotes":null}`,
		},
		{
			name:     "addRemotes set",
			in:       UpdateRemoteShard{AddRemotes: NewOpt([]string{"remote3"})},
			wantJSON: `{"addRemotes":["remote3"]}`,
		},
		{
			name:     "addRemotes null",
			in:       UpdateRemoteShard{AddRemotes: Null[[]string]()},
			wantJSON: `{"addRemotes":null}`,
		},
		{
			name:     "removeRemotes set",
			in:       UpdateRemoteShard{RemoveRemotes: NewOpt([]string{"remote1"})},
			wantJSON: `{"removeRemotes":["remote1"]}`,
		},
		{
			name:     "removeRemotes null",
			in:       UpdateRemoteShard{RemoveRemotes: Null[[]string]()},
			wantJSON: `{"removeRemotes":null}`,
		},
		{
			name: "all fields set",
			in: UpdateRemoteShard{
				Remotes:       NewOpt([]string{"remote1", "remote2"}),
				AddRemotes:    NewOpt([]string{"remote3"}),
				RemoveRemotes: NewOpt([]string{"remote4"}),
			},
			wantJSON: `{"remotes":["remote1","remote2"],"addRemotes":["remote3"],"removeRemotes":["remote4"]}`,
		},
		{
			name: "mixed null and set",
			in: UpdateRemoteShard{
				Remotes:       Null[[]string](),
				AddRemotes:    NewOpt([]string{"remote1"}),
				RemoveRemotes: Null[[]string](),
			},
			wantJSON: `{"remotes":null,"addRemotes":["remote1"],"removeRemotes":null}`,
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

func TestUpdateNetwork_MarshalJSON_Shards(t *testing.T) {
	t.Parallel()

	type R = Opt[UpdateRemoteShard]

	tests := []struct {
		name     string
		in       UpdateNetworkRequest
		wantJSON string
	}{
		{
			name: "shards set with one shard",
			in: UpdateNetworkRequest{
				Shards: NewOpt(map[string]R{
					"shard1": NewOpt(UpdateRemoteShard{
						Remotes: NewOpt([]string{"remote1", "remote2"}),
					}),
				}),
			},
			wantJSON: `{"shards":{"shard1":{"remotes":["remote1","remote2"]}}}`,
		},
		{
			name: "shards set with multiple shards",
			in: UpdateNetworkRequest{
				Shards: NewOpt(map[string]R{
					"shard1": NewOpt(UpdateRemoteShard{
						Remotes: NewOpt([]string{"remote1"}),
					}),
					"shard2": NewOpt(UpdateRemoteShard{
						AddRemotes: NewOpt([]string{"remote2"}),
					}),
				}),
			},
			wantJSON: `{"shards":{"shard1":{"remotes":["remote1"]},"shard2":{"addRemotes":["remote2"]}}}`,
		},
		{
			name: "shards null",
			in: UpdateNetworkRequest{
				Shards: Null[map[string]R](),
			},
			wantJSON: `{"shards":null}`,
		},
		{
			name: "shards with null shard",
			in: UpdateNetworkRequest{
				Shards: NewOpt(map[string]R{
					"shard1": Null[UpdateRemoteShard](),
				}),
			},
			wantJSON: `{"shards":{"shard1":null}}`,
		},
		{
			name: "all fields including shards",
			in: UpdateNetworkRequest{
				Self:    String("primary"),
				Leader:  String("leader-uuid"),
				Version: String("v1"),
				Remotes: NewOpt(map[string]Opt[UpdateRemote]{
					"remote1": NewOpt(UpdateRemote{URL: String("https://remote1.example.com")}),
				}),
				Shards: NewOpt(map[string]R{
					"shard1": NewOpt(UpdateRemoteShard{
						Remotes: NewOpt([]string{"remote1"}),
					}),
				}),
			},
			wantJSON: `{
				"self":"primary",
				"leader":"leader-uuid",
				"remotes":{"remote1":{"url":"https://remote1.example.com"}},
				"shards":{"shard1":{"remotes":["remote1"]}},
				"version":"v1"
			}`,
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

func TestTimestampz_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		input          Timestampz
		expectedString string
		expectedTime   time.Time
	}{
		{
			name:           "positive timestamp",
			input:          1609459200, // 2021-01-01 00:00:00 UTC
			expectedString: "2021-01-01T00:00:00Z",
			expectedTime:   time.Unix(1609459200, 0).UTC(),
		},
		{
			name:           "large timestamp",
			input:          2147483647, // max int32
			expectedString: "2038-01-19T03:14:07Z",
			expectedTime:   time.Unix(2147483647, 0).UTC(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expectedString, tt.input.String())
			require.Equal(t, tt.expectedTime, tt.input.ToTime())
		})
	}
}

func TestSearchRequest_ValidateEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("Hybrid embedder already set to custom value", func(t *testing.T) {
		sr := &SearchRequest{
			Hybrid: &SearchRequestHybrid{
				Embedder:      "custom-embedder",
				SemanticRatio: 0.5,
			},
		}
		sr.validate()
		require.NotNil(t, sr.Hybrid)
		require.Equal(t, "custom-embedder", sr.Hybrid.Embedder)
		require.Equal(t, 0.5, sr.Hybrid.SemanticRatio)
	})

	t.Run("SearchRequest with no hybrid", func(t *testing.T) {
		sr := &SearchRequest{
			Query: "test query",
			Limit: 10,
		}
		sr.validate()
		require.Nil(t, sr.Hybrid)
	})
}

func TestUpdateNetworkRequest_RoundTrip(t *testing.T) {
	t.Parallel()

	original := UpdateNetworkRequest{
		Self:    String("primary-node"),
		Leader:  String("leader-uuid"),
		Version: String("v1.2.3"),
		Remotes: NewOpt(map[string]Opt[UpdateRemote]{
			"east": NewOpt(UpdateRemote{
				URL:          String("https://east.example.com"),
				SearchAPIKey: String("sek_east"),
				WriteAPIKey:  String("wek_east"),
			}),
			"west": NewOpt(UpdateRemote{
				URL:          String("https://west.example.com"),
				SearchAPIKey: Null[string](),
			}),
		}),
		Shards: NewOpt(map[string]Opt[UpdateRemoteShard]{
			"shard1": NewOpt(UpdateRemoteShard{
				Remotes: NewOpt([]string{"east", "west"}),
			}),
		}),
	}

	// Marshal
	data, err := original.MarshalJSON()
	require.NoError(t, err)

	// Verify JSON structure
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	// Check top-level fields
	require.Equal(t, "primary-node", result["self"])
	require.Equal(t, "leader-uuid", result["leader"])
	require.Equal(t, "v1.2.3", result["version"])
	require.Contains(t, result, "remotes")
	require.Contains(t, result, "shards")
}

func TestUpdateRemote_AllFieldsNull(t *testing.T) {
	t.Parallel()

	remote := UpdateRemote{
		URL:          Null[string](),
		SearchAPIKey: Null[string](),
		WriteAPIKey:  Null[string](),
	}

	data, err := remote.MarshalJSON()
	require.NoError(t, err)
	require.JSONEq(t, `{"url":null,"searchApiKey":null,"writeApiKey":null}`, string(data))
}

func TestUpdateRemoteShard_EmptySlices(t *testing.T) {
	t.Parallel()

	shard := UpdateRemoteShard{
		Remotes:       NewOpt([]string{}),
		AddRemotes:    NewOpt([]string{}),
		RemoveRemotes: NewOpt([]string{}),
	}

	data, err := shard.MarshalJSON()
	require.NoError(t, err)
	require.JSONEq(t, `{"remotes":[],"addRemotes":[],"removeRemotes":[]}`, string(data))
}

func TestMultiSearchFederation_Distinct(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		in       MultiSearchFederation
		wantJSON string
	}{
		{
			name:     "distinct omitted when empty",
			in:       MultiSearchFederation{},
			wantJSON: `{}`,
		},
		{
			name: "distinct set to attribute name",
			in: MultiSearchFederation{
				Distinct: "product_id",
			},
			wantJSON: `{"distinct":"product_id"}`,
		},
		{
			name: "distinct set alongside other federation options",
			in: MultiSearchFederation{
				Offset:   0,
				Limit:    20,
				Distinct: "sku",
			},
			wantJSON: `{"limit":20,"distinct":"sku"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.in)
			require.NoError(t, err)
			require.JSONEq(t, tt.wantJSON, string(got))
		})
	}
}
