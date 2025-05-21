package utils

import (
	"math"
	"strconv"
)

func ValidateDays(d string) (int, bool) {
	if day, err := strconv.Atoi(d); err != nil {
		return int(math.NaN()), false
	} else {
		return day, true
	}
}
