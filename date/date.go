package date

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jeangnc/financial-agent/regexp"
)

const REGEXP = `(?<day>[0-9]{2}) (?<month>[a-zA-Z]{3})`

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
	m, err := regexp.Match(REGEXP, dateStr)
	if err != nil {
		return time.Now(), err
	}

	monthString := strings.TrimSpace(m["month"])
	month, ok := months[monthString]
	if !ok {
		return time.Now(), fmt.Errorf("failed to convert month: %s", monthString)
	}

	dayString := strings.TrimSpace(m["day"])
	day, err := strconv.ParseInt(dayString, 10, 64)
	if err != nil {
		return time.Now(), err
	}

	return time.Date(2024, time.Month(month), int(day), 0, 0, 0, 0, time.Local), nil
}
