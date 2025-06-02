package summary

import (
	"fmt"
	"log"
	models "productivity-planner/summary-service/model"
)

type SummaryService struct {
	Repo models.Repository
}

func (s *SummaryService) GetDailySessionSummary(userId string, date string) (*DailySessionSummary, error) {

	startDate, err := StartOfDayUTC(date)

	if err != nil {
		return nil, err
	}

	endDate := EndOfDayUTC(startDate)
	log.Println("Processing dates. Start Date:", startDate, "End Date:", endDate)
	summaryDao := &models.Summary{
		UserId:    userId,
		StartTime: startDate,
		EndTime:   endDate,
	}

	log.Println("Repository:", s.Repo)
	sessions, err := s.Repo.FindAllSessionsBetweenDates(summaryDao)
	if err != nil {
		return nil, err
	}

	sessionSummary := CalculateSummary(sessions, startDate)
	log.Println("Session Summary:", sessionSummary)
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

	sessions, err := s.Repo.FindAllSessionsBetweenDates(summaryDao)
	if err != nil {
		return nil, fmt.Errorf("no sessions found")
	}

	weeklySessionSummary := CalculateWeeklySummary(sessions, startDate, endDate)

	return weeklySessionSummary, nil

}
