package meilisearch

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError_VersionErrorHintMessage(t *testing.T) {
	type args struct {
		request     *internalRequest
		mockedError error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "VersionErrorHintMessageGetDocument",
			args: args{
				request: &internalRequest{
					functionName: "GetDocuments",
				},
				mockedError: &Error{
					Endpoint:         "endpointForGetDocuments",
					Method:           http.MethodPost,
					Function:         "GetDocuments",
					RequestToString:  "empty request",
					ResponseToString: "empty response",
					APIError: meilisearchApiError{
						Message: "empty Meilisearch message",
					},
					StatusCode: 1,
					rawMessage: "empty message",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VersionErrorHintMessage(tt.args.mockedError, tt.args.request)
			require.Error(t, err)
			fmt.Println(err)
			require.Equal(t, tt.args.mockedError.Error()+". Hint: It might not be working because you're not up to date with the Meilisearch version that "+tt.args.request.functionName+" call requires", err.Error())
		})
	}
}

type failEncoder struct{}

func (f *failEncoder) Encode(r io.Reader) (io.ReadCloser, error) {
	return nil, nil
}
func (f *failEncoder) Decode(_ []byte, v interface{}) error {
	return fmt.Errorf("decode failed")
}

func (f *failEncoder) Decoder(r io.Reader) (streamDecoder, error) {
	return nil, fmt.Errorf("decoder failed")
}

func TestError_ErrorBody_WithEncoder(t *testing.T) {

	err := &Error{
		encoder: &mockEncoder{},
	}
	body := []byte(`{"message":"should not be used"}`)
	err.ErrorBody(body)
	require.Equal(t, "mocked message", err.APIError.Message)
	require.Equal(t, APIErrCode("mocked code"), err.APIError.Code)
	require.Equal(t, "mocked type", err.APIError.Type)
	require.Equal(t, "mocked link", err.APIError.Link)

	err2 := &Error{
		encoder: &failEncoder{},
	}
	body2 := []byte(`{"message":"should not be used"}`)
	err2.ErrorBody(body2)
	// Should not set APIError fields
	require.Empty(t, err2.APIError.Message)
	require.Empty(t, err2.APIError.Code)
	require.Empty(t, err2.APIError.Type)
	require.Empty(t, err2.APIError.Link)
}

func TestError_UnwrapAndHasCode(t *testing.T) {
	origin := fmt.Errorf("underlying database error")
	err := &Error{
		Endpoint: "endpoint",
		Method:   "GET",
		Function: "Test",
		APIError: APIErrorDetails{
			Message: "Index movies not found.",
			Code:    APIErrCodeIndexNotFound,
			Type:    "invalid_request",
			Link:    "https://docs.meilisearch.com/errors#index_not_found",
		},
		OriginError: origin,
	}

	// Test Unwrap
	require.Equal(t, origin, err.Unwrap())

	// Test HasCode
	require.True(t, err.HasCode(APIErrCodeIndexNotFound))
	require.False(t, err.HasCode(APIErrCodeAPIKeyNotFound))
}
