package genexcel

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/Dimensionexpert/payslip/internal/models"
)

func sanitizeFilename(s string) string {
	return strings.Join(strings.Fields(s), "_")
}

func GenerateMonthlyPayslip(
	templatePath string,
	outputDir string,
	data models.PayslipExport,
) (string, error) {

	f, err := excelize.OpenFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("opening template: %w", err)
	}
	defer f.Close()

	const sheet = "Sheet1"

	school := data.School
	employee := data.Employee
	payslip := data.Payslip

	// --------------------------------------------------
	// Employee / School information
	// --------------------------------------------------

	f.SetCellValue(sheet, "C8", employee.Name)
	f.SetCellValue(sheet, "H8", employee.ID)

	f.SetCellValue(sheet, "C9", employee.Designation)
	f.SetCellValue(sheet, "H9", school.UDISECode)

	f.SetCellValue(sheet, "C10", school.Name)
	f.SetCellValue(sheet, "H10", employee.PAN)

	// Employee can be either GPF or DCPS.
	retirementNo := employee.GPFNo
	if retirementNo == "" {
		retirementNo = employee.DCPSNo
	}

	f.SetCellValue(sheet, "C11", retirementNo)
	f.SetCellValue(sheet, "C12", employee.ShalarthID)

	// --------------------------------------------------
	// Payroll period
	// --------------------------------------------------

	monthName := time.Month(payslip.Month).String()

	f.SetCellValue(
		sheet,
		"D14",
		fmt.Sprintf("%s %d", monthName, payslip.Year),
	)

	// --------------------------------------------------
	// Earnings
	// --------------------------------------------------

	f.SetCellValue(sheet, "D17", payslip.BasicPay)
	f.SetCellValue(sheet, "D18", payslip.DA)
	f.SetCellValue(sheet, "D19", payslip.HRA)
	f.SetCellValue(sheet, "D20", payslip.TA)
	f.SetCellValue(sheet, "D21", payslip.TAArrear)
	f.SetCellValue(sheet, "D22", payslip.DAArrears)
	f.SetCellValue(sheet, "D23", payslip.BasicArrears)
	f.SetCellValue(sheet, "D24", payslip.NPSEmprAllow)

	// --------------------------------------------------
	// Deductions
	// --------------------------------------------------

	f.SetCellValue(sheet, "H17", payslip.FA)
	f.SetCellValue(sheet, "H18", payslip.GPF)
	f.SetCellValue(sheet, "H19", payslip.GPFAdvance)
	f.SetCellValue(sheet, "H20", payslip.PT)

	// GIS is represented by two DB fields,
	// but the template has a single GIS field.
	gis := payslip.GISZP + payslip.GISScout
	f.SetCellValue(sheet, "H21", gis)

	f.SetCellValue(sheet, "H22", payslip.RevenueStamp)
	f.SetCellValue(sheet, "H23", payslip.NPSEmprContri)
	f.SetCellValue(sheet, "H24", payslip.NPSEmpContri)
	f.SetCellValue(sheet, "H25", payslip.IncomeTax)
	f.SetCellValue(sheet, "H26", payslip.NGRSocietyLoan)

	// HOME LOAN was removed from the template.
	// Therefore H27 is intentionally untouched.

	// --------------------------------------------------
	// Totals
	// --------------------------------------------------

	f.SetCellValue(sheet, "D28", payslip.TotalPay)

	totalDeduction := payslip.TotalGovtDeductions +
		payslip.NGRTotalDeduction

	f.SetCellValue(sheet, "H28", totalDeduction)

	// --------------------------------------------------
	// Net salary
	// --------------------------------------------------

	f.SetCellValue(
		sheet,
		"A31",
		fmt.Sprintf("₹ %.0f", payslip.EmployeeNetSalary),
	)

	f.SetCellValue(
		sheet,
		"D31",
		AmountInWords(payslip.EmployeeNetSalary),
	)

	// --------------------------------------------------
	// Footer
	// --------------------------------------------------

	// A33 ("PLACE") remains untouched because
	// f.SetCellValue(sheet, "A33", school.Name)

	f.SetCellValue(sheet, "A34", school.Name)
	f.SetCellValue(
		sheet,
		"E34",
		time.Now().Format("02-Jan-06"),
	)

	// --------------------------------------------------
	// Output
	// --------------------------------------------------

	periodDir := filepath.Join(
		outputDir,
		fmt.Sprintf("%s_%d", monthName, payslip.Year),
	)

	clusterDir := filepath.Join(
		periodDir,
		sanitizeFilename(data.Cluster.Name),
	)

	schoolDir := filepath.Join(
		clusterDir,
		sanitizeFilename(school.Name),
	)

	if err := os.MkdirAll(schoolDir, 0755); err != nil {
		return "", fmt.Errorf("creating school directory: %w", err)
	}

	filename := fmt.Sprintf(
		"%s.xlsx",
		sanitizeFilename(employee.Name),
	)

	outputPath := filepath.Join(schoolDir, filename)

	if err := f.SaveAs(outputPath); err != nil {
		return "", fmt.Errorf("saving output: %w", err)
	}

	return outputPath, nil
}

func GenerateMonthlyPayslips(
	templatePath string,
	outputDir string,
	payslips []models.PayslipExport,
) error {
	for i, payslip := range payslips {
		_, err := GenerateMonthlyPayslip(
			templatePath,
			outputDir,
			payslip,
		)
		if err != nil {
			return fmt.Errorf(
				"generating payslip %d/%d for %s: %w",
				i+1,
				len(payslips),
				payslip.Employee.ShalarthID,
				err,
			)
		}
	}

	return nil
}
