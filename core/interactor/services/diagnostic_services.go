package services

import (
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/tools/utils"
)

// NewDiagnostic creates a new Diagnostic with generated ID and current timestamp
func NewDiagnostic(motorcycleID, branchID string, problemDescription *string) *domain.Diagnostic {
	return &domain.Diagnostic{
		ID:                 utils.Generate(),
		MotorcycleID:       motorcycleID,
		BranchID:           branchID,
		Date:               time.Now(),
		ProblemDescription: problemDescription,
		SentViaWhatsApp:    false,
	}
}

// NewDiagnosticEvidence creates a new DiagnosticEvidence with generated ID and timestamp
func NewDiagnosticEvidence(diagnosticID, imageURL string, description *string) *domain.DiagnosticEvidence {
	return &domain.DiagnosticEvidence{
		ID:           utils.Generate(),
		DiagnosticID: diagnosticID,
		ImageURL:     imageURL,
		Description:  description,
		CreatedAt:    time.Now(),
	}
}

// RefreshDiagnostic updates an existing diagnostic with a new problem description and resets
// its timestamp. Used when the same motorcycle+branch combination already has a diagnostic (UPSERT).
func RefreshDiagnostic(existing *domain.Diagnostic, problemDescription *string) {
	existing.ProblemDescription = problemDescription
	existing.Date = time.Now()
	existing.SentViaWhatsApp = false
}
