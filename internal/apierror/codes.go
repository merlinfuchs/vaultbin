package apierror

import "net/http"

const (
	ErrorCodeInternalServerError = "internal_server_error"
	ErrorCodeBadRequest          = "bad_request"

	ErrorCodeAuthRequired = "auth_required"
	ErrorCodeForbidden    = "forbidden"

	ErrorCodeMethodNotAllowed = "method_not_allowed"
	ErrorCodeNotFound         = "not_found"
)

var ErrorCodeToDefaultStatus = map[string]int{
	ErrorCodeInternalServerError: http.StatusInternalServerError,
	ErrorCodeBadRequest:          http.StatusBadRequest,

	ErrorCodeAuthRequired: http.StatusUnauthorized,
	ErrorCodeForbidden:    http.StatusForbidden,

	ErrorCodeMethodNotAllowed: http.StatusMethodNotAllowed,
	ErrorCodeNotFound:         http.StatusNotFound,
}
