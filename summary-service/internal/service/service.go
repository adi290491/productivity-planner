package service

import (
	"context"

	"github.com/adi290491/productivity-planner/summary-service/internal/model"
)

type Service interface {
	GetDailySessionSummary(ctx context.Context, userID string, date string) (*model.DailySessionSummary, error)
	GetWeeklySessionSummary(ctx context.Context, userID string, date string) (*model.WeeklySessionSummary, error)
}
