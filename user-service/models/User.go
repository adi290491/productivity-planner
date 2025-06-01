package models

import (
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"primaryKey"`
	Email        string
	PasswordHash string
	Name         string
}

func (u User) String() string {
	return fmt.Sprintf(
		"User{ID: %s, Email: %s, Name: %s}",
		u.ID, u.Email, u.Name)
}
