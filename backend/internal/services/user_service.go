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
	GetUserByID(ctx context.Context, id string) (*models.UserFetchDTO, error)
	AuthenticateUser(ctx context.Context, authRequest *models.AuthRequest) (*models.UserFetchDTO, error)
	UpdateUser(ctx context.Context, id string, user *models.UserUpdateDTO) (*models.UserFetchDTO, error)
	SetDefaultProject(ctx context.Context, id string, projectID string) (*models.UserFetchDTO, error)
	GetUsers(ctx context.Context) ([]*models.UserFetchDTO, error)
	ListUsersByTeamID(ctx context.Context, teamID string) ([]*models.UserFetchDTO, error)
	DeleteUser(ctx context.Context, id string) error
	ChangePassword(ctx context.Context, id string, request *models.ChangePasswordRequest) error
	ChangeImage(ctx context.Context, id string, image *multipart.FileHeader) error
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

func (s *userService) GetUserByID(ctx context.Context, id string) (*models.UserFetchDTO, error) {
	claims := auth.GetUserClaims(ctx)
	u, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !(claims.IsSuperadmin() || (claims.IsAdmin() && u.TeamID != nil && claims.TeamID == *u.TeamID) || (claims.Subject == id)) {
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
	if err := s.userRepo.CreateUser(ctx, u); err != nil {
		return err
	}
	return nil
}

func (s *userService) UpdateUser(ctx context.Context, id string, user *models.UserUpdateDTO) (*models.UserFetchDTO, error) {
	u, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, apperrors.New(http.StatusNotFound, "user not found")
	}
	u.Role = user.Role
	u.Disabled = user.Disabled
	if user.TeamID != "" {
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

func (s *userService) GetUsers(ctx context.Context) ([]*models.UserFetchDTO, error) {
	claims := auth.GetUserClaims(ctx)
	var users []*models.User
	var err error

	if !claims.IsSuperadmin() {
		users, err = s.userRepo.ListUsersByTeamID(ctx, claims.TeamID)
		if err != nil {
			return nil, apperrors.New(http.StatusInternalServerError, "failed to list users")
		}
	} else {
		users, err = s.userRepo.ListUsers(ctx)
		if err != nil {
			return nil, apperrors.New(http.StatusInternalServerError, "failed to list users")
		}
	}

	var userDTOs []*models.UserFetchDTO
	for _, u := range users {
		userDTOs = append(userDTOs, u.ToFetchDTO())
	}
	return userDTOs, nil
}

func (s *userService) ListUsersByTeamID(ctx context.Context, teamID string) ([]*models.UserFetchDTO, error) {
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

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.DeleteUser(ctx, id)
}

func (s *userService) SetDefaultProject(ctx context.Context, id string, projectID string) (*models.UserFetchDTO, error) {
	claims := auth.GetUserClaims(ctx)
	if claims.Subject != id {
		return nil, apperrors.New(http.StatusUnauthorized, "unauthorized access")
	}

	p, err := s.projectRepo.GetProjectByID(ctx, projectID)
	if p == nil || err != nil {
		return nil, apperrors.New(http.StatusNotFound, "project not found")
	}

	if !claims.IsSuperadmin() && p.TeamID != claims.TeamID {
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

func (s *userService) ChangePassword(ctx context.Context, id string, request *models.ChangePasswordRequest) error {
	claims := auth.GetUserClaims(ctx)
	if claims.Subject != id {
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

func (s *userService) ChangeImage(ctx context.Context, id string, image *multipart.FileHeader) error {
	claims := auth.GetUserClaims(ctx)
	if claims.Subject != id {
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
