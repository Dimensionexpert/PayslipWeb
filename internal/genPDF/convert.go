package pdfgen

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ConvertToPDF converts an Excel file to PDF using LibreOffice.
func ConvertToPDF(inputPath, outputDir string, workerID int) error {
	profileDir := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("lo_profile_%d", workerID),
	)

	// Convert the profile directory to a URI.
	profilePath := filepath.ToSlash(profileDir)

	var profileURI string

	if runtime.GOOS == "windows" {
		// Example: file:///C:/Users/user/AppData/Local/Temp/lo_profile_0
		profileURI = "file:///" + profilePath
	} else {
		profileURI = (&url.URL{
			Scheme: "file",
			Path:   profilePath,
		}).String()
	}

	var soffice string

	switch runtime.GOOS {
	case "linux":
		soffice = "soffice"

	case "windows":
		soffice = `C:\Program Files\LibreOffice\program\soffice.com`

	default:
		return fmt.Errorf(
			"unsupported operating system: %s",
			runtime.GOOS,
		)
	}

	cmd := exec.Command(
		soffice,
		"-env:UserInstallation="+profileURI,
		"--headless",
		"--convert-to",
		"pdf",
		"--outdir",
		outputDir,
		inputPath,
	)

	// Diagnostic output.
	fmt.Printf("LibreOffice executable: %s\n", soffice)
	fmt.Printf("LibreOffice profile URI: %s\n", profileURI)
	fmt.Printf("LibreOffice arguments: %#v\n", cmd.Args)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"failed to convert %q to PDF: %w\nLibreOffice output:\n%s",
			inputPath,
			err,
			strings.TrimSpace(string(output)),
		)
	}

	// Verify that the expected PDF was generated.
	pdfPath := filepath.Join(
		outputDir,
		strings.TrimSuffix(
			filepath.Base(inputPath),
			filepath.Ext(inputPath),
		)+".pdf",
	)

	info, err := os.Stat(pdfPath)
	if err != nil {
		return fmt.Errorf(
			"LibreOffice exited successfully but PDF was not found at %q: %w",
			pdfPath,
			err,
		)
	}

	if info.IsDir() || info.Size() == 0 {
		return fmt.Errorf(
			"generated PDF is missing or empty: %s",
			pdfPath,
		)
	}

	return nil
}
