package meilisearch

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchRequest_ShowPerformanceDetails(t *testing.T) {
	tests := []struct {
		name     string
		request  SearchRequest
		expected string
	}{
		{
			name: "ShowPerformanceDetails is true",
			request: SearchRequest{
				Query:                  "test",
				ShowPerformanceDetails: true,
			},
			expected: `{"showPerformanceDetails":true,"q":"test","hybrid":null}`,
		},
		{
			name: "ShowPerformanceDetails is false (omitted)",
			request: SearchRequest{
				Query: "test",
			},
			expected: `{"q":"test","hybrid":null}`,
		},
		{
			name: "ShowPerformanceDetails with other show options",
			request: SearchRequest{
				Query:                   "test",
				ShowPerformanceDetails:  true,
				ShowRankingScore:        true,
				ShowRankingScoreDetails: true,
			},
			expected: `{"showRankingScore":true,"showRankingScoreDetails":true,"showPerformanceDetails":true,"q":"test","hybrid":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.request)
			require.NoError(t, err)
			require.JSONEq(t, tt.expected, string(data))
		})
	}
}

func TestSearchResponse_PerformanceDetails(t *testing.T) {
	tests := []struct {
		name           string
		jsonResponse   string
		expectedHasKey bool
	}{
		{
			name: "Response with performanceDetails",
			jsonResponse: `{
				"hits": [],
				"processingTimeMs": 10,
				"query": "test",
				"performanceDetails": {"step1": {"time": 5}, "step2": {"time": 3}}
			}`,
			expectedHasKey: true,
		},
		{
			name: "Response without performanceDetails",
			jsonResponse: `{
				"hits": [],
				"processingTimeMs": 10,
				"query": "test"
			}`,
			expectedHasKey: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp SearchResponse
			err := json.Unmarshal([]byte(tt.jsonResponse), &resp)
			require.NoError(t, err)

			if tt.expectedHasKey {
				require.NotNil(t, resp.PerformanceDetails)
				require.NotEmpty(t, resp.PerformanceDetails)
			} else {
				require.Nil(t, resp.PerformanceDetails)
			}
		})
	}
}

func TestSimilarDocumentQuery_ShowPerformanceDetails(t *testing.T) {
	tests := []struct {
		name     string
		query    SimilarDocumentQuery
		expected string
	}{
		{
			name: "ShowPerformanceDetails is true",
			query: SimilarDocumentQuery{
				Id:                     "1",
				Embedder:               "default",
				ShowPerformanceDetails: true,
			},
			expected: `{"id":"1","embedder":"default","showPerformanceDetails":true}`,
		},
		{
			name: "ShowPerformanceDetails is false (omitted)",
			query: SimilarDocumentQuery{
				Id:       "1",
				Embedder: "default",
			},
			expected: `{"id":"1","embedder":"default"}`,
		},
		{
			name: "ShowPerformanceDetails with other show options",
			query: SimilarDocumentQuery{
				Id:                      "1",
				Embedder:                "default",
				ShowPerformanceDetails:  true,
				ShowRankingScore:        true,
				ShowRankingScoreDetails: true,
			},
			expected: `{"id":"1","embedder":"default","showRankingScore":true,"showRankingScoreDetails":true,"showPerformanceDetails":true}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.query)
			require.NoError(t, err)
			require.JSONEq(t, tt.expected, string(data))
		})
	}
}
