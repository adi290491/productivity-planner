package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/adi290491/productivity-planner/session-service/internal/service"
	"github.com/adi290491/productivity-planner/session-service/pkg/httperr"
)

// Handler handles HTTP requests for sessions
type Handler struct {
	sessionService service.SessionServiceInterface
}

// NewHandler creates a new HTTP handler
func NewHandler(sessionService service.SessionServiceInterface) *Handler {
	return &Handler{
		sessionService: sessionService,
	}
}

// StartSession handles POST /sessions/v1/start-session
func (h *Handler) StartSession(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from header (set by gateway)

	if r.Method != http.MethodPost {
		httperr.WriteError(w, fmt.Errorf("method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	userID := strings.TrimSpace(r.Header.Get("X-USER-ID"))
	if userID == "" {
		httperr.WriteError(w, fmt.Errorf("missing user ID"), http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		httperr.WriteError(w, fmt.Errorf("failed to read request body: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse request body
	var req service.StartSessionRequest
	if err := json.Unmarshal(body, &req); err != nil {
		httperr.WriteError(w, fmt.Errorf("invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate session type
	if !req.SessionType.IsValid() {
		httperr.WriteError(w, fmt.Errorf("invalid session type"), http.StatusBadRequest)
		return
	}

	// Start session
	response, err := h.sessionService.StartSession(req, userID)
	if err != nil {
		httperr.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	// Return response
	httperr.WriteJSON(w, response, http.StatusOK)
}

// StopSession handles PATCH /sessions/v1/stop-session
func (h *Handler) StopSession(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPatch {
		httperr.WriteError(w, fmt.Errorf("method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from header (set by gateway)
	userID := strings.TrimSpace(r.Header.Get("X-USER-ID"))
	if userID == "" {
		httperr.WriteError(w, fmt.Errorf("missing user ID"), http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		httperr.WriteError(w, fmt.Errorf("failed to read request body: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse request body
	var req service.StopSessionRequest
	if err := json.Unmarshal(body, &req); err != nil {
		httperr.WriteError(w, fmt.Errorf("invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate session type
	if !req.SessionType.IsValid() {
		httperr.WriteError(w, fmt.Errorf("invalid session type"), http.StatusBadRequest)
		return
	}

	// Stop session
	response, err := h.sessionService.StopSession(req, userID)
	if err != nil {
		httperr.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	// Return response
	httperr.WriteJSON(w, response, http.StatusOK)
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httperr.WriteError(w, fmt.Errorf("method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"status":  "healthy",
		"service": "session-service",
	}
	httperr.WriteJSON(w, response, http.StatusOK)
}

// ReadyCheck handles GET /ready
func (h *Handler) ReadyCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httperr.WriteError(w, fmt.Errorf("method not allowed"), http.StatusMethodNotAllowed)
		return
	}
	response := map[string]string{
		"status": "ready",
	}
	httperr.WriteJSON(w, response, http.StatusOK)
}

// WriteJSON is a helper to write JSON responses (for future net/http migration)
func WriteJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
