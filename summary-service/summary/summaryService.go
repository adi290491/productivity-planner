package summary

import (
	"fmt"
	models "productivity-planner/summary-service/model"

	"time"
)

type SummaryService struct {
	repo models.Repository
}

func NewSummaryService(repo models.Repository) *SummaryService {
	return &SummaryService{
		repo: repo,
	}
}

func (s *SummaryService) GetDailySessionSummary(userId string, date time.Time) (*DailySessionSummary, error) {

	summaryDao := &models.Summary{
		UserId:    userId,
		StartTime: StartOfDayUTC(date),
		EndTime:   EndOfDayUTC(date),
	}

	sessions, err := s.repo.FindAllSessionsBetweenDates(summaryDao)
	if err != nil {
		return nil, err
	}

	sessionSummary := CalculateSummary(sessions, date)

	return sessionSummary, nil

}

func (s *SummaryService) GetWeeklySessionSummary(userId string, start string) (*WeeklySessionSummary, error) {

	startDate, err := StartOfWeekUTC(start)

	if err != nil {
		return nil, err
	}

	endDate := startDate.AddDate(0, 0, 7)

	summaryDao := &models.Summary{
		UserId:    userId,
		StartTime: startDate,
		EndTime:   endDate,
	}

	sessions, err := s.repo.FindAllSessionsBetweenDates(summaryDao)
	if err != nil {
		return nil, fmt.Errorf("no sessions found")
	}

	weeklySessionSummary := CalculateWeeklySummary(sessions, startDate, endDate)

	return weeklySessionSummary, nil

}
