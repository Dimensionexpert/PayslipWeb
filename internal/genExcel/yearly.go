package genexcel

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"

	"github.com/Dimensionexpert/payslip/internal/models"
)

func GenerateYearlyPayslip(
	templatePath string,
	outputDir string,
	data models.PayslipExportYear,
	startYear int,
) (string, error) {
	f, err := excelize.OpenFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("opening yearly template: %w", err)
	}
	defer f.Close()

	const sheet = "Sheet1"

	employee := data.Employee
	school := data.School

	// --------------------------------------------------
	// Financial year
	//
	// Example:
	//     April 2026 to March 2027
	// --------------------------------------------------

	endYear := startYear + 1

	f.SetCellValue(
		sheet,
		"A6",
		fmt.Sprintf("Financial Year - %d - %d", startYear, endYear),
	)

	// --------------------------------------------------
	// Employee information
	// --------------------------------------------------

	f.SetCellValue(sheet, "B1", employee.Name)
	f.SetCellValue(sheet, "B2", "") // DOB is not available yet
	f.SetCellValue(sheet, "B3", data.Cluster.Name)
	f.SetCellValue(sheet, "B4", school.Name)

	f.SetCellValue(sheet, "G2", employee.Mobile)
	f.SetCellValue(sheet, "G3", school.Block) // Taluka
	f.SetCellValue(sheet, "G4", employee.Aadhaar)

	f.SetCellValue(sheet, "K2", employee.Email)
	f.SetCellValue(sheet, "K3", "Pune") // Fixed district
	f.SetCellValue(sheet, "K4", employee.PAN)

	f.SetCellValue(sheet, "P1", employee.BankName)
	f.SetCellValue(sheet, "P2", employee.BranchName)
	f.SetCellValue(sheet, "P3", employee.BankAccount)
	f.SetCellValue(sheet, "P4", employee.BankIFSC)

	// --------------------------------------------------
	// Financial-year month order
	//
	// April to March
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

	// --------------------------------------------------
	// Find the payslip for a specific month
	// --------------------------------------------------

	findPayslip := func(month int) *models.PayslipRecord {
		for i := range data.Payslips {
			payslip := &data.Payslips[i]

			if payslip.Month == month {
				return payslip
			}
		}

		return nil
	}

	// --------------------------------------------------
	// Write monthly values
	//
	// Columns:
	//
	// C  Basic
	// D  DA
	// E  HRA
	// F  TA
	// G  DA Arrears
	// H  NPS Employer Allowance
	// I  Total Pay
	// J  GPF
	// K  PT
	// L  GIS ZP + Scout
	// M  Revenue
	// N  NPS Employer Contribution
	// O  NPS Employee Contribution
	// P  NPS Employer Contribution Arrears
	// Q  NPS Employee Contribution Arrears
	// R  Group Accidental Policy
	// S  Income Tax
	// --------------------------------------------------

	for index, month := range months {
		row := 10 + index

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("%02d", index+1))
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), month.name)

		payslip := findPayslip(month.month)
		if payslip == nil {
			continue
		}

		gis := payslip.GISZP + payslip.GISScout

		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), payslip.BasicPay)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), payslip.DA)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), payslip.HRA)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), payslip.TA)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), payslip.DAArrears)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), payslip.NPSEmprAllow)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), payslip.TotalPay)

		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), payslip.GPF)
		f.SetCellValue(sheet, fmt.Sprintf("K%d", row), payslip.PT)
		f.SetCellValue(sheet, fmt.Sprintf("L%d", row), gis)
		f.SetCellValue(sheet, fmt.Sprintf("M%d", row), payslip.RevenueStamp)

		f.SetCellValue(sheet, fmt.Sprintf("N%d", row), payslip.NPSEmprContri)
		f.SetCellValue(sheet, fmt.Sprintf("O%d", row), payslip.NPSEmpContri)
		f.SetCellValue(sheet, fmt.Sprintf("P%d", row), payslip.NPSEmprContriArr)
		f.SetCellValue(sheet, fmt.Sprintf("Q%d", row), payslip.NPSEmpContriArr)

		f.SetCellValue(
			sheet,
			fmt.Sprintf("R%d", row),
			payslip.GroupAccidentalPolicy,
		)

		f.SetCellValue(sheet, fmt.Sprintf("S%d", row), payslip.IncomeTax)
	}

	// --------------------------------------------------
	// Total row
	// --------------------------------------------------

	f.SetCellValue(sheet, "A22", "15")
	f.SetCellValue(sheet, "B22", "Total")

	totalColumns := []string{
		"C", "D", "E", "F", "G", "H", "I",
		"J", "K", "L", "M", "N", "O", "P",
		"Q", "R", "S",
	}

	for _, column := range totalColumns {
		formula := fmt.Sprintf("SUM(%s10:%s21)", column, column)

		if err := f.SetCellFormula(sheet, column+"22", formula); err != nil {
			return "", fmt.Errorf(
				"setting total formula for column %s: %w",
				column,
				err,
			)
		}
	}

	// --------------------------------------------------
	// Output
	// --------------------------------------------------

	financialYearDir := filepath.Join(
		outputDir,
		fmt.Sprintf("Financial_Year_%d_%d", startYear, endYear),
	)

	clusterDir := filepath.Join(
		financialYearDir,
		sanitizeFilename(data.Cluster.Name),
	)

	schoolDir := filepath.Join(
		clusterDir,
		sanitizeFilename(school.Name),
	)

	if err := os.MkdirAll(schoolDir, 0755); err != nil {
		return "", fmt.Errorf(
			"creating yearly school directory: %w",
			err,
		)
	}

	filename := fmt.Sprintf(
		"%s_%d_%d.xlsx",
		sanitizeFilename(employee.Name),
		startYear,
		endYear,
	)

	outputPath := filepath.Join(schoolDir, filename)

	if err := f.SaveAs(outputPath); err != nil {
		return "", fmt.Errorf("saving yearly payslip: %w", err)
	}

	return outputPath, nil
}
