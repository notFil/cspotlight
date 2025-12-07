package repositories

import (
	"context"

	"github.com/notFil/cspotlight/internal/models"
	"gorm.io/gorm"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *models.Team) error
	GetTeamByID(ctx context.Context, id string) (*models.Team, error)
	UpdateTeam(ctx context.Context, team *models.Team) error
	DeleteTeam(ctx context.Context, id string) error
	ListTeams(ctx context.Context) ([]*models.Team, error)
}

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) CreateTeam(ctx context.Context, team *models.Team) error {
	return r.db.WithContext(ctx).Create(team).Error
}

func (r *teamRepository) GetTeamByID(ctx context.Context, id string) (team *models.Team, err error) {
	err = r.db.WithContext(ctx).First(&team, "id = ?", id).Error
	return team, err
}

func (r *teamRepository) UpdateTeam(ctx context.Context, team *models.Team) error {
	return r.db.WithContext(ctx).Save(team).Error
}

func (r *teamRepository) DeleteTeam(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Team{}, "id = ?", id).Error
}

func (r *teamRepository) ListTeams(ctx context.Context) (teams []*models.Team, err error) {
	if err = r.db.WithContext(ctx).Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}
