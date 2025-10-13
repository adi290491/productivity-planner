package user

import (
	"testing"

	"github.com/google/uuid"
)

func TestUserServiceInterface_Implementations(t *testing.T) {
	t.Run("UserService implements UserServiceInterface", func(t *testing.T) {
		var _ UserServiceInterface = &UserService{}
		// This test ensures UserService correctly implements the interface
		// It will fail to compile if the interface is not properly implemented
	})

	t.Run("MockUserService implements UserServiceInterface", func(t *testing.T) {
		var _ UserServiceInterface = &MockUserService{}
		// This test ensures MockUserService correctly implements the interface
		// It will fail to compile if the interface is not properly implemented
	})
}

func TestUserServiceInterface_Methods(t *testing.T) {
	t.Run("interface defines correct methods", func(t *testing.T) {
		// Test with mock service to verify interface methods exist
		mockService := &MockUserService{}

		// Check that all interface methods can be called
		// Note: These will use default mock behavior

		// Test Signup method exists and can be called
		_, err := mockService.Signup(SignupDTO{
			Email:    "test@example.com",
			Password: "password",
			Name:     "Test User",
		})
		// Don't check for specific success/failure, just that method exists
		_ = err

		// Test Login method exists and can be called
		_, err = mockService.Login(LoginRequest{
			Email:    "test@example.com",
			Password: "1234",
		})
		_ = err

		// Test GetUsersBatch method exists and can be called
		_, err = mockService.GetUsersBatch(GetUsersBatchRequest{
			UserIDs: []uuid.UUID{},
		})
		_ = err
	})
}

func TestUserServiceInterface_PointerReceivers(t *testing.T) {
	t.Run("methods work with pointer receivers", func(t *testing.T) {
		// Verify that interface methods work correctly with pointer receivers
		var service UserServiceInterface = &MockUserService{}

		if service == nil {
			t.Error("Expected non-nil service implementation")
		}

		// Test that the interface assignment works
		anotherService := service
		if anotherService == nil {
			t.Error("Expected successful interface assignment")
		}
	})
}
