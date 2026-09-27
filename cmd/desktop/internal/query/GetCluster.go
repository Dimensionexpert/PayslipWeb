package query

import (
	"database/sql"

	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/dto"
)

func GetClusters(db *sql.DB) ([]dto.ClusterSummary, error) {
	rows, err := db.Query(`
		SELECT
			id,
			name
		FROM clusters
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clusters []dto.ClusterSummary

	for rows.Next() {
		var cluster dto.ClusterSummary

		err := rows.Scan(
			&cluster.ID,
			&cluster.Name,
		)
		if err != nil {
			return nil, err
		}

		clusters = append(clusters, cluster)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clusters, nil
}
