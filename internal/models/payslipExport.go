package models

type PayslipExport struct {
	Cluster  Cluster
	School   School
	Employee Employee
	Payslip  PayslipRecord
}

type PayslipExportYear struct {
	Cluster  Cluster
	School   School
	Employee Employee
	Payslips []PayslipRecord
}
