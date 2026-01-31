package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/adi290491/productivity-planner/user-service/internal/model"
	"github.com/google/uuid"
)

// Mock repository for testing
type mockUserRepository struct {
	users            map[string]*model.User
	shouldFailCreate bool
	shouldFailGet    bool
	shouldFailBatch  bool
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*model.User),
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *model.User) error {
	if m.shouldFailCreate {
		return errors.New("create failed")
	}
	if _, exists := m.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	if m.shouldFailGet {
		return nil, errors.New("get failed")
	}
	user, exists := m.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *mockUserRepository) GetByIDs(ctx context.Context, userIDs []uuid.UUID) ([]model.UserInfo, error) {
	if m.shouldFailBatch {
		return nil, errors.New("batch failed")
	}

	var result []model.UserInfo
	for _, user := range m.users {
		for _, id := range userIDs {
			if user.ID == id {
				result = append(result, model.UserInfo{
					ID:    user.ID,
					Email: user.Email,
					Name:  user.Name,
				})
				break
			}
		}
	}
	return result, nil
}

func TestNewUserService(t *testing.T) {
	repo := newMockUserRepository()
	svc := NewUserService(repo)

	if svc == nil {
		t.Error("NewUserService() returned nil")
	}
}

func TestUserService_Signup(t *testing.T) {
	tests := []struct {
		name    string
		req     *model.SignupRequest
		wantErr bool
	}{
		{
			name: "valid signup",
			req: &model.SignupRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			req: &model.SignupRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			wantErr: true,
		},
	}

	repo := newMockUserRepository()
	svc := NewUserService(repo)
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := svc.Signup(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Signup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if user == nil {
					t.Error("Signup() returned nil user")
				}
				if user.Email != tt.req.Email {
					t.Errorf("Signup() email = %v, want %v", user.Email, tt.req.Email)
				}
				if user.Name != tt.req.Name {
					t.Errorf("Signup() name = %v, want %v", user.Name, tt.req.Name)
				}
				if user.ID == uuid.Nil {
					t.Error("Signup() returned nil UUID")
				}
				if user.PasswordHash == tt.req.Password {
					t.Error("Signup() did not hash password")
				}
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	repo := newMockUserRepository()
	svc := NewUserService(repo)
	ctx := context.Background()

	// Create a test user
	signupReq := &model.SignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	_, err := svc.Signup(ctx, signupReq)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name    string
		req     *model.LoginRequest
		wantErr bool
	}{
		{
			name: "valid login",
			req: &model.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "wrong password",
			req: &model.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			wantErr: true,
		},
		{
			name: "user not found",
			req: &model.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := svc.Login(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && user == nil {
				t.Error("Login() returned nil user")
			}
		})
	}
}

func TestUserService_GetUsersBatch(t *testing.T) {
	repo := newMockUserRepository()
	svc := NewUserService(repo)
	ctx := context.Background()

	// Create test users
	user1ID := uuid.New()
	user2ID := uuid.New()

	repo.users["user1@example.com"] = &model.User{
		ID:    user1ID,
		Email: "user1@example.com",
		Name:  "User One",
	}
	repo.users["user2@example.com"] = &model.User{
		ID:    user2ID,
		Email: "user2@example.com",
		Name:  "User Two",
	}

	tests := []struct {
		name      string
		req       *model.GetUsersBatchRequest
		wantCount int
		wantErr   bool
	}{
		{
			name: "get two users",
			req: &model.GetUsersBatchRequest{
				UserIDs: []uuid.UUID{user1ID, user2ID},
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "get one user",
			req: &model.GetUsersBatchRequest{
				UserIDs: []uuid.UUID{user1ID},
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "empty request",
			req: &model.GetUsersBatchRequest{
				UserIDs: []uuid.UUID{},
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "non-existent user",
			req: &model.GetUsersBatchRequest{
				UserIDs: []uuid.UUID{uuid.New()},
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, err := svc.GetUsersBatch(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUsersBatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(users) != tt.wantCount {
				t.Errorf("GetUsersBatch() count = %v, want %v", len(users), tt.wantCount)
			}
		})
	}
}

func TestUserService_WithRepositoryErrors(t *testing.T) {
	ctx := context.Background()

	t.Run("signup with repository error", func(t *testing.T) {
		repo := newMockUserRepository()
		repo.shouldFailCreate = true
		svc := NewUserService(repo)

		req := &model.SignupRequest{
			Email:    "test@example.com",
			Password: "password123",
			Name:     "Test User",
		}

		_, err := svc.Signup(ctx, req)
		if err == nil {
			t.Error("Signup() should fail with repository error")
		}
	})

	t.Run("login with repository error", func(t *testing.T) {
		repo := newMockUserRepository()
		repo.shouldFailGet = true
		svc := NewUserService(repo)

		req := &model.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		_, err := svc.Login(ctx, req)
		if err == nil {
			t.Error("Login() should fail with repository error")
		}
	})

	t.Run("batch with repository error", func(t *testing.T) {
		repo := newMockUserRepository()
		repo.shouldFailBatch = true
		svc := NewUserService(repo)

		req := &model.GetUsersBatchRequest{
			UserIDs: []uuid.UUID{uuid.New()},
		}

		_, err := svc.GetUsersBatch(ctx, req)
		if err == nil {
			t.Error("GetUsersBatch() should fail with repository error")
		}
	})
}

func TestUserService_ContextTimeout(t *testing.T) {
	repo := newMockUserRepository()
	svc := NewUserService(repo)

	// Create context with immediate timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(2 * time.Millisecond) // Ensure timeout

	req := &model.SignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	// This may or may not fail depending on timing, but should not panic
	_, _ = svc.Signup(ctx, req)
}
