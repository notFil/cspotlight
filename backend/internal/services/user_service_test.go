package services

import (
	"context"
	"errors"
	"testing"

	"cspotlight/internal/auth"
	"cspotlight/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a manual mock for UserRepository
type MockUserRepository struct {
	users map[uuid.UUID]*models.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[uuid.UUID]*models.User),
	}
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	if _, exists := m.users[user.ID]; exists {
		return errors.New("user already exists")
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
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

func (m *MockUserRepository) GetUserCount(ctx context.Context) (int64, error) {
	return int64(len(m.users)), nil
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	if _, exists := m.users[user.ID]; exists {
		m.users[user.ID] = user
		return nil
	}
	return errors.New("user not found")
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
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

func (m *MockUserRepository) ListUsersByTeamID(ctx context.Context, teamID uuid.UUID) ([]*models.User, error) {
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
	service := NewUserService(mockRepo)

	userDTO := &models.UserRegisterDTO{
		FirstName:       "John",
		LastName:        "Doe",
		Username:        "johndoe",
		Email:           "john@example.com",
		Password:        "P@ssword123",
		ConfirmPassword: "P@ssword123",
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
	if createdUser.PasswordHash == "P@ssword123" {
		t.Fatalf("expected password to be hashed")
	}

	// Test second user (should be disabled)
	userDTO2 := &models.UserRegisterDTO{
		FirstName:       "Jane",
		LastName:        "Doe",
		Username:        "janedoe",
		Email:           "jane@example.com",
		Password:        "P@ssword123",
		ConfirmPassword: "P@ssword123",
	}

	err = service.RegisterUser(context.Background(), userDTO2)
	if err != nil {
		t.Fatalf("expected no error for second user, got %v", err)
	}

	createdUser2, err := mockRepo.GetUserByUsername(context.Background(), "janedoe")
	if err != nil {
		t.Fatalf("expected to find second user, got error: %v", err)
	}
	if createdUser2.Disabled != nil && !*createdUser2.Disabled {
		t.Fatalf("expected second user to be disabled")
	}
}

func TestGetUserByID(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewUserService(mockRepo)

	userID := uuid.New()
	disabled := false
	user := &models.User{
		ID:       userID,
		Username: "testuser",
		Role:     "user",
		Disabled: &disabled,
	}
	mockRepo.CreateUser(context.Background(), user)

	// Test authorized access (same user)
	userData := auth.UserContext{ID: userID, Role: "user"}
	ctx := auth.ContextWithUser(context.Background(), &userData)
	fetchedUser, err := service.GetUserByID(ctx, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if fetchedUser.ID != userID {
		t.Errorf("expected user ID %s, got %s", userID, fetchedUser.ID)
	}

	// Test unauthorized access
	otherID := uuid.New()
	otherData := auth.UserContext{ID: otherID, Role: "user"}
	ctx = auth.ContextWithUser(context.Background(), &otherData)
	_, err = service.GetUserByID(ctx, userID)
	if err == nil {
		t.Fatal("expected unauthorized error, got nil")
	}
	if err.Error() != "unauthorized access" {
		t.Errorf("expected 'Unauthorized' error, got %v", err)
	}

	// Test admin access (should be allowed)
	adminID := uuid.New()
	adminData := auth.UserContext{ID: adminID, Role: "superadmin"}
	ctx = auth.ContextWithUser(context.Background(), &adminData)
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
	service := NewUserService(mockRepo)

	password := "P@ssword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	disabled := false

	user := &models.User{
		ID:           uuid.New(),
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
		Disabled:     &disabled,
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

func TestRegisterUser_PasswordComplexity(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewUserService(mockRepo)

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"valid complex password", "P@ssword123", false},
		{"too short", "P@s1", true},
		{"no uppercase", "p@ssword123", true},
		{"no lowercase", "P@SSWORD123", true},
		{"no digit", "P@ssword!!!", true},
		{"no special char", "Password123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userDTO := &models.UserRegisterDTO{
				FirstName:       "John",
				LastName:        "Doe",
				Username:        uuid.New().String(), // Random username to avoid conflict
				Email:           "john@example.com",
				Password:        tt.password,
				ConfirmPassword: tt.password,
			}
			err := service.RegisterUser(context.Background(), userDTO)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChangePassword_PasswordComplexity(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewUserService(mockRepo)

	userID := uuid.New()
	currentPassword := "OldP@ssword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(currentPassword), bcrypt.DefaultCost)
	disabled := false

	user := &models.User{
		ID:           userID,
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
		Disabled:     &disabled,
	}
	mockRepo.CreateUser(context.Background(), user)

	userData := auth.UserContext{ID: userID, Role: "user"}
	ctx := auth.ContextWithUser(context.Background(), &userData)

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"valid complex password", "NewP@ssword123", false},
		{"too short", "P@s1", true},
		{"no uppercase", "p@ssword123", true},
		{"no lowercase", "P@SSWORD123", true},
		{"no digit", "P@ssword!!!", true},
		{"no special char", "Password123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &models.ChangePasswordRequest{
				CurrentPassword:    currentPassword,
				NewPassword:        tt.password,
				ConfirmNewPassword: tt.password,
			}
			err := service.ChangePassword(ctx, userID, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
