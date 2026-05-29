package meilisearch

import (
	"context"
	"net/http"
	"strconv"
)

func (m *meilisearch) GetTaskDocuments(taskUID int64, dst interface{}) error {
	return m.GetTaskDocumentsWithContext(context.Background(), taskUID, dst)
}

func (m *meilisearch) GetTaskDocumentsWithContext(ctx context.Context, taskUID int64, dst interface{}) error {
	req := &internalRequest{
		endpoint:             "/tasks/" + strconv.FormatInt(taskUID, 10) + "/documents",
		method:               http.MethodGet,
		withRequest:          nil,
		withResponse:         dst,
		withQueryParams:      nil,
		withResponseEncoding: true,
		acceptedStatusCodes:  []int{http.StatusOK},
		acceptedContentType:  contentTypeNDJSON,
		functionName:         "GetTaskDocuments",
	}
	return m.client.executeRequest(ctx, req)
}
