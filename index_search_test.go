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
						"book_id": float64(123), "title": "Pride and Prejudice",
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
			name: "TestIndexSearchWithFilters",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "and",
				request: SearchRequest{
					Filters: "tag = \"Romance\"",
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
				NbHits:           4,
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
				request: SearchRequest{
					AttributesToHighlight: []string{"*"},
					Filters:               "book_id > 10",
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
		UID        string
		PrimaryKey string
		client     *Client
		query      string
		request    SearchRequest
		facet      []string
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
				facet: []string{"tag"},
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
							"Crime fiction":        float64(0),
							"Epic":                 float64(0),
							"Epic fantasy":         float64(1),
							"Historical fiction":   float64(0),
							"Modernist literature": float64(0),
							"Novel":                float64(0),
							"Tale":                 float64(1),
							"Romance":              float64(0),
							"Satiric":              float64(0),
							"Tragedy":              float64(0),
						},
					}),
				ExhaustiveFacetsCount: interface{}(true),
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
				facet: []string{"tag"},
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
							"Crime fiction":        float64(0),
							"Epic":                 float64(0),
							"Epic fantasy":         float64(1),
							"Historical fiction":   float64(0),
							"Modernist literature": float64(0),
							"Novel":                float64(0),
							"Tale":                 float64(1),
							"Romance":              float64(0),
							"Satiric":              float64(0),
							"Tragedy":              float64(0),
						},
					}),
				ExhaustiveFacetsCount: interface{}(true),
			},
		},
		{
			name: "TestIndexSearchWithFacetsDistributionWithTag",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "prince",
				request: SearchRequest{
					FacetFilters: []string{"tag:Epic fantasy"},
				},
				facet: []string{"tag", "title"},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				NbHits:                1,
				Offset:                0,
				Limit:                 20,
				ExhaustiveNbHits:      false,
				FacetsDistribution:    nil,
				ExhaustiveFacetsCount: interface{}(nil),
			},
		},
		{
			name: "TestIndexSearchWithFacetsDistributionWithTagAndOneFacet",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				query:  "prince",
				request: SearchRequest{
					FacetFilters: []string{"tag:Epic fantasy"},
				},
				facet: []string{"tag"},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				NbHits:                1,
				Offset:                0,
				Limit:                 20,
				ExhaustiveNbHits:      false,
				FacetsDistribution:    nil,
				ExhaustiveFacetsCount: interface{}(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			update, _ := i.UpdateAttributesForFaceting(&tt.args.facet)
			i.DefaultWaitForPendingUpdate(update)

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
			if got.FacetsDistribution != nil {
				require.Equal(t, tt.want.FacetsDistribution.(map[string]interface{})["tag"].(map[string]interface{})["Epic fantasy"], got.FacetsDistribution.(map[string]interface{})["tag"].(map[string]interface{})["Epic fantasy"])
				require.Equal(t, tt.want.FacetsDistribution.(map[string]interface{})["tag"].(map[string]interface{})["Tragedy"], got.FacetsDistribution.(map[string]interface{})["tag"].(map[string]interface{})["Tragedy"])
				require.Equal(t, tt.want.FacetsDistribution.(map[string]interface{})["tag"].(map[string]interface{})["Romance"], got.FacetsDistribution.(map[string]interface{})["tag"].(map[string]interface{})["Romance"])
			}
			require.Equal(t, tt.want.ExhaustiveFacetsCount, got.ExhaustiveFacetsCount)

			deleteAllIndexes(c)
		})
	}
}
