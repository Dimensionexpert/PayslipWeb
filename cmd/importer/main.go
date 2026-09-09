package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	// yearlyXlsxPath := "output/Financial_Year_2026_2027"
	monthly_template := "data/monthly_template.xlsx"
	yearly_template := "data/yearly_template.xlsx"

	const financialYearStart = 2026
	const pdfWorkers = 12

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
	//
	// Yearly files are skipped.
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

	success := 0
	failed := 0

	success, failed = generator.CountConversionResults(results)

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
