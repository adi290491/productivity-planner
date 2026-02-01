package service

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/adi290491/productivity-planner/summary-service/internal/model"
	"github.com/adi290491/productivity-planner/summary-service/internal/repository"
	"github.com/adi290491/productivity-planner/summary-service/pkg/timeutil"
)

type SummaryService struct {
	repo repository.Repository
}

func NewSummaryService(repo repository.Repository) *SummaryService {
	return &SummaryService{
		repo: repo,
	}
}

func (s *SummaryService) GetDailySessionSummary(ctx context.Context, userID string, date string) (*model.DailySessionSummary, error) {

	startDate, err := timeutil.StartOfDayUTC(date)

	if err != nil {
		slog.Warn("Invalid date format for daily summary",
			"user_id", userID,
			"start_date", startDate,
			"error", err,
		)
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	endDate := timeutil.EndOfDayUTC(startDate)

	slog.Debug("Fetching daily sessions",
		"user_id", userID,
		"start_date", startDate,
		"end_date", endDate,
	)

	summaryDao := &model.Summary{
		UserId:    userID,
		StartTime: startDate,
		EndTime:   endDate,
	}

	sessions, err := s.repo.FindSessionsBetweenDates(ctx, summaryDao)
	if err != nil {
		slog.Error("Failed to fetch sessions from repository",
			"user_id", userID,
			"start_date", startDate,
			"end_date", endDate,
			"error", err,
		)
		return nil, fmt.Errorf("no sessions found: %w", err)
	}

	slog.Debug("Sessions retrieved from repository",
		"user_id", userID,
		"session_count", len(sessions),
	)

	return calculateDailySummary(sessions, startDate), nil

}

func (s *SummaryService) GetWeeklySessionSummary(ctx context.Context, userID string, date string) (*model.WeeklySessionSummary, error) {

	startDate, err := timeutil.StartOfDayUTC(date)

	if err != nil {
		slog.Warn("Invalid date format for weekly summary",
			"user_id", userID,
			"start_date", startDate,
			"error", err,
		)
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	endDate := startDate.AddDate(0, 0, 7)

	slog.Debug("Fetching weekly sessions",
		"user_id", userID,
		"week_start", startDate,
		"week_end", endDate,
	)

	summaryDao := &model.Summary{
		UserId:    userID,
		StartTime: startDate,
		EndTime:   endDate,
	}

	sessions, err := s.repo.FindSessionsBetweenDates(ctx, summaryDao)
	if err != nil {
		slog.Error("Failed to fetch sessions from repository",
			"user_id", userID,
			"week_start", startDate,
			"week_end", endDate,
			"error", err,
		)
		return nil, fmt.Errorf("no sessions found: %w", err)
	}

	slog.Debug("Sessions retrieved from repository",
		"user_id", userID,
		"session_count", len(sessions),
	)

	return calculateWeeklySummary(sessions, startDate, endDate), nil
}

func calculateDailySummary(sessions []model.Session, date time.Time) *model.DailySessionSummary {

	breakdown := make(map[string]time.Duration)
	var totalDuration time.Duration

	for _, s := range sessions {
		if s.EndTime != nil {
			duration := s.EndTime.Sub(s.StartTime)
			breakdown[s.SessionType] += duration
			totalDuration += duration
		}
	}

	out := map[string]string{}

	for k, v := range breakdown {
		out[strings.ToLower(k)] = timeutil.DurationToHumanFormat(v)
	}

	sessionSummary := &model.DailySessionSummary{
		Date:      date.Format("2006-01-02"),
		TotalTime: timeutil.DurationToHumanFormat(totalDuration),
		Breakdown: out,
	}

	slog.Debug("Daily summary calculated",
		"date", sessionSummary.Date,
		"total_time", sessionSummary.TotalTime,
		"session_types", len(out),
	)

	return sessionSummary
}

func calculateWeeklySummary(sessions []model.Session, startDate, endDate time.Time) *model.WeeklySessionSummary {

	var totalDuration time.Duration
	dailyMap := make(map[string][]model.Session)

	for _, s := range sessions {

		if s.EndTime != nil {
			date := s.StartTime.UTC().Format("2006-01-02")
			dailyMap[date] = append(dailyMap[date], s)
			totalDuration += s.EndTime.Sub(s.StartTime)
		}
	}

	var dailySummaries []*model.DailySessionSummary

	for dateStr, sessionGroup := range dailyMap {
		parsedDate, _ := time.Parse("2006-01-02", dateStr)
		dailySummary := calculateDailySummary(sessionGroup, parsedDate)
		dailySummaries = append(dailySummaries, dailySummary)
	}

	sort.Slice(dailySummaries, func(i, j int) bool {
		return dailySummaries[i].Date < dailySummaries[j].Date
	})

	weeklySessionSummary := &model.WeeklySessionSummary{
		StartDate:      startDate.Format("2006-01-02"),
		EndDate:        endDate.Format("2006-01-02"),
		TotalTime:      timeutil.DurationToHumanFormat(totalDuration),
		DailySummaries: dailySummaries,
	}

	slog.Debug("Weekly summary calculated",
		"start_date", weeklySessionSummary.StartDate,
		"end_date", weeklySessionSummary.EndDate,
		"total_time", weeklySessionSummary.TotalTime,
		"days_with_activity", len(dailySummaries),
	)

	return weeklySessionSummary
}
