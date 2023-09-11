package apierror

import (
	"encoding/json"
	"fmt"
)

type APIErrorResponse struct {
	Success bool     `json:"success" tstype:"false"`
	Error   APIError `json:"error" tstype:",required"`
}

type APIError struct {
	ErrorCode string `json:"code"`

	// Human-friendly description that can be shown to the user.
	Detail string `json:"message"`

	Data json.RawMessage `json:"data,omitempty"`

	// The original error that caused this response, not forwarded to the end user.
	// Can be nil
	RawError error `json:"-"`
	// HTTP status that the error should have.
	Status int `json:"-"`
}

func Error(code string, detail string) *APIError {
	return &APIError{
		ErrorCode: code,
		Detail:    detail,
	}
}

func InternalError(err error, detail string) *APIError {
	return &APIError{
		ErrorCode: ErrorCodeInternalServerError,
		RawError:  err,
		Detail:    detail,
	}
}

func ErrorWithData(code string, detail string, data json.RawMessage) *APIError {
	return &APIError{
		ErrorCode: code,
		Detail:    detail,
		Data:      data,
	}
}

func (e *APIError) WithError(err error) *APIError {
	e.RawError = err
	return e
}

func (e *APIError) WithStatus(status int) *APIError {
	e.Status = status
	return e
}

func (e *APIError) Error() string {
	errMsg := ""
	if e.RawError != nil {
		errMsg = e.RawError.Error()
	}
	return fmt.Sprintf("%s (detail=%s) %s", e.ErrorCode, e.Detail, errMsg)
}

func IsAPIErrorWithCode(err error, errorCode string) bool {
	aerr, ok := err.(*APIError)
	return ok && aerr.ErrorCode == errorCode
}
