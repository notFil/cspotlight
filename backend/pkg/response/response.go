package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/pagination"
	"github.com/notFil/cspotlight/pkg/errs"
)

type Response struct {
	Message    string                 `json:"message,omitempty"`
	Data       interface{}            `json:"data,omitempty"`
	Pagination *pagination.Pagination `json:"pagination,omitempty"`
	Error      *ErrorDetails          `json:"error,omitempty"`
}

type ErrorDetails struct {
	Details string `json:"details"`
}

func Success(c *gin.Context, httpStatus int, message string, data interface{}) {
	c.JSON(httpStatus, Response{
		Message: message,
		Data:    data,
	})
}

func SuccessPagedResponse(c *gin.Context, httpStatus int, message string, data interface{}, p *pagination.Pagination) {
	c.JSON(httpStatus, Response{
		Message:    message,
		Data:       data,
		Pagination: p,
	})
}

func ErrorResponse(c *gin.Context, httpStatus int, details string) {
	c.JSON(httpStatus, Response{
		Error: &ErrorDetails{
			Details: details,
		},
	})
}
func Error(c *gin.Context, err error) {
	statusCode := http.StatusInternalServerError
	message := "internal server error"

	switch {
	case errors.Is(err, errs.ErrUnauthorized):
		statusCode = http.StatusForbidden
		message = "unauthorized"
	case errors.Is(err, errs.ErrNotFound):
		statusCode = http.StatusNotFound
		message = "resource not found"
	case errors.Is(err, errs.ErrNoMatchOnNewPasswords):
		statusCode = http.StatusBadRequest
		message = "new password and confirm new password do not match"
	case errors.Is(err, errs.ErrInvalidInput):
		statusCode = http.StatusBadRequest
		message = "invalid input"
	}

	ErrorResponse(c, statusCode, message)
}
