package parser

import (
	"strconv"
	"strings"

	"github.com/Dimensionexpert/payslip/internal/excel"
	"github.com/Dimensionexpert/payslip/internal/models"
)

func BuildPayslip(
	row []string,
	headerMap map[string]int,
	shalarthID string,
	month int,
	year int,

) models.PayslipRecord {
	return models.PayslipRecord{
		Month: month,
		Year:  year,

		BasicPay:         parseMoney(excel.Get(row, headerMap, "BASIC PAY")),
		DA:               parseMoney(excel.Get(row, headerMap, "D.A")),
		HRA:              parseMoney(excel.Get(row, headerMap, "HRA")),
		HRAArrear:        parseMoney(excel.Get(row, headerMap, "HRA ARREAR")),
		TA:               parseMoney(excel.Get(row, headerMap, "T.A")),
		TAArrear:         parseMoney(excel.Get(row, headerMap, "T.A ARREAR")),
		TribalAllowance:  parseMoney(excel.Get(row, headerMap, "TRIBAL ALLOWANCE")),
		WashingAllowance: parseMoney(excel.Get(row, headerMap, "WASHING ALLOWANCE")),
		DAArrears:        parseMoney(excel.Get(row, headerMap, "DA ARREARS")),
		BasicArrears:     parseMoney(excel.Get(row, headerMap, "BASIC ARREARS")),
		CLA:              parseMoney(excel.Get(row, headerMap, "CLA")),
		NPSEmprAllow:     parseMoney(excel.Get(row, headerMap, "NPS EMPR ALLOW")),

		TotalPay:     parseMoney(excel.Get(row, headerMap, "TOTAL PAY")),
		FA:           parseMoney(excel.Get(row, headerMap, "F A")),
		GrossAfterFA: parseMoney(excel.Get(row, headerMap, "GROSS AFTER DEDUCTING FA")),

		GPF:                   parseMoney(excel.Get(row, headerMap, "GPF")),
		GPFAdvance:            parseMoney(excel.Get(row, headerMap, "GPF ADV")),
		PT:                    parseMoney(excel.Get(row, headerMap, "PT")),
		GISZP:                 parseMoney(excel.Get(row, headerMap, "GIS(ZP)")),
		GISScout:              parseMoney(excel.Get(row, headerMap, "GIS SCOUT")),
		DCPSRegular:           parseMoney(excel.Get(row, headerMap, "DCPS REGULAR")),
		DCPSDelayed:           parseMoney(excel.Get(row, headerMap, "DCPS DELAYED")),
		DCPSPayArrears:        parseMoney(excel.Get(row, headerMap, "DCPS PAY ARREARS RECOVERY")),
		RevenueStamp:          parseMoney(excel.Get(row, headerMap, "REVENUE STAMP")),
		DCPSDAArrears:         parseMoney(excel.Get(row, headerMap, "DCPS DA ARREARS RECOVERY")),
		GroupAccidentalPolicy: parseMoney(excel.Get(row, headerMap, "GROUP ACCIDENTAL POLICY")),
		NAA:                   parseMoney(excel.Get(row, headerMap, "NAA")),
		TotalGovtDeductions:   parseMoney(excel.Get(row, headerMap, "TOTAL GOVT DEDUCTIONS")),
		GrossAfterGovtDeductions: parseMoney(
			excel.Get(row, headerMap, "GROSS PAYMENT AFTER GOVT DEDUCTIONS"),
		),

		NPSEmprContri:    parseMoney(excel.Get(row, headerMap, "NPS EMPR CONTRI")),
		NPSEmpContri:     parseMoney(excel.Get(row, headerMap, "NPS EMP CONTRI")),
		NPSEmprContriArr: parseMoney(excel.Get(row, headerMap, "NPS EMPR CONTRI ARR")),
		NPSEmpContriArr:  parseMoney(excel.Get(row, headerMap, "NPS EMP CONTRI ARR")),
		NPSTotal:         parseMoney(excel.Get(row, headerMap, "NPS TOTAL")),
		GrossAfterNPS: parseMoney(
			excel.Get(row, headerMap, "GROSS PAYMENT AFTER NPS DEDUCTIONS"),
		),

		IncomeTax:         parseMoney(excel.Get(row, headerMap, "INCOME TAX")),
		CoopBank:          parseMoney(excel.Get(row, headerMap, "CO-OP BANK")),
		NGRLIC:            parseMoney(excel.Get(row, headerMap, "NGR(LIC)")),
		NGRSocietyLoan:    parseMoney(excel.Get(row, headerMap, "NGR(SOCIETY LOAN)")),
		NGRMisc:           parseMoney(excel.Get(row, headerMap, "NGR(MISC)")),
		NGROtherRecovery:  parseMoney(excel.Get(row, headerMap, "NGR(OTHER RECOVERY)")),
		NGRRD:             parseMoney(excel.Get(row, headerMap, "NGR(RD)")),
		NGROtherDeduction: parseMoney(excel.Get(row, headerMap, "NGR(OTHER DEDUCTION)")),
		NGRTotalDeduction: parseMoney(excel.Get(row, headerMap, "NGR(TOTAL DEDUCTIONS)")),

		EmployeeNetSalary: parseMoney(
			excel.Get(row, headerMap, "EMPLOYEE NET SALARY"),
		),

		Remarks:    excel.Get(row, headerMap, "REMARKS"),
		ShalarthID: shalarthID,
	}
}

func parseMoney(value string) float64 {
	value = strings.TrimSpace(value)

	if value == "" {
		return 0
	}

	value = strings.ReplaceAll(value, ",", "")
	value = strings.ReplaceAll(value, "₹", "")

	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}

	return number
}
