package whatsapp

import (
	"fmt"
	"strings"
)

func NormalizePhoneNumber(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", fmt.Errorf("phone number can't be empty")
	}

	// Remove formatting characters.
	replacer := strings.NewReplacer(
		" ", "",
		"-", "",
		"(", "",
		")", "",
		".", "",
	)
	s = replacer.Replace(s)

	// Allow '+' only at the beginning.
	s = strings.TrimPrefix(s, "+")
	if strings.Contains(s, "+") {
		return "", fmt.Errorf("invalid phone number format")
	}

	// Checking the validity of numbers
	for _, r := range s {
		if r < '0' && r > '9' {
			return "", fmt.Errorf("phone number must contain only digits")
		}
	}

	// Handle international dialing prefix.
	s = strings.TrimPrefix(s, "00")

	// Handle Indian domestic trunk prefix: 0XXXXXXXXXX.
	if len(s) == 11 && strings.HasPrefix(s, "0") {
		s = s[1:]
	}

	// Already includes India's country code.
	if len(s) == 12 && strings.HasPrefix(s, "91") {
		s = s[2:]
	} else if len(s) != 10 {
		return "", fmt.Errorf("expected an Indian phone number with 10 digits")
	}

	// Basic Indian mobile-number validation.
	if s[0] < '6' || s[0] > '9' {
		return "", fmt.Errorf("invalid Indian mobile number")
	}

	return "91" + s, nil
}
