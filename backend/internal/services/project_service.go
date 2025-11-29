package services

import (
	"errors"

	"github.com/notFil/cspotlight/pkg/auth"

	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/repositories"
)

type ProjectService interface {
	GetProjectByID(id string, claims auth.Claims) (*models.ProjectFetchDTO, error)
	CreateProject(project *models.ProjectUpsertDTO, claims auth.Claims) (bool, error)
	UpdateProject(id string, project *models.ProjectUpsertDTO, claims auth.Claims) (*models.ProjectFetchDTO, error)
	DeleteProject(id string, claims auth.Claims) error
	ListProjects(claims auth.Claims) ([]*models.ProjectFetchDTO, error)
}

type projectService struct {
	projectRepo repositories.ProjectRepository
}

func NewProjectService(projectRepo repositories.ProjectRepository) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
	}
}

func (s *projectService) GetProjectByID(id string, claims auth.Claims) (*models.ProjectFetchDTO, error) {
	p, err := s.projectRepo.GetProjectByID(id)
	if err != nil {
		return nil, err
	}
	if !claims.IsSuperadmin() && p.TeamID != claims.TeamID {
		return nil, errors.New("unauthorized access")
	}

	return p.ToFetchDTO(), nil
}

func (s *projectService) CreateProject(project *models.ProjectUpsertDTO, claims auth.Claims) (bool, error) {
	p := project.ToProject()
	if !claims.IsSuperadmin() {
		p.TeamID = claims.TeamID
	}
	if err := s.projectRepo.CreateProject(p); err != nil {
		return false, err
	}
	return true, nil
}

func (s *projectService) UpdateProject(id string, project *models.ProjectUpsertDTO, claims auth.Claims) (*models.ProjectFetchDTO, error) {
	p, err := s.projectRepo.GetProjectByID(id)
	if err != nil {
		return nil, err
	}
	if !claims.IsSuperadmin() && p.TeamID != claims.TeamID {
		return nil, errors.New("unauthorized access")
	}
	p.Name = project.Name
	p.Description = project.Description
	if err := s.projectRepo.UpdateProject(p); err != nil {
		return nil, err
	}
	return p.ToFetchDTO(), nil
}

func (s *projectService) DeleteProject(id string, claims auth.Claims) error {
	p, err := s.projectRepo.GetProjectByID(id)
	if err != nil {
		return err
	}
	if !claims.IsSuperadmin() && p.TeamID != claims.TeamID {
		return errors.New("unauthorized access")
	}
	return s.projectRepo.DeleteProject(id)
}

func (s *projectService) ListProjects(claims auth.Claims) ([]*models.ProjectFetchDTO, error) {
	var ps []*models.Project
	var err error
	if claims.IsSuperadmin() {
		ps, err = s.projectRepo.ListProjects()
	} else {
		ps, err = s.projectRepo.ListProjectsByTeamID(claims.TeamID)
	}
	if err != nil {
		return nil, err
	}

	var projects []*models.ProjectFetchDTO
	for _, p := range ps {
		projects = append(projects, p.ToFetchDTO())
	}
	return projects, nil
}
