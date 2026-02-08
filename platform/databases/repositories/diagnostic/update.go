package diagnostic

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// Update modifies an existing diagnostic record
func (r *repository) Update(ctx context.Context, tx output.Tx, diagnostic *domain.Diagnostic) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	dbDiag := FromDomain(diagnostic)
	_, err := sqlTx.ExecContext(ctx, queryUpdate,
		dbDiag.ProblemDescription,
		dbDiag.PossibleSolution,
		dbDiag.LaborQuote,
		dbDiag.PartsQuote,
		dbDiag.EstimatedTime,
		dbDiag.SentViaWhatsApp,
		dbDiag.ID,
	)
	if err != nil {
		log.Error(logger.LogDiagnosticRepoUpdateError, err)
		return domain.ErrDiagnosticCannotUpdate
	}

	return nil
}
