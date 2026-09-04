package parser

import (
	"github.com/Dimensionexpert/payslip/internal/excel"
	"github.com/Dimensionexpert/payslip/internal/models"
)

func BuildSchool(
	row []string,
	headerMap map[string]int,
	clusterMap map[string]string,
) models.School {
	udiseCode := excel.Get(row, headerMap, "SCHOOL UDISE CODE")

	return models.School{
		UDISECode: udiseCode,
		Name:      excel.Get(row, headerMap, "NAME OF SCHOOL"),
		DDOCode:   excel.Get(row, headerMap, "SCHOOL SHALARTH DDO CODE"),
		Block:     excel.Get(row, headerMap, "BLOCK / TALUKA"),
		Cluster:   clusterMap[udiseCode],
	}
}
