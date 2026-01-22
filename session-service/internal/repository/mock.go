package repository

import (
	"errors"
	"fmt"

	"github.com/adi290491/productivity-planner/session-service/internal/model"
)

// MockSessionRepository is a mock implementation for testing
type MockSessionRepository struct {
	ActiveSessions map[string]*model.Session
}

// NewMockSessionRepository creates a new mock repository
func NewMockSessionRepository() *MockSessionRepository {
	return &MockSessionRepository{
		ActiveSessions: make(map[string]*model.Session),
	}
}

// Create simulates creating a session
func (r *MockSessionRepository) Create(session *model.Session) (*model.Session, error) {
	// Simulate: if user already has an active session, return error
	if existing, ok := r.ActiveSessions[session.UserID.String()]; ok && existing.EndTime == nil {
		return nil, errors.New("user already has an active session — please end it before starting a new one")
	}

	// Simulate: create a new session
	r.ActiveSessions[session.UserID.String()] = session
	return session, nil
}

// Stop simulates stopping a session
func (r *MockSessionRepository) Stop(session *model.Session) (*model.Session, error) {
	// Simulate: no active session found
	existing, ok := r.ActiveSessions[session.UserID.String()]
	if !ok || existing.EndTime != nil {
		return nil, fmt.Errorf("no active session found")
	}

	// Simulate: stop session
	existing.EndTime = session.EndTime
	return existing, nil
}
