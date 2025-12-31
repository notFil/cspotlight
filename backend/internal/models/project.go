package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string    `gorm:"type:varchar(20)"`
	Description string    `gorm:"type:varchar(500)"`
	TeamID      uuid.UUID `gorm:"type:uuid"`
	Team        Team      `gorm:"foreignKey:TeamID"`
	Disabled    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type ProjectFetchDTO struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Disabled     bool       `json:"disabled"`
	TeamID       uuid.UUID  `json:"teamId"`
	TeamName     string     `json:"teamName"`
	ReportingURL string     `json:"reportingUrl"`
	LastActive   *time.Time `json:"lastActive"`
}
type ProjectUpsertDTO struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	TeamID      uuid.UUID `json:"teamId" binding:"required"`
	Disabled    bool      `json:"disabled"`
}

func (p *Project) ToFetchDTO() *ProjectFetchDTO {
	return &ProjectFetchDTO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		TeamID:      p.TeamID,
		TeamName:    p.Team.Name,
		Disabled:    p.Disabled,
	}
}

func (p *Project) BelongsToTeam(teamID uuid.UUID) bool {
	return p.TeamID == teamID
}

func (p *ProjectUpsertDTO) ToProject() *Project {
	return &Project{
		Name:        p.Name,
		Description: p.Description,
		TeamID:      p.TeamID,
	}
}

func (p *Project) AfterDelete(tx *gorm.DB) (err error) {
	tx.Delete(&CSPReport{}, "project_id = ?", p.ID)
	return
}
