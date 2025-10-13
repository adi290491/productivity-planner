package user

import (
	"fmt"
	"strings"
	"testing"

	"github.com/adi290491/productivity-planner/user-service/models"

	"github.com/google/uuid"
)

func TestUserService_Signup(t *testing.T) {
	svc := &UserService{Repo: &models.TestDBRepo{}}

	dto := SignupDTO{
		Email:    "alice@example.com",
		Name:     "Alice",
		Password: "password",
	}

	u, err := svc.Signup(dto)
	if err != nil {
		t.Fatalf("Signup failed: %v", err)
	}

	if u.Email != dto.Email {
		t.Errorf("expected email %s, got %s", dto.Email, u.Email)
	}
	if u.Name != dto.Name {
		t.Errorf("expected name %s, got %s", dto.Name, u.Name)
	}
	if u.ID == uuid.Nil {
		t.Error("expected non-zero user ID")
	}
}

func TestUserService_Login_Success(t *testing.T) {
	svc := &UserService{Repo: &models.TestDBRepo{}}

	dto := LoginRequest{
		Email:    "alice@example.com",
		Password: "1234", // matches correctHash
	}

	u, err := svc.Login(dto)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if u.Email != dto.Email {
		t.Errorf("expected email %s, got %s", dto.Email, u.Email)
	}
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	svc := &UserService{Repo: &models.TestDBRepo{}}

	dto := LoginRequest{
		Email:    "alice@example.com",
		Password: "wrongpassword",
	}

	_, err := svc.Login(dto)
	if err == nil {
		t.Fatal("expected error for invalid password, got nil")
	}
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	svc := &UserService{Repo: &models.TestDBRepo{}}

	dto := LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "1234",
	}

	_, err := svc.Login(dto)
	if err == nil {
		t.Fatal("expected error for non-existent user, got nil")
	}
}

func TestUserService_GetUsersBatch(t *testing.T) {
	// Mock the repository to return test data
	mockRepo := &MockRepository{
		users: []models.UserInfo{
			{
				ID:    uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				Email: "user1@example.com",
				Name:  "User 1",
			},
			{
				ID:    uuid.MustParse("22222222-2222-2222-2222-222222222222"),
				Email: "user2@example.com",
				Name:  "User 2",
			},
		},
	}

	svc := &UserService{Repo: mockRepo}

	dto := GetUsersBatchRequest{
		UserIDs: []uuid.UUID{
			uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		},
	}

	response, err := svc.GetUsersBatch(dto)
	if err != nil {
		t.Fatalf("GetUsersBatch failed: %v", err)
	}

	if response == nil {
		t.Fatal("expected response, got nil")
	}

	if len(*response) != 2 {
		t.Errorf("expected 2 users, got %d", len(*response))
	}

	// Verify user data
	users := *response
	if users[0].ID != dto.UserIDs[0] {
		t.Errorf("expected user ID %s, got %s", dto.UserIDs[0], users[0].ID)
	}
	if users[0].Email != "user1@example.com" {
		t.Errorf("expected user email user1@example.com, got %s", users[0].Email)
	}
}

func TestUserService_GetUsersBatch_Empty(t *testing.T) {
	mockRepo := &MockRepository{users: []models.UserInfo{}}
	svc := &UserService{Repo: mockRepo}

	dto := GetUsersBatchRequest{UserIDs: []uuid.UUID{}}

	response, err := svc.GetUsersBatch(dto)
	if err != nil {
		t.Fatalf("GetUsersBatch failed: %v", err)
	}

	if response == nil {
		t.Fatal("expected response, got nil")
	}

	if len(*response) != 0 {
		t.Errorf("expected 0 users, got %d", len(*response))
	}
}

func TestUserService_GetUsersBatch_RepositoryError(t *testing.T) {
	mockRepo := &MockRepository{shouldReturnError: true}
	svc := &UserService{Repo: mockRepo}

	dto := GetUsersBatchRequest{
		UserIDs: []uuid.UUID{
			uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		},
	}

	_, err := svc.GetUsersBatch(dto)
	if err == nil {
		t.Fatal("expected error from repository, got nil")
	}

	if !strings.Contains(err.Error(), "repository failed to get users by id") {
		t.Errorf("expected specific error message, got %v", err)
	}
}

func TestUserService_Signup_RepositoryError(t *testing.T) {
	mockRepo := &MockRepository{shouldReturnError: true}
	svc := &UserService{Repo: mockRepo}

	dto := SignupDTO{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	_, err := svc.Signup(dto)
	if err == nil {
		t.Fatal("expected error from repository, got nil")
	}

	if !strings.Contains(err.Error(), "user creation failed") {
		t.Errorf("expected specific error message, got %v", err)
	}
}

// Mock repository for comprehensive testing
type MockRepository struct {
	shouldReturnError bool
	users             []models.UserInfo
}

func (m *MockRepository) CreateUser(user *models.User) (*models.User, error) {
	if m.shouldReturnError {
		return nil, fmt.Errorf("database error")
	}
	return user, nil
}

func (m *MockRepository) GetUser(user *models.User) (*models.User, error) {
	if m.shouldReturnError {
		return nil, fmt.Errorf("database error")
	}

	if user.Email == "nonexistent@example.com" {
		return nil, fmt.Errorf("user not found")
	}

	return &models.User{
		ID:           uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Email:        user.Email,
		Name:         "Test User",
		PasswordHash: "$2a$10$DTHWAFgobsSCeqip6vROy.b8S0alUnaN7ickVmju2o52v8GhfNi1O", // "1234"
	}, nil
}

func (m *MockRepository) GetUsersById(users *models.UserBatch) (*[]models.UserInfo, error) {
	if m.shouldReturnError {
		return nil, fmt.Errorf("database error")
	}
	return &m.users, nil
}
