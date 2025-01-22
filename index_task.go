package meilisearch

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (i *index) GetTask(taskUID int64) (*Task, error) {
	return i.GetTaskWithContext(context.Background(), taskUID)
}

func (i *index) GetTaskWithContext(ctx context.Context, taskUID int64) (*Task, error) {
	return getTask(ctx, i.client, taskUID)
}

func (i *index) GetTasks(param *TasksQuery) (*TaskResult, error) {
	return i.GetTasksWithContext(context.Background(), param)
}

func (i *index) GetTasksWithContext(ctx context.Context, param *TasksQuery) (*TaskResult, error) {
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
	if param != nil {
		if param.Limit != 0 {
			req.withQueryParams["limit"] = strconv.FormatInt(param.Limit, 10)
		}
		if param.From != 0 {
			req.withQueryParams["from"] = strconv.FormatInt(param.From, 10)
		}
		if len(param.Statuses) != 0 {
			statuses := make([]string, len(param.Statuses))
			for i, status := range param.Statuses {
				statuses[i] = string(status)
			}
			req.withQueryParams["statuses"] = strings.Join(statuses, ",")
		}

		if len(param.Types) != 0 {
			types := make([]string, len(param.Types))
			for i, t := range param.Types {
				types[i] = string(t)
			}
			req.withQueryParams["types"] = strings.Join(types, ",")
		}
		if len(param.IndexUIDS) != 0 {
			param.IndexUIDS = append(param.IndexUIDS, i.uid)
			req.withQueryParams["indexUids"] = strings.Join(param.IndexUIDS, ",")
		} else {
			req.withQueryParams["indexUids"] = i.uid
		}

		if param.Reverse {
			req.withQueryParams["reverse"] = "true"
		}
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) WaitForTask(taskUID int64, interval time.Duration) (*Task, error) {
	return waitForTask(context.Background(), i.client, taskUID, interval)
}

func (i *index) WaitForTaskWithContext(ctx context.Context, taskUID int64, interval time.Duration) (*Task, error) {
	return waitForTask(ctx, i.client, taskUID, interval)
}
