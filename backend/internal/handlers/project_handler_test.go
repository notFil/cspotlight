package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/notFil/cspotlight/internal/auth"
	"github.com/notFil/cspotlight/internal/constants"
	"github.com/notFil/cspotlight/internal/models"
)

type MockProjectService struct {
	GetProjectByIDFunc func(ctx context.Context, id string) (*models.ProjectFetchDTO, error)
	CreateProjectFunc  func(ctx context.Context, project *models.ProjectUpsertDTO) error
	UpdateProjectFunc  func(ctx context.Context, id string, project *models.ProjectUpsertDTO) (*models.ProjectFetchDTO, error)
	DeleteProjectFunc  func(ctx context.Context, id string) error
	ListProjectsFunc   func(ctx context.Context) ([]*models.ProjectFetchDTO, error)
}

func (m *MockProjectService) GetProjectByID(ctx context.Context, id string) (*models.ProjectFetchDTO, error) {
	if m.GetProjectByIDFunc != nil {
		return m.GetProjectByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockProjectService) CreateProject(ctx context.Context, project *models.ProjectUpsertDTO) error {
	if m.CreateProjectFunc != nil {
		return m.CreateProjectFunc(ctx, project)
	}
	return nil
}

func (m *MockProjectService) UpdateProject(ctx context.Context, id string, project *models.ProjectUpsertDTO) (*models.ProjectFetchDTO, error) {
	if m.UpdateProjectFunc != nil {
		return m.UpdateProjectFunc(ctx, id, project)
	}
	return nil, nil
}

func (m *MockProjectService) DeleteProject(ctx context.Context, id string) error {
	if m.DeleteProjectFunc != nil {
		return m.DeleteProjectFunc(ctx, id)
	}
	return nil
}

func (m *MockProjectService) ListProjects(ctx context.Context) ([]*models.ProjectFetchDTO, error) {
	if m.ListProjectsFunc != nil {
		return m.ListProjectsFunc(ctx)
	}
	return nil, nil
}

func TestProjectHandler_GetProjectByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockProjectService{
			GetProjectByIDFunc: func(ctx context.Context, id string) (*models.ProjectFetchDTO, error) {
				return &models.ProjectFetchDTO{ID: id, Name: "Test Project"}, nil
			},
		}
		handler := NewProjectHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "proj-1"}}
		c.Set(constants.ClaimsContextKey, &auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"}})
		c.Request = httptest.NewRequest("GET", "/projects/proj-1", nil)

		handler.GetProjectByID(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService := &MockProjectService{
			GetProjectByIDFunc: func(ctx context.Context, id string) (*models.ProjectFetchDTO, error) {
				return nil, nil
			},
		}
		handler := NewProjectHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "proj-1"}}
		c.Set(constants.ClaimsContextKey, &auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"}})
		c.Request = httptest.NewRequest("GET", "/projects/proj-1", nil)

		handler.GetProjectByID(c)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})
}

func TestProjectHandler_CreateProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockProjectService{
			CreateProjectFunc: func(ctx context.Context, project *models.ProjectUpsertDTO) error {
				return nil
			},
		}
		handler := NewProjectHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(constants.ClaimsContextKey, &auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"}})
		body := `{"name": "New Project", "teamId": "team-1"}`
		c.Request = httptest.NewRequest("POST", "/projects", bytes.NewBufferString(body))

		handler.CreateProject(c)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}
	})

	t.Run("InvalidPayload", func(t *testing.T) {
		handler := NewProjectHandler(&MockProjectService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(constants.ClaimsContextKey, &auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"}})
		c.Request = httptest.NewRequest("POST", "/projects", bytes.NewBufferString("invalid"))

		handler.CreateProject(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestProjectHandler_UpdateProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockProjectService{
			UpdateProjectFunc: func(ctx context.Context, id string, project *models.ProjectUpsertDTO) (*models.ProjectFetchDTO, error) {
				return &models.ProjectFetchDTO{ID: id, Name: project.Name}, nil
			},
		}
		handler := NewProjectHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "proj-1"}}
		c.Set(constants.ClaimsContextKey, &auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"}})
		body := `{"name": "Updated Project", "teamId": "team-1"}`
		c.Request = httptest.NewRequest("PUT", "/projects/proj-1", bytes.NewBufferString(body))

		handler.UpdateProject(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestProjectHandler_DeleteProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockProjectService{
			DeleteProjectFunc: func(ctx context.Context, id string) error {
				return nil
			},
		}
		handler := NewProjectHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "proj-1"}}
		c.Set(constants.ClaimsContextKey, &auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"}})
		c.Request = httptest.NewRequest("DELETE", "/projects/proj-1", nil)

		handler.DeleteProject(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("Forbidden", func(t *testing.T) {
		mockService := &MockProjectService{
			DeleteProjectFunc: func(ctx context.Context, id string) error {
				return errors.New("unauthorized")
			},
		}
		handler := NewProjectHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "proj-1"}}
		c.Set(constants.ClaimsContextKey, &auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"}})
		c.Request = httptest.NewRequest("DELETE", "/projects/proj-1", nil)

		handler.DeleteProject(c)

		if w.Code != http.StatusForbidden {
			t.Errorf("expected status 403, got %d", w.Code)
		}
	})
}

func TestProjectHandler_ListProjects(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockProjectService{
			ListProjectsFunc: func(ctx context.Context) ([]*models.ProjectFetchDTO, error) {
				return []*models.ProjectFetchDTO{{ID: "proj-1"}}, nil
			},
		}
		handler := NewProjectHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(constants.ClaimsContextKey, &auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"}})
		c.Request = httptest.NewRequest("GET", "/projects", nil)

		handler.ListProjects(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}
