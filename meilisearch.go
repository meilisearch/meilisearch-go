package meilisearch

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strconv"
	"time"
)

type meilisearch struct {
	client *client
}

type ServiceManager interface {
	// Index retrieves an IndexManager for a specific index.
	Index(uid string) IndexManager

	// GetIndex fetches the details of a specific index.
	GetIndex(indexID string) (*IndexResult, error)

	// GetIndexWithContext fetches the details of a specific index with a context for cancellation.
	GetIndexWithContext(ctx context.Context, indexID string) (*IndexResult, error)

	// GetRawIndex fetches the raw JSON representation of a specific index.
	GetRawIndex(uid string) (map[string]interface{}, error)

	// GetRawIndexWithContext fetches the raw JSON representation of a specific index with a context for cancellation.
	GetRawIndexWithContext(ctx context.Context, uid string) (map[string]interface{}, error)

	// ListIndexes lists all indexes.
	ListIndexes(param *IndexesQuery) (*IndexesResults, error)

	// ListIndexesWithContext lists all indexes with a context for cancellation.
	ListIndexesWithContext(ctx context.Context, param *IndexesQuery) (*IndexesResults, error)

	// GetRawIndexes fetches the raw JSON representation of all indexes.
	GetRawIndexes(param *IndexesQuery) (map[string]interface{}, error)

	// GetRawIndexesWithContext fetches the raw JSON representation of all indexes with a context for cancellation.
	GetRawIndexesWithContext(ctx context.Context, param *IndexesQuery) (map[string]interface{}, error)

	// CreateIndex creates a new index.
	CreateIndex(config *IndexConfig) (*TaskInfo, error)

	// CreateIndexWithContext creates a new index with a context for cancellation.
	CreateIndexWithContext(ctx context.Context, config *IndexConfig) (*TaskInfo, error)

	// DeleteIndex deletes a specific index.
	DeleteIndex(uid string) (*TaskInfo, error)

	// DeleteIndexWithContext deletes a specific index with a context for cancellation.
	DeleteIndexWithContext(ctx context.Context, uid string) (*TaskInfo, error)

	// MultiSearch performs a multi-index search.
	MultiSearch(queries *MultiSearchRequest) (*MultiSearchResponse, error)

	// MultiSearchWithContext performs a multi-index search with a context for cancellation.
	MultiSearchWithContext(ctx context.Context, queries *MultiSearchRequest) (*MultiSearchResponse, error)

	// CreateKey creates a new API key.
	CreateKey(request *Key) (*Key, error)

	// CreateKeyWithContext creates a new API key with a context for cancellation.
	CreateKeyWithContext(ctx context.Context, request *Key) (*Key, error)

	// GetKey fetches the details of a specific API key.
	GetKey(identifier string) (*Key, error)

	// GetKeyWithContext fetches the details of a specific API key with a context for cancellation.
	GetKeyWithContext(ctx context.Context, identifier string) (*Key, error)

	// GetKeys lists all API keys.
	GetKeys(param *KeysQuery) (*KeysResults, error)

	// GetKeysWithContext lists all API keys with a context for cancellation.
	GetKeysWithContext(ctx context.Context, param *KeysQuery) (*KeysResults, error)

	// UpdateKey updates a specific API key.
	UpdateKey(keyOrUID string, request *Key) (*Key, error)

	// UpdateKeyWithContext updates a specific API key with a context for cancellation.
	UpdateKeyWithContext(ctx context.Context, keyOrUID string, request *Key) (*Key, error)

	// DeleteKey deletes a specific API key.
	DeleteKey(keyOrUID string) (bool, error)

	// DeleteKeyWithContext deletes a specific API key with a context for cancellation.
	DeleteKeyWithContext(ctx context.Context, keyOrUID string) (bool, error)

	// GetTask fetches the details of a specific task.
	GetTask(taskUID int64) (*Task, error)

	// GetTaskWithContext fetches the details of a specific task with a context for cancellation.
	GetTaskWithContext(ctx context.Context, taskUID int64) (*Task, error)

	// GetTasks lists all tasks.
	GetTasks(param *TasksQuery) (*TaskResult, error)

	// GetTasksWithContext lists all tasks with a context for cancellation.
	GetTasksWithContext(ctx context.Context, param *TasksQuery) (*TaskResult, error)

	// CancelTasks cancels specific tasks.
	CancelTasks(param *CancelTasksQuery) (*TaskInfo, error)

	// CancelTasksWithContext cancels specific tasks with a context for cancellation.
	CancelTasksWithContext(ctx context.Context, param *CancelTasksQuery) (*TaskInfo, error)

	// DeleteTasks deletes specific tasks.
	DeleteTasks(param *DeleteTasksQuery) (*TaskInfo, error)

	// DeleteTasksWithContext deletes specific tasks with a context for cancellation.
	DeleteTasksWithContext(ctx context.Context, param *DeleteTasksQuery) (*TaskInfo, error)

	// WaitForTask waits for a specific task to complete.
	WaitForTask(taskUID int64, interval time.Duration) (*Task, error)

	// WaitForTaskWithContext waits for a specific task to complete with a context for cancellation.
	WaitForTaskWithContext(ctx context.Context, taskUID int64, interval time.Duration) (*Task, error)

	// SwapIndexes swaps the positions of two indexes.
	SwapIndexes(param []*SwapIndexesParams) (*TaskInfo, error)

	// SwapIndexesWithContext swaps the positions of two indexes with a context for cancellation.
	SwapIndexesWithContext(ctx context.Context, param []*SwapIndexesParams) (*TaskInfo, error)

	// GenerateTenantToken generates a tenant token for multi-tenancy.
	GenerateTenantToken(apiKeyUID string, searchRules map[string]interface{}, options *TenantTokenOptions) (string, error)

	// GetStats fetches global stats.
	GetStats() (*Stats, error)

	// GetStatsWithContext fetches global stats with a context for cancellation.
	GetStatsWithContext(ctx context.Context) (*Stats, error)

	// CreateDump creates a database dump.
	CreateDump() (*TaskInfo, error)

	// CreateDumpWithContext creates a database dump with a context for cancellation.
	CreateDumpWithContext(ctx context.Context) (*TaskInfo, error)

	// Version fetches the version of the Meilisearch server.
	Version() (*Version, error)

	// VersionWithContext fetches the version of the Meilisearch server with a context for cancellation.
	VersionWithContext(ctx context.Context) (*Version, error)

	// Health checks the health of the Meilisearch server.
	Health() (*Health, error)

	// HealthWithContext checks the health of the Meilisearch server with a context for cancellation.
	HealthWithContext(ctx context.Context) (*Health, error)

	// IsHealthy checks if the Meilisearch server is healthy.
	IsHealthy() bool

	// CreateSnapshot create database snapshot from meilisearch
	CreateSnapshot() (*TaskInfo, error)

	// CreateSnapshotWithContext create database snapshot from meilisearch and support parent context
	CreateSnapshotWithContext(ctx context.Context) (*TaskInfo, error)

	// Close closes the connection to the Meilisearch server.
	Close()
}

// New create new service manager for operating on meilisearch
func New(host string, options ...Option) ServiceManager {
	defOpt := defaultMeiliOpt

	for _, opt := range options {
		opt(defOpt)
	}

	return &meilisearch{
		client: newClient(
			defOpt.client,
			host,
			defOpt.apiKey,
		),
	}
}

// Connect create service manager and check connection with meilisearch
func Connect(host string, options ...Option) (ServiceManager, error) {
	meili := New(host, options...)

	if !meili.IsHealthy() {
		return nil, ErrConnectingFailed
	}

	return meili, nil
}

func (m *meilisearch) Index(uid string) IndexManager {
	return newIndex(m.client, uid)
}

func (m *meilisearch) GetIndex(indexID string) (*IndexResult, error) {
	return m.GetIndexWithContext(context.Background(), indexID)
}

func (m *meilisearch) GetIndexWithContext(ctx context.Context, indexID string) (*IndexResult, error) {
	return newIndex(m.client, indexID).FetchInfoWithContext(ctx)
}

func (m *meilisearch) GetRawIndex(uid string) (map[string]interface{}, error) {
	return m.GetRawIndexWithContext(context.Background(), uid)
}

func (m *meilisearch) GetRawIndexWithContext(ctx context.Context, uid string) (map[string]interface{}, error) {
	resp := map[string]interface{}{}
	req := &internalRequest{
		endpoint:            "/indexes/" + uid,
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetRawIndex",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) ListIndexes(param *IndexesQuery) (*IndexesResults, error) {
	return m.ListIndexesWithContext(context.Background(), param)
}

func (m *meilisearch) ListIndexesWithContext(ctx context.Context, param *IndexesQuery) (*IndexesResults, error) {
	resp := new(IndexesResults)
	req := &internalRequest{
		endpoint:            "/indexes",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		withQueryParams:     map[string]string{},
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetIndexes",
	}
	if param != nil && param.Limit != 0 {
		req.withQueryParams["limit"] = strconv.FormatInt(param.Limit, 10)
	}
	if param != nil && param.Offset != 0 {
		req.withQueryParams["offset"] = strconv.FormatInt(param.Offset, 10)
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}

	for i := range resp.Results {
		resp.Results[i].IndexManager = newIndex(m.client, resp.Results[i].UID)
	}

	return resp, nil
}

func (m *meilisearch) GetRawIndexes(param *IndexesQuery) (map[string]interface{}, error) {
	return m.GetRawIndexesWithContext(context.Background(), param)
}

func (m *meilisearch) GetRawIndexesWithContext(ctx context.Context, param *IndexesQuery) (map[string]interface{}, error) {
	resp := map[string]interface{}{}
	req := &internalRequest{
		endpoint:            "/indexes",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		withQueryParams:     map[string]string{},
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetRawIndexes",
	}
	if param != nil && param.Limit != 0 {
		req.withQueryParams["limit"] = strconv.FormatInt(param.Limit, 10)
	}
	if param != nil && param.Offset != 0 {
		req.withQueryParams["offset"] = strconv.FormatInt(param.Offset, 10)
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) CreateIndex(config *IndexConfig) (*TaskInfo, error) {
	return m.CreateIndexWithContext(context.Background(), config)
}

func (m *meilisearch) CreateIndexWithContext(ctx context.Context, config *IndexConfig) (*TaskInfo, error) {
	request := &CreateIndexRequest{
		UID:        config.Uid,
		PrimaryKey: config.PrimaryKey,
	}
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "CreateIndex",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) DeleteIndex(uid string) (*TaskInfo, error) {
	return m.DeleteIndexWithContext(context.Background(), uid)
}

func (m *meilisearch) DeleteIndexWithContext(ctx context.Context, uid string) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + uid,
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "DeleteIndex",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) MultiSearch(queries *MultiSearchRequest) (*MultiSearchResponse, error) {
	return m.MultiSearchWithContext(context.Background(), queries)
}

func (m *meilisearch) MultiSearchWithContext(ctx context.Context, queries *MultiSearchRequest) (*MultiSearchResponse, error) {
	resp := new(MultiSearchResponse)

	for i := 0; i < len(queries.Queries); i++ {
		queries.Queries[i].validate()
	}

	req := &internalRequest{
		endpoint:            "/multi-search",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         queries,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "MultiSearch",
	}

	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *meilisearch) CreateKey(request *Key) (*Key, error) {
	return m.CreateKeyWithContext(context.Background(), request)
}

func (m *meilisearch) CreateKeyWithContext(ctx context.Context, request *Key) (*Key, error) {
	parsedRequest := convertKeyToParsedKey(*request)
	resp := new(Key)
	req := &internalRequest{
		endpoint:            "/keys",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         &parsedRequest,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusCreated},
		functionName:        "CreateKey",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) GetKey(identifier string) (*Key, error) {
	return m.GetKeyWithContext(context.Background(), identifier)
}

func (m *meilisearch) GetKeyWithContext(ctx context.Context, identifier string) (*Key, error) {
	resp := new(Key)
	req := &internalRequest{
		endpoint:            "/keys/" + identifier,
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetKey",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) GetKeys(param *KeysQuery) (*KeysResults, error) {
	return m.GetKeysWithContext(context.Background(), param)
}

func (m *meilisearch) GetKeysWithContext(ctx context.Context, param *KeysQuery) (*KeysResults, error) {
	resp := new(KeysResults)
	req := &internalRequest{
		endpoint:            "/keys",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		withQueryParams:     map[string]string{},
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetKeys",
	}
	if param != nil && param.Limit != 0 {
		req.withQueryParams["limit"] = strconv.FormatInt(param.Limit, 10)
	}
	if param != nil && param.Offset != 0 {
		req.withQueryParams["offset"] = strconv.FormatInt(param.Offset, 10)
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) UpdateKey(keyOrUID string, request *Key) (*Key, error) {
	return m.UpdateKeyWithContext(context.Background(), keyOrUID, request)
}

func (m *meilisearch) UpdateKeyWithContext(ctx context.Context, keyOrUID string, request *Key) (*Key, error) {
	parsedRequest := KeyUpdate{Name: request.Name, Description: request.Description}
	resp := new(Key)
	req := &internalRequest{
		endpoint:            "/keys/" + keyOrUID,
		method:              http.MethodPatch,
		contentType:         contentTypeJSON,
		withRequest:         &parsedRequest,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "UpdateKey",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) DeleteKey(keyOrUID string) (bool, error) {
	return m.DeleteKeyWithContext(context.Background(), keyOrUID)
}

func (m *meilisearch) DeleteKeyWithContext(ctx context.Context, keyOrUID string) (bool, error) {
	req := &internalRequest{
		endpoint:            "/keys/" + keyOrUID,
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        nil,
		acceptedStatusCodes: []int{http.StatusNoContent},
		functionName:        "DeleteKey",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return false, err
	}
	return true, nil
}

func (m *meilisearch) GetTask(taskUID int64) (*Task, error) {
	return m.GetTaskWithContext(context.Background(), taskUID)
}

func (m *meilisearch) GetTaskWithContext(ctx context.Context, taskUID int64) (*Task, error) {
	return getTask(ctx, m.client, taskUID)
}

func (m *meilisearch) GetTasks(param *TasksQuery) (*TaskResult, error) {
	return m.GetTasksWithContext(context.Background(), param)
}

func (m *meilisearch) GetTasksWithContext(ctx context.Context, param *TasksQuery) (*TaskResult, error) {
	resp := new(TaskResult)
	req := &internalRequest{
		endpoint:            "/tasks",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		withQueryParams:     map[string]string{},
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetTasks",
	}
	if param != nil {
		encodeTasksQuery(param, req)
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) CancelTasks(param *CancelTasksQuery) (*TaskInfo, error) {
	return m.CancelTasksWithContext(context.Background(), param)
}

func (m *meilisearch) CancelTasksWithContext(ctx context.Context, param *CancelTasksQuery) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/tasks/cancel",
		method:              http.MethodPost,
		withRequest:         nil,
		withResponse:        &resp,
		withQueryParams:     map[string]string{},
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "CancelTasks",
	}
	if param != nil {
		paramToSend := &TasksQuery{
			UIDS:             param.UIDS,
			IndexUIDS:        param.IndexUIDS,
			Statuses:         param.Statuses,
			Types:            param.Types,
			BeforeEnqueuedAt: param.BeforeEnqueuedAt,
			AfterEnqueuedAt:  param.AfterEnqueuedAt,
			BeforeStartedAt:  param.BeforeStartedAt,
			AfterStartedAt:   param.AfterStartedAt,
		}
		encodeTasksQuery(paramToSend, req)
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) DeleteTasks(param *DeleteTasksQuery) (*TaskInfo, error) {
	return m.DeleteTasksWithContext(context.Background(), param)
}

func (m *meilisearch) DeleteTasksWithContext(ctx context.Context, param *DeleteTasksQuery) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/tasks",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        &resp,
		withQueryParams:     map[string]string{},
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "DeleteTasks",
	}
	if param != nil {
		paramToSend := &TasksQuery{
			UIDS:             param.UIDS,
			IndexUIDS:        param.IndexUIDS,
			Statuses:         param.Statuses,
			Types:            param.Types,
			CanceledBy:       param.CanceledBy,
			BeforeEnqueuedAt: param.BeforeEnqueuedAt,
			AfterEnqueuedAt:  param.AfterEnqueuedAt,
			BeforeStartedAt:  param.BeforeStartedAt,
			AfterStartedAt:   param.AfterStartedAt,
			BeforeFinishedAt: param.BeforeFinishedAt,
			AfterFinishedAt:  param.AfterFinishedAt,
		}
		encodeTasksQuery(paramToSend, req)
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) SwapIndexes(param []*SwapIndexesParams) (*TaskInfo, error) {
	return m.SwapIndexesWithContext(context.Background(), param)
}

func (m *meilisearch) SwapIndexesWithContext(ctx context.Context, param []*SwapIndexesParams) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/swap-indexes",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         param,
		withResponse:        &resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "SwapIndexes",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) WaitForTask(taskUID int64, interval time.Duration) (*Task, error) {
	return waitForTask(context.Background(), m.client, taskUID, interval)
}

func (m *meilisearch) WaitForTaskWithContext(ctx context.Context, taskUID int64, interval time.Duration) (*Task, error) {
	return waitForTask(ctx, m.client, taskUID, interval)
}

func (m *meilisearch) GenerateTenantToken(
	apiKeyUID string,
	searchRules map[string]interface{},
	options *TenantTokenOptions,
) (string, error) {
	// validate the arguments
	if searchRules == nil {
		return "", fmt.Errorf("GenerateTenantToken: The search rules added in the token generation " +
			"must be of type array or object")
	}
	if (options == nil || options.APIKey == "") && m.client.apiKey == "" {
		return "", fmt.Errorf("GenerateTenantToken: The API key used for the token " +
			"generation must exist and be a valid meilisearch key")
	}
	if apiKeyUID == "" || !IsValidUUID(apiKeyUID) {
		return "", fmt.Errorf("GenerateTenantToken: The uid used for the token " +
			"generation must exist and comply to uuid4 format")
	}
	if options != nil && !options.ExpiresAt.IsZero() && options.ExpiresAt.Before(time.Now()) {
		return "", fmt.Errorf("GenerateTenantToken: When the expiresAt field in " +
			"the token generation has a value, it must be a date set in the future")
	}

	var secret string
	if options == nil || options.APIKey == "" {
		secret = m.client.apiKey
	} else {
		secret = options.APIKey
	}

	// For HMAC signing method, the key should be any []byte
	hmacSampleSecret := []byte(secret)

	// Create the claims
	claims := TenantTokenClaims{}
	if options != nil && !options.ExpiresAt.IsZero() {
		claims.RegisteredClaims = jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(options.ExpiresAt),
		}
	}
	claims.APIKeyUID = apiKeyUID
	claims.SearchRules = searchRules

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)

	return tokenString, err
}

func (m *meilisearch) GetStats() (*Stats, error) {
	return m.GetStatsWithContext(context.Background())
}

func (m *meilisearch) GetStatsWithContext(ctx context.Context) (*Stats, error) {
	resp := new(Stats)
	req := &internalRequest{
		endpoint:            "/stats",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetStats",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) CreateDump() (*TaskInfo, error) {
	return m.CreateDumpWithContext(context.Background())
}

func (m *meilisearch) CreateDumpWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/dumps",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "CreateDump",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) Version() (*Version, error) {
	return m.VersionWithContext(context.Background())
}

func (m *meilisearch) VersionWithContext(ctx context.Context) (*Version, error) {
	resp := new(Version)
	req := &internalRequest{
		endpoint:            "/version",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Version",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) Health() (*Health, error) {
	return m.HealthWithContext(context.Background())
}

func (m *meilisearch) HealthWithContext(ctx context.Context) (*Health, error) {
	resp := new(Health)
	req := &internalRequest{
		endpoint:            "/health",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Health",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) CreateSnapshot() (*TaskInfo, error) {
	return m.CreateSnapshotWithContext(context.Background())
}

func (m *meilisearch) CreateSnapshotWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/snapshots",
		method:              http.MethodPost,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		contentType:         contentTypeJSON,
		functionName:        "CreateSnapshot",
	}

	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) IsHealthy() bool {
	res, err := m.HealthWithContext(context.Background())
	return err == nil && res.Status == "available"
}

func (m *meilisearch) Close() {
	m.client.client.CloseIdleConnections()
}

func getTask(ctx context.Context, cli *client, taskUID int64) (*Task, error) {
	resp := new(Task)
	req := &internalRequest{
		endpoint:            "/tasks/" + strconv.FormatInt(taskUID, 10),
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetTask",
	}
	if err := cli.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func waitForTask(ctx context.Context, cli *client, taskUID int64, interval time.Duration) (*Task, error) {
	if interval == 0 {
		interval = 50 * time.Millisecond
	}

	// extract closure to get the task and check the status first before the ticker
	fn := func() (*Task, error) {
		getTask, err := getTask(ctx, cli, taskUID)
		if err != nil {
			return nil, err
		}

		if getTask.Status != TaskStatusEnqueued && getTask.Status != TaskStatusProcessing {
			return getTask, nil
		}
		return nil, nil
	}

	// run first before the ticker, we do not want to wait for the first interval
	task, err := fn()
	if err != nil {
		// Return error if it exists
		return nil, err
	}

	// Return task if it exists
	if task != nil {
		return task, nil
	}

	// Create a ticker to check the task status, because our initial check was not successful
	ticker := time.NewTicker(interval)

	// Defer the stop of the ticker, help GC to cleanup
	defer func() {
		// we might want to revist this, go.mod now is 1.16
		// however I still encouter the issue on go 1.22.2
		// there are 2 issues regarding tickers
		// https://go-review.googlesource.com/c/go/+/512355
		// https://github.com/golang/go/issues/61542
		ticker.Stop()
		ticker = nil
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			task, err := fn()
			if err != nil {
				return nil, err
			}

			if task != nil {
				return task, nil
			}
		}
	}
}
