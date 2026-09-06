package generator

import (
	"fmt"
	"os"
	"time"

	genexcel "github.com/Dimensionexpert/payslip/internal/genExcel"
	"github.com/Dimensionexpert/payslip/internal/models"
)

func GenerateMonthlyExcel(
	templatePath string,
	outputDir string,
	exportPayslips []models.PayslipExport,
) error {
	start := time.Now()

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory failed: %w", err)
	}

	if err := genexcel.GenerateMonthlyPayslips(
		templatePath,
		outputDir,
		exportPayslips,
	); err != nil {
		return fmt.Errorf("generating monthly Excel files failed: %w", err)
	}

	fmt.Printf(
		"Time taken to generate Excel files: %v\n",
		time.Since(start),
	)

	fmt.Printf(
		"Generated %d monthly payslips in %s\n",
		len(exportPayslips),
		outputDir,
	)

	return nil
}
