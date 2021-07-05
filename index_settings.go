package meilisearch

import (
	"net/http"
)

func (i Index) GetSettings() (resp *Settings, err error) {
	resp = &Settings{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetSettings",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) UpdateSettings(request *Settings) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateSettings",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) ResetSettings() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetSettings",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) GetRankingRules() (resp *[]string, err error) {
	resp = &[]string{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings/ranking-rules",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetRankingRules",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) UpdateRankingRules(request *[]string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings/ranking-rules",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateRankingRules",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) ResetRankingRules() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings/ranking-rules",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetRankingRules",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) GetDistinctAttribute() (resp *string, err error) {
	empty := ""
	resp = &empty
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings/distinct-attribute",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetDistinctAttribute",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) UpdateDistinctAttribute(request string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings/distinct-attribute",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateDistinctAttribute",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) ResetDistinctAttribute() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings/distinct-attribute",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetDistinctAttribute",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) GetSearchableAttributes() (resp *[]string, err error) {
	resp = &[]string{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings/searchable-attributes",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetSearchableAttributes",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) UpdateSearchableAttributes(request *[]string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings/searchable-attributes",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateSearchableAttributes",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) ResetSearchableAttributes() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings/searchable-attributes",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetSearchableAttributes",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) GetDisplayedAttributes() (resp *[]string, err error) {
	resp = &[]string{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings/displayed-attributes",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetDisplayedAttributes",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) UpdateDisplayedAttributes(request *[]string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings/displayed-attributes",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateDisplayedAttributes",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) ResetDisplayedAttributes() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings/displayed-attributes",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetDisplayedAttributes",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) GetStopWords() (resp *[]string, err error) {
	resp = &[]string{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings/stop-words",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetStopWords",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) UpdateStopWords(request *[]string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings/stop-words",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateStopWords",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) ResetStopWords() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings/stop-words",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetStopWords",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) GetSynonyms() (resp *map[string][]string, err error) {
	resp = &map[string][]string{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings/synonyms",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetSynonyms",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) UpdateSynonyms(request *map[string][]string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings/synonyms",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateSynonyms",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) ResetSynonyms() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings/synonyms",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetSynonyms",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) GetFilterableAttributes() (resp *[]string, err error) {
	resp = &[]string{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings/filterable-attributes",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetFilterableAttributes",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) UpdateFilterableAttributes(request *[]string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings/filterable-attributes",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateFilterableAttributes",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (i Index) ResetFilterableAttributes() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + i.UID + "/settings/filterable-attributes",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetFilterableAttributes",
	}
	if err := i.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}
