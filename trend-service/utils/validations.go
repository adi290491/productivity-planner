package utils

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

func ValidateDays(d string) (int, error) {
	if day, err := strconv.Atoi(d); err != nil {
		return int(math.NaN()), errors.New("invalid value for 'days'. A positive integer is required")
	} else if day < MIN_DAYS || day > MAX_DAYS {
		return int(math.NaN()), errors.New("invalid query parameter: days must be between 1 and 30")
	} else {
		return day, nil
	}
}

func ValidateWeeks(w string) (int, error) {
	if week, err := strconv.Atoi(w); err != nil {
		return int(math.NaN()), errors.New("invalid value for 'weeks'. A positive integer is required")
	} else if week < MIN_WEEKS || week > MAX_WEEKS {
		return int(math.NaN()), errors.New("invalid query parameter: weeks must be between 1 and 12")
	} else {
		return week, nil
	}
}

func FormatTimeToHrMin(t float64) string { //time in minutes

	//convert to hours and minutes
	hours := int(t) / 60
	minutes := int(t) % 60
	return fmt.Sprintf("%02dh%02dm", hours, minutes)
}
