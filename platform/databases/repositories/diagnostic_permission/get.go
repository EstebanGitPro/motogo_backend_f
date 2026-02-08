package diagnostic_permission

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// GetByMotorcycleAndBranch retrieves an active permission for a specific motorcycle-branch pair
func (r *repository) GetByMotorcycleAndBranch(ctx context.Context, motorcycleID, branchID string) (*domain.DiagnosticPermission, error) {
	var dbPerm DiagnosticPermission
	err := r.stmtGetByMotorcycleAndBranch.QueryRowContext(ctx, motorcycleID, branchID).Scan(
		&dbPerm.ID,
		&dbPerm.MotorcycleID,
		&dbPerm.BranchID,
		&dbPerm.Active,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrPermissionNotFound
		}
		log.Error(logger.LogDiagPermRepoGetError, err)
		return nil, err
	}

	result := dbPerm.ToDomain()
	return &result, nil
}

// GetByMotorcycleID retrieves all active permissions for a motorcycle
func (r *repository) GetByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.DiagnosticPermission, error) {
	rows, err := r.stmtGetByMotorcycleID.QueryContext(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogDiagPermRepoListError, err)
		return nil, err
	}
	defer rows.Close()

	var permissions []domain.DiagnosticPermission
	for rows.Next() {
		var dbPerm DiagnosticPermission
		err := rows.Scan(
			&dbPerm.ID,
			&dbPerm.MotorcycleID,
			&dbPerm.BranchID,
			&dbPerm.Active,
		)
		if err != nil {
			log.Error(logger.LogDiagPermRepoScanError, err)
			return nil, err
		}
		permissions = append(permissions, dbPerm.ToDomain())
	}

	return permissions, nil
}
