package meilisearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func TestClient_CreateIndex(t *testing.T) {
	type args struct {
		config IndexConfig
	}
	tests := []struct {
		name          string
		client        *Client
		args          args
		wantResp      *Index
		wantErr       bool
		expectedError Error
	}{
		{
			name:   "TestBasicCreateIndex",
			client: defaultClient,
			args: args{
				config: IndexConfig{
					Uid: "TestBasicCreateIndex",
				},
			},
			wantResp: &Index{
				UID: "TestBasicCreateIndex",
			},
			wantErr: false,
		},
		{
			name:   "TestCreateIndexWithCustomClient",
			client: customClient,
			args: args{
				config: IndexConfig{
					Uid: "TestBasicCreateIndex",
				},
			},
			wantResp: &Index{
				UID: "TestBasicCreateIndex",
			},
			wantErr: false,
		},
		{
			name:   "TestCreateIndexWithCustomClient",
			client: customClient,
			args: args{
				config: IndexConfig{
					Uid: "TestBasicCreateIndex",
				},
			},
			wantResp: &Index{
				UID: "TestBasicCreateIndex",
			},
			wantErr: false,
		},
		{
			name:   "TestCreateIndexWithPrimaryKey",
			client: defaultClient,
			args: args{
				config: IndexConfig{
					Uid:        "TestCreateIndexWithPrimaryKey",
					PrimaryKey: "PrimaryKey",
				},
			},
			wantResp: &Index{
				UID:        "TestCreateIndexWithPrimaryKey",
				PrimaryKey: "PrimaryKey",
			},
			wantErr: false,
		},
		{
			name:   "TestCreateIndexInvalidUid",
			client: defaultClient,
			args: args{
				config: IndexConfig{
					Uid: "TestCreateIndexInvalidUid*",
				},
			},
			wantErr: true,
			expectedError: Error(Error{
				Endpoint:         "/indexes",
				Method:           "POST",
				Function:         "CreateIndex",
				RequestToString:  "{\"uid\":\"TestCreateIndexInvalidUid*\"}",
				ResponseToString: "{\"message\":\"Index must have a valid uid; Index uid can be of type integer or string only composed of alphanumeric characters, hyphens (-) and underscores (_).\",\"errorCode\":\"invalid_index_uid\",\"errorType\":\"invalid_request_error\",\"errorLink\":\"https://docs.meilisearch.com/errors#invalid_index_uid\"}",
				MeilisearchApiMessage: meilisearchApiMessage{
					Message:   "Index must have a valid uid; Index uid can be of type integer or string only composed of alphanumeric characters, hyphens (-) and underscores (_).",
					ErrorCode: "invalid_index_uid",
					ErrorType: "invalid_request_error",
					ErrorLink: "https://docs.meilisearch.com/errors#invalid_index_uid",
				},
				StatusCode:         400,
				StatusCodeExpected: []int{201},
				rawMessage:         "unaccepted status code found: ${statusCode} expected: ${statusCodeExpected}, MeilisearchApiError Message: ${message}, ErrorCode: ${errorCode}, ErrorType: ${errorType}, ErrorLink: ${errorLink} (path \"${method} ${endpoint}\" with method \"${function}\")",
				OriginError:        error(nil),
				ErrCode:            4,
			}),
		},
		{
			name:   "TestCreateIndexAlreadyExist",
			client: defaultClient,
			args: args{
				config: IndexConfig{
					Uid: "indexUID",
				},
			},
			wantErr: true,
			expectedError: Error(Error{
				Endpoint:         "/indexes",
				Method:           "POST",
				Function:         "CreateIndex",
				RequestToString:  "{\"uid\":\"indexUID\"}",
				ResponseToString: "{\"message\":\"Index already exists.\",\"errorCode\":\"index_already_exists\",\"errorType\":\"invalid_request_error\",\"errorLink\":\"https://docs.meilisearch.com/errors#index_already_exists\"}",
				MeilisearchApiMessage: meilisearchApiMessage{
					Message:   "Index already exists.",
					ErrorCode: "index_already_exists",
					ErrorType: "invalid_request_error",
					ErrorLink: "https://docs.meilisearch.com/errors#index_already_exists",
				},
				StatusCode:         400,
				StatusCodeExpected: []int{201},
				rawMessage:         "unaccepted status code found: ${statusCode} expected: ${statusCodeExpected}, MeilisearchApiError Message: ${message}, ErrorCode: ${errorCode}, ErrorType: ${errorType}, ErrorLink: ${errorLink} (path \"${method} ${endpoint}\" with method \"${function}\")",
				OriginError:        error(nil),
				ErrCode:            4,
			}),
		},
		{
			name:   "TestCreateIndexTimeout",
			client: timeoutClient,
			args: args{
				config: IndexConfig{
					Uid: "indexUID",
				},
			},
			wantErr: true,
			expectedError: Error(Error{
				Endpoint:         "/indexes",
				Method:           "POST",
				Function:         "CreateIndex",
				RequestToString:  "{\"uid\":\"indexUID\"}",
				ResponseToString: "empty response",
				MeilisearchApiMessage: meilisearchApiMessage{
					Message:   "empty meilisearch message",
					ErrorCode: "",
					ErrorType: "",
					ErrorLink: "",
				},
				StatusCode:         0,
				StatusCodeExpected: []int{201},
				rawMessage:         "MeilisearchTimeoutError (path \"${method} ${endpoint}\" with method \"${function}\")",
				OriginError:        fasthttp.ErrTimeout,
				ErrCode:            6,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client
			SetUpBasicIndex()

			gotResp, err := c.CreateIndex(&tt.args.config)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, &tt.expectedError, err)
			} else {
				require.NoError(t, err)
				if assert.NotNil(t, gotResp) {
					require.Equal(t, tt.wantResp.UID, gotResp.UID)
					require.Equal(t, tt.wantResp.PrimaryKey, gotResp.PrimaryKey)
				}
			}

			deleteAllIndexes(c)
		})
	}
}

func TestClient_DeleteIndex(t *testing.T) {
	type args struct {
		createUid []string
		deleteUid []string
	}
	tests := []struct {
		name          string
		client        *Client
		args          args
		wantErr       bool
		expectedError []Error
	}{
		{
			name:   "TestBasicDeleteIndex",
			client: defaultClient,
			args: args{
				createUid: []string{"1"},
				deleteUid: []string{"1"},
			},
			wantErr: false,
		},
		{
			name:   "TestDeleteIndexWithCustomClient",
			client: customClient,
			args: args{
				createUid: []string{"1"},
				deleteUid: []string{"1"},
			},
			wantErr: false,
		},
		{
			name:   "TestMultipleDeleteIndex",
			client: defaultClient,
			args: args{
				createUid: []string{"2", "3", "4", "5"},
				deleteUid: []string{"2", "3", "4", "5"},
			},
			wantErr: false,
		},
		{
			name:   "TestNotExistingDeleteIndex",
			client: defaultClient,
			args: args{
				deleteUid: []string{"1"},
			},
			wantErr: true,
			expectedError: []Error{
				Error(Error{
					Endpoint:         "/indexes/1",
					Method:           "DELETE",
					Function:         "DeleteIndex",
					RequestToString:  "empty request",
					ResponseToString: "{\"message\":\"Index \\\"1\\\" not found.\",\"errorCode\":\"index_not_found\",\"errorType\":\"invalid_request_error\",\"errorLink\":\"https://docs.meilisearch.com/errors#index_not_found\"}",
					MeilisearchApiMessage: meilisearchApiMessage{
						Message:   "Index \"1\" not found.",
						ErrorCode: "index_not_found",
						ErrorType: "invalid_request_error",
						ErrorLink: "https://docs.meilisearch.com/errors#index_not_found",
					},
					StatusCode:         404,
					StatusCodeExpected: []int{204},
					rawMessage:         "unaccepted status code found: ${statusCode} expected: ${statusCodeExpected}, MeilisearchApiError Message: ${message}, ErrorCode: ${errorCode}, ErrorType: ${errorType}, ErrorLink: ${errorLink} (path \"${method} ${endpoint}\" with method \"${function}\")",
					OriginError:        error(nil),
					ErrCode:            4,
				})},
		},
		{
			name:   "TestMultipleNotExistingDeleteIndex",
			client: defaultClient,
			args: args{
				deleteUid: []string{"2", "3", "4", "5"},
			},
			wantErr: true,
			expectedError: []Error{
				Error(Error{
					Endpoint:         "/indexes/2",
					Method:           "DELETE",
					Function:         "DeleteIndex",
					RequestToString:  "empty request",
					ResponseToString: "{\"message\":\"Index \\\"2\\\" not found.\",\"errorCode\":\"index_not_found\",\"errorType\":\"invalid_request_error\",\"errorLink\":\"https://docs.meilisearch.com/errors#index_not_found\"}",
					MeilisearchApiMessage: meilisearchApiMessage{
						Message:   "Index \"2\" not found.",
						ErrorCode: "index_not_found",
						ErrorType: "invalid_request_error",
						ErrorLink: "https://docs.meilisearch.com/errors#index_not_found"},
					StatusCode:         404,
					StatusCodeExpected: []int{204},
					rawMessage:         "unaccepted status code found: ${statusCode} expected: ${statusCodeExpected}, MeilisearchApiError Message: ${message}, ErrorCode: ${errorCode}, ErrorType: ${errorType}, ErrorLink: ${errorLink} (path \"${method} ${endpoint}\" with method \"${function}\")",
					OriginError:        error(nil),
					ErrCode:            4,
				}),
				Error(Error{
					Endpoint:         "/indexes/3",
					Method:           "DELETE",
					Function:         "DeleteIndex",
					RequestToString:  "empty request",
					ResponseToString: "{\"message\":\"Index \\\"3\\\" not found.\",\"errorCode\":\"index_not_found\",\"errorType\":\"invalid_request_error\",\"errorLink\":\"https://docs.meilisearch.com/errors#index_not_found\"}",
					MeilisearchApiMessage: meilisearchApiMessage{
						Message:   "Index \"3\" not found.",
						ErrorCode: "index_not_found",
						ErrorType: "invalid_request_error",
						ErrorLink: "https://docs.meilisearch.com/errors#index_not_found"},
					StatusCode:         404,
					StatusCodeExpected: []int{204},
					rawMessage:         "unaccepted status code found: ${statusCode} expected: ${statusCodeExpected}, MeilisearchApiError Message: ${message}, ErrorCode: ${errorCode}, ErrorType: ${errorType}, ErrorLink: ${errorLink} (path \"${method} ${endpoint}\" with method \"${function}\")", OriginError: error(nil), ErrCode: 4}),
				Error(Error{
					Endpoint:         "/indexes/4",
					Method:           "DELETE",
					Function:         "DeleteIndex",
					RequestToString:  "empty request",
					ResponseToString: "{\"message\":\"Index \\\"4\\\" not found.\",\"errorCode\":\"index_not_found\",\"errorType\":\"invalid_request_error\",\"errorLink\":\"https://docs.meilisearch.com/errors#index_not_found\"}",
					MeilisearchApiMessage: meilisearchApiMessage{
						Message:   "Index \"4\" not found.",
						ErrorCode: "index_not_found",
						ErrorType: "invalid_request_error",
						ErrorLink: "https://docs.meilisearch.com/errors#index_not_found"},
					StatusCode:         404,
					StatusCodeExpected: []int{204},
					rawMessage:         "unaccepted status code found: ${statusCode} expected: ${statusCodeExpected}, MeilisearchApiError Message: ${message}, ErrorCode: ${errorCode}, ErrorType: ${errorType}, ErrorLink: ${errorLink} (path \"${method} ${endpoint}\" with method \"${function}\")",
					OriginError:        error(nil),
					ErrCode:            4,
				}),
				Error(Error{
					Endpoint:         "/indexes/5",
					Method:           "DELETE",
					Function:         "DeleteIndex",
					RequestToString:  "empty request",
					ResponseToString: "{\"message\":\"Index \\\"5\\\" not found.\",\"errorCode\":\"index_not_found\",\"errorType\":\"invalid_request_error\",\"errorLink\":\"https://docs.meilisearch.com/errors#index_not_found\"}",
					MeilisearchApiMessage: meilisearchApiMessage{
						Message:   "Index \"5\" not found.",
						ErrorCode: "index_not_found",
						ErrorType: "invalid_request_error",
						ErrorLink: "https://docs.meilisearch.com/errors#index_not_found"},
					StatusCode:         404,
					StatusCodeExpected: []int{204},
					rawMessage:         "unaccepted status code found: ${statusCode} expected: ${statusCodeExpected}, MeilisearchApiError Message: ${message}, ErrorCode: ${errorCode}, ErrorType: ${errorType}, ErrorLink: ${errorLink} (path \"${method} ${endpoint}\" with method \"${function}\")",
					OriginError:        error(nil),
					ErrCode:            4,
				}),
			},
		},
		{
			name:   "TestDeleteIndexTimeout",
			client: timeoutClient,
			args: args{
				deleteUid: []string{"1"},
			},
			wantErr: true,
			expectedError: []Error{
				Error(Error{
					Endpoint:         "/indexes/1",
					Method:           "DELETE",
					Function:         "DeleteIndex",
					RequestToString:  "empty request",
					ResponseToString: "empty response",
					MeilisearchApiMessage: meilisearchApiMessage{
						Message:   "empty meilisearch message",
						ErrorCode: "",
						ErrorType: "",
						ErrorLink: "",
					},
					StatusCode:         0,
					StatusCodeExpected: []int{204},
					rawMessage:         "MeilisearchTimeoutError (path \"${method} ${endpoint}\" with method \"${function}\")",
					OriginError:        fasthttp.ErrTimeout,
					ErrCode:            6,
				})},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client

			for _, uid := range tt.args.createUid {
				_, err := c.CreateIndex(&IndexConfig{Uid: uid})
				require.NoError(t, err, "CreateIndex() in TestDeleteIndex error should be nil")
			}
			for k := range tt.args.deleteUid {
				gotOk, err := c.DeleteIndex(tt.args.deleteUid[k])
				if tt.wantErr {
					require.Error(t, err)
					require.Equal(t, &tt.expectedError[k], err)
				} else {
					require.NoError(t, err)
					require.True(t, gotOk)
				}
			}

			deleteAllIndexes(c)
		})
	}
}

func TestClient_GetAllIndexes(t *testing.T) {
	type args struct {
		uid []string
	}
	tests := []struct {
		name     string
		client   *Client
		args     args
		wantResp []Index
	}{
		{
			name:   "TestGelAllIndexesOnNoIndexes",
			client: defaultClient,
			args: args{
				uid: []string{},
			},
			wantResp: []Index{},
		},
		{
			name:   "TestBasicGelAllIndexes",
			client: defaultClient,
			args: args{
				uid: []string{"1"},
			},
			wantResp: []Index{
				{
					UID: "1",
				},
			},
		},
		{
			name:   "TestGelAllIndexesWithCustomClient",
			client: customClient,
			args: args{
				uid: []string{"1"},
			},
			wantResp: []Index{
				{
					UID: "1",
				},
			},
		},
		{
			name:   "TestGelAllIndexesOnMultipleIndex",
			client: defaultClient,
			args: args{
				uid: []string{"1", "2", "3"},
			},
			wantResp: []Index{
				{
					UID: "1",
				},
				{
					UID: "2",
				},
				{
					UID: "3",
				},
			},
		},
		{
			name:   "TestGelAllIndexesOnMultipleIndexWithPrimaryKey",
			client: defaultClient,
			args: args{
				uid: []string{"1", "2", "3"},
			},
			wantResp: []Index{
				{
					UID:        "1",
					PrimaryKey: "PrimaryKey1",
				},
				{
					UID:        "2",
					PrimaryKey: "PrimaryKey2",
				},
				{
					UID:        "3",
					PrimaryKey: "PrimaryKey3",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client

			for _, uid := range tt.args.uid {
				_, err := c.CreateIndex(&IndexConfig{Uid: uid})
				require.NoError(t, err, "CreateIndex() in TestGetAllIndexes error should be nil")
			}
			gotResp, err := c.GetAllIndexes()
			require.NoError(t, err)
			require.Equal(t, len(tt.wantResp), len(gotResp))

			deleteAllIndexes(c)
		})
	}
}

func TestClient_GetIndex(t *testing.T) {
	type args struct {
		config     IndexConfig
		createdUid string
		uid        string
	}
	tests := []struct {
		name     string
		client   *Client
		args     args
		wantResp *Index
		wantErr  bool
		wantCmp  bool
	}{
		{
			name:   "TestBasicGetIndex",
			client: defaultClient,
			args: args{
				config: IndexConfig{
					Uid: "1",
				},
				uid: "1",
			},
			wantResp: &Index{
				UID: "1",
			},
			wantErr: false,
			wantCmp: false,
		},
		{
			name:   "TestGetIndexWithCustomClient",
			client: customClient,
			args: args{
				config: IndexConfig{
					Uid: "1",
				},
				uid: "1",
			},
			wantResp: &Index{
				UID: "1",
			},
			wantErr: false,
			wantCmp: false,
		},
		{
			name:   "TestGetIndexWithPrimaryKey",
			client: defaultClient,
			args: args{
				config: IndexConfig{
					Uid:        "1",
					PrimaryKey: "PrimaryKey",
				},
				uid: "1",
			},
			wantResp: &Index{
				UID:        "1",
				PrimaryKey: "PrimaryKey",
			},
			wantErr: false,
			wantCmp: false,
		},
		{
			name:   "TestGetIndexOnNotExistingIndex",
			client: defaultClient,
			args: args{
				config: IndexConfig{},
				uid:    "1",
			},
			wantResp: nil,
			wantErr:  true,
			wantCmp:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client

			gotCreatedResp, err := c.CreateIndex(&tt.args.config)
			gotResp, err := c.GetIndex(tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantResp.UID, gotResp.UID)
				require.Equal(t, gotCreatedResp.UID, gotResp.UID)
				require.Equal(t, tt.wantResp.PrimaryKey, gotResp.PrimaryKey)
				require.Equal(t, gotCreatedResp.PrimaryKey, gotResp.PrimaryKey)
			}

			deleteAllIndexes(c)
		})
	}
}

func TestClient_GetOrCreateIndex(t *testing.T) {
	type args struct {
		config IndexConfig
	}
	tests := []struct {
		name     string
		client   *Client
		args     args
		wantResp *Index
	}{
		{
			name:   "TestBasicGetOrCreateIndex",
			client: defaultClient,
			args: args{
				config: IndexConfig{
					Uid: "TestBasicGetOrCreateIndex",
				},
			},
			wantResp: &Index{
				UID: "TestBasicGetOrCreateIndex",
			},
		},
		{
			name:   "TestGetOrCreateIndexWithCustomClient",
			client: customClient,
			args: args{
				config: IndexConfig{
					Uid: "TestBasicGetOrCreateIndex",
				},
			},
			wantResp: &Index{
				UID: "TestBasicGetOrCreateIndex",
			},
		},
		{
			name:   "TestGetOrCreateIndexWithPrimaryKey",
			client: defaultClient,
			args: args{
				config: IndexConfig{
					Uid:        "TestGetOrCreateIndexWithPrimaryKey",
					PrimaryKey: "PrimaryKey",
				},
			},
			wantResp: &Index{
				UID:        "TestGetOrCreateIndexWithPrimaryKey",
				PrimaryKey: "PrimaryKey",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client

			gotResp, err := c.GetOrCreateIndex(&tt.args.config)
			require.NoError(t, err)
			if assert.NotNil(t, gotResp) {
				require.Equal(t, tt.wantResp.UID, gotResp.UID)
				require.Equal(t, tt.wantResp.PrimaryKey, gotResp.PrimaryKey)
			}

			deleteAllIndexes(c)
		})
	}
}

func TestClient_Index(t *testing.T) {
	type args struct {
		uid string
	}
	tests := []struct {
		name   string
		client *Client
		args   args
		want   Index
	}{
		{
			name:   "TestBasicIndex",
			client: defaultClient,
			args: args{
				uid: "1",
			},
			want: Index{
				UID: "1",
			},
		},
		{
			name:   "TestIndexWithCustomClient",
			client: customClient,
			args: args{
				uid: "1",
			},
			want: Index{
				UID: "1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.client.Index(tt.args.uid)
			require.NotNil(t, got)
			require.Equal(t, tt.want.UID, got.UID)
		})
	}
}
