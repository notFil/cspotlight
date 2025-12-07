package repositories

import (
	"context"

	"github.com/notFil/cspotlight/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context) ([]*models.User, error)
	ListUsersByTeamID(ctx context.Context, teamID string) ([]*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (user *models.User, err error) {
	if err = r.db.WithContext(ctx).Preload("Team").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (user *models.User, err error) {
	if err = r.db.WithContext(ctx).Preload("Team").Where("email = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}

func (r *userRepository) ListUsers(ctx context.Context) (users []*models.User, err error) {
	if err = r.db.WithContext(ctx).Preload("Team").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) ListUsersByTeamID(ctx context.Context, teamID string) (users []*models.User, err error) {
	if err = r.db.WithContext(ctx).Preload("Team").Where("team_id = ?", teamID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
