package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Dimensionexpert/payslip/internal/importer"
)

func ImportPayroll(
	truthPath string,
	month int,
	year int,
	sourceDir string,
	resourcesDir string,
	dbPath string,
	cacheDir string,
) (importer.ImportReport, error) {
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		return importer.ImportReport{}, fmt.Errorf(
			"[file] creating source directory: %w",
			err,
		)
	}

	destination := filepath.Join(
		sourceDir,
		filepath.Base(truthPath),
	)

	sourceAbs, err := filepath.Abs(truthPath)
	if err != nil {
		return importer.ImportReport{}, fmt.Errorf(
			"[path] resolving payroll file path: %w",
			err,
		)
	}

	destinationAbs, err := filepath.Abs(destination)
	if err != nil {
		return importer.ImportReport{}, fmt.Errorf(
			"[path] resolving source path: %w",
			err,
		)
	}

	if sourceAbs != destinationAbs {
		src, err := os.Open(truthPath)
		if err != nil {
			return importer.ImportReport{}, fmt.Errorf(
				"[file] opening selected payroll file: %w",
				err,
			)
		}
		defer src.Close()

		dst, err := os.Create(destination)
		if err != nil {
			return importer.ImportReport{}, fmt.Errorf(
				"[file] creating source payroll file: %w",
				err,
			)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return importer.ImportReport{}, fmt.Errorf(
				"[file] copying payroll file: %w",
				err,
			)
		}
	}

	clusterPath := filepath.Join(
		resourcesDir,
		"clusters.xlsx",
	)

	report, _, _, err := importer.ImportPayroll(
		clusterPath,
		destination,
		dbPath,
		month,
		year,
	)
	if err != nil {
		return report, fmt.Errorf(
			"[import] importing payroll: %w",
			err,
		)
	}

	monthlyCacheRoot := filepath.Join(
		cacheDir,
		fmt.Sprintf("%s_%d", time.Month(month), year),
	)

	if err := os.RemoveAll(monthlyCacheRoot); err != nil {
		return report, fmt.Errorf(
			"[file] invalidating monthly XLSX cache: %w",
			err,
		)
	}

	return report, nil
}
