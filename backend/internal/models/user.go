package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           string `gorm:"primaryKey;default:uuid_generate_v4()"`
	FirstName    string `gorm:"type:varchar(50)"`
	LastName     string `gorm:"type:varchar(50)"`
	Username     string `gorm:"type:varchar(50)"`
	Email        string `gorm:"type:varchar(100)"`
	PasswordHash string `gorm:"type:varchar(50)"`
	Role         string `gorm:"type:varchar(10)"`
	Disabled     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	TeamID       *string        `gorm:"type:uuid"`
	Team         Team           `gorm:"foreignKey:TeamID"`
}

type UserFetchDTO struct {
	ID        string    `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	Email     string    `json:"email"`
	Disabled  bool      `json:"disabled"`
	TeamID    string    `json:"teamId"`
	TeamName  string    `json:"teamName"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserCreateDTO struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	Role      string `json:"role" binding:"required"`
	TeamID    string `json:"teamId"`
}

type UserUpdateDTO struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Role      string `json:"role" binding:"required"`
	TeamID    string `json:"teamId"`
}

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

func (u *User) ToFetchDTO() *UserFetchDTO {
	dto := UserFetchDTO{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Email:     u.Email,
		Disabled:  u.Disabled,
		Role:      u.Role,
	}
	if u.TeamID != nil {
		dto.TeamID = *u.TeamID
		dto.TeamName = u.Team.Name
	}
	return &dto
}

func (u *UserCreateDTO) ToUser() *User {
	var teamID *string
	if u.TeamID != "" {
		teamID = &u.TeamID
	}
	return &User{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		TeamID:    teamID,
	}
}

func (u *UserUpdateDTO) ToUser() *User {
	var teamID *string
	if u.TeamID != "" {
		teamID = &u.TeamID
	}
	return &User{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		TeamID:    teamID,
	}
}
