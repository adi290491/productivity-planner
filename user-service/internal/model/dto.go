package model

import "github.com/google/uuid"

// SignupRequest represents the request body for user signup
type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// GetUsersBatchRequest represents request to fetch multiple users
type GetUsersBatchRequest struct {
	UserIDs []uuid.UUID `json:"user_ids"`
}

// UserResponse represents the response body for user data
type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}

// LoginResponse represents the response body for login
type LoginResponse struct {
	Token string `json:"token"`
}

// ErrorResponse represents error response body
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
