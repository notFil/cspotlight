package handlers

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/config"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/pkg/auth"
	"github.com/notFil/cspotlight/pkg/logger"
)

func init() {
	logger.InitializeLogger("test")
}

type MockUserService struct {
	RegisterUserFunc      func(user *models.UserRegisterDTO) error
	GetUserByIDFunc       func(id string, claims auth.Claims) (*models.UserFetchDTO, error)
	GetUserByUsernameFunc func(username string) (*models.UserFetchDTO, error)
	AuthenticateUserFunc  func(authRequest *models.AuthRequest) (*models.UserFetchDTO, error)
	UpdateUserFunc        func(id string, user *models.UserUpdateDTO, claims auth.Claims) (*models.UserFetchDTO, error)
	SetDefaultProjectFunc func(id string, projectID string, claims auth.Claims) (*models.UserFetchDTO, error)
	GetUsersFunc          func(claims auth.Claims) ([]*models.UserFetchDTO, error)
	ListUsersByTeamIDFunc func(teamID string, claims auth.Claims) ([]*models.UserFetchDTO, error)
	DeleteUserFunc        func(id string, claims auth.Claims) error
	ChangePasswordFunc    func(id string, request *models.ChangePasswordRequest, claims auth.Claims) error
	ChangeImageFunc       func(id string, image *multipart.FileHeader, claims auth.Claims) error
}

func (m *MockUserService) RegisterUser(user *models.UserRegisterDTO) error {
	if m.RegisterUserFunc != nil {
		return m.RegisterUserFunc(user)
	}
	return nil
}

func (m *MockUserService) GetUserByID(id string, claims auth.Claims) (*models.UserFetchDTO, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(id, claims)
	}
	return nil, nil
}

func (m *MockUserService) GetUserByUsername(username string) (*models.UserFetchDTO, error) {
	if m.GetUserByUsernameFunc != nil {
		return m.GetUserByUsernameFunc(username)
	}
	return nil, nil
}

func (m *MockUserService) AuthenticateUser(authRequest *models.AuthRequest) (*models.UserFetchDTO, error) {
	if m.AuthenticateUserFunc != nil {
		return m.AuthenticateUserFunc(authRequest)
	}
	return nil, nil
}

func (m *MockUserService) UpdateUser(id string, user *models.UserUpdateDTO, claims auth.Claims) (*models.UserFetchDTO, error) {
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(id, user, claims)
	}
	return nil, nil
}

func (m *MockUserService) SetDefaultProject(id string, projectID string, claims auth.Claims) (*models.UserFetchDTO, error) {
	if m.SetDefaultProjectFunc != nil {
		return m.SetDefaultProjectFunc(id, projectID, claims)
	}
	return nil, nil
}

func (m *MockUserService) GetUsers(claims auth.Claims) ([]*models.UserFetchDTO, error) {
	if m.GetUsersFunc != nil {
		return m.GetUsersFunc(claims)
	}
	return nil, nil
}

func (m *MockUserService) ListUsersByTeamID(teamID string, claims auth.Claims) ([]*models.UserFetchDTO, error) {
	if m.ListUsersByTeamIDFunc != nil {
		return m.ListUsersByTeamIDFunc(teamID, claims)
	}
	return nil, nil
}

func (m *MockUserService) DeleteUser(id string, claims auth.Claims) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(id, claims)
	}
	return nil
}

func (m *MockUserService) ChangePassword(id string, request *models.ChangePasswordRequest, claims auth.Claims) error {
	if m.ChangePasswordFunc != nil {
		return m.ChangePasswordFunc(id, request, claims)
	}
	return nil
}

func (m *MockUserService) ChangeImage(id string, image *multipart.FileHeader, claims auth.Claims) error {
	if m.ChangeImageFunc != nil {
		return m.ChangeImageFunc(id, image, claims)
	}
	return nil
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			AuthenticateUserFunc: func(authRequest *models.AuthRequest) (*models.UserFetchDTO, error) {
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
			AuthenticateUserFunc: func(authRequest *models.AuthRequest) (*models.UserFetchDTO, error) {
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
			RegisterUserFunc: func(user *models.UserRegisterDTO) error {
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
			RegisterUserFunc: func(user *models.UserRegisterDTO) error {
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
