package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TestDBRepo struct {
	DB *gorm.DB
}

func (p *TestDBRepo) FindAllSessionsBetweenDates(summaryDao *Summary) ([]Session, error) {

	return []Session{
		{
			ID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			UserId:      uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			SessionType: "focus",
			StartTime:   time.Now().Add(-48 * time.Hour),
			EndTime:     ptrTime(time.Now().Add(-46 * time.Hour)),
		},
	}, nil
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
