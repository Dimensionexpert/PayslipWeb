package query

import (
	"database/sql"

	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/dto"
)

func GetEmployee(
	db *sql.DB,
	shalarthID string,
) (dto.EmployeeSummary, error) {
	var employee dto.EmployeeSummary

	err := db.QueryRow(`
		SELECT
			e.shalarth_id,
			e.name,
			e.designation,
			s.school_name,
			s.udise
		FROM employees e
		LEFT JOIN schools s
			ON s.id = e.school_id
		WHERE e.shalarth_id = ?
	`, shalarthID).Scan(
		&employee.ShalarthID,
		&employee.Name,
		&employee.Designation,
		&employee.SchoolName,
		&employee.UDISECode,
	)

	return employee, err
}
