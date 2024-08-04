package meilisearch

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock structures for testing
type mockResponse struct {
	Message string `json:"message"`
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
		expectedErr  error
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
			expectedErr:  nil,
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
			expectedErr:  nil,
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
			expectedErr:  &Error{StatusCode: http.StatusNotFound},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.executeRequest(context.Background(), tt.internalReq)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				var apiErr *Error
				assert.ErrorAs(t, err, &apiErr)
				assert.Equal(t, tt.expectedErr.(*Error).StatusCode, apiErr.StatusCode)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, tt.internalReq.withResponse)
			}
		})
	}
}
