package diagnostic

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// Save inserts a new diagnostic record
func (r *repository) Save(ctx context.Context, tx output.Tx, diagnostic *domain.Diagnostic) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	dbDiag := FromDomain(diagnostic)
	_, err := sqlTx.ExecContext(ctx, queryInsert,
		dbDiag.ID,
		dbDiag.MotorcycleID,
		dbDiag.BranchID,
		dbDiag.Date,
		dbDiag.ProblemDescription,
		dbDiag.PossibleSolution,
	)
	if err != nil {
		log.Error(logger.LogDiagnosticRepoSaveError, err)
		return domain.ErrDiagnosticCannotSave
	}

	return nil
}
