package excel

import (
	"fmt"
)

func FindHeaderRow(rows [][]string) (int, map[string]int, error) {
	requiredHeaders := []string{
		"SCHOOL UDISE CODE",
		"NAME OF SCHOOL",
		"SHALARTH ID",
	}

	for i, row := range rows {
		headerMap := BuildHeaderMap(row)

		matches := 0
		for _, header := range requiredHeaders {
			if _, exists := headerMap[header]; exists {
				matches++
			}
		}

		if matches == len(requiredHeaders) {
			return i, headerMap, nil
		}
	}

	return -1, nil, fmt.Errorf(
		"could not find payroll header row: required columns are missing",
	)
}
