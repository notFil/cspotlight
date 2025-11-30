package services

import (
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/pkg/auth"
)

// MockProjectRepository is a manual mock for ProjectRepository
type MockProjectRepository struct {
	projects map[string]*models.Project
}

func NewMockProjectRepository() *MockProjectRepository {
	return &MockProjectRepository{
		projects: make(map[string]*models.Project),
	}
}

func (m *MockProjectRepository) CreateProject(project *models.Project) error {
	if project.ID == "" {
		project.ID = uuid.New().String()
	}
	if _, exists := m.projects[project.ID]; exists {
		return errors.New("project already exists")
	}
	m.projects[project.ID] = project
	return nil
}

func (m *MockProjectRepository) GetProjectByID(id string) (*models.Project, error) {
	if project, exists := m.projects[id]; exists {
		return project, nil
	}
	return nil, errors.New("project not found")
}

func (m *MockProjectRepository) UpdateProject(project *models.Project) error {
	if _, exists := m.projects[project.ID]; exists {
		m.projects[project.ID] = project
		return nil
	}
	return errors.New("project not found")
}

func (m *MockProjectRepository) DeleteProject(id string) error {
	if _, exists := m.projects[id]; exists {
		delete(m.projects, id)
		return nil
	}
	return errors.New("project not found")
}

func (m *MockProjectRepository) ListProjects() ([]*models.ProjectFetchDTO, error) {
	var projects []*models.ProjectFetchDTO
	for _, project := range m.projects {
		projects = append(projects, project.ToFetchDTO())
	}
	return projects, nil
}

func (m *MockProjectRepository) ListProjectsByTeamID(teamID string) ([]*models.ProjectFetchDTO, error) {
	var projects []*models.ProjectFetchDTO
	for _, project := range m.projects {
		if project.TeamID == teamID {
			projects = append(projects, project.ToFetchDTO())
		}
	}
	return projects, nil
}

func TestCreateProject(t *testing.T) {
	mockRepo := NewMockProjectRepository()
	service := NewProjectService(mockRepo)

	projectDTO := &models.ProjectUpsertDTO{
		Name:        "Test Project",
		Description: "A test project",
		TeamID:      "team-1",
	}

	// Test as superadmin
	adminClaims := auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "admin"}, Role: "superadmin"}
	success, err := service.CreateProject(projectDTO, adminClaims)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !success {
		t.Fatalf("expected success to be true")
	}

	// Test as user (should force team ID)
	userTeamID := "user-team-id"
	userClaims := auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user"}, Role: "user", TeamID: userTeamID}
	projectDTOUser := &models.ProjectUpsertDTO{
		Name:        "User Project",
		Description: "A user project",
		TeamID:      "some-other-team", // Should be ignored/overwritten
	}

	success, err = service.CreateProject(projectDTOUser, userClaims)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !success {
		t.Fatalf("expected success to be true")
	}

	// Verify the user project was created with the correct team ID
	// Since we don't have the ID returned, we can list projects by team ID to verify
	projects, err := mockRepo.ListProjectsByTeamID(userTeamID)
	if err != nil {
		t.Fatalf("expected no error listing projects, got %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("expected 1 project for team %s, got %d", userTeamID, len(projects))
	}
	if projects[0].Name != "User Project" {
		t.Errorf("expected project name 'User Project', got %s", projects[0].Name)
	}
}

func TestGetProjectByID(t *testing.T) {
	mockRepo := NewMockProjectRepository()
	service := NewProjectService(mockRepo)

	projectID := "proj-1"
	teamID := "team-1"
	project := &models.Project{
		ID:     projectID,
		Name:   "Test Project",
		TeamID: teamID,
	}
	mockRepo.projects[projectID] = project

	// Test authorized access (same team)
	claims := auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user"}, Role: "user", TeamID: teamID}
	fetchedProject, err := service.GetProjectByID(projectID, claims)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if fetchedProject.ID != projectID {
		t.Errorf("expected project ID %s, got %s", projectID, fetchedProject.ID)
	}

	// Test unauthorized access (different team)
	otherTeamID := "team-2"
	otherClaims := auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "other"}, Role: "user", TeamID: otherTeamID}
	_, err = service.GetProjectByID(projectID, otherClaims)
	if err == nil {
		t.Fatal("expected unauthorized error, got nil")
	}

	// Test superadmin access
	adminClaims := auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "admin"}, Role: "superadmin"}
	fetchedProjectAdmin, err := service.GetProjectByID(projectID, adminClaims)
	if err != nil {
		t.Fatalf("expected no error for admin, got %v", err)
	}
	if fetchedProjectAdmin.ID != projectID {
		t.Errorf("expected project ID %s, got %s", projectID, fetchedProjectAdmin.ID)
	}
}

func TestUpdateProject(t *testing.T) {
	mockRepo := NewMockProjectRepository()
	service := NewProjectService(mockRepo)

	projectID := "proj-update"
	teamID := "team-1"
	project := &models.Project{
		ID:     projectID,
		Name:   "Original Name",
		TeamID: teamID,
	}
	mockRepo.projects[projectID] = project

	updateDTO := &models.ProjectUpsertDTO{
		Name: "Updated Name",
	}

	// Test authorized update
	claims := auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user"}, Role: "user", TeamID: teamID}
	updatedProject, err := service.UpdateProject(projectID, updateDTO, claims)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updatedProject.Name != "Updated Name" {
		t.Errorf("expected updated name, got %s", updatedProject.Name)
	}

	// Test unauthorized update
	otherTeamID := "team-2"
	otherClaims := auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "other"}, Role: "user", TeamID: otherTeamID}
	_, err = service.UpdateProject(projectID, updateDTO, otherClaims)
	if err == nil {
		t.Fatal("expected unauthorized error, got nil")
	}
}

func TestDeleteProject(t *testing.T) {
	mockRepo := NewMockProjectRepository()
	service := NewProjectService(mockRepo)

	projectID := "proj-delete"
	teamID := "team-1"
	project := &models.Project{
		ID:     projectID,
		Name:   "To Delete",
		TeamID: teamID,
	}
	mockRepo.projects[projectID] = project

	// Test unauthorized delete
	otherTeamID := "team-2"
	otherClaims := auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "other"}, Role: "user", TeamID: otherTeamID}
	err := service.DeleteProject(projectID, otherClaims)
	if err == nil {
		t.Fatal("expected unauthorized error, got nil")
	}

	// Test authorized delete
	claims := auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user"}, Role: "user", TeamID: teamID}
	err = service.DeleteProject(projectID, claims)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify deletion
	_, err = mockRepo.GetProjectByID(projectID)
	if err == nil {
		t.Fatal("expected error getting deleted project, got nil")
	}
}
