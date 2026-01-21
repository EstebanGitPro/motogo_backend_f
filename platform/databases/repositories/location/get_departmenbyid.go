package location

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// GetDepartmentByID retrieves a department by its ID
func (r *repository) GetDepartmentByID(ctx context.Context, departmentID string) (*domain.Department, error) {
	var dept domain.Department
	err := r.stmtGetDepartmentByID.QueryRowContext(ctx, departmentID).Scan(&dept.ID, &dept.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrDepartmentNotFound
		}
		log.Error(logger.LogLocationRepoGetDeptByIDError, "error", err, "department_id", departmentID)
		return nil, err
	}
	return &dept, nil
}
