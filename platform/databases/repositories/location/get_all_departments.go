package location

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) GetAllDepartments(ctx context.Context) ([]domain.Department, error) {
	rows, err := r.stmtGetAllDepartments.QueryContext(ctx)
	if err != nil {
		log.Error(logger.LogLocationRepoGetDepartmentsError, "error", err)
		return nil, err
	}
	defer rows.Close()

	var departments []domain.Department
	for rows.Next() {
		var dept domain.Department
		if err := rows.Scan(&dept.ID, &dept.Name); err != nil {
			log.Error(logger.LogLocationRepoGetDepartmentsScanError, "error", err)
			continue
		}
		departments = append(departments, dept)
	}

	if err := rows.Err(); err != nil {
		log.Error(logger.LogLocationRepoGetDepartmentsIterError, "error", err)
		return nil, err
	}

	return departments, nil
}