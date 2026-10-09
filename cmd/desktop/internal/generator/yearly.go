package generator

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/query"
	"github.com/Dimensionexpert/payslip/internal/concurrency"
	"github.com/Dimensionexpert/payslip/internal/database"
	genexcel "github.com/Dimensionexpert/payslip/internal/genExcel"
	genPDF "github.com/Dimensionexpert/payslip/internal/genPDF"
	"github.com/Dimensionexpert/payslip/internal/generator"
)

func yResolveSchoolScopeRoot(
	db *sql.DB,
	udise string,
	financialYearStart int,
	baseDir string,
) (string, error) {
	employees, err := query.GetEmployeesBySchool(db, udise)
	if err != nil {
		return "", fmt.Errorf(
			"[query] getting employees for school %s: %w",
			udise,
			err,
		)
	}

	if len(employees) == 0 {
		return "", fmt.Errorf(
			"[query] no employees found for school %s",
			udise,
		)
	}

	payslip, err := database.GetYearlyPayslip(
		db,
		employees[0].ShalarthID,
		financialYearStart,
	)
	if err != nil {
		return "", fmt.Errorf(
			"[query] getting payslip for school %s: %w",
			udise,
			err,
		)
	}
	pdfPath := genexcel.YearlyPayslipPDFPath(
		baseDir,
		payslip,
		financialYearStart,
	)

	return filepath.Dir(filepath.Dir(pdfPath)), nil
}

func yResolveClusterScopeRoot(
	db *sql.DB,
	cluster string,
	financialYearStart int,
	baseDir string,
) (string, error) {
	schools, err := query.GetSchoolsByCluster(db, cluster)
	if err != nil {
		return "", fmt.Errorf(
			"[query] getting schools for cluster %s: %w",
			cluster,
			err,
		)
	}

	if len(schools) == 0 {
		return "", fmt.Errorf(
			"[query] no schools found in cluster %s",
			cluster,
		)
	}

	for _, school := range schools {
		employees, err := query.GetEmployeesBySchool(
			db,
			school.UDISECode,
		)
		if err != nil {
			return "", fmt.Errorf(
				"[query] getting payslip for school %s in cluster %s: %w",
				school.UDISECode,
				cluster,
				err,
			)
		}

		if len(employees) == 0 {
			continue
		}

		payslip, err := database.GetYearlyPayslip(
			db,
			employees[0].ShalarthID,
			financialYearStart,
		)
		if err != nil {
			return "", fmt.Errorf(
				"[query] getting payslip for school %s: %w",
				school.UDISECode,
				err,
			)
		}

		pdfPath := genexcel.YearlyPayslipPDFPath(
			baseDir,
			payslip,
			financialYearStart,
		)

		return filepath.Dir(
			filepath.Dir(
				filepath.Dir(pdfPath),
			),
		), nil
	}

	return "", fmt.Errorf(
		"[query] no employees found in cluster %s",
		cluster,
	)
}

func YearlyPayslipRoot(
	db *sql.DB,
	cluster string,
	scope string,
	udise string,
	financialYearStart int,
	baseDir string,
) (string, error) {
	if baseDir == "" {
		return "", fmt.Errorf("[path] base directory not configured")
	}

	root := filepath.Join(
		baseDir,
		fmt.Sprintf(
			"Financial_Year_%d_%d",
			financialYearStart,
			financialYearStart+1,
		),
	)

	switch scope {
	case "all":
		return root, nil

	case "school":
		return yResolveSchoolScopeRoot(
			db,
			udise,
			financialYearStart,
			baseDir,
		)

	case "cluster":
		return yResolveClusterScopeRoot(
			db,
			cluster,
			financialYearStart,
			baseDir,
		)

	default:
		return "", fmt.Errorf(
			"invalid generation scope: %s",
			scope,
		)
	}
}

// Generate Yearly payslip.
// intentionally not derived from the cache, it is to quickly apply arbitary changes requested.

func GenerateYearlyPayslip(
	db *sql.DB,
	shalarthID string,
	financialYearStart int,
	outputDir string,
	yearlyTemplate string,
) (string, error) {
	if outputDir == "" {
		return "", fmt.Errorf(
			"[path] output directory is not configured",
		)
	}

	yearly, err := database.GetYearlyPayslip(
		db,
		shalarthID,
		financialYearStart,
	)
	if err != nil {
		return "", fmt.Errorf(
			"[query] fetching yearly payslip for employee %s: %w",
			shalarthID,
			err,
		)
	}

	xlsxPath, err := genexcel.GenerateYearlyPayslip(
		yearlyTemplate,
		outputDir,
		yearly,
		financialYearStart,
	)
	if err != nil {
		return "", fmt.Errorf(
			"[excel] generating yearly payslip for employee %s: %w",
			shalarthID,
			err,
		)
	}

	pdfDir := filepath.Join(
		filepath.Dir(xlsxPath),
		"PDF",
	)

	if err := genPDF.ConvertToPDF(
		xlsxPath,
		pdfDir,
		0,
	); err != nil {
		return "", fmt.Errorf(
			"[pdf] converting yearly payslip for employee %s: %w",
			shalarthID,
			err,
		)
	}

	pdfFilename := strings.TrimSuffix(
		filepath.Base(xlsxPath),
		filepath.Ext(xlsxPath),
	) + ".pdf"

	pdfPath := filepath.Join(
		pdfDir,
		pdfFilename,
	)

	return pdfPath, nil
}

func GenerateYearlyPayslips(
	db *sql.DB,
	financialYearStart int,
	cacheDir string,
	resourcesDir string,
	outputDir string,
) error {
	yearlyTemplate := filepath.Join(
		resourcesDir,
		"yearly_template.xlsx",
	)

	yearlyRoot, err := EnsureYearlyXLSX(
		financialYearStart,
		db,
		cacheDir,
		yearlyTemplate,
	)
	if err != nil {
		return fmt.Errorf(
			"[excel] ensuring yearly XLSX files: %w",
			err,
		)
	}

	if outputDir == "" {
		return fmt.Errorf(
			"[path] output directory is not configured",
		)
	}

	yearlyOutputRoot := filepath.Join(
		outputDir,
		fmt.Sprintf(
			"Financial_Year_%d_%d",
			financialYearStart,
			financialYearStart+1,
		),
	)

	conversionJobs, err := generator.CollectMissingPDFJobs(
		yearlyRoot,
		yearlyOutputRoot,
	)
	if err != nil {
		return fmt.Errorf(
			"[file] collecting yearly PDF jobs: %w",
			err,
		)
	}

	results := concurrency.RunPDFConversion(
		conversionJobs,
		8,
		func(result concurrency.ConversionResult) {
			if result.Err != nil {
				fmt.Printf(
					"YEARLY PDF FAILED: %s: %v\n",
					result.Filepath,
					result.Err,
				)
			}
		},
	)

	success, failed := generator.CountConversionResults(results)

	fmt.Printf(
		"Yearly PDF conversion: %d succeeded, %d failed\n",
		success,
		failed,
	)

	if failed > 0 {
		return fmt.Errorf(
			"[pdf] yearly conversion failed for %d file(s)",
			failed,
		)
	}

	return nil
}
