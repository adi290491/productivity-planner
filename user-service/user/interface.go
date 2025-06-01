package user

import "productivity-planner/user-service/models"

type UserServiceInterface interface {
	Signup(userDto SignupDTO) (*models.User, error)
	Login(loginDto LoginRequest) (*models.User, error)
}
