package genexcel

import (
	"math"
	"strconv"
	"strings"
)

func formatINR(amount float64) string {
	n := int64(math.Round(amount))
	negative := n < 0

	if negative {
		n = -n
	}

	s := strconv.FormatInt(n, 10)

	if len(s) > 3 {
		lastThree := s[len(s)-3:]
		remaining := s[:len(s)-3]

		var groups []string
		for len(remaining) > 2 {
			groups = append([]string{remaining[len(remaining)-2:]}, groups...)
			remaining = remaining[:len(remaining)-2]
		}

		if remaining != "" {
			groups = append([]string{remaining}, groups...)
		}

		s = strings.Join(groups, ",") + "," + lastThree
	}

	if negative {
		s = "-" + s
	}

	return "₹ " + s
}

func groupIndian(s string) string {
	if len(s) <= 3 {
		return s
	}
	return s[:2] + "," + groupIndian(s[2:])
}
