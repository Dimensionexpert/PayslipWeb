package generator

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Dimensionexpert/payslip/internal/database"
	genexcel "github.com/Dimensionexpert/payslip/internal/genExcel"
	"github.com/Dimensionexpert/payslip/internal/generator"
)

func EnsureMonthlyXLSX(
	month int,
	year int,
	db *sql.DB,
	cacheDir string,
	monthlyTemplate string,
) (string, error) {
	payslips, err := database.GetPayslips(db, month, year)
	if err != nil {
		return "", fmt.Errorf(
			"[query] fetching payslips for %d/%d: %w",
			month,
			year,
			err,
		)
	}

	if len(payslips) == 0 {
		return "", fmt.Errorf(
			"[query] no payslips found for %d/%d",
			month,
			year,
		)
	}

	monthlyRoot := filepath.Join(
		cacheDir,
		fmt.Sprintf("%s_%d", time.Month(month), year),
	)

	readyMarker := filepath.Join(
		monthlyRoot,
		".xlsx-ready",
	)

	if _, err := os.Stat(readyMarker); err == nil {
		return monthlyRoot, nil
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf(
			"[file] checking monthly XLSX cache: %w",
			err,
		)
	}

	if err := genexcel.GenerateMonthlyPayslips(
		monthlyTemplate,
		cacheDir,
		payslips,
	); err != nil {
		return "", fmt.Errorf(
			"[excel] generating monthly Excel files: %w",
			err,
		)
	}

	if err := os.WriteFile(
		readyMarker,
		[]byte{},
		0644,
	); err != nil {
		return "", fmt.Errorf(
			"[file] writing monthly XLSX ready marker: %w",
			err,
		)
	}

	return monthlyRoot, nil
}

func EnsureYearlyXLSX(
	financialYearStart int,
	db *sql.DB,
	cacheDir string,
	yearlyTemplate string,
) (string, error) {
	yearlyRoot := filepath.Join(
		cacheDir,
		fmt.Sprintf(
			"Financial_Year_%d_%d",
			financialYearStart,
			financialYearStart+1,
		),
	)

	readyMarker := filepath.Join(
		yearlyRoot,
		".xlsx-ready",
	)

	if _, err := os.Stat(readyMarker); err == nil {
		return yearlyRoot, nil
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf(
			"[file] checking yearly XLSX cache: %w",
			err,
		)
	}

	if err := generator.ExportYearlyEmployeePayslips(
		db,
		yearlyTemplate,
		cacheDir,
		financialYearStart,
	); err != nil {
		return "", fmt.Errorf(
			"[excel] generating yearly Excel files: %w",
			err,
		)
	}

	if err := os.WriteFile(
		readyMarker,
		[]byte{},
		0644,
	); err != nil {
		return "", fmt.Errorf(
			"[file] writing yearly XLSX ready marker: %w",
			err,
		)
	}

	return yearlyRoot, nil
}
