package services

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"slices"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/notFil/cspotlight/internal/auth"
	apperrors "github.com/notFil/cspotlight/internal/errors"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

const ImageMaxSize = 2 << 20

var ImageAllowedExts = []string{"image/jpeg", "image/png", "image/gif"}

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
	ChangeImage(ctx context.Context, id uuid.UUID, image *multipart.FileHeader) error
}

type userService struct {
	userRepo    repositories.UserRepository
	projectRepo repositories.ProjectRepository
}

func NewUserService(userRepo repositories.UserRepository, projectRepo repositories.ProjectRepository) UserService {
	return &userService{
		userRepo:    userRepo,
		projectRepo: projectRepo,
	}
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserFetchDTO, error) {
	userContext := auth.GetUserContext(ctx)

	u, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !(userContext.IsSuperadmin() || (userContext.IsAdmin() && u.TeamID != nil && userContext.SameTeam(*u.TeamID)) || userContext.SameUser(id)) {
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
	if !userContext.SameUser(id) {
		return nil, apperrors.New(http.StatusUnauthorized, "unauthorized access")
	}

	project, err := s.projectRepo.GetProjectByID(ctx, projectID)
	if project == nil || err != nil {
		return nil, apperrors.New(http.StatusNotFound, "project not found")
	}

	if !userContext.IsSuperadmin() && !project.BelongsToTeam(userContext.TeamID) {
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
	if !userContext.SameUser(id) {
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

func (s *userService) ChangeImage(ctx context.Context, id uuid.UUID, image *multipart.FileHeader) error {
	userContext := auth.GetUserContext(ctx)
	if !userContext.SameUser(id) {
		return apperrors.New(http.StatusUnauthorized, "unauthorized access")
	}

	if image.Size > ImageMaxSize {
		return apperrors.New(http.StatusBadRequest, "image too large")
	}

	f, err := image.Open()
	if err != nil {
		return apperrors.New(http.StatusInternalServerError, "failed to open image")
	}
	defer f.Close()

	mimeType, err := validateImageMimeType(f)
	if err != nil {
		return apperrors.New(http.StatusBadRequest, "invalid image format")
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return apperrors.New(http.StatusInternalServerError, "failed to read image")
	}

	strFile := base64.StdEncoding.EncodeToString(data)

	u, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return apperrors.New(http.StatusNotFound, "user not found")
	}
	u.Image = fmt.Sprintf("data:%s;base64,%s", mimeType, strFile)
	if err := s.userRepo.UpdateUser(ctx, u); err != nil {
		return apperrors.New(http.StatusInternalServerError, "failed to update user")
	}
	return nil
}

func validateImageMimeType(f multipart.File) (string, error) {
	mimeType, err := mimetype.DetectReader(f)
	if err != nil {
		return "", err
	}
	if !slices.Contains(ImageAllowedExts, mimeType.String()) {
		return "", errors.New("invalid image format")
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	return mimeType.String(), nil
}
