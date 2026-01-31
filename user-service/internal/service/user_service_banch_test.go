package service

import (
	"context"
	"testing"

	"github.com/adi290491/productivity-planner/user-service/internal/model"
	"github.com/google/uuid"
)

func BenchmarkUserService_Signup(b *testing.B) {
	repo := newMockUserRepository()
	svc := NewUserService(repo)
	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := &model.SignupRequest{
			Email:    "bench@example.com",
			Password: "password123",
			Name:     "Bench User",
		}
		_, _ = svc.Signup(ctx, req)

		delete(repo.users, req.Email)
	}
}

func BenchmarkUserService_Login(b *testing.B) {
	repo := newMockUserRepository()
	svc := NewUserService(repo)
	ctx := context.Background()

	signupReq := &model.SignupRequest{
		Email:    "benchlogin@example.com",
		Password: "password123",
		Name:     "Bench Login User",
	}

	_, err := svc.Signup(ctx, signupReq)
	if err != nil {
		b.Fatalf("Failed to create test users: %v", err)
	}

	loginReq := &model.LoginRequest{
		Email:    "benchlogin@example.com",
		Password: "password123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.Login(ctx, loginReq)
	}
}

func BenchmarkUserService_GetUsersBatch(b *testing.B) {
	repo := newMockUserRepository()
	svc := NewUserService(repo)
	ctx := context.Background()

	userIDs := make([]uuid.UUID, 10)

	for i := 0; i < 10; i++ {
		id := uuid.New()
		userIDs[i] = id
		repo.users[id.String()] = &model.User{
			ID:    id,
			Email: id.String() + "@example.com",
			Name:  "Bench User",
		}
	}

	req := &model.GetUsersBatchRequest{
		UserIDs: userIDs,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.GetUsersBatch(ctx, req)
	}
}

func BenchmarkUserService_GetUsersBatch_Parallel(b *testing.B) {
	repo := newMockUserRepository()
	svc := NewUserService(repo)
	ctx := context.Background()

	// Pre-create users
	userIDs := make([]uuid.UUID, 10)
	for i := 0; i < 10; i++ {
		id := uuid.New()
		userIDs[i] = id
		repo.users[id.String()] = &model.User{
			ID:    id,
			Email: id.String() + "@example.com",
			Name:  "Bench User",
		}
	}

	req := &model.GetUsersBatchRequest{
		UserIDs: userIDs,
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = svc.GetUsersBatch(ctx, req)
		}
	})
}
