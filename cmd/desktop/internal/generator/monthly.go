package generator

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/osutils"
	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/query"
	"github.com/Dimensionexpert/payslip/internal/concurrency"
	"github.com/Dimensionexpert/payslip/internal/database"
	genexcel "github.com/Dimensionexpert/payslip/internal/genExcel"
	genPDF "github.com/Dimensionexpert/payslip/internal/genPDF"
	"github.com/Dimensionexpert/payslip/internal/generator"
)

// Resolve the school's output directory from one existing payslip.
func mResolveSchoolScopeRoot(
	db *sql.DB,
	udise string,
	month int,
	year int,
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

	payslip, err := database.GetPayslip(
		db,
		employees[0].ShalarthID,
		month,
		year,
	)
	if err != nil {
		return "", fmt.Errorf(
			"[query] getting payslip for school %s: %w",
			udise,
			err,
		)
	}

	pdfPath := genexcel.MonthlyPayslipPDFPath(baseDir, payslip)

	return filepath.Dir(filepath.Dir(pdfPath)), nil
}

func mResolveClusterScopeRoot(
	db *sql.DB,
	cluster string,
	month int,
	year int,
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
				"[query] getting employees for school %s in cluster %s: %w",
				school.UDISECode,
				cluster,
				err,
			)
		}

		if len(employees) == 0 {
			continue
		}

		payslip, err := database.GetPayslip(
			db,
			employees[0].ShalarthID,
			month,
			year,
		)
		if err != nil {
			return "", fmt.Errorf(
				"[query] getting payslip for school %s in cluster %s: %w",
				school.UDISECode,
				cluster,
				err,
			)
		}

		pdfPath := genexcel.MonthlyPayslipPDFPath(
			baseDir,
			payslip,
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

// ResolveMonthlyScopeRoot resolves the output directory for a monthly generation scope.
func ResolveMonthlyScopeRoot(
	db *sql.DB,
	baseDir string,
	scope string,
	cluster string,
	udise string,
	month int,
	year int,
) (string, error) {
	if baseDir == "" {
		return "", fmt.Errorf("[path] base directory not configured")
	}

	root := filepath.Join(
		baseDir,
		fmt.Sprintf("%s_%d", time.Month(month), year),
	)

	switch scope {
	case "all":
		return root, nil

	case "school":
		return mResolveSchoolScopeRoot(
			db,
			udise,
			month,
			year,
			baseDir,
		)

	case "cluster":
		return mResolveClusterScopeRoot(
			db,
			cluster,
			month,
			year,
			baseDir,
		)

	default:
		return "", fmt.Errorf(
			"invalid generation scope: %s",
			scope,
		)
	}
}

// Generate Monthly payslip.
// intentionally not derived from the cache, it is to quickly apply arbitary changes requested.

func GenerateMonthlyPayslip(
	db *sql.DB,
	shalarthID string,
	month int,
	year int,
	outputDir string,
	monthlyTemplate string,
) (string, error) {
	if outputDir == "" {
		return "", fmt.Errorf("[path] output directory is not configured")
	}

	payslip, err := database.GetPayslip(
		db,
		shalarthID,
		month,
		year,
	)
	if err != nil {
		return "", fmt.Errorf(
			"[query] fetching monthly payslip for employee %s: %w",
			shalarthID,
			err,
		)
	}

	xlsxPath, err := genexcel.GenerateMonthlyPayslip(
		monthlyTemplate,
		outputDir,
		payslip,
	)
	if err != nil {
		return "", fmt.Errorf(
			"[excel] generating monthly payslip for employee %s: %w",
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
			"[pdf] converting monthly payslip for employee %s: %w",
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

// Generating bulk payslips from cache, to save time on xlsx gen on every generate operatin
func GenerateMonthlyPayslips(
	db *sql.DB,
	month int,
	year int,
	cacheDir string,
	resourcesDir string,
	outputDir string,
) error {
	monthlyTemplate := filepath.Join(
		resourcesDir,
		"monthly_template.xlsx",
	)

	monthlyRoot, err := EnsureMonthlyXLSX(
		month,
		year,
		db,
		cacheDir,
		monthlyTemplate,
	)
	if err != nil {
		return fmt.Errorf(
			"[excel] ensuring monthly XLSX files: %w",
			err,
		)
	}

	if outputDir == "" {
		return fmt.Errorf(
			"[path] output directory is not configured",
		)
	}

	monthlyOutputRoot := filepath.Join(
		outputDir,
		fmt.Sprintf("%s_%d", time.Month(month), year),
	)

	conversionJobs, err := generator.CollectMissingPDFJobs(
		monthlyRoot,
		monthlyOutputRoot,
	)
	if err != nil {
		return fmt.Errorf(
			"[file] collecting monthly PDF jobs: %w",
			err,
		)
	}

	results := concurrency.RunPDFConversion(
		conversionJobs,
		8,
		func(result concurrency.ConversionResult) {
			if result.Err != nil {
				fmt.Printf(
					"MONTHLY PDF FAILED: %s: %v\n",
					result.Filepath,
					result.Err,
				)
			}
		},
	)

	success, failed := generator.CountConversionResults(results)

	fmt.Printf(
		"Monthly PDF conversion: %d succeeded, %d failed\n",
		success,
		failed,
	)

	if failed > 0 {
		return fmt.Errorf(
			"[pdf] monthly conversion failed for %d file(s)",
			failed,
		)
	}

	return nil
}

// Opening monthly Payslips inf default applicatin

func OpenMonthlyPayslip(
	db *sql.DB,
	shalarthID string,
	month int,
	year int,
	outputDir string,
) error {
	payslip, err := database.GetPayslip(
		db,
		shalarthID,
		month,
		year,
	)
	if err != nil {
		return fmt.Errorf(
			"[query] fetching monthly payslip: %w",
			err,
		)
	}

	if outputDir == "" {
		return fmt.Errorf(
			"[path] output directory is not configured",
		)
	}

	pdfPath := genexcel.MonthlyPayslipPDFPath(
		outputDir,
		payslip,
	)

	if _, err := os.Stat(pdfPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(
				"[file] payslip PDF not found: %s",
				pdfPath,
			)
		}

		return fmt.Errorf(
			"[file] checking payslip PDF: %w",
			err,
		)
	}

	if err := osutils.OpenDir(pdfPath); err != nil {
		return fmt.Errorf(
			"[file] opening monthly payslip PDF: %w",
			err,
		)
	}

	fmt.Println(pdfPath)
	return nil
}

func GenerateMonthlyPayslipsForScope(
	db *sql.DB,
	cacheDir string,
	resourcesDir string,
	outputDir string,
	scope string,
	cluster string,
	udise string,
	month int,
	year int,
) (int, error) {
	monthlyTemplate := filepath.Join(
		resourcesDir,
		"monthly_template.xlsx",
	)

	if _, err := EnsureMonthlyXLSX(
		month,
		year,
		db,
		cacheDir,
		monthlyTemplate,
	); err != nil {
		return 0, fmt.Errorf(
			"[excel] ensuring monthly XLSX files: %w",
			err,
		)
	}

	root, err := ResolveMonthlyScopeRoot(
		db,
		cacheDir,
		scope,
		cluster,
		udise,
		month,
		year,
	)
	if err != nil {
		return 0, fmt.Errorf(
			"[path] resolving monthly cache scope: %w",
			err,
		)
	}

	outputRoot, err := ResolveMonthlyScopeRoot(
		db,
		outputDir,
		scope,
		cluster,
		udise,
		month,
		year,
	)
	if err != nil {
		return 0, fmt.Errorf(
			"[path] resolving monthly output scope: %w",
			err,
		)
	}

	jobs, err := generator.CollectMissingPDFJobs(
		root,
		outputRoot,
	)
	if err != nil {
		return 0, fmt.Errorf(
			"[file] collecting monthly PDF jobs: %w",
			err,
		)
	}

	if len(jobs) == 0 {
		return 0, nil
	}

	results := concurrency.RunPDFConversion(
		jobs,
		8,
		func(result concurrency.ConversionResult) {
			if result.Err != nil {
				fmt.Printf(
					"MONTHLY PDF FAILED: %s: %v\n",
					result.Filepath,
					result.Err,
				)
			}
		},
	)

	_, failed := generator.CountConversionResults(results)

	if failed > 0 {
		return 0, fmt.Errorf(
			"[pdf] monthly conversion failed for %d file(s)",
			failed,
		)
	}

	return len(jobs), nil
}
