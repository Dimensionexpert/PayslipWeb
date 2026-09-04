package models

// PayslipRecord represents an employee's salary record for one month.
type PayslipRecord struct {
	Month int
	Year  int

	ID         int
	EmployeeID int

	BasicPay         float64
	DA               float64
	HRA              float64
	HRAArrear        float64
	TA               float64
	TAArrear         float64
	TribalAllowance  float64
	WashingAllowance float64
	DAArrears        float64
	BasicArrears     float64
	CLA              float64
	NPSEmprAllow     float64

	TotalPay              float64
	FA                    float64
	GrossAfterFA          float64
	GPF                   float64
	GPFAdvance            float64
	PT                    float64
	GISZP                 float64
	GISScout              float64
	DCPSRegular           float64
	DCPSDelayed           float64
	DCPSPayArrears        float64
	RevenueStamp          float64
	DCPSDAArrears         float64
	GroupAccidentalPolicy float64
	NAA                   float64

	TotalGovtDeductions      float64
	GrossAfterGovtDeductions float64

	NPSEmprContri    float64
	NPSEmpContri     float64
	NPSEmprContriArr float64
	NPSEmpContriArr  float64
	NPSTotal         float64
	GrossAfterNPS    float64

	IncomeTax         float64
	CoopBank          float64
	NGRLIC            float64
	NGRSocietyLoan    float64
	NGRMisc           float64
	NGROtherRecovery  float64
	NGRRD             float64
	NGROtherDeduction float64
	NGRTotalDeduction float64

	EmployeeNetSalary float64
	Remarks           string

	UDISECode  string
	ShalarthID string
}
