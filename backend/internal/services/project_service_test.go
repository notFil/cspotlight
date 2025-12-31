package services

import (
	"context"
	"errors"
	"testing"

	"cspotlight/internal/auth"
	"cspotlight/internal/models"

	"github.com/google/uuid"
)

type MockProjectRepository struct {
	projects map[uuid.UUID]*models.Project
}

func NewMockProjectRepository() *MockProjectRepository {
	return &MockProjectRepository{
		projects: make(map[uuid.UUID]*models.Project),
	}
}

func (m *MockProjectRepository) CreateProject(ctx context.Context, project *models.Project) error {
	if project.ID == uuid.Nil {
		project.ID = uuid.New()
	}
	if _, exists := m.projects[project.ID]; exists {
		return errors.New("project already exists")
	}
	m.projects[project.ID] = project
	return nil
}

func (m *MockProjectRepository) GetProjectByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
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

func (m *MockProjectRepository) DeleteProject(ctx context.Context, id uuid.UUID) error {
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

func (m *MockProjectRepository) ListProjectsByTeamID(ctx context.Context, teamID uuid.UUID) ([]*models.ProjectFetchDTO, error) {
	var projects []*models.ProjectFetchDTO
	for _, project := range m.projects {
		if project.TeamID == teamID {
			projects = append(projects, project.ToFetchDTO())
		}
	}
	return projects, nil
}

func (m *MockProjectRepository) CreateProjectWithID(ctx context.Context, project *models.Project) error {
	if project.ID == uuid.Nil {
		project.ID = uuid.New()
	}
	m.projects[project.ID] = project
	return nil
}

func TestCreateProject(t *testing.T) {
	mockRepo := NewMockProjectRepository()
	service := NewProjectService(mockRepo, "http://localhost:8080")

	projectDTO := &models.ProjectUpsertDTO{
		Name:        "Test Project",
		Description: "A test project",
		TeamID:      uuid.New(),
	}

	adminData := auth.UserContext{ID: uuid.New(), Role: "superadmin"}
	ctx := auth.ContextWithUser(context.Background(), &adminData)
	err := service.CreateProject(ctx, projectDTO)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	userTeamID := uuid.New()
	userData := auth.UserContext{ID: uuid.New(), Role: "user", TeamID: userTeamID}
	projectDTOUser := &models.ProjectUpsertDTO{
		Name:        "User Project",
		Description: "A user project",
		TeamID:      userTeamID,
	}

	ctx = auth.ContextWithUser(context.Background(), &userData)
	err = service.CreateProject(ctx, projectDTOUser)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

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
	service := NewProjectService(mockRepo, "http://localhost:8080")

	projectID := uuid.New()
	teamID := uuid.New()
	project := &models.Project{
		ID:     projectID,
		Name:   "Test Project",
		TeamID: teamID,
	}
	mockRepo.projects[projectID] = project

	userData := auth.UserContext{ID: uuid.New(), Role: "user", TeamID: teamID}
	ctx := auth.ContextWithUser(context.Background(), &userData)
	fetchedProject, err := service.GetProjectByID(ctx, projectID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if fetchedProject.ID != projectID {
		t.Errorf("expected project ID %s, got %s", projectID, fetchedProject.ID)
	}

	otherTeamID := uuid.New()
	otherData := auth.UserContext{ID: uuid.New(), Role: "user", TeamID: otherTeamID}
	ctx = auth.ContextWithUser(context.Background(), &otherData)
	_, err = service.GetProjectByID(ctx, projectID)
	if err == nil {
		t.Fatal("expected unauthorized error, got nil")
	}

	// Test superadmin access
	adminData := auth.UserContext{ID: uuid.New(), Role: "superadmin"}
	ctx = auth.ContextWithUser(context.Background(), &adminData)
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
	service := NewProjectService(mockRepo, "http://localhost:8080")

	projectID := uuid.New()
	teamID := uuid.New()
	project := &models.Project{
		ID:     projectID,
		Name:   "Original Name",
		TeamID: teamID,
	}
	mockRepo.projects[projectID] = project

	updateDTO := &models.ProjectUpsertDTO{
		Name: "Updated Name",
	}

	userData := auth.UserContext{ID: uuid.New(), Role: "user", TeamID: teamID}
	ctx := auth.ContextWithUser(context.Background(), &userData)
	updatedProject, err := service.UpdateProject(ctx, projectID, updateDTO)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updatedProject.Name != "Updated Name" {
		t.Errorf("expected updated name, got %s", updatedProject.Name)
	}

	otherTeamID := uuid.New()
	otherData := auth.UserContext{ID: uuid.New(), Role: "user", TeamID: otherTeamID}
	ctx = auth.ContextWithUser(context.Background(), &otherData)
	_, err = service.UpdateProject(ctx, projectID, updateDTO)
	if err == nil {
		t.Fatal("expected unauthorized error, got nil")
	}
}

func TestDeleteProject(t *testing.T) {
	mockRepo := NewMockProjectRepository()
	service := NewProjectService(mockRepo, "http://localhost:8080")

	projectID := uuid.New()
	teamID := uuid.New()
	project := &models.Project{
		ID:     projectID,
		Name:   "To Delete",
		TeamID: teamID,
	}
	mockRepo.projects[projectID] = project

	otherTeamID := uuid.New()
	otherData := auth.UserContext{ID: uuid.New(), Role: "user", TeamID: otherTeamID}
	ctx := auth.ContextWithUser(context.Background(), &otherData)
	err := service.DeleteProject(ctx, projectID)
	if err == nil {
		t.Fatal("expected unauthorized error, got nil")
	}

	userData := auth.UserContext{ID: uuid.New(), Role: "user", TeamID: teamID}
	ctx = auth.ContextWithUser(context.Background(), &userData)
	err = service.DeleteProject(ctx, projectID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = mockRepo.GetProjectByID(context.Background(), projectID)
	if err == nil {
		t.Fatal("expected error getting deleted project, got nil")
	}

	err = service.DeleteProject(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent project, got nil")
	}
}
