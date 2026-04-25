package meilisearch

import (
	"context"
	"fmt"
	"net/http"
)

func (m *meilisearch) Delete(uid string) error {
	return m.DeleteWithContext(context.Background(), uid)
}

func (m *meilisearch) DeleteWithContext(ctx context.Context, uid string) error {
	req := &internalRequest{
		endpoint:            fmt.Sprintf("/dynamic-search-rules/%s", uid),
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        nil,
		acceptedStatusCodes: []int{http.StatusNoContent},
		functionName:        "DeleteSearchRule",
	}
	return m.client.executeRequest(ctx, req)
}
