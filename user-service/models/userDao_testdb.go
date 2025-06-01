package models

import (
	"fmt"

	"github.com/google/uuid"
)

type TestDBRepo struct {
}

const correctHash = "$2a$10$DTHWAFgobsSCeqip6vROy.b8S0alUnaN7ickVmju2o52v8GhfNi1O"

func (r *TestDBRepo) CreateUser(user *User) (*User, error) {

	return &User{
		ID:           uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Email:        "alice@example.com",
		Name:         "Alice",
		PasswordHash: correctHash,
	}, nil
}

func (r *TestDBRepo) GetUser(user *User) (*User, error) {

	if user.Email == "nonexistent@example.com" {
		return nil, fmt.Errorf("user not found")
	}

	return &User{
		ID:           uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Email:        "alice@example.com",
		Name:         "Alice",
		PasswordHash: correctHash,
	}, nil
}
