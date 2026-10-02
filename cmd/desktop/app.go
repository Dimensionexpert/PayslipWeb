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
	"strings"
	"time"

	goRuntime "runtime"

	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/config"
	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/dto"
	osutil "github.com/Dimensionexpert/payslip/cmd/desktop/internal/osutils"
	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/paths"
	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/query"
	"github.com/Dimensionexpert/payslip/cmd/desktop/resources"
	"github.com/Dimensionexpert/payslip/internal/concurrency"
	"github.com/Dimensionexpert/payslip/internal/database"
	genexcel "github.com/Dimensionexpert/payslip/internal/genExcel"
	genPDF "github.com/Dimensionexpert/payslip/internal/genPDF"
	"github.com/Dimensionexpert/payslip/internal/generator"
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

func (a *App) monthlyScopeRoot(
	baseDir string,
	scope string,
	cluster string,
	udise string,
	month int,
	year int,
) (string, error) {
	if baseDir == "" {
		return "", fmt.Errorf("base directory is not configured")
	}

	root := filepath.Join(
		baseDir,
		fmt.Sprintf("%s_%d", time.Month(month), year),
	)

	switch scope {
	case "all":
		return root, nil

	case "school":
		employees, err := query.GetEmployeesBySchool(a.db, udise)
		if err != nil {
			return "", err
		}

		if len(employees) == 0 {
			return "", fmt.Errorf("no employees found for school %s", udise)
		}

		payslip, err := database.GetPayslip(
			a.db,
			employees[0].ShalarthID,
			month,
			year,
		)
		if err != nil {
			return "", err
		}

		pdfPath := genexcel.MonthlyPayslipPDFPath(
			baseDir,
			payslip,
		)

		return filepath.Dir(filepath.Dir(pdfPath)), nil

	case "cluster":
		schools, err := query.GetSchoolsByCluster(a.db, cluster)
		if err != nil {
			return "", err
		}

		for _, school := range schools {
			employees, err := query.GetEmployeesBySchool(
				a.db,
				school.UDISECode,
			)
			if err != nil {
				return "", err
			}

			if len(employees) == 0 {
				continue
			}

			payslip, err := database.GetPayslip(
				a.db,
				employees[0].ShalarthID,
				month,
				year,
			)
			if err != nil {
				return "", err
			}

			pdfPath := genexcel.MonthlyPayslipPDFPath(
				baseDir,
				payslip,
			)

			return filepath.Dir(
				filepath.Dir(
					filepath.Dir(pdfPath),
				),
			), nil
		}

		return "", fmt.Errorf("no employees found in cluster %s", cluster)

	default:
		return "", fmt.Errorf("invalid generation scope: %s", scope)
	}
}

func (a *App) GenerateMonthlyPayslip(
	shalarthID string,
	month int,
	year int,
) (string, error) {

	outputDir := a.config.OutputDir

	if outputDir == "" {
		return "", fmt.Errorf("output directory is not configured")
	}

	// XLSX generation

	payslip, err := database.GetPayslip(
		a.db,
		shalarthID,
		month,
		year,
	)
	if err != nil {
		return "", err
	}

	monthlyTemplate := filepath.Join(a.paths.ResourcesDir, "monthly_template.xlsx")

	xlsxPath, err := genexcel.GenerateMonthlyPayslip(
		monthlyTemplate,
		outputDir,
		payslip,
	)
	if err != nil {
		return "", err
	}

	// PDF conversion

	pdfDir := filepath.Join(
		filepath.Dir(xlsxPath),
		"PDF",
	)

	if err := genPDF.ConvertToPDF(
		xlsxPath,
		pdfDir,
		0,
	); err != nil {
		return "", err
	}

	pdfFilename := strings.TrimSuffix(
		filepath.Base(xlsxPath),
		filepath.Ext(xlsxPath),
	) + ".pdf"

	pdfPath := filepath.Join(
		pdfDir,
		pdfFilename,
	)

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

	// XLSX generation

	yearly, err := database.GetYearlyPayslip(
		a.db,
		shalarthID,
		financialYearStart,
	)
	if err != nil {
		return "", err
	}

	yearlyTemplate := filepath.Join(a.paths.ResourcesDir, "yearly_template.xlsx")
	outputDir := a.config.OutputDir

	if outputDir == "" {
		return "", fmt.Errorf("output directory is not configured")
	}

	xlsxPath, err := genexcel.GenerateYearlyPayslip(
		yearlyTemplate,
		outputDir,
		yearly,
		financialYearStart,
	)
	if err != nil {
		return "", err
	}

	// PDF conversion

	pdfDir := filepath.Join(
		filepath.Dir(xlsxPath),
		"PDF",
	)

	if err := genPDF.ConvertToPDF(
		xlsxPath,
		pdfDir,
		0,
	); err != nil {
		return "", err
	}

	pdfFilename := strings.TrimSuffix(
		filepath.Base(xlsxPath),
		filepath.Ext(xlsxPath),
	) + ".pdf"

	pdfPath := filepath.Join(
		pdfDir,
		pdfFilename,
	)

	return pdfPath, nil
}

func (a *App) GenerateMonthlyPayslips(
	month int,
	year int,
) error {

	monthlyRoot, err := a.EnsureMonthlyXLSX(
		month,
		year,
	)
	if err != nil {
		return err
	}

	if a.config.OutputDir == "" {
		return fmt.Errorf("output directory is not configured")
	}

	monthlyOutputRoot := filepath.Join(
		a.config.OutputDir,
		fmt.Sprintf("%s_%d", time.Month(month), year),
	)

	conversionJobs, err := generator.CollectMissingPDFJobs(
		monthlyRoot,
		monthlyOutputRoot,
	)
	if err != nil {
		return fmt.Errorf(
			"collecting monthly PDF jobs: %w",
			err,
		)
	}

	results := concurrency.RunPDFConversion(
		conversionJobs,
		8,
		func(result concurrency.ConversionResult) {
			if result.Err != nil {
				fmt.Printf(
					"MONTHLY PDF FAILED: %s: %v\n",
					result.Filepath,
					result.Err,
				)
			}
		},
	)

	success, failed := generator.CountConversionResults(results)

	fmt.Printf(
		"Monthly PDF conversion: %d succeeded, %d failed\n",
		success,
		failed,
	)

	if failed > 0 {
		return fmt.Errorf(
			"monthly PDF conversion failed for %d file(s)",
			failed,
		)
	}

	return nil
}

func (a *App) GenerateYearlyPayslips(
	financialYearStart int,
) error {

	yearlyRoot, err := a.EnsureYearlyXLSX(
		financialYearStart,
	)
	if err != nil {
		return err
	}

	if a.config.OutputDir == "" {
		return fmt.Errorf("output directory is not configured")
	}

	yearlyOutputRoot := filepath.Join(
		a.config.OutputDir,
		fmt.Sprintf(
			"Financial_Year_%d_%d",
			financialYearStart,
			financialYearStart+1,
		),
	)

	conversionJobs, err := generator.CollectMissingPDFJobs(
		yearlyRoot,
		yearlyOutputRoot,
	)
	if err != nil {
		return fmt.Errorf(
			"collecting yearly PDF jobs: %w",
			err,
		)
	}

	results := concurrency.RunPDFConversion(
		conversionJobs,
		8,
		func(result concurrency.ConversionResult) {
			if result.Err != nil {
				fmt.Printf(
					"YEARLY PDF FAILED: %s: %v\n",
					result.Filepath,
					result.Err,
				)
			}
		},
	)

	success, failed := generator.CountConversionResults(results)

	fmt.Printf(
		"Yearly PDF conversion: %d succeeded, %d failed\n",
		success,
		failed,
	)

	if failed > 0 {
		return fmt.Errorf(
			"yearly PDF conversion failed for %d file(s)",
			failed,
		)
	}

	return nil
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

	payslip, err := database.GetPayslip(
		a.db,
		shalarthID,
		month,
		year,
	)
	if err != nil {
		return fmt.Errorf("fetching payslip: %w", err)
	}

	outputDir := a.config.OutputDir

	if outputDir == "" {
		return fmt.Errorf("output directory is not configured")
	}

	pdfPath := genexcel.MonthlyPayslipPDFPath(
		outputDir,
		payslip,
	)

	if _, err := os.Stat(pdfPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(
				"payslip PDF not found: %s",
				pdfPath,
			)
		}

		return fmt.Errorf(
			"checking payslip PDF: %w",
			err,
		)
	}

	// Native system viewer launcher replacing BrowserOpenURL
	var cmd *exec.Cmd
	switch goRuntime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", pdfPath)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", pdfPath)
	case "darwin":
		cmd = exec.Command("open", pdfPath)
	default:
		return fmt.Errorf("unsupported platform for opening files")
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open PDF launcher: %w", err)
	}

	fmt.Println(pdfPath)

	return nil
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

func (a *App) EnsureMonthlyXLSX(
	month int,
	year int,
) (string, error) {

	payslips, err := database.GetPayslips(
		a.db,
		month,
		year,
	)
	if err != nil {
		return "", fmt.Errorf("fetching payslips: %w", err)
	}

	if len(payslips) == 0 {
		return "", fmt.Errorf("no payslips found")
	}

	cacheDir := a.paths.CacheDir

	monthlyRoot := filepath.Join(
		cacheDir,
		fmt.Sprintf("%s_%d", time.Month(month), year),
	)

	readyMarker := filepath.Join(
		monthlyRoot,
		".xlsx-ready",
	)

	if _, err := os.Stat(readyMarker); err == nil {
		return monthlyRoot, nil
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("checking monthly XLSX cache: %w", err)
	}

	monthlyTemplate := filepath.Join(
		a.paths.ResourcesDir,
		"monthly_template.xlsx",
	)

	if err := genexcel.GenerateMonthlyPayslips(
		monthlyTemplate,
		cacheDir,
		payslips,
	); err != nil {
		return "", fmt.Errorf(
			"generating monthly Excel files: %w",
			err,
		)
	}

	if err := os.WriteFile(
		readyMarker,
		[]byte{},
		0644,
	); err != nil {
		return "", fmt.Errorf(
			"writing monthly XLSX ready marker: %w",
			err,
		)
	}

	return monthlyRoot, nil
}

func (a *App) EnsureYearlyXLSX(
	financialYearStart int,
) (string, error) {

	cacheDir := a.paths.CacheDir

	yearlyRoot := filepath.Join(
		cacheDir,
		fmt.Sprintf(
			"Financial_Year_%d_%d",
			financialYearStart,
			financialYearStart+1,
		),
	)

	readyMarker := filepath.Join(
		yearlyRoot,
		".xlsx-ready",
	)

	if _, err := os.Stat(readyMarker); err == nil {
		return yearlyRoot, nil
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("checking yearly XLSX cache: %w", err)
	}

	yearlyTemplate := filepath.Join(
		a.paths.ResourcesDir,
		"yearly_template.xlsx",
	)

	if err := generator.ExportYearlyEmployeePayslips(
		a.db,
		yearlyTemplate,
		cacheDir,
		financialYearStart,
	); err != nil {
		return "", fmt.Errorf(
			"generating yearly Excel files: %w",
			err,
		)
	}

	if err := os.WriteFile(
		readyMarker,
		[]byte{},
		0644,
	); err != nil {
		return "", fmt.Errorf(
			"writing yearly XLSX ready marker: %w",
			err,
		)
	}

	return yearlyRoot, nil
}

func (a *App) GenerateMonthlyPayslipsForScope(
	scope string,
	cluster string,
	udise string,
	month int,
	year int,
) (int, error) {
	// Make sure all XLSX files for this month exist.
	if _, err := a.EnsureMonthlyXLSX(month, year); err != nil {
		return 0, err
	}

	// Pick the directory based on the requested scope.
	root, err := a.monthlyScopeRoot(
		a.paths.CacheDir,
		scope,
		cluster,
		udise,
		month,
		year,
	)
	if err != nil {
		return 0, err
	}

	outputRoot, err := a.monthlyScopeRoot(
		a.config.OutputDir,
		scope,
		cluster,
		udise,
		month,
		year,
	)
	if err != nil {
		return 0, err
	}

	// Convert only PDFs that don't exist yet.
	jobs, err := generator.CollectMissingPDFJobs(
		root,
		outputRoot,
	)
	if err != nil {
		return 0, fmt.Errorf("collecting PDF jobs: %w", err)
	}
	if len(jobs) == 0 {
		return 0, nil
	}

	results := concurrency.RunPDFConversion(
		jobs,
		8,
		func(result concurrency.ConversionResult) {
			if result.Err != nil {
				fmt.Printf(
					"MONTHLY PDF FAILED: %s: %v\n",
					result.Filepath,
					result.Err,
				)
			}
		},
	)

	_, failed := generator.CountConversionResults(results)

	if failed > 0 {
		return 0, fmt.Errorf(
			"monthly PDF conversion failed for %d file(s)",
			failed,
		)
	}

	return len(jobs), nil
}

func (a *App) yearlyScopeRoot(
	baseDir string,
	scope string,
	cluster string,
	udise string,
	financialYearStart int,
) (string, error) {
	if baseDir == "" {
		return "", fmt.Errorf("base directory is not configured")
	}

	root := filepath.Join(
		baseDir,
		fmt.Sprintf(
			"Financial_Year_%d_%d",
			financialYearStart,
			financialYearStart+1,
		),
	)

	switch scope {
	case "all":
		return root, nil

	case "school":
		employees, err := query.GetEmployeesBySchool(a.db, udise)
		if err != nil {
			return "", err
		}

		if len(employees) == 0 {
			return "", fmt.Errorf("no employees found for school %s", udise)
		}

		yearly, err := database.GetYearlyPayslip(
			a.db,
			employees[0].ShalarthID,
			financialYearStart,
		)
		if err != nil {
			return "", err
		}

		pdfPath := genexcel.YearlyPayslipPDFPath(
			baseDir,
			yearly,
			financialYearStart,
		)

		return filepath.Dir(filepath.Dir(pdfPath)), nil

	case "cluster":
		schools, err := query.GetSchoolsByCluster(a.db, cluster)
		if err != nil {
			return "", err
		}

		for _, school := range schools {
			employees, err := query.GetEmployeesBySchool(
				a.db,
				school.UDISECode,
			)
			if err != nil {
				return "", err
			}

			if len(employees) == 0 {
				continue
			}

			yearly, err := database.GetYearlyPayslip(
				a.db,
				employees[0].ShalarthID,
				financialYearStart,
			)
			if err != nil {
				continue
			}

			pdfPath := genexcel.YearlyPayslipPDFPath(
				baseDir,
				yearly,
				financialYearStart,
			)

			return filepath.Dir(
				filepath.Dir(
					filepath.Dir(pdfPath),
				),
			), nil
		}

		return "", fmt.Errorf("no employees found in cluster %s", cluster)

	default:
		return "", fmt.Errorf("invalid generation scope: %s", scope)
	}
}

func (a *App) GenerateYearlyPayslipsForScope(
	scope string,
	cluster string,
	udise string,
	financialYearStart int,
) (int, error) {
	if _, err := a.EnsureYearlyXLSX(financialYearStart); err != nil {
		return 0, err
	}

	root, err := a.yearlyScopeRoot(
		a.paths.CacheDir,
		scope,
		cluster,
		udise,
		financialYearStart,
	)
	if err != nil {
		return 0, err
	}

	outputRoot, err := a.yearlyScopeRoot(
		a.config.OutputDir,
		scope,
		cluster,
		udise,
		financialYearStart,
	)
	if err != nil {
		return 0, err
	}

	jobs, err := generator.CollectMissingPDFJobs(
		root,
		outputRoot,
	)
	if err != nil {
		return 0, fmt.Errorf(
			"collecting PDF jobs: %w",
			err,
		)
	}
	if len(jobs) == 0 {
		return 0, nil
	}

	results := concurrency.RunPDFConversion(
		jobs,
		8,
		func(result concurrency.ConversionResult) {
			if result.Err != nil {
				fmt.Printf(
					"YEARLY PDF FAILED: %s: %v\n",
					result.Filepath,
					result.Err,
				)
			}
		},
	)

	_, failed := generator.CountConversionResults(results)

	if failed > 0 {
		return 0, fmt.Errorf(
			"yearly PDF conversion failed for %d file(s)",
			failed,
		)
	}

	return len(jobs), nil
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
	return a.monthlyScopeRoot(
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
	return a.yearlyScopeRoot(
		a.config.OutputDir,
		scope,
		cluster,
		udise,
		financialYearStart,
	)
}
