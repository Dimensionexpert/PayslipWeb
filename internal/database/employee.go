package database

import (
	"database/sql"
	"fmt"

	"github.com/Dimensionexpert/payslip/internal/models"
)

func GetEmployees(db *sql.DB) ([]models.Employee, error) {
	rows, err := db.Query(`
		SELECT
			e.id,
			e.school_id,
			e.name,
			e.shalarth_id,
			e.gender,
			e.designation,

			COALESCE(e.gpf_no,''),
			COALESCE(e.dcps_no,''),
			COALESCE(e.pran_no,''),
			e.pan,
			e.aadhaar,

			e.mobile,
			e.email,

			e.ddo_bank_name,
			e.ddo_bank_account,
			e.ddo_bank_ifsc,

			e.bank_name,
			e.bank_account,
			e.bank_ifsc,
			e.branch_name,

			e.pay_matrix,

			s.udise
		FROM employees e
		LEFT JOIN schools s
			ON s.id = e.school_id
		ORDER BY e.id
	`)
	if err != nil {
		return nil, fmt.Errorf("querying employees: %w", err)
	}
	defer rows.Close()

	employees := make([]models.Employee, 0)

	for rows.Next() {
		var employee models.Employee

		err := rows.Scan(
			&employee.ID,
			&employee.SchoolID,
			&employee.Name,
			&employee.ShalarthID,
			&employee.Gender,
			&employee.Designation,

			&employee.GPFNo,
			&employee.DCPSNo,
			&employee.PRANNo,
			&employee.PAN,
			&employee.Aadhaar,

			&employee.Mobile,
			&employee.Email,

			&employee.DDOBankName,
			&employee.DDOBankAccount,
			&employee.DDOBankIFSC,

			&employee.BankName,
			&employee.BankAccount,
			&employee.BankIFSC,
			&employee.BranchName,

			&employee.PayMatrix,

			&employee.UDISECode,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning employee: %w", err)
		}

		employees = append(employees, employee)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating employees: %w", err)
	}

	return employees, nil
}
