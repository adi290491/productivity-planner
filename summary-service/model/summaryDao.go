package models

import (
	"context"
	"fmt"
	"time"
)

func (p *PostgresRepository) FindAllSessionsBetweenDates(summaryDao *Summary) ([]Session, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var sessions []Session

	err := p.DB.WithContext(ctx).
		Where("user_id = ? AND start_time >= ? AND end_time <= ?", summaryDao.UserId, summaryDao.StartTime, summaryDao.EndTime).
		Find(&sessions).Error

	if len(sessions) == 0 {
		return nil, fmt.Errorf("no sessions found for the given day")
	}

	if err != nil {
		return nil, err
	}

	return sessions, nil
}
