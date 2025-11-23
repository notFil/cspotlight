package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID          string `gorm:"primaryKey;default:uuid_generate_v4()"`
	Name        string `gorm:"type:varchar(20)"`
	Description string `gorm:"type:varchar(500)"`
	TeamID      string `gorm:"type:uuid"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type ProjectFetchDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TeamID      string `json:"teamId"`
}

type ProjectUpsertDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	TeamID      string `json:"teamId" binding:"required"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

func (p *Project) ToFetchDTO() *ProjectFetchDTO {
	return &ProjectFetchDTO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		TeamID:      p.TeamID,
	}
}

func (p *ProjectUpsertDTO) ToProject() *Project {
	return &Project{
		Name:        p.Name,
		Description: p.Description,
		TeamID:      p.TeamID,
	}
}
