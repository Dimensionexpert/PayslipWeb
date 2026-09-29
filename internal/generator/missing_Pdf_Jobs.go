package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Dimensionexpert/payslip/internal/concurrency"
)

func CollectMissingPDFJobs(
	xlsxRoot string,
	outputRoot string,
) ([]concurrency.ConversionJob, error) {
	var jobs []concurrency.ConversionJob

	err := filepath.WalkDir(xlsxRoot, func(
		path string,
		entry os.DirEntry,
		walkErr error,
	) error {
		if walkErr != nil {
			return walkErr
		}

		if entry.IsDir() {
			return nil
		}

		if !strings.EqualFold(filepath.Ext(path), ".xlsx") {
			return nil
		}

		// Find the XLSX path relative to the cache root.
		relativePath, err := filepath.Rel(xlsxRoot, path)
		if err != nil {
			return fmt.Errorf(
				"getting relative XLSX path: %w",
				err,
			)
		}

		// Build the corresponding output directory.
		outputDir := filepath.Join(
			outputRoot,
			filepath.Dir(relativePath),
			"PDF",
		)

		pdfPath := filepath.Join(
			outputDir,
			strings.TrimSuffix(
				entry.Name(),
				filepath.Ext(entry.Name()),
			)+".pdf",
		)

		if _, err := os.Stat(pdfPath); err == nil {
			return nil
		} else if !os.IsNotExist(err) {
			return err
		}

		jobs = append(jobs, concurrency.ConversionJob{
			Filepath: path,
			OutDir:   outputDir,
		})

		return nil
	})

	return jobs, err
}
