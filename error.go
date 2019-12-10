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

// Error is the internal error structure that all exposed method use.
// So ALL errors returned by this library can be cast to this struct (as a pointer)
type Error struct {
	// Endpoint is the path of the request (host is not in)
	Endpoint string

	// Method is the HTTP verb of the request
	Method string

	// Function name used
	Function string

	// APIName is which part/module of the api
	APIName string

	// RequestToString is the raw request into string ('empty request' if not present)
	RequestToString string

	// RequestToString is the raw request into string ('empty response' if not present)
	ResponseToString string

	// MeilisearchMessage is the raw request into string ('empty meilisearch message' if not present)
	MeilisearchMessage string

	// StatusCode of the request
	StatusCode int

	// StatusCode expected by the endpoint to be considered as a success
	StatusCodeExpected []int

	rawMessage string

	// OriginError is the origin error that produce the current Error. It can be nil in case of a bad status code.
	OriginError error

	// ErrCode is the internal error code that represent the different step when executing a request that can produce
	// an error.
	ErrCode ErrCode
}

// Error return a well human formatted message.
func (e Error) Error() string {
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

func (e *Error) WithMessage(str string, errs ...error) *Error {
	if errs != nil {
		e.OriginError = errs[0]
	}

	e.rawMessage = str
	e.ErrCode = ErrCodeUnknown
	return e
}

func (e *Error) WithErrCode(err ErrCode, errs ...error) *Error {
	if errs != nil {
		e.OriginError = errs[0]
	}

	e.rawMessage = err.rawMessage()
	e.ErrCode = err
	return e
}

func (e *Error) ErrorBody(body []byte) {
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
