package pdfgen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ConvertToPDF(
	inputPath string,
	outputDir string,
	workerID int,
) error {

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	profileDir := fmt.Sprintf(
		"/tmp/lo_profile_%d",
		workerID,
	)

	cmd := exec.Command(
		"soffice",
		"--headless",
		"-env:UserInstallation=file://"+profileDir,
		"--convert-to",
		"pdf",
		"--outdir",
		outputDir,
		inputPath,
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"converting %s to PDF: %w",
			filepath.Base(inputPath),
			err,
		)
	}

	return nil
}
