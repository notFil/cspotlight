package services

import (
	"errors"

	"github.com/notFil/cspotlight/internal/common"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(user *models.UserCreateDTO) (bool, error)
	GetUserByID(id string, claims common.AuthClaims) (*models.UserFetchDTO, error)
	GetUserByUsername(username string) (*models.UserFetchDTO, error)
	AuthenticateUser(authRequest *models.AuthRequest) (*models.UserFetchDTO, error)
	UpdateUser(id string, user *models.UserUpdateDTO) (*models.UserFetchDTO, error)
	ListUsers(common.AuthClaims) ([]*models.UserFetchDTO, error)
	ListUsersByTeamID(teamID string, claims common.AuthClaims) ([]*models.UserFetchDTO, error)
	DeleteUser(id string, claims common.AuthClaims) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUserByID(id string, claims common.AuthClaims) (*models.UserFetchDTO, error) {
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}
	u.PasswordHash = string(hashedPassword)
	if err := s.userRepo.CreateUser(u); err != nil {
		return false, err
	}
	return true, nil
}

func (s *userService) UpdateUser(id string, user *models.UserUpdateDTO) (*models.UserFetchDTO, error) {
	u, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
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

func (s *userService) ListUsers(claims common.AuthClaims) ([]*models.UserFetchDTO, error) {
	var users []*models.User
	var err error

	if !claims.IsSuperadmin() {
		users, err = s.userRepo.ListUsersByTeamID(*claims.TeamID)
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

func (s *userService) ListUsersByTeamID(teamID string, claims common.AuthClaims) ([]*models.UserFetchDTO, error) {
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

func (s *userService) DeleteUser(id string, claims common.AuthClaims) error {
	u, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}
	if !claims.IsSuperadmin() && claims.TeamID != u.TeamID {
		return errors.New("Unauthorized")
	}

	return s.userRepo.DeleteUser(id)
}
