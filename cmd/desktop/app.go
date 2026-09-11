package main

import (
	"context"
	"database/sql"
	"fmt"

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

	return &App{
		db: db,
	}, nil
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
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

func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			fmt.Println("closing db error:", err)
		}
	}
}
