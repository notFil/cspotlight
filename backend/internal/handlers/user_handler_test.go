package handlers

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/notFil/cspotlight/internal/auth"
	"github.com/notFil/cspotlight/internal/middleware"
	"github.com/notFil/cspotlight/internal/models"
)

const staticPath = "../static"

func TestUserHandler_GetUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.UserFetchDTO, error) {
				return &models.UserFetchDTO{ID: id, Username: "testuser"}, nil
			},
		}
		mockProjectService := &MockProjectService{}
		handler := NewUserHandler(mockService, mockProjectService, staticPath)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userID", Value: userID.String()}}
		c.Request = httptest.NewRequest("GET", "/users/"+userID.String(), nil)
		userCtx := auth.NewUserContext(userID.String(), uuid.New().String(), "user")
		c.Request = c.Request.WithContext(auth.ContextWithUser(c.Request.Context(), userCtx))

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
		mockProjectService := &MockProjectService{}
		handler := NewUserHandler(mockService, mockProjectService, staticPath)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userID", Value: userID.String()}}
		body := `{"username": "updateduser", "firstName": "Updated", "lastName": "User", "email": "updated@example.com", "role": "user"}`
		c.Request = httptest.NewRequest("PUT", "/users/"+userID.String(), bytes.NewBufferString(body))
		userCtx := auth.NewUserContext(userID.String(), uuid.New().String(), "user")
		c.Request = c.Request.WithContext(auth.ContextWithUser(c.Request.Context(), userCtx))

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
		mockProjectService := &MockProjectService{}
		handler := NewUserHandler(mockService, mockProjectService, staticPath)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userID", Value: userID.String()}}
		c.Request = httptest.NewRequest("DELETE", "/users/"+userID.String(), nil)
		userCtx := auth.NewUserContext(userID.String(), uuid.New().String(), "user")
		c.Request = c.Request.WithContext(auth.ContextWithUser(c.Request.Context(), userCtx))

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
		mockProjectService := &MockProjectService{}
		handler := NewUserHandler(mockService, mockProjectService, staticPath)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/users", nil)
		userCtx := auth.NewUserContext(uuid.New().String(), uuid.New().String(), "user")
		c.Request = c.Request.WithContext(auth.ContextWithUser(c.Request.Context(), userCtx))

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
		mockProjectService := &MockProjectService{}
		handler := NewUserHandler(mockService, mockProjectService, staticPath)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "teamID", Value: teamID.String()}}
		c.Request = httptest.NewRequest("GET", "/teams/"+teamID.String()+"/users", nil)
		userCtx := auth.NewUserContext(uuid.NewString(), uuid.NewString(), "user")
		c.Request = c.Request.WithContext(auth.ContextWithUser(c.Request.Context(), userCtx))

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
				return &models.UserFetchDTO{ID: id, DefaultProjectID: &projectID}, nil
			},
		}
		mockProjectService := &MockProjectService{
			GetProjectByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.ProjectFetchDTO, error) {
				return &models.ProjectFetchDTO{ID: id}, nil
			},
		}
		handler := NewUserHandler(mockService, mockProjectService, staticPath)

		w := httptest.NewRecorder()
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.Use(func(c *gin.Context) {
			userCtx := auth.NewUserContext(userID.String(), uuid.New().String(), "user")
			c.Request = c.Request.WithContext(auth.ContextWithUser(c.Request.Context(), userCtx))
			c.Next()
		})
		router.PATCH("/users/:userID/default-project", handler.SetDefaultProject)

		body := `{"projectID": "` + projectID.String() + `"}`
		req := httptest.NewRequest("PATCH", "/users/"+userID.String()+"/default-project", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

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
		mockProjectService := &MockProjectService{
			GetProjectByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.ProjectFetchDTO, error) {
				return &models.ProjectFetchDTO{ID: id}, nil
			},
		}
		handler := NewUserHandler(mockService, mockProjectService, staticPath)

		w := httptest.NewRecorder()
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.Use(func(c *gin.Context) {
			userCtx := auth.NewUserContext(userID.String(), uuid.New().String(), "user")
			c.Request = c.Request.WithContext(auth.ContextWithUser(c.Request.Context(), userCtx))
			c.Next()
		})
		router.PATCH("/users/:userID/default-project", handler.SetDefaultProject)

		body := `{"projectID": "` + projectID.String() + `"}`
		req := httptest.NewRequest("PATCH", "/users/"+userID.String()+"/default-project", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

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
			ChangeImageFunc: func(ctx context.Context, id uuid.UUID, path string) error {
				return nil
			},
		}
		mockProjectService := &MockProjectService{}
		handler := NewUserHandler(mockService, mockProjectService, staticPath)

		w := httptest.NewRecorder()
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.Use(func(c *gin.Context) {
			userCtx := auth.NewUserContext(userID.String(), uuid.New().String(), "user")
			c.Request = c.Request.WithContext(auth.ContextWithUser(c.Request.Context(), userCtx))
			c.Next()
		})
		router.PATCH("/users/:userID/image", handler.ChangeImage)

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("image", "test.png")

		img := image.NewRGBA(image.Rect(0, 0, 1, 1))
		png.Encode(part, img)

		writer.Close()

		req := httptest.NewRequest("PATCH", "/users/"+userID.String()+"/image", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}
