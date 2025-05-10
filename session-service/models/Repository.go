package models

import "gorm.io/gorm"

type Repository interface {
	CreateSession(session *Session) (*Session, error)
	StopSession(session *Session) (*Session, error)
}

type postgresRepository struct {
	DB *gorm.DB
}

func NewPostgresRepository(conn *gorm.DB) Repository {
	return &postgresRepository{
		DB: conn,
	}
}
