package repository

import "github.com/adi290491/productivity-planner/session-service/internal/model"

// SessionRepository defines the interface for session data access
type SessionRepository interface {
	Create(session *model.Session) (*model.Session, error)
	Stop(session *model.Session) (*model.Session, error)
}
