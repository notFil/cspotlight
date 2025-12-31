package services

import (
	"context"
	"net/http"
	"path/filepath"

	"cspotlight/internal/auth"
	apperrors "cspotlight/internal/errors"
	"cspotlight/internal/models"
	"cspotlight/internal/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(ctx context.Context, user *models.UserRegisterDTO) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserFetchDTO, error)
	AuthenticateUser(ctx context.Context, authRequest *models.AuthRequest) (*models.UserFetchDTO, error)
	UpdateUser(ctx context.Context, id uuid.UUID, user *models.UserUpdateDTO) (*models.UserFetchDTO, error)
	SetDefaultProject(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (*models.UserFetchDTO, error)
	ListUsers(ctx context.Context) ([]*models.UserFetchDTO, error)
	ListUsersByTeamID(ctx context.Context, teamID uuid.UUID) ([]*models.UserFetchDTO, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ChangePassword(ctx context.Context, id uuid.UUID, request *models.ChangePasswordRequest) error
	ChangeImage(ctx context.Context, id uuid.UUID, path string) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserFetchDTO, error) {
	userContext := auth.GetUserContext(ctx)

	u, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !(userContext.IsSuperadmin() || (userContext.IsAdmin() && u.TeamID != nil && userContext.HasSameTeam(*u.TeamID)) || userContext.IsSameUser(id)) {
		return nil, apperrors.New(http.StatusUnauthorized, "unauthorized access")
	}

	if u == nil {
		return nil, apperrors.New(http.StatusNotFound, "user not found")
	}

	return u.ToFetchDTO(), nil
}

func (s *userService) RegisterUser(ctx context.Context, user *models.UserRegisterDTO) error {
	if user.Password != user.ConfirmPassword {
		return apperrors.New(http.StatusBadRequest, "passwords do not match")
	}
	u := user.ToUser()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	count, err := s.userRepo.GetUserCount(ctx)
	if err != nil {
		return err
	}
	if count == 0 {
		u.Role = "superadmin"
		u.Disabled = false
	} else {
		u.Disabled = true
	}
	if err := s.userRepo.CreateUser(ctx, u); err != nil {
		return err
	}
	return nil
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, user *models.UserUpdateDTO) (*models.UserFetchDTO, error) {
	u, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, apperrors.New(http.StatusNotFound, "user not found")
	}
	u.Role = user.Role
	u.Disabled = user.Disabled
	if user.TeamID != uuid.Nil {
		u.TeamID = &user.TeamID
	} else {
		u.TeamID = nil
	}
	if err := s.userRepo.UpdateUser(ctx, u); err != nil {
		return nil, apperrors.New(http.StatusInternalServerError, "failed to update user")
	}
	return u.ToFetchDTO(), nil
}

func (s *userService) AuthenticateUser(ctx context.Context, authRequest *models.AuthRequest) (*models.UserFetchDTO, error) {
	u, err := s.userRepo.GetUserByUsername(ctx, authRequest.Username)
	if u.Disabled {
		return nil, apperrors.New(http.StatusUnauthorized, "user is disabled")
	}
	if err != nil {
		return nil, apperrors.New(http.StatusUnauthorized, "invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(authRequest.Password)); err != nil {
		return nil, apperrors.New(http.StatusUnauthorized, "invalid credentials")
	}
	return u.ToFetchDTO(), nil
}

func (s *userService) ListUsers(ctx context.Context) ([]*models.UserFetchDTO, error) {
	userContext := auth.GetUserContext(ctx)
	var users []*models.User
	var err error

	if !userContext.IsSuperadmin() {
		users, err = s.userRepo.ListUsersByTeamID(ctx, userContext.TeamID)
	} else {
		users, err = s.userRepo.ListUsers(ctx)
	}

	if err != nil {
		return nil, apperrors.New(http.StatusInternalServerError, "failed to list users")
	}

	var userDTOs []*models.UserFetchDTO
	for _, u := range users {
		userDTOs = append(userDTOs, u.ToFetchDTO())
	}
	return userDTOs, nil
}

func (s *userService) ListUsersByTeamID(ctx context.Context, teamID uuid.UUID) ([]*models.UserFetchDTO, error) {
	users, err := s.userRepo.ListUsersByTeamID(ctx, teamID)
	if err != nil {
		return nil, apperrors.New(http.StatusInternalServerError, "failed to list users")
	}

	var userDTOs []*models.UserFetchDTO
	for _, u := range users {
		userDTOs = append(userDTOs, u.ToFetchDTO())
	}
	return userDTOs, nil
}

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.DeleteUser(ctx, id)
}

func (s *userService) SetDefaultProject(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (*models.UserFetchDTO, error) {
	userContext := auth.GetUserContext(ctx)
	if !userContext.IsSameUser(id) {
		return nil, apperrors.New(http.StatusUnauthorized, "unauthorized access")
	}
	u, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, apperrors.New(http.StatusNotFound, "user not found")
	}
	u.DefaultProjectID = &projectID
	if err := s.userRepo.UpdateUser(ctx, u); err != nil {
		return nil, apperrors.New(http.StatusInternalServerError, "failed to update user")
	}
	return u.ToFetchDTO(), nil
}

func (s *userService) ChangePassword(ctx context.Context, id uuid.UUID, request *models.ChangePasswordRequest) error {
	userContext := auth.GetUserContext(ctx)
	if !userContext.IsSameUser(id) {
		return apperrors.New(http.StatusUnauthorized, "unauthorized access")
	}

	if request.NewPassword != request.ConfirmNewPassword {
		return apperrors.New(http.StatusBadRequest, "new password and confirm new password do not match")
	}

	u, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return apperrors.New(http.StatusNotFound, "user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(request.CurrentPassword)); err != nil {
		return apperrors.New(http.StatusUnauthorized, "invalid password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.New(http.StatusInternalServerError, "failed to update user")
	}
	u.PasswordHash = string(hashedPassword)
	if err := s.userRepo.UpdateUser(ctx, u); err != nil {
		return apperrors.New(http.StatusInternalServerError, "failed to update user")
	}
	return nil
}

func (s *userService) ChangeImage(ctx context.Context, id uuid.UUID, path string) error {
	userContext := auth.GetUserContext(ctx)
	if !userContext.IsSameUser(id) {
		return apperrors.New(http.StatusUnauthorized, "unauthorized access")
	}

	u, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return apperrors.New(http.StatusNotFound, "user not found")
	}

	path = filepath.Clean(path)

	u.Image = path
	if err := s.userRepo.UpdateUser(ctx, u); err != nil {
		return apperrors.New(http.StatusInternalServerError, "failed to update user")
	}
	return nil
}
