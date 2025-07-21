package meilisearch

import (
	"context"
	"net/http"
)

func (i *index) GetSettings() (*Settings, error) {
	return i.GetSettingsWithContext(context.Background())
}

func (i *index) GetSettingsWithContext(ctx context.Context) (*Settings, error) {
	resp := new(Settings)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetSettings",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) UpdateSettings(request *Settings) (*TaskInfo, error) {
	return i.UpdateSettingsWithContext(context.Background(), request)
}

func (i *index) UpdateSettingsWithContext(ctx context.Context, request *Settings) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings",
		method:              http.MethodPatch,
		contentType:         contentTypeJSON,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateSettings",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) ResetSettings() (*TaskInfo, error) {
	return i.ResetSettingsWithContext(context.Background())
}

func (i *index) ResetSettingsWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetSettings",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetRankingRules() (*[]string, error) {
	return i.GetRankingRulesWithContext(context.Background())
}

func (i *index) GetRankingRulesWithContext(ctx context.Context) (*[]string, error) {
	resp := &[]string{}
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/ranking-rules",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetRankingRules",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) UpdateRankingRules(request *[]string) (*TaskInfo, error) {
	return i.UpdateRankingRulesWithContext(context.Background(), request)
}

func (i *index) UpdateRankingRulesWithContext(ctx context.Context, request *[]string) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/ranking-rules",
		method:              http.MethodPut,
		contentType:         contentTypeJSON,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateRankingRules",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) ResetRankingRules() (*TaskInfo, error) {
	return i.ResetRankingRulesWithContext(context.Background())
}

func (i *index) ResetRankingRulesWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/ranking-rules",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetRankingRules",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetDistinctAttribute() (*string, error) {
	return i.GetDistinctAttributeWithContext(context.Background())
}

func (i *index) GetDistinctAttributeWithContext(ctx context.Context) (*string, error) {
	resp := new(string)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/distinct-attribute",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetDistinctAttribute",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) UpdateDistinctAttribute(request string) (*TaskInfo, error) {
	return i.UpdateDistinctAttributeWithContext(context.Background(), request)
}

func (i *index) UpdateDistinctAttributeWithContext(ctx context.Context, request string) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/distinct-attribute",
		method:              http.MethodPut,
		contentType:         contentTypeJSON,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateDistinctAttribute",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) ResetDistinctAttribute() (*TaskInfo, error) {
	return i.ResetDistinctAttributeWithContext(context.Background())
}

func (i *index) ResetDistinctAttributeWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/distinct-attribute",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetDistinctAttribute",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetSearchableAttributes() (*[]string, error) {
	return i.GetSearchableAttributesWithContext(context.Background())
}

func (i *index) GetSearchableAttributesWithContext(ctx context.Context) (*[]string, error) {
	resp := &[]string{}
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/searchable-attributes",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetSearchableAttributes",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) UpdateSearchableAttributes(request *[]string) (*TaskInfo, error) {
	return i.UpdateSearchableAttributesWithContext(context.Background(), request)
}

func (i *index) UpdateSearchableAttributesWithContext(ctx context.Context, request *[]string) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/searchable-attributes",
		method:              http.MethodPut,
		contentType:         contentTypeJSON,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateSearchableAttributes",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) ResetSearchableAttributes() (*TaskInfo, error) {
	return i.ResetSearchableAttributesWithContext(context.Background())
}

func (i *index) ResetSearchableAttributesWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/searchable-attributes",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetSearchableAttributes",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetDisplayedAttributes() (*[]string, error) {
	return i.GetDisplayedAttributesWithContext(context.Background())
}

func (i *index) GetDisplayedAttributesWithContext(ctx context.Context) (*[]string, error) {
	resp := &[]string{}
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/displayed-attributes",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetDisplayedAttributes",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) UpdateDisplayedAttributes(request *[]string) (*TaskInfo, error) {
	return i.UpdateDisplayedAttributesWithContext(context.Background(), request)
}

func (i *index) UpdateDisplayedAttributesWithContext(ctx context.Context, request *[]string) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/displayed-attributes",
		method:              http.MethodPut,
		contentType:         contentTypeJSON,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateDisplayedAttributes",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) ResetDisplayedAttributes() (*TaskInfo, error) {
	return i.ResetDisplayedAttributesWithContext(context.Background())
}

func (i *index) ResetDisplayedAttributesWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/displayed-attributes",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetDisplayedAttributes",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetStopWords() (*[]string, error) {
	return i.GetStopWordsWithContext(context.Background())
}

func (i *index) GetStopWordsWithContext(ctx context.Context) (*[]string, error) {
	resp := &[]string{}
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/stop-words",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetStopWords",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) UpdateStopWords(request *[]string) (*TaskInfo, error) {
	return i.UpdateStopWordsWithContext(context.Background(), request)
}

func (i *index) UpdateStopWordsWithContext(ctx context.Context, request *[]string) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/stop-words",
		method:              http.MethodPut,
		contentType:         contentTypeJSON,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateStopWords",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) ResetStopWords() (*TaskInfo, error) {
	return i.ResetStopWordsWithContext(context.Background())
}

func (i *index) ResetStopWordsWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/stop-words",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetStopWords",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetSynonyms() (*map[string][]string, error) {
	return i.GetSynonymsWithContext(context.Background())
}

func (i *index) GetSynonymsWithContext(ctx context.Context) (*map[string][]string, error) {
	resp := &map[string][]string{}
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/synonyms",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetSynonyms",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) UpdateSynonyms(request *map[string][]string) (*TaskInfo, error) {
	return i.UpdateSynonymsWithContext(context.Background(), request)
}

func (i *index) UpdateSynonymsWithContext(ctx context.Context, request *map[string][]string) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/synonyms",
		method:              http.MethodPut,
		contentType:         contentTypeJSON,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateSynonyms",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) ResetSynonyms() (*TaskInfo, error) {
	return i.ResetSynonymsWithContext(context.Background())
}

func (i *index) ResetSynonymsWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/synonyms",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetSynonyms",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetFilterableAttributes() (*[]interface{}, error) {
	return i.GetFilterableAttributesWithContext(context.Background())
}

func (i *index) GetFilterableAttributesWithContext(ctx context.Context) (*[]interface{}, error) {
	resp := &[]interface{}{}
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/filterable-attributes",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetFilterableAttributes",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) UpdateFilterableAttributes(request *[]interface{}) (*TaskInfo, error) {
	return i.UpdateFilterableAttributesWithContext(context.Background(), request)
}

func (i *index) UpdateFilterableAttributesWithContext(ctx context.Context, request *[]interface{}) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/filterable-attributes",
		method:              http.MethodPut,
		contentType:         contentTypeJSON,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateFilterableAttributes",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) ResetFilterableAttributes() (*TaskInfo, error) {
	return i.ResetFilterableAttributesWithContext(context.Background())
}

func (i *index) ResetFilterableAttributesWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/filterable-attributes",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetFilterableAttributes",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetSortableAttributes() (*[]string, error) {
	return i.GetSortableAttributesWithContext(context.Background())
}

func (i *index) GetSortableAttributesWithContext(ctx context.Context) (*[]string, error) {
	resp := &[]string{}
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/sortable-attributes",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetSortableAttributes",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) UpdateSortableAttributes(request *[]string) (*TaskInfo, error) {
	return i.UpdateSortableAttributesWithContext(context.Background(), request)
}

func (i *index) UpdateSortableAttributesWithContext(ctx context.Context, request *[]string) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/sortable-attributes",
		method:              http.MethodPut,
		contentType:         contentTypeJSON,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateSortableAttributes",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) ResetSortableAttributes() (*TaskInfo, error) {
	return i.ResetSortableAttributesWithContext(context.Background())
}

func (i *index) ResetSortableAttributesWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/sortable-attributes",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetSortableAttributes",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetTypoTolerance() (*TypoTolerance, error) {
	return i.GetTypoToleranceWithContext(context.Background())
}

func (i *index) GetTypoToleranceWithContext(ctx context.Context) (*TypoTolerance, error) {
	resp := new(TypoTolerance)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/typo-tolerance",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetTypoTolerance",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) UpdateTypoTolerance(request *TypoTolerance) (*TaskInfo, error) {
	return i.UpdateTypoToleranceWithContext(context.Background(), request)
}

func (i *index) UpdateTypoToleranceWithContext(ctx context.Context, request *TypoTolerance) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/typo-tolerance",
		method:              http.MethodPatch,
		contentType:         contentTypeJSON,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateTypoTolerance",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) ResetTypoTolerance() (*TaskInfo, error) {
	return i.ResetTypoToleranceWithContext(context.Background())
}

func (i *index) ResetTypoToleranceWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/typo-tolerance",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetTypoTolerance",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetPagination() (*Pagination, error) {
	return i.GetPaginationWithContext(context.Background())
}

func (i *index) GetPaginationWithContext(ctx context.Context) (*Pagination, error) {
	resp := new(Pagination)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/pagination",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetPagination",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) UpdatePagination(request *Pagination) (*TaskInfo, error) {
	return i.UpdatePaginationWithContext(context.Background(), request)
}

func (i *index) UpdatePaginationWithContext(ctx context.Context, request *Pagination) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/pagination",
		method:              http.MethodPatch,
		contentType:         contentTypeJSON,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdatePagination",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) ResetPagination() (*TaskInfo, error) {
	return i.ResetPaginationWithContext(context.Background())
}

func (i *index) ResetPaginationWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/pagination",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetPagination",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetFaceting() (*Faceting, error) {
	return i.GetFacetingWithContext(context.Background())
}

func (i *index) GetFacetingWithContext(ctx context.Context) (*Faceting, error) {
	resp := new(Faceting)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/faceting",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetFaceting",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) UpdateFaceting(request *Faceting) (*TaskInfo, error) {
	return i.UpdateFacetingWithContext(context.Background(), request)
}

func (i *index) UpdateFacetingWithContext(ctx context.Context, request *Faceting) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/faceting",
		method:              http.MethodPatch,
		contentType:         contentTypeJSON,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateFaceting",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) ResetFaceting() (*TaskInfo, error) {
	return i.ResetFacetingWithContext(context.Background())
}

func (i *index) ResetFacetingWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/faceting",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetFaceting",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetEmbedders() (map[string]Embedder, error) {
	return i.GetEmbeddersWithContext(context.Background())
}

func (i *index) GetEmbeddersWithContext(ctx context.Context) (map[string]Embedder, error) {
	resp := make(map[string]Embedder)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/embedders",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetEmbedders",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) UpdateEmbedders(request map[string]Embedder) (*TaskInfo, error) {
	return i.UpdateEmbeddersWithContext(context.Background(), request)
}

func (i *index) UpdateEmbeddersWithContext(ctx context.Context, request map[string]Embedder) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/embedders",
		method:              http.MethodPatch,
		contentType:         contentTypeJSON,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateEmbedders",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) ResetEmbedders() (*TaskInfo, error) {
	return i.ResetEmbeddersWithContext(context.Background())
}

func (i *index) ResetEmbeddersWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/embedders",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetEmbedders",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetSearchCutoffMs() (int64, error) {
	return i.GetSearchCutoffMsWithContext(context.Background())
}

func (i *index) GetSearchCutoffMsWithContext(ctx context.Context) (int64, error) {
	var resp int64
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/search-cutoff-ms",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetSearchCutoffMs",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return 0, err
	}
	return resp, nil
}

func (i *index) UpdateSearchCutoffMs(request int64) (*TaskInfo, error) {
	return i.UpdateSearchCutoffMsWithContext(context.Background(), request)
}

func (i *index) UpdateSearchCutoffMsWithContext(ctx context.Context, request int64) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/search-cutoff-ms",
		method:              http.MethodPut,
		contentType:         contentTypeJSON,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateSearchCutoffMs",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) ResetSearchCutoffMs() (*TaskInfo, error) {
	return i.ResetSearchCutoffMsWithContext(context.Background())
}

func (i *index) ResetSearchCutoffMsWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/search-cutoff-ms",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetSearchCutoffMs",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetDictionary() ([]string, error) {
	return i.GetDictionaryWithContext(context.Background())
}

func (i *index) GetDictionaryWithContext(ctx context.Context) ([]string, error) {
	resp := make([]string, 0)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/dictionary",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetDictionary",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) UpdateDictionary(words []string) (*TaskInfo, error) {
	return i.UpdateDictionaryWithContext(context.Background(), words)
}

func (i *index) UpdateDictionaryWithContext(ctx context.Context, words []string) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/dictionary",
		method:              http.MethodPut,
		contentType:         contentTypeJSON,
		withRequest:         &words,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateDictionary",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) ResetDictionary() (*TaskInfo, error) {
	return i.ResetDictionaryWithContext(context.Background())
}

func (i *index) ResetDictionaryWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/dictionary",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetDictionary",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetSeparatorTokens() ([]string, error) {
	return i.GetSeparatorTokensWithContext(context.Background())
}

func (i *index) GetSeparatorTokensWithContext(ctx context.Context) ([]string, error) {
	resp := make([]string, 0)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/separator-tokens",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetSeparatorTokens",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) UpdateSeparatorTokens(req []string) (*TaskInfo, error) {
	return i.UpdateSeparatorTokensWithContext(context.Background(), req)
}

func (i *index) UpdateSeparatorTokensWithContext(ctx context.Context, tokens []string) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/separator-tokens",
		method:              http.MethodPut,
		withRequest:         &tokens,
		withResponse:        resp,
		contentType:         contentTypeJSON,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateSeparatorTokens",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) ResetSeparatorTokens() (*TaskInfo, error) {
	return i.ResetSeparatorTokensWithContext(context.Background())
}

func (i *index) ResetSeparatorTokensWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/separator-tokens",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetSeparatorTokens",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetNonSeparatorTokens() ([]string, error) {
	return i.GetNonSeparatorTokensWithContext(context.Background())
}

func (i *index) GetNonSeparatorTokensWithContext(ctx context.Context) ([]string, error) {
	resp := make([]string, 0)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/non-separator-tokens",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetNonSeparatorTokens",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) UpdateNonSeparatorTokens(req []string) (*TaskInfo, error) {
	return i.UpdateNonSeparatorTokensWithContext(context.Background(), req)
}

func (i *index) UpdateNonSeparatorTokensWithContext(ctx context.Context, tokens []string) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/non-separator-tokens",
		method:              http.MethodPut,
		withRequest:         &tokens,
		withResponse:        resp,
		contentType:         contentTypeJSON,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateNonSeparatorTokens",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) ResetNonSeparatorTokens() (*TaskInfo, error) {
	return i.ResetNonSeparatorTokensWithContext(context.Background())
}

func (i *index) ResetNonSeparatorTokensWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/non-separator-tokens",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetNonSeparatorTokens",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetProximityPrecision() (ProximityPrecisionType, error) {
	return i.GetProximityPrecisionWithContext(context.Background())
}

func (i *index) GetProximityPrecisionWithContext(ctx context.Context) (ProximityPrecisionType, error) {
	resp := new(ProximityPrecisionType)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/proximity-precision",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetProximityPrecision",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return "", err
	}
	return *resp, nil
}

func (i *index) UpdateProximityPrecision(proximityType ProximityPrecisionType) (*TaskInfo, error) {
	return i.UpdateProximityPrecisionWithContext(context.Background(), proximityType)
}

func (i *index) UpdateProximityPrecisionWithContext(ctx context.Context, proximityType ProximityPrecisionType) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/proximity-precision",
		method:              http.MethodPut,
		withRequest:         &proximityType,
		withResponse:        resp,
		contentType:         contentTypeJSON,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateProximityPrecision",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) ResetProximityPrecision() (*TaskInfo, error) {
	return i.ResetProximityPrecisionWithContext(context.Background())
}

func (i *index) ResetProximityPrecisionWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/proximity-precision",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetProximityPrecision",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetLocalizedAttributes() ([]*LocalizedAttributes, error) {
	return i.GetLocalizedAttributesWithContext(context.Background())
}

func (i *index) GetLocalizedAttributesWithContext(ctx context.Context) ([]*LocalizedAttributes, error) {
	resp := make([]*LocalizedAttributes, 0)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/localized-attributes",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetLocalizedAttributes",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) UpdateLocalizedAttributes(request []*LocalizedAttributes) (*TaskInfo, error) {
	return i.UpdateLocalizedAttributesWithContext(context.Background(), request)
}

func (i *index) UpdateLocalizedAttributesWithContext(ctx context.Context, request []*LocalizedAttributes) (*TaskInfo, error) {

	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/localized-attributes",
		method:              http.MethodPut,
		withRequest:         request,
		withResponse:        resp,
		contentType:         contentTypeJSON,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateLocalizedAttributes",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) ResetLocalizedAttributes() (*TaskInfo, error) {
	return i.ResetLocalizedAttributesWithContext(context.Background())
}

func (i *index) ResetLocalizedAttributesWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/localized-attributes",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetLocalizedAttributes",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetPrefixSearch() (*string, error) {
	return i.GetPrefixSearchWithContext(context.Background())
}

func (i *index) GetPrefixSearchWithContext(ctx context.Context) (*string, error) {
	resp := new(string)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/prefix-search",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetPrefixSearch",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) UpdatePrefixSearch(request string) (*TaskInfo, error) {
	return i.UpdatePrefixSearchWithContext(context.Background(), request)
}

func (i *index) UpdatePrefixSearchWithContext(ctx context.Context, request string) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/prefix-search",
		method:              http.MethodPut,
		contentType:         contentTypeJSON,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdatePrefixSearch",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) ResetPrefixSearch() (*TaskInfo, error) {
	return i.ResetPrefixSearchWithContext(context.Background())
}

func (i *index) ResetPrefixSearchWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/prefix-search",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetPrefixSearch",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) GetFacetSearch() (bool, error) {
	return i.GetFacetSearchWithContext(context.Background())
}

func (i *index) GetFacetSearchWithContext(ctx context.Context) (bool, error) {
	var resp bool
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/facet-search",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetFacetSearch",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return false, err
	}
	return resp, nil
}

func (i *index) UpdateFacetSearch(request bool) (*TaskInfo, error) {
	return i.UpdateFacetSearchWithContext(context.Background(), request)
}

func (i *index) UpdateFacetSearchWithContext(ctx context.Context, request bool) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/facet-search",
		method:              http.MethodPut,
		contentType:         contentTypeJSON,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateFacetSearch",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *index) ResetFacetSearch() (*TaskInfo, error) {
	return i.ResetFacetSearchWithContext(context.Background())
}

func (i *index) ResetFacetSearchWithContext(ctx context.Context) (*TaskInfo, error) {
	resp := new(TaskInfo)
	req := &internalRequest{
		endpoint:            "/indexes/" + i.uid + "/settings/facet-search",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetFacetSearch",
	}
	if err := i.client.executeRequest(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}
