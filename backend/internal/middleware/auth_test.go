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
	"github.com/stretchr/testify/assert"
)

func changeSession(c *gin.Context) {
	session := sessions.Default(c)
	// Helper to populate session for tests
	if uid, ok := c.Get("test_userID"); ok {
		session.Set("userID", uid)
	}
	if tid, ok := c.Get("test_teamID"); ok {
		session.Set("teamID", tid)
	}
	if role, ok := c.Get("test_role"); ok {
		session.Set("role", role)
	}
	// session.Save() // No need to save to cookie for same-request reading
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
				c.Set("test_userID", "123e4567-e89b-12d3-a456-426614174000")
				c.Set("test_teamID", "123e4567-e89b-12d3-a456-426614174001")
				c.Set("test_role", "user")
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Missing TeamID (Should Not Panic)",
			setupSession: func(c *gin.Context) {
				c.Set("test_userID", "123e4567-e89b-12d3-a456-426614174000")
				// No teamID set (simulating nil in session)
				c.Set("test_role", "user")
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid UserID Type",
			setupSession: func(c *gin.Context) {
				c.Set("test_userID", 123) // Not a string
				c.Set("test_teamID", "123e4567-e89b-12d3-a456-426614174001")
				c.Set("test_role", "user")
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Missing UserID",
			setupSession: func(c *gin.Context) {
				// No userID
				c.Set("test_teamID", "123e4567-e89b-12d3-a456-426614174001")
				c.Set("test_role", "user")
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupRouter()

			// Setup route with middleware needed for test
			r.GET("/test", func(c *gin.Context) {
				// Inject test data into context so changeSession can pick it up
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
