package meilisearch

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndex_FacetSearch(t *testing.T) {
	type args struct {
		UID                  string
		PrimaryKey           string
		client               *Client
		request              *FacetSearchRequest
		filterableAttributes []string
	}

	tests := []struct {
		name    string
		args    args
		want    *FacetSearchResponse
		wantErr bool
	}{
		{
			name: "TestIndexBasicFacetSearch",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: &FacetSearchRequest{
					FacetName:  "tag",
					FacetQuery: "Novel",
				},
				filterableAttributes: []string{"tag"},
			},
			want: &FacetSearchResponse{
				FacetHits: []interface{}{
					map[string]interface{}{
						"value": "Novel", "count": float64(5),
					},
				},
				FacetQuery: "Novel",
			},
			wantErr: false,
		},
		{
			name: "TestIndexFacetSearchWithFilter",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: &FacetSearchRequest{
					FacetName:  "tag",
					FacetQuery: "Novel",
					Filter:     "tag = 'Novel'",
				},
				filterableAttributes: []string{"tag"},
			},
			want: &FacetSearchResponse{
				FacetHits: []interface{}{
					map[string]interface{}{
						"value": "Novel", "count": float64(5),
					},
				},
				FacetQuery: "Novel",
			},
			wantErr: false,
		},
		{
			name: "TestIndexFacetSearchWithMatchingStrategy",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: &FacetSearchRequest{
					FacetName:        "tag",
					FacetQuery:       "Novel",
					MatchingStrategy: "frequency",
				},
				filterableAttributes: []string{"tag"},
			},
			want: &FacetSearchResponse{
				FacetHits: []interface{}{
					map[string]interface{}{
						"value": "Novel", "count": float64(5),
					},
				},
				FacetQuery: "Novel",
			},
			wantErr: false,
		},
		{
			name: "TestIndexFacetSearchWithAttributesToSearchOn",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: &FacetSearchRequest{
					FacetName:            "tag",
					FacetQuery:           "Novel",
					AttributesToSearchOn: []string{"tag"},
				},
				filterableAttributes: []string{"tag"},
			},
			want: &FacetSearchResponse{
				FacetHits: []interface{}{
					map[string]interface{}{
						"value": "Novel", "count": float64(5),
					},
				},
				FacetQuery: "Novel",
			},
			wantErr: false,
		},
		{
			name: "TestIndexFacetSearchWithNoFacetSearchRequest",
			args: args{
				UID:     "indexUID",
				client:  defaultClient,
				request: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "TestIndexFacetSearchWithNoFacetName",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: &FacetSearchRequest{
					FacetQuery: "Novel",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "TestIndexFacetSearchWithNoFacetQuery",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: &FacetSearchRequest{
					FacetName: "tag",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "TestIndexFacetSearchWithNoFilterableAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: &FacetSearchRequest{
					FacetName:  "tag",
					FacetQuery: "Novel",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "TestIndexFacetSearchWithQ",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: &FacetSearchRequest{
					Q:         "query",
					FacetName: "tag",
				},
				filterableAttributes: []string{"tag"},
			},
			want: &FacetSearchResponse{
				FacetHits:  []interface{}{},
				FacetQuery: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			if len(tt.args.filterableAttributes) > 0 {
				updateFilter, err := i.UpdateFilterableAttributes(&tt.args.filterableAttributes)
				require.NoError(t, err)
				testWaitForTask(t, i, updateFilter)
			}

			gotRaw, err := i.FacetSearch(tt.args.request)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, gotRaw)
				return
			}

			require.NoError(t, err)
			// Unmarshall the raw response from FacetSearch into a FacetSearchResponse
			var got FacetSearchResponse
			err = json.Unmarshal(*gotRaw, &got)
			require.NoError(t, err, "error unmarshalling raw got FacetSearchResponse")

			require.Equal(t, len(tt.want.FacetHits), len(got.FacetHits))
			for len := range got.FacetHits {
				require.Equal(t, tt.want.FacetHits[len].(map[string]interface{})["value"], got.FacetHits[len].(map[string]interface{})["value"])
				require.Equal(t, tt.want.FacetHits[len].(map[string]interface{})["count"], got.FacetHits[len].(map[string]interface{})["count"])
			}
			require.Equal(t, tt.want.FacetQuery, got.FacetQuery)
		})
	}
}
