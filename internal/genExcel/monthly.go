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

	// Helper for consistent cell-writing error handling.
	setValue := func(cell string, value any) error {
		if err := f.SetCellValue(sheet, cell, value); err != nil {
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

	if err := setValue("A3", header); err != nil {
		return "", err
	}

	// --------------------------------------------------
	// Employee information
	// --------------------------------------------------

	caser := cases.Title(language.English)
	employeeName := caser.String(employee.Name)

	// Employee can have either GPF or DCPS.
	retirementNo := employee.GPFNo
	if retirementNo == "" {
		retirementNo = employee.DCPSNo
	}

	employeeCells := []struct {
		cell  string
		value any
	}{
		{"C8", employeeName},
		{"C9", school.Name},
		{"C10", employee.Designation},
		{"C11", retirementNo},
		{
			"C12",
			fmt.Sprintf(
				"%s %d",
				time.Month(payslip.Month).String(),
				payslip.Year,
			),
		},

		{"F10", employee.ShalarthID},
		{"F11", school.UDISECode},
		{"F12", employee.PAN},
	}

	for _, item := range employeeCells {
		if err := setValue(item.cell, item.value); err != nil {
			return "", err
		}
	}

	// --------------------------------------------------
	// Earnings
	//
	// Write numeric values, not formatted currency strings.
	// The Excel template controls the INR display format.
	// --------------------------------------------------

	earnings := []struct {
		cell  string
		value float64
	}{
		{"C15", payslip.BasicPay},
		{"C16", payslip.DA},
		{"C17", payslip.HRA},
		{"C18", payslip.TA},
		{"C19", payslip.TAArrear},
		{"C20", payslip.DAArrears},
		{"C21", payslip.BasicArrears},
		{"C22", payslip.NPSEmprAllow},
		{"C23", 0},
		{"C26", payslip.TotalPay},
	}

	for _, item := range earnings {
		if err := setValue(item.cell, item.value); err != nil {
			return "", err
		}
	}

	// --------------------------------------------------
	// Deductions
	// --------------------------------------------------

	// GIS is represented by two database fields,
	// but the template has a single GIS field.
	gis := payslip.GISZP + payslip.GISScout

	deductions := []struct {
		cell  string
		value float64
	}{
		{"F15", payslip.GPF},
		{"F16", payslip.GPFAdvance},
		{"F17", payslip.PT},
		{"F18", gis},
		{"F19", payslip.RevenueStamp},
		{"F20", payslip.NPSEmprContri},
		{"F21", payslip.NPSEmpContri},
		{"F22", payslip.NPSEmprContriArr},
		{"F23", payslip.NPSEmpContriArr},
		{"F24", payslip.IncomeTax},
		{"F25", payslip.NGRSocietyLoan},
	}

	for _, item := range deductions {
		if err := setValue(item.cell, item.value); err != nil {
			return "", err
		}
	}

	// --------------------------------------------------
	// Total deductions
	// --------------------------------------------------

	// Preserve the existing deduction calculation.
	totalDeduction := payslip.TotalGovtDeductions +
		payslip.NPSEmpContri +
		payslip.NPSEmprContri +
		payslip.NPSEmpContriArr +
		payslip.NPSEmprContriArr +
		payslip.NGRTotalDeduction

	if err := setValue("F26", totalDeduction); err != nil {
		return "", err
	}

	// --------------------------------------------------
	// Net salary
	// --------------------------------------------------

	if err := setValue("A29", payslip.EmployeeNetSalary); err != nil {
		return "", err
	}

	if err := setValue(
		"E29",
		AmountInWords(payslip.EmployeeNetSalary),
	); err != nil {
		return "", err
	}

	// --------------------------------------------------
	// Signature / footer information
	// --------------------------------------------------

	footerCells := []struct {
		cell  string
		value any
	}{
		{"E36", data.Cluster.Name},
		{"A32", school.Name},
		{"B36", school.Name},
		{"F32", time.Now().Format("02-Jan-06")},
	}

	for _, item := range footerCells {
		if err := setValue(item.cell, item.value); err != nil {
			return "", err
		}
	}

	// --------------------------------------------------
	// Output directory and filename
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
		return "", fmt.Errorf(
			"creating school directory: %w",
			err,
		)
	}

	filename := fmt.Sprintf(
		"%s.xlsx",
		SanitizeFilename(employee.Name),
	)

	outputPath := filepath.Join(schoolDir, filename)

	// --------------------------------------------------
	// Save
	// --------------------------------------------------

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
