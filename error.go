package meilisearch

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ErrCode are all possible errors found during requests
type ErrCode int

const (
	// ErrCodeUnknown default error code, undefined
	ErrCodeUnknown ErrCode = 0
	// ErrCodeMarshalRequest impossible to serialize request body
	ErrCodeMarshalRequest ErrCode = iota + 1
	// ErrCodeResponseUnmarshalBody impossible deserialize the response body
	ErrCodeResponseUnmarshalBody
	// MeilisearchApiError send by the Meilisearch api
	MeilisearchApiError
	// MeilisearchApiError send by the Meilisearch api
	MeilisearchApiErrorWithoutMessage
	// MeilisearchTimeoutError
	MeilisearchTimeoutError
	// MeilisearchCommunicationError impossible execute a request
	MeilisearchCommunicationError
)

const (
	rawStringCtx                               = `(path "${method} ${endpoint}" with method "${function}")`
	rawStringMarshalRequest                    = `unable to marshal body from request: '${request}'`
	rawStringResponseUnmarshalBody             = `unable to unmarshal body from response: '${response}' status code: ${statusCode}`
	rawStringMeilisearchApiError               = `unaccepted status code found: ${statusCode} expected: ${statusCodeExpected}, MeilisearchApiError Message: ${message}, Code: ${code}, Type: ${type}, Link: ${link}`
	rawStringMeilisearchApiErrorWithoutMessage = `unaccepted status code found: ${statusCode} expected: ${statusCodeExpected}, MeilisearchApiError Message: ${message}`
	rawStringMeilisearchTimeoutError           = `MeilisearchTimeoutError`
	rawStringMeilisearchCommunicationError     = `MeilisearchCommunicationError unable to execute request`
)

func (e ErrCode) rawMessage() string {
	switch e {
	case ErrCodeMarshalRequest:
		return rawStringMarshalRequest + " " + rawStringCtx
	case ErrCodeResponseUnmarshalBody:
		return rawStringResponseUnmarshalBody + " " + rawStringCtx
	case MeilisearchApiError:
		return rawStringMeilisearchApiError + " " + rawStringCtx
	case MeilisearchApiErrorWithoutMessage:
		return rawStringMeilisearchApiErrorWithoutMessage + " " + rawStringCtx
	case MeilisearchTimeoutError:
		return rawStringMeilisearchTimeoutError + " " + rawStringCtx
	case MeilisearchCommunicationError:
		return rawStringMeilisearchCommunicationError + " " + rawStringCtx
	default:
		return rawStringCtx
	}
}

type meilisearchApiError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Type    string `json:"type"`
	Link    string `json:"link"`
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

	// RequestToString is the raw request into string ('empty request' if not present)
	RequestToString string

	// RequestToString is the raw request into string ('empty response' if not present)
	ResponseToString string

	// Error info from Meilisearch api
	// Message is the raw request into string ('empty Meilisearch message' if not present)
	MeilisearchApiError meilisearchApiError

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
		"request":            e.RequestToString,
		"response":           e.ResponseToString,
		"statusCodeExpected": e.StatusCodeExpected,
		"statusCode":         e.StatusCode,
		"message":            e.MeilisearchApiError.Message,
		"code":               e.MeilisearchApiError.Code,
		"type":               e.MeilisearchApiError.Type,
		"link":               e.MeilisearchApiError.Link,
	})
	if e.OriginError != nil {
		return fmt.Sprintf("%s: %s", message, e.OriginError.Error())
	}

	return message
}

// WithErrCode add an error code to an error
func (e *Error) WithErrCode(err ErrCode, errs ...error) *Error {
	if errs != nil {
		e.OriginError = errs[0]
	}

	e.rawMessage = err.rawMessage()
	e.ErrCode = err
	return e
}

// ErrorBody add a body to an error
func (e *Error) ErrorBody(body []byte) {
	e.ResponseToString = string(body)
	msg := meilisearchApiError{}
	err := json.Unmarshal(body, &msg)
	if err == nil {
		e.MeilisearchApiError.Message = msg.Message
		e.MeilisearchApiError.Code = msg.Code
		e.MeilisearchApiError.Type = msg.Type
		e.MeilisearchApiError.Link = msg.Link
	}
}

// Added a hint to the error message if it may come from a version incompatibility with Meilisearch
func VersionErrorHintMessage(err error, req *internalRequest) error {
	return fmt.Errorf("%w. Hint: It might not be working because you're not up to date with the Meilisearch version that %s call requires.", err, req.functionName)
}

func namedSprintf(format string, params map[string]interface{}) string {
	for key, val := range params {
		format = strings.ReplaceAll(format, "${"+key+"}", fmt.Sprintf("%v", val))
	}
	return format
}
