package services

import (
	"context"
	"net/http"

	"github.com/notFil/cspotlight/internal/auth"

	apperrors "github.com/notFil/cspotlight/internal/errors"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/repositories"
)

type ProjectService interface {
	GetProjectByID(ctx context.Context, id string) (*models.ProjectFetchDTO, error)
	CreateProject(ctx context.Context, project *models.ProjectUpsertDTO) error
	UpdateProject(ctx context.Context, id string, project *models.ProjectUpsertDTO) (*models.ProjectFetchDTO, error)
	DeleteProject(ctx context.Context, id string) error
	ListProjects(ctx context.Context) ([]*models.ProjectFetchDTO, error)
}

type projectService struct {
	projectRepo repositories.ProjectRepository
}

func NewProjectService(projectRepo repositories.ProjectRepository) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
	}
}

func (s *projectService) GetProjectByID(ctx context.Context, id string) (*models.ProjectFetchDTO, error) {
	p, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	claims := auth.GetUserClaims(ctx)

	if !claims.IsSuperadmin() && p.TeamID != claims.TeamID {
		return nil, apperrors.New(http.StatusUnauthorized, "unauthorized access")
	}

	if p == nil {
		return nil, apperrors.New(http.StatusNotFound, "project not found")
	}

	return p.ToFetchDTO(), nil
}

func (s *projectService) CreateProject(ctx context.Context, project *models.ProjectUpsertDTO) error {
	claims := auth.GetUserClaims(ctx)
	p := project.ToProject()
	if !claims.IsSuperadmin() {
		p.TeamID = claims.TeamID
	}
	if err := s.projectRepo.CreateProject(ctx, p); err != nil {
		return apperrors.New(http.StatusInternalServerError, "failed to create project")
	}
	return nil
}

func (s *projectService) UpdateProject(ctx context.Context, id string, project *models.ProjectUpsertDTO) (*models.ProjectFetchDTO, error) {
	claims := auth.GetUserClaims(ctx)
	p, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return nil, apperrors.New(http.StatusNotFound, "failed to update project")
	}

	if !claims.IsSuperadmin() && p.TeamID != claims.TeamID {
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

func (s *projectService) DeleteProject(ctx context.Context, id string) error {
	claims := auth.GetUserClaims(ctx)
	p, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return apperrors.New(http.StatusNotFound, "project not found")
	}
	if !claims.IsSuperadmin() && p.TeamID != claims.TeamID {
		return apperrors.New(http.StatusUnauthorized, "unauthorized access")
	}
	if err := s.projectRepo.DeleteProject(ctx, id); err != nil {
		return apperrors.New(http.StatusInternalServerError, "failed to delete project")
	}
	return nil
}

func (s *projectService) ListProjects(ctx context.Context) ([]*models.ProjectFetchDTO, error) {
	var ps []*models.ProjectFetchDTO
	var err error

	claims := auth.GetUserClaims(ctx)

	if claims.IsSuperadmin() {
		ps, err = s.projectRepo.ListProjects(ctx)
	} else {
		ps, err = s.projectRepo.ListProjectsByTeamID(ctx, claims.TeamID)
	}
	if err != nil {
		return nil, apperrors.New(http.StatusInternalServerError, "failed to list projects")
	}
	if ps == nil {
		ps = []*models.ProjectFetchDTO{}
	}
	return ps, nil
}
