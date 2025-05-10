package user

import (
	"fmt"
	"log"
	"productivity-planner/user-service/models"
	"productivity-planner/user-service/utils"

	"github.com/google/uuid"
)

type UserService struct {
	repo    models.Repository

}

func NewUserService(repo models.Repository) *UserService {
	return &UserService{
		repo:    repo,
	}
}

func (u *UserService) Signup(userDto SignupDTO) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(userDto.Password)

	if err != nil {
		log.Panic(err)
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        userDto.Email,
		PasswordHash: hashedPassword,
		Name:         userDto.Name,
	}

	response, err := u.repo.Create(user)

	if err != nil {
		return nil, fmt.Errorf("user creation failed: %v", err)
	}

	return response, nil
}

func (u *UserService) Login(loginDto LoginRequest) (*models.User, error) {

	userDao := &models.User{
		Email: loginDto.Email,
	}

	userEntity, err := u.repo.FetchUser(userDao)

	if err != nil {
		return nil, err
	}

	err = utils.VerifyPassword(loginDto.Password, userEntity.PasswordHash)

	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return userEntity, nil
}
