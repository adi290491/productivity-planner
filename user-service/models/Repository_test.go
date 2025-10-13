package models

import (
	"testing"

	"github.com/google/uuid"
)

// Test that Repository interface is properly defined
func TestRepositoryInterface(t *testing.T) {
	t.Run("repository interface exists and is implementable", func(t *testing.T) {
		// This test ensures that the Repository interface is properly defined
		// and can be implemented by various types

		// Create a mock repository to test interface compliance
		mockRepo := &MockRepository{
			Users:        make(map[string]*User),
			ShouldFail:   false,
			FailureError: nil,
		}

		// Test that mockRepo implements Repository interface
		var repo Repository = mockRepo

		if repo == nil {
			t.Error("Repository interface is not properly implemented")
		}
	})
}

// MockRepository for testing purposes
type MockRepository struct {
	Users        map[string]*User
	ShouldFail   bool
	FailureError error
}

func (m *MockRepository) CreateUser(user *User) (*User, error) {
	if m.ShouldFail {
		return nil, m.FailureError
	}

	// Assign ID if not set
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	m.Users[user.Email] = user
	return user, nil
}

func (m *MockRepository) GetUser(user *User) (*User, error) {
	if m.ShouldFail {
		return nil, m.FailureError
	}

	for _, storedUser := range m.Users {
		if storedUser.Email == user.Email {
			return storedUser, nil
		}
	}

	return nil, nil
}

func (m *MockRepository) GetUsersById(batch *UserBatch) (*[]UserInfo, error) {
	if m.ShouldFail {
		return nil, m.FailureError
	}

	var userInfos []UserInfo
	for _, id := range batch.UserIDs {
		for _, user := range m.Users {
			if user.ID == id {
				userInfo := UserInfo{
					ID:    user.ID,
					Email: user.Email,
					Name:  user.Name,
				}
				userInfos = append(userInfos, userInfo)
				break
			}
		}
	}

	return &userInfos, nil
}

func TestMockRepository(t *testing.T) {
	t.Run("mock repository creates users", func(t *testing.T) {
		repo := &MockRepository{
			Users:      make(map[string]*User),
			ShouldFail: false,
		}

		user := &User{
			ID:           uuid.New(),
			Email:        "mock@example.com",
			Name:         "Mock User",
			PasswordHash: "hashed",
		}

		createdUser, err := repo.CreateUser(user)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if createdUser == nil {
			t.Error("Expected created user, got nil")
		}

		searchUser := &User{Email: "mock@example.com"}
		retrievedUser, err := repo.GetUser(searchUser)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if retrievedUser == nil {
			t.Error("Expected user, got nil")
		}
		if retrievedUser.Email != "mock@example.com" {
			t.Errorf("Expected email 'mock@example.com', got %s", retrievedUser.Email)
		}
	})

	t.Run("mock repository handles batch retrieval", func(t *testing.T) {
		repo := &MockRepository{
			Users:      make(map[string]*User),
			ShouldFail: false,
		}

		user1 := &User{
			ID:           uuid.New(),
			Email:        "user1@example.com",
			Name:         "User 1",
			PasswordHash: "hashed1",
		}
		user2 := &User{
			ID:           uuid.New(),
			Email:        "user2@example.com",
			Name:         "User 2",
			PasswordHash: "hashed2",
		}

		repo.CreateUser(user1)
		repo.CreateUser(user2)

		batch := &UserBatch{
			UserIDs: []uuid.UUID{user1.ID, user2.ID},
		}

		userInfos, err := repo.GetUsersById(batch)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if userInfos == nil {
			t.Error("Expected user infos, got nil")
		}
		if len(*userInfos) != 2 {
			t.Errorf("Expected 2 users, got %d", len(*userInfos))
		}
	})

	t.Run("mock repository simulates failures", func(t *testing.T) {
		expectedError := Error("simulated failure")
		repo := &MockRepository{
			Users:        make(map[string]*User),
			ShouldFail:   true,
			FailureError: expectedError,
		}

		user := &User{
			ID:           uuid.New(),
			Email:        "fail@example.com",
			Name:         "Fail User",
			PasswordHash: "hashed",
		}

		_, err := repo.CreateUser(user)
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err != expectedError {
			t.Errorf("Expected specific error, got %v", err)
		}

		searchUser := &User{Email: "fail@example.com"}
		_, err = repo.GetUser(searchUser)
		if err == nil {
			t.Error("Expected error, got nil")
		}

		batch := &UserBatch{UserIDs: []uuid.UUID{user.ID}}
		_, err = repo.GetUsersById(batch)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

// Error implements a simple error type for testing
type Error string

func (e Error) Error() string {
	return string(e)
}
