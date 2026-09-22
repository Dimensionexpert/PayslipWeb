package genexcel

import (
	"fmt"
	"path/filepath"

	"github.com/Dimensionexpert/payslip/internal/models"
)

func YearlyPayslipPDFPath(
	outputDir string,
	data models.PayslipExportYear,
	startYear int,
) string {

	endYear := startYear + 1

	financialYearDir := filepath.Join(
		outputDir,
		fmt.Sprintf(
			"Financial_Year_%d_%d",
			startYear,
			endYear,
		),
	)

	clusterDir := filepath.Join(
		financialYearDir,
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
			"%s_%d_%d.pdf",
			SanitizeFilename(data.Employee.Name),
			startYear,
			endYear,
		),
	)
}
