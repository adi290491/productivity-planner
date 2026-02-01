package repository

import (
	"context"

	"github.com/adi290491/productivity-planner/summary-service/internal/model"
)

type Repository interface {
	FindSessionsBetweenDates(ctx context.Context, summary *model.Summary) ([]model.Session, error)
	Close() error
}
