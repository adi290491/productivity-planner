package summary

import (
	"fmt"
	"log"
	models "productivity-planner/summary-service/model"
	"sort"
	"strings"
	"time"
)

func StartOfDayUTC(date string) (time.Time, error) {

	log.Println("Date I received:", date)
	var t time.Time
	var err error
	if len(date) == 0 {
		t = time.Now().UTC()
	} else {
		t, err = time.Parse("2006-01-02", date)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid date format")
		}
	}

	y, m, d := t.Date()

	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC), nil
}

func EndOfDayUTC(t time.Time) time.Time {
	y, m, d := t.Date()

	return time.Date(y, m, d, 23, 59, 59, 0, time.UTC)
}

func StartOfWeekUTC(date string) (time.Time, error) {

	var base time.Time
	var err error

	if date == "" {
		base = time.Now().UTC()
	} else {
		base, err = time.Parse("2006-01-02", date)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid date format")
		}
	}

	// Normalize to Monday
	weekday := int(base.Weekday())
	if weekday == 0 {
		weekday = 7 // treat Sunday as 7
	}
	monday := base.AddDate(0, 0, -weekday+1)
	return StartOfDayUTC(monday.Format("2006-01-02"))
}

func CalculateSummary(sessions []models.Session, date time.Time) *DailySessionSummary {

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
		out[strings.ToLower(k)] = durationToHumanFormat(v)
	}

	sessionSummary := &DailySessionSummary{
		Date:      date.Format("2006-01-02"),
		TotalTime: durationToHumanFormat(totalDuration),
		Breakdown: out,
	}

	return sessionSummary
}

func CalculateWeeklySummary(sessions []models.Session, startDate, endDate time.Time) *WeeklySessionSummary {

	var totalDuration time.Duration
	dailyMap := make(map[string][]models.Session)

	for _, s := range sessions {

		if s.EndTime != nil {
			date := s.StartTime.UTC().Format("2006-01-02")
			dailyMap[date] = append(dailyMap[date], s)
			totalDuration += s.EndTime.Sub(s.StartTime)
		}
	}

	var dailySummaries []*DailySessionSummary

	for dateStr, sessionGroup := range dailyMap {
		parsedDate, _ := time.Parse("2006-01-02", dateStr)
		dailySummary := CalculateSummary(sessionGroup, parsedDate)
		dailySummaries = append(dailySummaries, dailySummary)
	}

	sort.Slice(dailySummaries, func(i, j int) bool {
		return dailySummaries[i].Date < dailySummaries[j].Date
	})

	weeklySessionSummary := &WeeklySessionSummary{
		StartDate:      startDate.Format("2006-01-02"),
		EndDate:        endDate.Format("2006-01-02"),
		TotalTime:      durationToHumanFormat(totalDuration),
		DailySummaries: dailySummaries,
	}

	return weeklySessionSummary
}

func durationToHumanFormat(d time.Duration) string {
	return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
}
