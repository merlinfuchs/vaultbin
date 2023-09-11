package apierror

import (
	"net/http"

	"log/slog"

	"github.com/labstack/echo/v4"
)

func ErrorRewriteMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				return errorHandler(c, err)
			}
			return err
		}
	}
}

func shouldLogError(errorCode string) bool {
	return errorCode == ErrorCodeInternalServerError
}

func errorHandler(c echo.Context, err error) error {
	statusCode := http.StatusInternalServerError
	if e, ok := err.(*APIError); ok {
		if e.Status != 0 {
			statusCode = e.Status
		} else {
			statusCode = ErrorCodeToDefaultStatus[e.ErrorCode]
		}

		if shouldLogError(e.ErrorCode) {
			slog.With("error_code", e.ErrorCode).With("error", err).Error("")
		}

		return c.JSON(statusCode, APIErrorResponse{
			Success: false,
			Error:   *e,
		})
	}

	errorCode := ErrorCodeInternalServerError
	if e, ok := err.(*echo.HTTPError); ok {
		statusCode = e.Code
		if e.Code == http.StatusMethodNotAllowed {
			errorCode = ErrorCodeMethodNotAllowed
		} else if e.Code == http.StatusNotFound {
			errorCode = ErrorCodeNotFound
		}
	}

	if shouldLogError(errorCode) {
		slog.With("error", err).Error("")
	}

	return c.JSON(statusCode, APIErrorResponse{
		Success: false,
		Error: APIError{
			ErrorCode: errorCode,
			Detail:    "",
		},
	})
}
