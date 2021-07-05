package meilisearch

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func TestClient_Version(t *testing.T) {
	tests := []struct {
		name   string
		client *Client
	}{
		{
			name:   "TestVersion",
			client: defaultClient,
		},
		{
			name:   "TestVersionWithCustomClient",
			client: customClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.client.GetVersion()
			require.NoError(t, err)
			require.NotNil(t, gotResp, "Version() should not return nil value")
			require.Equal(t, "0.21.0", gotResp.PkgVersion)
		})
	}
}

func TestClient_TimeoutError(t *testing.T) {
	tests := []struct {
		name          string
		client        *Client
		expectedError Error
	}{
		{
			name:   "TestTimeoutError",
			client: timeoutClient,
			expectedError: Error(Error{
				Endpoint:         "/version",
				Method:           "GET",
				Function:         "Version",
				RequestToString:  "empty request",
				ResponseToString: "empty response",
				MeilisearchApiMessage: meilisearchApiMessage{
					Message:   "empty meilisearch message",
					ErrorCode: "",
					ErrorType: "",
					ErrorLink: "",
				},
				StatusCode:         0,
				StatusCodeExpected: []int{200},
				rawMessage:         "MeilisearchTimeoutError (path \"${method} ${endpoint}\" with method \"${function}\")",
				OriginError:        fasthttp.ErrTimeout,
				ErrCode:            6,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.client.GetVersion()
			require.Error(t, err)
			require.Nil(t, gotResp)
			require.Equal(t, &tt.expectedError, err)
		})
	}
}

func TestClient_GetAllStats(t *testing.T) {
	tests := []struct {
		name   string
		client *Client
	}{
		{
			name:   "TestGetAllStats",
			client: defaultClient,
		},
		{
			name:   "TestGetAllStatsWithCustomClient",
			client: customClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.client.GetAllStats()
			require.NoError(t, err)
			require.NotNil(t, gotResp, "GetAllStats() should not return nil value")
		})
	}
}

func TestClient_GetKeys(t *testing.T) {
	tests := []struct {
		name   string
		client *Client
	}{
		{
			name:   "TestGetKeys",
			client: defaultClient,
		},
		{
			name:   "TestGetKeysWithCustomClient",
			client: defaultClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.client.GetKeys()
			require.NoError(t, err)
			require.NotNil(t, gotResp, "GetKeys() should not return nil value")
		})
	}
}

func TestClient_Health(t *testing.T) {
	tests := []struct {
		name          string
		client        *Client
		wantResp      *Health
		wantErr       bool
		expectedError Error
	}{
		{
			name:   "TestHealth",
			client: defaultClient,
			wantResp: &Health{
				Status: "available",
			},
			wantErr: false,
		},
		{
			name:   "TestHealthWithCustomClient",
			client: customClient,
			wantResp: &Health{
				Status: "available",
			},
			wantErr: false,
		},
		{
			name: "TestHealthWIthBadUrl",
			client: &Client{
				config: ClientConfig{
					Host:   "http://wrongurl:1234",
					APIKey: masterKey,
				},
				httpClient: &fasthttp.Client{
					Name: "meilsearch-client",
				},
			},
			wantErr: true,
			expectedError: Error(Error{
				Endpoint:         "/health",
				Method:           "GET",
				Function:         "Health",
				RequestToString:  "empty request",
				ResponseToString: "empty response",
				MeilisearchApiMessage: meilisearchApiMessage{
					Message:   "empty meilisearch message",
					ErrorCode: "",
					ErrorType: "",
					ErrorLink: "",
				},
				StatusCode:         0,
				StatusCodeExpected: []int{200},
				rawMessage:         "MeilisearchCommunicationError unable to execute request (path \"${method} ${endpoint}\" with method \"${function}\")",
				OriginError: &net.DNSError{Err: "Temporary failure in name resolution",
					Name:        "wrongurl",
					Server:      "",
					IsTimeout:   false,
					IsTemporary: true,
					IsNotFound:  false},
				ErrCode: 7,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.client.Health()
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, &tt.expectedError, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantResp, gotResp, "Health() got response %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestClient_IsHealthy(t *testing.T) {
	tests := []struct {
		name   string
		client *Client
		want   bool
	}{
		{
			name:   "TestIsHealthy",
			client: defaultClient,
			want:   true,
		},
		{
			name:   "TestIsHealthyWithCustomClient",
			client: customClient,
			want:   true,
		},
		{
			name: "TestIsHealthyWIthBadUrl",
			client: &Client{
				config: ClientConfig{
					Host:   "http://wrongurl:1234",
					APIKey: masterKey,
				},
				httpClient: &fasthttp.Client{
					Name: "meilsearch-client",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.client.IsHealthy()
			require.Equal(t, tt.want, got, "IsHealthy() got response %v, want %v", got, tt.want)
		})
	}
}

func TestClient_CreateDump(t *testing.T) {
	tests := []struct {
		name     string
		client   *Client
		wantResp *Dump
	}{
		{
			name:   "TestCreateDump",
			client: defaultClient,
			wantResp: &Dump{
				Status: "in_progress",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client

			gotResp, err := c.CreateDump()
			require.NoError(t, err)
			if assert.NotNil(t, gotResp, "CreateDump() should not return nil value") {
				require.Equal(t, tt.wantResp.Status, gotResp.Status, "CreateDump() got response status %v, want: %v", gotResp.Status, tt.wantResp.Status)
			}

			// Waiting for CreateDump() to finished
			for {
				gotResp, _ := c.GetDumpStatus(gotResp.UID)
				if gotResp.Status == "done" {
					break
				}
			}
		})
	}
}

func TestClient_GetDumpStatus(t *testing.T) {
	tests := []struct {
		name     string
		client   *Client
		wantResp []string
		wantErr  bool
	}{
		{
			name:     "TestGetDumpStatus",
			client:   defaultClient,
			wantResp: []string{"in_progress", "failed", "done"},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client

			dump, err := c.CreateDump()
			require.NoError(t, err, "CreateDump() in TestGetDumpStatus error should be nil")

			gotResp, err := c.GetDumpStatus(dump.UID)
			require.NoError(t, err)
			require.Contains(t, tt.wantResp, gotResp.Status, "GetDumpStatus() got response status %v", gotResp.Status)
			require.NotEqual(t, "failed", gotResp.Status, "GetDumpStatus() response status should not be failed")
		})
	}
}
