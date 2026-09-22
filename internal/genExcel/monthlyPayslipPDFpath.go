package genexcel

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/Dimensionexpert/payslip/internal/models"
)

func MonthlyPayslipPDFPath(
	outputDir string,
	data models.PayslipExport,
) string {
	monthName := time.Month(data.Payslip.Month)

	periodDir := filepath.Join(
		outputDir,
		fmt.Sprintf("%s_%d", monthName, data.Payslip.Year),
	)

	clusterDir := filepath.Join(
		periodDir,
		SanitizeFilename(data.Cluster.Name),
	)

	schoolDir := filepath.Join(
		clusterDir,
		SanitizeFilename(data.School.Name),
	)

	return filepath.Join(
		schoolDir,
		"PDF",
		fmt.Sprintf(
			"%s.pdf",
			SanitizeFilename(data.Employee.Name),
		),
	)
}
