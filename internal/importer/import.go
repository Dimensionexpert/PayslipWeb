package importer

import (
	"fmt"
	"path/filepath"

	"github.com/Dimensionexpert/payslip/internal/database"
	"github.com/Dimensionexpert/payslip/internal/excel"
)

func ImportPayroll(
	clusterPath string,
	truthPath string,
	dbPath string,
) (ImportReport, int, int, error) {

	// Read cluster workbook
	clusterFile, err := excel.Open(clusterPath)
	if err != nil {
		return ImportReport{}, 0, 0,
			fmt.Errorf("opening cluster file: %w", err)
	}
	defer clusterFile.Close()

	clusterRows, err := excel.GetRows(clusterFile, "Sheet1")
	if err != nil {
		return ImportReport{}, 0, 0,
			fmt.Errorf("reading cluster file: %w", err)
	}

	clusterMap, duplicateUDISECount := GetClusterMap(clusterRows)

	// Read payroll workbook
	truthFile, err := excel.Open(truthPath)
	if err != nil {
		return ImportReport{}, 0, 0,
			fmt.Errorf("opening payroll file: %w", err)
	}
	defer truthFile.Close()

	truthRows, err := excel.GetRows(truthFile, "abstract")
	if err != nil {
		return ImportReport{}, 0, 0,
			fmt.Errorf("reading payroll file: %w", err)
	}

	// Detect payroll period
	month, year, err := ParsePeriod(filepath.Base(truthPath))
	if err != nil {
		return ImportReport{}, 0, 0,
			fmt.Errorf("detecting payroll period: %w", err)
	}

	// Convert rows into Go structs
	employees, schools, payslips, report, err := ImportTruth(
		truthRows,
		clusterMap,
		month,
		year,
	)
	if err != nil {
		return report, month, year,
			fmt.Errorf("importing payroll rows: %w", err)
	}

	report.DuplicateUDISECount = duplicateUDISECount

	// Open database
	db, err := database.Open(dbPath)
	if err != nil {
		return report, month, year,
			fmt.Errorf("opening database: %w", err)
	}
	defer db.Close()

	// Insert master data
	if err := database.InsertClusters(db, clusterMap); err != nil {
		return report, month, year,
			fmt.Errorf("inserting clusters: %w", err)
	}

	if err := database.InsertSchools(db, schools); err != nil {
		return report, month, year,
			fmt.Errorf("inserting schools: %w", err)
	}

	if err := database.InsertEmployees(db, employees); err != nil {
		return report, month, year,
			fmt.Errorf("inserting employees: %w", err)
	}

	// Insert payslip records
	if err := database.InsertPayslipRecords(db, payslips); err != nil {
		return report, month, year,
			fmt.Errorf("inserting payslip records: %w", err)
	}

	return report, month, year, nil
}
