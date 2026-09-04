package database

import (
	"database/sql"
	"fmt"

	"github.com/Dimensionexpert/payslip/internal/models"
)

func GetPayslip(
	db *sql.DB,
	shalarthID string,
	month int,
	year int,
) (models.PayslipExport, error) {

	var export models.PayslipExport

	err := db.QueryRow(`
		SELECT
			-- Cluster
			c.id,
			c.name,

			-- School
			s.id,
			s.udise,
			s.school_name,
			COALESCE(s.ddo, ''),
			COALESCE(s.block, ''),

			-- Employee
			e.id,
			e.school_id,
			COALESCE(e.employee_serial_no, ''),
			e.shalarth_id,
			e.name,
			COALESCE(e.gender, ''),
			COALESCE(e.designation, ''),
			COALESCE(e.gpf_no, ''),
			COALESCE(e.dcps_no, ''),
			COALESCE(e.pran_no, ''),
			COALESCE(e.pan, ''),
			COALESCE(e.aadhaar, ''),
			COALESCE(e.mobile, ''),
			COALESCE(e.email, ''),
			COALESCE(e.ddo_bank_name, ''),
			COALESCE(e.ddo_bank_account, ''),
			COALESCE(e.ddo_bank_ifsc, ''),
			COALESCE(e.bank_name, ''),
			COALESCE(e.bank_account, ''),
			COALESCE(e.bank_ifsc, ''),
			COALESCE(e.branch_name, ''),
			COALESCE(e.pay_matrix, ''),

			-- Payslip
			p.id,
			p.employee_id,
			p.month,
			p.year,
			p.basic_pay,
			p.da,
			p.hra,
			p.hra_arrear,
			p.ta,
			p.ta_arrear,
			p.tribal_allowance,
			p.washing_allowance,
			p.da_arrears,
			p.basic_arrears,
			p.cla,
			p.nps_empr_allow,
			p.total_pay,
			p.fa,
			p.gross_after_fa,
			p.gpf,
			p.gpf_advance,
			p.pt,
			p.gis_zp,
			p.gis_scout,
			p.dcps_regular,
			p.dcps_delayed,
			p.dcps_pay_arrears_recovery,
			p.revenue_stamp,
			p.dcps_da_arrears_recovery,
			p.group_accidental_policy,
			p.naa,
			p.total_govt_deductions,
			p.gross_after_govt_deductions,
			p.nps_empr_contri,
			p.nps_emp_contri,
			p.nps_empr_contri_arr,
			p.nps_emp_contri_arr,
			p.nps_total,
			p.gross_after_nps_deductions,
			p.income_tax,
			p.coop_bank,
			p.ngr_lic,
			p.ngr_society_loan,
			p.ngr_misc,
			p.ngr_other_recovery,
			p.ngr_rd,
			p.ngr_other_deduction,
			p.ngr_total_deductions,
			p.employee_net_salary,
			COALESCE(p.remarks, '')

		FROM payslip_records p

		JOIN employees e
			ON p.employee_id = e.id

		JOIN schools s
			ON e.school_id = s.id

		JOIN clusters c
			ON s.cluster_id = c.id

		WHERE e.shalarth_id = ?
		  AND p.month = ?
		  AND p.year = ?
	`,
		shalarthID,
		month,
		year,
	).Scan(
		// Cluster
		&export.Cluster.ID,
		&export.Cluster.Name,

		// School
		&export.School.ID,
		&export.School.UDISECode,
		&export.School.Name,
		&export.School.DDOCode,
		&export.School.Block,

		// Employee
		&export.Employee.ID,
		&export.Employee.SchoolID,
		&export.Employee.SerialNo,
		&export.Employee.ShalarthID,
		&export.Employee.Name,
		&export.Employee.Gender,
		&export.Employee.Designation,
		&export.Employee.GPFNo,
		&export.Employee.DCPSNo,
		&export.Employee.PRANNo,
		&export.Employee.PAN,
		&export.Employee.Aadhaar,
		&export.Employee.Mobile,
		&export.Employee.Email,
		&export.Employee.DDOBankName,
		&export.Employee.DDOBankAccount,
		&export.Employee.DDOBankIFSC,
		&export.Employee.BankName,
		&export.Employee.BankAccount,
		&export.Employee.BankIFSC,
		&export.Employee.BranchName,
		&export.Employee.PayMatrix,

		// Payslip
		&export.Payslip.ID,
		&export.Payslip.EmployeeID,
		&export.Payslip.Month,
		&export.Payslip.Year,
		&export.Payslip.BasicPay,
		&export.Payslip.DA,
		&export.Payslip.HRA,
		&export.Payslip.HRAArrear,
		&export.Payslip.TA,
		&export.Payslip.TAArrear,
		&export.Payslip.TribalAllowance,
		&export.Payslip.WashingAllowance,
		&export.Payslip.DAArrears,
		&export.Payslip.BasicArrears,
		&export.Payslip.CLA,
		&export.Payslip.NPSEmprAllow,
		&export.Payslip.TotalPay,
		&export.Payslip.FA,
		&export.Payslip.GrossAfterFA,
		&export.Payslip.GPF,
		&export.Payslip.GPFAdvance,
		&export.Payslip.PT,
		&export.Payslip.GISZP,
		&export.Payslip.GISScout,
		&export.Payslip.DCPSRegular,
		&export.Payslip.DCPSDelayed,
		&export.Payslip.DCPSPayArrears,
		&export.Payslip.RevenueStamp,
		&export.Payslip.DCPSDAArrears,
		&export.Payslip.GroupAccidentalPolicy,
		&export.Payslip.NAA,
		&export.Payslip.TotalGovtDeductions,
		&export.Payslip.GrossAfterGovtDeductions,
		&export.Payslip.NPSEmprContri,
		&export.Payslip.NPSEmpContri,
		&export.Payslip.NPSEmprContriArr,
		&export.Payslip.NPSEmpContriArr,
		&export.Payslip.NPSTotal,
		&export.Payslip.GrossAfterNPS,
		&export.Payslip.IncomeTax,
		&export.Payslip.CoopBank,
		&export.Payslip.NGRLIC,
		&export.Payslip.NGRSocietyLoan,
		&export.Payslip.NGRMisc,
		&export.Payslip.NGROtherRecovery,
		&export.Payslip.NGRRD,
		&export.Payslip.NGROtherDeduction,
		&export.Payslip.NGRTotalDeduction,
		&export.Payslip.EmployeeNetSalary,
		&export.Payslip.Remarks,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return export, fmt.Errorf(
				"payslip not found: Shalarth ID %s, %02d/%d",
				shalarthID,
				month,
				year,
			)
		}

		return export, err
	}

	return export, nil
}
