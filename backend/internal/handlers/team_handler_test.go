package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/notFil/cspotlight/internal/middleware"
	"github.com/notFil/cspotlight/internal/models"
)

type MockTeamService struct {
	CreateTeamFunc  func(ctx context.Context, team *models.TeamUpsertDTO) error
	GetTeamByIDFunc func(ctx context.Context, id uuid.UUID) (*models.TeamFetchDTO, error)
	UpdateTeamFunc  func(ctx context.Context, id uuid.UUID, team *models.TeamUpsertDTO) (*models.TeamFetchDTO, error)
	DeleteTeamFunc  func(ctx context.Context, id uuid.UUID) error
	ListTeamsFunc   func(ctx context.Context) ([]*models.TeamFetchDTO, error)
}

func (m *MockTeamService) CreateTeam(ctx context.Context, team *models.TeamUpsertDTO) error {
	if m.CreateTeamFunc != nil {
		return m.CreateTeamFunc(ctx, team)
	}
	return nil
}

func (m *MockTeamService) GetTeamByID(ctx context.Context, id uuid.UUID) (*models.TeamFetchDTO, error) {
	if m.GetTeamByIDFunc != nil {
		return m.GetTeamByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockTeamService) UpdateTeam(ctx context.Context, id uuid.UUID, team *models.TeamUpsertDTO) (*models.TeamFetchDTO, error) {
	if m.UpdateTeamFunc != nil {
		return m.UpdateTeamFunc(ctx, id, team)
	}
	return nil, nil
}

func (m *MockTeamService) DeleteTeam(ctx context.Context, id uuid.UUID) error {
	if m.DeleteTeamFunc != nil {
		return m.DeleteTeamFunc(ctx, id)
	}
	return nil
}

func (m *MockTeamService) ListTeams(ctx context.Context) ([]*models.TeamFetchDTO, error) {
	if m.ListTeamsFunc != nil {
		return m.ListTeamsFunc(ctx)
	}
	return nil, nil
}

func TestTeamHandler_GetTeamByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	teamID := uuid.New()
	t.Run("Success", func(t *testing.T) {
		mockService := &MockTeamService{
			GetTeamByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.TeamFetchDTO, error) {
				return &models.TeamFetchDTO{ID: id, Name: "Test Team"}, nil
			},
		}
		handler := NewTeamHandler(mockService)

		w := httptest.NewRecorder()
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.GET("/teams/:teamID", handler.GetTeamByID)

		req := httptest.NewRequest("GET", "/teams/"+teamID.String(), nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockService := &MockTeamService{
			GetTeamByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.TeamFetchDTO, error) {
				return nil, errors.New("team not found")
			},
		}
		handler := NewTeamHandler(mockService)

		w := httptest.NewRecorder()
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.GET("/teams/:teamID", handler.GetTeamByID)

		req := httptest.NewRequest("GET", "/teams/"+teamID.String(), nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}

func TestTeamHandler_CreateTeam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockTeamService{
			CreateTeamFunc: func(ctx context.Context, team *models.TeamUpsertDTO) error {
				return nil
			},
		}
		handler := NewTeamHandler(mockService)

		w := httptest.NewRecorder()
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.POST("/teams", handler.CreateTeam)

		body := `{"name": "New Team"}`
		req := httptest.NewRequest("POST", "/teams", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}
	})

	t.Run("InvalidPayload", func(t *testing.T) {
		handler := NewTeamHandler(&MockTeamService{})
		w := httptest.NewRecorder()
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.POST("/teams", handler.CreateTeam)

		req := httptest.NewRequest("POST", "/teams", bytes.NewBufferString("invalid"))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestTeamHandler_UpdateTeam(t *testing.T) {
	gin.SetMode(gin.TestMode)
	teamID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockService := &MockTeamService{
			UpdateTeamFunc: func(ctx context.Context, id uuid.UUID, team *models.TeamUpsertDTO) (*models.TeamFetchDTO, error) {
				return &models.TeamFetchDTO{ID: id, Name: team.Name}, nil
			},
		}
		handler := NewTeamHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "teamID", Value: teamID.String()}}
		body := `{"name": "Updated Team"}`
		c.Request = httptest.NewRequest("PUT", "/teams/"+teamID.String(), bytes.NewBufferString(body))

		handler.UpdateTeam(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestTeamHandler_DeleteTeam(t *testing.T) {
	gin.SetMode(gin.TestMode)
	teamID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockService := &MockTeamService{
			DeleteTeamFunc: func(ctx context.Context, id uuid.UUID) error {
				return nil
			},
		}
		handler := NewTeamHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "teamID", Value: teamID.String()}}
		c.Request = httptest.NewRequest("DELETE", "/teams/"+teamID.String(), nil)

		handler.DeleteTeam(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestTeamHandler_ListTeams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockTeamService{
			ListTeamsFunc: func(ctx context.Context) ([]*models.TeamFetchDTO, error) {
				return []*models.TeamFetchDTO{{ID: uuid.New()}}, nil
			},
		}
		handler := NewTeamHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/teams", nil)

		handler.ListTeams(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}
