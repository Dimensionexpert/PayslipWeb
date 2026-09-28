package generator

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Dimensionexpert/payslip/internal/concurrency"
)

func CollectMissingPDFJobs(root string) ([]concurrency.ConversionJob, error) {
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

		pdfPath := filepath.Join(
			filepath.Dir(path),
			"PDF",
			strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))+".pdf",
		)

		if _, err := os.Stat(pdfPath); err == nil {
			return nil
		} else if !os.IsNotExist(err) {
			return err
		}

		jobs = append(jobs, concurrency.ConversionJob{
			Filepath: path,
			OutDir:   filepath.Join(filepath.Dir(path), "PDF"),
		})

		return nil
	})

	return jobs, err
}
