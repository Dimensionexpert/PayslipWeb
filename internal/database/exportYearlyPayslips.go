package database

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/Dimensionexpert/payslip/internal/models"
)

func GetYearlyPayslips(
	conn *sql.DB,
	year int,
) ([]models.PayslipExportYear, error) {

	rows, err := conn.Query(`
		SELECT
			-- Cluster
			c.id,
			c.name,

			-- School
			s.id,
			s.school_name,
			s.udise,
			s.block,

			-- Employee
			e.id,
			-- e.employee_serial_no
			e.shalarth_id,
			e.name,
			e.gender,
			e.designation,
			e.pan,
			e.aadhaar,
			e.mobile,
			e.email,
			e.bank_name,
			e.bank_account,
			e.bank_ifsc,
			e.branch_name,

			-- Payslip
			p.id,
			p.employee_id,
			p.month,
			p.year,

			-- Earnings
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

			-- Government deductions
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

			-- NPS
			p.nps_empr_contri,
			p.nps_emp_contri,
			p.nps_empr_contri_arr,
			p.nps_emp_contri_arr,
			p.nps_total,
			p.gross_after_nps_deductions,

			-- Other deductions
			p.income_tax,
			p.coop_bank,
			p.ngr_lic,
			p.ngr_society_loan,
			p.ngr_misc,
			p.ngr_other_recovery,
			p.ngr_rd,
			p.ngr_other_deduction,
			p.ngr_total_deductions,

			-- Net
			p.employee_net_salary,
			p.remarks

		FROM payslip_records p

		JOIN employees e
			ON e.id = p.employee_id

		JOIN schools s
			ON s.id = e.school_id

		JOIN clusters c
			ON c.id = s.cluster_id

		WHERE p.year = ?

		ORDER BY
			e.shalarth_id,
			p.month;
	`, year)

	if err != nil {
		return nil, fmt.Errorf("querying yearly payslips: %w", err)
	}
	defer rows.Close()

	grouped := make(map[string]*models.PayslipExportYear)

	for rows.Next() {
		var export models.PayslipExportYear
		var payslip models.PayslipRecord

		err := rows.Scan(
			// Cluster
			&export.Cluster.ID,
			&export.Cluster.Name,

			// School
			&export.School.ID,
			&export.School.Name,
			&export.School.UDISECode,
			&export.School.Block,

			// Employee
			&export.Employee.ID,
			// &export.Employee.SerialNo,
			&export.Employee.ShalarthID,
			&export.Employee.Name,
			&export.Employee.Gender,
			&export.Employee.Designation,
			&export.Employee.PAN,
			&export.Employee.Aadhaar,
			&export.Employee.Mobile,
			&export.Employee.Email,
			&export.Employee.BankName,
			&export.Employee.BankAccount,
			&export.Employee.BankIFSC,
			&export.Employee.BranchName,

			// Payslip
			&payslip.ID,
			&payslip.EmployeeID,
			&payslip.Month,
			&payslip.Year,

			// Earnings
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

			// Government deductions
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

			// NPS
			&payslip.NPSEmprContri,
			&payslip.NPSEmpContri,
			&payslip.NPSEmprContriArr,
			&payslip.NPSEmpContriArr,
			&payslip.NPSTotal,
			&payslip.GrossAfterNPS,

			// Other deductions
			&payslip.IncomeTax,
			&payslip.CoopBank,
			&payslip.NGRLIC,
			&payslip.NGRSocietyLoan,
			&payslip.NGRMisc,
			&payslip.NGROtherRecovery,
			&payslip.NGRRD,
			&payslip.NGROtherDeduction,
			&payslip.NGRTotalDeduction,

			// Net
			&payslip.EmployeeNetSalary,
			&payslip.Remarks,
		)

		if err != nil {
			return nil, fmt.Errorf("scanning yearly payslip: %w", err)
		}

		shalarthID := export.Employee.ShalarthID

		existing, exists := grouped[shalarthID]

		if !exists {
			grouped[shalarthID] = &models.PayslipExportYear{
				Cluster:  export.Cluster,
				School:   export.School,
				Employee: export.Employee,
				Payslips: []models.PayslipRecord{payslip},
			}
			continue
		}

		existing.Payslips = append(existing.Payslips, payslip)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating yearly payslips: %w", err)
	}

	result := make([]models.PayslipExportYear, 0, len(grouped))

	for _, export := range grouped {
		result = append(result, *export)
	}

	// Map iteration order is not guaranteed.
	sort.Slice(result, func(i, j int) bool {
		return result[i].Employee.ShalarthID < result[j].Employee.ShalarthID
	})

	return result, nil
}
