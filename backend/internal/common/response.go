package common

import "github.com/gin-gonic/gin"

type Response struct {
	Message string        `json:"message,omitempty"`
	Data    interface{}   `json:"data,omitempty"`
	Error   *ErrorDetails `json:"error,omitempty"`
}

type ErrorDetails struct {
	Details string `json:"details"`
}

func SuccessResponse(c *gin.Context, httpStatus int, message string, data interface{}) {
	c.JSON(httpStatus, Response{
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, httpStatus int, details string) {
	c.JSON(httpStatus, Response{
		Error: &ErrorDetails{
			Details: details,
		},
	})
}
