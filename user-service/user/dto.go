package user

import "fmt"

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
