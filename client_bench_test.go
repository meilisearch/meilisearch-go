package meilisearch

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Benchmark_ExecuteRequest(b *testing.B) {
	b.ReportAllocs()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/test" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message":"get successful"}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := newClient(&http.Client{}, ts.URL, "testApiKey", clientConfig{
		disableRetry: true,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := c.executeRequest(context.Background(), &internalRequest{
			endpoint:            "/test",
			method:              http.MethodGet,
			withResponse:        &mockResponse{},
			acceptedStatusCodes: []int{http.StatusOK},
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_ExecuteRequestWithEncoding(b *testing.B) {
	b.ReportAllocs()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/test" {
			accept := r.Header.Get("Accept-Encoding")
			ce := r.Header.Get("Content-Encoding")

			reqEnc := newEncoding(ContentEncoding(ce), DefaultCompression)
			respEnc := newEncoding(ContentEncoding(accept), DefaultCompression)
			req := new(mockData)

			if len(ce) != 0 {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					b.Fatal(err)
				}

				err = reqEnc.Decode(body, req)
				if err != nil {
					b.Fatal(err)
				}
			}

			if len(accept) != 0 {
				d, err := json.Marshal(req)
				if err != nil {
					b.Fatal(err)
				}
				res, err := respEnc.Encode(bytes.NewReader(d))
				if err != nil {
					b.Fatal(err)
				}
				_, _ = w.Write(res.Bytes())
				w.WriteHeader(http.StatusOK)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := newClient(&http.Client{}, ts.URL, "testApiKey", clientConfig{
		disableRetry:             true,
		contentEncoding:          GzipEncoding,
		encodingCompressionLevel: DefaultCompression,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := c.executeRequest(context.Background(), &internalRequest{
			endpoint:            "/test",
			method:              http.MethodPost,
			contentType:         contentTypeJSON,
			withRequest:         &mockData{Name: "foo", Age: 30},
			withResponse:        &mockData{},
			acceptedStatusCodes: []int{http.StatusOK},
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_ExecuteRequestWithoutRetries(b *testing.B) {
	b.ReportAllocs()
	retryCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/test" {
			if retryCount == 2 {
				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusBadGateway)
			retryCount++
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := newClient(&http.Client{}, ts.URL, "testApiKey", clientConfig{
		disableRetry: false,
		maxRetries:   3,
		retryOnStatus: map[int]bool{
			502: true,
			503: true,
			504: true,
		},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := c.executeRequest(context.Background(), &internalRequest{
			endpoint:            "/test",
			method:              http.MethodGet,
			withResponse:        nil,
			withRequest:         nil,
			acceptedStatusCodes: []int{http.StatusOK},
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}
