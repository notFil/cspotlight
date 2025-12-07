package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/notFil/cspotlight/internal/auth"
	"github.com/notFil/cspotlight/internal/models"
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

func (m *MockProjectRepository) CreateProject(ctx context.Context, project *models.Project) error {
	if project.ID == "" {
		project.ID = uuid.New().String()
	}
	if _, exists := m.projects[project.ID]; exists {
		return errors.New("project already exists")
	}
	m.projects[project.ID] = project
	return nil
}

func (m *MockProjectRepository) GetProjectByID(ctx context.Context, id string) (*models.Project, error) {
	if project, exists := m.projects[id]; exists {
		return project, nil
	}
	return nil, errors.New("project not found")
}

func (m *MockProjectRepository) UpdateProject(ctx context.Context, project *models.Project) error {
	if _, exists := m.projects[project.ID]; exists {
		m.projects[project.ID] = project
		return nil
	}
	return errors.New("project not found")
}

func (m *MockProjectRepository) DeleteProject(ctx context.Context, id string) error {
	if _, exists := m.projects[id]; exists {
		delete(m.projects, id)
		return nil
	}
	return errors.New("project not found")
}

func (m *MockProjectRepository) ListProjects(ctx context.Context) ([]*models.ProjectFetchDTO, error) {
	var projects []*models.ProjectFetchDTO
	for _, project := range m.projects {
		projects = append(projects, project.ToFetchDTO())
	}
	return projects, nil
}

func (m *MockProjectRepository) ListProjectsByTeamID(ctx context.Context, teamID string) ([]*models.ProjectFetchDTO, error) {
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
	ctx := auth.ContextWithClaims(context.Background(), &adminClaims)
	err := service.CreateProject(ctx, projectDTO)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test as user (should force team ID)
	userTeamID := "user-team-id"
	userClaims := auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user"}, Role: "user", TeamID: userTeamID}
	projectDTOUser := &models.ProjectUpsertDTO{
		Name:        "User Project",
		Description: "A user project",
		TeamID:      "some-other-team", // Should be ignored/overwritten
	}

	ctx = auth.ContextWithClaims(context.Background(), &userClaims)
	err = service.CreateProject(ctx, projectDTOUser)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify the user project was created with the correct team ID
	// Since we don't have the ID returned, we can list projects by team ID to verify
	projects, err := mockRepo.ListProjectsByTeamID(context.Background(), userTeamID)
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
	ctx := auth.ContextWithClaims(context.Background(), &claims)
	fetchedProject, err := service.GetProjectByID(ctx, projectID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if fetchedProject.ID != projectID {
		t.Errorf("expected project ID %s, got %s", projectID, fetchedProject.ID)
	}

	// Test unauthorized access (different team)
	otherTeamID := "team-2"
	otherClaims := auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "other"}, Role: "user", TeamID: otherTeamID}
	ctx = auth.ContextWithClaims(context.Background(), &otherClaims)
	_, err = service.GetProjectByID(ctx, projectID)
	if err == nil {
		t.Fatal("expected unauthorized error, got nil")
	}

	// Test superadmin access
	adminClaims := auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "admin"}, Role: "superadmin"}
	ctx = auth.ContextWithClaims(context.Background(), &adminClaims)
	fetchedProjectAdmin, err := service.GetProjectByID(ctx, projectID)
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
	ctx := auth.ContextWithClaims(context.Background(), &claims)
	updatedProject, err := service.UpdateProject(ctx, projectID, updateDTO)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updatedProject.Name != "Updated Name" {
		t.Errorf("expected updated name, got %s", updatedProject.Name)
	}

	// Test unauthorized update
	otherTeamID := "team-2"
	otherClaims := auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "other"}, Role: "user", TeamID: otherTeamID}
	ctx = auth.ContextWithClaims(context.Background(), &otherClaims)
	_, err = service.UpdateProject(ctx, projectID, updateDTO)
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
	ctx := auth.ContextWithClaims(context.Background(), &otherClaims)
	err := service.DeleteProject(ctx, projectID)
	if err == nil {
		t.Fatal("expected unauthorized error, got nil")
	}

	// Test authorized delete
	claims := auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user"}, Role: "user", TeamID: teamID}
	ctx = auth.ContextWithClaims(context.Background(), &claims)
	err = service.DeleteProject(ctx, projectID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify deletion
	_, err = mockRepo.GetProjectByID(context.Background(), projectID)
	if err == nil {
		t.Fatal("expected error getting deleted project, got nil")
	}
}
