package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Dimensionexpert/payslip/internal/concurrency"
	"github.com/Dimensionexpert/payslip/internal/database"
	genexcel "github.com/Dimensionexpert/payslip/internal/genExcel"
	"github.com/Dimensionexpert/payslip/internal/importer"
)

func main() {
	totalTime := time.Now()

	// ==================================================
	// 1. Configure input files
	// ==================================================

	clusterPath := "./Data/clusters.xlsx"
	truthPath := "./Source/August_2026_School_All Formate Maval copy.xlsx"
	dbPath := "payslip.db"

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

	fmt.Println("Database opened successfully")

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

	outputDir := "output"

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Println("creating output directory:", err)
		return
	}

	excelStart := time.Now()

	if err := genexcel.GenerateMonthlyPayslips(
		"Data/template.xlsx",
		outputDir,
		exportPayslips,
	); err != nil {
		fmt.Println("generating monthly Excel files:", err)
		return
	}

	fmt.Printf(
		"Time taken to generate Excel files: %v\n",
		time.Since(excelStart),
	)

	fmt.Printf(
		"Generated %d monthly payslips in %s\n",
		len(exportPayslips),
		outputDir,
	)

	// ==================================================
	// 6. Collect monthly Excel files recursively
	//
	// Yearly files are skipped.
	// ==================================================

	var conversionJobs []concurrency.ConversionJob

	err = filepath.WalkDir(
		outputDir,
		func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			if entry.IsDir() {
				return nil
			}

			if !strings.EqualFold(filepath.Ext(entry.Name()), ".xlsx") {
				return nil
			}

			// Do not convert yearly XLSX files as monthly PDFs.
			if strings.Contains(path, "Financial_Year_") {
				return nil
			}

			schoolDir := filepath.Dir(path)
			pdfDir := filepath.Join(schoolDir, "PDF")

			conversionJobs = append(
				conversionJobs,
				concurrency.ConversionJob{
					Filepath: path,
					OutDir:   pdfDir,
				},
			)

			return nil
		},
	)

	if err != nil {
		fmt.Println("walking output directory:", err)
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

	const pdfWorkers = 12

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

	success := 0
	failed := 0

	for _, result := range results {
		if result.Err != nil {
			failed++
		} else {
			success++
		}
	}

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

	const financialYearStart = 2026

	yearlyStart := time.Now()

	yearlyEmployees, err := database.GetEmployees(db)
	if err != nil {
		fmt.Println("fetching employees for yearly payslips:", err)
		return
	}

	if len(yearlyEmployees) == 0 {
		fmt.Println("No employees found for yearly payslips")
		return
	}

	yearlySuccess := 0
	yearlyFailed := 0

	for index, employee := range yearlyEmployees {
		fmt.Printf(
			"Generating yearly payslip %d/%d: %s\n",
			index+1,
			len(yearlyEmployees),
			employee.Name,
		)

		yearly, err := database.GetYearlyPayslip(
			db,
			employee.ShalarthID,
			financialYearStart,
		)
		if err != nil {
			fmt.Printf(
				"FAILED: fetching yearly payslip for %s: %v\n",
				employee.Name,
				err,
			)

			yearlyFailed++
			continue
		}

		if len(yearly.Payslips) == 0 {
			fmt.Printf(
				"SKIPPED: no yearly payslip records for %s\n",
				employee.Name,
			)

			yearlyFailed++
			continue
		}

		yearlyPath, err := genexcel.GenerateYearlyPayslip(
			"Data/yearly_template.xlsx",
			outputDir,
			yearly,
			financialYearStart,
		)
		if err != nil {
			fmt.Printf(
				"FAILED: generating yearly payslip for %s: %v\n",
				employee.Name,
				err,
			)

			yearlyFailed++
			continue
		}

		fmt.Printf(
			"Yearly payslip generated: %s\n",
			yearlyPath,
		)

		yearlySuccess++
	}

	fmt.Printf(
		"Yearly payslip generation: %d succeeded, %d failed in %v\n",
		yearlySuccess,
		yearlyFailed,
		time.Since(yearlyStart),
	)

	// ==================================================
	// 9. Convert yearly Excel files to PDFs concurrently
	//
	// Output:
	//
	// output/Financial_Year_2026_2027/Cluster/School/PDF/Employee.pdf
	// ==================================================

	yearlyPDFStart := time.Now()

	yearlyPDFRoot := filepath.Join(
		outputDir,
		fmt.Sprintf(
			"Financial_Year_%d_%d",
			financialYearStart,
			financialYearStart+1,
		),
	)

	var yearlyConversionJobs []concurrency.ConversionJob

	err = filepath.WalkDir(
		yearlyPDFRoot,
		func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			if entry.IsDir() {
				return nil
			}

			if !strings.EqualFold(filepath.Ext(entry.Name()), ".xlsx") {
				return nil
			}

			schoolDir := filepath.Dir(path)
			pdfDir := filepath.Join(schoolDir, "PDF")

			yearlyConversionJobs = append(
				yearlyConversionJobs,
				concurrency.ConversionJob{
					Filepath: path,
					OutDir:   pdfDir,
				},
			)

			return nil
		},
	)

	if err != nil {
		fmt.Println("walking yearly output directory:", err)
		return
	}

	fmt.Printf(
		"Yearly XLSX files found for PDF conversion: %d\n",
		len(yearlyConversionJobs),
	)

	const yearlyPDFWorkers = 12

	yearlyPDFResults := concurrency.RunPDFConversion(
		yearlyConversionJobs,
		yearlyPDFWorkers,
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

	yearlyPDFSuccess := 0
	yearlyPDFFailed := 0

	for _, result := range yearlyPDFResults {
		if result.Err != nil {
			yearlyPDFFailed++
		} else {
			yearlyPDFSuccess++
		}
	}

	fmt.Printf(
		"Yearly PDF conversion: %d succeeded, %d failed in %v\n",
		yearlyPDFSuccess,
		yearlyPDFFailed,
		time.Since(yearlyPDFStart),
	)

	// ==================================================
	// 10. Total execution time
	// ==================================================

	fmt.Printf(
		"Total time: %v\n",
		time.Since(totalTime),
	)
}
