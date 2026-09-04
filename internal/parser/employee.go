package parser

import (
	"github.com/Dimensionexpert/payslip/internal/excel"
	"github.com/Dimensionexpert/payslip/internal/models"
)

func BuildEmployee(
	row []string,
	headerMap map[string]int,
) models.Employee {
	return models.Employee{
		SerialNo:    excel.Get(row, headerMap, "S.R NO OF EMPL"),
		Name:        excel.Get(row, headerMap, "EMPLOYEE NAME"),
		ShalarthID:  excel.Get(row, headerMap, "SHALARTH ID"),
		Gender:      excel.Get(row, headerMap, "GENDER M/F"),
		Designation: excel.Get(row, headerMap, "DESIGNATION"),

		GPFNo:   excel.Get(row, headerMap, "GPF NO"),
		DCPSNo:  excel.Get(row, headerMap, "DCPS NO"),
		PRANNo:  excel.Get(row, headerMap, "PRAN NO"),
		PAN:     excel.Get(row, headerMap, "PAN NO"),
		Aadhaar: excel.Get(row, headerMap, "ADHAR NO"),

		Mobile: excel.Get(row, headerMap, "MOB NO"),
		Email:  excel.Get(row, headerMap, "EMAIL ID"),

		DDOBankName:    excel.Get(row, headerMap, "DDO BANK NAME"),
		DDOBankAccount: excel.Get(row, headerMap, "DDO BANK ACCOUNT NUMBER"),
		DDOBankIFSC:    excel.Get(row, headerMap, "DDO BANK IFSC CODE"),

		BankName:    excel.Get(row, headerMap, "BANK NAME"),
		BankAccount: excel.Get(row, headerMap, "BANK ACCOUNT NUMBER"),
		BankIFSC:    excel.Get(row, headerMap, "BANK IFSC CODE"),
		BranchName:  excel.Get(row, headerMap, "BRANCH NAME"),

		PayMatrix: excel.Get(row, headerMap, "PAY MATRIX"),

		UDISECode: excel.Get(row, headerMap, "SCHOOL UDISE CODE"),
	}
}
