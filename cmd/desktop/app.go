package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"path/filepath"

	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/config"
	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/dto"
	desktopGenerator "github.com/Dimensionexpert/payslip/cmd/desktop/internal/generator"
	osutil "github.com/Dimensionexpert/payslip/cmd/desktop/internal/osutils"
	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/paths"
	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/query"
	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/services"
	"github.com/Dimensionexpert/payslip/cmd/desktop/resources"
	"github.com/Dimensionexpert/payslip/internal/database"
	"github.com/Dimensionexpert/payslip/internal/importer"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/Dimensionexpert/payslip/internal/models"
)

type App struct {
	ctx    context.Context
	db     *sql.DB
	config config.Config
	paths  dto.AppPaths
}

func NewApp() (*App, error) {
	// Resolve application paths.
	appPaths, err := paths.New()
	if err != nil {
		return nil, fmt.Errorf("resolving application paths: %w", err)
	}

	// Materialize bundled resources.
	if err := resources.Materialize(appPaths.ResourcesDir); err != nil {
		return nil, fmt.Errorf("materializing resources: %w", err)
	}

	// Open the application database.
	db, err := database.Open(appPaths.DBPath)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	// Load application configuration.
	cfg, err := config.Load()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("loading config: %w", err)
	}

	// Save configuration.
	if err := config.Save(cfg); err != nil {
		db.Close()
		return nil, fmt.Errorf("saving config: %w", err)
	}

	return &App{
		db:     db,
		paths:  appPaths,
		config: cfg,
	}, nil
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	if a.db == nil {
		return
	}

	if err := a.db.Close(); err != nil {
		log.Printf("closing database: %v", err)
	}
}

func (a *App) GetEmployee(
	shalarthID string,
) (dto.EmployeeSummary, error) {
	return query.GetEmployee(a.db, shalarthID)
}

func (a *App) GetEmployees() ([]models.Employee, error) {
	return database.GetEmployees(a.db)
}

func (a *App) GetPayslip(
	shalarthID string,
	month int,
	year int,
) (models.PayslipExport, error) {
	return database.GetPayslip(a.db, shalarthID, month, year)
}

func (a *App) GetYearlyPayslip(
	shalarthID string,
	year int,
) (models.PayslipExportYear, error) {
	return database.GetYearlyPayslip(a.db, shalarthID, year)
}

func (a *App) GetSchool(
	udise string,
) (dto.SchoolSummary, error) {
	return query.GetSchool(a.db, udise)
}

func (a *App) GetEmployeesBySchool(
	udise string,
) ([]dto.EmployeeSummary, error) {
	return query.GetEmployeesBySchool(a.db, udise)
}

func (a *App) GetSchoolsByCluster(
	cluster string,
) ([]dto.SchoolSummary, error) {
	return query.GetSchoolsByCluster(a.db, cluster)
}

func (a *App) GenerateMonthlyPayslip(
	shalarthID string,
	month int,
	year int,
) (string, error) {
	monthlyTemplate := filepath.Join(
		a.paths.ResourcesDir,
		"monthly_template.xlsx",
	)

	pdfPath, err := desktopGenerator.GenerateMonthlyPayslip(
		a.db,
		shalarthID,
		month,
		year,
		a.config.OutputDir,
		monthlyTemplate,
	)
	if err != nil {
		return "", fmt.Errorf(
			"[excel] generating individual monthly payslip: %w",
			err,
		)
	}

	return pdfPath, nil
}

func (a *App) ChooseOutputDirectory() (string, error) {
	if a.ctx == nil {
		return "", fmt.Errorf("app context is not initialized")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}

	defaultDirectory := filepath.Join(homeDir, "Downloads")

	path, err := runtime.OpenDirectoryDialog(
		a.ctx,
		runtime.OpenDialogOptions{
			DefaultDirectory:     defaultDirectory,
			Title:                "Choose output folder",
			CanCreateDirectories: true,
		},
	)
	if err != nil {
		return "", fmt.Errorf("opening output directory dialog: %w", err)
	}

	return path, nil
}

func (a *App) GenerateYearlyPayslip(
	shalarthID string,
	financialYearStart int,
) (string, error) {
	yearlyTemplate := filepath.Join(
		a.paths.ResourcesDir,
		"yearly_template.xlsx",
	)

	pdfPath, err := desktopGenerator.GenerateYearlyPayslip(
		a.db,
		shalarthID,
		financialYearStart,
		a.config.OutputDir,
		yearlyTemplate,
	)
	if err != nil {
		return "", fmt.Errorf(
			"[excel] generating individual yearly payslip: %w",
			err,
		)
	}

	return pdfPath, nil
}

func (a *App) GenerateMonthlyPayslips(
	month int,
	year int,
) error {
	return desktopGenerator.GenerateMonthlyPayslips(
		a.db,
		month,
		year,
		a.paths.CacheDir,
		a.paths.ResourcesDir,
		a.config.OutputDir,
	)
}

func (a *App) GenerateYearlyPayslips(
	financialYearStart int,
) error {
	return desktopGenerator.GenerateYearlyPayslips(
		a.db,
		financialYearStart,
		a.paths.CacheDir,
		a.paths.ResourcesDir,
		a.config.OutputDir,
	)
}

func (a *App) SetOutputDirectory() (string, error) {
	if a.ctx == nil {
		return "", fmt.Errorf("app context is not initialized")
	}

	path, err := runtime.OpenDirectoryDialog(
		a.ctx,
		runtime.OpenDialogOptions{
			Title:                "Choose Payslip Output Folder",
			CanCreateDirectories: true,
		},
	)
	if err != nil {
		return "", fmt.Errorf("opening output directory dialog: %w", err)
	}

	if path == "" {
		return "", nil
	}

	a.config.OutputDir = path

	if err := config.Save(a.config); err != nil {
		return "", fmt.Errorf("saving config: %w", err)
	}

	return path, nil
}

func (a *App) GetOutputDirectory() string {
	return a.config.OutputDir
}

func (a *App) OpenMonthlyPayslip(
	shalarthID string,
	month int,
	year int,
) error {
	return desktopGenerator.OpenMonthlyPayslip(
		a.db,
		shalarthID,
		month,
		year,
		a.config.OutputDir,
	)
}

func (a *App) OpenYearlyPayslip(
	shalarthID string,
	financialYearStart int,
) error {
	return desktopGenerator.OpenYearlyPayslip(
		a.db,
		shalarthID,
		financialYearStart,
		a.config.OutputDir,
	)
}

func (a *App) GetClusters() ([]dto.ClusterSummary, error) {
	return query.GetClusters(a.db)
}

func (a *App) GenerateMonthlyPayslipsForScope(
	scope string,
	cluster string,
	udise string,
	month int,
	year int,
) (int, error) {
	return desktopGenerator.GenerateMonthlyPayslipsForScope(
		a.db,
		a.paths.CacheDir,
		a.paths.ResourcesDir,
		a.config.OutputDir,
		scope,
		cluster,
		udise,
		month,
		year,
	)
}

func (a *App) GenerateYearlyPayslipsForScope(
	scope string,
	cluster string,
	udise string,
	financialYearStart int,
) (int, error) {
	return desktopGenerator.GenerateYearlyPayslipsForScope(
		a.db,
		a.paths.CacheDir,
		a.paths.ResourcesDir,
		a.config.OutputDir,
		scope,
		cluster,
		udise,
		financialYearStart,
	)
}

func (a *App) ImportPayroll(
	truthPath string,
	month int,
	year int,
) (importer.ImportReport, error) {
	return services.ImportPayroll(
		truthPath,
		month,
		year,
		a.paths.SourceDir,
		a.paths.ResourcesDir,
		a.paths.DBPath,
		a.paths.CacheDir,
	)
}

func (a *App) OpenDir(path string) error {
	return osutil.OpenDir(path)
}

func (a *App) GetMonthlyBulkOutputDirectory(
	scope string,
	cluster string,
	udise string,
	month int,
	year int,
) (string, error) {
	return desktopGenerator.ResolveMonthlyScopeRoot(
		a.db,
		a.config.OutputDir,
		scope,
		cluster,
		udise,
		month,
		year,
	)
}

func (a *App) GetYearlyBulkOutputDirectory(
	scope string,
	cluster string,
	udise string,
	financialYearStart int,
) (string, error) {
	return desktopGenerator.YearlyPayslipRoot(
		a.db,
		cluster,
		scope,
		udise,
		financialYearStart,
		a.config.OutputDir,
	)
}

func (a *App) SelectPayrollFile() (string, error) {
	if a.ctx == nil {
		return "", fmt.Errorf("app context is not initialized")
	}

	path, err := runtime.OpenFileDialog(
		a.ctx,
		runtime.OpenDialogOptions{
			Title: "Select monthly payroll file",
			Filters: []runtime.FileFilter{
				{
					DisplayName: "Excel Files (*.xlsx)",
					Pattern:     "*.xlsx",
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("opening payroll file dialog: %w", err)
	}

	return path, nil
}
