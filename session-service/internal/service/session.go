package service

import (
	"fmt"
	"time"

	"github.com/adi290491/productivity-planner/session-service/internal/model"
	"github.com/adi290491/productivity-planner/session-service/internal/repository"
	"github.com/google/uuid"
)

// SessionService handles session business logic
type SessionService struct {
	repo repository.SessionRepository
}

// NewSessionService creates a new session service
func NewSessionService(repo repository.SessionRepository) *SessionService {
	return &SessionService{repo: repo}
}

// StartSessionRequest represents the request to start a session
type StartSessionRequest struct {
	SessionType model.SessionType `json:"session_type"`
}

// StopSessionRequest represents the request to stop a session
type StopSessionRequest struct {
	SessionType model.SessionType `json:"session_type"`
}

// SessionResponse represents the response for session operations
type SessionResponse struct {
	Status  SessionStatus `json:"status"`
	Session SessionDTO    `json:"session"`
}

// SessionDTO is the data transfer object for sessions
type SessionDTO struct {
	SessionID   string `json:"sessionId"`
	SessionType string `json:"type"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time,omitempty"`
}

// SessionStatus represents the status of a session
type SessionStatus string

const (
	SessionStatusStarted SessionStatus = "started"
	SessionStatusEnded   SessionStatus = "ended"
)

// StartSession starts a new session for a user
func (s *SessionService) StartSession(req StartSessionRequest, userID string) (*SessionResponse, error) {
	// Parse user ID
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Validate session type
	if !req.SessionType.IsValid() {
		return nil, fmt.Errorf("invalid session type: %s", req.SessionType)
	}

	// Create session model
	session := &model.Session{
		ID:          uuid.New(),
		UserID:      uid,
		SessionType: req.SessionType,
		StartTime:   time.Now().UTC(),
		EndTime:     nil,
	}

	// Save to repository
	created, err := s.repo.Create(session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Build response
	response := &SessionResponse{
		Status: SessionStatusStarted,
		Session: SessionDTO{
			SessionID:   created.ID.String(),
			SessionType: created.SessionType.String(),
			StartTime:   created.StartTime.Format(time.RFC3339),
			EndTime:     "",
		},
	}

	return response, nil
}

// StopSession stops an active session for a user
func (s *SessionService) StopSession(req StopSessionRequest, userID string) (*SessionResponse, error) {
	// Parse user ID
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Validate session type
	if !req.SessionType.IsValid() {
		return nil, fmt.Errorf("invalid session type: %s", req.SessionType)
	}

	// Create session with end time
	endTime := time.Now().UTC()
	session := &model.Session{
		UserID:      uid,
		SessionType: req.SessionType,
		EndTime:     &endTime,
	}

	// Update in repository
	stopped, err := s.repo.Stop(session)
	if err != nil {
		return nil, fmt.Errorf("failed to stop session: %w", err)
	}

	// Build response
	response := &SessionResponse{
		Status: SessionStatusEnded,
		Session: SessionDTO{
			SessionID:   stopped.ID.String(),
			SessionType: stopped.SessionType.String(),
			StartTime:   stopped.StartTime.Format(time.RFC3339),
			EndTime:     stopped.EndTime.Format(time.RFC3339),
		},
	}

	return response, nil
}
