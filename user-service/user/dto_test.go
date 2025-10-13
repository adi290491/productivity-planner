package user

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestSignupDTO_SafeString(t *testing.T) {
	t.Run("returns safe string for valid DTO", func(t *testing.T) {
		dto := &SignupDTO{
			Email:    "test@example.com",
			Password: "secretpassword",
			Name:     "Test User",
		}

		result := dto.SafeString()

		if !strings.Contains(result, "SignupDTO") {
			t.Error("Expected result to contain 'SignupDTO'")
		}
		if !strings.Contains(result, "redacted") {
			t.Error("Expected result to contain 'redacted'")
		}
		if strings.Contains(result, "secretpassword") {
			t.Error("Expected password to be redacted")
		}
		if strings.Contains(result, "test@example.com") {
			t.Error("Expected email to be redacted")
		}
	})

	t.Run("returns safe string for nil DTO", func(t *testing.T) {
		var dto *SignupDTO = nil

		result := dto.SafeString()

		if result != "<nil SignupDTO>" {
			t.Errorf("Expected '<nil SignupDTO>', got %s", result)
		}
	})
}

func TestUserInfoResponse(t *testing.T) {
	t.Run("creates user info response correctly", func(t *testing.T) {
		userID := uuid.New()
		response := UserInfoResponse{
			ID:    userID,
			Email: "user@example.com",
			Name:  "Test User",
		}

		if response.ID != userID {
			t.Errorf("Expected ID %s, got %s", userID, response.ID)
		}
		if response.Email != "user@example.com" {
			t.Errorf("Expected email 'user@example.com', got %s", response.Email)
		}
		if response.Name != "Test User" {
			t.Errorf("Expected name 'Test User', got %s", response.Name)
		}
	})
}

func TestGetUsersBatchRequest(t *testing.T) {
	t.Run("creates batch request correctly", func(t *testing.T) {
		userIDs := []uuid.UUID{
			uuid.New(),
			uuid.New(),
		}

		request := GetUsersBatchRequest{
			UserIDs: userIDs,
		}

		if len(request.UserIDs) != 2 {
			t.Errorf("Expected 2 user IDs, got %d", len(request.UserIDs))
		}
		if request.UserIDs[0] != userIDs[0] {
			t.Error("First user ID doesn't match")
		}
		if request.UserIDs[1] != userIDs[1] {
			t.Error("Second user ID doesn't match")
		}
	})

	t.Run("handles empty user IDs slice", func(t *testing.T) {
		request := GetUsersBatchRequest{
			UserIDs: []uuid.UUID{},
		}

		if len(request.UserIDs) != 0 {
			t.Errorf("Expected 0 user IDs, got %d", len(request.UserIDs))
		}
	})
}

func TestLoginRequest(t *testing.T) {
	t.Run("creates login request correctly", func(t *testing.T) {
		request := LoginRequest{
			Email:    "user@example.com",
			Password: "testpassword",
		}

		if request.Email != "user@example.com" {
			t.Errorf("Expected email 'user@example.com', got %s", request.Email)
		}
		if request.Password != "testpassword" {
			t.Errorf("Expected password 'testpassword', got %s", request.Password)
		}
	})
}
