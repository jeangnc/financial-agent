package date

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"regexp"
)

var months = map[string]int{
	"Jan": 1,
	"Fev": 2,
	"Mar": 3,
	"Abr": 4,
	"Mai": 5,
	"Jun": 6,
	"Jul": 7,
	"Ago": 8,
	"Set": 9,
	"Out": 10,
	"Nov": 11,
	"Dez": 12,
}

func ParseBrlDate(dateStr string) (time.Time, error) {

	var nonAlpha = regexp.MustCompile(`\d`)
	monthString := strings.TrimSpace(nonAlpha.ReplaceAllString(dateStr, ""))
	month, ok := months[monthString]
	if !ok {
		return time.Now(), fmt.Errorf("failed to convert month: %s", monthString)
	}

	var nonNumeric = regexp.MustCompile(`[a-zA-Z ]`)
	dayString := strings.TrimSpace(nonNumeric.ReplaceAllString(dateStr, ""))
	day, err := strconv.ParseInt(dayString, 10, 64)
	if err != nil {
		return time.Now(), err
	}

	return time.Date(2024, time.Month(month), int(day), 0, 0, 0, 0, time.Local), nil
}
