package services

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"slices"

	"github.com/gabriel-vasile/mimetype"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/repositories"
	"github.com/notFil/cspotlight/pkg/auth"
	"github.com/notFil/cspotlight/pkg/errs"
	"golang.org/x/crypto/bcrypt"
)

const ImageMaxSize = 2 << 20

var ImageAllowedExts = []string{"image/jpeg", "image/png", "image/gif"}

type UserService interface {
	RegisterUser(user *models.UserRegisterDTO) error
	GetUserByID(id string, claims auth.Claims) (*models.UserFetchDTO, error)
	AuthenticateUser(authRequest *models.AuthRequest) (*models.UserFetchDTO, error)
	UpdateUser(id string, user *models.UserUpdateDTO) (*models.UserFetchDTO, error)
	SetDefaultProject(id string, projectID string, claims auth.Claims) (*models.UserFetchDTO, error)
	GetUsers(claims auth.Claims) ([]*models.UserFetchDTO, error)
	ListUsersByTeamID(teamID string, claims auth.Claims) ([]*models.UserFetchDTO, error)
	DeleteUser(id string) error
	ChangePassword(id string, request *models.ChangePasswordRequest, claims auth.Claims) error
	ChangeImage(id string, image *multipart.FileHeader, claims auth.Claims) error
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

func (s *userService) GetUserByID(id string, claims auth.Claims) (*models.UserFetchDTO, error) {
	u, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if !(claims.IsSuperadmin() || (claims.IsAdmin() && u.TeamID != nil && claims.TeamID == *u.TeamID) || (claims.Subject == id)) {
		return nil, errs.ErrUnauthorized
	}

	if u == nil {
		return nil, errs.ErrNotFound
	}

	return u.ToFetchDTO(), nil
}

func (s *userService) RegisterUser(user *models.UserRegisterDTO) error {
	u := user.ToUser()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	if err := s.userRepo.CreateUser(u); err != nil {
		return err
	}
	return nil
}

func (s *userService) UpdateUser(id string, user *models.UserUpdateDTO) (*models.UserFetchDTO, error) {
	u, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	u.Role = user.Role
	u.Disabled = user.Disabled
	if user.TeamID != "" {
		u.TeamID = &user.TeamID
	} else {
		u.TeamID = nil
	}
	if err := s.userRepo.UpdateUser(u); err != nil {
		return nil, err
	}
	return u.ToFetchDTO(), nil
}

func (s *userService) AuthenticateUser(authRequest *models.AuthRequest) (*models.UserFetchDTO, error) {
	u, err := s.userRepo.GetUserByUsername(authRequest.Username)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(authRequest.Password)); err != nil {
		return nil, err
	}
	return u.ToFetchDTO(), nil
}

func (s *userService) GetUsers(claims auth.Claims) ([]*models.UserFetchDTO, error) {
	var users []*models.User
	var err error

	if !claims.IsSuperadmin() {
		users, err = s.userRepo.ListUsersByTeamID(claims.TeamID)
		if err != nil {
			return nil, err
		}
	} else {
		users, err = s.userRepo.ListUsers()
		if err != nil {
			return nil, err
		}
	}

	var userDTOs []*models.UserFetchDTO
	for _, u := range users {
		userDTOs = append(userDTOs, u.ToFetchDTO())
	}
	return userDTOs, nil
}

func (s *userService) ListUsersByTeamID(teamID string, claims auth.Claims) ([]*models.UserFetchDTO, error) {
	users, err := s.userRepo.ListUsersByTeamID(teamID)
	if err != nil {
		return nil, err
	}

	var userDTOs []*models.UserFetchDTO
	for _, u := range users {
		userDTOs = append(userDTOs, u.ToFetchDTO())
	}
	return userDTOs, nil
}

func (s *userService) DeleteUser(id string) error {
	return s.userRepo.DeleteUser(id)
}

func (s *userService) SetDefaultProject(id string, projectID string, claims auth.Claims) (*models.UserFetchDTO, error) {
	if claims.Subject != id {
		return nil, errs.ErrUnauthorized
	}

	p, err := s.projectRepo.GetProjectByID(projectID)
	if p == nil || err != nil {
		return nil, errs.ErrNotFound
	}

	if !claims.IsSuperadmin() && p.TeamID != claims.TeamID {
		return nil, errs.ErrUnauthorized
	}
	u, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	u.DefaultProjectID = &projectID
	if err := s.userRepo.UpdateUser(u); err != nil {
		return nil, err
	}
	return u.ToFetchDTO(), nil
}

func (s *userService) ChangePassword(id string, request *models.ChangePasswordRequest, claims auth.Claims) error {
	if claims.Subject != id {
		return errs.ErrUnauthorized
	}
	u, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return errs.ErrNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(request.CurrentPassword)); err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return s.userRepo.UpdateUser(u)
}

func (s *userService) ChangeImage(id string, image *multipart.FileHeader, claims auth.Claims) error {
	if claims.Subject != id {
		return errs.ErrUnauthorized
	}

	if image.Size > ImageMaxSize {
		return errs.ErrImageTooLarge
	}

	f, err := image.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	mimeType, err := validateImageMimeType(f)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	strFile := base64.StdEncoding.EncodeToString(data)

	u, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}
	u.Image = fmt.Sprintf("data:%s;base64,%s", mimeType, strFile)
	return s.userRepo.UpdateUser(u)
}

func validateImageMimeType(f multipart.File) (string, error) {
	mimeType, err := mimetype.DetectReader(f)
	if err != nil {
		return "", err
	}
	if !slices.Contains(ImageAllowedExts, mimeType.String()) {
		return "", errs.ErrInvalidImageFormat
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	return mimeType.String(), nil
}
