package main

import (
	"fmt"
	"time"

	"github.com/Dimensionexpert/payslip/internal/concurrency"
	"github.com/Dimensionexpert/payslip/internal/database"
	"github.com/Dimensionexpert/payslip/internal/generator"
	"github.com/Dimensionexpert/payslip/internal/importer"
)

func main() {
	totalTime := time.Now()

	// ==================================================
	// 1. Configure input files
	// ==================================================

	clusterPath := "./data/clusters.xlsx"
	truthPath := "./source/August_2026_School_All Formate Maval copy.xlsx"
	dbPath := "payslip.db"

	outputDir := "output"

	monthlyXlsxPath := "output/August_2026"
	yearlyXlsxPath := "output/Financial_Year_2026_2027"
	monthly_template := "data/monthly_template.xlsx"
	yearly_template := "data/yearly_template.xlsx"

	const financialYearStart = 2026
	const pdfWorkers = 8

	// ==================================================
	// 2. Import payroll data into the database
	// ==================================================

	report, month, year, err := importer.ImportPayroll(
		clusterPath,
		truthPath,
		dbPath,
	)
	if err != nil {
		fmt.Println("Import failed:", err)
		return
	}

	report.Print()

	// ==================================================
	// 3. Open database for generation workflow
	// ==================================================

	db, err := database.Open(dbPath)
	if err != nil {
		fmt.Println("opening database:", err)
		return
	}
	defer db.Close()

	// ==================================================
	// 4. Fetch imported monthly payslips
	// ==================================================

	exportPayslips, err := database.GetPayslips(db, month, year)
	if err != nil {
		fmt.Println("fetching monthly payslips:", err)
		return
	}

	fmt.Printf(
		"Payslips found for %02d/%d: %d\n",
		month,
		year,
		len(exportPayslips),
	)

	if len(exportPayslips) == 0 {
		fmt.Println("No payslips found")
		return
	}

	// ==================================================
	// 5. Generate monthly Excel files
	//
	// Output:
	//
	// output/Month_Year/Cluster/School/Employee.xlsx
	// ==================================================

	if err := generator.GenerateMonthlyExcel(
		monthly_template,
		outputDir,
		exportPayslips,
	); err != nil {
		fmt.Println("generating monthly Excel files:", err)
		return
	}

	fmt.Printf("Generated payslips: %d\n", len(exportPayslips))

	// ==================================================
	// 6. Collect monthly Excel files recursively
	// ==================================================

	conversionJobs, err := generator.CollectPDFJobs(monthlyXlsxPath)
	if err != nil {
		fmt.Println("collecting monthly PDF jobs:", err)
		return
	}

	fmt.Printf(
		"XLSX files found for PDF conversion: %d\n",
		len(conversionJobs),
	)

	// ==================================================
	// 7. Convert monthly Excel files to PDFs concurrently
	//
	// Output:
	//
	// output/Month_Year/Cluster/School/PDF/Employee.pdf
	// ==================================================

	pdfStart := time.Now()

	results := concurrency.RunPDFConversion(
		conversionJobs,
		pdfWorkers,
		func(result concurrency.ConversionResult) {
			if result.Err != nil {
				fmt.Printf(
					"FAILED: %s: %v\n",
					result.Filepath,
					result.Err,
				)
				return
			}

			fmt.Printf(
				"PDF generated: %s\n",
				result.Filepath,
			)
		},
	)

	success, failed := generator.CountConversionResults(results)

	fmt.Printf(
		"PDF conversion: %d succeeded, %d failed in %v\n",
		success,
		failed,
		time.Since(pdfStart),
	)

	// ==================================================
	// 8. Generate yearly payslips for every employee
	//
	// Financial year:
	// April 2026 to March 2027
	// ==================================================

	yearlyStart := time.Now()
	err = generator.ExportYearlyEmployeePayslips(
		db,
		yearly_template,
		outputDir,
		financialYearStart,
	)
	if err != nil {
		fmt.Println("yearly payslip export failed:", err)
		return
	}
	fmt.Println(time.Since(yearlyStart))

	// ==================================================
	// 9. Collect yearly Excel files recursively
	// ==================================================

	yearlyConversionJobs, err := generator.CollectPDFJobs(yearlyXlsxPath)
	if err != nil {
		fmt.Println("collecting yearly PDF jobs:", err)
		return
	}

	fmt.Printf(
		"Yearly XLSX files found for PDF conversion: %d\n",
		len(yearlyConversionJobs),
	)

	// ==================================================
	// 10. Convert yearly Excel files to PDFs concurrently
	//
	// Output:
	//
	// output/Financial_Year_2026_2027/Cluster/School/PDF/Employee.pdf
	// ==================================================

	yearlyPDFStart := time.Now()

	yearlyPDFResults := concurrency.RunPDFConversion(
		yearlyConversionJobs,
		pdfWorkers,
		func(result concurrency.ConversionResult) {
			if result.Err != nil {
				fmt.Printf(
					"YEARLY PDF FAILED: %s: %v\n",
					result.Filepath,
					result.Err,
				)
				return
			}

			fmt.Printf(
				"Yearly PDF generated: %s\n",
				result.Filepath,
			)
		},
	)

	yearlyPDFSuccess, yearlyPDFFailed :=
		generator.CountConversionResults(yearlyPDFResults)

	fmt.Printf(
		"Yearly PDF conversion: %d succeeded, %d failed in %v\n",
		yearlyPDFSuccess,
		yearlyPDFFailed,
		time.Since(yearlyPDFStart),
	)

	// ==================================================
	// 11. Total execution time
	// ==================================================

	fmt.Printf(
		"Total time: %v\n",
		time.Since(totalTime),
	)
}
