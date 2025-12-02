package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/pkg/auth"
	"github.com/notFil/cspotlight/pkg/constants"
)

type MockProjectService struct {
	GetProjectByIDFunc func(id string, claims auth.Claims) (*models.ProjectFetchDTO, error)
	CreateProjectFunc  func(project *models.ProjectUpsertDTO, claims auth.Claims) (bool, error)
	UpdateProjectFunc  func(id string, project *models.ProjectUpsertDTO, claims auth.Claims) (*models.ProjectFetchDTO, error)
	DeleteProjectFunc  func(id string, claims auth.Claims) error
	ListProjectsFunc   func(claims auth.Claims) ([]*models.ProjectFetchDTO, error)
}

func (m *MockProjectService) GetProjectByID(id string, claims auth.Claims) (*models.ProjectFetchDTO, error) {
	if m.GetProjectByIDFunc != nil {
		return m.GetProjectByIDFunc(id, claims)
	}
	return nil, nil
}

func (m *MockProjectService) CreateProject(project *models.ProjectUpsertDTO, claims auth.Claims) (bool, error) {
	if m.CreateProjectFunc != nil {
		return m.CreateProjectFunc(project, claims)
	}
	return true, nil
}

func (m *MockProjectService) UpdateProject(id string, project *models.ProjectUpsertDTO, claims auth.Claims) (*models.ProjectFetchDTO, error) {
	if m.UpdateProjectFunc != nil {
		return m.UpdateProjectFunc(id, project, claims)
	}
	return nil, nil
}

func (m *MockProjectService) DeleteProject(id string, claims auth.Claims) error {
	if m.DeleteProjectFunc != nil {
		return m.DeleteProjectFunc(id, claims)
	}
	return nil
}

func (m *MockProjectService) ListProjects(claims auth.Claims) ([]*models.ProjectFetchDTO, error) {
	if m.ListProjectsFunc != nil {
		return m.ListProjectsFunc(claims)
	}
	return nil, nil
}

func TestProjectHandler_GetProjectByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockProjectService{
			GetProjectByIDFunc: func(id string, claims auth.Claims) (*models.ProjectFetchDTO, error) {
				return &models.ProjectFetchDTO{ID: id, Name: "Test Project"}, nil
			},
		}
		handler := NewProjectHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "proj-1"}}
		c.Set(constants.ClaimsContextKey, &auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"}})

		handler.GetProjectByID(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService := &MockProjectService{
			GetProjectByIDFunc: func(id string, claims auth.Claims) (*models.ProjectFetchDTO, error) {
				return nil, nil
			},
		}
		handler := NewProjectHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "proj-1"}}
		c.Set(constants.ClaimsContextKey, &auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"}})

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
			CreateProjectFunc: func(project *models.ProjectUpsertDTO, claims auth.Claims) (bool, error) {
				return true, nil
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
			UpdateProjectFunc: func(id string, project *models.ProjectUpsertDTO, claims auth.Claims) (*models.ProjectFetchDTO, error) {
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
			DeleteProjectFunc: func(id string, claims auth.Claims) error {
				return nil
			},
		}
		handler := NewProjectHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "proj-1"}}
		c.Set(constants.ClaimsContextKey, &auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"}})

		handler.DeleteProject(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("Forbidden", func(t *testing.T) {
		mockService := &MockProjectService{
			DeleteProjectFunc: func(id string, claims auth.Claims) error {
				return errors.New("unauthorized")
			},
		}
		handler := NewProjectHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "proj-1"}}
		c.Set(constants.ClaimsContextKey, &auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"}})

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
			ListProjectsFunc: func(claims auth.Claims) ([]*models.ProjectFetchDTO, error) {
				return []*models.ProjectFetchDTO{{ID: "proj-1"}}, nil
			},
		}
		handler := NewProjectHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(constants.ClaimsContextKey, &auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1"}})

		handler.ListProjects(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}
