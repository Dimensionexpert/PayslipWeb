package genexcel

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/Dimensionexpert/payslip/internal/models"
)

func SanitizeFilename(s string) string {
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

	const sheet = "Monthly"

	school := data.School
	employee := data.Employee
	payslip := data.Payslip

	// --------------------------------------------------
	// Employee information
	// --------------------------------------------------

	f.SetCellValue(
		sheet,
		"A3",
		fmt.Sprintf(
			"CLUSTER : %s   |   TAL : %s   |   DIST : %s",
			data.Cluster.Name,
			data.School.Block,
			"Pune",
		),
	)

	// --------------------------------------------------
	// Employee information
	// --------------------------------------------------
	caser := cases.Title(language.English)
	employeeName := caser.String(employee.Name)
	f.SetCellValue(sheet, "C8", employeeName)
	f.SetCellValue(sheet, "C9", school.Name)
	f.SetCellValue(sheet, "C10", employee.Designation)
	f.SetCellValue(sheet, "C11", employee.GPFNo)
	f.SetCellValue(sheet, "C12", fmt.Sprintf(
		"%s %d",
		time.Month(payslip.Month).String(),
		payslip.Year,
	))

	f.SetCellValue(sheet, "F10", employee.ShalarthID)
	f.SetCellValue(sheet, "F11", school.UDISECode)
	f.SetCellValue(sheet, "F12", employee.PAN)

	// Employee can have either GPF or DCPS.
	retirementNo := employee.GPFNo
	if retirementNo == "" {
		retirementNo = employee.DCPSNo
	}

	f.SetCellValue(sheet, "C11", retirementNo)

	// --------------------------------------------------
	// Earnings
	// --------------------------------------------------

	f.SetCellValue(sheet, "C15", formatINR(payslip.BasicPay))
	f.SetCellValue(sheet, "C16", formatINR(payslip.DA))
	f.SetCellValue(sheet, "C17", formatINR(payslip.HRA))
	f.SetCellValue(sheet, "C18", formatINR(payslip.TA))
	f.SetCellValue(sheet, "C19", formatINR(payslip.TAArrear))
	f.SetCellValue(sheet, "C20", formatINR(payslip.DAArrears))
	f.SetCellValue(sheet, "C21", formatINR(payslip.BasicArrears))
	f.SetCellValue(sheet, "C22", formatINR(payslip.NPSEmprAllow))
	f.SetCellValue(sheet, "C23", formatINR(0))

	// --------------------------------------------------
	// Deductions
	// --------------------------------------------------

	f.SetCellValue(sheet, "F15", formatINR(payslip.GPF))
	f.SetCellValue(sheet, "F16", formatINR(payslip.GPFAdvance))
	f.SetCellValue(sheet, "F17", formatINR(payslip.PT))

	// GIS is represented by two DB fields,
	// but the template has a single GIS field.
	gis := payslip.GISZP + payslip.GISScout
	f.SetCellValue(sheet, "F18", formatINR(gis))

	f.SetCellValue(sheet, "F19", formatINR(payslip.RevenueStamp))
	f.SetCellValue(sheet, "F20", formatINR(payslip.NPSEmprContri))
	f.SetCellValue(sheet, "F21", formatINR(payslip.NPSEmpContri))
	f.SetCellValue(sheet, "F22", formatINR(payslip.NPSEmprContriArr))
	f.SetCellValue(sheet, "F23", formatINR(payslip.NPSEmpContriArr))
	f.SetCellValue(sheet, "F24", formatINR(payslip.IncomeTax))
	f.SetCellValue(sheet, "F25", formatINR(payslip.NGRSocietyLoan))

	// --------------------------------------------------
	// Totals
	// --------------------------------------------------

	f.SetCellValue(sheet, "C26", formatINR(payslip.TotalPay))

	totalDeduction := payslip.TotalGovtDeductions +
		payslip.NPSEmpContri +
		payslip.NPSEmprContri +
		payslip.NPSEmpContriArr +
		payslip.NPSEmprContriArr +
		payslip.NGRTotalDeduction

	f.SetCellValue(sheet, "F26", formatINR(totalDeduction))

	// --------------------------------------------------
	// Net salary
	// --------------------------------------------------

	f.SetCellValue(sheet, "A29", formatINR(payslip.EmployeeNetSalary))

	f.SetCellValue(
		sheet,
		"E29",
		AmountInWords(payslip.EmployeeNetSalary),
	)

	// --------------------------------------------------
	// Signature / footer information
	// --------------------------------------------------

	f.SetCellValue(sheet, "E36", data.Cluster.Name)
	f.SetCellValue(sheet, "A32", school.Name)
	f.SetCellValue(sheet, "B36", school.Name)

	f.SetCellValue(
		sheet,
		"F32",
		time.Now().Format("02-Jan-06"),
	)

	// --------------------------------------------------
	// Output
	// --------------------------------------------------

	monthName := time.Month(payslip.Month).String()

	periodDir := filepath.Join(
		outputDir,
		fmt.Sprintf("%s_%d", monthName, payslip.Year),
	)

	clusterDir := filepath.Join(
		periodDir,
		SanitizeFilename(data.Cluster.Name),
	)

	schoolDir := filepath.Join(
		clusterDir,
		SanitizeFilename(school.Name),
	)

	if err := os.MkdirAll(schoolDir, 0755); err != nil {
		return "", fmt.Errorf("creating school directory: %w", err)
	}

	filename := fmt.Sprintf(
		"%s.xlsx",
		SanitizeFilename(employee.Name),
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
