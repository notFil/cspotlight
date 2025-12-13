package services

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	apperrors "github.com/notFil/cspotlight/internal/errors"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/repositories"
)

type TeamService interface {
	CreateTeam(ctx context.Context, user *models.TeamUpsertDTO) error
	GetTeamByID(ctx context.Context, id uuid.UUID) (*models.TeamFetchDTO, error)
	UpdateTeam(ctx context.Context, id uuid.UUID, team *models.TeamUpsertDTO) (*models.TeamFetchDTO, error)
	DeleteTeam(ctx context.Context, id uuid.UUID) error
	ListTeams(ctx context.Context) ([]*models.TeamFetchDTO, error)
}

type teamService struct {
	teamRepo repositories.TeamRepository
}

func NewTeamService(teamRepo repositories.TeamRepository) TeamService {
	return &teamService{
		teamRepo: teamRepo,
	}
}

func (s *teamService) CreateTeam(ctx context.Context, team *models.TeamUpsertDTO) error {
	t := team.ToTeam()
	if err := s.teamRepo.CreateTeam(ctx, t); err != nil {
		return apperrors.New(http.StatusInternalServerError, "failed to create team")
	}
	return nil
}

func (s *teamService) GetTeamByID(ctx context.Context, id uuid.UUID) (*models.TeamFetchDTO, error) {
	t, err := s.teamRepo.GetTeamByID(ctx, id)
	if err != nil {
		return nil, apperrors.New(http.StatusNotFound, "team not found")
	}
	if t == nil {
		return nil, apperrors.New(http.StatusNotFound, "team not found")
	}
	return t.ToFetchDTO(), nil
}

func (s *teamService) UpdateTeam(ctx context.Context, id uuid.UUID, team *models.TeamUpsertDTO) (*models.TeamFetchDTO, error) {
	t, err := s.teamRepo.GetTeamByID(ctx, id)
	if err != nil {
		return nil, apperrors.New(http.StatusNotFound, "team not found")
	}
	t.Name = team.Name
	t.Description = team.Description
	if err := s.teamRepo.UpdateTeam(ctx, t); err != nil {
		return nil, apperrors.New(http.StatusInternalServerError, "failed to update team")
	}
	return t.ToFetchDTO(), nil
}

func (s *teamService) DeleteTeam(ctx context.Context, id uuid.UUID) error {
	if err := s.teamRepo.DeleteTeam(ctx, id); err != nil {
		return apperrors.New(http.StatusInternalServerError, "failed to delete team")
	}
	return nil
}

func (s *teamService) ListTeams(ctx context.Context) (teams []*models.TeamFetchDTO, err error) {
	ts, err := s.teamRepo.ListTeams(ctx)
	if err != nil {
		return nil, apperrors.New(http.StatusInternalServerError, "failed to list teams")
	}
	var teamDTOs []*models.TeamFetchDTO
	for _, t := range ts {
		teamDTOs = append(teamDTOs, t.ToFetchDTO())
	}
	return teamDTOs, nil
}
