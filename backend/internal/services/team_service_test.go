package services

import (
	"context"
	"errors"
	"testing"

	"cspotlight/internal/models"

	"github.com/google/uuid"
)

type MockTeamRepository struct {
	teams map[uuid.UUID]*models.Team
}

func NewMockTeamRepository() *MockTeamRepository {
	return &MockTeamRepository{
		teams: make(map[uuid.UUID]*models.Team),
	}
}

func (m *MockTeamRepository) CreateTeam(ctx context.Context, team *models.Team) error {
	if team.ID == uuid.Nil {
		team.ID = uuid.New()
	}
	if _, exists := m.teams[team.ID]; exists {
		return errors.New("team already exists")
	}
	m.teams[team.ID] = team
	return nil
}

func (m *MockTeamRepository) GetTeamByID(ctx context.Context, id uuid.UUID) (*models.Team, error) {
	if team, exists := m.teams[id]; exists {
		return team, nil
	}
	return nil, errors.New("team not found")
}

func (m *MockTeamRepository) UpdateTeam(ctx context.Context, team *models.Team) error {
	if _, exists := m.teams[team.ID]; exists {
		m.teams[team.ID] = team
		return nil
	}
	return errors.New("team not found")
}

func (m *MockTeamRepository) DeleteTeam(ctx context.Context, id uuid.UUID) error {
	if _, exists := m.teams[id]; exists {
		delete(m.teams, id)
		return nil
	}
	return errors.New("team not found")
}

func (m *MockTeamRepository) ListTeams(ctx context.Context) ([]*models.Team, error) {
	var teams []*models.Team
	for _, team := range m.teams {
		teams = append(teams, team)
	}
	return teams, nil
}

func TestCreateTeam(t *testing.T) {
	mockRepo := NewMockTeamRepository()
	service := NewTeamService(mockRepo)

	teamDTO := &models.TeamUpsertDTO{
		Name:        "Test Team",
		Description: "A test team",
	}

	err := service.CreateTeam(context.Background(), teamDTO)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify team creation
	teams, err := mockRepo.ListTeams(context.Background())
	if err != nil {
		t.Fatalf("expected no error listing teams, got %v", err)
	}
	if len(teams) != 1 {
		t.Fatalf("expected 1 team, got %d", len(teams))
	}
	if teams[0].Name != "Test Team" {
		t.Errorf("expected team name 'Test Team', got %s", teams[0].Name)
	}
}

func TestGetTeamByID(t *testing.T) {
	mockRepo := NewMockTeamRepository()
	service := NewTeamService(mockRepo)

	teamID := uuid.New()
	team := &models.Team{
		ID:   teamID,
		Name: "Test Team",
	}
	mockRepo.teams[teamID] = team

	fetchedTeam, err := service.GetTeamByID(context.Background(), teamID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if fetchedTeam.ID != teamID {
		t.Errorf("expected team ID %s, got %s", teamID, fetchedTeam.ID)
	}

	_, err = service.GetTeamByID(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent team, got nil")
	}
}

func TestUpdateTeam(t *testing.T) {
	mockRepo := NewMockTeamRepository()
	service := NewTeamService(mockRepo)

	teamID := uuid.New()
	team := &models.Team{
		ID:   teamID,
		Name: "Original Name",
	}
	mockRepo.teams[teamID] = team

	updateDTO := &models.TeamUpsertDTO{
		Name: "Updated Name",
	}

	updatedTeam, err := service.UpdateTeam(context.Background(), teamID, updateDTO)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updatedTeam.Name != "Updated Name" {
		t.Errorf("expected updated name, got %s", updatedTeam.Name)
	}
}

func TestDeleteTeam(t *testing.T) {
	mockRepo := NewMockTeamRepository()
	service := NewTeamService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		teamID := uuid.New()
		team := &models.Team{
			ID:   teamID,
			Name: "To Delete",
		}
		mockRepo.teams[teamID] = team

		err := service.DeleteTeam(context.Background(), teamID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		_, err = mockRepo.GetTeamByID(context.Background(), teamID)
		if err == nil {
			t.Fatal("expected error getting deleted team, got nil")
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		err := service.DeleteTeam(context.Background(), uuid.New())
		if err == nil {
			t.Fatal("expected error for non-existent team, got nil")
		}
	})
}
