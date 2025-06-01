package models

import (
	"github.com/google/uuid"
)

type TestDBRepo struct {
}

func (p *TestDBRepo) CreateUser(user *User) (*User, error) {

	return &User{
		ID:           uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Email:        "alice@example.com",
		Name:         "Alice",
		PasswordHash: "hashed_password_1",
	}, nil
}

func (p *TestDBRepo) GetUser(userDao *User) (*User, error) {
	return &User{
		ID:           uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Email:        "alice@example.com",
		Name:         "Alice",
		PasswordHash: "hashed_password_1",
	}, nil
}
