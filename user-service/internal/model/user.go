package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// User represents a user entity in the system
type User struct {
	ID           uuid.UUID
	Email        string
	Name         string
	PasswordHash string
	CreatedAt    time.Time
}

// String returns a string representation of the User (excludes password hash)
func (u User) String() string {
	return fmt.Sprintf(
		"User{ID: %s, Email: %s, Name: %s, CreatedAt: %s}",
		u.ID, u.Email, u.Name, u.CreatedAt.Format(time.RFC3339))
}

// UserInfo represents public user information (without sensitive data)
type UserInfo struct {
	ID    uuid.UUID
	Email string
	Name  string
}
