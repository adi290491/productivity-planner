package repository

import (
	"context"

	"github.com/adi290491/productivity-planner/user-service/internal/model"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user in the database
	Create(ctx context.Context, user *model.User) error
	
	// GetByEmail retrieves a user by their email address
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	
	// GetByIDs retrieves multiple users by their IDs
	GetByIDs(ctx context.Context, userIDs []uuid.UUID) ([]model.UserInfo, error)
}
