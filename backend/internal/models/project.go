package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID          string `gorm:"primaryKey;default:uuid_generate_v4()"`
	Name        string `gorm:"type:varchar(20)"`
	Description string `gorm:"type:varchar(500)"`
	TeamID      string `gorm:"type:uuid"`
	Team        Team   `gorm:"foreignKey:TeamID"`
	Disabled    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type ProjectFetchDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Disabled    bool
	TeamID      string `json:"teamId"`
	TeamName    string `json:"teamName"`
	LastActive  string `json:"lastActive"`
}
type ProjectUpsertDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	TeamID      string `json:"teamId" binding:"required"`
	Disabled    bool   `json:"disabled"`
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

func (p *ProjectUpsertDTO) ToProject() *Project {
	return &Project{
		Name:        p.Name,
		Description: p.Description,
		TeamID:      p.TeamID,
	}
}
