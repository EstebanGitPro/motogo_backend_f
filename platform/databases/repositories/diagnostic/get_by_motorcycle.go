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

// GetEvidenceByDiagnosticID retrieves all evidence photos for a diagnostic
func (r *repository) GetEvidenceByDiagnosticID(ctx context.Context, diagnosticID string) ([]domain.DiagnosticEvidence, error) {
	rows, err := r.stmtGetEvidenceByDiagnosticID.QueryContext(ctx, diagnosticID)
	if err != nil {
		log.Error(logger.LogDiagnosticRepoListEvidenceError, err)
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var evidences []domain.DiagnosticEvidence
	for rows.Next() {
		var evid DiagnosticEvidence
		if err := rows.Scan(
			&evid.ID,
			&evid.DiagnosticID,
			&evid.ImageURL,
			&evid.Description,
			&evid.CreatedAt,
		); err != nil {
			log.Error(logger.LogDiagnosticRepoScanError, err)
			return nil, err
		}
		evidences = append(evidences, evid.EvidenceToDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return evidences, nil
}
