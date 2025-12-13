package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	FirstName        string     `gorm:"type:varchar(50)"`
	LastName         string     `gorm:"type:varchar(50)"`
	Username         string     `gorm:"type:varchar(50)"`
	Email            string     `gorm:"type:varchar(100)"`
	PasswordHash     string     `gorm:"type:varchar(255)"`
	Role             string     `gorm:"type:varchar(10)"`
	Image            string     `gorm:"type:varchar(255)"`
	DefaultProjectID *uuid.UUID `gorm:"type:uuid"`
	DefaultProject   Project    `gorm:"foreignKey:DefaultProjectID"`
	Disabled         bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	TeamID           *uuid.UUID     `gorm:"type:uuid"`
	Team             Team           `gorm:"foreignKey:TeamID"`
}

type UserFetchDTO struct {
	ID               uuid.UUID `json:"id"`
	FirstName        string    `json:"firstName"`
	LastName         string    `json:"lastName"`
	Username         string    `json:"username"`
	Role             string    `json:"role"`
	DefaultProjectID uuid.UUID `json:"defaultProjectId,omitempty"`
	Image            string    `json:"image,omitempty"`
	Email            string    `json:"email"`
	Disabled         bool      `json:"disabled"`
	TeamID           uuid.UUID `json:"teamId"`
	TeamName         string    `json:"teamName"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type UserRegisterDTO struct {
	FirstName       string `json:"firstName" binding:"required"`
	LastName        string `json:"lastName" binding:"required"`
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

type UserUpdateDTO struct {
	FirstName string    `json:"firstName" binding:"required"`
	LastName  string    `json:"lastName" binding:"required"`
	Username  string    `json:"username" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	Role      string    `json:"role" binding:"required"`
	TeamID    uuid.UUID `json:"teamId"`
	Disabled  bool      `json:"disabled"`
}

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword    string `json:"currentPassword" binding:"required"`
	NewPassword        string `json:"newPassword" binding:"required"`
	ConfirmNewPassword string `json:"confirmNewPassword" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type SetDefaultProjectRequest struct {
	ProjectID string `json:"projectID" binding:"required"`
}

func (u *User) ToFetchDTO() *UserFetchDTO {
	dto := UserFetchDTO{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Email:     u.Email,
		Image:     u.Image,
		Disabled:  u.Disabled,
		Role:      u.Role,
		UpdatedAt: u.UpdatedAt,
	}
	if u.TeamID != nil {
		dto.TeamID = *u.TeamID
		dto.TeamName = u.Team.Name
	}
	if u.DefaultProjectID != nil {
		dto.DefaultProjectID = *u.DefaultProjectID
	}
	return &dto
}

func (u *UserUpdateDTO) ToUser() *User {
	var teamID *uuid.UUID
	if u.TeamID != uuid.Nil {
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

func (u *UserRegisterDTO) ToUser() *User {
	return &User{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Email:     u.Email,
		Role:      "user",
	}
}
