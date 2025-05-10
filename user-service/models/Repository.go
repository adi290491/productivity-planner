package models

import "gorm.io/gorm"

type Repository interface {
	Create(user *User) (*User, error)
	FetchUser(user *User) (*User, error)
}

type postgresRepository struct {
	DB *gorm.DB
}

func NewPostgresRepository(conn *gorm.DB) Repository {
	return &postgresRepository{
		DB: conn,
	}
}
