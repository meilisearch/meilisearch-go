package meilisearch

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndex_Search(t *testing.T) {
	type args struct {
		UID        string
		PrimaryKey string
		client     *Client
		query      string
		request    SearchRequest
	}
	tests := []struct {
		name string
		args args
		want *SearchResponse
	}{
		{
			name: "TestIndexBasicSearch",
			args: args{
				UID:     "indexUID",
				client:  defaultClient,
				query:   "prince",
				request: SearchRequest{},
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
				NbHits:           2,
				Offset:           0,
				Limit:            20,
				ExhaustiveNbHits: false,
			},
		},
		{
			name: "TestIndexSearchWithCustomClient",
			args: args{
				UID:     "indexUID",
				client:  customClient,
				query:   "prince",
				request: SearchRequest{},
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
				NbHits:           2,
				Offset:           0,
				Limit:            20,
				ExhaustiveNbHits: false,
			},
		},
		{
			name: "TestIndexSearchWithLimit",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "prince",
				request: SearchRequest{
					Limit: 1,
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(456), "title": "Le Petit Prince",
					},
				},
				NbHits:           2,
				Offset:           0,
				Limit:            1,
				ExhaustiveNbHits: false,
			},
		},
		{
			name: "TestIndexSearchWithPlaceholderSearch",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: SearchRequest{
					PlaceholderSearch: true,
					Limit:             1,
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(1), "title": "Alice In Wonderland",
					},
				},
				NbHits:           20,
				Offset:           0,
				Limit:            1,
				ExhaustiveNbHits: false,
			},
		},
		{
			name: "TestIndexSearchWithOffset",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "prince",
				request: SearchRequest{
					Offset: 1,
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				NbHits:           2,
				Offset:           1,
				Limit:            20,
				ExhaustiveNbHits: false,
			},
		},
		{
			name: "TestIndexSearchWithAttributeToRetrieve",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "prince",
				request: SearchRequest{
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
				NbHits:           2,
				Offset:           0,
				Limit:            20,
				ExhaustiveNbHits: false,
			},
		},
		{
			name: "TestIndexSearchWithAttributesToCrop",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "to",
				request: SearchRequest{
					AttributesToCrop: []string{"title"},
					CropLength:       7,
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(42), "title": "The Hitchhiker's Guide to the Galaxy",
					},
				},
				NbHits:           1,
				Offset:           0,
				Limit:            20,
				ExhaustiveNbHits: false,
			},
		},
		{
			name: "TestIndexSearchWithMatches",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "and",
				request: SearchRequest{
					Matches: true,
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
				NbHits:           4,
				Offset:           0,
				Limit:            20,
				ExhaustiveNbHits: false,
			},
		},
		{
			name: "TestIndexSearchWithQuoteInQUery",
			args: args{
				UID:     "indexUID",
				client:  defaultClient,
				query:   "and \"harry\"",
				request: SearchRequest{},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				NbHits:           1,
				Offset:           0,
				Limit:            20,
				ExhaustiveNbHits: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			got, err := i.Search(tt.args.query, &tt.args.request)
			require.NoError(t, err)
			require.Equal(t, len(tt.want.Hits), len(got.Hits))
			for len := range got.Hits {
				require.Equal(t, tt.want.Hits[len].(map[string]interface{})["title"], got.Hits[len].(map[string]interface{})["title"])
				require.Equal(t, tt.want.Hits[len].(map[string]interface{})["book_id"], got.Hits[len].(map[string]interface{})["book_id"])
			}
			require.Equal(t, tt.want.NbHits, got.NbHits)
			require.Equal(t, tt.want.Offset, got.Offset)
			require.Equal(t, tt.want.Limit, got.Limit)
			require.Equal(t, tt.want.ExhaustiveNbHits, got.ExhaustiveNbHits)
			require.Equal(t, tt.want.FacetsDistribution, got.FacetsDistribution)
			require.Equal(t, tt.want.ExhaustiveFacetsCount, got.ExhaustiveFacetsCount)

			deleteAllIndexes(c)
		})
	}
}

func TestIndex_SearchFacets(t *testing.T) {
	type args struct {
		UID                  string
		PrimaryKey           string
		client               *Client
		query                string
		request              SearchRequest
		filterableAttributes []string
	}
	tests := []struct {
		name string
		args args
		want *SearchResponse
	}{
		{
			name: "TestIndexSearchWithFacetsDistribution",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "prince",
				request: SearchRequest{
					FacetsDistribution: []string{"*"},
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
				NbHits:           2,
				Offset:           0,
				Limit:            20,
				ExhaustiveNbHits: false,
				FacetsDistribution: map[string]interface{}(
					map[string]interface{}{
						"tag": map[string]interface{}{
							"Epic fantasy": float64(1),
							"Tale":         float64(1),
						},
					}),
				ExhaustiveFacetsCount: interface{}(false),
			},
		},
		{
			name: "TestIndexSearchWithFacetsDistributionWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
				query:  "prince",
				request: SearchRequest{
					FacetsDistribution: []string{"*"},
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
				NbHits:           2,
				Offset:           0,
				Limit:            20,
				ExhaustiveNbHits: false,
				FacetsDistribution: map[string]interface{}(
					map[string]interface{}{
						"tag": map[string]interface{}{
							"Epic fantasy": float64(1),
							"Tale":         float64(1),
						},
					}),
				ExhaustiveFacetsCount: interface{}(false),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			updateFilter, err := i.UpdateFilterableAttributes(&tt.args.filterableAttributes)
			require.NoError(t, err)
			i.DefaultWaitForPendingUpdate(updateFilter)

			got, err := i.Search(tt.args.query, &tt.args.request)
			require.NoError(t, err)
			require.Equal(t, len(tt.want.Hits), len(got.Hits))
			for len := range got.Hits {
				require.Equal(t, tt.want.Hits[len].(map[string]interface{})["title"], got.Hits[len].(map[string]interface{})["title"])
				require.Equal(t, tt.want.Hits[len].(map[string]interface{})["book_id"], got.Hits[len].(map[string]interface{})["book_id"])
			}
			require.Equal(t, tt.want.NbHits, got.NbHits)
			require.Equal(t, tt.want.Offset, got.Offset)
			require.Equal(t, tt.want.Limit, got.Limit)
			require.Equal(t, tt.want.ExhaustiveNbHits, got.ExhaustiveNbHits)

			require.Equal(t, tt.want.FacetsDistribution, got.FacetsDistribution)
			require.Equal(t, tt.want.ExhaustiveFacetsCount, got.ExhaustiveFacetsCount)

			deleteAllIndexes(c)
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
		request              SearchRequest
	}
	tests := []struct {
		name string
		args args
		want *SearchResponse
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
				request: SearchRequest{
					Filter: "tag = romance",
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(123), "title": "Pride and Prejudice",
					},
				},
				NbHits:           1,
				Offset:           0,
				Limit:            20,
				ExhaustiveNbHits: false,
			},
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
				request: SearchRequest{
					Filter: "year = 2005",
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				NbHits:           1,
				Offset:           0,
				Limit:            20,
				ExhaustiveNbHits: false,
			},
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
				request: SearchRequest{
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
				NbHits:           1,
				Offset:           0,
				Limit:            20,
				ExhaustiveNbHits: false,
			},
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
				request: SearchRequest{
					Filter: [][]string{
						[]string{"year < 1850"},
						[]string{"tag = romance"},
					},
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(123), "title": "Pride and Prejudice",
					},
				},
				NbHits:           1,
				Offset:           0,
				Limit:            20,
				ExhaustiveNbHits: false,
			},
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
				request: SearchRequest{
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
				NbHits:           2,
				Offset:           0,
				Limit:            20,
				ExhaustiveNbHits: false,
			},
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
				request: SearchRequest{
					Filter: "year < 1930 AND year > 1910",
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(17), "title": "In Search of Lost Time",
					},
					map[string]interface{}{
						"book_id": float64(204), "title": "Ulysses",
					},
					map[string]interface{}{
						"book_id": float64(742), "title": "The Great Gatsby",
					},
				},
				NbHits:           3,
				Offset:           0,
				Limit:            20,
				ExhaustiveNbHits: false,
			},
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
				request: SearchRequest{
					Filter: "year < 1930 AND tag = Tale",
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(1), "title": "Alice In Wonderland",
					},
				},
				NbHits:           1,
				Offset:           0,
				Limit:            20,
				ExhaustiveNbHits: false,
			},
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
				request: SearchRequest{
					Filter: "year > 2000 OR tag = Tale",
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(1), "title": "Alice In Wonderland",
					},
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
					map[string]interface{}{
						"book_id": float64(456), "title": "Le Petit Prince",
					},
				},
				NbHits:           3,
				Offset:           0,
				Limit:            20,
				ExhaustiveNbHits: false,
			},
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
				request: SearchRequest{
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
				NbHits:           1,
				Offset:           0,
				Limit:            20,
				ExhaustiveNbHits: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			updateFilter, err := i.UpdateFilterableAttributes(&tt.args.filterableAttributes)
			require.NoError(t, err)
			i.DefaultWaitForPendingUpdate(updateFilter)

			got, err := i.Search(tt.args.query, &tt.args.request)
			require.NoError(t, err)
			require.Equal(t, len(tt.want.Hits), len(got.Hits))

			for len := range got.Hits {
				require.Equal(t, tt.want.Hits[len].(map[string]interface{})["title"], got.Hits[len].(map[string]interface{})["title"])
				require.Equal(t, tt.want.Hits[len].(map[string]interface{})["book_id"], got.Hits[len].(map[string]interface{})["book_id"])
			}
			require.Equal(t, tt.want.NbHits, got.NbHits)
			require.Equal(t, tt.want.Offset, got.Offset)
			require.Equal(t, tt.want.Limit, got.Limit)
			require.Equal(t, tt.want.ExhaustiveNbHits, got.ExhaustiveNbHits)
			require.Equal(t, tt.want.FacetsDistribution, got.FacetsDistribution)
			require.Equal(t, tt.want.ExhaustiveFacetsCount, got.ExhaustiveFacetsCount)

			deleteAllIndexes(c)
		})
	}
}
