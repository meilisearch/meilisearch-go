package meilisearch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type benchDoc struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Categories  []string `json:"categories"`
	Price       float64  `json:"price"`
}

func generatePayload(size int) []benchDoc {
	docs := make([]benchDoc, size)
	for i := 0; i < size; i++ {
		docs[i] = benchDoc{
			ID:          fmt.Sprintf("doc_%d", i),
			Title:       "High Performance Go Architecture",
			Description: "An in-depth look into SDK design, benchmarking, and real-world system integrations.",
			Categories:  []string{"programming", "go", "architecture"},
			Price:       49.99,
		}
	}
	return docs
}

func setupBenchServer() *httptest.Server {
	retryCount := 0
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/get-fast":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ok"}`))

		case "/post-fast":
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"taskUid": 1}`))

		case "/ndjson-large":
			w.Header().Set("Content-Type", "application/x-ndjson")
			w.WriteHeader(http.StatusOK)
			for i := 0; i < 500; i++ {
				_, _ = w.Write([]byte(`{"id":"1","title":"test"}` + "\n"))
			}

		case "/retry-sim":
			if retryCount%3 != 0 {
				retryCount++
				w.WriteHeader(http.StatusBadGateway)
				return
			}
			retryCount++
			w.WriteHeader(http.StatusOK)

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func BenchmarkClient_ExecuteRequest(b *testing.B) {
	ts := setupBenchServer()
	defer ts.Close()

	payloadSmall := generatePayload(10)
	payloadLarge := generatePayload(1000)

	baseCfg := &clientConfig{
		disableRetry:  false,
		maxRetries:    3,
		retryOnStatus: map[int]bool{502: true, 503: true, 504: true},
		jsonMarshal:   json.Marshal,
		jsonUnmarshal: json.Unmarshal,
	}

	b.Run("GET_Simple", func(b *testing.B) {
		c := newClient(&http.Client{}, ts.URL, "key", baseCfg)
		ctx := context.Background()
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = c.executeRequest(ctx, &internalRequest{
				endpoint:            "/get-fast",
				method:              http.MethodGet,
				acceptedStatusCodes: []int{http.StatusOK},
			})
		}
	})

	b.Run("POST_SmallPayload_Uncompressed", func(b *testing.B) {
		c := newClient(&http.Client{}, ts.URL, "key", baseCfg)
		ctx := context.Background()
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = c.executeRequest(ctx, &internalRequest{
				endpoint:            "/post-fast",
				method:              http.MethodPost,
				contentType:         contentTypeJSON,
				withRequest:         payloadSmall,
				acceptedStatusCodes: []int{http.StatusCreated},
			})
		}
	})

	b.Run("POST_LargePayload_Uncompressed", func(b *testing.B) {
		c := newClient(&http.Client{}, ts.URL, "key", baseCfg)
		ctx := context.Background()
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = c.executeRequest(ctx, &internalRequest{
				endpoint:            "/post-fast",
				method:              http.MethodPost,
				contentType:         contentTypeJSON,
				withRequest:         payloadLarge,
				acceptedStatusCodes: []int{http.StatusCreated},
			})
		}
	})

	b.Run("POST_LargePayload_Gzip", func(b *testing.B) {
		cfgGzip := *baseCfg
		cfgGzip.contentEncoding = GzipEncoding
		cfgGzip.encodingCompressionLevel = DefaultCompression
		c := newClient(&http.Client{}, ts.URL, "key", &cfgGzip)
		ctx := context.Background()
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = c.executeRequest(ctx, &internalRequest{
				endpoint:            "/post-fast",
				method:              http.MethodPost,
				contentType:         contentTypeJSON,
				withRequest:         payloadLarge,
				acceptedStatusCodes: []int{http.StatusCreated},
			})
		}
	})

	b.Run("GET_NDJSON_Decoding", func(b *testing.B) {
		c := newClient(&http.Client{}, ts.URL, "key", baseCfg)
		ctx := context.Background()
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			var resp []benchDoc
			_ = c.executeRequest(ctx, &internalRequest{
				endpoint:            "/ndjson-large",
				method:              http.MethodGet,
				acceptedContentType: contentTypeNDJSON,
				withResponse:        &resp,
				acceptedStatusCodes: []int{http.StatusOK},
			})
		}
	})

	b.Run("GET_WithRetries", func(b *testing.B) {
		c := newClient(&http.Client{}, ts.URL, "key", baseCfg)

		c.retryBackoff = func(attempt uint8) time.Duration {
			return 0
		}

		ctx := context.Background()
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = c.executeRequest(ctx, &internalRequest{
				endpoint:            "/retry-sim",
				method:              http.MethodGet,
				acceptedStatusCodes: []int{http.StatusOK},
			})
		}
	})
}
