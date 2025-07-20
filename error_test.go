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
					MeilisearchApiError: meilisearchApiError{
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

type mockEncoder struct{}

func (m *mockEncoder) Encode(r io.Reader) (io.ReadCloser, error) {
	return nil, nil
}

func (m *mockEncoder) Decode(data []byte, v interface{}) error {
	msg, ok := v.(*meilisearchApiError)
	if !ok {
		return fmt.Errorf("wrong type")
	}
	msg.Message = "mocked message"
	msg.Code = "mocked code"
	msg.Type = "mocked type"
	msg.Link = "mocked link"
	return nil
}

type failEncoder struct{}

func (f *failEncoder) Encode(r io.Reader) (io.ReadCloser, error) {
	return nil, nil
}
func (f *failEncoder) Decode(_ []byte, v interface{}) error {
	return fmt.Errorf("decode failed")
}

func TestError_ErrorBody_WithEncoder(t *testing.T) {

	err := &Error{
		encoder: &mockEncoder{},
	}
	body := []byte(`{"message":"should not be used"}`)
	err.ErrorBody(body)
	require.Equal(t, "mocked message", err.MeilisearchApiError.Message)
	require.Equal(t, "mocked code", err.MeilisearchApiError.Code)
	require.Equal(t, "mocked type", err.MeilisearchApiError.Type)
	require.Equal(t, "mocked link", err.MeilisearchApiError.Link)

	err2 := &Error{
		encoder: &failEncoder{},
	}
	body2 := []byte(`{"message":"should not be used"}`)
	err2.ErrorBody(body2)
	// Should not set MeilisearchApiError fields
	require.Empty(t, err2.MeilisearchApiError.Message)
	require.Empty(t, err2.MeilisearchApiError.Code)
	require.Empty(t, err2.MeilisearchApiError.Type)
	require.Empty(t, err2.MeilisearchApiError.Link)
}
