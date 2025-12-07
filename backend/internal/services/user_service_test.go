package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/notFil/cspotlight/internal/auth"
	"github.com/notFil/cspotlight/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a manual mock for UserRepository
type MockUserRepository struct {
	users map[string]*models.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*models.User),
	}
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	if _, exists := m.users[user.ID]; exists {
		return errors.New("user already exists")
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	if _, exists := m.users[user.ID]; exists {
		m.users[user.ID] = user
		return nil
	}
	return errors.New("user not found")
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id string) error {
	if _, exists := m.users[id]; exists {
		delete(m.users, id)
		return nil
	}
	return errors.New("user not found")
}

func (m *MockUserRepository) ListUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserRepository) ListUsersByTeamID(ctx context.Context, teamID string) ([]*models.User, error) {
	var users []*models.User
	for _, user := range m.users {
		if user.TeamID != nil && *user.TeamID == teamID {
			users = append(users, user)
		}
	}
	return users, nil
}

func TestRegisterUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockProjectRepo := NewMockProjectRepository()
	service := NewUserService(mockRepo, mockProjectRepo)

	userDTO := &models.UserRegisterDTO{
		FirstName:       "John",
		LastName:        "Doe",
		Username:        "johndoe",
		Email:           "john@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	err := service.RegisterUser(context.Background(), userDTO)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify user was created with hashed password
	createdUser, err := mockRepo.GetUserByUsername(context.Background(), "johndoe")
	if err != nil {
		t.Fatalf("expected to find user, got error: %v", err)
	}
	if createdUser.PasswordHash == "password123" {
		t.Fatalf("expected password to be hashed")
	}
}

func TestGetUserByID(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockProjectRepo := NewMockProjectRepository()
	service := NewUserService(mockRepo, mockProjectRepo)

	userID := "test-id"
	user := &models.User{
		ID:       userID,
		Username: "testuser",
		Role:     "user",
	}
	mockRepo.CreateUser(context.Background(), user)

	// Test authorized access (same user)
	claims := auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: userID}, Role: "user"}
	ctx := auth.ContextWithClaims(context.Background(), &claims)
	fetchedUser, err := service.GetUserByID(ctx, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if fetchedUser.ID != userID {
		t.Errorf("expected user ID %s, got %s", userID, fetchedUser.ID)
	}

	// Test unauthorized access
	otherClaims := auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "other-id"}, Role: "user"}
	ctx = auth.ContextWithClaims(context.Background(), &otherClaims)
	_, err = service.GetUserByID(ctx, userID)
	if err == nil {
		t.Fatal("expected unauthorized error, got nil")
	}
	if err.Error() != "unauthorized access" {
		t.Errorf("expected 'Unauthorized' error, got %v", err)
	}

	// Test admin access (should be allowed)
	adminClaims := auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "admin-id"}, Role: "superadmin"}
	ctx = auth.ContextWithClaims(context.Background(), &adminClaims)
	fetchedUserAdmin, err := service.GetUserByID(ctx, userID)
	if err != nil {
		t.Fatalf("expected no error for admin, got %v", err)
	}
	if fetchedUserAdmin.ID != userID {
		t.Errorf("expected user ID %s, got %s", userID, fetchedUserAdmin.ID)
	}
}

func TestAuthenticateUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockProjectRepo := NewMockProjectRepository()
	service := NewUserService(mockRepo, mockProjectRepo)

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &models.User{
		ID:           "test-id",
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
	}
	mockRepo.CreateUser(context.Background(), user)

	// Test correct password
	authRequest := &models.AuthRequest{
		Username: "testuser",
		Password: password,
	}
	authenticatedUser, err := service.AuthenticateUser(context.Background(), authRequest)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if authenticatedUser.Username != "testuser" {
		t.Errorf("expected username testuser, got %s", authenticatedUser.Username)
	}

	// Test incorrect password
	badAuthRequest := &models.AuthRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}
	_, err = service.AuthenticateUser(context.Background(), badAuthRequest)
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}
