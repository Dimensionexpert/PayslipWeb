package importer

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParsePeriod(filename string) (month, year int, err error) {
	filename = strings.TrimSpace(filename)
	filename = strings.TrimSuffix(filename, ".xlsx")

	parts := strings.Split(filename, "_")
	if len(parts) < 2 {
		return 0, 0, fmt.Errorf(
			"invalid filename %q: expected Month_Year format",
			filename,
		)
	}

	monthStr := parts[0]
	yearStr := parts[1]

	parsedMonth, err := time.Parse("January", monthStr)
	if err != nil {
		return 0, 0, fmt.Errorf(
			"invalid month %q: %w",
			monthStr,
			err,
		)
	}

	year, err = strconv.Atoi(yearStr)
	if err != nil {
		return 0, 0, fmt.Errorf(
			"invalid year %q: %w",
			yearStr,
			err,
		)
	}

	return int(parsedMonth.Month()), year, nil
}
