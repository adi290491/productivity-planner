package util

import (
	"testing"

	"github.com/adi290491/productivity-planner/user-service/internal/model"
	"github.com/google/uuid"
)

func BenchmarkJWT_GenerateToken(b *testing.B) {
	jwtUtil := NewJWTUtil([]byte("test-secret"))
	user := &model.User{
		ID:    uuid.New(),
		Email: "bench@example.com",
		Name:  "Bench User",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = jwtUtil.GenerateToken(user)
	}
}

func BenchmarkJWT_ValidateToken(b *testing.B) {
	jwtUtil := NewJWTUtil([]byte("test-secret"))
	user := &model.User{
		ID:    uuid.New(),
		Email: "bench@example.com",
		Name:  "Bench User",
	}

	token, err := jwtUtil.GenerateToken(user)
	if err != nil {
		b.Fatalf("Failed to generate token: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = jwtUtil.ValidateToken(token)
	}
}

func BenchmarkJWT_GenerateToken_Parallel(b *testing.B) {
	jwtUtil := NewJWTUtil([]byte("test-secret"))
	user := &model.User{
		ID:    uuid.New(),
		Email: "bench@example.com",
		Name:  "Bench User",
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = jwtUtil.GenerateToken(user)
		}
	})
}

func BenchmarkJWT_ValidateToken_Parallel(b *testing.B) {
	jwtUtil := NewJWTUtil([]byte("test-secret"))
	user := &model.User{
		ID:    uuid.New(),
		Email: "bench@example.com",
		Name:  "Bench User",
	}

	token, err := jwtUtil.GenerateToken(user)
	if err != nil {
		b.Fatalf("Failed to generate token: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = jwtUtil.ValidateToken(token)
		}
	})
}
