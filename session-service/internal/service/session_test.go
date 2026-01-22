package service

import (
	"strings"
	"testing"
	"time"

	"github.com/adi290491/productivity-planner/session-service/internal/model"
	"github.com/adi290491/productivity-planner/session-service/internal/repository"
	"github.com/google/uuid"
)

func TestSessionService_StartSession_Success(t *testing.T) {
	// Arrange
	userID := uuid.New()
	repo := repository.NewMockSessionRepository()
	svc := NewSessionService(repo)
	req := StartSessionRequest{SessionType: model.SessionTypeFocus}

	// Act
	resp, err := svc.StartSession(req, userID.String())

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response, got nil")
	}
	if resp.Status != SessionStatusStarted {
		t.Errorf("expected status %v, got %v", SessionStatusStarted, resp.Status)
	}
	if resp.Session.SessionType != "focus" {
		t.Errorf("expected session type focus, got %v", resp.Session.SessionType)
	}
	if resp.Session.SessionID == "" {
		t.Errorf("expected non-empty session id")
	}
}

func TestSessionService_StartSession_InvalidUUID(t *testing.T) {
	// Arrange
	repo := repository.NewMockSessionRepository()
	svc := NewSessionService(repo)
	req := StartSessionRequest{SessionType: model.SessionTypeFocus}
	invalidUserID := "not-a-uuid"

	// Act
	resp, err := svc.StartSession(req, invalidUserID)

	// Assert
	if err == nil {
		t.Fatalf("expected error for invalid uuid, got nil")
	}
	if resp != nil {
		t.Errorf("expected nil response, got %+v", resp)
	}
	if !strings.Contains(err.Error(), "invalid user ID") {
		t.Errorf("expected error to contain 'invalid user ID', got %q", err.Error())
	}
}

func TestSessionService_StartSession_AlreadyActiveSession(t *testing.T) {
	// Arrange
	userID := uuid.New()
	repo := repository.NewMockSessionRepository()
	// Simulate an active session for this user
	activeSession := &model.Session{
		ID:          uuid.New(),
		UserID:      userID,
		SessionType: model.SessionTypeFocus,
		StartTime:   time.Now().Add(-1 * time.Hour),
		EndTime:     nil,
	}
	repo.ActiveSessions[userID.String()] = activeSession

	svc := NewSessionService(repo)
	req := StartSessionRequest{SessionType: model.SessionTypeFocus}

	// Act
	resp, err := svc.StartSession(req, userID.String())

	// Assert
	if err == nil {
		t.Fatalf("expected error for already active session, got nil")
	}
	if resp != nil {
		t.Errorf("expected nil response, got %+v", resp)
	}
	if !strings.Contains(err.Error(), "failed to create session") {
		t.Errorf("expected error to contain 'failed to create session', got %q", err.Error())
	}
}

func TestSessionService_StopSession_Success(t *testing.T) {
	// Arrange
	userID := uuid.New()
	startTime := time.Now().Add(-30 * time.Minute).UTC()
	activeSession := &model.Session{
		ID:          uuid.New(),
		UserID:      userID,
		SessionType: model.SessionTypeFocus,
		StartTime:   startTime,
		EndTime:     nil,
	}
	repo := repository.NewMockSessionRepository()
	repo.ActiveSessions[userID.String()] = activeSession

	svc := NewSessionService(repo)
	req := StopSessionRequest{SessionType: model.SessionTypeFocus}

	// Act
	resp, err := svc.StopSession(req, userID.String())

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response, got nil")
	}
	if resp.Status != SessionStatusEnded {
		t.Errorf("expected status %v, got %v", SessionStatusEnded, resp.Status)
	}
	if resp.Session.SessionType != "focus" {
		t.Errorf("expected session type focus, got %v", resp.Session.SessionType)
	}
	if resp.Session.SessionID == "" {
		t.Errorf("expected non-empty session id")
	}
	if resp.Session.EndTime == "" {
		t.Errorf("expected non-empty end time")
	}
}

func TestSessionService_StopSession_InvalidUUID(t *testing.T) {
	// Arrange
	repo := repository.NewMockSessionRepository()
	svc := NewSessionService(repo)
	req := StopSessionRequest{SessionType: model.SessionTypeFocus}
	invalidUserID := "not-a-uuid"

	// Act
	resp, err := svc.StopSession(req, invalidUserID)

	// Assert
	if err == nil {
		t.Fatalf("expected error for invalid uuid, got nil")
	}
	if resp != nil {
		t.Errorf("expected nil response, got %+v", resp)
	}
	if !strings.Contains(err.Error(), "invalid user ID") {
		t.Errorf("expected error to contain 'invalid user ID', got %q", err.Error())
	}
}

func TestSessionService_StopSession_NoActiveSession(t *testing.T) {
	// Arrange
	userID := uuid.New()
	repo := repository.NewMockSessionRepository()
	svc := NewSessionService(repo)
	req := StopSessionRequest{SessionType: model.SessionTypeFocus}

	// Act
	resp, err := svc.StopSession(req, userID.String())

	// Assert
	if err == nil {
		t.Fatalf("expected error for no active session, got nil")
	}
	if resp != nil {
		t.Errorf("expected nil response, got %+v", resp)
	}
	if !strings.Contains(err.Error(), "failed to stop session") {
		t.Errorf("expected error to contain 'failed to stop session', got %q", err.Error())
	}
}
