package services

import (
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// NewDiagnostic creates a new Diagnostic with generated ID and current timestamp
func NewDiagnostic(motorcycleID, branchID string, problemDescription *string) *domain.Diagnostic {
	d := &domain.Diagnostic{
		MotorcycleID:       motorcycleID,
		BranchID:           branchID,
		Date:               time.Now(),
		ProblemDescription: problemDescription,
	}
	d.SetID()
	return d
}

// NewDiagnosticEvidence creates a new DiagnosticEvidence with generated ID and timestamp
func NewDiagnosticEvidence(diagnosticID, imageURL string, description *string) *domain.DiagnosticEvidence {
	e := &domain.DiagnosticEvidence{
		DiagnosticID: diagnosticID,
		ImageURL:     imageURL,
		Description:  description,
		CreatedAt:    time.Now(),
	}
	e.SetID()
	return e
}

// RefreshDiagnostic updates an existing diagnostic with a new problem description and resets
// its timestamp. Used when the same motorcycle+branch combination already has a diagnostic (UPSERT).
func RefreshDiagnostic(existing *domain.Diagnostic, problemDescription *string) {
	existing.ProblemDescription = problemDescription
	existing.Date = time.Now()
}
