package models

import "gorm.io/gorm"

type Repository interface {
	FindAllSessionsBetweenDates(summaryDao *Summary) ([]Session, error)
}

type PostgresRepository struct {
	DB *gorm.DB
}

func NewPostgresRepository(conn *gorm.DB) Repository {
	return &PostgresRepository{
		DB: conn,
	}
}
