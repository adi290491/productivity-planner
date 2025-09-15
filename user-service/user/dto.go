package user

import "github.com/google/uuid"

type SignupDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// SafeString returns a safe, redacted representation of SignupDTO for logging
func (s *SignupDTO) SafeString() string {
	if s == nil {
		return "<nil SignupDTO>"
	}
	return "{SignupDTO: <redacted>}"
}

type GetUsersBatchRequest struct {
	UserIDs []uuid.UUID `json:"user_ids"`
}

type UserInfoResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}
