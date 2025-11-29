package services

import (
	"errors"

	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/repositories"
	"github.com/notFil/cspotlight/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(user *models.UserCreateDTO) (bool, error)
	GetUserByID(id string, claims auth.Claims) (*models.UserFetchDTO, error)
	GetUserByUsername(username string) (*models.UserFetchDTO, error)
	AuthenticateUser(authRequest *models.AuthRequest) (*models.UserFetchDTO, error)
	UpdateUser(id string, user *models.UserUpdateDTO, claims auth.Claims) (*models.UserFetchDTO, error)
	SetDefaultProject(projectID string, claims *auth.Claims) error
	GetUsers(claims auth.Claims) ([]*models.UserFetchDTO, error)
	ListUsersByTeamID(teamID string, claims auth.Claims) ([]*models.UserFetchDTO, error)
	DeleteUser(id string, claims auth.Claims) error
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
	if !claims.IsAdmin() && claims.UserID != id {
		return nil, errors.New("Unauthorized")
	}
	u, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return u.ToFetchDTO(), nil
}

func (s *userService) GetUserByUsername(username string) (*models.UserFetchDTO, error) {
	u, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return u.ToFetchDTO(), nil
}

func (s *userService) RegisterUser(user *models.UserCreateDTO) (bool, error) {
	u := user.ToUser()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}
	u.PasswordHash = string(hashedPassword)
	if err := s.userRepo.CreateUser(u); err != nil {
		return false, err
	}
	return true, nil
}

func (s *userService) UpdateUser(id string, user *models.UserUpdateDTO, claims auth.Claims) (*models.UserFetchDTO, error) {
	u, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	if !claims.IsAdmin() && claims.UserID != id {
		return nil, errors.New("Unauthorized")
	}
	u.Username = user.Username
	u.Email = user.Email
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

func (s *userService) DeleteUser(id string, claims auth.Claims) error {
	u, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}
	if !claims.IsAdmin() {
		if u.TeamID == nil || claims.TeamID != *u.TeamID {
			return errors.New("unauthorized")
		}
	}

	return s.userRepo.DeleteUser(id)
}

func (s *userService) SetDefaultProject(projectID string, claims *auth.Claims) error {
	p, err := s.projectRepo.GetProjectByID(projectID)
	if p == nil || err != nil {
		return errors.New("project not found")
	}
	if p.TeamID != claims.TeamID {
		return errors.New("unauthorized")
	}
	u, err := s.userRepo.GetUserByID(claims.UserID)
	if err != nil {
		return err
	}
	u.DefaultProjectID = &projectID
	return s.userRepo.UpdateUser(u)
}
