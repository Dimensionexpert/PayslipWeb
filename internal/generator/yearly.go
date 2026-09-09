package generator

import (
	"database/sql"
	"fmt"

	"github.com/Dimensionexpert/payslip/internal/database"
	genexcel "github.com/Dimensionexpert/payslip/internal/genExcel"
	"github.com/Dimensionexpert/payslip/internal/models"
)

// GenerateYearlyPayslipForEmployee generates a yearly payslip for a single employee.
func GenerateYearlyPayslipForEmployee(
	db *sql.DB,
	yearlyTemplatePath string,
	outputDir string,
	employee models.Employee,
	financialYearStart int,
) (string, error) {
	yearly, err := database.GetYearlyPayslip(
		db,
		employee.ShalarthID,
		financialYearStart,
	)
	if err != nil {
		return "", fmt.Errorf("fetching yearly payslip: %w", err)
	}

	if len(yearly.Payslips) == 0 {
		return "", fmt.Errorf("no yearly payslip records")
	}

	yearlyPath, err := genexcel.GenerateYearlyPayslip(
		yearlyTemplatePath,
		outputDir,
		yearly,
		financialYearStart,
	)
	if err != nil {
		return "", fmt.Errorf("generating yearly payslip: %w", err)
	}

	return yearlyPath, nil
}

// ExportYearlyEmployeePayslips generates yearly payslips for all employees in the database.
func ExportYearlyEmployeePayslips(
	db *sql.DB,
	yearlyTemplatePath string,
	outputDir string,
	financialYearStart int,
) error {
	yearlyEmployees, err := database.GetEmployees(db)
	if err != nil {
		return fmt.Errorf("fetching employees for yearly payslips: %w", err)
	}

	if len(yearlyEmployees) == 0 {
		return fmt.Errorf("no employees found for yearly payslips")
	}

	yearlySuccess := 0
	yearlyFailed := 0

	for index, employee := range yearlyEmployees {
		fmt.Printf(
			"Generating yearly payslip %d/%d: %s\n",
			index+1,
			len(yearlyEmployees),
			employee.Name,
		)

		yearlyPath, err := GenerateYearlyPayslipForEmployee(
			db,
			yearlyTemplatePath,
			outputDir,
			employee,
			financialYearStart,
		)
		if err != nil {
			fmt.Printf(
				"FAILED: generating yearly payslip for %s: %v\n",
				employee.Name,
				err,
			)

			yearlyFailed++
			continue
		}

		fmt.Printf("Yearly payslip generated: %s\n", yearlyPath)
		yearlySuccess++
	}

	fmt.Printf(
		"Yearly payslip generation: %d succeeded, %d failed\n",
		yearlySuccess,
		yearlyFailed,
	)

	return nil
}
