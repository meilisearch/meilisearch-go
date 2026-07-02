package meilisearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockResponse struct {
	Message string `json:"message"`
}

type mockData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type mockJsonMarshaller struct {
	valid bool
	null  bool
	Foo   string `json:"foo"`
	Bar   string `json:"bar"`
}

func (m mockJsonMarshaller) MarshalJSON() ([]byte, error) {
	if !m.valid {
		return nil, errors.New("mockJsonMarshaller not valid")
	}
	if m.null {
		return nil, nil
	}
	return json.Marshal(map[string]string{"foo": m.Foo, "bar": m.Bar})
}

type failingEncoder struct{}

func (fe failingEncoder) Encode(r io.Reader) (io.ReadCloser, error) {
	return nil, errors.New("dummy encoding failure")
}
func (fe failingEncoder) Decode(b []byte, v interface{}) error {
	return errors.New("dummy decode failure")
}
func (fe failingEncoder) Decoder(r io.Reader) (streamDecoder, error) {
	return nil, errors.New("dummy decoder failure")
}

func setupMockServer(t *testing.T) (*httptest.Server, *int) {
	retryCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/success-get":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message":"get successful"}`))

		case "/success-post":
			b, _ := io.ReadAll(r.Body)
			require.NotEmpty(t, b, "POST body should not be empty")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"message":"post successful"}`))

		case "/ndjson-success":
			w.Header().Set("Content-Type", "application/x-ndjson")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"name":"Alice","age":30}` + "\n" + `{"name":"Bob","age":25}`))

		case "/ndjson-malformed":
			w.Header().Set("Content-Type", "application/x-ndjson")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"name":"Alice","age":30` + "\n" + `{"name":"Bob",}`))

		case "/retry-success":
			if retryCount < 2 {
				retryCount++
				w.WriteHeader(http.StatusBadGateway)
				return
			}
			b, _ := io.ReadAll(r.Body)
			if r.Method == http.MethodPost {
				require.NotEmpty(t, b, "POST body must survive retries via GetBody")
			}
			w.WriteHeader(http.StatusOK)

		case "/timeout":
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)

		case "/bad-request":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message":"bad request", "code": "bad_request"}`))

		case "/encoded-post":
			enc := r.Header.Get("Content-Encoding")
			require.NotEmpty(t, enc)
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"message":"encoded successful"}`))

		case "/always-502":
			w.WriteHeader(http.StatusBadGateway)

		case "/return-null":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`null`))

		case "/bad-json":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message": "incomplete JSON`))

		case "/query-params":
			require.Equal(t, "meilisearch", r.URL.Query().Get("q"))
			w.WriteHeader(http.StatusOK)

		case "/bad-request-no-code":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message": "error without code field"}`))

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	return ts, &retryCount
}

func TestClient_ExecuteRequest(t *testing.T) {
	ts, retryCount := setupMockServer(t)
	defer ts.Close()

	tests := []struct {
		name         string
		req          *internalRequest
		setupCtx     func() (context.Context, context.CancelFunc)
		cfg          *clientConfig
		expectedResp interface{}
		wantErr      bool
		errTypeCheck func(err error)
	}{
		{
			name: "GET Success",
			req: &internalRequest{
				endpoint:            "/success-get",
				method:              http.MethodGet,
				withResponse:        &mockResponse{},
				acceptedStatusCodes: []int{http.StatusOK},
			},
			expectedResp: &mockResponse{Message: "get successful"},
		},
		{
			name: "POST Success with raw struct",
			req: &internalRequest{
				endpoint:            "/success-post",
				method:              http.MethodPost,
				contentType:         contentTypeJSON,
				withRequest:         &mockData{Name: "John", Age: 40},
				withResponse:        &mockResponse{},
				acceptedStatusCodes: []int{http.StatusCreated},
			},
			expectedResp: &mockResponse{Message: "post successful"},
		},
		{
			name: "POST Success with io.Reader",
			req: &internalRequest{
				endpoint:            "/success-post",
				method:              http.MethodPost,
				contentType:         contentTypeJSON,
				withRequest:         bytes.NewReader([]byte(`{"name":"John"}`)),
				withResponse:        &mockResponse{},
				acceptedStatusCodes: []int{http.StatusCreated},
			},
			expectedResp: &mockResponse{Message: "post successful"},
		},
		{
			name: "POST Success with []byte",
			req: &internalRequest{
				endpoint:            "/success-post",
				method:              http.MethodPost,
				contentType:         contentTypeJSON,
				withRequest:         []byte(`{"name":"John"}`),
				withResponse:        &mockResponse{},
				acceptedStatusCodes: []int{http.StatusCreated},
			},
			expectedResp: &mockResponse{Message: "post successful"},
		},
		{
			name: "NDJSON Success",
			req: &internalRequest{
				endpoint:            "/ndjson-success",
				method:              http.MethodGet,
				withResponse:        &[]mockData{},
				acceptedContentType: contentTypeNDJSON,
				acceptedStatusCodes: []int{http.StatusOK},
			},
			expectedResp: &[]mockData{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 25}},
		},
		{
			name: "NDJSON Malformed Error",
			req: &internalRequest{
				endpoint:            "/ndjson-malformed",
				method:              http.MethodGet,
				withResponse:        &[]mockData{},
				acceptedContentType: contentTypeNDJSON,
				acceptedStatusCodes: []int{http.StatusOK},
			},
			wantErr: true,
		},
		{
			name: "NDJSON Invalid Destination",
			req: &internalRequest{
				endpoint:            "/ndjson-success",
				method:              http.MethodGet,
				withResponse:        &mockData{},
				acceptedContentType: contentTypeNDJSON,
				acceptedStatusCodes: []int{http.StatusOK},
			},
			wantErr: true,
		},
		{
			name: "Context Timeout",
			setupCtx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Millisecond)
			},
			req: &internalRequest{
				endpoint:            "/timeout",
				method:              http.MethodGet,
				acceptedStatusCodes: []int{http.StatusOK},
			},
			wantErr: true,
			errTypeCheck: func(err error) {
				var e *Error
				require.ErrorAs(t, err, &e)
				assert.Equal(t, MeilisearchTimeoutError, e.ErrCode)
			},
		},
		{
			name: "Context Canceled Manually",
			setupCtx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx, cancel
			},
			req: &internalRequest{
				endpoint:            "/timeout",
				method:              http.MethodGet,
				acceptedStatusCodes: []int{http.StatusOK},
			},
			wantErr: true,
			errTypeCheck: func(err error) {
				var e *Error
				require.ErrorAs(t, err, &e)
				assert.Equal(t, MeilisearchTimeoutError, e.ErrCode)
			},
		},
		{
			name: "Retry Logic Success (GetBody survival)",
			req: &internalRequest{
				endpoint:            "/retry-success",
				method:              http.MethodPost,
				contentType:         contentTypeJSON,
				withRequest:         &mockData{Name: "Retry", Age: 1},
				acceptedStatusCodes: []int{http.StatusOK},
			},
		},
		{
			name: "Disable Retry Fails Immediately",
			cfg: &clientConfig{
				disableRetry: true,
			},
			req: &internalRequest{
				endpoint:            "/retry-success",
				method:              http.MethodGet,
				acceptedStatusCodes: []int{http.StatusOK},
			},
			wantErr: true,
		},
		{
			name: "Bad Request API Error",
			req: &internalRequest{
				endpoint:            "/bad-request",
				method:              http.MethodGet,
				acceptedStatusCodes: []int{http.StatusOK},
			},
			wantErr: true,
			errTypeCheck: func(err error) {
				var e *Error
				require.ErrorAs(t, err, &e)
				assert.Equal(t, MeilisearchApiError, e.ErrCode)
				assert.Equal(t, "bad_request", e.MeilisearchApiError.Code)
			},
		},
		{
			name: "Validation Error: GET with Body",
			req: &internalRequest{
				endpoint:            "/success-get",
				method:              http.MethodGet,
				withRequest:         &mockData{},
				acceptedStatusCodes: []int{http.StatusOK},
			},
			wantErr: true,
			errTypeCheck: func(err error) {
				assert.ErrorIs(t, err, ErrInvalidRequestMethod)
			},
		},
		{
			name: "Max Retries Exceeded",
			req: &internalRequest{
				endpoint:            "/always-502",
				method:              http.MethodGet,
				acceptedStatusCodes: []int{http.StatusOK},
			},
			wantErr: true,
			errTypeCheck: func(err error) {
				var e *Error
				require.ErrorAs(t, err, &e)
				assert.Equal(t, MeilisearchMaxRetriesExceeded, e.ErrCode)
			},
		},
		{
			name: "Response Null Body",
			req: &internalRequest{
				endpoint:            "/return-null",
				method:              http.MethodGet,
				withResponse:        &mockResponse{},
				acceptedStatusCodes: []int{http.StatusOK},
			},
			expectedResp: nil, // Triggers req.withResponse = nil
		},
		{
			name: "JSON Unmarshal Error",
			req: &internalRequest{
				endpoint:            "/bad-json",
				method:              http.MethodGet,
				withResponse:        &mockResponse{},
				acceptedStatusCodes: []int{http.StatusOK},
			},
			wantErr: true,
			errTypeCheck: func(err error) {
				var e *Error
				require.ErrorAs(t, err, &e)
				assert.Equal(t, ErrCodeResponseUnmarshalBody, e.ErrCode)
			},
		},
		{
			name: "JSON Marshal Error in buildBody",
			req: &internalRequest{
				endpoint:            "/success-post",
				method:              http.MethodPost,
				contentType:         contentTypeJSON,
				withRequest:         mockJsonMarshaller{valid: false}, // Triggers marshal error
				acceptedStatusCodes: []int{http.StatusCreated},
			},
			wantErr: true,
			errTypeCheck: func(err error) {
				var e *Error
				require.ErrorAs(t, err, &e)
				assert.Equal(t, ErrCodeMarshalRequest, e.ErrCode)
			},
		},
		{
			name: "Query Parameters check",
			req: &internalRequest{
				endpoint:            "/query-params",
				method:              http.MethodGet,
				withQueryParams:     map[string]string{"q": "meilisearch"}, // Triggers query encoder
				acceptedStatusCodes: []int{http.StatusOK},
			},
		},
		{
			name: "API Error Without Code",
			req: &internalRequest{
				endpoint:            "/bad-request-no-code",
				method:              http.MethodGet,
				acceptedStatusCodes: []int{http.StatusOK},
			},
			wantErr: true,
			errTypeCheck: func(err error) {
				var e *Error
				require.ErrorAs(t, err, &e)
				assert.Equal(t, MeilisearchApiErrorWithoutMessage, e.ErrCode)
			},
		},
		{
			name: "Encoder Failure: Request Compress",
			cfg: &clientConfig{
				contentEncoding: GzipEncoding,
			},
			req: &internalRequest{
				endpoint:            "/encoded-post",
				method:              http.MethodPost,
				contentType:         contentTypeJSON,
				withRequest:         &mockData{Name: "Gzip", Age: 10},
				acceptedStatusCodes: []int{http.StatusCreated},
			},
			wantErr: true,
		},
		{
			name: "Encoder Failure: Response Decode",
			cfg: &clientConfig{
				contentEncoding: GzipEncoding,
			},
			req: &internalRequest{
				endpoint:            "/success-get",
				method:              http.MethodGet,
				withResponse:        &mockResponse{},
				acceptedStatusCodes: []int{http.StatusOK},
			},
			wantErr: true,
		},
		{
			name: "Validation Error: No Content-Type for POST body",
			req: &internalRequest{
				endpoint:            "/success-post",
				method:              http.MethodPost,
				withRequest:         &mockData{},
				acceptedStatusCodes: []int{http.StatusCreated},
			},
			wantErr: true,
			errTypeCheck: func(err error) {
				assert.ErrorIs(t, err, ErrRequestBodyWithoutContentType)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			*retryCount = 0

			cfg := tt.cfg
			if cfg == nil {
				cfg = &clientConfig{
					maxRetries:    3,
					retryOnStatus: map[int]bool{502: true, 503: true, 504: true},
					jsonMarshal:   json.Marshal,
					jsonUnmarshal: json.Unmarshal,
				}
			} else {
				if cfg.jsonMarshal == nil {
					cfg.jsonMarshal = json.Marshal
					cfg.jsonUnmarshal = json.Unmarshal
				}
			}

			c := newClient(&http.Client{}, ts.URL, "testApiKey", cfg)

			if strings.Contains(tt.name, "Encoder Failure") {
				c.encoder = failingEncoder{}
			}

			ctx := context.Background()
			if tt.setupCtx != nil {
				var cancel context.CancelFunc
				ctx, cancel = tt.setupCtx()
				defer cancel()
			}

			err := c.executeRequest(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errTypeCheck != nil {
					tt.errTypeCheck(err)
				}
			} else {
				require.NoError(t, err)

				if tt.name == "Response Null Body" {
					assert.Nil(t, tt.req.withResponse)
				} else if tt.expectedResp != nil {
					assert.Equal(t, tt.expectedResp, tt.req.withResponse)
				}
			}
		})
	}
}

type mockRoundTripper struct {
	fn func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.fn(req)
}

// errorReadCloser simulates an io.Reader that fails midway through reading
type errorReadCloser struct{}

func (e errorReadCloser) Read(p []byte) (n int, err error) {
	return 0, errors.New("mock read error")
}

func (e errorReadCloser) Close() error {
	return nil
}

func TestClient_Coverage_EdgeCases(t *testing.T) {
	mockHTTP := &http.Client{
		Transport: &mockRoundTripper{
			fn: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"message": "ok"}`))),
				}, nil
			},
		},
	}
	c := newClient(mockHTTP, "http://localhost", "key", &clientConfig{
		disableRetry: true,
	})

	t.Run("sendRequest - URL Parse Error", func(t *testing.T) {
		err := c.executeRequest(context.Background(), &internalRequest{
			endpoint: ":\x7f//invalid", // Invalid control character forces url.Parse to fail
		})
		require.Error(t, err)
	})

	t.Run("sendRequest - http.NewRequestWithContext Error", func(t *testing.T) {
		err := c.executeRequest(context.Background(), &internalRequest{
			endpoint: "/test",
			method:   "B@D M3THOD", // Invalid HTTP method syntax
		})
		require.Error(t, err)
	})

	t.Run("sendRequest - Body ReadAll Error", func(t *testing.T) {
		err := c.executeRequest(context.Background(), &internalRequest{
			endpoint:    "/test",
			method:      http.MethodPost,
			contentType: contentTypeJSON,
			withRequest: errorReadCloser{}, // Fails when sendRequest calls io.ReadAll()
		})
		require.Error(t, err)
	})

	t.Run("executeRequest - Response ReadAll Error", func(t *testing.T) {
		errHTTP := &http.Client{
			Transport: &mockRoundTripper{
				fn: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       errorReadCloser{}, // Fails when client reads the response body
					}, nil
				},
			},
		}
		cErrResp := newClient(errHTTP, "http://localhost", "key", &clientConfig{
			disableRetry: true,
		})
		err := cErrResp.executeRequest(context.Background(), &internalRequest{
			endpoint:            "/test",
			method:              http.MethodGet,
			acceptedStatusCodes: []int{http.StatusOK},
		})
		require.Error(t, err)
	})

	t.Run("do - Backoff Context Cancellation", func(t *testing.T) {
		retryHTTP := &http.Client{
			Transport: &mockRoundTripper{
				fn: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusBadGateway,
						Body:       io.NopCloser(bytes.NewReader([]byte{})),
					}, nil
				},
			},
		}
		cRetry := newClient(retryHTTP, "http://localhost", "key", &clientConfig{
			maxRetries:    3,
			retryOnStatus: map[int]bool{http.StatusBadGateway: true},
		})

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := cRetry.executeRequest(ctx, &internalRequest{
			endpoint:            "/test",
			method:              http.MethodGet,
			acceptedStatusCodes: []int{http.StatusOK},
		})

		require.Error(t, err)
		var e *Error
		require.ErrorAs(t, err, &e)
		require.Equal(t, MeilisearchTimeoutError, e.ErrCode)
	})

	t.Run("do - GetBody Rewind Error", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "http://localhost", nil)
		req.GetBody = func() (io.ReadCloser, error) {
			return nil, errors.New("mock getbody error")
		}
		cRetry := newClient(mockHTTP, "http://localhost", "key", &clientConfig{
			maxRetries:    3,
			retryOnStatus: map[int]bool{http.StatusOK: true}, // Force retry to trigger GetBody
		})
		cRetry.retryBackoff = func(attempt uint8) time.Duration { return 0 }
		_, err := cRetry.do(req, &Error{})
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to rewind body")
	})

	t.Run("executeRequest - handleContentType Error", func(t *testing.T) {
		xmlHTTP := &http.Client{
			Transport: &mockRoundTripper{
				fn: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"Content-Type": []string{"application/json"}},
						Body:       io.NopCloser(bytes.NewReader([]byte(`{}`))),
					}, nil
				},
			},
		}
		cXML := newClient(xmlHTTP, "http://localhost", "key", &clientConfig{
			disableRetry: true,
		})
		err := cXML.executeRequest(context.Background(), &internalRequest{
			endpoint:            "/test",
			method:              http.MethodGet,
			acceptedContentType: "text/xml", // Forces content type mismatch error
			acceptedStatusCodes: []int{http.StatusOK},
		})
		require.Error(t, err)
	})

	t.Run("executeRequest - handleResponse Compressed Unmarshal Error", func(t *testing.T) {
		cGzip := newClient(mockHTTP, "http://localhost", "key", &clientConfig{
			disableRetry:    true,
			contentEncoding: GzipEncoding,
		})
		err := cGzip.executeRequest(context.Background(), &internalRequest{
			endpoint:            "/test",
			method:              http.MethodGet,
			withResponse:        &mockResponse{},
			acceptedStatusCodes: []int{http.StatusOK},
		})
		require.Error(t, err) // Fails because mockHTTP returns plain JSON, not a valid gzip stream
	})

	t.Run("handleStatusCode - Nil Accepted Status Codes", func(t *testing.T) {
		err := c.executeRequest(context.Background(), &internalRequest{
			endpoint: "/test",
			method:   http.MethodGet,
			// acceptedStatusCodes intentionally omitted (nil)
		})
		require.NoError(t, err)
	})

	t.Run("handleStatusCode - API Error Without Code Field", func(t *testing.T) {
		noCodeHTTP := &http.Client{
			Transport: &mockRoundTripper{
				fn: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       io.NopCloser(bytes.NewReader([]byte(`{"message":"error missing code field"}`))),
					}, nil
				},
			},
		}
		cNoCode := newClient(noCodeHTTP, "http://localhost", "key", &clientConfig{
			disableRetry: true,
		})
		err := cNoCode.executeRequest(context.Background(), &internalRequest{
			endpoint:            "/test",
			method:              http.MethodGet,
			acceptedStatusCodes: []int{http.StatusOK},
		})

		require.Error(t, err)

		var e *Error
		require.ErrorAs(t, err, &e)
		require.Equal(t, MeilisearchApiErrorWithoutMessage, e.ErrCode)
	})

	t.Run("handleResponse - Null Body Literal", func(t *testing.T) {
		nullHTTP := &http.Client{
			Transport: &mockRoundTripper{
				fn: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`null`))),
					}, nil
				},
			},
		}
		cNull := newClient(nullHTTP, "http://localhost", "key", &clientConfig{
			disableRetry: true,
		})
		var target mockResponse
		err := cNull.executeRequest(context.Background(), &internalRequest{
			endpoint:            "/test",
			method:              http.MethodGet,
			withResponse:        &target,
			acceptedStatusCodes: []int{http.StatusOK},
		})
		require.NoError(t, err)
	})
}
