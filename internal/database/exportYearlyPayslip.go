package database

import (
	"database/sql"
	"fmt"

	"github.com/Dimensionexpert/payslip/internal/models"
)

func GetYearlyPayslip(
	db *sql.DB,
	shalarthID string,
	year int,
) (models.PayslipExportYear, error) {

	rows, err := db.Query(`
		SELECT
			c.id,
			c.name,

			s.id,
			s.udise,
			s.school_name,
			COALESCE(s.ddo, ''),
			COALESCE(s.block, ''),

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
		  AND p.year = ?

		ORDER BY p.month
	`, shalarthID, year)

	if err != nil {
		return models.PayslipExportYear{}, err
	}
	defer rows.Close()

	var export models.PayslipExportYear

	for rows.Next() {
		var payslip models.PayslipRecord

		err := rows.Scan(
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
			&payslip.ID,
			&payslip.EmployeeID,
			&payslip.Month,
			&payslip.Year,
			&payslip.BasicPay,
			&payslip.DA,
			&payslip.HRA,
			&payslip.HRAArrear,
			&payslip.TA,
			&payslip.TAArrear,
			&payslip.TribalAllowance,
			&payslip.WashingAllowance,
			&payslip.DAArrears,
			&payslip.BasicArrears,
			&payslip.CLA,
			&payslip.NPSEmprAllow,
			&payslip.TotalPay,
			&payslip.FA,
			&payslip.GrossAfterFA,
			&payslip.GPF,
			&payslip.GPFAdvance,
			&payslip.PT,
			&payslip.GISZP,
			&payslip.GISScout,
			&payslip.DCPSRegular,
			&payslip.DCPSDelayed,
			&payslip.DCPSPayArrears,
			&payslip.RevenueStamp,
			&payslip.DCPSDAArrears,
			&payslip.GroupAccidentalPolicy,
			&payslip.NAA,
			&payslip.TotalGovtDeductions,
			&payslip.GrossAfterGovtDeductions,
			&payslip.NPSEmprContri,
			&payslip.NPSEmpContri,
			&payslip.NPSEmprContriArr,
			&payslip.NPSEmpContriArr,
			&payslip.NPSTotal,
			&payslip.GrossAfterNPS,
			&payslip.IncomeTax,
			&payslip.CoopBank,
			&payslip.NGRLIC,
			&payslip.NGRSocietyLoan,
			&payslip.NGRMisc,
			&payslip.NGROtherRecovery,
			&payslip.NGRRD,
			&payslip.NGROtherDeduction,
			&payslip.NGRTotalDeduction,
			&payslip.EmployeeNetSalary,
			&payslip.Remarks,
		)

		if err != nil {
			return models.PayslipExportYear{}, fmt.Errorf(
				"scan yearly payslip row: %w",
				err,
			)
		}

		export.Payslips = append(export.Payslips, payslip)
	}

	if err := rows.Err(); err != nil {
		return models.PayslipExportYear{}, fmt.Errorf(
			"iterate yearly payslip rows: %w",
			err,
		)
	}

	return export, nil
}
