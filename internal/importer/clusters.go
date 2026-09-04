package importer

// GetClusterMap builds a UDISE-to-cluster lookup from rows
// that have already been read from the cluster XLSX file.
func GetClusterMap(rows [][]string) (map[string]string, int) {
	clusterMap := make(map[string]string)
	duplicateUDISE := make(map[string]struct{})

	for i := 1; i < len(rows); i++ {
		row := rows[i]

		if len(row) < 2 {
			continue
		}

		clusterName := row[0]
		udiseCode := row[1]

		if clusterName == "" || udiseCode == "" {
			continue
		}

		if _, exists := clusterMap[udiseCode]; exists {
			duplicateUDISE[udiseCode] = struct{}{}
		}

		clusterMap[udiseCode] = clusterName
	}

	return clusterMap, len(duplicateUDISE)
}
