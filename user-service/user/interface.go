package user

import "github.com/adi290491/productivity-planner/user-service/models"

type UserServiceInterface interface {
	Signup(userDto SignupDTO) (*models.User, error)
	Login(loginDto LoginRequest) (*models.User, error)
	GetUsersBatch(batchDto GetUsersBatchRequest) (*[]UserInfoResponse, error)
}
