package user

import (
	"fmt"
	"log/slog"

	"github.com/adi290491/productivity-planner/user-service/models"
	"github.com/adi290491/productivity-planner/user-service/utils"
	"github.com/google/uuid"
)

type UserService struct {
	Repo models.Repository
}

func (u *UserService) Signup(userDto SignupDTO) (*models.User, error) {
	slog.Info("UserService.Signup called", "email", userDto.Email)

	hashedPassword, err := utils.HashPassword(userDto.Password)
	if err != nil {
		slog.Error("Password hashing failed", "error", err)
		return nil, fmt.Errorf("hashing password failed: %w", err)
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        userDto.Email,
		PasswordHash: hashedPassword,
		Name:         userDto.Name,
	}

	response, err := u.Repo.CreateUser(user)

	if err != nil {
		slog.Error("User creation failed in repository",
			"email", userDto.Email,
			"error", err,
		)
		return nil, fmt.Errorf("user creation failed: %+w", err)
	}
	slog.Info("User created successfully", "userId", response.ID, "email", response.Email)
	return response, nil
}

func (u *UserService) Login(loginDto LoginRequest) (*models.User, error) {
	slog.Info("UserService.Login called", "email", loginDto.Email)

	userDao := &models.User{
		Email: loginDto.Email,
	}

	userEntity, err := u.Repo.GetUser(userDao)

	if err != nil {
		slog.Error("Failed to fetch user from repository",
			"email", loginDto.Email,
			"error", err,
		)
		return nil, fmt.Errorf("user not found: %w", err)
	}
	slog.Debug("User found, verifying password", "userId", userEntity.ID)
	err = utils.VerifyPassword(loginDto.Password, userEntity.PasswordHash)

	if err != nil {
		slog.Warn("Password verification failed",
			"email", loginDto.Email,
			"userId", userEntity.ID,
		)
		return nil, fmt.Errorf("invalid password")
	}
	slog.Info("Password verified successfully", "userId", userEntity.ID)
	return userEntity, nil
}

func (u *UserService) GetUsersBatch(batchDto GetUsersBatchRequest) (*[]UserInfoResponse, error) {
	slog.Info("UserService.GetUsersBatch called", "userCount", len(batchDto.UserIDs))

	batchUserDao := &models.UserBatch{
		UserIDs: batchDto.UserIDs,
	}
	userInfoResult, err := u.Repo.GetUsersById(batchUserDao)

	if err != nil {
		slog.Error("Repository failed to get users by id", "error", err)
		return nil, fmt.Errorf("repository failed to get users by id: %w", err)
	}

	var response []UserInfoResponse
	for _, res := range *userInfoResult {
		response = append(response, UserInfoResponse{
			ID:    res.ID,
			Email: res.Email,
			Name:  res.Name,
		})
	}

	slog.Info("Batch fetch successful", "foundUsers", len(response))
	return &response, nil
}
