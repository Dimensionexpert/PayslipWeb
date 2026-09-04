package models

// Employee represents an employee.

type Employee struct {
	ID          int
	SchoolID    int
	SerialNo    string
	Name        string
	ShalarthID  string
	Gender      string
	Designation string

	GPFNo   string
	DCPSNo  string
	PRANNo  string
	PAN     string
	Aadhaar string

	Mobile string
	Email  string

	DDOBankName    string
	DDOBankAccount string
	DDOBankIFSC    string

	BankName    string
	BankAccount string
	BankIFSC    string
	BranchName  string

	PayMatrix string
	UDISECode string
}
