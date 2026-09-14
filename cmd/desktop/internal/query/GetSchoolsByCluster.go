package query

import (
	"database/sql"

	"github.com/Dimensionexpert/payslip/cmd/desktop/internal/dto"
)

func GetSchoolsByCluster(
	db *sql.DB,
	cluster string,
) ([]dto.SchoolSummary, error) {
	rows, err := db.Query(`
		SELECT
			s.udise,
			s.school_name
		FROM schools s
		JOIN clusters c
			ON c.id = s.cluster_id
		WHERE c.name = ?
		ORDER BY s.school_name
	`, cluster)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schools []dto.SchoolSummary

	for rows.Next() {
		var school dto.SchoolSummary

		err := rows.Scan(
			&school.UDISECode,
			&school.Name,
		)
		if err != nil {
			return nil, err
		}

		schools = append(schools, school)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schools, nil
}
