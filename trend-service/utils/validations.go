package utils

import (
	"fmt"
	"math"
	"strconv"
)

func ValidateDays(d string) (int, bool) {
	if day, err := strconv.Atoi(d); err != nil {
		return int(math.NaN()), false
	} else if day < 0 {
		return int(math.NaN()), false
	} else {
		return day, true
	}
}

func FormatTimeToHrMin(t float64) string { //time in minutes

	//convert to hours and minutes
	hours := int(t) / 60
	minutes := int(t) % 60
	return fmt.Sprintf("%02dh%02dm", hours, minutes)
}
