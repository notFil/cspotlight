package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/pkg/auth"
	"github.com/notFil/cspotlight/pkg/constants"
)

func TestUserHandler_GetUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			GetUserByIDFunc: func(id string, claims auth.Claims) (*models.UserFetchDTO, error) {
				return &models.UserFetchDTO{ID: id, Username: "testuser"}, nil
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "user-1"}}
		c.Set(constants.ClaimsContextKey, &auth.Claims{UserID: "user-1"})

		handler.GetUserByID(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestUserHandler_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			RegisterUserFunc: func(user *models.UserCreateDTO) (bool, error) {
				return true, nil
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := `{"email": "newuser@example.com", "password": "password", "firstName": "New", "lastName": "User", "username": "newuser", "role": "user"}`
		c.Request = httptest.NewRequest("POST", "/users", bytes.NewBufferString(body))

		handler.CreateUser(c)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}
	})
}

func TestUserHandler_UpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			UpdateUserFunc: func(id string, user *models.UserUpdateDTO, claims auth.Claims) (*models.UserFetchDTO, error) {
				return &models.UserFetchDTO{ID: id, Username: user.Username}, nil
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "user-1"}}
		c.Set(constants.ClaimsContextKey, &auth.Claims{UserID: "user-1"})
		body := `{"username": "updateduser", "firstName": "Updated", "lastName": "User", "email": "updated@example.com", "role": "user"}`
		c.Request = httptest.NewRequest("PUT", "/users/user-1", bytes.NewBufferString(body))

		handler.UpdateUser(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestUserHandler_DeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			DeleteUserFunc: func(id string, claims auth.Claims) error {
				return nil
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "user-1"}}
		c.Set(constants.ClaimsContextKey, &auth.Claims{UserID: "user-1"})

		handler.DeleteUser(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestUserHandler_ListUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			GetUsersFunc: func(claims auth.Claims) ([]*models.UserFetchDTO, error) {
				return []*models.UserFetchDTO{{ID: "user-1"}}, nil
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(constants.ClaimsContextKey, &auth.Claims{UserID: "user-1"})

		handler.ListUsers(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestUserHandler_ListUsersByTeamID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			ListUsersByTeamIDFunc: func(teamID string, claims auth.Claims) ([]*models.UserFetchDTO, error) {
				return []*models.UserFetchDTO{{ID: "user-1"}}, nil
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "teamID", Value: "team-1"}}
		c.Set(constants.ClaimsContextKey, &auth.Claims{UserID: "user-1"})

		handler.ListUsersByTeamID(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestUserHandler_SetDefaultProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			SetDefaultProjectFunc: func(projectID string, claims *auth.Claims) error {
				return nil
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(constants.ClaimsContextKey, &auth.Claims{UserID: "user-1"})
		body := `{"projectID": "proj-1"}`
		c.Request = httptest.NewRequest("POST", "/users/default-project", bytes.NewBufferString(body))

		handler.SetDefaultProject(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockService := &MockUserService{
			SetDefaultProjectFunc: func(projectID string, claims *auth.Claims) error {
				return errors.New("failed")
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(constants.ClaimsContextKey, &auth.Claims{UserID: "user-1"})
		body := `{"projectID": "proj-1"}`
		c.Request = httptest.NewRequest("POST", "/users/default-project", bytes.NewBufferString(body))

		handler.SetDefaultProject(c)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}
