package integration

import (
	"crypto/tls"
	"encoding/json"
	"github.com/meilisearch/meilisearch-go"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndex_SearchWithContentEncoding(t *testing.T) {
	tests := []struct {
		Name            string
		ContentEncoding meilisearch.ContentEncoding
		Query           string
		Request         *meilisearch.SearchRequest
		FacetRequest    *meilisearch.FacetSearchRequest
		Response        *meilisearch.SearchResponse
		FacetResponse   *meilisearch.FacetSearchResponse
	}{
		{
			Name:            "SearchResultWithGzipEncoding",
			ContentEncoding: meilisearch.GzipEncoding,
			Query:           "prince",
			Request: &meilisearch.SearchRequest{
				IndexUID: "indexUID",
				Limit:    20,
				Offset:   0,
			},
			FacetRequest: &meilisearch.FacetSearchRequest{
				FacetName:  "tag",
				FacetQuery: "Novel",
			},
			FacetResponse: &meilisearch.FacetSearchResponse{
				FacetHits: meilisearch.Hits{
					{"value": json.RawMessage(`"Novel"`), "count": json.RawMessage(`5`)},
				},
				FacetQuery: "Novel",
			},
			Response: &meilisearch.SearchResponse{
				Hits: meilisearch.Hits{
					{"book_id": json.RawMessage(`456`), "title": json.RawMessage(`"Le Petit Prince"`)},
					{"Tag": json.RawMessage(`"Epic fantasy"`), "book_id": json.RawMessage(`4`), "title": json.RawMessage(`"Harry Potter and the Half-Blood Prince"`)},
				},
				EstimatedTotalHits: 2,
				Offset:             0,
				Limit:              20,
			},
		},
		{
			Name:            "SearchResultWithDeflateEncoding",
			ContentEncoding: meilisearch.DeflateEncoding,
			Query:           "prince",
			Request: &meilisearch.SearchRequest{
				IndexUID: "indexUID",
				Limit:    20,
				Offset:   0,
			},
			Response: &meilisearch.SearchResponse{
				Hits: meilisearch.Hits{
					{"book_id": json.RawMessage(`456`), "title": json.RawMessage(`"Le Petit Prince"`)},
					{"Tag": json.RawMessage(`"Epic fantasy"`), "book_id": json.RawMessage(`4`), "title": json.RawMessage(`"Harry Potter and the Half-Blood Prince"`)},
				},
				EstimatedTotalHits: 2,
				Offset:             0,
				Limit:              20,
			},
			FacetRequest: &meilisearch.FacetSearchRequest{
				FacetName:  "tag",
				FacetQuery: "Novel",
			},
			FacetResponse: &meilisearch.FacetSearchResponse{
				FacetHits: meilisearch.Hits{
					{"value": json.RawMessage(`"Novel"`), "count": json.RawMessage(`5`)},
				},
				FacetQuery: "Novel",
			},
		},
		{
			Name:            "SearchResultWithBrotliEncoding",
			ContentEncoding: meilisearch.BrotliEncoding,
			Query:           "prince",
			Request: &meilisearch.SearchRequest{
				IndexUID: "indexUID",
				Limit:    20,
				Offset:   0,
			},
			Response: &meilisearch.SearchResponse{
				Hits: meilisearch.Hits{
					{"book_id": json.RawMessage(`456`), "title": json.RawMessage(`"Le Petit Prince"`)},
					{"Tag": json.RawMessage(`"Epic fantasy"`), "book_id": json.RawMessage(`4`), "title": json.RawMessage(`"Harry Potter and the Half-Blood Prince"`)},
				},
				EstimatedTotalHits: 2,
				Offset:             0,
				Limit:              20,
			},
			FacetRequest: &meilisearch.FacetSearchRequest{
				FacetName:  "tag",
				FacetQuery: "Novel",
			},
			FacetResponse: &meilisearch.FacetSearchResponse{
				FacetHits: meilisearch.Hits{
					{"value": json.RawMessage(`"Novel"`), "count": json.RawMessage(`5`)},
				},
				FacetQuery: "Novel",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			sv := setup(t, "", meilisearch.WithContentEncoding(tt.ContentEncoding, meilisearch.DefaultCompression))
			setUpIndexForFaceting(sv)
			i := sv.Index(tt.Request.IndexUID)
			t.Cleanup(cleanup(sv))

			require.NotNil(t, tt.Request)
			require.NotEmpty(t, tt.Request.IndexUID)
			got, err := i.Search(tt.Query, tt.Request)
			require.NoError(t, err, "Search request failed unexpectedly")
			require.Equal(t, len(tt.Response.Hits), len(got.Hits))

			gotJson, err := i.SearchRaw(tt.Query, tt.Request)
			require.NoError(t, err)

			var resp meilisearch.SearchResponse
			err = json.Unmarshal(*gotJson, &resp)
			require.NoError(t, err, "error unmarshalling raw got meilisearch.SearchResponse")
			require.Equal(t, len(tt.Response.Hits), len(resp.Hits))

			filterableAttrs := []string{"tag"}
			task, err := i.UpdateFilterableAttributes(&filterableAttrs)
			require.NoError(t, err)
			testWaitForTask(t, i, task)

			gotJson, err = i.FacetSearch(tt.FacetRequest)
			require.NoError(t, err)
			var gotFacet meilisearch.FacetSearchResponse
			err = json.Unmarshal(*gotJson, &gotFacet)
			require.NoError(t, err, "error unmarshalling raw got meilisearch.FacetSearchResponse")
			require.Equal(t, len(gotFacet.FacetHits), len(tt.FacetResponse.FacetHits))
		})
	}
}

func TestIndex_SearchRaw(t *testing.T) {
	sv := setup(t, "")

	type args struct {
		UID        string
		PrimaryKey string
		client     meilisearch.ServiceManager
		query      string
		request    *meilisearch.SearchRequest
	}

	tests := []struct {
		name    string
		args    args
		want    *meilisearch.SearchResponse
		wantErr bool
	}{
		{
			name: "TestIndexBasicSearch",
			args: args{
				UID:    "indexUID",
				client: sv,
				query:  "prince",
				request: &meilisearch.SearchRequest{
					IndexUID: "foobar",
				},
			},
			want: &meilisearch.SearchResponse{
				Hits: meilisearch.Hits{
					{"book_id": json.RawMessage(`456`), "title": json.RawMessage(`"Le Petit Prince"`)},
					{"Tag": json.RawMessage(`"Epic fantasy"`), "book_id": json.RawMessage(`4`), "title": json.RawMessage(`"Harry Potter and the Half-Blood Prince"`)},
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

			// Unmarshal the raw response from SearchRaw into a meilisearch.SearchResponse
			var got meilisearch.SearchResponse
			err = json.Unmarshal(*gotRaw, &got)
			require.NoError(t, err, "error unmarshalling raw got meilisearch.SearchResponse")

			// Check meilisearch.Hits length
			require.Equal(t, len(tt.want.Hits), len(got.Hits))

			// Compare each hit in meilisearch.Hits
			for idx := range got.Hits {
				expectedHit := tt.want.Hits[idx]
				actualHit := got.Hits[idx]

				require.Equal(t, expectedHit["title"], actualHit["title"])
				require.Equal(t, expectedHit["book_id"], actualHit["book_id"])
			}

			// Check if `_formatted` exists before comparison
			if _, ok := tt.want.Hits[0]["_formatted"]; ok {
				require.Equal(t, tt.want.Hits[0]["_formatted"], got.Hits[0]["_formatted"])
			}

			// Check other response fields
			require.Equal(t, tt.want.EstimatedTotalHits, got.EstimatedTotalHits)
			require.Equal(t, tt.want.Offset, got.Offset)
			require.Equal(t, tt.want.Limit, got.Limit)
			require.Equal(t, tt.want.FacetDistribution, got.FacetDistribution)
		})
	}
}

func TestIndex_Search(t *testing.T) {
	sv := setup(t, "")

	type args struct {
		UID        string
		PrimaryKey string
		client     meilisearch.ServiceManager
		query      string
		request    *meilisearch.SearchRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *meilisearch.SearchResponse
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
				request: &meilisearch.SearchRequest{},
			},
			want: &meilisearch.SearchResponse{
				Hits: meilisearch.Hits{
					{"book_id": json.RawMessage(`456`), "title": json.RawMessage(`"Le Petit Prince"`)},
					{"Tag": json.RawMessage(`"Epic fantasy"`), "book_id": json.RawMessage(`4`), "title": json.RawMessage(`"Harry Potter and the Half-Blood Prince"`)},
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

			for idx := range got.Hits {
				expectedHit := tt.want.Hits[idx]
				actualHit := got.Hits[idx]

				require.Equal(t, expectedHit["title"], actualHit["title"])
				require.Equal(t, expectedHit["book_id"], actualHit["book_id"])
			}

			if _, ok := tt.want.Hits[0]["_formatted"]; ok {
				require.Equal(t, tt.want.Hits[0]["_formatted"], got.Hits[0]["_formatted"])
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
	customSv := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID                  string
		PrimaryKey           string
		client               meilisearch.ServiceManager
		query                string
		request              *meilisearch.SearchRequest
		filterableAttributes []string
	}
	tests := []struct {
		name    string
		args    args
		want    *meilisearch.SearchResponse
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
				request: &meilisearch.SearchRequest{
					Facets: []string{"*"},
				},
				filterableAttributes: []string{"tag"},
			},
			want: &meilisearch.SearchResponse{
				Hits: meilisearch.Hits{
					{"book_id": json.RawMessage(`456`), "title": json.RawMessage(`"Le Petit Prince"`)},
					{"book_id": json.RawMessage(`4`), "title": json.RawMessage(`"Harry Potter and the Half-Blood Prince"`)},
				},
				EstimatedTotalHits: 2,
				Offset:             0,
				Limit:              20,
				FacetDistribution: toRawMessage(map[string]map[string]float64{
					"tag": {
						"Epic fantasy": 1,
						"Tale":         1,
					},
				}),
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithFacetsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customSv,
				query:  "prince",
				request: &meilisearch.SearchRequest{
					Facets: []string{"*"},
				},
				filterableAttributes: []string{"tag"},
			},
			want: &meilisearch.SearchResponse{
				Hits: meilisearch.Hits{
					{"book_id": json.RawMessage(`456`), "title": json.RawMessage(`"Le Petit Prince"`)},
					{"book_id": json.RawMessage(`4`), "title": json.RawMessage(`"Harry Potter and the Half-Blood Prince"`)},
				},
				EstimatedTotalHits: 2,
				Offset:             0,
				Limit:              20,
				FacetDistribution: toRawMessage(map[string]map[string]float64{
					"tag": {
						"Epic fantasy": 1,
						"Tale":         1,
					},
				}),
			},
			wantErr: false,
		},
		{
			name: "TestIndexSearchWithFacetsAndFacetsStats",
			args: args{
				UID:    "indexUID",
				client: sv,
				query:  "prince",
				request: &meilisearch.SearchRequest{
					Facets: []string{"book_id"},
				},
				filterableAttributes: []string{"book_id"},
			},
			want: &meilisearch.SearchResponse{
				Hits: meilisearch.Hits{
					{"book_id": json.RawMessage(`456`), "title": json.RawMessage(`"Le Petit Prince"`)},
					{"book_id": json.RawMessage(`4`), "title": json.RawMessage(`"Harry Potter and the Half-Blood Prince"`)},
				},
				EstimatedTotalHits: 2,
				Offset:             0,
				Limit:              20,
				FacetDistribution: toRawMessage(map[string]map[string]float64{
					"book_id": {
						"4":   1,
						"456": 1,
					},
				}),
				FacetStats: toRawMessage(map[string]map[string]float64{
					"book_id": {
						"max": 456,
						"min": 4,
					},
				}),
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

			for idx := range got.Hits {
				expectedHit := tt.want.Hits[idx]
				actualHit := got.Hits[idx]

				require.Equal(t, expectedHit["title"], actualHit["title"])
				require.Equal(t, expectedHit["book_id"], actualHit["book_id"])
			}

			require.Equal(t, tt.want.EstimatedTotalHits, got.EstimatedTotalHits)
			require.Equal(t, tt.want.Offset, got.Offset)
			require.Equal(t, tt.want.Limit, got.Limit)
			require.Equal(t, tt.want.FacetDistribution, got.FacetDistribution)

			if tt.want.FacetStats != nil {
				require.NotNil(t, got.FacetStats)
			}
		})
	}
}

func TestIndex_SearchWithFilters(t *testing.T) {
	sv := setup(t, "")

	type args struct {
		UID                  string
		PrimaryKey           string
		client               meilisearch.ServiceManager
		query                string
		filterableAttributes []string
		request              *meilisearch.SearchRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *meilisearch.SearchResponse
		wantErr bool
	}{
		{
			name: "TestIndexBasicSearchWithFilter",
			args: args{
				UID:                  "indexUID",
				client:               sv,
				query:                "and",
				filterableAttributes: []string{"tag"},
				request: &meilisearch.SearchRequest{
					Filter: "tag = romance",
				},
			},
			want: &meilisearch.SearchResponse{
				Hits: meilisearch.Hits{
					{"book_id": json.RawMessage(`123`), "title": json.RawMessage(`"Pride and Prejudice"`)},
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
				UID:                  "indexUID",
				client:               sv,
				query:                "and",
				filterableAttributes: []string{"year"},
				request: &meilisearch.SearchRequest{
					Filter: "year = 2005",
				},
			},
			want: &meilisearch.SearchResponse{
				Hits: meilisearch.Hits{
					{"book_id": json.RawMessage(`4`), "title": json.RawMessage(`"Harry Potter and the Half-Blood Prince"`)},
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

			for idx := range got.Hits {
				expectedHit := tt.want.Hits[idx]
				actualHit := got.Hits[idx]

				require.Equal(t, expectedHit["title"], actualHit["title"])
				require.Equal(t, expectedHit["book_id"], actualHit["book_id"])
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
		client             meilisearch.ServiceManager
		query              string
		sortableAttributes []string
		request            *meilisearch.SearchRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *meilisearch.SearchResponse
		wantErr bool
	}{
		{
			name: "TestIndexBasicSearchWithSortIntParameter",
			args: args{
				UID:                "indexUID",
				client:             sv,
				query:              "and",
				sortableAttributes: []string{"year"},
				request: &meilisearch.SearchRequest{
					Sort: []string{"year:asc"},
				},
			},
			want: &meilisearch.SearchResponse{
				Hits: meilisearch.Hits{
					{"book_id": toRawMessage(123), "title": toRawMessage("Pride and Prejudice")},
					{"book_id": toRawMessage(730), "title": toRawMessage("War and Peace")},
					{"book_id": toRawMessage(1032), "title": toRawMessage("Crime and Punishment")},
					{"book_id": toRawMessage(4), "title": toRawMessage("Harry Potter and the Half-Blood Prince")},
				},
				EstimatedTotalHits: 4,
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

			for idx := range got.Hits {
				expectedHit := tt.want.Hits[idx]
				actualHit := got.Hits[idx]

				require.Equal(t, expectedHit["title"], actualHit["title"])
				require.Equal(t, expectedHit["book_id"], actualHit["book_id"])
			}

			require.Equal(t, tt.want.EstimatedTotalHits, got.EstimatedTotalHits)
			require.Equal(t, tt.want.Offset, got.Offset)
			require.Equal(t, tt.want.Limit, got.Limit)
			require.Equal(t, tt.want.FacetDistribution, got.FacetDistribution)
		})
	}
}

func TestIndex_SearchOnNestedFields(t *testing.T) {
	sv := setup(t, "")

	type args struct {
		UID                 string
		PrimaryKey          string
		client              meilisearch.ServiceManager
		query               string
		request             *meilisearch.SearchRequest
		searchableAttribute []string
		sortableAttribute   []string
	}
	tests := []struct {
		name    string
		args    args
		want    *meilisearch.SearchResponse
		wantErr bool
	}{
		{
			name: "TestIndexBasicSearchOnNestedFields",
			args: args{
				UID:     "TestIndexBasicSearchOnNestedFields",
				client:  sv,
				query:   "An awesome",
				request: &meilisearch.SearchRequest{},
			},
			want: &meilisearch.SearchResponse{
				Hits: meilisearch.Hits{
					{
						"id": toRawMessage(5), "title": toRawMessage("The Hobbit"),
						"info": toRawMessage(map[string]interface{}{
							"comment": "An awesome book", "reviewNb": 900,
						}),
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

			for idx := range got.Hits {
				expectedHit := tt.want.Hits[idx]
				actualHit := got.Hits[idx]

				require.Equal(t, expectedHit["title"], actualHit["title"])
				require.Equal(t, expectedHit["id"], actualHit["id"])
				require.Equal(t, expectedHit["info"], actualHit["info"])
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
		client     meilisearch.ServiceManager
		query      string
		request    *meilisearch.SearchRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *meilisearch.SearchResponse
		wantErr bool
	}{
		{
			name: "TestIndexBasicSearchWithHitsPerPage",
			args: args{
				UID:    "indexUID",
				client: sv,
				query:  "and",
				request: &meilisearch.SearchRequest{
					HitsPerPage: 10,
				},
			},
			want: &meilisearch.SearchResponse{
				Hits: meilisearch.Hits{
					{"book_id": toRawMessage(123), "title": toRawMessage("Pride and Prejudice")},
					{"book_id": toRawMessage(730), "title": toRawMessage("War and Peace")},
					{"book_id": toRawMessage(1032), "title": toRawMessage("Crime and Punishment")},
					{"book_id": toRawMessage(4), "title": toRawMessage("Harry Potter and the Half-Blood Prince")},
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
				request: &meilisearch.SearchRequest{
					Page: 1,
				},
			},
			want: &meilisearch.SearchResponse{
				Hits: meilisearch.Hits{
					{"book_id": toRawMessage(123), "title": toRawMessage("Pride and Prejudice")},
					{"book_id": toRawMessage(730), "title": toRawMessage("War and Peace")},
					{"book_id": toRawMessage(1032), "title": toRawMessage("Crime and Punishment")},
					{"book_id": toRawMessage(4), "title": toRawMessage("Harry Potter and the Half-Blood Prince")},
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
				request: &meilisearch.SearchRequest{
					HitsPerPage: 10,
					Page:        1,
				},
			},
			want: &meilisearch.SearchResponse{
				Hits: meilisearch.Hits{
					{"book_id": toRawMessage(123), "title": toRawMessage("Pride and Prejudice")},
					{"book_id": toRawMessage(730), "title": toRawMessage("War and Peace")},
					{"book_id": toRawMessage(1032), "title": toRawMessage("Crime and Punishment")},
					{"book_id": toRawMessage(4), "title": toRawMessage("Harry Potter and the Half-Blood Prince")},
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
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, tt.want)
				return
			}

			require.NoError(t, err)
			require.Equal(t, len(tt.want.Hits), len(got.Hits))

			for idx := range got.Hits {
				expectedHit := tt.want.Hits[idx]
				actualHit := got.Hits[idx]

				require.Equal(t, expectedHit["title"], actualHit["title"])
				require.Equal(t, expectedHit["book_id"], actualHit["book_id"])
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
		client     meilisearch.ServiceManager
		query      string
		request    meilisearch.SearchRequest
	}

	testArg := args{
		UID:    "indexUID",
		client: sv,
		query:  "and",
		request: meilisearch.SearchRequest{
			ShowRankingScore: true,
		},
	}

	setUpIndexForFaceting(testArg.client)
	c := testArg.client
	i := c.Index(testArg.UID)
	t.Cleanup(cleanup(c))

	got, err := i.Search(testArg.query, &testArg.request)
	require.NoError(t, err)
	require.NotNil(t, got)

	// Ensure at least one result is present
	require.Greater(t, len(got.Hits), 0)

	// Convert to a structured format and verify _rankingScore presence
	var result map[string]json.RawMessage
	err = got.Hits[0].Decode(&result)
	require.NoError(t, err)
	require.Contains(t, got.Hits[0], "_rankingScore")
}

func TestIndex_SearchWithShowRankingScoreDetails(t *testing.T) {
	sv := setup(t, "")

	type args struct {
		UID        string
		PrimaryKey string
		client     meilisearch.ServiceManager
		query      string
		request    meilisearch.SearchRequest
	}

	testArg := args{
		UID:    "indexUID",
		client: sv,
		query:  "and",
		request: meilisearch.SearchRequest{
			ShowRankingScoreDetails: true,
		},
	}

	setUpIndexForFaceting(testArg.client)
	c := testArg.client
	i := c.Index(testArg.UID)
	t.Cleanup(cleanup(c))

	got, err := i.Search(testArg.query, &testArg.request)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Greater(t, len(got.Hits), 0, "expected at least one hit")

	// Convert first hit to structured format and check for _rankingScoreDetails
	var result map[string]json.RawMessage
	err = got.Hits[0].Decode(&result)
	require.NoError(t, err)
	require.Contains(t, result, "_rankingScoreDetails", "expected _rankingScoreDetails to be present in the search result")
}

func TestIndex_SearchWithVectorStore(t *testing.T) {
	sv := setup(t, "")

	tests := []struct {
		name       string
		UID        string
		PrimaryKey string
		client     meilisearch.ServiceManager
		query      string
		request    meilisearch.SearchRequest
	}{
		{
			name:   "basic hybrid test",
			UID:    "indexUID",
			client: sv,
			query:  "Pride and Prejudice",
			request: meilisearch.SearchRequest{
				Hybrid: &meilisearch.SearchRequestHybrid{
					SemanticRatio: 0.5,
					Embedder:      "default",
				},
				RetrieveVectors: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, err := setUpIndexWithVector(tt.client, tt.UID)
			require.NoError(t, err)

			c := tt.client
			t.Cleanup(cleanup(c))

			got, err := i.Search(tt.query, &tt.request)
			require.NoError(t, err)
			require.Greater(t, len(got.Hits), 0, "expected at least one hit")

			for _, hit := range got.Hits {
				var hitMap map[string]json.RawMessage
				err := hit.Decode(&hitMap)
				require.NoError(t, err)
				require.Contains(t, hitMap, "_vectors", "expected _vectors field in search result")
			}
		})
	}
}

func TestIndex_SearchWithDistinct(t *testing.T) {
	sv := setup(t, "")

	tests := []struct {
		UID        string
		PrimaryKey string
		client     meilisearch.ServiceManager
		query      string
		request    meilisearch.SearchRequest
	}{
		{
			UID:    "indexUID",
			client: sv,
			query:  "white shirt",
			request: meilisearch.SearchRequest{
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
		client     meilisearch.ServiceManager
		request    *meilisearch.SimilarDocumentQuery
		resp       *meilisearch.SimilarDocumentResult
		wantErr    bool
	}{
		{
			UID:    "indexUID",
			client: sv,
			request: &meilisearch.SimilarDocumentQuery{
				Id:       "123",
				Embedder: "default",
			},
			resp:    new(meilisearch.SimilarDocumentResult),
			wantErr: false,
		},
		{
			UID:    "indexUID",
			client: sv,
			request: &meilisearch.SimilarDocumentQuery{
				Embedder: "default",
			},
			resp:    new(meilisearch.SimilarDocumentResult),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.UID, func(t *testing.T) {
			i, err := setUpIndexWithVector(tt.client, tt.UID)
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
		client               meilisearch.ServiceManager
		request              *meilisearch.FacetSearchRequest
		filterableAttributes []string
	}

	tests := []struct {
		name    string
		args    args
		want    *meilisearch.FacetSearchResponse
		wantErr bool
	}{
		{
			name: "TestIndexBasicFacetSearch",
			args: args{
				UID:    "indexUID",
				client: sv,
				request: &meilisearch.FacetSearchRequest{
					FacetName:  "tag",
					FacetQuery: "Novel",
				},
				filterableAttributes: []string{"tag"},
			},
			want: &meilisearch.FacetSearchResponse{
				FacetHits: meilisearch.Hits{
					{"value": toRawMessage("Novel"), "count": toRawMessage(5)},
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
				request: &meilisearch.FacetSearchRequest{
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
				request: &meilisearch.FacetSearchRequest{
					Q:         "query",
					FacetName: "tag",
				},
				filterableAttributes: []string{"tag"},
			},
			want: &meilisearch.FacetSearchResponse{
				FacetHits:  meilisearch.Hits{},
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
			require.NotNil(t, gotRaw)

			// Unmarshal the raw response into a meilisearch.FacetSearchResponse
			var got meilisearch.FacetSearchResponse
			err = json.Unmarshal(*gotRaw, &got)
			require.NoError(t, err, "error unmarshalling raw got meilisearch.FacetSearchResponse")

			require.Equal(t, len(tt.want.FacetHits), len(got.FacetHits))

			for idx := range got.FacetHits {
				expectedHit := tt.want.FacetHits[idx]
				actualHit := got.FacetHits[idx]

				require.Equal(t, expectedHit["value"], actualHit["value"])
				require.Equal(t, expectedHit["count"], actualHit["count"])
			}

			require.Equal(t, tt.want.FacetQuery, got.FacetQuery)
		})
	}
}
