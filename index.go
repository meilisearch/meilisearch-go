package meilisearch

import (
	"context"
	"net/http"
)

// index is the type that represent an index in meilisearch
type index struct {
	uid        string
	primaryKey string
	client     *client
}

func newIndex(cli *client, uid string) IndexManager {
	return &index{
		client: cli,
		uid:    uid,
	}
}

func (i *index) GetTaskReader() TaskReader {
	return i
}

func (i *index) GetDocumentManager() DocumentManager {
	return i
}

func (i *index) GetDocumentReader() DocumentReader {
	return i
}

func (i *index) GetSettingsManager() SettingsManager {
	return i
}

func (i *index) GetSettingsReader() SettingsReader {
	return i
}

func (i *index) GetSearch() SearchReader {
	return i
}

func (i *index) GetIndexReader() IndexReader {
	return i
}

func (i *index) FetchInfo() (*IndexResult, error) {
	return i.FetchInfoWithContext(context.Background())
}

func (i *index) FetchInfoWithContext(ctx context.Context) (*IndexResult, error) {
	resp := new(IndexResult)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid,
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "FetchInfo",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	if resp.PrimaryKey != "" {
		i.primaryKey = resp.PrimaryKey
	}
	resp.IndexManager = i
	return resp, nil
}

func (i *index) FetchPrimaryKey() (*string, error) {
	return i.FetchPrimaryKeyWithContext(context.Background())
}

func (i *index) FetchPrimaryKeyWithContext(ctx context.Context) (*string, error) {
	idx, err := i.FetchInfoWithContext(ctx)
	if err != nil {
		return nil, err
	}
	i.primaryKey = idx.PrimaryKey
	return &idx.PrimaryKey, nil
}

func (i *index) UpdateIndex(params *UpdateIndexRequestParams) (*TaskInfo, error) {
	return i.UpdateIndexWithContext(context.Background(), params)
}

func (i *index) UpdateIndexWithContext(ctx context.Context, params *UpdateIndexRequestParams) (*TaskInfo, error) {
	resp := new(TaskInfo)

	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid,
		method:              http.MethodPatch,
		contentType:         contentTypeJSON,
		withRequest:         params,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateIndex",
	}

	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}

	if params.PrimaryKey != "" {
		i.primaryKey = params.PrimaryKey
	}

	if params.UID != "" {
		i.uid = params.UID
	}

	return resp, nil
}

func (i *index) Delete(uid string) (bool, error) {
	return i.DeleteWithContext(context.Background(), uid)
}

func (i *index) DeleteWithContext(ctx context.Context, uid string) (bool, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + uid,
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "Delete",
	}
	// err is not nil if status code is not 204 StatusNoContent
	if err := i.client.executeRequest(ctx, req); err != nil {
		return false, err
	}
	i.primaryKey = ""
	return true, nil
}

func (i *index) GetStats() (*StatsIndex, error) {
	return i.GetStatsWithContext(context.Background())
}

func (i *index) GetStatsWithContext(ctx context.Context) (*StatsIndex, error) {
	resp := new(StatsIndex)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/stats",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetStats",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) Compact() (*TaskInfo, error) {
	return i.CompactWithContext(context.Background())
}

func (i *index) CompactWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/compact",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "Compact",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}
