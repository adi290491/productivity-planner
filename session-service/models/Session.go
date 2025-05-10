package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	UserId      uuid.UUID
	SessionType string
	StartTime   time.Time
	EndTime     *time.Time
}
