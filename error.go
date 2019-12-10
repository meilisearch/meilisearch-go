package meilisearch

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
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
	rawStringCtx                   = `(path "${method} ${endpoint}" with method "${apiName}.${function})`
	rawStringMarshalRequest        = `unable to marshal body from request: ${request}`
	rawStringRequestCreation       = `unable to create new request`
	rawStringRequestExecution      = `unable to execute request`
	rawStringResponseStatusCode    = `unaccepted status code found: ${statusCode} expected: ${statusCodeExpected}, message from api: '${meilisearchMessage}'`
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
		return rawStringCtx
	}
}

type apiMessage struct {
	Message string `json:"message"`
}

type MeiliError struct {
	Endpoint string
	Method   string
	Function string
	APIName  string

	RequestToString  string
	ResponseToString string

	MeilisearchMessage string

	StatusCode         int
	StatusCodeExpected []int

	rawMessage string

	OriginError error

	ErrCode ErrCode
}

func (e MeiliError) Error() string {
	message := namedSprintf(e.rawMessage, map[string]interface{}{
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
	if e.OriginError != nil {
		return errors.Wrap(e.OriginError, message).Error()
	}

	return message
}

func (e *MeiliError) WithMessage(str string, errs ...error) *MeiliError {
	if errs != nil {
		e.OriginError = errs[0]
	}

	e.rawMessage = str
	e.ErrCode = ErrCodeUnknown
	return e
}

func (e *MeiliError) WithErrCode(err ErrCode, errs ...error) *MeiliError {
	if errs != nil {
		e.OriginError = errs[0]
	}

	e.rawMessage = err.rawMessage()
	e.ErrCode = err
	return e
}

func (e *MeiliError) ErrorBody(body []byte) {
	e.ResponseToString = string(body)
	msg := apiMessage{}
	err := json.Unmarshal(body, &msg)
	if err != nil {
		e.MeilisearchMessage = msg.Message
	}
}

func namedSprintf(format string, params map[string]interface{}) string {
	for key, val := range params {
		format = strings.ReplaceAll(format, "${"+key+"}", fmt.Sprintf("%v", val))
	}
	return format
}
