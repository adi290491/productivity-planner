package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/adi290491/productivity-planner/user-service/internal/model"
	"github.com/adi290491/productivity-planner/user-service/internal/service"
	"github.com/adi290491/productivity-planner/user-service/pkg/httperr"
	"github.com/adi290491/productivity-planner/user-service/pkg/util"
)

type Handler struct {
	userService service.UserService
	jwtUtil     *util.JWTUtil
}

var (
	dbInitialized atomic.Bool
	dbInitMutex   sync.Mutex
)

func NewHandler(svc service.UserService, jwtUtil *util.JWTUtil) *Handler {
	return &Handler{
		userService: svc,
		jwtUtil:     jwtUtil,
	}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"service":   "user-service",
		"timestamp": time.Now().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	response := map[string]any{
		"status": "ready",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	slog.Info("Signup request received")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var req model.SignupRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode signup request", "error", err)
		httperr.WriteErrorResponse(w, fmt.Errorf("invalid request body"), http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		httperr.WriteErrorResponse(w, fmt.Errorf("Email, password, and name are required"), http.StatusBadRequest)
		return
	}

	user, err := h.userService.Signup(ctx, &req)

	if err != nil {
		slog.Error("Signup failed", "error", err, "email", req.Email)
		httperr.WriteErrorResponse(w, fmt.Errorf("Failed to create user"), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]any{
		"user": model.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	slog.Info("Login request received")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var req model.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode login request", "error", err)
		httperr.WriteErrorResponse(w, fmt.Errorf("invalid request body"), http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		httperr.WriteErrorResponse(w, fmt.Errorf("email and password are required"), http.StatusBadRequest)
		return
	}

	if !isValidEmail(req.Email) {
		httperr.WriteErrorResponse(w, fmt.Errorf("invalid email format"), http.StatusBadRequest)
		return
	}

	user, err := h.userService.Login(ctx, &req)

	if err != nil {
		slog.Error("Login service failed",
			"email", req.Email,
			"error", err,
		)
		httperr.WriteErrorResponse(w, fmt.Errorf("invalid credentials"), http.StatusUnauthorized)
		return
	}

	token, err := h.jwtUtil.GenerateToken(user)

	if err != nil {
		slog.Error("Token generation failed",
			"userId", user.ID,
			"error", err,
		)
		httperr.WriteErrorResponse(w, fmt.Errorf("failed to generate token"), http.StatusInternalServerError)
		return
	}

	response := model.LoginResponse{
		Token: token,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetUsersBatch(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var req model.GetUsersBatchRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode batch request", "error", err)
		httperr.WriteErrorResponse(w, fmt.Errorf("invalid request body"), http.StatusBadRequest)
		return
	}

	if len(req.UserIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]model.UserResponse{})
		return
	}

	userInfos, err := h.userService.GetUsersBatch(ctx, &req)

	if err != nil {
		slog.Error("Get users batch failed", "error", err, "count", len(req.UserIDs))
		httperr.WriteErrorResponse(w, fmt.Errorf("Failed to retrieve users"), http.StatusInternalServerError)
		return
	}

	response := make([]model.UserResponse, len(userInfos))
	for i, info := range userInfos {
		response[i] = model.UserResponse{
			ID:    info.ID,
			Email: info.Email,
			Name:  info.Name,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func isValidEmail(email string) bool {
	if len(email) < 3 {
		return false
	}
	atIdx := -1
	for i, ch := range email {
		if ch == '@' {
			if atIdx != -1 {
				return false
			}
			atIdx = i
		}
	}

	if atIdx == -1 || atIdx == 0 || atIdx == len(email)-1 {
		return false
	}

	hasDotAfter := false
	for i := atIdx + 1; i < len(email); i++ {
		if email[i] == '.' {
			hasDotAfter = true
			break
		}
	}

	return hasDotAfter
}
