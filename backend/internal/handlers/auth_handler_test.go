package handlers

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/config"
	"github.com/notFil/cspotlight/internal/logger"
	"github.com/notFil/cspotlight/internal/models"
)

func init() {
	logger.InitializeLogger("test")
}

type MockUserService struct {
	RegisterUserFunc      func(ctx context.Context, user *models.UserRegisterDTO) error
	GetUserByIDFunc       func(ctx context.Context, id string) (*models.UserFetchDTO, error)
	GetUserByUsernameFunc func(ctx context.Context, username string) (*models.UserFetchDTO, error)
	AuthenticateUserFunc  func(ctx context.Context, authRequest *models.AuthRequest) (*models.UserFetchDTO, error)
	UpdateUserFunc        func(ctx context.Context, id string, user *models.UserUpdateDTO) (*models.UserFetchDTO, error)
	SetDefaultProjectFunc func(ctx context.Context, id string, projectID string) (*models.UserFetchDTO, error)
	GetUsersFunc          func(ctx context.Context) ([]*models.UserFetchDTO, error)
	ListUsersByTeamIDFunc func(ctx context.Context, teamID string) ([]*models.UserFetchDTO, error)
	DeleteUserFunc        func(ctx context.Context, id string) error
	ChangePasswordFunc    func(ctx context.Context, id string, request *models.ChangePasswordRequest) error
	ChangeImageFunc       func(ctx context.Context, id string, image *multipart.FileHeader) error
}

func (m *MockUserService) RegisterUser(ctx context.Context, user *models.UserRegisterDTO) error {
	if m.RegisterUserFunc != nil {
		return m.RegisterUserFunc(ctx, user)
	}
	return nil
}

func (m *MockUserService) GetUserByID(ctx context.Context, id string) (*models.UserFetchDTO, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockUserService) GetUserByUsername(ctx context.Context, username string) (*models.UserFetchDTO, error) {
	if m.GetUserByUsernameFunc != nil {
		return m.GetUserByUsernameFunc(ctx, username)
	}
	return nil, nil
}

func (m *MockUserService) AuthenticateUser(ctx context.Context, authRequest *models.AuthRequest) (*models.UserFetchDTO, error) {
	if m.AuthenticateUserFunc != nil {
		return m.AuthenticateUserFunc(ctx, authRequest)
	}
	return nil, nil
}

func (m *MockUserService) UpdateUser(ctx context.Context, id string, user *models.UserUpdateDTO) (*models.UserFetchDTO, error) {
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(ctx, id, user)
	}
	return nil, nil
}

func (m *MockUserService) SetDefaultProject(ctx context.Context, id string, projectID string) (*models.UserFetchDTO, error) {
	if m.SetDefaultProjectFunc != nil {
		return m.SetDefaultProjectFunc(ctx, id, projectID)
	}
	return nil, nil
}

func (m *MockUserService) GetUsers(ctx context.Context) ([]*models.UserFetchDTO, error) {
	if m.GetUsersFunc != nil {
		return m.GetUsersFunc(ctx)
	}
	return nil, nil
}

func (m *MockUserService) ListUsersByTeamID(ctx context.Context, teamID string) ([]*models.UserFetchDTO, error) {
	if m.ListUsersByTeamIDFunc != nil {
		return m.ListUsersByTeamIDFunc(ctx, teamID)
	}
	return nil, nil
}

func (m *MockUserService) DeleteUser(ctx context.Context, id string) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(ctx, id)
	}
	return nil
}

func (m *MockUserService) ChangePassword(ctx context.Context, id string, request *models.ChangePasswordRequest) error {
	if m.ChangePasswordFunc != nil {
		return m.ChangePasswordFunc(ctx, id, request)
	}
	return nil
}

func (m *MockUserService) ChangeImage(ctx context.Context, id string, image *multipart.FileHeader) error {
	if m.ChangeImageFunc != nil {
		return m.ChangeImageFunc(ctx, id, image)
	}
	return nil
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			AuthenticateUserFunc: func(ctx context.Context, authRequest *models.AuthRequest) (*models.UserFetchDTO, error) {
				return &models.UserFetchDTO{ID: "user-1", Username: "testuser"}, nil
			},
		}
		handler := NewAuthHandler(mockService, config.JWTConfig{SecretKey: "secret"})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := `{"username": "testuser", "password": "password"}`
		c.Request = httptest.NewRequest("POST", "/login", bytes.NewBufferString(body))

		handler.Login(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("InvalidPayload", func(t *testing.T) {
		handler := NewAuthHandler(&MockUserService{}, config.JWTConfig{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/login", bytes.NewBufferString("invalid"))

		handler.Login(c)

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
		handler := NewAuthHandler(mockService, config.JWTConfig{})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := `{"username": "testuser", "password": "wrongpassword"}`
		c.Request = httptest.NewRequest("POST", "/login", bytes.NewBufferString(body))

		handler.Login(c)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
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
		handler := NewAuthHandler(mockService, config.JWTConfig{})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := `{"email": "test@example.com", "password": "password", "confirmPassword": "password", "firstName": "Test", "lastName": "User", "username": "testuser", "role": "user"}`
		c.Request = httptest.NewRequest("POST", "/register", bytes.NewBufferString(body))

		handler.Register(c)

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
		handler := NewAuthHandler(mockService, config.JWTConfig{})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := `{"email": "test@example.com", "password": "password", "confirmPassword": "password", "firstName": "Test", "lastName": "User", "username": "testuser", "role": "user"}`
		c.Request = httptest.NewRequest("POST", "/register", bytes.NewBufferString(body))

		handler.Register(c)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}

func TestAuthHandler_SignOut(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(&MockUserService{}, config.JWTConfig{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/signout", nil)

	handler.SignOut(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
