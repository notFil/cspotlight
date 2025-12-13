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
	"github.com/google/uuid"
	"github.com/notFil/cspotlight/internal/auth"
	"github.com/notFil/cspotlight/internal/constants"
	"github.com/notFil/cspotlight/internal/models"
)

func TestUserHandler_GetUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.UserFetchDTO, error) {
				return &models.UserFetchDTO{ID: id, Username: "testuser"}, nil
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: userID.String()}}
		c.Set(constants.ClaimsContextKey, &auth.AuthContext{Subject: userID})
		c.Request = httptest.NewRequest("GET", "/users/"+userID.String(), nil)

		handler.GetUserByID(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestUserHandler_UpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			UpdateUserFunc: func(ctx context.Context, id uuid.UUID, user *models.UserUpdateDTO) (*models.UserFetchDTO, error) {
				return &models.UserFetchDTO{ID: id, Username: user.Username}, nil
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: userID.String()}}
		c.Set(constants.ClaimsContextKey, &auth.AuthContext{Subject: userID})
		body := `{"username": "updateduser", "firstName": "Updated", "lastName": "User", "email": "updated@example.com", "role": "user"}`
		c.Request = httptest.NewRequest("PUT", "/users/"+userID.String(), bytes.NewBufferString(body))

		handler.UpdateUser(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestUserHandler_DeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			DeleteUserFunc: func(ctx context.Context, id uuid.UUID) error {
				return nil
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: userID.String()}}
		c.Set(constants.ClaimsContextKey, &auth.AuthContext{Subject: userID})
		c.Request = httptest.NewRequest("DELETE", "/users/"+userID.String(), nil)

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
				return []*models.UserFetchDTO{{ID: uuid.New()}}, nil
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(constants.ClaimsContextKey, &auth.AuthContext{Subject: uuid.New()})
		c.Request = httptest.NewRequest("GET", "/users", nil)

		handler.ListUsers(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestUserHandler_ListUsersByTeamID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	teamID := uuid.New()
	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			ListUsersByTeamIDFunc: func(ctx context.Context, teamID uuid.UUID) ([]*models.UserFetchDTO, error) {
				return []*models.UserFetchDTO{{ID: uuid.New()}}, nil
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "teamID", Value: teamID.String()}}
		c.Set(constants.ClaimsContextKey, &auth.AuthContext{Subject: uuid.New()})
		c.Request = httptest.NewRequest("GET", "/teams/"+teamID.String()+"/users", nil)

		handler.ListUsersByTeamID(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestUserHandler_SetDefaultProject(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	projectID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			SetDefaultProjectFunc: func(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (*models.UserFetchDTO, error) {
				return &models.UserFetchDTO{ID: id, DefaultProjectID: projectID}, nil
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(constants.ClaimsContextKey, &auth.AuthContext{Subject: userID})
		c.Params = gin.Params{{Key: "id", Value: userID.String()}}
		body := `{"projectID": "` + projectID.String() + `"}`
		c.Request = httptest.NewRequest("PATCH", "/users/"+userID.String()+"/default-project", bytes.NewBufferString(body))

		handler.SetDefaultProject(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockService := &MockUserService{
			SetDefaultProjectFunc: func(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (*models.UserFetchDTO, error) {
				return nil, errors.New("failed")
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(constants.ClaimsContextKey, &auth.AuthContext{Subject: userID})
		c.Params = gin.Params{{Key: "id", Value: userID.String()}}
		body := `{"projectID": "` + projectID.String() + `"}`
		c.Request = httptest.NewRequest("PATCH", "/users/default-project", bytes.NewBufferString(body))

		handler.SetDefaultProject(c)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}

func TestUserHandler_ChangeImage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			ChangeImageFunc: func(ctx context.Context, id uuid.UUID, image *multipart.FileHeader) error {
				return nil
			},
		}
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: userID.String()}}
		c.Set(constants.ClaimsContextKey, &auth.AuthContext{Subject: userID})

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("image", "test.jpg")
		part.Write([]byte("image content"))
		writer.Close()

		c.Request = httptest.NewRequest("PATCH", "/users/"+userID.String()+"/image", body)
		c.Request.Header.Set("Content-Type", writer.FormDataContentType())

		handler.ChangeImage(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}
