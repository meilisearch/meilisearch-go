package meilisearch

import (
	"fmt"
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
