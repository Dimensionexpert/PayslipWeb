package excel

import (
	"github.com/xuri/excelize/v2"
)

// Open opens an Excel file at the given path and returns the file handle and any error that occurred.
func Open(path string) (*excelize.File, error) {
	return excelize.OpenFile(path)
}
