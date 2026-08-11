package meilisearch

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStatsIndexSizeFields(t *testing.T) {
	tests := []struct {
		name          string
		sizeFormat    string
		response      string
		indexSize     any
		usedIndexSize any
	}{
		{
			name:          "raw",
			sizeFormat:    "raw",
			response:      `{"numberOfDocuments": 6, "indexSize": 4096, "usedIndexSize": 2048}`,
			indexSize:     float64(4096),
			usedIndexSize: float64(2048),
		},
		{
			name:          "human",
			sizeFormat:    "human",
			response:      `{"numberOfDocuments": 6, "indexSize": "4 KiB", "usedIndexSize": "2 KiB"}`,
			indexSize:     "4 KiB",
			usedIndexSize: "2 KiB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if sizeFormat := r.URL.Query().Get("sizeFormat"); sizeFormat != tt.sizeFormat {
					t.Errorf("sizeFormat = %q, want %q", sizeFormat, tt.sizeFormat)
				}
				w.Header().Set("Content-Type", "application/json")

				switch r.URL.Path {
				case "/indexes/movies/stats":
					_, _ = w.Write([]byte(tt.response))
				case "/stats":
					_, _ = fmt.Fprintf(w, `{"indexes":{"movies":%s}}`, tt.response)
				default:
					http.NotFound(w, r)
				}
			}))
			t.Cleanup(server.Close)

			client := New(server.URL)
			params := &StatsParams{SizeFormat: tt.sizeFormat}

			indexStats, err := client.Index("movies").GetStats(params)
			require.NoError(t, err)
			require.Equal(t, tt.indexSize, indexStats.IndexSize)
			require.Equal(t, tt.usedIndexSize, indexStats.UsedIndexSize)

			allStats, err := client.GetStats(params)
			require.NoError(t, err)
			require.Equal(t, tt.indexSize, allStats.Indexes["movies"].IndexSize)
			require.Equal(t, tt.usedIndexSize, allStats.Indexes["movies"].UsedIndexSize)
		})
	}
}
