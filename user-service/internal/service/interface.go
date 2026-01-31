package service

import (
	"context"

	"github.com/adi290491/productivity-planner/user-service/internal/model"
)

// UserService defines the interface for user business logic
type UserService interface {
	// Signup creates a new user account
	Signup(ctx context.Context, req *model.SignupRequest) (*model.User, error)
	
	// Login authenticates a user and returns user details
	Login(ctx context.Context, req *model.LoginRequest) (*model.User, error)
	
	// GetUsersBatch retrieves multiple users by their IDs
	GetUsersBatch(ctx context.Context, req *model.GetUsersBatchRequest) ([]model.UserInfo, error)
}
