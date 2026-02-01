package timeutil

import (
	"fmt"

	"time"
)

func StartOfDayUTC(date string) (time.Time, error) {
	var t time.Time
	var err error

	if len(date) == 0 {
		t = time.Now().UTC()
	} else {
		t, err = time.Parse("2006-01-02", date)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid date format: %w", err)
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

func DurationToHumanFormat(d time.Duration) string {
	return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
}
