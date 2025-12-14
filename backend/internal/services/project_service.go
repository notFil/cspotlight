package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/notFil/cspotlight/internal/auth"

	apperrors "github.com/notFil/cspotlight/internal/errors"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/repositories"
)

type ProjectService interface {
	GetProjectByID(ctx context.Context, id uuid.UUID) (*models.ProjectFetchDTO, error)
	CreateProject(ctx context.Context, project *models.ProjectUpsertDTO) error
	UpdateProject(ctx context.Context, id uuid.UUID, project *models.ProjectUpsertDTO) (*models.ProjectFetchDTO, error)
	DeleteProject(ctx context.Context, id uuid.UUID) error
	ListProjects(ctx context.Context) ([]*models.ProjectFetchDTO, error)
}

type projectService struct {
	projectRepo repositories.ProjectRepository
	baseURL     string
}

func NewProjectService(projectRepo repositories.ProjectRepository, baseURL string) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
		baseURL:     baseURL,
	}
}

func (s *projectService) GetProjectByID(ctx context.Context, id uuid.UUID) (*models.ProjectFetchDTO, error) {
	p, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	userContext := auth.GetUserContext(ctx)

	if !userContext.IsSuperadmin() && !userContext.HasSameTeam(p.TeamID) {
		return nil, apperrors.New(http.StatusUnauthorized, "unauthorized access")
	}

	if p == nil {
		return nil, apperrors.New(http.StatusNotFound, "project not found")
	}

	project := p.ToFetchDTO()
	project.ReportingURL = fmt.Sprintf("%s/api/v1/reports/%s/endpoint", s.baseURL, p.ID.String())

	return project, nil
}

func (s *projectService) CreateProject(ctx context.Context, project *models.ProjectUpsertDTO) error {
	userContext := auth.GetUserContext(ctx)

	p := project.ToProject()

	if !userContext.IsSuperadmin() && !userContext.HasSameTeam(p.TeamID) {
		return apperrors.New(http.StatusUnauthorized, "unauthorized access")
	}

	if err := s.projectRepo.CreateProject(ctx, p); err != nil {
		return apperrors.New(http.StatusInternalServerError, "failed to create project")
	}
	return nil
}

func (s *projectService) UpdateProject(ctx context.Context, id uuid.UUID, project *models.ProjectUpsertDTO) (*models.ProjectFetchDTO, error) {
	userContext := auth.GetUserContext(ctx)

	p, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return nil, apperrors.New(http.StatusNotFound, "failed to update project")
	}

	if !userContext.IsSuperadmin() && !userContext.HasSameTeam(p.TeamID) {
		return nil, apperrors.New(http.StatusUnauthorized, "unauthorized access")
	}

	if p == nil {
		return nil, apperrors.New(http.StatusNotFound, "project not found")
	}
	p.Name = project.Name
	p.Description = project.Description
	if err := s.projectRepo.UpdateProject(ctx, p); err != nil {
		return nil, apperrors.New(http.StatusInternalServerError, "failed to update project")
	}
	return p.ToFetchDTO(), nil
}

func (s *projectService) DeleteProject(ctx context.Context, id uuid.UUID) error {
	userContext := auth.GetUserContext(ctx)

	p, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return apperrors.New(http.StatusNotFound, "project not found")
	}

	if !userContext.IsSuperadmin() && !userContext.HasSameTeam(p.TeamID) {
		return apperrors.New(http.StatusUnauthorized, "unauthorized access")
	}

	if err := s.projectRepo.DeleteProject(ctx, id); err != nil {
		return apperrors.New(http.StatusInternalServerError, "failed to delete project")
	}
	return nil
}

func (s *projectService) ListProjects(ctx context.Context) ([]*models.ProjectFetchDTO, error) {
	var projects []*models.ProjectFetchDTO
	var err error

	userContext := auth.GetUserContext(ctx)

	if userContext.IsSuperadmin() {
		projects, err = s.projectRepo.ListProjects(ctx)
	} else {
		projects, err = s.projectRepo.ListProjectsByTeamID(ctx, userContext.TeamID)
	}
	if err != nil {
		return nil, apperrors.New(http.StatusInternalServerError, "failed to list projects")
	}
	for _, p := range projects {
		p.ReportingURL = fmt.Sprintf("%s/api/v1/reports/%s/endpoint", s.baseURL, p.ID.String())
	}

	if projects == nil {
		return []*models.ProjectFetchDTO{}, nil
	}
	return projects, nil
}
