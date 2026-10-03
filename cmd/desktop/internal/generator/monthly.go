package generator

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/query"
	"github.com/Dimensionexpert/payslip/internal/database"
	genexcel "github.com/Dimensionexpert/payslip/internal/genExcel"
)

// Resolve the school's output directory from one existing payslip.
func ResolveSchoolScopeRoot(
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

func ResolveClusterScopeRoot(
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
		return ResolveSchoolScopeRoot(
			db,
			udise,
			month,
			year,
			baseDir,
		)

	case "cluster":
		return ResolveClusterScopeRoot(
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
