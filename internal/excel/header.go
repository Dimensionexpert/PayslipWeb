package excel

import "strings"

// BuildHeaderMap builds a map of header names to their index in the row.
func BuildHeaderMap(row []string) map[string]int {
	headerMap := make(map[string]int)
	for i, header := range row {
		headerMap[strings.TrimSpace(header)] = i
	}

	return headerMap

}
