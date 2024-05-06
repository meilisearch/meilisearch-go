package meilisearch

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndex_SearchRaw(t *testing.T) {
	type args struct {
		UID        string
		PrimaryKey string
		client     *Client
		query      string
		request    *SearchRequest
	}

	tests := []struct {
		name    string
		args    args
		want    *SearchResponse
		wantErr bool
	}{
		{
			name: "TestIndexBasicSearch",
			args: args{
				UID:     "indexUID",
				client:  defaultClient,
				query:   "prince",
				request: &SearchRequest{},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(456), "title": "Le Petit Prince",
					},
					map[string]interface{}{
						"Tag": "Epic fantasy", "book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				EstimatedTotalHits: 2,
				Offset:             0,
				Limit:              20,
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

			gotRaw, err := i.SearchRaw(tt.args.query, tt.args.request)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, tt.want)
				return
			}

			require.NoError(t, err)
			// Unmarshall the raw response from SearchRaw into a SearchResponse
			var got SearchResponse
			err = json.Unmarshal(*gotRaw, &got)
			require.NoError(t, err, "error unmarshalling raw got SearchResponse")

			require.Equal(t, len(tt.want.Hits), len(got.Hits))
			for len := range got.Hits {
				require.Equal(t, tt.want.Hits[len].(map[string]interface{})["title"], got.Hits[len].(map[string]interface{})["title"])
				require.Equal(t, tt.want.Hits[len].(map[string]interface{})["book_id"], got.Hits[len].(map[string]interface{})["book_id"])
			}
			if tt.want.Hits[0].(map[string]interface{})["_formatted"] != nil {
				require.Equal(t, tt.want.Hits[0].(map[string]interface{})["_formatted"], got.Hits[0].(map[string]interface{})["_formatted"])
			}
			require.Equal(t, tt.want.EstimatedTotalHits, got.EstimatedTotalHits)
			require.Equal(t, tt.want.Offset, got.Offset)
			require.Equal(t, tt.want.Limit, got.Limit)
			require.Equal(t, tt.want.FacetDistribution, got.FacetDistribution)
		})
	}
}

func TestIndex_Search(t *testing.T) {
	type args struct {
		UID        string
		PrimaryKey string
		client     *Client
		query      string
		request    *SearchRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *SearchResponse
		wantErr bool
	}{
		{
			name: "TestIndexSearchWithEmptyRequest",
			args: args{
				UID:     "indexUID",
				client:  defaultClient,
				query:   "prince",
				request: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "TestIndexBasicSearch",
			args: args{
				UID:     "indexUID",
				client:  defaultClient,
				query:   "prince",
				request: &SearchRequest{},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(456), "title": "Le Petit Prince",
					},
					map[string]interface{}{
						"Tag": "Epic fantasy", "book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				EstimatedTotalHits: 2,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithCustomClient",
			args: args{
				UID:     "indexUID",
				client:  customClient,
				query:   "prince",
				request: &SearchRequest{},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(456), "title": "Le Petit Prince",
					},
					map[string]interface{}{
						"Tag": "Epic fantasy", "book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				EstimatedTotalHits: 2,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithLimit",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "prince",
				request: &SearchRequest{
					Limit: 1,
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(456), "title": "Le Petit Prince",
					},
				},
				EstimatedTotalHits: 2,
				Offset:             0,
				Limit:              1,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithPlaceholderSearch",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: &SearchRequest{
					PlaceholderSearch: true,
					Limit:             1,
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(123), "title": "Pride and Prejudice",
					},
				},
				EstimatedTotalHits: 20,
				Offset:             0,
				Limit:              1,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithOffset",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "prince",
				request: &SearchRequest{
					Offset: 1,
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				EstimatedTotalHits: 2,
				Offset:             1,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithAttributeToRetrieve",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "prince",
				request: &SearchRequest{
					AttributesToRetrieve: []string{"book_id", "title"},
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(456), "title": "Le Petit Prince",
					},
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				EstimatedTotalHits: 2,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithAttributeToSearchOn",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "prince",
				request: &SearchRequest{
					AttributesToSearchOn: []string{"title"},
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(456), "title": "Le Petit Prince",
					},
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				EstimatedTotalHits: 2,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithAttributesToCrop",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "to",
				request: &SearchRequest{
					AttributesToCrop: []string{"title"},
					CropLength:       2,
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(42), "title": "The Hitchhiker's Guide to the Galaxy",
						"_formatted": map[string]interface{}{
							"book_id": "42", "tag": "Epic fantasy", "title": "…Guide to…", "year": "1978",
						},
					},
				},
				EstimatedTotalHits: 1,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithAttributesToCropAndCustomCropMarker",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "to",
				request: &SearchRequest{
					AttributesToCrop: []string{"title"},
					CropLength:       2,
					CropMarker:       "(ꈍᴗꈍ)",
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(42), "title": "The Hitchhiker's Guide to the Galaxy",
						"_formatted": map[string]interface{}{
							"book_id": "42", "tag": "Epic fantasy", "title": "(ꈍᴗꈍ)Guide to(ꈍᴗꈍ)", "year": "1978",
						},
					},
				},
				EstimatedTotalHits: 1,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithAttributeToHighlight",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "prince",
				request: &SearchRequest{
					Limit:                 1,
					AttributesToHighlight: []string{"*"},
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(456), "title": "Le Petit Prince",
						"_formatted": map[string]interface{}{
							"book_id": "456", "tag": "Tale", "title": "Le Petit <em>Prince</em>", "year": "1943",
						},
					},
				},
				EstimatedTotalHits: 2,
				Offset:             0,
				Limit:              1,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithCustomPreAndPostHighlightTags",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "prince",
				request: &SearchRequest{
					Limit:                 1,
					AttributesToHighlight: []string{"*"},
					HighlightPreTag:       "(⊃｡•́‿•̀｡)⊃ ",
					HighlightPostTag:      " ⊂(´• ω •`⊂)",
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(456), "title": "Le Petit Prince",
						"_formatted": map[string]interface{}{
							"book_id": "456", "tag": "Tale", "title": "Le Petit (⊃｡•́‿•̀｡)⊃ Prince ⊂(´• ω •`⊂)", "year": "1943",
						},
					},
				},
				EstimatedTotalHits: 2,
				Offset:             0,
				Limit:              1,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithShowMatchesPosition",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "and",
				request: &SearchRequest{
					ShowMatchesPosition: true,
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(123), "title": "Pride and Prejudice",
					},
					map[string]interface{}{
						"book_id": float64(730), "title": "War and Peace",
					},
					map[string]interface{}{
						"book_id": float64(1032), "title": "Crime and Punishment",
					},
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				EstimatedTotalHits: 4,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithQuoteInQUery",
			args: args{
				UID:     "indexUID",
				client:  defaultClient,
				query:   "and \"harry\"",
				request: &SearchRequest{},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				EstimatedTotalHits: 1,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithCustomMatchingStrategyAll",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "le prince",
				request: &SearchRequest{
					Limit:                10,
					AttributesToRetrieve: []string{"book_id", "title"},
					MatchingStrategy:     "all",
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(456), "title": "Le Petit Prince",
					},
				},
				EstimatedTotalHits: 1,
				Offset:             0,
				Limit:              10,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithCustomMatchingStrategyLast",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "prince",
				request: &SearchRequest{
					Limit:                10,
					AttributesToRetrieve: []string{"book_id", "title"},
					MatchingStrategy:     "last",
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(456), "title": "Le Petit Prince",
					},
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				EstimatedTotalHits: 2,
				Offset:             0,
				Limit:              10,
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

			got, err := i.Search(tt.args.query, tt.args.request)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, tt.want)
				return
			}

			require.NoError(t, err)
			require.Equal(t, len(tt.want.Hits), len(got.Hits))
			for len := range got.Hits {
				require.Equal(t, tt.want.Hits[len].(map[string]interface{})["title"], got.Hits[len].(map[string]interface{})["title"])
				require.Equal(t, tt.want.Hits[len].(map[string]interface{})["book_id"], got.Hits[len].(map[string]interface{})["book_id"])
			}
			if tt.want.Hits[0].(map[string]interface{})["_formatted"] != nil {
				require.Equal(t, tt.want.Hits[0].(map[string]interface{})["_formatted"], got.Hits[0].(map[string]interface{})["_formatted"])
			}
			require.Equal(t, tt.want.EstimatedTotalHits, got.EstimatedTotalHits)
			require.Equal(t, tt.want.Offset, got.Offset)
			require.Equal(t, tt.want.Limit, got.Limit)
			require.Equal(t, tt.want.FacetDistribution, got.FacetDistribution)
		})
	}
}

func TestIndex_SearchFacets(t *testing.T) {
	type args struct {
		UID                  string
		PrimaryKey           string
		client               *Client
		query                string
		request              *SearchRequest
		filterableAttributes []string
	}
	tests := []struct {
		name    string
		args    args
		want    *SearchResponse
		wantErr bool
	}{
		{
			name: "TestIndexSearchWithEmptyRequest",
			args: args{
				UID:     "indexUID",
				client:  defaultClient,
				query:   "prince",
				request: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "TestIndexSearchWithFacets",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "prince",
				request: &SearchRequest{
					Facets: []string{"*"},
				},
				filterableAttributes: []string{"tag"},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(456), "title": "Le Petit Prince",
					},
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				EstimatedTotalHits: 2,
				Offset:             0,
				Limit:              20,
				FacetDistribution: map[string]interface{}(
					map[string]interface{}{
						"tag": map[string]interface{}{
							"Epic fantasy": float64(1),
							"Tale":         float64(1),
						},
					}),
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithFacetsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
				query:  "prince",
				request: &SearchRequest{
					Facets: []string{"*"},
				},
				filterableAttributes: []string{"tag"},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(456), "title": "Le Petit Prince",
					},
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				EstimatedTotalHits: 2,
				Offset:             0,
				Limit:              20,
				FacetDistribution: map[string]interface{}(
					map[string]interface{}{
						"tag": map[string]interface{}{
							"Epic fantasy": float64(1),
							"Tale":         float64(1),
						},
					}),
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithFacetsAndFacetsStats",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "prince",
				request: &SearchRequest{
					Facets: []string{"book_id"},
				},
				filterableAttributes: []string{"book_id"},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(456), "title": "Le Petit Prince",
					},
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				EstimatedTotalHits: 2,
				Offset:             0,
				Limit:              20,
				FacetDistribution: map[string]interface{}(
					map[string]interface{}{
						"book_id": map[string]interface{}{
							"4":   float64(1),
							"456": float64(1),
						},
					}),
				FacetStats: map[string]interface{}(
					map[string]interface{}{
						"book_id": map[string]interface{}{
							"max": float64(456),
							"min": float64(4),
						},
					}),
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

			updateFilter, err := i.UpdateFilterableAttributes(&tt.args.filterableAttributes)
			require.NoError(t, err)
			testWaitForTask(t, i, updateFilter)

			got, err := i.Search(tt.args.query, tt.args.request)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, tt.want)
				return
			}

			require.NoError(t, err)
			require.Equal(t, len(tt.want.Hits), len(got.Hits))

			for len := range got.Hits {
				require.Equal(t, tt.want.Hits[len].(map[string]interface{})["title"], got.Hits[len].(map[string]interface{})["title"])
				require.Equal(t, tt.want.Hits[len].(map[string]interface{})["book_id"], got.Hits[len].(map[string]interface{})["book_id"])
			}
			require.Equal(t, tt.want.EstimatedTotalHits, got.EstimatedTotalHits)
			require.Equal(t, tt.want.Offset, got.Offset)
			require.Equal(t, tt.want.Limit, got.Limit)
			require.Equal(t, tt.want.FacetDistribution, got.FacetDistribution)
			if tt.want.FacetStats != nil {
				require.Equal(t, tt.want.FacetStats, got.FacetStats)
			}
		})
	}
}

func TestIndex_SearchWithFilters(t *testing.T) {
	type args struct {
		UID                  string
		PrimaryKey           string
		client               *Client
		query                string
		filterableAttributes []string
		request              *SearchRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *SearchResponse
		wantErr bool
	}{
		{
			name: "TestIndexBasicSearchWithFilter",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "and",
				filterableAttributes: []string{
					"tag",
				},
				request: &SearchRequest{
					Filter: "tag = romance",
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(123), "title": "Pride and Prejudice",
					},
				},
				EstimatedTotalHits: 1,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithFilterInInt",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "and",
				filterableAttributes: []string{
					"year",
				},
				request: &SearchRequest{
					Filter: "year = 2005",
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				EstimatedTotalHits: 1,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithFilterArray",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "and",
				filterableAttributes: []string{
					"year",
				},
				request: &SearchRequest{
					Filter: []string{
						"year = 2005",
					},
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				EstimatedTotalHits: 1,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithFilterMultipleArray",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "and",
				filterableAttributes: []string{
					"year",
					"tag",
				},
				request: &SearchRequest{
					Filter: [][]string{
						{"year < 1850"},
						{"tag = romance"},
					},
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(123), "title": "Pride and Prejudice",
					},
				},
				EstimatedTotalHits: 1,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithMultipleFilter",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "prince",
				filterableAttributes: []string{
					"tag",
					"year",
				},
				request: &SearchRequest{
					Filter: "year > 1930",
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(456), "title": "Le Petit Prince",
					},
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				EstimatedTotalHits: 2,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithOneFilterAnd",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "",
				filterableAttributes: []string{
					"year",
				},
				request: &SearchRequest{
					Filter: "year < 1930 AND year > 1910",
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(742), "title": "The Great Gatsby",
					},
					map[string]interface{}{
						"book_id": float64(17), "title": "In Search of Lost Time",
					},
					map[string]interface{}{
						"book_id": float64(204), "title": "Ulysses",
					},
				},
				EstimatedTotalHits: 3,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithMultipleFilterAnd",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "",
				filterableAttributes: []string{
					"tag",
					"year",
				},
				request: &SearchRequest{
					Filter: "year < 1930 AND tag = Tale",
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(1), "title": "Alice In Wonderland",
					},
				},
				EstimatedTotalHits: 1,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithFilterOr",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "",
				filterableAttributes: []string{
					"year",
					"tag",
				},
				request: &SearchRequest{
					Filter: "year > 2000 OR tag = Tale",
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(456), "title": "Le Petit Prince",
					},
					map[string]interface{}{
						"book_id": float64(1), "title": "Alice In Wonderland",
					},
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				EstimatedTotalHits: 3,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithAttributeToHighlight",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "prince",
				filterableAttributes: []string{
					"book_id",
				},
				request: &SearchRequest{
					AttributesToHighlight: []string{"*"},
					Filter:                "book_id > 10",
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(456), "title": "Le Petit Prince",
					},
				},
				EstimatedTotalHits: 1,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithFilterContainingSpaces",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "and",
				filterableAttributes: []string{
					"tag",
				},
				request: &SearchRequest{
					Filter: "tag = 'Crime fiction'",
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(1032), "title": "Crime and Punishment",
					},
				},
				EstimatedTotalHits: 1,
				Offset:             0,
				Limit:              20,
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

			updateFilter, err := i.UpdateFilterableAttributes(&tt.args.filterableAttributes)
			require.NoError(t, err)
			testWaitForTask(t, i, updateFilter)

			got, err := i.Search(tt.args.query, tt.args.request)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, tt.want)
				return
			}

			require.NoError(t, err)
			require.Equal(t, len(tt.want.Hits), len(got.Hits))

			for len := range got.Hits {
				require.Equal(t, tt.want.Hits[len].(map[string]interface{})["title"], got.Hits[len].(map[string]interface{})["title"])
				require.Equal(t, tt.want.Hits[len].(map[string]interface{})["book_id"], got.Hits[len].(map[string]interface{})["book_id"])
			}
			require.Equal(t, tt.args.query, got.Query)
			require.Equal(t, tt.want.EstimatedTotalHits, got.EstimatedTotalHits)
			require.Equal(t, tt.want.Offset, got.Offset)
			require.Equal(t, tt.want.Limit, got.Limit)
			require.Equal(t, tt.want.FacetDistribution, got.FacetDistribution)
		})
	}
}

func TestIndex_SearchWithSort(t *testing.T) {
	type args struct {
		UID                string
		PrimaryKey         string
		client             *Client
		query              string
		sortableAttributes []string
		request            *SearchRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *SearchResponse
		wantErr bool
	}{
		{
			name: "TestIndexBasicSearchWithSortIntParameter",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "and",
				sortableAttributes: []string{
					"year",
				},
				request: &SearchRequest{
					Sort: []string{
						"year:asc",
					},
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(123), "title": "Pride and Prejudice",
					},
					map[string]interface{}{
						"book_id": float64(730), "title": "War and Peace",
					},
					map[string]interface{}{
						"book_id": float64(1032), "title": "Crime and Punishment",
					},
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				EstimatedTotalHits: 4,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexBasicSearchWithSortStringParameter",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "and",
				sortableAttributes: []string{
					"title",
				},
				request: &SearchRequest{
					Sort: []string{
						"title:asc",
					},
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(1032), "title": "Crime and Punishment",
					},
					map[string]interface{}{
						"book_id": float64(123), "title": "Pride and Prejudice",
					},
					map[string]interface{}{
						"book_id": float64(730), "title": "War and Peace",
					},
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				EstimatedTotalHits: 4,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexBasicSearchWithSortMultipleParameter",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "and",
				sortableAttributes: []string{
					"title",
					"year",
				},
				request: &SearchRequest{
					Sort: []string{
						"title:asc",
						"year:asc",
					},
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(1032), "title": "Crime and Punishment",
					},
					map[string]interface{}{
						"book_id": float64(123), "title": "Pride and Prejudice",
					},
					map[string]interface{}{
						"book_id": float64(730), "title": "War and Peace",
					},
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				EstimatedTotalHits: 4,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexBasicSearchWithSortMultipleParameterReverse",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "and",
				sortableAttributes: []string{
					"title",
					"year",
				},
				request: &SearchRequest{
					Sort: []string{
						"year:asc",
						"title:asc",
					},
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(123), "title": "Pride and Prejudice",
					},
					map[string]interface{}{
						"book_id": float64(730), "title": "War and Peace",
					},
					map[string]interface{}{
						"book_id": float64(1032), "title": "Crime and Punishment",
					},
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				EstimatedTotalHits: 4,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexBasicSearchWithSortMultipleParameterPlaceHolder",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "",
				sortableAttributes: []string{
					"title",
					"year",
				},
				request: &SearchRequest{
					Sort: []string{
						"year:asc",
						"title:asc",
					},
					Limit: 4,
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(56), "title": "The Divine Comedy",
					},
					map[string]interface{}{
						"book_id": float64(32), "title": "The Odyssey",
					},
					map[string]interface{}{
						"book_id": float64(69), "title": "Hamlet",
					},
					map[string]interface{}{
						"book_id": float64(7), "title": "Don Quixote",
					},
				},
				EstimatedTotalHits: 20,
				Offset:             0,
				Limit:              4,
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

			updateFilter, err := i.UpdateSortableAttributes(&tt.args.sortableAttributes)
			require.NoError(t, err)
			testWaitForTask(t, i, updateFilter)

			got, err := i.Search(tt.args.query, tt.args.request)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, tt.want)
				return
			}

			require.NoError(t, err)
			require.Equal(t, len(tt.want.Hits), len(got.Hits))

			for len := range got.Hits {
				require.Equal(t, tt.want.Hits[len].(map[string]interface{})["title"], got.Hits[len].(map[string]interface{})["title"])
				require.Equal(t, tt.want.Hits[len].(map[string]interface{})["book_id"], got.Hits[len].(map[string]interface{})["book_id"])
			}
			require.Equal(t, tt.want.EstimatedTotalHits, got.EstimatedTotalHits)
			require.Equal(t, tt.want.Offset, got.Offset)
			require.Equal(t, tt.want.Limit, got.Limit)
			require.Equal(t, tt.want.FacetDistribution, got.FacetDistribution)
		})
	}
}

func TestIndex_SearchOnNestedFileds(t *testing.T) {
	type args struct {
		UID                 string
		PrimaryKey          string
		client              *Client
		query               string
		request             *SearchRequest
		searchableAttribute []string
		sortableAttribute   []string
	}
	tests := []struct {
		name    string
		args    args
		want    *SearchResponse
		wantErr bool
	}{
		{
			name: "TestIndexBasicSearchOnNestedFields",
			args: args{
				UID:     "TestIndexBasicSearchOnNestedFields",
				client:  defaultClient,
				query:   "An awesome",
				request: &SearchRequest{},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"id": float64(5), "title": "The Hobbit",
						"info": map[string]interface{}{
							"comment":  "An awesome book",
							"reviewNb": float64(900),
						},
					},
				},
				EstimatedTotalHits: 1,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexBasicSearchOnNestedFieldsWithCustomClient",
			args: args{
				UID:     "TestIndexBasicSearchOnNestedFieldsWithCustomClient",
				client:  customClient,
				query:   "An awesome",
				request: &SearchRequest{},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"id": float64(5), "title": "The Hobbit",
						"info": map[string]interface{}{
							"comment":  "An awesome book",
							"reviewNb": float64(900),
						},
					},
				},
				EstimatedTotalHits: 1,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchOnMultipleNestedFields",
			args: args{
				UID:     "TestIndexSearchOnMultipleNestedFields",
				client:  defaultClient,
				query:   "french",
				request: &SearchRequest{},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"id": float64(2), "title": "Le Petit Prince",
						"info": map[string]interface{}{
							"comment":  "A french book",
							"reviewNb": float64(600),
						},
					},
					map[string]interface{}{
						"id": float64(3), "title": "Le Rouge et le Noir",
						"info": map[string]interface{}{
							"comment":  "Another french book",
							"reviewNb": float64(700),
						},
					},
				},
				EstimatedTotalHits: 2,
				Offset:             0,
				Limit:              20,
			},
		},
		{
			name: "TestIndexSearchOnNestedFieldsWithSearchableAttribute",
			args: args{
				UID:     "TestIndexSearchOnNestedFieldsWithSearchableAttribute",
				client:  defaultClient,
				query:   "An awesome",
				request: &SearchRequest{},
				searchableAttribute: []string{
					"title", "info.comment",
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"id": float64(5), "title": "The Hobbit",
						"info": map[string]interface{}{
							"comment":  "An awesome book",
							"reviewNb": float64(900),
						},
					},
				},
				EstimatedTotalHits: 1,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchOnNestedFieldsWithSortableAttribute",
			args: args{
				UID:    "TestIndexSearchOnNestedFieldsWithSortableAttribute",
				client: defaultClient,
				query:  "An awesome",
				request: &SearchRequest{
					Sort: []string{
						"info.reviewNb:desc",
					},
				},
				searchableAttribute: []string{
					"title", "info.comment",
				},
				sortableAttribute: []string{
					"info.reviewNb",
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"id": float64(5), "title": "The Hobbit",
						"info": map[string]interface{}{
							"comment":  "An awesome book",
							"reviewNb": float64(900),
						},
					},
				},
				EstimatedTotalHits: 1,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexWithNestedFields(tt.args.UID)
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			if tt.args.searchableAttribute != nil {
				gotTask, err := i.UpdateSearchableAttributes(&tt.args.searchableAttribute)
				require.NoError(t, err)
				testWaitForTask(t, i, gotTask)
			}

			if tt.args.sortableAttribute != nil {
				gotTask, err := i.UpdateSortableAttributes(&tt.args.sortableAttribute)
				require.NoError(t, err)
				testWaitForTask(t, i, gotTask)
			}

			got, err := i.Search(tt.args.query, tt.args.request)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, tt.want)
				return
			}

			require.NoError(t, err)
			require.Equal(t, len(tt.want.Hits), len(got.Hits))
			for len := range got.Hits {
				require.Equal(t, tt.want.Hits[len], got.Hits[len])
			}
			require.Equal(t, tt.want.EstimatedTotalHits, got.EstimatedTotalHits)
			require.Equal(t, tt.want.Offset, got.Offset)
			require.Equal(t, tt.want.Limit, got.Limit)
			require.Equal(t, tt.want.FacetDistribution, got.FacetDistribution)
		})
	}
}

func TestIndex_SearchWithPagination(t *testing.T) {
	type args struct {
		UID        string
		PrimaryKey string
		client     *Client
		query      string
		request    *SearchRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *SearchResponse
		wantErr bool
	}{
		{
			name: "TestIndexBasicSearchWithHitsPerPage",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "and",
				request: &SearchRequest{
					HitsPerPage: 10,
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(123), "title": "Pride and Prejudice",
					},
					map[string]interface{}{
						"book_id": float64(730), "title": "War and Peace",
					},
					map[string]interface{}{
						"book_id": float64(1032), "title": "Crime and Punishment",
					},
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				HitsPerPage: 10,
				Page:        1,
				TotalHits:   4,
				TotalPages:  1,
			},
			wantErr: false,
		},
		{
			name: "TestIndexBasicSearchWithPage",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "and",
				request: &SearchRequest{
					Page: 1,
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(123), "title": "Pride and Prejudice",
					},
					map[string]interface{}{
						"book_id": float64(730), "title": "War and Peace",
					},
					map[string]interface{}{
						"book_id": float64(1032), "title": "Crime and Punishment",
					},
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				HitsPerPage: 20,
				Page:        1,
				TotalHits:   4,
				TotalPages:  1,
			},
			wantErr: false,
		},
		{
			name: "TestIndexBasicSearchWithPageAndHitsPerPage",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "and",
				request: &SearchRequest{
					HitsPerPage: 10,
					Page:        1,
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(123), "title": "Pride and Prejudice",
					},
					map[string]interface{}{
						"book_id": float64(730), "title": "War and Peace",
					},
					map[string]interface{}{
						"book_id": float64(1032), "title": "Crime and Punishment",
					},
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				HitsPerPage: 10,
				Page:        1,
				TotalHits:   4,
				TotalPages:  1,
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

			got, err := i.Search(tt.args.query, tt.args.request)
			require.NoError(t, err)
			require.Equal(t, len(tt.want.Hits), len(got.Hits))

			for len := range got.Hits {
				require.Equal(t, tt.want.Hits[len].(map[string]interface{})["title"], got.Hits[len].(map[string]interface{})["title"])
				require.Equal(t, tt.want.Hits[len].(map[string]interface{})["book_id"], got.Hits[len].(map[string]interface{})["book_id"])
			}
			require.Equal(t, tt.args.query, got.Query)
			require.Equal(t, tt.want.HitsPerPage, got.HitsPerPage)
			require.Equal(t, tt.want.Page, got.Page)
			require.Equal(t, tt.want.TotalHits, got.TotalHits)
			require.Equal(t, tt.want.TotalPages, got.TotalPages)
		})
	}
}

func TestIndex_SearchWithShowRankingScore(t *testing.T) {
	type args struct {
		UID        string
		PrimaryKey string
		client     *Client
		query      string
		request    SearchRequest
	}
	testArg := args{
		UID:    "indexUID",
		client: defaultClient,
		query:  "and",
		request: SearchRequest{
			ShowRankingScore: true,
		},
	}
	SetUpIndexForFaceting()
	c := testArg.client
	i := c.Index(testArg.UID)
	t.Cleanup(cleanup(c))

	got, err := i.Search(testArg.query, &testArg.request)
	require.NoError(t, err)
	require.NotNil(t, got.Hits[0].(map[string]interface{})["_rankingScore"])
}

func TestIndex_SearchWithShowRankingScoreDetails(t *testing.T) {
	type args struct {
		UID        string
		PrimaryKey string
		client     *Client
		query      string
		request    SearchRequest
	}
	testArg := args{
		UID:    "indexUID",
		client: defaultClient,
		query:  "and",
		request: SearchRequest{
			ShowRankingScoreDetails: true,
		},
	}
	SetUpIndexForFaceting()
	c := testArg.client
	i := c.Index(testArg.UID)
	t.Cleanup(cleanup(c))

	got, err := i.Search(testArg.query, &testArg.request)
	require.NoError(t, err)
	require.NotNil(t, got.Hits[0].(map[string]interface{})["_rankingScoreDetails"])
}

func TestIndex_SearchWithVectorStore(t *testing.T) {
	type args struct {
		UID        string
		PrimaryKey string
		client     *Client
		query      string
		request    SearchRequest
	}
	testArg := args{
		UID:    "indexUID",
		client: defaultClient,
		query:  "Pride and Prejudice",
		request: SearchRequest{
			Vector: []float32{0.1, 0.2, 0.3},
			Hybrid: &SearchRequestHybrid{
				SemanticRatio: 0.5,
				Embedder:      "default",
			},
		},
	}

	i, err := SetUpIndexWithVector(testArg.UID)
	if err != nil {
		t.Fatal(err)
	}

	c := testArg.client
	t.Cleanup(cleanup(c))

	got, err := i.Search(testArg.query, &testArg.request)
	require.NoError(t, err)

	for _, hit := range got.Hits {
		hit := hit.(map[string]interface{})
		require.NotNil(t, hit["_vectors"])
		vectors := hit["_vectors"].(map[string]interface{})

		require.NotNil(t, vectors["default"])
		def := vectors["default"].([]interface{})
		require.Equal(t, 3, len(def))
	}
}
