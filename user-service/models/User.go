package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"primaryKey"`
	Email        string
	PasswordHash string
	Name         string
}
