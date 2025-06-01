package user

import (
	"fmt"
	"productivity-planner/user-service/models"

	"github.com/google/uuid"
)

type MockUserService struct{}

func (m *MockUserService) Signup(req SignupDTO) (*models.User, error) {
	return &models.User{
		ID:    uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Email: req.Email,
		Name:  req.Name,
	}, nil
}

func (m *MockUserService) Login(req LoginRequest) (*models.User, error) {
	if req.Email == "test@example.com" && req.Password == "1234" {
		return &models.User{
			ID:    uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			Email: req.Email,
			Name:  "Mock User",
		}, nil
	}
	return nil, fmt.Errorf("invalid credentials")
}
