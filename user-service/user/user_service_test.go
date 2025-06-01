package user

import (
	"productivity-planner/user-service/models"
	"testing"
)

func TestUserService_Signup(t *testing.T) {
	svc := &UserService{Repo: &models.TestDBRepo{}}

	dto := SignupDTO{
		Email:    "alice@example.com",
		Name:     "Alice",
		Password: "password",
	}

	u, err := svc.Signup(dto)
	if err != nil {
		t.Fatalf("Signup failed: %v", err)
	}

	if u.Email != dto.Email {
		t.Errorf("expected email %s, got %s", dto.Email, u.Email)
	}
}

func TestUserService_Login_Success(t *testing.T) {
	svc := &UserService{Repo: &models.TestDBRepo{}}

	dto := LoginRequest{
		Email:    "alice@example.com",
		Password: "1234", // matches correctHash
	}

	u, err := svc.Login(dto)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if u.Email != dto.Email {
		t.Errorf("expected email %s, got %s", dto.Email, u.Email)
	}
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	svc := &UserService{Repo: &models.TestDBRepo{}}

	dto := LoginRequest{
		Email:    "alice@example.com",
		Password: "wrongpassword",
	}

	_, err := svc.Login(dto)
	if err == nil {
		t.Fatal("expected error for invalid password, got nil")
	}
}
