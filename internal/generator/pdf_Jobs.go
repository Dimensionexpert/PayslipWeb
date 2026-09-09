package generator

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Dimensionexpert/payslip/internal/concurrency"
)

// internal/generator/pdf.go

func CollectPDFJobs(root string) ([]concurrency.ConversionJob, error) {
	var jobs []concurrency.ConversionJob

	err := filepath.WalkDir(root, func(
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

		jobs = append(jobs, concurrency.ConversionJob{
			Filepath: path,
			OutDir:   filepath.Join(filepath.Dir(path), "PDF"),
		})

		return nil
	})

	return jobs, err
}
