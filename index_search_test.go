package meilisearch

import (
	"crypto/tls"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIndex_SearchWithContentEncoding(t *testing.T) {
	tests := []struct {
		Name            string
		ContentEncoding ContentEncoding
		Query           string
		Request         *SearchRequest
		FacetRequest    *FacetSearchRequest
		Response        *SearchResponse
		FacetResponse   *FacetSearchResponse
	}{
		{
			Name:            "SearchResultWithGzipEncoding",
			ContentEncoding: GzipEncoding,
			Query:           "prince",
			Request: &SearchRequest{
				IndexUID: "indexUID",
			},
			FacetRequest: &FacetSearchRequest{
				FacetName:  "tag",
				FacetQuery: "Novel",
			},
			FacetResponse: &FacetSearchResponse{
				FacetHits: []interface{}{
					map[string]interface{}{
						"value": "Novel", "count": float64(5),
					},
				},
				FacetQuery: "Novel",
			},
			Response: &SearchResponse{
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
		},
		{
			Name:            "SearchResultWithDeflateEncoding",
			ContentEncoding: DeflateEncoding,
			Query:           "prince",
			Request: &SearchRequest{
				IndexUID: "indexUID",
			},
			Response: &SearchResponse{
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
			FacetRequest: &FacetSearchRequest{
				FacetName:  "tag",
				FacetQuery: "Novel",
			},
			FacetResponse: &FacetSearchResponse{
				FacetHits: []interface{}{
					map[string]interface{}{
						"value": "Novel", "count": float64(5),
					},
				},
				FacetQuery: "Novel",
			},
		},
		{
			Name:            "SearchResultWithBrotliEncoding",
			ContentEncoding: BrotliEncoding,
			Query:           "prince",
			Request: &SearchRequest{
				IndexUID: "indexUID",
			},
			Response: &SearchResponse{
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
			FacetRequest: &FacetSearchRequest{
				FacetName:  "tag",
				FacetQuery: "Novel",
			},
			FacetResponse: &FacetSearchResponse{
				FacetHits: []interface{}{
					map[string]interface{}{
						"value": "Novel", "count": float64(5),
					},
				},
				FacetQuery: "Novel",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			sv := setup(t, "", WithContentEncoding(tt.ContentEncoding, DefaultCompression))
			setUpIndexForFaceting(sv)
			i := sv.Index(tt.Request.IndexUID)
			t.Cleanup(cleanup(sv))

			got, err := i.Search(tt.Query, tt.Request)
			require.NoError(t, err)
			require.Equal(t, len(tt.Response.Hits), len(got.Hits))

			gotJson, err := i.SearchRaw(tt.Query, tt.Request)
			require.NoError(t, err)

			var resp SearchResponse
			err = json.Unmarshal(*gotJson, &resp)
			require.NoError(t, err, "error unmarshalling raw got SearchResponse")
			require.Equal(t, len(tt.Response.Hits), len(resp.Hits))

			task, err := i.UpdateFilterableAttributes(&[]string{"tag"})
			require.NoError(t, err)
			testWaitForTask(t, i, task)

			gotJson, err = i.FacetSearch(tt.FacetRequest)
			require.NoError(t, err)
			var gotFacet FacetSearchResponse
			err = json.Unmarshal(*gotJson, &gotFacet)
			require.NoError(t, err, "error unmarshalling raw got SearchResponse")
			require.NoError(t, err)
			require.Equal(t, len(gotFacet.FacetHits), len(tt.FacetResponse.FacetHits))
		})
	}
}

func TestIndex_SearchRaw(t *testing.T) {
	sv := setup(t, "")

	type args struct {
		UID        string
		PrimaryKey string
		client     ServiceManager
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
				UID:    "indexUID",
				client: sv,
				query:  "prince",
				request: &SearchRequest{
					IndexUID: "foobar",
				},
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
			name: "TestNullRequestInSearchRow",
			args: args{
				UID:     "indexUID",
				client:  sv,
				query:   "prince",
				request: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID        string
		PrimaryKey string
		client     ServiceManager
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
				client:  sv,
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
				client:  sv,
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
			name: "TestIndexBasicSearchWithIndexUIDInRequest",
			args: args{
				UID:    "indexUID",
				client: sv,
				query:  "prince",
				request: &SearchRequest{
					IndexUID: "foobar",
				},
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
				client:  customSv,
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
				client: sv,
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
				client: sv,
				request: &SearchRequest{
					Limit: 1,
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(123), "title": "Pride and Prejudice",
					},
				},
				EstimatedTotalHits: 22,
				Offset:             0,
				Limit:              1,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithOffset",
			args: args{
				UID:    "indexUID",
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client:  sv,
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
				client: sv,
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
				client: sv,
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
		{
			name: "TestIndexSearchWithRankingScoreThreshold",
			args: args{
				UID:    "indexUID",
				client: sv,
				query:  "pri",
				request: &SearchRequest{
					Limit:                 10,
					AttributesToRetrieve:  []string{"book_id", "title"},
					RankingScoreThreshold: 0.2,
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(123), "title": "Pride and Prejudice",
					},
					map[string]interface{}{
						"book_id": float64(456), "title": "Le Petit Prince",
					},
					map[string]interface{}{
						"book_id": float64(4), "title": "Harry Potter and the Half-Blood Prince",
					},
				},
				EstimatedTotalHits: 3,
				Offset:             0,
				Limit:              10,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithMatchStrategyFrequency",
			args: args{
				UID:    "indexUID",
				client: sv,
				query:  "white shirt",
				request: &SearchRequest{
					MatchingStrategy: Frequency,
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(1039), "title": "The Girl in the white shirt",
					},
				},
				EstimatedTotalHits: 1,
				Offset:             0,
				Limit:              20,
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithInvalidIndex",
			args: args{
				UID:    "invalidIndex",
				client: sv,
				query:  "pri",
				request: &SearchRequest{
					Limit:                 10,
					AttributesToRetrieve:  []string{"book_id", "title"},
					RankingScoreThreshold: 0.2,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "TestIndexSearchWithLocate",
			args: args{
				UID:    "indexUID",
				client: sv,
				query:  "王子",
				request: &SearchRequest{
					Locates: []string{"jpn"},
				},
			},
			want: &SearchResponse{
				Hits: []interface{}{
					map[string]interface{}{
						"book_id": float64(1050), "title": "星の王子さま",
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
			setUpIndexForFaceting(tt.args.client)
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
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID                  string
		PrimaryKey           string
		client               ServiceManager
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
				client:  sv,
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
				client: sv,
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
				FacetDistribution: map[string]interface{}{
					"tag": map[string]interface{}{
						"Epic fantasy": float64(1),
						"Tale":         float64(1),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithFacetsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customSv,
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
				FacetDistribution: map[string]interface{}{
					"tag": map[string]interface{}{
						"Epic fantasy": float64(1),
						"Tale":         float64(1),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithFacetsAndFacetsStats",
			args: args{
				UID:    "indexUID",
				client: sv,
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
				FacetDistribution: map[string]interface{}{
					"book_id": map[string]interface{}{
						"4":   float64(1),
						"456": float64(1),
					},
				},
				FacetStats: map[string]interface{}{
					"book_id": map[string]interface{}{
						"max": float64(456),
						"min": float64(4),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	sv := setup(t, "")

	type args struct {
		UID                  string
		PrimaryKey           string
		client               ServiceManager
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
			setUpIndexForFaceting(tt.args.client)
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
	sv := setup(t, "")

	type args struct {
		UID                string
		PrimaryKey         string
		client             ServiceManager
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				EstimatedTotalHits: 22,
				Offset:             0,
				Limit:              4,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID                 string
		PrimaryKey          string
		client              ServiceManager
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
				client:  sv,
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
				client:  customSv,
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
				client:  sv,
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
				client:  sv,
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
				client: sv,
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
			setUpIndexWithNestedFields(tt.args.client, tt.args.UID)
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
	sv := setup(t, "")

	type args struct {
		UID        string
		PrimaryKey string
		client     ServiceManager
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
			setUpIndexForFaceting(tt.args.client)
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
	sv := setup(t, "")

	type args struct {
		UID        string
		PrimaryKey string
		client     ServiceManager
		query      string
		request    SearchRequest
	}
	testArg := args{
		UID:    "indexUID",
		client: sv,
		query:  "and",
		request: SearchRequest{
			ShowRankingScore: true,
		},
	}
	setUpIndexForFaceting(testArg.client)
	c := testArg.client
	i := c.Index(testArg.UID)
	t.Cleanup(cleanup(c))

	got, err := i.Search(testArg.query, &testArg.request)
	require.NoError(t, err)
	require.NotNil(t, got.Hits[0].(map[string]interface{})["_rankingScore"])
}

func TestIndex_SearchWithShowRankingScoreDetails(t *testing.T) {
	sv := setup(t, "")

	type args struct {
		UID        string
		PrimaryKey string
		client     ServiceManager
		query      string
		request    SearchRequest
	}
	testArg := args{
		UID:    "indexUID",
		client: sv,
		query:  "and",
		request: SearchRequest{
			ShowRankingScoreDetails: true,
		},
	}
	setUpIndexForFaceting(testArg.client)
	c := testArg.client
	i := c.Index(testArg.UID)
	t.Cleanup(cleanup(c))

	got, err := i.Search(testArg.query, &testArg.request)
	require.NoError(t, err)
	require.NotNil(t, got.Hits[0].(map[string]interface{})["_rankingScoreDetails"])
}

func TestIndex_SearchWithVectorStore(t *testing.T) {
	sv := setup(t, "")

	tests := []struct {
		name       string
		UID        string
		PrimaryKey string
		client     ServiceManager
		query      string
		request    SearchRequest
	}{
		{
			name:   "basic hybrid test",
			UID:    "indexUID",
			client: sv,
			query:  "Pride and Prejudice",
			request: SearchRequest{
				Hybrid: &SearchRequestHybrid{
					SemanticRatio: 0.5,
					Embedder:      "default",
				},
				RetrieveVectors: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, err := setUpIndexWithVector(tt.client.(*meilisearch), tt.UID)
			if err != nil {
				t.Fatal(err)
			}

			c := tt.client
			t.Cleanup(cleanup(c))

			got, err := i.Search(tt.query, &tt.request)
			require.NoError(t, err)

			for _, hit := range got.Hits {
				hit := hit.(map[string]interface{})
				require.NotNil(t, hit["_vectors"])
			}
		})
	}
}

func TestIndex_SearchWithDistinct(t *testing.T) {
	sv := setup(t, "")

	tests := []struct {
		UID        string
		PrimaryKey string
		client     ServiceManager
		query      string
		request    SearchRequest
	}{
		{
			UID:    "indexUID",
			client: sv,
			query:  "white shirt",
			request: SearchRequest{
				Distinct: "sku",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.UID, func(t *testing.T) {
			setUpDistinctIndex(tt.client, tt.UID)
			c := tt.client
			t.Cleanup(cleanup(c))
			i := c.Index(tt.UID)

			got, err := i.Search(tt.query, &tt.request)
			require.NoError(t, err)
			require.NotNil(t, got.Hits)
		})
	}
}

func TestIndex_SearchSimilarDocuments(t *testing.T) {
	sv := setup(t, "")

	tests := []struct {
		UID        string
		PrimaryKey string
		client     ServiceManager
		request    *SimilarDocumentQuery
		resp       *SimilarDocumentResult
		wantErr    bool
	}{
		{
			UID:    "indexUID",
			client: sv,
			request: &SimilarDocumentQuery{
				Id: "123",
				Embedder: "default",
			},
			resp:    new(SimilarDocumentResult),
			wantErr: false,
		},
		{
			UID:     "indexUID",
			client:  sv,
			request: &SimilarDocumentQuery{
				Embedder: "default",
			},
			resp:    new(SimilarDocumentResult),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.UID, func(t *testing.T) {
			i, err := setUpIndexWithVector(tt.client.(*meilisearch), tt.UID)
			require.NoError(t, err)
			c := tt.client
			t.Cleanup(cleanup(c))

			err = i.SearchSimilarDocuments(tt.request, tt.resp)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, tt.resp)
		})
	}
}

func TestIndex_FacetSearch(t *testing.T) {
	sv := setup(t, "")

	type args struct {
		UID                  string
		PrimaryKey           string
		client               ServiceManager
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client:  sv,
				request: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "TestIndexFacetSearchWithNoFacetName",
			args: args{
				UID:    "indexUID",
				client: sv,
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
				client: sv,
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
				client: sv,
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
				client: sv,
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
			setUpIndexForFaceting(tt.args.client)
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
