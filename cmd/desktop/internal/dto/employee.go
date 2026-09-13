package dto

// dto means data transfer object

type EmployeeSummary struct {
	ShalarthID  string `json:"shalarthId"`
	Name        string `json:"name"`
	Designation string `json:"designation"`
	SchoolName  string `json:"schoolName"`
	UDISECode   string `json:"udiseCode"`
}
