package diagnostic

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// GetByID retrieves a diagnostic by its ID
func (r *repository) GetByID(ctx context.Context, diagnosticID string) (*domain.Diagnostic, error) {
	var diag Diagnostic

	err := r.stmtGetByID.QueryRowContext(ctx, diagnosticID).Scan(
		&diag.ID,
		&diag.MotorcycleID,
		&diag.BranchID,
		&diag.Date,
		&diag.ProblemDescription,
		&diag.PossibleSolution,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrDiagnosticNotFound
		}
		log.Error(logger.LogDiagnosticRepoGetByIDError, err)
		return nil, err
	}

	result := diag.ToDomain()
	return &result, nil
}
