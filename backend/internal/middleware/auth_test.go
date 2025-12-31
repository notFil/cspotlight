package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cspotlight/internal/logger"
	"cspotlight/internal/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func changeSession(c *gin.Context) {
	session := sessions.Default(c)
	if uid, ok := c.Get("test_userID"); ok {
		session.Set("userID", uid)
	}
	if tid, ok := c.Get("test_teamID"); ok {
		session.Set("teamID", tid)
	}
	if role, ok := c.Get("test_role"); ok {
		session.Set("role", role)
	}
}

func setupRouter() *gin.Engine {
	logger.InitializeLogger(false)
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(middleware.ErrorHandler())
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	return r
}

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		setupSession   func(c *gin.Context)
		expectedStatus int
	}{
		{
			name: "Valid Session",
			setupSession: func(c *gin.Context) {
				c.Set("test_userID", uuid.New())
				c.Set("test_teamID", uuid.New())
				c.Set("test_role", "user")
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Missing TeamID (Should Not Panic)",
			setupSession: func(c *gin.Context) {
				c.Set("test_userID", uuid.New())
				c.Set("test_role", "user")
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid UserID Type",
			setupSession: func(c *gin.Context) {
				c.Set("test_userID", 123)
				c.Set("test_teamID", uuid.New())
				c.Set("test_role", "user")
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Missing UserID",
			setupSession: func(c *gin.Context) {
				c.Set("test_teamID", uuid.New())
				c.Set("test_role", "user")
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupRouter()

			r.GET("/test", func(c *gin.Context) {
				tt.setupSession(c)
				changeSession(c)
			}, middleware.AuthMiddleware, func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
