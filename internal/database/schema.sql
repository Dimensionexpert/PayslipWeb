PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS clusters (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS schools (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    cluster_id  INTEGER NOT NULL REFERENCES clusters(id),
    udise       TEXT NOT NULL UNIQUE,
    school_name TEXT NOT NULL,
    ddo         TEXT UNIQUE,
    block       TEXT
);

CREATE TABLE IF NOT EXISTS employees (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    school_id            INTEGER NOT NULL REFERENCES schools(id),

    employee_serial_no   TEXT,
    shalarth_id          TEXT NOT NULL UNIQUE,
    name                 TEXT NOT NULL,
    gender               TEXT,
    designation          TEXT,

    gpf_no               TEXT,
    dcps_no              TEXT,
    pran_no              TEXT,
    pan                  TEXT UNIQUE,
    aadhaar              TEXT UNIQUE,

    mobile               TEXT,
    email                TEXT,

    ddo_bank_name        TEXT,
    ddo_bank_account     TEXT,
    ddo_bank_ifsc        TEXT,

    bank_name            TEXT,
    bank_account         TEXT,
    bank_ifsc            TEXT,
    branch_name          TEXT,

    pay_matrix            TEXT,

    CHECK (
        (gpf_no IS NOT NULL AND dcps_no IS NULL)
        OR
        (gpf_no IS NULL AND dcps_no IS NOT NULL)
    )
);

CREATE TABLE IF NOT EXISTS payslip_records (
    id                            INTEGER PRIMARY KEY AUTOINCREMENT,
    employee_id                   INTEGER NOT NULL REFERENCES employees(id),

    month                         INTEGER NOT NULL CHECK (month BETWEEN 1 AND 12),
    year                          INTEGER NOT NULL CHECK (year >= 2020),

    basic_pay                     REAL NOT NULL DEFAULT 0,
    da                            REAL NOT NULL DEFAULT 0,
    hra                           REAL NOT NULL DEFAULT 0,
    hra_arrear                    REAL NOT NULL DEFAULT 0,
    ta                            REAL NOT NULL DEFAULT 0,
    ta_arrear                     REAL NOT NULL DEFAULT 0,
    tribal_allowance              REAL NOT NULL DEFAULT 0,
    washing_allowance             REAL NOT NULL DEFAULT 0,
    da_arrears                    REAL NOT NULL DEFAULT 0,
    basic_arrears                 REAL NOT NULL DEFAULT 0,
    cla                           REAL NOT NULL DEFAULT 0,
    nps_empr_allow                REAL NOT NULL DEFAULT 0,

    total_pay                     REAL NOT NULL DEFAULT 0,
    fa                            REAL NOT NULL DEFAULT 0,
    gross_after_fa                REAL NOT NULL DEFAULT 0,

    gpf                           REAL NOT NULL DEFAULT 0,
    gpf_advance                   REAL NOT NULL DEFAULT 0,
    pt                            REAL NOT NULL DEFAULT 0,
    gis_zp                        REAL NOT NULL DEFAULT 0,
    gis_scout                     REAL NOT NULL DEFAULT 0,

    dcps_regular                  REAL NOT NULL DEFAULT 0,
    dcps_delayed                  REAL NOT NULL DEFAULT 0,
    dcps_pay_arrears_recovery     REAL NOT NULL DEFAULT 0,
    revenue_stamp                 REAL NOT NULL DEFAULT 0,
    dcps_da_arrears_recovery      REAL NOT NULL DEFAULT 0,
    group_accidental_policy       REAL NOT NULL DEFAULT 0,
    naa                           REAL NOT NULL DEFAULT 0,

    total_govt_deductions         REAL NOT NULL DEFAULT 0,
    gross_after_govt_deductions   REAL NOT NULL DEFAULT 0,

    nps_empr_contri               REAL NOT NULL DEFAULT 0,
    nps_emp_contri                REAL NOT NULL DEFAULT 0,
    nps_empr_contri_arr           REAL NOT NULL DEFAULT 0,
    nps_emp_contri_arr            REAL NOT NULL DEFAULT 0,
    nps_total                     REAL NOT NULL DEFAULT 0,
    gross_after_nps_deductions    REAL NOT NULL DEFAULT 0,

    income_tax                    REAL NOT NULL DEFAULT 0,
    coop_bank                     REAL NOT NULL DEFAULT 0,
    ngr_lic                       REAL NOT NULL DEFAULT 0,
    ngr_society_loan              REAL NOT NULL DEFAULT 0,
    ngr_misc                      REAL NOT NULL DEFAULT 0,
    ngr_other_recovery            REAL NOT NULL DEFAULT 0,
    ngr_rd                        REAL NOT NULL DEFAULT 0,
    ngr_other_deduction           REAL NOT NULL DEFAULT 0,
    ngr_total_deductions          REAL NOT NULL DEFAULT 0,

    employee_net_salary           REAL NOT NULL DEFAULT 0,

    remarks                       TEXT,

    UNIQUE (employee_id, month, year)
);


-- Indexes for foreign-key lookups

CREATE INDEX IF NOT EXISTS idx_schools_cluster_id
ON schools(cluster_id);

CREATE INDEX IF NOT EXISTS idx_employees_school_id
ON employees(school_id);

CREATE INDEX IF NOT EXISTS idx_payslip_employee_id
ON payslip_records(employee_id);
