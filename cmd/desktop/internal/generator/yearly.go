package generator

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/query"
	"github.com/Dimensionexpert/payslip/internal/database"
	genexcel "github.com/Dimensionexpert/payslip/internal/genExcel"
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
