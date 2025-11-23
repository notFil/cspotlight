package repositories

import (
	"github.com/notFil/cspotlight/internal/models"
	"gorm.io/gorm"
)

type TeamRepository interface {
	CreateTeam(team *models.Team) error
	GetTeamByID(id string) (*models.Team, error)
	UpdateTeam(team *models.Team) error
	DeleteTeam(id string) error
	ListTeams() ([]*models.Team, error)
}

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) CreateTeam(team *models.Team) error {
	return r.db.Create(team).Error
}

func (r *teamRepository) GetTeamByID(id string) (team *models.Team, err error) {
	err = r.db.First(&team, "id = ?", id).Error
	return team, err
}

func (r *teamRepository) UpdateTeam(team *models.Team) error {
	return r.db.Save(team).Error
}

func (r *teamRepository) DeleteTeam(id string) error {
	return r.db.Delete(&models.Team{}, "id = ?", id).Error
}

func (r *teamRepository) ListTeams() (teams []*models.Team, err error) {
	if err = r.db.Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}
