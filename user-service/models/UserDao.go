package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	DB *gorm.DB
}

func (p *PostgresRepository) CreateUser(user *User) (*User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result := p.DB.WithContext(ctx).Create(user)

	if result.Error != nil {
		return nil, fmt.Errorf("user creation failed. %v", result.Error.Error())
	}

	return user, nil
}

func (p *PostgresRepository) GetUser(userDao *User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	result := p.DB.WithContext(ctx).First(&user, "email = ?", userDao.Email)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, fmt.Errorf("error when fetching user: %v", result.Error)
	}

	return &user, nil
}
