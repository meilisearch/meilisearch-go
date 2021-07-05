package meilisearch

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

// IndexConfig configure the Index
type IndexConfig struct {

	// Host is the host of your meilisearch database
	// Example: 'http://localhost:7700'
	Uid string

	// PrimaryKey is optional
	PrimaryKey string

	client *Client
}

type IndexInterface interface {
	FetchInfo() (resp *Index, err error)
	FetchPrimaryKey() (resp *string, err error)
	UpdateIndex(primaryKey string) (resp *Index, err error)
	Delete(uid string) (ok bool, err error)
	GetStats() (resp *StatsIndex, err error)

	AddDocuments(documentsPtr interface{}) (resp *AsyncUpdateID, err error)
	UpdateDocuments(documentsPtr interface{}) (resp *AsyncUpdateID, err error)
	AddDocumentsWithPrimaryKey(documentsPtr interface{}, primaryKey string) (resp *AsyncUpdateID, err error)
	UpdateDocumentsWithPrimaryKey(documentsPtr interface{}, primaryKey string) (resp *AsyncUpdateID, err error)
	GetDocument(uid string, documentPtr interface{}) error
	GetDocuments(request *DocumentsRequest, resp interface{}) error
	DeleteDocument(uid string) (resp *AsyncUpdateID, err error)
	DeleteDocuments(uid []string) (resp *AsyncUpdateID, err error)
	DeleteAllDocuments() (resp *AsyncUpdateID, err error)
	Search(query string, request *SearchRequest) (*SearchResponse, error)

	GetUpdateStatus(updateID int64) (resp *Update, err error)
	GetAllUpdateStatus() (resp *[]Update, err error)

	GetSettings() (resp *Settings, err error)
	UpdateSettings(request *Settings) (resp *AsyncUpdateID, err error)
	ResetSettings() (resp *AsyncUpdateID, err error)
	GetRankingRules() (resp *[]string, err error)
	UpdateRankingRules(request *[]string) (resp *AsyncUpdateID, err error)
	ResetRankingRules() (resp *AsyncUpdateID, err error)
	GetDistinctAttribute() (resp *string, err error)
	UpdateDistinctAttribute(request string) (resp *AsyncUpdateID, err error)
	ResetDistinctAttribute() (resp *AsyncUpdateID, err error)
	GetSearchableAttributes() (resp *[]string, err error)
	UpdateSearchableAttributes(request *[]string) (resp *AsyncUpdateID, err error)
	ResetSearchableAttributes() (resp *AsyncUpdateID, err error)
	GetDisplayedAttributes() (resp *[]string, err error)
	UpdateDisplayedAttributes(request *[]string) (resp *AsyncUpdateID, err error)
	ResetDisplayedAttributes() (resp *AsyncUpdateID, err error)
	GetStopWords() (resp *[]string, err error)
	UpdateStopWords(request *[]string) (resp *AsyncUpdateID, err error)
	ResetStopWords() (resp *AsyncUpdateID, err error)
	GetSynonyms() (resp *map[string][]string, err error)
	UpdateSynonyms(request *map[string][]string) (resp *AsyncUpdateID, err error)
	ResetSynonyms() (resp *AsyncUpdateID, err error)
	GetFilterableAttributes() (resp *[]string, err error)
	UpdateFilterableAttributes(request *[]string) (resp *AsyncUpdateID, err error)
	ResetFilterableAttributes() (resp *AsyncUpdateID, err error)

	WaitForPendingUpdate(ctx context.Context, interval time.Duration, updateID *AsyncUpdateID) (UpdateStatus, error)
	DefaultWaitForPendingUpdate(updateID *AsyncUpdateID) (UpdateStatus, error)
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
	i.PrimaryKey = resp.PrimaryKey
	return resp, nil
}

func (i Index) FetchPrimaryKey() (resp *string, err error) {
	index, err := i.FetchInfo()
	if err != nil {
		return nil, err
	}
	return &index.PrimaryKey, nil
}

func (i Index) UpdateIndex(primaryKey string) (resp *Index, err error) {
	request := &UpdateIndexRequest{
		PrimaryKey: primaryKey,
	}
	i.PrimaryKey = primaryKey

	req := internalRequest{
		endpoint:            "/indexes/" + i.UID,
		method:              http.MethodPut,
		withRequest:         request,
		withResponse:        &i,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "UpdateIndex",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return &i, nil
}

func (i Index) Delete(uid string) (ok bool, err error) {
	req := internalRequest{
		endpoint:            "/indexes/" + uid,
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        nil,
		acceptedStatusCodes: []int{http.StatusNoContent},
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

func (i Index) GetUpdateStatus(updateID int64) (resp *Update, err error) {
	resp = &Update{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/updates/" + strconv.FormatInt(updateID, 10),
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetUpdateStatus",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) GetAllUpdateStatus() (resp *[]Update, err error) {
	resp = &[]Update{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/updates",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetAllUpdateStatus",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

// DefaultWaitForPendingUpdate checks each 50ms the status of a WaitForPendingUpdate.
// This is a default implementation of WaitForPendingUpdate.
func (i Index) DefaultWaitForPendingUpdate(updateID *AsyncUpdateID) (UpdateStatus, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	return i.WaitForPendingUpdate(ctx, time.Millisecond*50, updateID)
}

// WaitForPendingUpdate waits for the end of an update.
// The function will check by regular interval provided in parameter interval
// the UpdateStatus. If it is not UpdateStatusEnqueued or the ctx cancelled
// we return the UpdateStatus.
func (i Index) WaitForPendingUpdate(
	ctx context.Context,
	interval time.Duration,
	updateID *AsyncUpdateID) (UpdateStatus, error) {
	for {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		update, err := i.GetUpdateStatus(updateID.UpdateID)
		if err != nil {
			return UpdateStatusUnknown, nil
		}
		if update.Status != UpdateStatusEnqueued && update.Status != UpdateStatusProcessing {
			return update.Status, nil
		}
		time.Sleep(interval)
	}
}
