package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adi290491/productivity-planner/session-service/internal/model"
	"github.com/adi290491/productivity-planner/session-service/internal/repository"
	"github.com/adi290491/productivity-planner/session-service/internal/service"
)

func setupTestHandler() (*http.ServeMux, *Handler) {
	repo := repository.NewMockSessionRepository()
	svc := service.NewSessionService(repo)
	handler := NewHandler(svc)

	mux := http.NewServeMux()
	RegisterRoutes(mux, handler)

	return mux, handler
}

func TestHealthCheck(t *testing.T) {
	mux, _ := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("expected status 'healthy', got %v", response["status"])
	}
}

func TestReadyCheck(t *testing.T) {
	mux, _ := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] != "ready" {
		t.Errorf("expected status 'ready', got %v", response["status"])
	}
}

func TestStartSession_Success(t *testing.T) {
	mux, _ := setupTestHandler()

	body := map[string]string{
		"session_type": "focus",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/sessions/v1/start-session", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-USER-ID", "11111111-1111-1111-1111-111111111111")

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response service.SessionResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != service.SessionStatusStarted {
		t.Errorf("expected status 'started', got %v", response.Status)
	}
}

func TestStartSession_MissingUserID(t *testing.T) {
	mux, _ := setupTestHandler()

	body := map[string]string{
		"session_type": "focus",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/sessions/v1/start-session", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	// No X-USER-ID header

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestStartSession_InvalidBody(t *testing.T) {
	mux, _ := setupTestHandler()

	req := httptest.NewRequest(http.MethodPost, "/sessions/v1/start-session", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-USER-ID", "11111111-1111-1111-1111-111111111111")

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestStartSession_InvalidSessionType(t *testing.T) {
	mux, _ := setupTestHandler()

	body := map[string]string{
		"session_type": "invalid",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/sessions/v1/start-session", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-USER-ID", "11111111-1111-1111-1111-111111111111")

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestStopSession_Success(t *testing.T) {
	mux, _ := setupTestHandler()

	// First start a session
	startBody := map[string]string{
		"session_type": "focus",
	}
	startJSON, _ := json.Marshal(startBody)

	startReq := httptest.NewRequest(http.MethodPost, "/sessions/v1/start-session", bytes.NewBuffer(startJSON))
	startReq.Header.Set("Content-Type", "application/json")
	startReq.Header.Set("X-USER-ID", "11111111-1111-1111-1111-111111111111")

	startW := httptest.NewRecorder()
	mux.ServeHTTP(startW, startReq)

	// Then stop it
	stopBody := map[string]string{
		"session_type": "focus",
	}
	stopJSON, _ := json.Marshal(stopBody)

	stopReq := httptest.NewRequest(http.MethodPatch, "/sessions/v1/stop-session", bytes.NewBuffer(stopJSON))
	stopReq.Header.Set("Content-Type", "application/json")
	stopReq.Header.Set("X-USER-ID", "11111111-1111-1111-1111-111111111111")

	stopW := httptest.NewRecorder()
	mux.ServeHTTP(stopW, stopReq)

	if stopW.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", stopW.Code)
	}

	var response service.SessionResponse
	if err := json.NewDecoder(stopW.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != service.SessionStatusEnded {
		t.Errorf("expected status 'ended', got %v", response.Status)
	}
}

func TestStopSession_MissingUserID(t *testing.T) {
	mux, _ := setupTestHandler()

	body := map[string]string{
		"session_type": "focus",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPatch, "/sessions/v1/stop-session", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	// No X-USER-ID header

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestStopSession_InvalidSessionType(t *testing.T) {
	mux, _ := setupTestHandler()

	body := map[string]string{
		"session_type": "invalid",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPatch, "/sessions/v1/stop-session", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-USER-ID", "11111111-1111-1111-1111-111111111111")

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestSessionTypeValidation(t *testing.T) {
	tests := []struct {
		sessionType model.SessionType
		valid       bool
	}{
		{model.SessionTypeFocus, true},
		{model.SessionTypeMeeting, true},
		{model.SessionTypeBreak, true},
		{model.SessionType("invalid"), false},
		{model.SessionType(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.sessionType), func(t *testing.T) {
			if got := tt.sessionType.IsValid(); got != tt.valid {
				t.Errorf("SessionType(%q).IsValid() = %v, want %v", tt.sessionType, got, tt.valid)
			}
		})
	}
}
