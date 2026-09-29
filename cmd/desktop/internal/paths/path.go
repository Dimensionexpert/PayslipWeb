package paths

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/dto"
)

func New() (dto.AppPaths, error) {

	// Get the OS-specific directory intended for application data/configuration.
	configDir, err := os.UserConfigDir()
	if err != nil {
		return dto.AppPaths{}, fmt.Errorf(
			"getting user config directory: %w",
			err,
		)
	}

	// Get the OS-specific directory intended for cache data.
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return dto.AppPaths{}, fmt.Errorf(
			"getting user cache directory: %w",
			err,
		)
	}

	// Application data directory.
	// Linux:
	//   ~/.config/Payslip
	// Windows:
	//   %APPDATA%\Payslip
	dataDir := filepath.Join(configDir, "Payslip")

	// Application cache directory.
	// Linux:
	//   ~/.cache/Payslip
	// Windows:
	//   %LocalAppData%\Payslip
	appCacheDir := filepath.Join(cacheDir, "Payslip")

	sourceDir := filepath.Join(dataDir, "source")
	dbPath := filepath.Join(dataDir, "payslip.db")
	resourceDir := filepath.Join(dataDir, "resources")

	// Create directories required at runtime.

	if err := os.MkdirAll(resourceDir, 0755); err != nil {
		return dto.AppPaths{}, fmt.Errorf("creating resource directory: %w", err)
	}

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		return dto.AppPaths{}, fmt.Errorf(
			"creating source directory: %w",
			err,
		)
	}

	if err := os.MkdirAll(appCacheDir, 0755); err != nil {
		return dto.AppPaths{}, fmt.Errorf(
			"creating cache directory: %w",
			err,
		)
	}

	return dto.AppPaths{
		ResourcesDir: resourceDir,
		DataDir:      dataDir,
		DBPath:       dbPath,
		SourceDir:    sourceDir,
		CacheDir:     appCacheDir,
	}, nil
}
