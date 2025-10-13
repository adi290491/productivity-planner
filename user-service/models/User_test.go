package models

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestUser_String(t *testing.T) {
	t.Run("formats user string correctly", func(t *testing.T) {
		userID := uuid.New()
		user := User{
			ID:           userID,
			Email:        "test@example.com",
			Name:         "Test User",
			PasswordHash: "hashedpassword",
		}

		result := user.String()

		if !strings.Contains(result, "User{") {
			t.Error("Expected result to contain 'User{'")
		}
		if !strings.Contains(result, userID.String()) {
			t.Error("Expected result to contain user ID")
		}
		if !strings.Contains(result, "test@example.com") {
			t.Error("Expected result to contain email")
		}
		if !strings.Contains(result, "Test User") {
			t.Error("Expected result to contain name")
		}
		// Password hash should not be in the string representation
		if strings.Contains(result, "hashedpassword") {
			t.Error("Expected password hash to not be in string representation")
		}
	})

	t.Run("handles empty fields", func(t *testing.T) {
		user := User{
			ID:    uuid.Nil,
			Email: "",
			Name:  "",
		}

		result := user.String()

		if !strings.Contains(result, "User{") {
			t.Error("Expected result to contain 'User{'")
		}
		if !strings.Contains(result, "00000000-0000-0000-0000-000000000000") {
			t.Error("Expected result to contain nil UUID")
		}
	})
}

func TestUserBatch(t *testing.T) {
	t.Run("creates user batch correctly", func(t *testing.T) {
		userIDs := []uuid.UUID{
			uuid.New(),
			uuid.New(),
			uuid.New(),
		}

		batch := UserBatch{
			UserIDs: userIDs,
		}

		if len(batch.UserIDs) != 3 {
			t.Errorf("Expected 3 user IDs, got %d", len(batch.UserIDs))
		}

		for i, expectedID := range userIDs {
			if batch.UserIDs[i] != expectedID {
				t.Errorf("Expected user ID at index %d to be %s, got %s", i, expectedID, batch.UserIDs[i])
			}
		}
	})

	t.Run("handles empty user IDs", func(t *testing.T) {
		batch := UserBatch{
			UserIDs: []uuid.UUID{},
		}

		if len(batch.UserIDs) != 0 {
			t.Errorf("Expected 0 user IDs, got %d", len(batch.UserIDs))
		}
	})
}

func TestUserInfo(t *testing.T) {
	t.Run("creates user info correctly", func(t *testing.T) {
		userID := uuid.New()
		userInfo := UserInfo{
			ID:    userID,
			Email: "info@example.com",
			Name:  "Info User",
		}

		if userInfo.ID != userID {
			t.Errorf("Expected ID %s, got %s", userID, userInfo.ID)
		}
		if userInfo.Email != "info@example.com" {
			t.Errorf("Expected email 'info@example.com', got %s", userInfo.Email)
		}
		if userInfo.Name != "Info User" {
			t.Errorf("Expected name 'Info User', got %s", userInfo.Name)
		}
	})
}

func TestUser_Fields(t *testing.T) {
	t.Run("user struct has all required fields", func(t *testing.T) {
		userID := uuid.New()
		user := User{
			ID:           userID,
			Email:        "complete@example.com",
			Name:         "Complete User",
			PasswordHash: "hashed_password_here",
		}

		// Test that all fields are accessible and correctly set
		if user.ID != userID {
			t.Error("ID field not correctly set")
		}
		if user.Email != "complete@example.com" {
			t.Error("Email field not correctly set")
		}
		if user.Name != "Complete User" {
			t.Error("Name field not correctly set")
		}
		if user.PasswordHash != "hashed_password_here" {
			t.Error("PasswordHash field not correctly set")
		}
	})
}
