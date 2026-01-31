package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/adi290491/productivity-planner/user-service/internal/model"
	"github.com/adi290491/productivity-planner/user-service/internal/repository"
	"github.com/adi290491/productivity-planner/user-service/pkg/util"
	"github.com/google/uuid"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Signup(ctx context.Context, req *model.SignupRequest) (*model.User, error) {
	slog.Info("Processing signup request", "email", req.Email)

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		ID:           uuid.New(),
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		slog.Error("Failed to create uer in repository", "error", err, "email", req.Email)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	slog.Info("User signup successful", "userId", user.ID, "email", user.Email)
	return user, nil
}

func (s *userService) Login(ctx context.Context, req *model.LoginRequest) (*model.User, error) {
	slog.Info("Processing login request", "email", req.Email)
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		slog.Error("User not found for login", "email", req.Email)
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := util.VerifyPassword(req.Password, user.PasswordHash); err != nil {
		slog.Warn("Invalid password attempt", "email", req.Email, "userId", user.ID)
		return nil, fmt.Errorf("invalid credentials")
	}

	slog.Info("User login successful", "userId", user.ID, "email", user.Email)
	return user, nil
}

func (s *userService) GetUsersBatch(ctx context.Context, req *model.GetUsersBatchRequest) ([]model.UserInfo, error) {
	slog.Info("Processing batch request", "count", len(req.UserIDs))
	if len(req.UserIDs) == 0 {
		slog.Warn("Empty batch request")
		return []model.UserInfo{}, nil
	}

	userInfos, err := s.repo.GetByIDs(ctx, req.UserIDs)
	if err != nil {
		slog.Error("Failed to get users batch", "error", err, "count", len(req.UserIDs))
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	slog.Info("Batch request successful", "count", len(userInfos))
	return userInfos, nil
}
