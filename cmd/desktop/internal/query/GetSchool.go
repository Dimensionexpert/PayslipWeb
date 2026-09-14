package query

import (
	"database/sql"

	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/dto"
)

func GetSchool(db *sql.DB, udise string) (dto.SchoolSummary, error) {
	var school dto.SchoolSummary

	err := db.QueryRow(`
		SELECT
			udise,
			school_name
		FROM schools
		WHERE udise = ?
	`, udise).Scan(
		&school.UDISECode,
		&school.Name,
	)

	return school, err
}
