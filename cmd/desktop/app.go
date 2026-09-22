package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	goRuntime "runtime"

	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/config"
	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/dto"
	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/query"
	"github.com/Dimensionexpert/payslip/internal/concurrency"
	"github.com/Dimensionexpert/payslip/internal/database"
	genexcel "github.com/Dimensionexpert/payslip/internal/genExcel"
	genPDF "github.com/Dimensionexpert/payslip/internal/genPDF"
	"github.com/Dimensionexpert/payslip/internal/generator"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/Dimensionexpert/payslip/internal/models"
)

type App struct {
	ctx    context.Context
	db     *sql.DB
	config config.Config
}

func NewApp() (*App, error) {
	db, err := database.Open("../../payslip.db")
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	cfg, err := config.Load()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("loading config: %w", err)
	}

	if err := config.Save(cfg); err != nil {
		db.Close()
		return nil, fmt.Errorf("saving config: %w", err)
	}

	return &App{
		db:     db,
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
		fmt.Println("closing database:", err)
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

	monthlyTemplate := "../../data/monthly_template.xlsx"

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

	yearlyTemplate := "../../data/yearly_template.xlsx"
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

	outputDir := a.config.OutputDir

	if outputDir == "" {
		return fmt.Errorf("output directory is not configured")
	}

	payslips, err := database.GetPayslips(
		a.db,
		month,
		year,
	)
	if err != nil {
		return fmt.Errorf("fetching payslips: %w", err)
	}

	if len(payslips) == 0 {
		return fmt.Errorf("no payslips found")
	}

	monthlyTemplate := "../../data/monthly_template.xlsx"

	if err := genexcel.GenerateMonthlyPayslips(
		monthlyTemplate,
		outputDir,
		payslips,
	); err != nil {
		return fmt.Errorf(
			"generating monthly Excel files: %w",
			err,
		)
	}

	monthlyXlsxPath := filepath.Join(
		outputDir,
		fmt.Sprintf("%s_%d", time.Month(month), year),
	)

	conversionJobs, err := generator.CollectPDFJobs(
		monthlyXlsxPath,
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

	outputDir := a.config.OutputDir

	if outputDir == "" {
		return fmt.Errorf("output directory is not configured")
	}

	yearlyTemplate := "../../data/yearly_template.xlsx"

	if err := generator.ExportYearlyEmployeePayslips(
		a.db,
		yearlyTemplate,
		outputDir,
		financialYearStart,
	); err != nil {
		return fmt.Errorf(
			"generating yearly Excel files: %w",
			err,
		)
	}

	yearlyXlsxPath := filepath.Join(
		outputDir,
		fmt.Sprintf(
			"Financial_Year_%d_%d",
			financialYearStart,
			financialYearStart+1,
		),
	)

	conversionJobs, err := generator.CollectPDFJobs(
		yearlyXlsxPath,
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
