package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/adi290491/productivity-planner/summary-service/internal/model"
	"github.com/google/uuid"
)

// MockRepository is a mock implementation for testing
type MockRepository struct {
	Sessions []model.Session
	Err      error
}

// FindSessionsBetweenDates returns mock sessions
func (m *MockRepository) FindSessionsBetweenDates(ctx context.Context, summary *model.Summary) ([]model.Session, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	if summary.UserId == "notfound" {
		return nil, fmt.Errorf("no sessions found for the given day")
	}

	if len(m.Sessions) > 0 {
		return m.Sessions, nil
	}

	// Default mock data
	endTime := time.Now().Add(-46 * time.Hour)
	return []model.Session{
		{
			ID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			UserId:      uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			SessionType: "focus",
			StartTime:   time.Now().Add(-48 * time.Hour),
			EndTime:     &endTime,
		},
	}, nil
}

func (m *MockRepository) Close() error {
	return nil
}
