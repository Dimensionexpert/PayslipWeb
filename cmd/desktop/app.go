package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/dto"
	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/query"
	"github.com/Dimensionexpert/payslip/internal/database"
	genexcel "github.com/Dimensionexpert/payslip/internal/genExcel"
	genPDF "github.com/Dimensionexpert/payslip/internal/genPDF"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/Dimensionexpert/payslip/internal/models"
)

type App struct {
	ctx context.Context
	db  *sql.DB
}

func NewApp() (*App, error) {
	db, err := database.Open("../../payslip.db")
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	return &App{db: db}, nil
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
	outputDir string,
) (string, error) {

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
	outputDir string,
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
