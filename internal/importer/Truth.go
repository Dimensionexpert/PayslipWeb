package importer

import (
	"fmt"

	"github.com/Dimensionexpert/payslip/internal/excel"
	"github.com/Dimensionexpert/payslip/internal/models"
	"github.com/Dimensionexpert/payslip/internal/parser"
)

func ImportTruth(
	rows [][]string,
	clusterMap map[string]string,
	month int,
	year int,
) (
	[]models.Employee,
	map[string]models.School,
	[]models.PayslipRecord,
	ImportReport,
	error,
) {
	report := ImportReport{
		ClusterMappingCount: len(clusterMap),
	}

	headerMap := excel.BuildHeaderMap(rows[4])

	employees := []models.Employee{}
	schools := make(map[string]models.School)
	payslips := []models.PayslipRecord{}

	missingClusters := make(map[string]string)
	unknownUDISE := make(map[string]struct{})

	for i := 5; i < len(rows); i++ {
		row := rows[i]

		udiseCode := excel.Get(row, headerMap, "SCHOOL UDISE CODE")
		schoolName := excel.Get(row, headerMap, "NAME OF SCHOOL")

		// Ignore completely empty rows.
		if udiseCode == "" && schoolName == "" {
			continue
		}

		// Check whether the UDISE exists in the cluster mapping.
		clusterName, ok := clusterMap[udiseCode]

		if !ok {
			if udiseCode != "" {
				unknownUDISE[udiseCode] = struct{}{}
			}
		} else if clusterName == "" {
			missingClusters[udiseCode] = schoolName
		}

		// Build school.
		school := parser.BuildSchool(
			row,
			headerMap,
			clusterMap,
		)

		schools[udiseCode] = school

		// Build employee.
		employee := parser.BuildEmployee(
			row,
			headerMap,
		)

		employees = append(employees, employee)

		// Build payslip.
		payslip := parser.BuildPayslip(
			row,
			headerMap,
			employee.ShalarthID,
			month,
			year,
		)

		payslips = append(payslips, payslip)
	}

	report.EmployeeCount = len(employees)
	report.SchoolCount = len(schools)
	report.MissingClusterCount = len(missingClusters)
	report.UnknownUDISECount = len(unknownUDISE)

	fmt.Println("\nMissing cluster mappings:")

	for udise, schoolName := range missingClusters {
		fmt.Printf("%s - %s\n", udise, schoolName)
	}

	fmt.Println("\nUnknown UDISE:")

	for udise := range unknownUDISE {
		fmt.Println(udise)
	}

	return employees, schools, payslips, report, nil
}
