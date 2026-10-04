package meilisearch

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type meilisearch struct {
	client *client
}

// New create new service manager for operating on meilisearch
func New(host string, options ...Option) ServiceManager {
	opts := _defaultOpts()

	for _, opt := range options {
		opt(opts)
	}

	return &meilisearch{
		client: newClient(
			opts.client,
			host,
			opts.apiKey,
			&clientConfig{
				contentEncoding:          opts.contentEncoding.encodingType,
				encodingCompressionLevel: opts.contentEncoding.level,
				disableRetry:             opts.disableRetry,
				retryOnStatus:            opts.retryOnStatus,
				maxRetries:               opts.maxRetries,
				jsonMarshal:              opts.jsonMarshaler,
				jsonUnmarshal:            opts.jsonUnmarshaler,
			},
		),
	}
}

// Connect create service manager and check connection with meilisearch
func Connect(host string, options ...Option) (ServiceManager, error) {
	meili := New(host, options...)

	resp, err := meili.Health(context.Background())

	if err != nil {
		return nil, err
	}

	if resp.Status != "available" {
		return nil, ErrMeilisearchNotAvailable
	}

	return meili, nil
}

func (m *meilisearch) ServiceReader() ServiceReader {
	return m
}

func (m *meilisearch) TaskManager() TaskManager {
	return m
}

func (m *meilisearch) TaskReader() TaskReader {
	return m
}

func (m *meilisearch) KeyManager() KeyManager {
	return m
}

func (m *meilisearch) KeyReader() KeyReader {
	return m
}

func (m *meilisearch) WebhookManager() WebhookManager { return m }

func (m *meilisearch) WebhookReader() WebhookReader { return m }

func (m *meilisearch) SearchRulesManager() SearchRulesManager { return m }

func (m *meilisearch) SearchRulesReader() SearchRulesReader { return m }

func (m *meilisearch) ChatManager() ChatManager { return m }

func (m *meilisearch) ChatReader() ChatReader { return m }

func (m *meilisearch) Index(uid string) IndexManager {
	return newIndex(m.client, uid)
}

func (m *meilisearch) GetIndex(ctx context.Context, indexID string) (*IndexResult, error) {
	return newIndex(m.client, indexID).FetchInfo(ctx)
}

func (m *meilisearch) GetRawIndex(ctx context.Context, uid string) (map[string]interface{}, error) {
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

func (m *meilisearch) ListIndexes(ctx context.Context, param *IndexesQuery) (*IndexesResults, error) {
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

func (m *meilisearch) GetRawIndexes(ctx context.Context, param *IndexesQuery) (map[string]interface{}, error) {
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

func (m *meilisearch) CreateIndex(ctx context.Context, config *IndexConfig) (*TaskInfo, error) {
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

func (m *meilisearch) DeleteIndex(ctx context.Context, uid string) (*TaskInfo, error) {
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

func (m *meilisearch) MultiSearch(ctx context.Context, queries *MultiSearchRequest) (*MultiSearchResponse, error) {
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

func (m *meilisearch) CreateKey(ctx context.Context, request *Key) (*Key, error) {
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

func (m *meilisearch) GetKey(ctx context.Context, identifier string) (*Key, error) {
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

func (m *meilisearch) GetKeys(ctx context.Context, param *KeysQuery) (*KeysResults, error) {
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

func (m *meilisearch) UpdateKey(ctx context.Context, keyOrUID string, request *Key) (*Key, error) {
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

func (m *meilisearch) DeleteKey(ctx context.Context, keyOrUID string) (bool, error) {
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

func (m *meilisearch) GetTask(ctx context.Context, taskUID int64) (*Task, error) {
	return getTask(ctx, m.client, taskUID)
}

func (m *meilisearch) GetTasks(ctx context.Context, param *TasksQuery) (*TaskResult, error) {
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

func (m *meilisearch) CancelTasks(ctx context.Context, param *CancelTasksQuery) (*TaskInfo, error) {
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

func (m *meilisearch) DeleteTasks(ctx context.Context, param *DeleteTasksQuery) (*TaskInfo, error) {
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

func (m *meilisearch) SwapIndexes(ctx context.Context, param []*SwapIndexesParams) (*TaskInfo, error) {
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

func (m *meilisearch) WaitForTask(ctx context.Context, taskUID int64, interval time.Duration) (*Task, error) {
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

func (m *meilisearch) GetStats(ctx context.Context, param *StatsParams) (*Stats, error) {
	resp := new(Stats)
	req := &internalRequest{
		endpoint:            "/stats",
		method:              http.MethodGet,
		withQueryParams:     map[string]string{},
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetStats",
	}
	if param != nil && param.SizeFormat != "" {
		req.withQueryParams["sizeFormat"] = param.SizeFormat
	}
	if param != nil && param.ShowInternalDatabaseSizes {
		req.withQueryParams["showInternalDatabaseSizes"] = "true"
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) CreateDump(ctx context.Context) (*TaskInfo, error) {
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

func (m *meilisearch) Version(ctx context.Context) (*Version, error) {
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

func (m *meilisearch) RenderTemplate(ctx context.Context, params *RenderTemplateParams) (*RenderTemplateResponse, error) {
	resp := new(RenderTemplateResponse)
	req := &internalRequest{
		endpoint:            "/render-template",
		method:              http.MethodPost,
		withRequest:         params,
		contentType:         contentTypeJSON,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "RenderTemplate",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) Health(ctx context.Context) (*Health, error) {
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

func (m *meilisearch) CreateSnapshot(ctx context.Context) (*TaskInfo, error) {
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
	res, err := m.Health(context.Background())
	return err == nil && res.Status == "available"
}

func (m *meilisearch) Close() {
	m.client.client.CloseIdleConnections()
}

func (m *meilisearch) GetBatches(ctx context.Context, param *BatchesQuery) (*BatchesResults, error) {
	resp := new(BatchesResults)
	req := &internalRequest{
		endpoint:            "/batches",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		withQueryParams:     map[string]string{},
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetBatches",
	}
	if param != nil {
		if len(param.UIDs) > 0 {
			req.withQueryParams["uids"] = joinInt64(param.UIDs)
		}
		if len(param.BatchUIDs) > 0 {
			req.withQueryParams["batchUids"] = joinInt64(param.BatchUIDs)
		}
		if len(param.IndexUIDs) > 0 {
			req.withQueryParams["indexUids"] = joinString(param.IndexUIDs)
		}
		if len(param.Statuses) > 0 {
			req.withQueryParams["statuses"] = joinString(param.Statuses)
		}
		if len(param.Types) > 0 {
			req.withQueryParams["types"] = joinString(param.Types)
		}
		if param.Limit > 0 {
			req.withQueryParams["limit"] = strconv.FormatInt(param.Limit, 10)
		}
		if param.From > 0 {
			req.withQueryParams["from"] = strconv.FormatInt(param.From, 10)
		}
		if param.Reverse {
			req.withQueryParams["reverse"] = "true"
		}
		if !param.BeforeEnqueuedAt.IsZero() {
			req.withQueryParams["beforeEnqueuedAt"] = param.BeforeEnqueuedAt.String()
		}
		if !param.BeforeStartedAt.IsZero() {
			req.withQueryParams["beforeStartedAt"] = param.BeforeStartedAt.String()
		}
		if !param.BeforeFinishedAt.IsZero() {
			req.withQueryParams["beforeFinishedAt"] = param.BeforeFinishedAt.String()
		}
		if !param.AfterEnqueuedAt.IsZero() {
			req.withQueryParams["afterEnqueuedAt"] = param.AfterEnqueuedAt.String()
		}
		if !param.AfterStartedAt.IsZero() {
			req.withQueryParams["afterStartedAt"] = param.AfterStartedAt.String()
		}
		if !param.AfterFinishedAt.IsZero() {
			req.withQueryParams["afterFinishedAt"] = param.AfterFinishedAt.String()
		}
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) GetBatch(ctx context.Context, batchUID int) (*Batch, error) {
	resp := new(Batch)
	req := &internalRequest{
		endpoint:            fmt.Sprintf("/batches/%d", batchUID),
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		withQueryParams:     map[string]string{},
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetBatch",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *meilisearch) Export(ctx context.Context, param *ExportParams) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/export",
		method:              http.MethodPost,
		withRequest:         param,
		contentType:         contentTypeJSON,
		withResponse:        resp,
		withQueryParams:     map[string]string{},
		acceptedStatusCodes: []int{http.StatusOK, http.StatusAccepted},
		functionName:        "Export",
	}
	if err := m.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
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
