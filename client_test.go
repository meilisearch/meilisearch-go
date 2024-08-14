package meilisearch

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
	// Create a mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/test-get" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message":"get successful"}`))
		} else if r.Method == http.MethodPost && r.URL.Path == "/test-post" {
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"message":"post successful"}`))
		} else if r.URL.Path == "/test-bad-request" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message":"bad request"}`))
		} else if r.URL.Path == "/invalid-response-body" {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message":"bad response body"}`))
		} else if r.URL.Path == "/io-reader" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message":"io reader"}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	client := newClient(&http.Client{}, ts.URL, "testApiKey")

	tests := []struct {
		name         string
		internalReq  *internalRequest
		expectedResp interface{}
		wantErr      bool
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.executeRequest(context.Background(), tt.internalReq)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedResp, tt.internalReq.withResponse)
			}
		})
	}
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
