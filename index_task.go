package meilisearch

import (
	"context"
	"net/http"
	"time"
)

func (i *index) GetTask(ctx context.Context, taskUID int64) (*Task, error) {
	return getTask(ctx, i.client, taskUID)
}

func (i *index) GetTasks(ctx context.Context, param *TasksQuery) (*TaskResult, error) {
	resp := new(TaskResult)
	req := &internalRequest{
		endpoint:            "/tasks",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		withQueryParams:     map[string]string{},
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetTasks",
	}

	// Work on a copy so the caller's query is never modified, and always
	// scope the request to this index, even when no query is given.
	query := TasksQuery{}
	if param != nil {
		query = *param
	}
	indexUIDs := make([]string, 0, len(query.IndexUIDS)+1)
	indexUIDs = append(indexUIDs, query.IndexUIDS...)
	query.IndexUIDS = append(indexUIDs, i.uid)
	encodeTasksQuery(&query, req)

	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) WaitForTask(ctx context.Context, taskUID int64, interval time.Duration) (*Task, error) {
	return waitForTask(ctx, i.client, taskUID, interval)
}
