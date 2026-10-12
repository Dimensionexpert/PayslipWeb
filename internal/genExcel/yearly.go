package genexcel

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Dimensionexpert/payslip/internal/excel"
	"github.com/Dimensionexpert/payslip/internal/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func GenerateYearlyPayslip(
	templatePath string,
	outputDir string,
	data models.PayslipExportYear,
	startYear int,
) (string, error) {
	// Open the yearly template.
	yearlyTemplate, err := excel.Open(templatePath)
	if err != nil {
		return "", fmt.Errorf("opening yearly template: %w", err)
	}
	defer yearlyTemplate.Close()

	const sheet = "Yearly"

	employee := data.Employee
	school := data.School
	endYear := startYear + 1

	// Helper for consistent cell-writing error handling.
	setValue := func(cell string, value any) error {
		if err := yearlyTemplate.SetCellValue(sheet, cell, value); err != nil {
			return fmt.Errorf("setting cell %s: %w", cell, err)
		}
		return nil
	}

	// --------------------------------------------------
	// Header
	// --------------------------------------------------

	header := fmt.Sprintf(
		"CLUSTER : %s   |   TAL : %s   |   DIST : Pune",
		data.Cluster.Name,
		school.Block,
	)

	if err := setValue("I3", header); err != nil {
		return "", err
	}

	if err := setValue(
		"A13",
		fmt.Sprintf("FINANCIAL YEAR : %d-%d", startYear, endYear),
	); err != nil {
		return "", err
	}

	// --------------------------------------------------
	// Employee information
	// --------------------------------------------------

	caser := cases.Title(language.English)
	employeeName := caser.String(employee.Name)

	employeeCells := []struct {
		cell  string
		value any
	}{
		{"C8", employeeName},
		{"C9", employee.Aadhaar},
		{"C10", employee.PAN},

		{"J8", school.Name},
		{"J9", employee.Designation},
		{"J10", employee.Mobile},
		{"J11", employee.Email},

		{"Q8", employee.BankName},
		{"Q9", employee.BranchName},
		{"Q10", employee.BankAccount},
		{"Q11", employee.BankIFSC},
	}

	for _, item := range employeeCells {
		if err := setValue(item.cell, item.value); err != nil {
			return "", err
		}
	}

	// --------------------------------------------------
	// Financial-year month order: April to March
	// --------------------------------------------------

	months := []struct {
		month int
		name  string
	}{
		{4, "Apr"},
		{5, "May"},
		{6, "Jun"},
		{7, "Jul"},
		{8, "Aug"},
		{9, "Sep"},
		{10, "Oct"},
		{11, "Nov"},
		{12, "Dec"},
		{1, "Jan"},
		{2, "Feb"},
		{3, "Mar"},
	}

	// Find a payslip for a specific month and calendar year.
	findPayslip := func(month, year int) *models.PayslipRecord {
		for i := range data.Payslips {
			payslip := &data.Payslips[i]

			if payslip.Month == month && payslip.Year == year {
				return payslip
			}
		}

		return nil
	}

	const firstMonthRow = 17

	// Accumulate annual totals by spreadsheet column.
	// Values stay numeric, preserving the template's cell formatting.
	annualTotals := make(map[string]float64)

	// --------------------------------------------------
	// Populate monthly rows and calculate annual totals
	// --------------------------------------------------

	for index, month := range months {
		row := firstMonthRow + index

		if err := setValue(
			fmt.Sprintf("A%d", row),
			fmt.Sprintf("%02d", index+1),
		); err != nil {
			return "", err
		}

		if err := setValue(
			fmt.Sprintf("B%d", row),
			month.name,
		); err != nil {
			return "", err
		}

		// January through March belong to the following year.
		year := startYear
		if month.month <= 3 {
			year = endYear
		}

		payslip := findPayslip(month.month, year)
		if payslip == nil {
			// Missing months remain blank and contribute zero
			// to the annual totals.
			continue
		}

		// GIS combines two database fields.
		gis := payslip.GISZP + payslip.GISScout

		// Preserve the existing monthly deduction calculation.
		totalDeduction := payslip.TotalGovtDeductions +
			payslip.NPSEmpContri +
			payslip.NPSEmprContri +
			payslip.NPSEmpContriArr +
			payslip.NPSEmprContriArr +
			payslip.NGRTotalDeduction

		values := []struct {
			column string
			value  float64
		}{
			{"C", payslip.BasicPay},
			{"D", payslip.DA},
			{"E", payslip.HRA},
			{"F", payslip.TA},
			{"G", payslip.DAArrears},
			{"H", payslip.NPSEmprAllow},
			{"I", payslip.TotalPay},

			{"J", payslip.GPF},
			{"K", payslip.PT},
			{"L", gis},
			{"M", payslip.RevenueStamp},

			{"N", payslip.NPSEmprContri},
			{"O", payslip.NPSEmpContri},
			{"P", payslip.NPSEmprContriArr},
			{"Q", payslip.NPSEmpContriArr},

			{"R", payslip.GroupAccidentalPolicy},
			{"S", payslip.IncomeTax},

			{"T", totalDeduction},
			{"U", payslip.EmployeeNetSalary},
		}

		for _, item := range values {
			cell := fmt.Sprintf("%s%d", item.column, row)

			if err := setValue(cell, item.value); err != nil {
				return "", err
			}

			// Calculate annual totals in Go instead of using
			// spreadsheet formulas.
			annualTotals[item.column] += item.value
		}
	}

	// --------------------------------------------------
	// Annual total row
	//
	// Monthly rows: 17-28
	// Total row:    29
	// --------------------------------------------------

	if err := setValue("A29", "13"); err != nil {
		return "", err
	}

	if err := setValue("B29", "Total"); err != nil {
		return "", err
	}

	totalColumns := []string{
		"C", "D", "E", "F", "G", "H", "I",
		"J", "K", "L", "M", "N", "O", "P",
		"Q", "R", "S", "T", "U",
	}

	for _, column := range totalColumns {
		cell := fmt.Sprintf("%s29", column)

		if err := setValue(cell, annualTotals[column]); err != nil {
			return "", err
		}
	}

	// --------------------------------------------------
	// Signatory names
	// --------------------------------------------------

	signatories := []struct {
		cell  string
		value string
	}{
		{"B35", employeeName},
		{"I35", school.Name},
		{"P35", data.Cluster.Name},
	}

	for _, item := range signatories {
		if err := setValue(item.cell, item.value); err != nil {
			return "", err
		}
	}

	// --------------------------------------------------
	// Output directory and filename
	// --------------------------------------------------

	financialYearDir := filepath.Join(
		outputDir,
		fmt.Sprintf("Financial_Year_%d_%d", startYear, endYear),
	)

	clusterDir := filepath.Join(
		financialYearDir,
		SanitizeFilename(data.Cluster.Name),
	)

	schoolDir := filepath.Join(
		clusterDir,
		SanitizeFilename(school.Name),
	)

	if err := os.MkdirAll(schoolDir, 0755); err != nil {
		return "", fmt.Errorf(
			"creating yearly school directory: %w",
			err,
		)
	}

	filename := fmt.Sprintf(
		"%s_%d_%d.xlsx",
		SanitizeFilename(employee.Name),
		startYear,
		endYear,
	)

	outputPath := filepath.Join(schoolDir, filename)

	// --------------------------------------------------
	// Save
	// --------------------------------------------------

	if err := yearlyTemplate.SaveAs(outputPath); err != nil {
		return "", fmt.Errorf("saving yearly payslip: %w", err)
	}

	return outputPath, nil
}
