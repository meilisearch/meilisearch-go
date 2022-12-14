package meilisearch

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// IndexConfig configure the Index
type IndexConfig struct {

	// Uid is the unique identifier of a given index.
	Uid string

	// PrimaryKey is optional
	PrimaryKey string

	client *Client //nolint:golint,unused,structcheck
}

type IndexInterface interface {
	FetchInfo() (resp *Index, err error)
	FetchPrimaryKey() (resp *string, err error)
	UpdateIndex(primaryKey string) (resp *TaskInfo, err error)
	Delete(uid string) (ok bool, err error)
	GetStats() (resp *StatsIndex, err error)

	AddDocuments(documentsPtr interface{}, primaryKey ...string) (resp *TaskInfo, err error)
	AddDocumentsInBatches(documentsPtr interface{}, batchSize int, primaryKey ...string) (resp []TaskInfo, err error)
	AddDocumentsCsv(documents []byte, primaryKey ...string) (resp *TaskInfo, err error)
	AddDocumentsCsvInBatches(documents []byte, batchSize int, primaryKey ...string) (resp []TaskInfo, err error)
	AddDocumentsNdjson(documents []byte, primaryKey ...string) (resp *TaskInfo, err error)
	AddDocumentsNdjsonInBatches(documents []byte, batchSize int, primaryKey ...string) (resp []TaskInfo, err error)
	UpdateDocuments(documentsPtr interface{}, primaryKey ...string) (resp *TaskInfo, err error)
	GetDocument(uid string, request *DocumentQuery, documentPtr interface{}) error
	GetDocuments(param *DocumentsQuery, resp *DocumentsResult) error
	DeleteDocument(uid string) (resp *TaskInfo, err error)
	DeleteDocuments(uid []string) (resp *TaskInfo, err error)
	DeleteAllDocuments() (resp *TaskInfo, err error)
	Search(query string, request *SearchRequest) (*SearchResponse, error)
	SearchRaw(query string, request *SearchRequest) (*json.RawMessage, error)

	GetTask(taskUID int64) (resp *Task, err error)
	GetTasks(param *TasksQuery) (resp *TaskResult, err error)

	GetSettings() (resp *Settings, err error)
	UpdateSettings(request *Settings) (resp *TaskInfo, err error)
	ResetSettings() (resp *TaskInfo, err error)
	GetRankingRules() (resp *[]string, err error)
	UpdateRankingRules(request *[]string) (resp *TaskInfo, err error)
	ResetRankingRules() (resp *TaskInfo, err error)
	GetDistinctAttribute() (resp *string, err error)
	UpdateDistinctAttribute(request string) (resp *TaskInfo, err error)
	ResetDistinctAttribute() (resp *TaskInfo, err error)
	GetSearchableAttributes() (resp *[]string, err error)
	UpdateSearchableAttributes(request *[]string) (resp *TaskInfo, err error)
	ResetSearchableAttributes() (resp *TaskInfo, err error)
	GetDisplayedAttributes() (resp *[]string, err error)
	UpdateDisplayedAttributes(request *[]string) (resp *TaskInfo, err error)
	ResetDisplayedAttributes() (resp *TaskInfo, err error)
	GetStopWords() (resp *[]string, err error)
	UpdateStopWords(request *[]string) (resp *TaskInfo, err error)
	ResetStopWords() (resp *TaskInfo, err error)
	GetSynonyms() (resp *map[string][]string, err error)
	UpdateSynonyms(request *map[string][]string) (resp *TaskInfo, err error)
	ResetSynonyms() (resp *TaskInfo, err error)
	GetFilterableAttributes() (resp *[]string, err error)
	UpdateFilterableAttributes(request *[]string) (resp *TaskInfo, err error)
	ResetFilterableAttributes() (resp *TaskInfo, err error)

	WaitForTask(taskUID int64, options ...WaitParams) (*Task, error)
}

var _ IndexInterface = &Index{}

func newIndex(client *Client, uid string) *Index {
	return &Index{
		UID:    uid,
		client: client,
	}
}

func (i Index) FetchInfo() (resp *Index, err error) {
	resp = newIndex(i.client, i.UID)
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID,
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "FetchInfo",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	i.PrimaryKey = resp.PrimaryKey //nolint:golint,staticcheck
	return resp, nil
}

func (i Index) FetchPrimaryKey() (resp *string, err error) {
	index, err := i.FetchInfo()
	if err != nil {
		return nil, err
	}
	return &index.PrimaryKey, nil
}

func (i Index) UpdateIndex(primaryKey string) (resp *TaskInfo, err error) {
	request := &UpdateIndexRequest{
		PrimaryKey: primaryKey,
	}
	i.PrimaryKey = primaryKey //nolint:golint,staticcheck
	resp = &TaskInfo{}

	req := internalRequest{
		endpoint:            "/indexes/" + i.UID,
		method:              http.MethodPatch,
		contentType:         contentTypeJSON,
		withRequest:         request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateIndex",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) Delete(uid string) (ok bool, err error) {
	resp := &TaskInfo{}
	req := internalRequest{
		endpoint:            "/indexes/" + uid,
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "Delete",
	}
	// err is not nil if status code is not 204 StatusNoContent
	if err := i.client.executeRequest(req); err != nil {
		return false, err
	}
	return true, nil
}

func (i Index) GetStats() (resp *StatsIndex, err error) {
	resp = &StatsIndex{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/stats",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetStats",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) GetTask(taskUID int64) (resp *Task, err error) {
	return i.client.GetTask(taskUID)
}

func (i Index) GetTasks(param *TasksQuery) (resp *TaskResult, err error) {
	resp = &TaskResult{}
	req := internalRequest{
		endpoint:            "/tasks",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
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
			req.withQueryParams["statuses"] = strings.Join(param.Statuses, ",")
		}
		if len(param.Types) != 0 {
			req.withQueryParams["types"] = strings.Join(param.Types, ",")
		}
		if len(param.IndexUIDS) != 0 {
			param.IndexUIDS = append(param.IndexUIDS, i.UID)
			req.withQueryParams["indexUids"] = strings.Join(param.IndexUIDS, ",")
		} else {
			req.withQueryParams["indexUids"] = i.UID
		}
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

// WaitForTask waits for a task to be processed.
// The function will check by regular interval provided in parameter interval
// the TaskStatus.
// If no ctx and interval are provided WaitForTask will check each 50ms the
// status of a task.
func (i Index) WaitForTask(taskUID int64, options ...WaitParams) (*Task, error) {
	return i.client.WaitForTask(taskUID, options...)
}
