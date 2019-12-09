package meilisearch

import (
	"fmt"
	"strings"
)

type ErrCode int

const (
	ErrCodeUnknown        ErrCode = 0
	ErrCodeMarshalRequest ErrCode = iota + 1
	ErrCodeRequestCreation
	ErrCodeRequestExecution
	ErrCodeResponseStatusCode
	ErrCodeResponseReadBody
	ErrCodeResponseUnmarshalBody
)

const (
	rawStringCtx                   = `"ctx: Endpoint "${method} ${endpoint}" Function "${function}" Api "${apiName}"`
	rawStringMarshalRequest        = `unable to marshal body from request ${request}`
	rawStringRequestCreation       = `unable to create new request`
	rawStringRequestExecution      = `unable to execute request`
	rawStringResponseStatusCode    = `unaccepted status code found: ${statusCode} expected: ${statusCodeExpected}, message from api: '${RawMessage}'`
	rawStringResponseReadBody      = `unable to read body from response ${response}`
	rawStringResponseUnmarshalBody = `unable to unmarshal body from response ${response}`
)

func (e ErrCode) rawMessage() string {
	switch e {

	case ErrCodeMarshalRequest:
		return rawStringMarshalRequest + " " + rawStringCtx
	case ErrCodeRequestCreation:
		return rawStringRequestCreation + " " + rawStringCtx
	case ErrCodeRequestExecution:
		return rawStringRequestExecution + " " + rawStringCtx
	case ErrCodeResponseStatusCode:
		return rawStringResponseStatusCode + " " + rawStringCtx
	case ErrCodeResponseReadBody:
		return rawStringResponseReadBody + " " + rawStringCtx
	case ErrCodeResponseUnmarshalBody:
		return rawStringResponseUnmarshalBody + " " + rawStringCtx
	default:
		return "Unknown error " + rawStringCtx
	}
}

type apiMessage struct {
	Message string `json:"message"`
}

type Error struct {
	Endpoint string
	Method   string
	Function string
	APIName  string

	RequestToString  string
	ResponseToString string

	MeilisearchMessage string

	StatusCode         int
	StatusCodeExpected []int

	RawMessage string

	ErrCode ErrCode
}

func (e Error) Error() string {
	return namedSprintf(e.RawMessage, map[string]interface{}{
		"endpoint":           e.Endpoint,
		"method":             e.Method,
		"function":           e.Function,
		"apiName":            e.APIName,
		"request":            e.RequestToString,
		"response":           e.ResponseToString,
		"meilisearchMessage": e.MeilisearchMessage,
		"statusCodeExpected": e.StatusCodeExpected,
		"statusCode":         e.StatusCode,
	})
}

func namedSprintf(format string, params map[string]interface{}) string {
	for key, val := range params {
		format = strings.ReplaceAll(format, "${"+key+"}s", fmt.Sprintf("%s", val))
	}
	return format
}
