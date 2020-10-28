package meilisearch

import (
	"net/http"
)

type clientSettings struct {
	client   *Client
	indexUID string
}

func newClientSettings(client *Client, indexUID string) clientSettings {
	return clientSettings{client: client, indexUID: indexUID}
}

func (c clientSettings) IndexID() string {
	return c.indexUID
}

func (c clientSettings) Client() ClientInterface {
	return c.client
}

func (c clientSettings) GetAll() (resp *Settings, err error) {
	resp = &Settings{}
	req := internalRequest{
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

func (c clientSettings) UpdateAll(request Settings) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
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

func (c clientSettings) ResetAll() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
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

func (c clientSettings) GetRankingRules() (resp *[]string, err error) {
	resp = &[]string{}
	req := internalRequest{
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

func (c clientSettings) UpdateRankingRules(request []string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
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

func (c clientSettings) ResetRankingRules() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
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

func (c clientSettings) GetDistinctAttribute() (*Str, error) {
	emp := Str("")
	resp := &emp
	req := internalRequest{
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

func (c clientSettings) UpdateDistinctAttribute(request Str) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
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

func (c clientSettings) ResetDistinctAttribute() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
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

func (c clientSettings) GetSearchableAttributes() (resp *[]string, err error) {
	resp = &[]string{}
	req := internalRequest{
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

func (c clientSettings) UpdateSearchableAttributes(request []string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
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

func (c clientSettings) ResetSearchableAttributes() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
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

func (c clientSettings) GetDisplayedAttributes() (resp *[]string, err error) {
	resp = &[]string{}
	req := internalRequest{
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

func (c clientSettings) UpdateDisplayedAttributes(request []string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
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

func (c clientSettings) ResetDisplayedAttributes() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
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

func (c clientSettings) GetStopWords() (resp *[]string, err error) {
	resp = &[]string{}
	req := internalRequest{
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

func (c clientSettings) UpdateStopWords(request []string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
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

func (c clientSettings) ResetStopWords() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
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

func (c clientSettings) GetSynonyms() (resp *map[string][]string, err error) {
	resp = &map[string][]string{}
	req := internalRequest{
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

func (c clientSettings) UpdateSynonyms(request map[string][]string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
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

func (c clientSettings) ResetSynonyms() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
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

func (c clientSettings) GetAttributesForFaceting() (resp *[]string, err error) {
	resp = &[]string{}
	req := internalRequest{
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

func (c clientSettings) UpdateAttributesForFaceting(request []string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
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

func (c clientSettings) ResetAttributesForFaceting() (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
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
