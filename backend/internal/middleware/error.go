package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/notFil/cspotlight/internal/errors"
	"github.com/notFil/cspotlight/internal/response"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		switch e := err.(type) {
		case *apperrors.AppError:
			response.ErrorResponse(c, e.Code, e.Message)
			return
		default:
			response.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
			return
		}

	}
}
