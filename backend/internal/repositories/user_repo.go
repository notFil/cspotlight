package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/notFil/cspotlight/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserCount(ctx context.Context) (int64, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsers(ctx context.Context) ([]*models.User, error)
	ListUsersByTeamID(ctx context.Context, teamID uuid.UUID) ([]*models.User, error)
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

func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (user *models.User, err error) {
	if err = r.db.WithContext(ctx).Preload("Team").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (user *models.User, err error) {
	if err = r.db.WithContext(ctx).Preload("Team").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserCount(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}

func (r *userRepository) ListUsers(ctx context.Context) (users []*models.User, err error) {
	if err = r.db.WithContext(ctx).Preload("Team").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) ListUsersByTeamID(ctx context.Context, teamID uuid.UUID) (users []*models.User, err error) {
	if err = r.db.WithContext(ctx).Preload("Team").Where("team_id = ?", teamID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
