package models

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	DB *gorm.DB
}

func (p *PostgresRepository) CreateUser(user *User) (*User, error) {

	log.Println("---------Calling CreateUser---------")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if p.DB == nil {
		log.Printf("Database reference not available: %+v", p.DB)
		return nil, fmt.Errorf("database error: %+v", p.DB)
	}
	result := p.DB.WithContext(ctx).Create(user)

	if result.Error != nil {
		log.Println("--------User creation error--------")
		return nil, fmt.Errorf("user creation failed. %v", result.Error.Error())
	}

	return user, nil
}

func (p *PostgresRepository) GetUser(userDao *User) (*User, error) {
	log.Println("---------Calling GetUser---------")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	result := p.DB.WithContext(ctx).First(&user, "email = ?", userDao.Email)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println("--------Record not found--------")
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		log.Println("--------User fetch error--------")
		return nil, fmt.Errorf("error when fetching user: %w", result.Error)
	}

	return &user, nil
}
