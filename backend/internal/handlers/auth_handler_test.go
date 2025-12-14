package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/notFil/cspotlight/internal/auth"
	"github.com/notFil/cspotlight/internal/logger"
	"github.com/notFil/cspotlight/internal/middleware"
	"github.com/notFil/cspotlight/internal/models"
)

func init() {
	logger.InitializeLogger(false)
}

type MockUserService struct {
	RegisterUserFunc      func(ctx context.Context, user *models.UserRegisterDTO) error
	GetUserByIDFunc       func(ctx context.Context, id uuid.UUID) (*models.UserFetchDTO, error)
	GetUserByUsernameFunc func(ctx context.Context, username string) (*models.UserFetchDTO, error)
	AuthenticateUserFunc  func(ctx context.Context, authRequest *models.AuthRequest) (*models.UserFetchDTO, error)
	UpdateUserFunc        func(ctx context.Context, id uuid.UUID, user *models.UserUpdateDTO) (*models.UserFetchDTO, error)
	SetDefaultProjectFunc func(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (*models.UserFetchDTO, error)
	GetUsersFunc          func(ctx context.Context) ([]*models.UserFetchDTO, error)
	ListUsersByTeamIDFunc func(ctx context.Context, teamID uuid.UUID) ([]*models.UserFetchDTO, error)
	DeleteUserFunc        func(ctx context.Context, id uuid.UUID) error
	ChangePasswordFunc    func(ctx context.Context, id uuid.UUID, request *models.ChangePasswordRequest) error
	ChangeImageFunc       func(ctx context.Context, id uuid.UUID, path string) error
}

func (m *MockUserService) RegisterUser(ctx context.Context, user *models.UserRegisterDTO) error {
	if m.RegisterUserFunc != nil {
		return m.RegisterUserFunc(ctx, user)
	}
	return nil
}

func (m *MockUserService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserFetchDTO, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockUserService) AuthenticateUser(ctx context.Context, authRequest *models.AuthRequest) (*models.UserFetchDTO, error) {
	if m.AuthenticateUserFunc != nil {
		return m.AuthenticateUserFunc(ctx, authRequest)
	}
	return nil, nil
}

func (m *MockUserService) UpdateUser(ctx context.Context, id uuid.UUID, user *models.UserUpdateDTO) (*models.UserFetchDTO, error) {
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(ctx, id, user)
	}
	return nil, nil
}

func (m *MockUserService) SetDefaultProject(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (*models.UserFetchDTO, error) {
	if m.SetDefaultProjectFunc != nil {
		return m.SetDefaultProjectFunc(ctx, id, projectID)
	}
	return nil, nil
}

func (m *MockUserService) ListUsers(ctx context.Context) ([]*models.UserFetchDTO, error) {
	if m.GetUsersFunc != nil {
		return m.GetUsersFunc(ctx)
	}
	return nil, nil
}

func (m *MockUserService) ListUsersByTeamID(ctx context.Context, teamID uuid.UUID) ([]*models.UserFetchDTO, error) {
	if m.ListUsersByTeamIDFunc != nil {
		return m.ListUsersByTeamIDFunc(ctx, teamID)
	}
	return nil, nil
}

func (m *MockUserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(ctx, id)
	}
	return nil
}

func (m *MockUserService) ChangePassword(ctx context.Context, id uuid.UUID, request *models.ChangePasswordRequest) error {
	if m.ChangePasswordFunc != nil {
		return m.ChangePasswordFunc(ctx, id, request)
	}
	return nil
}

func (m *MockUserService) ChangeImage(ctx context.Context, id uuid.UUID, path string) error {
	if m.ChangeImageFunc != nil {
		return m.ChangeImageFunc(ctx, id, path)
	}
	return nil
}

type MockCache struct {
	SetFunc func(ctx context.Context, key string, value string, expiration time.Duration) error
	GetFunc func(ctx context.Context, key string) (string, error)
	DelFunc func(ctx context.Context, key string) error
}

func (m *MockCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	if m.SetFunc != nil {
		return m.SetFunc(ctx, key, value, expiration)
	}
	return nil
}

func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, key)
	}
	return "", nil
}

func (m *MockCache) Del(ctx context.Context, key string) error {
	if m.DelFunc != nil {
		return m.DelFunc(ctx, key)
	}
	return nil
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			AuthenticateUserFunc: func(ctx context.Context, authRequest *models.AuthRequest) (*models.UserFetchDTO, error) {
				return &models.UserFetchDTO{ID: uuid.New(), Username: "testuser"}, nil
			},
		}
		handler := NewAuthHandler(mockService)

		w := httptest.NewRecorder()
		router := gin.New()
		store := cookie.NewStore([]byte("secret"))
		router.Use(sessions.Sessions("session", store))
		router.Use(middleware.ErrorHandler())
		router.POST("/login", handler.Login)

		body := `{"username": "testuser", "password": "password"}`
		req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		if bytes.Contains(w.Body.Bytes(), []byte("testuser")) {
			t.Errorf("expected response to not contain 'testuser', got %s", w.Body.String())
		}
	})

	t.Run("InvalidPayload", func(t *testing.T) {
		handler := NewAuthHandler(&MockUserService{})
		w := httptest.NewRecorder()
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.POST("/login", handler.Login)

		req := httptest.NewRequest("POST", "/login", bytes.NewBufferString("invalid"))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("AuthenticationFailed", func(t *testing.T) {
		mockService := &MockUserService{
			AuthenticateUserFunc: func(ctx context.Context, authRequest *models.AuthRequest) (*models.UserFetchDTO, error) {
				return nil, errors.New("invalid credentials")
			},
		}
		handler := NewAuthHandler(mockService)

		w := httptest.NewRecorder()
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.POST("/login", handler.Login)

		body := `{"username": "testuser", "password": "wrongpassword"}`
		req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d. Body: %s", w.Code, w.Body.String())
		}
	})
}

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			RegisterUserFunc: func(ctx context.Context, user *models.UserRegisterDTO) error {
				return nil
			},
		}
		handler := NewAuthHandler(mockService)

		w := httptest.NewRecorder()
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.POST("/register", handler.Register)

		body := `{"email": "test@example.com", "password": "password", "confirmPassword": "password", "firstName": "Test", "lastName": "User", "username": "testuser", "role": "user"}`
		req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}
	})

	t.Run("Failure", func(t *testing.T) {
		mockService := &MockUserService{
			RegisterUserFunc: func(ctx context.Context, user *models.UserRegisterDTO) error {
				return errors.New("email exists")
			},
		}
		handler := NewAuthHandler(mockService)

		w := httptest.NewRecorder()
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.POST("/register", handler.Register)

		body := `{"email": "test@example.com", "password": "password", "confirmPassword": "password", "firstName": "Test", "lastName": "User", "username": "testuser", "role": "user"}`
		req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d. Body: %s", w.Code, w.Body.String())
		}
	})
}

func TestAuthHandler_SignOut(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(&MockUserService{})

	w := httptest.NewRecorder()
	router := gin.New()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("session", store))
	router.Use(middleware.ErrorHandler())

	// Inject user context
	router.Use(func(c *gin.Context) {
		userCtx := auth.NewUserContext(uuid.New().String(), uuid.New().String(), "user")
		c.Request = c.Request.WithContext(auth.ContextWithUser(c.Request.Context(), userCtx))
		c.Next()
	})

	router.POST("/signout", handler.SignOut)

	req := httptest.NewRequest("POST", "/signout", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	cookie := w.Result().Header.Get("Set-Cookie")
	if cookie == "" {
		t.Error("expected Set-Cookie header to be present")
	}
}
