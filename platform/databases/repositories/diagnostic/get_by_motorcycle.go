package diagnostic

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// GetByMotorcycleID retrieves all diagnostics for a motorcycle, ordered by date descending
func (r *repository) GetByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.Diagnostic, error) {
	rows, err := r.stmtGetByMotorcycleID.QueryContext(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogDiagnosticRepoListByMotoError, err)
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var diagnostics []domain.Diagnostic
	for rows.Next() {
		var diag Diagnostic
		if err := rows.Scan(
			&diag.ID,
			&diag.MotorcycleID,
			&diag.BranchID,
			&diag.Date,
			&diag.ProblemDescription,
			&diag.PossibleSolution,
		); err != nil {
			log.Error(logger.LogDiagnosticRepoScanError, err)
			return nil, err
		}
		diagnostics = append(diagnostics, diag.ToDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return diagnostics, nil
}
