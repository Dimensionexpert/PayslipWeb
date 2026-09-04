package excel

import "github.com/xuri/excelize/v2"

// GetRows returns the rows of the given sheet from the given file.
func GetRows(f *excelize.File, sheet string) ([][]string, error) {
	return f.GetRows(sheet)
}
