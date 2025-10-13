package user

import (
	"fmt"

	"github.com/adi290491/productivity-planner/user-service/models"
	"github.com/google/uuid"
)

// MockUserService implements UserServiceInterface for testing
type MockUserService struct {
	// Fields to control mock behavior
	ShouldFailSignup        bool
	ShouldFailLogin         bool
	ShouldFailGetUsersBatch bool
	LoginShouldReturnUser   *models.User
	SignupShouldReturnUser  *models.User
	GetUsersBatchResponse   *[]UserInfoResponse
}

func (m *MockUserService) Signup(userDto SignupDTO) (*models.User, error) {
	if m.ShouldFailSignup {
		return nil, fmt.Errorf("signup failed")
	}

	if m.SignupShouldReturnUser != nil {
		return m.SignupShouldReturnUser, nil
	}

	// Default successful response
	return &models.User{
		ID:           uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Email:        userDto.Email,
		Name:         userDto.Name,
		PasswordHash: "hashedpassword",
	}, nil
}

func (m *MockUserService) Login(loginDto LoginRequest) (*models.User, error) {
	if m.ShouldFailLogin {
		return nil, fmt.Errorf("login failed")
	}

	if m.LoginShouldReturnUser != nil {
		return m.LoginShouldReturnUser, nil
	}

	// Default successful response for test@example.com
	if loginDto.Email == "test@example.com" && loginDto.Password == "1234" {
		return &models.User{
			ID:           uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			Email:        loginDto.Email,
			Name:         "Test User",
			PasswordHash: "hashedpassword",
		}, nil
	}

	// Return error for invalid credentials
	return nil, fmt.Errorf("invalid credentials")
}

func (m *MockUserService) GetUsersBatch(batchDto GetUsersBatchRequest) (*[]UserInfoResponse, error) {
	if m.ShouldFailGetUsersBatch {
		return nil, fmt.Errorf("failed to get users batch")
	}

	if m.GetUsersBatchResponse != nil {
		return m.GetUsersBatchResponse, nil
	}

	// Default response with mock users
	response := make([]UserInfoResponse, 0, len(batchDto.UserIDs))
	for i, userID := range batchDto.UserIDs {
		response = append(response, UserInfoResponse{
			ID:    userID,
			Email: fmt.Sprintf("user%d@example.com", i+1),
			Name:  fmt.Sprintf("User %d", i+1),
		})
	}

	return &response, nil
}
