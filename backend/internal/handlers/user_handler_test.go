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
	"github.com/golang-jwt/jwt/v5"
	"github.com/notFil/cspotlight/internal/auth"
	"github.com/notFil/cspotlight/internal/constants"
	"github.com/notFil/cspotlight/internal/models"
)

func TestUserHandler_GetUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			GetUserByIDFunc: func(ctx context.Context, id string) (*models.UserFetchDTO, error) {
				return &models.UserFetchDTO{ID: id, Username: "testuser"}, nil
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "user-1"}}
		c.Set(constants.ClaimsContextKey, &auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"}})
		c.Request = httptest.NewRequest("GET", "/users/user-1", nil)

		handler.GetUserByID(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestUserHandler_UpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			UpdateUserFunc: func(ctx context.Context, id string, user *models.UserUpdateDTO) (*models.UserFetchDTO, error) {
				return &models.UserFetchDTO{ID: id, Username: user.Username}, nil
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "user-1"}}
		c.Set(constants.ClaimsContextKey, &auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"}})
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
			DeleteUserFunc: func(ctx context.Context, id string) error {
				return nil
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "user-1"}}
		c.Set(constants.ClaimsContextKey, &auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"}})
		c.Request = httptest.NewRequest("DELETE", "/users/user-1", nil)

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
			GetUsersFunc: func(ctx context.Context) ([]*models.UserFetchDTO, error) {
				return []*models.UserFetchDTO{{ID: "user-1"}}, nil
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(constants.ClaimsContextKey, &auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"}})
		c.Request = httptest.NewRequest("GET", "/users", nil)

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
			ListUsersByTeamIDFunc: func(ctx context.Context, teamID string) ([]*models.UserFetchDTO, error) {
				return []*models.UserFetchDTO{{ID: "user-1"}}, nil
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "teamID", Value: "team-1"}}
		c.Set(constants.ClaimsContextKey, &auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"}})
		c.Request = httptest.NewRequest("GET", "/teams/team-1/users", nil)

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
			SetDefaultProjectFunc: func(ctx context.Context, id string, projectID string) (*models.UserFetchDTO, error) {
				return &models.UserFetchDTO{ID: id, DefaultProjectID: projectID}, nil
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(constants.ClaimsContextKey, &auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"}})
		body := `{"projectID": "proj-1"}`
		c.Request = httptest.NewRequest("PATCH", "/users/user-1/default-project", bytes.NewBufferString(body))

		handler.SetDefaultProject(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockService := &MockUserService{
			SetDefaultProjectFunc: func(ctx context.Context, id string, projectID string) (*models.UserFetchDTO, error) {
				return nil, errors.New("failed")
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(constants.ClaimsContextKey, &auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"}})
		body := `{"projectID": "proj-1"}`
		c.Request = httptest.NewRequest("POST", "/users/default-project", bytes.NewBufferString(body))

		handler.SetDefaultProject(c)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}

func TestUserHandler_ChangeImage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			ChangeImageFunc: func(ctx context.Context, id string, image *multipart.FileHeader) error {
				return nil
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(constants.ClaimsContextKey, &auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"}})

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("image", "test.jpg")
		part.Write([]byte("image content"))
		writer.Close()

		c.Request = httptest.NewRequest("PATCH", "/users/user-1/image", body)
		c.Request.Header.Set("Content-Type", writer.FormDataContentType())

		handler.ChangeImage(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}
