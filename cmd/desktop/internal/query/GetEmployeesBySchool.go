package query

import (
	"database/sql"

	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/dto"
)

func GetEmployeesBySchool(
	db *sql.DB,
	udise string,
) ([]dto.EmployeeSummary, error) {
	rows, err := db.Query(`
		SELECT
			e.shalarth_id,
			e.name,
			e.designation,
			s.school_name,
			s.udise
		FROM employees e
		JOIN schools s
			ON s.id = e.school_id
		WHERE s.udise = ?
		ORDER BY e.name
	`, udise)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []dto.EmployeeSummary

	for rows.Next() {
		var employee dto.EmployeeSummary

		err := rows.Scan(
			&employee.ShalarthID,
			&employee.Name,
			&employee.Designation,
			&employee.SchoolName,
			&employee.UDISECode,
		)
		if err != nil {
			return nil, err
		}

		employees = append(employees, employee)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}
