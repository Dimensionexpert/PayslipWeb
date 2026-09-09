package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Dimensionexpert/payslip/internal/models"
)

// getEmployeeIDByShalarthID returns the database ID of an employee.
func getEmployeeIDByShalarthID(db DBTX, shalarthID string) (int, error) {
	var employeeID int

	err := db.QueryRow(`
		SELECT id
		FROM employees
		WHERE shalarth_id = ?
	`, shalarthID).Scan(&employeeID)

	return employeeID, err
}

// nullable returns the value as-is if it's not empty, otherwise returns nil.
func nullable(value string) any {
	value = strings.TrimSpace(value)

	if value == "" {
		return nil
	}

	return value
}

// DBTX represents anything that supports QueryRow.
// Both *sql.DB and *sql.Tx satisfy this interface.
type DBTX interface {
	QueryRow(query string, args ...any) *sql.Row
}

// InsertClusters inserts unique clusters from the cluster map.
func InsertClusters(db *sql.DB, clusterMap map[string]string) error {
	seen := make(map[string]struct{})

	for _, clusterName := range clusterMap {
		if _, exists := seen[clusterName]; exists {
			continue
		}

		seen[clusterName] = struct{}{}

		_, err := db.Exec(`
			INSERT INTO clusters (name)
			VALUES (?)
			ON CONFLICT(name) DO NOTHING
		`, clusterName)
		if err != nil {
			return err
		}
	}

	return nil
}

// getClusterID returns the ID of a cluster by name.
func getClusterID(db *sql.DB, clusterName string) (int, error) {
	var id int

	err := db.QueryRow(`
		SELECT id
		FROM clusters
		WHERE name = ?
	`, clusterName).Scan(&id)

	return id, err
}

// getSchoolIDByUDISE returns the ID of a school by UDISE code.
func getSchoolIDByUDISE(db DBTX, udise string) (int, error) {
	var schoolID int

	err := db.QueryRow(`
		SELECT id
		FROM schools
		WHERE udise = ?
	`, udise).Scan(&schoolID)

	return schoolID, err
}

// InsertSchools inserts the given schools into the database.
func InsertSchools(db *sql.DB, schools map[string]models.School) error {
	for _, school := range schools {
		clusterID, err := getClusterID(db, school.Cluster)
		if err != nil {
			return err
		}

		_, err = db.Exec(`
			INSERT INTO schools (
				cluster_id,
				udise,
				school_name,
				ddo,
				block
			)
			VALUES (?, ?, ?, ?, ?)
			ON CONFLICT(udise) DO NOTHING
		`,
			clusterID,
			school.UDISECode,
			school.Name,
			school.DDOCode,
			school.Block,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// InsertEmployees inserts the given employees into the database.
func InsertEmployees(db *sql.DB, employees []models.Employee) error {
	seenAadhaar := make(map[string]string)

	for _, employee := range employees {
		aadhaar := strings.TrimSpace(employee.Aadhaar)

		if aadhaar == "" {
			continue
		}

		if previous, exists := seenAadhaar[aadhaar]; exists {
			fmt.Printf(
				"Duplicate Aadhaar: %s | %s | %s\n",
				aadhaar,
				previous,
				employee.Name,
			)
		} else {
			seenAadhaar[aadhaar] = employee.Name
		}

		// existing insert...
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO employees (
			school_id,
			shalarth_id,
			name,
			gender,
			designation,
			gpf_no,
			dcps_no,
			pran_no,
			pan,
			aadhaar,
			mobile,
			email,
			ddo_bank_name,
			ddo_bank_account,
			ddo_bank_ifsc,
			bank_name,
			bank_account,
			bank_ifsc,
			branch_name,
			pay_matrix
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(shalarth_id) DO NOTHING
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, employee := range employees {
		schoolID, err := getSchoolIDByUDISE(tx, employee.UDISECode)
		if err != nil {
			return err
		}

		if employee.GPFNo == "" && employee.DCPSNo == "" {
			fmt.Printf("Invalid employee: %+v\n", employee)
		}

		_, err = stmt.Exec(
			schoolID,
			employee.ShalarthID,
			employee.Name,
			employee.Gender,
			employee.Designation,
			nullable(employee.GPFNo),
			nullable(employee.DCPSNo),
			nullable(employee.PRANNo),
			nullable(employee.PAN),
			nullable(employee.Aadhaar),
			employee.Mobile,
			employee.Email,
			nullable(employee.DDOBankName),
			nullable(employee.DDOBankAccount),
			nullable(employee.DDOBankIFSC),
			nullable(employee.BankName),
			nullable(employee.BankAccount),
			nullable(employee.BankIFSC),
			nullable(employee.BranchName),
			employee.PayMatrix,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// InsertPayslipRecords inserts payslip records into the database.
func InsertPayslipRecords(
	db *sql.DB,
	records []models.PayslipRecord,
) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO payslip_records (
			employee_id,
			month,
			year,

			basic_pay,
			da,
			hra,
			hra_arrear,
			ta,
			ta_arrear,
			tribal_allowance,
			washing_allowance,
			da_arrears,
			basic_arrears,
			cla,
			nps_empr_allow,

			total_pay,
			fa,
			gross_after_fa,

			gpf,
			gpf_advance,
			pt,
			gis_zp,
			gis_scout,
			dcps_regular,
			dcps_delayed,
			dcps_pay_arrears_recovery,
			revenue_stamp,
			dcps_da_arrears_recovery,
			group_accidental_policy,
			naa,

			total_govt_deductions,
			gross_after_govt_deductions,

			nps_empr_contri,
			nps_emp_contri,
			nps_empr_contri_arr,
			nps_emp_contri_arr,
			nps_total,
			gross_after_nps_deductions,

			income_tax,
			coop_bank,
			ngr_lic,
			ngr_society_loan,
			ngr_misc,
			ngr_other_recovery,
			ngr_rd,
			ngr_other_deduction,
			ngr_total_deductions,

			employee_net_salary,
			remarks
		)
		VALUES (
			?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?,
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?
		)
		ON CONFLICT(employee_id, month, year)
		DO UPDATE SET
    basic_pay = excluded.basic_pay,
    da = excluded.da,
    hra = excluded.hra,
    hra_arrear = excluded.hra_arrear,
    ta = excluded.ta,
    ta_arrear = excluded.ta_arrear,
    tribal_allowance = excluded.tribal_allowance,
    washing_allowance = excluded.washing_allowance,
    da_arrears = excluded.da_arrears,
    basic_arrears = excluded.basic_arrears,
    cla = excluded.cla,
    nps_empr_allow = excluded.nps_empr_allow,
    total_pay = excluded.total_pay,
    fa = excluded.fa,
    gross_after_fa = excluded.gross_after_fa,
    gpf = excluded.gpf,
    gpf_advance = excluded.gpf_advance,
    pt = excluded.pt,
    gis_zp = excluded.gis_zp,
    gis_scout = excluded.gis_scout,
    dcps_regular = excluded.dcps_regular,
    dcps_delayed = excluded.dcps_delayed,
    dcps_pay_arrears_recovery = excluded.dcps_pay_arrears_recovery,
    revenue_stamp = excluded.revenue_stamp,
    dcps_da_arrears_recovery = excluded.dcps_da_arrears_recovery,
    group_accidental_policy = excluded.group_accidental_policy,
    naa = excluded.naa,
    total_govt_deductions = excluded.total_govt_deductions,
    gross_after_govt_deductions = excluded.gross_after_govt_deductions,
    nps_empr_contri = excluded.nps_empr_contri,
    nps_emp_contri = excluded.nps_emp_contri,
    nps_empr_contri_arr = excluded.nps_empr_contri_arr,
    nps_emp_contri_arr = excluded.nps_emp_contri_arr,
    nps_total = excluded.nps_total,
    gross_after_nps_deductions = excluded.gross_after_nps_deductions,
    income_tax = excluded.income_tax,
    coop_bank = excluded.coop_bank,
    ngr_lic = excluded.ngr_lic,
    ngr_society_loan = excluded.ngr_society_loan,
    ngr_misc = excluded.ngr_misc,
    ngr_other_recovery = excluded.ngr_other_recovery,
    ngr_rd = excluded.ngr_rd,
    ngr_other_deduction = excluded.ngr_other_deduction,
    ngr_total_deductions = excluded.ngr_total_deductions,
    employee_net_salary = excluded.employee_net_salary,
    remarks = excluded.remarks;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, record := range records {
		employeeID, err := getEmployeeIDByShalarthID(tx, record.ShalarthID)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(
			employeeID,
			record.Month,
			record.Year,

			record.BasicPay,
			record.DA,
			record.HRA,
			record.HRAArrear,
			record.TA,
			record.TAArrear,
			record.TribalAllowance,
			record.WashingAllowance,
			record.DAArrears,
			record.BasicArrears,
			record.CLA,
			record.NPSEmprAllow,

			record.TotalPay,
			record.FA,
			record.GrossAfterFA,

			record.GPF,
			record.GPFAdvance,
			record.PT,
			record.GISZP,
			record.GISScout,
			record.DCPSRegular,
			record.DCPSDelayed,
			record.DCPSPayArrears,
			record.RevenueStamp,
			record.DCPSDAArrears,
			record.GroupAccidentalPolicy,
			record.NAA,

			record.TotalGovtDeductions,
			record.GrossAfterGovtDeductions,

			record.NPSEmprContri,
			record.NPSEmpContri,
			record.NPSEmprContriArr,
			record.NPSEmpContriArr,
			record.NPSTotal,
			record.GrossAfterNPS,

			record.IncomeTax,
			record.CoopBank,
			record.NGRLIC,
			record.NGRSocietyLoan,
			record.NGRMisc,
			record.NGROtherRecovery,
			record.NGRRD,
			record.NGROtherDeduction,
			record.NGRTotalDeduction,

			record.EmployeeNetSalary,
			record.Remarks,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
