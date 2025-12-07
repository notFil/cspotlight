package response

import (
	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/pagination"
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
