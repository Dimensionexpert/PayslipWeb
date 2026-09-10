package importer

import "fmt"

// ImportReport holds the results of an import operation.
type ImportReport struct {
	EmployeeCount       int
	SchoolCount         int
	ClusterMappingCount int
	MissingClusterCount int
	UnknownUDISECount   int
	DuplicateUDISECount int

	MissingClusters []string
	UnknownUDISE    []string
	DuplicateUDISE  []string
}

func (r ImportReport) Print() {

	fmt.Println("========== IMPORT REPORT ==========")
	fmt.Printf("Employees:          %d\n", r.EmployeeCount)
	fmt.Printf("Schools:            %d\n", r.SchoolCount)
	fmt.Printf("Cluster mappings:   %d\n", r.ClusterMappingCount)
	fmt.Printf("Missing clusters:   %d\n", r.MissingClusterCount)
	fmt.Printf("Unknown UDISE:      %d\n", r.UnknownUDISECount)
	fmt.Printf("Duplicate UDISE:    %d\n", r.DuplicateUDISECount)

	if len(r.MissingClusters) > 0 {
		fmt.Println()
		fmt.Println("Missing Cluster Mappings:")
		for _, value := range r.MissingClusters {
			fmt.Println(" -", value)
		}
	}

	if len(r.UnknownUDISE) > 0 {
		fmt.Println()
		fmt.Println("Unknown UDISE:")
		for _, value := range r.UnknownUDISE {
			fmt.Println(" -", value)
		}
	}

	if len(r.DuplicateUDISE) > 0 {
		fmt.Println()
		fmt.Println("Duplicate UDISE:")
		for _, value := range r.DuplicateUDISE {
			fmt.Println(" -", value)
		}
	}

	fmt.Println("===================================")
}
