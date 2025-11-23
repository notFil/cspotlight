package services

import (
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/repositories"
)

type TeamService interface {
	CreateTeam(user *models.TeamUpsertDTO) (bool, error)
	GetTeamByID(id string) (*models.TeamFetchDTO, error)
	UpdateTeam(id string, team *models.TeamUpsertDTO) (*models.TeamFetchDTO, error)
	DeleteTeam(id string) (bool, error)
	ListTeams() ([]*models.TeamFetchDTO, error)
}

type teamService struct {
	teamRepo repositories.TeamRepository
}

func NewTeamService(teamRepo repositories.TeamRepository) TeamService {
	return &teamService{
		teamRepo: teamRepo,
	}
}

func (s *teamService) CreateTeam(team *models.TeamUpsertDTO) (bool, error) {
	t := team.ToTeam()
	if err := s.teamRepo.CreateTeam(t); err != nil {
		return false, err
	}
	return true, nil
}

func (s *teamService) GetTeamByID(id string) (*models.TeamFetchDTO, error) {
	t, err := s.teamRepo.GetTeamByID(id)
	if err != nil {
		return nil, err
	}
	return t.ToFetchDTO(), nil
}

func (s *teamService) UpdateTeam(id string, team *models.TeamUpsertDTO) (*models.TeamFetchDTO, error) {
	t, err := s.teamRepo.GetTeamByID(id)
	if err != nil {
		return nil, err
	}
	t.Name = team.Name
	t.Description = team.Description
	if err := s.teamRepo.UpdateTeam(t); err != nil {
		return nil, err
	}
	return t.ToFetchDTO(), nil
}

func (s *teamService) DeleteTeam(id string) (bool, error) {
	if err := s.teamRepo.DeleteTeam(id); err != nil {
		return false, err
	}
	return true, nil
}

func (s *teamService) ListTeams() (teams []*models.TeamFetchDTO, err error) {
	ts, err := s.teamRepo.ListTeams()
	if err != nil {
		return nil, err
	}
	var teamDTOs []*models.TeamFetchDTO
	for _, t := range ts {
		teamDTOs = append(teamDTOs, t.ToFetchDTO())
	}
	return teamDTOs, nil
}
