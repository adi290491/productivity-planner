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

func (s *SignupDTO) String() string {
	return fmt.Sprintf("{Name: %s,\tEmail: %s}", s.Name, s.Email)
}
