package generator

import (
	"fmt"
	"os"

	genexcel "github.com/Dimensionexpert/payslip/internal/genExcel"
	"github.com/Dimensionexpert/payslip/internal/models"
)

func GenerateMonthlyExcel(
	templatePath string,
	outputDir string,
	exportPayslips []models.PayslipExport,
) error {

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

	return nil
}
