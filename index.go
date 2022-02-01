package meilisearch

import (
	"net/http"
	"strconv"
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
	UpdateIndex(primaryKey string) (resp *Task, err error)
	Delete(uid string) (ok bool, err error)
	GetStats() (resp *StatsIndex, err error)

	AddDocuments(documentsPtr interface{}, primaryKey ...string) (resp *Task, err error)
	AddDocumentsInBatches(documentsPtr interface{}, batchSize int, primaryKey ...string) (resp []Task, err error)
	AddDocumentsCsv(documents []byte, primaryKey ...string) (resp *Task, err error)
	AddDocumentsCsvInBatches(documents []byte, batchSize int, primaryKey ...string) (resp []Task, err error)
	AddDocumentsNdjson(documents []byte, primaryKey ...string) (resp *Task, err error)
	AddDocumentsNdjsonInBatches(documents []byte, batchSize int, primaryKey ...string) (resp []Task, err error)
	UpdateDocuments(documentsPtr interface{}, primaryKey ...string) (resp *Task, err error)
	GetDocument(uid string, documentPtr interface{}) error
	GetDocuments(request *DocumentsRequest, resp interface{}) error
	DeleteDocument(uid string) (resp *Task, err error)
	DeleteDocuments(uid []string) (resp *Task, err error)
	DeleteAllDocuments() (resp *Task, err error)
	Search(query string, request *SearchRequest) (*SearchResponse, error)

	GetTask(taskID int64) (resp *Task, err error)
	GetTasks() (resp *ResultTask, err error)

	GetSettings() (resp *Settings, err error)
	UpdateSettings(request *Settings) (resp *Task, err error)
	ResetSettings() (resp *Task, err error)
	GetRankingRules() (resp *[]string, err error)
	UpdateRankingRules(request *[]string) (resp *Task, err error)
	ResetRankingRules() (resp *Task, err error)
	GetDistinctAttribute() (resp *string, err error)
	UpdateDistinctAttribute(request string) (resp *Task, err error)
	ResetDistinctAttribute() (resp *Task, err error)
	GetSearchableAttributes() (resp *[]string, err error)
	UpdateSearchableAttributes(request *[]string) (resp *Task, err error)
	ResetSearchableAttributes() (resp *Task, err error)
	GetDisplayedAttributes() (resp *[]string, err error)
	UpdateDisplayedAttributes(request *[]string) (resp *Task, err error)
	ResetDisplayedAttributes() (resp *Task, err error)
	GetStopWords() (resp *[]string, err error)
	UpdateStopWords(request *[]string) (resp *Task, err error)
	ResetStopWords() (resp *Task, err error)
	GetSynonyms() (resp *map[string][]string, err error)
	UpdateSynonyms(request *map[string][]string) (resp *Task, err error)
	ResetSynonyms() (resp *Task, err error)
	GetFilterableAttributes() (resp *[]string, err error)
	UpdateFilterableAttributes(request *[]string) (resp *Task, err error)
	ResetFilterableAttributes() (resp *Task, err error)

	WaitForTask(task *Task, options ...waitParams) (*Task, error)
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

func (i Index) UpdateIndex(primaryKey string) (resp *Task, err error) {
	request := &UpdateIndexRequest{
		PrimaryKey: primaryKey,
	}
	i.PrimaryKey = primaryKey //nolint:golint,staticcheck
	resp = &Task{}

	req := internalRequest{
		endpoint:            "/indexes/" + i.UID,
		method:              http.MethodPut,
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
	resp := &Task{}
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

func (i Index) GetTask(taskID int64) (resp *Task, err error) {
	resp = &Task{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/tasks/" + strconv.FormatInt(taskID, 10),
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetTask",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) GetTasks() (resp *ResultTask, err error) {
	resp = &ResultTask{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/tasks",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetTasks",
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
func (i Index) WaitForTask(task *Task, options ...waitParams) (*Task, error) {
	return i.client.WaitForTask(task, options...)
}
