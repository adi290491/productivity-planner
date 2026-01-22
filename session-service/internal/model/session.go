package model

import (
	"time"

	"github.com/google/uuid"
)

// Session represents a user's focus/break/meeting session
type Session struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	SessionType SessionType
	StartTime   time.Time
	EndTime     *time.Time
}

// SessionType represents the type of session
type SessionType string

const (
	SessionTypeFocus   SessionType = "focus"
	SessionTypeMeeting SessionType = "meeting"
	SessionTypeBreak   SessionType = "break"
)

// IsValid checks if the session type is valid
func (t SessionType) IsValid() bool {
	switch t {
	case SessionTypeFocus, SessionTypeMeeting, SessionTypeBreak:
		return true
	default:
		return false
	}
}

// String returns the string representation of SessionType
func (t SessionType) String() string {
	return string(t)
}
