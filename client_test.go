package meilisearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Mock structures for testing
type mockResponse struct {
	Message string `json:"message"`
}

type mockJsonMarshaller struct {
	valid bool
	null  bool
	Foo   string `json:"foo"`
	Bar   string `json:"bar"`
}

func TestExecuteRequest(t *testing.T) {
	retryCount := 0

	// Create a mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/test-get" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message":"get successful"}`))
		} else if r.Method == http.MethodGet && r.URL.Path == "/test-get-encoding" {
			encode := r.Header.Get("Accept-Encoding")
			if len(encode) != 0 {
				enc := newEncoding(ContentEncoding(encode), DefaultCompression)
				d := &mockData{Name: "foo", Age: 30}

				b, err := json.Marshal(d)
				require.NoError(t, err)

				res, err := enc.Encode(bytes.NewReader(b))
				require.NoError(t, err)
				_, _ = w.Write(res.Bytes())
				w.WriteHeader(http.StatusOK)
				return
			}
			_, _ = w.Write([]byte("invalid message"))
			w.WriteHeader(http.StatusInternalServerError)
		} else if r.Method == http.MethodPost && r.URL.Path == "/test-req-resp-encoding" {
			accept := r.Header.Get("Accept-Encoding")
			ce := r.Header.Get("Content-Encoding")

			reqEnc := newEncoding(ContentEncoding(ce), DefaultCompression)
			respEnc := newEncoding(ContentEncoding(accept), DefaultCompression)
			req := new(mockData)

			if len(ce) != 0 {
				b, err := io.ReadAll(r.Body)
				require.NoError(t, err)

				err = reqEnc.Decode(b, req)
				require.NoError(t, err)
			}

			if len(accept) != 0 {
				d, err := json.Marshal(req)
				require.NoError(t, err)
				res, err := respEnc.Encode(bytes.NewReader(d))
				require.NoError(t, err)
				_, _ = w.Write(res.Bytes())
				w.WriteHeader(http.StatusOK)
			}
		} else if r.Method == http.MethodPost && r.URL.Path == "/test-post" {
			w.WriteHeader(http.StatusCreated)
			msg := []byte(`{"message":"post successful"}`)
			_, _ = w.Write(msg)

		} else if r.Method == http.MethodGet && r.URL.Path == "/test-null-body" {
			w.WriteHeader(http.StatusOK)
			msg := []byte(`null`)
			_, _ = w.Write(msg)
		} else if r.Method == http.MethodPost && r.URL.Path == "/test-post-encoding" {
			w.WriteHeader(http.StatusCreated)
			msg := []byte(`{"message":"post successful"}`)

			enc := r.Header.Get("Accept-Encoding")
			if len(enc) != 0 {
				e := newEncoding(ContentEncoding(enc), DefaultCompression)
				b, err := e.Encode(bytes.NewReader(msg))
				require.NoError(t, err)
				_, _ = w.Write(b.Bytes())
				return
			}
			_, _ = w.Write(msg)
		} else if r.URL.Path == "/test-bad-request" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message":"bad request"}`))
		} else if r.URL.Path == "/invalid-response-body" {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message":"bad response body"}`))
		} else if r.URL.Path == "/io-reader" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message":"io reader"}`))
		} else if r.URL.Path == "/failed-retry" {
			w.WriteHeader(http.StatusBadGateway)
		} else if r.URL.Path == "/success-retry" {
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

	tests := []struct {
		name            string
		internalReq     *internalRequest
		expectedResp    interface{}
		contentEncoding ContentEncoding
		withTimeout     bool
		disableRetry    bool
		wantErr         bool
	}{
		{
			name: "Successful GET request",
			internalReq: &internalRequest{
				endpoint:            "/test-get",
				method:              http.MethodGet,
				withResponse:        &mockResponse{},
				acceptedStatusCodes: []int{http.StatusOK},
			},
			expectedResp: &mockResponse{Message: "get successful"},
			wantErr:      false,
		},
		{
			name: "Successful POST request",
			internalReq: &internalRequest{
				endpoint:            "/test-post",
				method:              http.MethodPost,
				withRequest:         map[string]string{"key": "value"},
				contentType:         contentTypeJSON,
				withResponse:        &mockResponse{},
				acceptedStatusCodes: []int{http.StatusCreated},
			},
			expectedResp: &mockResponse{Message: "post successful"},
			wantErr:      false,
		},
		{
			name: "404 Not Found",
			internalReq: &internalRequest{
				endpoint:            "/not-found",
				method:              http.MethodGet,
				withResponse:        &mockResponse{},
				acceptedStatusCodes: []int{http.StatusOK},
			},
			expectedResp: nil,
			wantErr:      true,
		},
		{
			name: "Invalid URL",
			internalReq: &internalRequest{
				endpoint:            "/invalid-url$%^*()*#",
				method:              http.MethodGet,
				withResponse:        &mockResponse{},
				acceptedStatusCodes: []int{http.StatusOK},
			},
			expectedResp: nil,
			wantErr:      true,
		},
		{
			name: "Invalid response body",
			internalReq: &internalRequest{
				endpoint:            "/invalid-response-body",
				method:              http.MethodGet,
				withResponse:        struct{}{},
				acceptedStatusCodes: []int{http.StatusInternalServerError},
			},
			expectedResp: nil,
			wantErr:      true,
		},
		{
			name: "Invalid request method",
			internalReq: &internalRequest{
				endpoint:            "/invalid-request-method",
				method:              http.MethodGet,
				withResponse:        nil,
				withRequest:         struct{}{},
				acceptedStatusCodes: []int{http.StatusBadRequest},
			},
			expectedResp: nil,
			wantErr:      true,
		},
		{
			name: "Invalid request content type",
			internalReq: &internalRequest{
				endpoint:            "/invalid-request-content-type",
				method:              http.MethodPost,
				withResponse:        nil,
				contentType:         "",
				withRequest:         struct{}{},
				acceptedStatusCodes: []int{http.StatusBadRequest},
			},
			expectedResp: nil,
			wantErr:      true,
		},
		{
			name: "Invalid json marshaler",
			internalReq: &internalRequest{
				endpoint:     "/invalid-marshaler",
				method:       http.MethodPost,
				withResponse: nil,
				withRequest: &mockJsonMarshaller{
					valid: false,
				},
				contentType: "application/json",
			},
			expectedResp: nil,
			wantErr:      true,
		},
		{
			name: "Null data marshaler",
			internalReq: &internalRequest{
				endpoint:     "/null-data-marshaler",
				method:       http.MethodPost,
				withResponse: nil,
				withRequest: &mockJsonMarshaller{
					valid: true,
					null:  true,
				},
				contentType: "application/json",
			},
			expectedResp: nil,
			wantErr:      true,
		},
		{
			name: "Successful request with io.reader",
			internalReq: &internalRequest{
				endpoint:            "/io-reader",
				method:              http.MethodPost,
				withResponse:        nil,
				contentType:         "text/plain",
				withRequest:         strings.NewReader("foobar"),
				acceptedStatusCodes: []int{http.StatusOK},
			},
			expectedResp: nil,
			wantErr:      false,
		},
		{
			name: "Test null body response",
			internalReq: &internalRequest{
				endpoint:            "/test-null-body",
				method:              http.MethodGet,
				withResponse:        make([]byte, 0),
				contentType:         "application/json",
				acceptedStatusCodes: []int{http.StatusOK},
			},
			expectedResp: nil,
			wantErr:      false,
		},
		{
			name: "400 Bad Request",
			internalReq: &internalRequest{
				endpoint:            "/test-bad-request",
				method:              http.MethodGet,
				withResponse:        &mockResponse{},
				acceptedStatusCodes: []int{http.StatusOK},
			},
			expectedResp: nil,
			wantErr:      true,
		},
		{
			name: "Test request encoding gzip",
			internalReq: &internalRequest{
				endpoint:            "/test-post-encoding",
				method:              http.MethodPost,
				withRequest:         map[string]string{"key": "value"},
				contentType:         contentTypeJSON,
				withResponse:        &mockResponse{},
				acceptedStatusCodes: []int{http.StatusCreated},
			},
			expectedResp:    &mockResponse{Message: "post successful"},
			contentEncoding: GzipEncoding,
			wantErr:         false,
		},
		{
			name: "Test request encoding deflate",
			internalReq: &internalRequest{
				endpoint:            "/test-post-encoding",
				method:              http.MethodPost,
				withRequest:         map[string]string{"key": "value"},
				contentType:         contentTypeJSON,
				withResponse:        &mockResponse{},
				acceptedStatusCodes: []int{http.StatusCreated},
			},
			expectedResp:    &mockResponse{Message: "post successful"},
			contentEncoding: DeflateEncoding,
			wantErr:         false,
		},
		{
			name: "Test request encoding brotli",
			internalReq: &internalRequest{
				endpoint:            "/test-post-encoding",
				method:              http.MethodPost,
				withRequest:         map[string]string{"key": "value"},
				contentType:         contentTypeJSON,
				withResponse:        &mockResponse{},
				acceptedStatusCodes: []int{http.StatusCreated},
			},
			expectedResp:    &mockResponse{Message: "post successful"},
			contentEncoding: BrotliEncoding,
			wantErr:         false,
		},
		{
			name: "Test response decoding gzip",
			internalReq: &internalRequest{
				endpoint:            "/test-get-encoding",
				method:              http.MethodGet,
				withRequest:         nil,
				withResponse:        &mockData{},
				acceptedStatusCodes: []int{http.StatusOK},
			},
			expectedResp:    &mockData{Name: "foo", Age: 30},
			contentEncoding: GzipEncoding,
			wantErr:         false,
		},
		{
			name: "Test response decoding deflate",
			internalReq: &internalRequest{
				endpoint:            "/test-get-encoding",
				method:              http.MethodGet,
				withRequest:         nil,
				withResponse:        &mockData{},
				acceptedStatusCodes: []int{http.StatusOK},
			},
			expectedResp:    &mockData{Name: "foo", Age: 30},
			contentEncoding: DeflateEncoding,
			wantErr:         false,
		},
		{
			name: "Test response decoding brotli",
			internalReq: &internalRequest{
				endpoint:            "/test-get-encoding",
				method:              http.MethodGet,
				withRequest:         nil,
				withResponse:        &mockData{},
				acceptedStatusCodes: []int{http.StatusOK},
			},
			expectedResp:    &mockData{Name: "foo", Age: 30},
			contentEncoding: BrotliEncoding,
			wantErr:         false,
		},
		{
			name: "Test request and response encoding",
			internalReq: &internalRequest{
				endpoint:            "/test-req-resp-encoding",
				method:              http.MethodPost,
				contentType:         contentTypeJSON,
				withRequest:         &mockData{Name: "foo", Age: 30},
				withResponse:        &mockData{},
				acceptedStatusCodes: []int{http.StatusOK},
			},
			expectedResp:    &mockData{Name: "foo", Age: 30},
			contentEncoding: GzipEncoding,
			wantErr:         false,
		},
		{
			name: "Test successful retries",
			internalReq: &internalRequest{
				endpoint:            "/success-retry",
				method:              http.MethodGet,
				withResponse:        nil,
				withRequest:         nil,
				acceptedStatusCodes: []int{http.StatusOK},
			},
			expectedResp: nil,
			wantErr:      false,
		},
		{
			name: "Test failed retries",
			internalReq: &internalRequest{
				endpoint:            "/failed-retry",
				method:              http.MethodGet,
				withResponse:        nil,
				withRequest:         nil,
				acceptedStatusCodes: []int{http.StatusOK},
			},
			expectedResp: nil,
			wantErr:      true,
		},
		{
			name: "Test disable retries",
			internalReq: &internalRequest{
				endpoint:            "/test-get",
				method:              http.MethodGet,
				withResponse:        nil,
				withRequest:         nil,
				acceptedStatusCodes: []int{http.StatusOK},
			},
			expectedResp: nil,
			disableRetry: true,
			wantErr:      false,
		},
		{
			name: "Test request timeout on retries",
			internalReq: &internalRequest{
				endpoint:            "/failed-retry",
				method:              http.MethodGet,
				withResponse:        nil,
				withRequest:         nil,
				acceptedStatusCodes: []int{http.StatusOK},
			},
			expectedResp: nil,
			withTimeout:  true,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newClient(&http.Client{}, ts.URL, "testApiKey", clientConfig{
				contentEncoding:          tt.contentEncoding,
				encodingCompressionLevel: DefaultCompression,
				maxRetries:               3,
				disableRetry:             tt.disableRetry,
				retryOnStatus: map[int]bool{
					502: true,
					503: true,
					504: true,
				},
			})

			ctx := context.Background()

			if tt.withTimeout {
				timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
				ctx = timeoutCtx
				defer cancel()
			}

			err := c.executeRequest(ctx, tt.internalReq)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedResp, tt.internalReq.withResponse)
			}
		})
	}
}

func TestNewClientNilRetryOnStatus(t *testing.T) {
	c := newClient(&http.Client{}, "", "", clientConfig{
		maxRetries:    3,
		retryOnStatus: nil,
	})

	require.NotNil(t, c.retryOnStatus)
}

func (m mockJsonMarshaller) MarshalJSON() ([]byte, error) {
	type Alias mockJsonMarshaller

	if !m.valid {
		return nil, errors.New("mockJsonMarshaller not valid")
	}

	if m.null {
		return nil, nil
	}

	return json.Marshal(&struct {
		Alias
	}{
		Alias: Alias(m),
	})
}
