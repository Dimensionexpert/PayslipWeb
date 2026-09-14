package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/dto"
	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/query"
	"github.com/Dimensionexpert/payslip/internal/database"
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
