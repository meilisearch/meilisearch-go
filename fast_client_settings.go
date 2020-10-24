package meilisearch

import (
	"net/http"
)

type fastClientSettings struct {
	client   *FastHTTPClient
	indexUID string
}

func newFastClientSettings(client *FastHTTPClient, indexUID string) fastClientSettings {
	return fastClientSettings{client: client, indexUID: indexUID}
}

func (c fastClientSettings) IndexID() string {
	return c.indexUID
}

func (c fastClientSettings) Client() ClientInterface {
	return c.client
}

func (c fastClientSettings) GetAll() (resp *Settings, err error) {
	resp = &Settings{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetAll",
		apiName:             "Settings",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c fastClientSettings) UpdateAll(request Settings) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateAll",
		apiName:             "Documents",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c fastClientSettings) ResetAll() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetAll",
		apiName:             "Documents",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c fastClientSettings) GetRankingRules() (resp *[]string, err error) {
	resp = &[]string{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings/ranking-rules",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetRankingRules",
		apiName:             "Settings",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c fastClientSettings) UpdateRankingRules(request []string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings/ranking-rules",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateRankingRules",
		apiName:             "Documents",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c fastClientSettings) ResetRankingRules() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings/ranking-rules",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetRankingRules",
		apiName:             "Documents",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c fastClientSettings) GetDistinctAttribute() (*Str, error) {
	emp := Str("")
	resp := &emp
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings/distinct-attribute",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetDistinctAttribute",
		apiName:             "Settings",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c fastClientSettings) UpdateDistinctAttribute(request Str) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings/distinct-attribute",
		method:              http.MethodPost,
		withRequest:         request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateDistinctAttribute",
		apiName:             "Documents",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c fastClientSettings) ResetDistinctAttribute() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings/distinct-attribute",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetDistinctAttribute",
		apiName:             "Documents",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c fastClientSettings) GetSearchableAttributes() (resp *[]string, err error) {
	resp = &[]string{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings/searchable-attributes",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetSearchableAttributes",
		apiName:             "Settings",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c fastClientSettings) UpdateSearchableAttributes(request []string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings/searchable-attributes",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateSearchableAttributes",
		apiName:             "Documents",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c fastClientSettings) ResetSearchableAttributes() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings/searchable-attributes",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetSearchableAttributes",
		apiName:             "Documents",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c fastClientSettings) GetDisplayedAttributes() (resp *[]string, err error) {
	resp = &[]string{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings/displayed-attributes",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetDisplayedAttributes",
		apiName:             "Settings",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c fastClientSettings) UpdateDisplayedAttributes(request []string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings/displayed-attributes",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateDisplayedAttributes",
		apiName:             "Documents",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c fastClientSettings) ResetDisplayedAttributes() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings/displayed-attributes",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetDisplayedAttributes",
		apiName:             "Documents",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c fastClientSettings) GetStopWords() (resp *[]string, err error) {
	resp = &[]string{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings/stop-words",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetStopWords",
		apiName:             "Settings",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c fastClientSettings) UpdateStopWords(request []string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings/stop-words",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateStopWords",
		apiName:             "Documents",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c fastClientSettings) ResetStopWords() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings/stop-words",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetStopWords",
		apiName:             "Documents",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c fastClientSettings) GetSynonyms() (resp *map[string][]string, err error) {
	resp = &map[string][]string{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings/synonyms",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetSynonyms",
		apiName:             "Settings",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c fastClientSettings) UpdateSynonyms(request map[string][]string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings/synonyms",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateSynonyms",
		apiName:             "Documents",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c fastClientSettings) ResetSynonyms() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings/synonyms",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetSynonyms",
		apiName:             "Documents",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c fastClientSettings) GetAttributesForFaceting() (resp *[]string, err error) {
	resp = &[]string{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings/attributes-for-faceting",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetAttributesForFaceting",
		apiName:             "Settings",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c fastClientSettings) UpdateAttributesForFaceting(request []string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings/attributes-for-faceting",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "UpdateAttributesForFaceting",
		apiName:             "Settings",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c fastClientSettings) ResetAttributesForFaceting() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRawRequest{
		endpoint:            "/indexes/" + c.indexUID + "/settings/attributes-for-faceting",
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "ResetAttributesForFaceting",
		apiName:             "Settings",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}
