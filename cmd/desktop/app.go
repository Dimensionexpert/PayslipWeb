package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"path/filepath"
	"time"

	goRuntime "runtime"

	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/config"
	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/dto"
	desktopGenerator "github.com/Dimensionexpert/payslip/cmd/desktop/internal/generator"
	osutil "github.com/Dimensionexpert/payslip/cmd/desktop/internal/osutils"
	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/paths"
	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/query"
	"github.com/Dimensionexpert/payslip/cmd/desktop/resources"
	"github.com/Dimensionexpert/payslip/internal/database"
	genexcel "github.com/Dimensionexpert/payslip/internal/genExcel"
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

	yearly, err := database.GetYearlyPayslip(
		a.db,
		shalarthID,
		financialYearStart,
	)
	if err != nil {
		return fmt.Errorf("fetching yearly payslip: %w", err)
	}

	outputDir := a.config.OutputDir

	if outputDir == "" {
		return fmt.Errorf("output directory is not configured")
	}

	pdfPath := genexcel.YearlyPayslipPDFPath(
		outputDir,
		yearly,
		financialYearStart,
	)

	if _, err := os.Stat(pdfPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(
				"yearly payslip PDF not found: %s",
				pdfPath,
			)
		}

		return fmt.Errorf(
			"checking yearly payslip PDF: %w",
			err,
		)
	}

	var cmd *exec.Cmd

	switch goRuntime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", pdfPath)
	case "windows":
		cmd = exec.Command(
			"rundll32",
			"url.dll,FileProtocolHandler",
			pdfPath,
		)
	case "darwin":
		cmd = exec.Command("open", pdfPath)
	default:
		return fmt.Errorf(
			"unsupported platform for opening files",
		)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf(
			"failed to open PDF launcher: %w",
			err,
		)
	}

	fmt.Println(pdfPath)

	return nil
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
	if err := os.MkdirAll(a.paths.SourceDir, 0755); err != nil {
		return importer.ImportReport{},
			fmt.Errorf("creating source directory: %w", err)
	}

	destination := filepath.Join(
		a.paths.SourceDir,
		filepath.Base(truthPath),
	)

	sourceAbs, err := filepath.Abs(truthPath)
	if err != nil {
		return importer.ImportReport{},
			fmt.Errorf("resolving payroll file path: %w", err)
	}

	destinationAbs, err := filepath.Abs(destination)
	if err != nil {
		return importer.ImportReport{},
			fmt.Errorf("resolving source path: %w", err)
	}

	if sourceAbs != destinationAbs {
		src, err := os.Open(truthPath)
		if err != nil {
			return importer.ImportReport{},
				fmt.Errorf("opening selected payroll file: %w", err)
		}
		defer src.Close()

		dst, err := os.Create(destination)
		if err != nil {
			return importer.ImportReport{},
				fmt.Errorf("creating source payroll file: %w", err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return importer.ImportReport{},
				fmt.Errorf("copying payroll file: %w", err)
		}
	}

	clusterPath := filepath.Join(
		a.paths.ResourcesDir,
		"clusters.xlsx",
	)

	report, _, _, err := importer.ImportPayroll(
		clusterPath,
		destination,
		a.paths.DBPath,
		month,
		year,
	)
	if err != nil {
		return report, err
	}

	// Invalidate cached XLSX files for the imported month.
	monthlyCacheRoot := filepath.Join(
		a.paths.CacheDir,
		fmt.Sprintf("%s_%d", time.Month(month), year),
	)

	if err := os.RemoveAll(monthlyCacheRoot); err != nil {
		return report, fmt.Errorf(
			"invalidating monthly XLSX cache: %w",
			err,
		)
	}

	return report, nil
}

func (a *App) SelectPayrollFile() (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Payroll File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Excel Files (*.xlsx)",
				Pattern:     "*.xlsx",
			},
		},
	})
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
