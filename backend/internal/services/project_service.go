package services

import (
	"errors"

	"github.com/notFil/cspotlight/internal/common"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/repositories"
)

type ProjectService interface {
	GetProjectByID(id string, claims common.AuthClaims) (*models.ProjectFetchDTO, error)
	CreateProject(project *models.ProjectUpsertDTO, claims common.AuthClaims) (bool, error)
	UpdateProject(id string, project *models.ProjectUpsertDTO, claims common.AuthClaims) (*models.ProjectFetchDTO, error)
	DeleteProject(id string, claims common.AuthClaims) error
	ListProjects(claims common.AuthClaims) ([]*models.ProjectFetchDTO, error)
}

type projectService struct {
	projectRepo repositories.ProjectRepository
}

func NewProjectService(projectRepo repositories.ProjectRepository) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
	}
}

func (s *projectService) GetProjectByID(id string, claims common.AuthClaims) (*models.ProjectFetchDTO, error) {
	p, err := s.projectRepo.GetProjectByID(id)
	if err != nil {
		return nil, err
	}
	if !claims.IsSuperadmin() && p.TeamID != *claims.TeamID {
		return nil, errors.New("unauthorized access")
	}

	return p.ToFetchDTO(), nil
}

func (s *projectService) CreateProject(project *models.ProjectUpsertDTO, claims common.AuthClaims) (bool, error) {
	p := project.ToProject()
	if !claims.IsSuperadmin() {
		p.TeamID = *claims.TeamID
	}
	if err := s.projectRepo.CreateProject(p); err != nil {
		return false, err
	}
	return true, nil
}

func (s *projectService) UpdateProject(id string, project *models.ProjectUpsertDTO, claims common.AuthClaims) (*models.ProjectFetchDTO, error) {
	p, err := s.projectRepo.GetProjectByID(id)
	if err != nil {
		return nil, err
	}
	if !claims.IsSuperadmin() && p.TeamID != *claims.TeamID {
		return nil, errors.New("unauthorized access")
	}
	p.Name = project.Name
	p.Description = project.Description
	if err := s.projectRepo.UpdateProject(p); err != nil {
		return nil, err
	}
	return p.ToFetchDTO(), nil
}

func (s *projectService) DeleteProject(id string, claims common.AuthClaims) error {
	p, err := s.projectRepo.GetProjectByID(id)
	if err != nil {
		return err
	}
	if !claims.IsSuperadmin() && p.TeamID != *claims.TeamID {
		return errors.New("unauthorized access")
	}
	return s.projectRepo.DeleteProject(id)
}

func (s *projectService) ListProjects(claims common.AuthClaims) (projects []*models.ProjectFetchDTO, err error) {
	var ps []*models.Project
	if claims.IsSuperadmin() {
		ps, err = s.projectRepo.ListProjects()
	} else {
		ps, err = s.projectRepo.ListProjectsByTeamID(*claims.TeamID)
	}
	if err != nil {
		return nil, err
	}

	for _, p := range ps {
		projects = append(projects, p.ToFetchDTO())
	}
	return projects, nil
}
