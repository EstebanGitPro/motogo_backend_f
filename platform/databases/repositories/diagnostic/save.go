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
		dbDiag.SentViaWhatsApp,
	)
	if err != nil {
		log.Error(logger.LogDiagnosticRepoSaveError, err)
		return domain.ErrDiagnosticCannotSave
	}

	return nil
}

// SaveEvidence inserts a new diagnostic evidence record
func (r *repository) SaveEvidence(ctx context.Context, tx output.Tx, evidence *domain.DiagnosticEvidence) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	dbEvid := EvidenceFromDomain(evidence)
	_, err := sqlTx.ExecContext(ctx, queryInsertEvidence,
		dbEvid.ID,
		dbEvid.DiagnosticID,
		dbEvid.ImageURL,
		dbEvid.Description,
		dbEvid.CreatedAt,
	)
	if err != nil {
		log.Error(logger.LogDiagnosticRepoSaveEvidenceError, err)
		return domain.ErrDiagnosticCannotSave
	}

	return nil
}
