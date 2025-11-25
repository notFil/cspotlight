package services

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/notFil/cspotlight/internal/models"
)

// MockTeamRepository is a manual mock for TeamRepository
type MockTeamRepository struct {
	teams map[string]*models.Team
}

func NewMockTeamRepository() *MockTeamRepository {
	return &MockTeamRepository{
		teams: make(map[string]*models.Team),
	}
}

func (m *MockTeamRepository) CreateTeam(team *models.Team) error {
	if team.ID == "" {
		team.ID = uuid.New().String()
	}
	if _, exists := m.teams[team.ID]; exists {
		return errors.New("team already exists")
	}
	m.teams[team.ID] = team
	return nil
}

func (m *MockTeamRepository) GetTeamByID(id string) (*models.Team, error) {
	if team, exists := m.teams[id]; exists {
		return team, nil
	}
	return nil, errors.New("team not found")
}

func (m *MockTeamRepository) UpdateTeam(team *models.Team) error {
	if _, exists := m.teams[team.ID]; exists {
		m.teams[team.ID] = team
		return nil
	}
	return errors.New("team not found")
}

func (m *MockTeamRepository) DeleteTeam(id string) error {
	if _, exists := m.teams[id]; exists {
		delete(m.teams, id)
		return nil
	}
	return errors.New("team not found")
}

func (m *MockTeamRepository) ListTeams() ([]*models.Team, error) {
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

	success, err := service.CreateTeam(teamDTO)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !success {
		t.Fatalf("expected success to be true")
	}

	// Verify team creation
	teams, err := mockRepo.ListTeams()
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

	teamID := "team-1"
	team := &models.Team{
		ID:   teamID,
		Name: "Test Team",
	}
	mockRepo.teams[teamID] = team

	fetchedTeam, err := service.GetTeamByID(teamID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if fetchedTeam.ID != teamID {
		t.Errorf("expected team ID %s, got %s", teamID, fetchedTeam.ID)
	}

	_, err = service.GetTeamByID("non-existent")
	if err == nil {
		t.Fatal("expected error for non-existent team, got nil")
	}
}

func TestUpdateTeam(t *testing.T) {
	mockRepo := NewMockTeamRepository()
	service := NewTeamService(mockRepo)

	teamID := "team-update"
	team := &models.Team{
		ID:   teamID,
		Name: "Original Name",
	}
	mockRepo.teams[teamID] = team

	updateDTO := &models.TeamUpsertDTO{
		Name: "Updated Name",
	}

	updatedTeam, err := service.UpdateTeam(teamID, updateDTO)
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

	teamID := "team-delete"
	team := &models.Team{
		ID:   teamID,
		Name: "To Delete",
	}
	mockRepo.teams[teamID] = team

	success, err := service.DeleteTeam(teamID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !success {
		t.Fatalf("expected success to be true")
	}

	_, err = mockRepo.GetTeamByID(teamID)
	if err == nil {
		t.Fatal("expected error getting deleted team, got nil")
	}
}
