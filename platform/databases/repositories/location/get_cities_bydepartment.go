package location

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) GetCitiesByDepartment(ctx context.Context, departmentID string) ([]domain.City, error) {
	rows, err := r.stmtGetCitiesByDepartment.QueryContext(ctx, departmentID)
	if err != nil {
		log.Error(logger.LogLocationRepoGetCitiesError, "error", err, "department_id", departmentID)
		return nil, err
	}
	defer rows.Close()

	var cities []domain.City
	for rows.Next() {
		var city domain.City
		if err := rows.Scan(&city.ID, &city.Name, &city.DepartmentID); err != nil {
			log.Error(logger.LogLocationRepoGetCitiesScanError, "error", err)
			continue
		}
		cities = append(cities, city)
	}

	if err := rows.Err(); err != nil {
		log.Error(logger.LogLocationRepoGetCitiesIterError, "error", err)
		return nil, err
	}

	return cities, nil
}
