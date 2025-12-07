package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/models"
)

type MockTeamService struct {
	CreateTeamFunc  func(ctx context.Context, team *models.TeamUpsertDTO) error
	GetTeamByIDFunc func(ctx context.Context, id string) (*models.TeamFetchDTO, error)
	UpdateTeamFunc  func(ctx context.Context, id string, team *models.TeamUpsertDTO) (*models.TeamFetchDTO, error)
	DeleteTeamFunc  func(ctx context.Context, id string) error
	ListTeamsFunc   func(ctx context.Context) ([]*models.TeamFetchDTO, error)
}

func (m *MockTeamService) CreateTeam(ctx context.Context, team *models.TeamUpsertDTO) error {
	if m.CreateTeamFunc != nil {
		return m.CreateTeamFunc(ctx, team)
	}
	return nil
}

func (m *MockTeamService) GetTeamByID(ctx context.Context, id string) (*models.TeamFetchDTO, error) {
	if m.GetTeamByIDFunc != nil {
		return m.GetTeamByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockTeamService) UpdateTeam(ctx context.Context, id string, team *models.TeamUpsertDTO) (*models.TeamFetchDTO, error) {
	if m.UpdateTeamFunc != nil {
		return m.UpdateTeamFunc(ctx, id, team)
	}
	return nil, nil
}

func (m *MockTeamService) DeleteTeam(ctx context.Context, id string) error {
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

	t.Run("Success", func(t *testing.T) {
		mockService := &MockTeamService{
			GetTeamByIDFunc: func(ctx context.Context, id string) (*models.TeamFetchDTO, error) {
				return &models.TeamFetchDTO{ID: id, Name: "Test Team"}, nil
			},
		}
		handler := NewTeamHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "team-1"}}
		c.Request = httptest.NewRequest("GET", "/teams/team-1", nil)

		handler.GetTeamByID(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockService := &MockTeamService{
			GetTeamByIDFunc: func(ctx context.Context, id string) (*models.TeamFetchDTO, error) {
				return nil, errors.New("team not found")
			},
		}
		handler := NewTeamHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "team-1"}}
		c.Request = httptest.NewRequest("GET", "/teams/team-1", nil)

		handler.GetTeamByID(c)

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
		c, _ := gin.CreateTestContext(w)
		body := `{"name": "New Team"}`
		c.Request = httptest.NewRequest("POST", "/teams", bytes.NewBufferString(body))

		handler.CreateTeam(c)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}
	})

	t.Run("InvalidPayload", func(t *testing.T) {
		handler := NewTeamHandler(&MockTeamService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/teams", bytes.NewBufferString("invalid"))

		handler.CreateTeam(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestTeamHandler_UpdateTeam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockTeamService{
			UpdateTeamFunc: func(ctx context.Context, id string, team *models.TeamUpsertDTO) (*models.TeamFetchDTO, error) {
				return &models.TeamFetchDTO{ID: id, Name: team.Name}, nil
			},
		}
		handler := NewTeamHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "team-1"}}
		body := `{"name": "Updated Team"}`
		c.Request = httptest.NewRequest("PUT", "/teams/team-1", bytes.NewBufferString(body))

		handler.UpdateTeam(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestTeamHandler_DeleteTeam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockTeamService{
			DeleteTeamFunc: func(ctx context.Context, id string) error {
				return nil
			},
		}
		handler := NewTeamHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "team-1"}}
		c.Request = httptest.NewRequest("DELETE", "/teams/team-1", nil)

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
				return []*models.TeamFetchDTO{{ID: "team-1"}}, nil
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
