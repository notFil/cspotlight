package models

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	ID          string `gorm:"primaryKey;default:uuid_generate_v4()"`
	Name        string `gorm:"type:varchar(20)"`
	Description string `gorm:"type:varchar(500)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type TeamFetchDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TeamUpsertDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (t *TeamUpsertDTO) ToTeam() *Team {
	return &Team{
		Name:        t.Name,
		Description: t.Description,
	}
}

func (t *Team) ToFetchDTO() *TeamFetchDTO {
	return &TeamFetchDTO{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
	}
}
