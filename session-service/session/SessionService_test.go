package session

import (
	"productivity-planner/task-service/models"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSessionService_StartSession_Success(t *testing.T) {
	// Arrange
	userID := uuid.New()
	repo := &models.TestDBRepo{ActiveSessions: make(map[string]*models.Session)}
	service := &SessionService{Repo: repo}
	req := SessionRequest{SessionType: "FOCUS"}

	// Act
	resp, err := service.StartSession(req, userID.String())

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response, got nil")
	}
	if resp.Status != STARTED {
		t.Errorf("expected status %v, got %v", STARTED, resp.Status)
	}
	if resp.Session.SessionType != "FOCUS" {
		t.Errorf("expected session type FOCUS, got %v", resp.Session.SessionType)
	}
	if resp.Session.SessionId == "" {
		t.Errorf("expected non-empty session id")
	}
}

func TestSessionService_StartSession_InvalidUUID(t *testing.T) {
	// Arrange
	repo := &models.TestDBRepo{ActiveSessions: make(map[string]*models.Session)}
	service := &SessionService{Repo: repo}
	req := SessionRequest{SessionType: "FOCUS"}
	invalidUserID := "not-a-uuid"

	// Act
	resp, err := service.StartSession(req, invalidUserID)

	// Assert
	if err == nil {
		t.Fatalf("expected error for invalid uuid, got nil")
	}
	if resp != nil {
		t.Errorf("expected nil response, got %+v", resp)
	}
	if got, want := err.Error(), "uuid parse error"; !contains(got, want) {
		t.Errorf("expected error to contain %q, got %q", want, got)
	}
}

func TestSessionService_StartSession_AlreadyActiveSession(t *testing.T) {
	// Arrange
	userID := uuid.New()
	repo := &models.TestDBRepo{ActiveSessions: make(map[string]*models.Session)}
	// Simulate an active session for this user
	activeSession := &models.Session{
		ID:        uuid.New(),
		UserId:    userID,
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   nil,
	}
	repo.ActiveSessions[userID.String()] = activeSession

	service := &SessionService{Repo: repo}
	req := SessionRequest{SessionType: "FOCUS"}

	// Act
	resp, err := service.StartSession(req, userID.String())

	// Assert
	if err == nil {
		t.Fatalf("expected error for already active session, got nil")
	}
	if resp != nil {
		t.Errorf("expected nil response, got %+v", resp)
	}
	if got, want := err.Error(), "session creation error"; !contains(got, want) {
		t.Errorf("expected error to contain %q, got %q", want, got)
	}
}

func TestSessionService_StopSession_Success(t *testing.T) {
	// Arrange
	userID := uuid.New()
	startTime := time.Now().Add(-30 * time.Minute).UTC()
	activeSession := &models.Session{
		ID:          uuid.New(),
		UserId:      userID,
		SessionType: "FOCUS",
		StartTime:   startTime,
		EndTime:     nil,
	}
	repo := &models.TestDBRepo{ActiveSessions: map[string]*models.Session{
		userID.String(): activeSession,
	}}
	service := &SessionService{Repo: repo}
	req := SessionRequest{SessionType: "FOCUS"}

	// Act
	resp, err := service.StopSession(req, userID.String())

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response, got nil")
	}
	if resp.Status != ENDED {
		t.Errorf("expected status %v, got %v", ENDED, resp.Status)
	}
	if resp.Session.SessionType != "FOCUS" {
		t.Errorf("expected session type FOCUS, got %v", resp.Session.SessionType)
	}
	if resp.Session.SessionId == "" {
		t.Errorf("expected non-empty session id")
	}
	if resp.Session.EndTime == "" {
		t.Errorf("expected non-empty end time")
	}
}

func TestSessionService_StopSession_InvalidUUID(t *testing.T) {
	// Arrange
	repo := &models.TestDBRepo{ActiveSessions: make(map[string]*models.Session)}
	service := &SessionService{Repo: repo}
	req := SessionRequest{SessionType: "FOCUS"}
	invalidUserID := "not-a-uuid"

	// Act
	resp, err := service.StopSession(req, invalidUserID)

	// Assert
	if err == nil {
		t.Fatalf("expected error for invalid uuid, got nil")
	}
	if resp != nil {
		t.Errorf("expected nil response, got %+v", resp)
	}
	if got, want := err.Error(), "uuid parse error"; !contains(got, want) {
		t.Errorf("expected error to contain %q, got %q", want, got)
	}
}

func TestSessionService_StopSession_NoActiveSession(t *testing.T) {
	// Arrange
	userID := uuid.New()
	repo := &models.TestDBRepo{ActiveSessions: make(map[string]*models.Session)}
	service := &SessionService{Repo: repo}
	req := SessionRequest{SessionType: "FOCUS"}

	// Act
	resp, err := service.StopSession(req, userID.String())

	// Assert
	if err == nil {
		t.Fatalf("expected error for no active session, got nil")
	}
	if resp != nil {
		t.Errorf("expected nil response, got %+v", resp)
	}
	if got, want := err.Error(), "error while ending session"; !contains(got, want) {
		t.Errorf("expected error to contain %q, got %q", want, got)
	}
}

// contains checks if substr is within s.
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
