package diagnostic

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// GetByMotorcycleAndBranch retrieves the most recent diagnostic for a motorcycle+branch combo
// Returns nil, nil if no diagnostic exists for that combination
func (r *repository) GetByMotorcycleAndBranch(ctx context.Context, motorcycleID, branchID string) (*domain.Diagnostic, error) {
	row := r.stmtGetByMotorcycleAndBranch.QueryRowContext(ctx, motorcycleID, branchID)

	var diag Diagnostic
	err := row.Scan(
		&diag.ID,
		&diag.MotorcycleID,
		&diag.BranchID,
		&diag.BranchName,
		&diag.Date,
		&diag.ProblemDescription,
		&diag.PossibleSolution,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No existing diagnostic for this combo
		}
		log.Error(logger.LogDiagnosticRepoGetByMotoBranchError, err)
		return nil, err
	}

	result := diag.ToDomain()
	return &result, nil
}
