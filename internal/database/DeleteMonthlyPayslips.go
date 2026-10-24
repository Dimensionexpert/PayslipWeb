package database

import (
	"database/sql"
	"fmt"
)

func DeleteMonthlyPayslips(
	db *sql.DB,
	year int,
	month int,
) (int64, error) {

	result, err := db.Exec(`
		DELETE FROM payslip_records
		WHERE year=? AND month=?
		`, year, month)
	if err != nil {
		return 0, fmt.Errorf("deleting monthly payslip records: %w", err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("getting deleted payslip record count: %w", err)
	}

	return affectedRows, nil
}
