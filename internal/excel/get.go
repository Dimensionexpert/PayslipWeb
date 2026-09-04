package excel

// Get returns the value of the field with the given name from the given row using the header map.
func Get(row []string, headerMap map[string]int, fieldName string) string {
	index, ok := headerMap[fieldName]
	if !ok || index >= len(row) {
		return "'"
	}

	return row[index]
}
