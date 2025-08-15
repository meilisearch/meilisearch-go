package meilisearch

import (
	"context"
	"net/http"
)

type ExperimentalFeatures struct {
	client *client
	ExperimentalFeaturesBase
}

func (m *meilisearch) ExperimentalFeatures() *ExperimentalFeatures {
	return &ExperimentalFeatures{client: m.client}
}

func (ef *ExperimentalFeatures) SetLogsRoute(logsRoute bool) *ExperimentalFeatures {
	ef.LogsRoute = &logsRoute
	return ef
}

func (ef *ExperimentalFeatures) SetMetrics(metrics bool) *ExperimentalFeatures {
	ef.Metrics = &metrics
	return ef
}

func (ef *ExperimentalFeatures) SetEditDocumentsByFunction(editDocumentsByFunction bool) *ExperimentalFeatures {
	ef.EditDocumentsByFunction = &editDocumentsByFunction
	return ef
}

func (ef *ExperimentalFeatures) SetContainsFilter(containsFilter bool) *ExperimentalFeatures {
	ef.ContainsFilter = &containsFilter
	return ef
}

func (ef *ExperimentalFeatures) SetNetwork(enable bool) *ExperimentalFeatures {
	ef.Network = &enable
	return ef
}

func (ef *ExperimentalFeatures) SetCompositeEmbedders(enable bool) *ExperimentalFeatures {
	ef.CompositeEmbedders = &enable
	return ef
}

func (ef *ExperimentalFeatures) SetChatCompletions(enable bool) *ExperimentalFeatures {
	ef.ChatCompletions = &enable
	return ef
}

func (ef *ExperimentalFeatures) SetMultiModal(enable bool) *ExperimentalFeatures {
	ef.MultiModal = &enable
	return ef
}

func (ef *ExperimentalFeatures) Get() (*ExperimentalFeaturesResult, error) {
	return ef.GetWithContext(context.Background())
}

func (ef *ExperimentalFeatures) GetWithContext(ctx context.Context) (*ExperimentalFeaturesResult, error) {
	resp := new(ExperimentalFeaturesResult)
	req := &internalRequest{
		endpoint:            "/experimental-features",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		withQueryParams:     map[string]string{},
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetExperimentalFeatures",
	}

	if err := ef.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (ef *ExperimentalFeatures) Update() (*ExperimentalFeaturesResult, error) {
	return ef.UpdateWithContext(context.Background())
}

func (ef *ExperimentalFeatures) UpdateWithContext(ctx context.Context) (*ExperimentalFeaturesResult, error) {
	request := ExperimentalFeaturesBase{
		LogsRoute:               ef.LogsRoute,
		Metrics:                 ef.Metrics,
		EditDocumentsByFunction: ef.EditDocumentsByFunction,
		ContainsFilter:          ef.ContainsFilter,
		Network:                 ef.Network,
		CompositeEmbedders:      ef.CompositeEmbedders,
		ChatCompletions:         ef.ChatCompletions,
		MultiModal:              ef.MultiModal,
	}
	resp := new(ExperimentalFeaturesResult)
	req := &internalRequest{
		endpoint:            "/experimental-features",
		method:              http.MethodPatch,
		contentType:         contentTypeJSON,
		withRequest:         request,
		withResponse:        resp,
		withQueryParams:     nil,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "UpdateExperimentalFeatures",
	}
	if err := ef.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}
