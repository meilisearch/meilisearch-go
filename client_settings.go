package meilisearch

import (
	"net/http"
)

type clientSettings struct {
	client  *Client
	indexID string
}

func (c clientSettings) IndexID() string {
	return c.indexID
}

func (c clientSettings) Client() *Client {
	return c.client
}

func newClientSettings(client *Client, indexID string) clientSettings {
	return clientSettings{client: client, indexID: indexID}
}

func (c clientSettings) GetAll() (resp *Settings, err error) {
	resp = &Settings{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/settings",
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

func (c clientSettings) AddOrUpdateAll(request Settings) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/settings",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "AddOrUpdateAll",
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
		endpoint:            "/indexes/" + c.indexID + "/settings",
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
		endpoint:            "/indexes/" + c.indexID + "/settings/ranking-rules",
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

func (c clientSettings) SetRankingRules(request []string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/settings/ranking-rules",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "SetRankingRules",
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
		endpoint:            "/indexes/" + c.indexID + "/settings/ranking-rules",
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

func (c clientSettings) GetDistinctAttribute() (resp *string, err error) {
	empty := ""
	resp = &empty
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/settings/distinct-attribute",
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

func (c clientSettings) SetDistinctAttribute(request string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/settings/distinct-attribute",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "SetDistinctAttribute",
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
		endpoint:            "/indexes/" + c.indexID + "/settings/distinct-attribute",
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
		endpoint:            "/indexes/" + c.indexID + "/settings/searchable-attributes",
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

func (c clientSettings) SetSearchableAttributes(request []string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/settings/searchable-attributes",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "SetSearchableAttributes",
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
		endpoint:            "/indexes/" + c.indexID + "/settings/searchable-attributes",
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
		endpoint:            "/indexes/" + c.indexID + "/settings/displayed-attributes",
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

func (c clientSettings) SetDisplayedAttributes(request []string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/settings/displayed-attributes",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "SetDisplayedAttributes",
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
		endpoint:            "/indexes/" + c.indexID + "/settings/displayed-attributes",
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
		endpoint:            "/indexes/" + c.indexID + "/settings/stop-words",
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

func (c clientSettings) SetStopWords(request []string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/settings/stop-words",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "SetStopWords",
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
		endpoint:            "/indexes/" + c.indexID + "/settings/stop-words",
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
		endpoint:            "/indexes/" + c.indexID + "/settings/synonyms",
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

func (c clientSettings) SetSynonyms(request map[string][]string) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/settings/synonyms",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "SetSynonyms",
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
		endpoint:            "/indexes/" + c.indexID + "/settings/synonyms",
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

func (c clientSettings) GetAcceptNewFields() (resp *bool, err error) {
	empty := true
	resp = &empty
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/settings/accept-new-fields",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetAcceptNewFields",
		apiName:             "Settings",
	}

	if err := c.client.executeRequest(req); err != nil {
		return resp, err
	}
	return resp, nil
}

func (c clientSettings) SetAcceptNewFields(request bool) (resp *AsyncUpdateID, err error) {
	resp = &AsyncUpdateID{}
	req := internalRequest{
		endpoint:            "/indexes/" + c.indexID + "/settings/accept-new-fields",
		method:              http.MethodPost,
		withRequest:         &request,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "SetAcceptNewFields",
		apiName:             "Documents",
	}

	if err := c.client.executeRequest(req); err != nil {
		return nil, err
	}

	return resp, nil
}
