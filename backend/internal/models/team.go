package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string    `gorm:"type:varchar(20)"`
	Description string    `gorm:"type:varchar(500)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type TeamFetchDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UpdatedAt   string    `json:"updatedAt"`
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
		UpdatedAt:   t.UpdatedAt.Format("02 Jan 06 15:04 MST"),
	}
}
