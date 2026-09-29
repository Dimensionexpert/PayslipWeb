package resources

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// ResourceFiles contains application resources embedded into the binary.
//
// These files are shipped with the application and will later be
// materialized into the application's runtime resource directory.

//go:embed clusters.xlsx monthly_template.xlsx yearly_template.xlsx
var ResourceFiles embed.FS

// Materialize prepares the destination directory and reads the
// embedded application resources.
//
// path should be the runtime ResourcesDir, for example:
//
//	~/.config/Payslip/resources
func Materialize(path string) error {
	// Make sure the destination directory exists.
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf(
			"creating resource directory %q: %w",
			path,
			err,
		)
	}

	// Read the embedded resources into memory.
	clusters, err := fs.ReadFile(ResourceFiles, "clusters.xlsx")
	if err != nil {
		return fmt.Errorf(
			"reading embedded clusters.xlsx: %w",
			err,
		)
	}

	monthlyTemplate, err := fs.ReadFile(ResourceFiles, "monthly_template.xlsx")
	if err != nil {
		return fmt.Errorf(
			"reading embedded monthly_template.xlsx: %w",
			err,
		)
	}

	yearlyTemplate, err := fs.ReadFile(ResourceFiles, "yearly_template.xlsx")
	if err != nil {
		return fmt.Errorf(
			"reading embedded yearly_template.xlsx: %w",
			err,
		)
	}

	// Build the destination paths for the embedded resources.
	clusterDST := filepath.Join(path, "clusters.xlsx")
	monthlyTemplateDST := filepath.Join(path, "monthly_template.xlsx")
	yearlyTemplateDST := filepath.Join(path, "yearly_template.xlsx")

	// Write the embedded files to disk.
	if err := os.WriteFile(clusterDST, clusters, 0644); err != nil {
		return fmt.Errorf("writing clusters.xlsx: %w", err)
	}

	if err := os.WriteFile(monthlyTemplateDST, monthlyTemplate, 0644); err != nil {
		return fmt.Errorf("writing monthly_template.xlsx: %w", err)
	}

	if err := os.WriteFile(yearlyTemplateDST, yearlyTemplate, 0644); err != nil {
		return fmt.Errorf("writing yearly_template.xlsx: %w", err)
	}
	return nil
}
