package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// EvidenceInteractor handles motorcycle evidence-related use cases (HU16-19)
type EvidenceInteractor struct {
	evidenceSvc input.EvidenceService
}

// NewEvidenceInteractor creates a new EvidenceInteractor instance
func NewEvidenceInteractor(evidenceSvc input.EvidenceService) *EvidenceInteractor {
	return &EvidenceInteractor{
		evidenceSvc: evidenceSvc,
	}
}

// CreateEvidence creates a new photographic evidence for a motorcycle (HU16)
func (i *EvidenceInteractor) CreateEvidence(ctx context.Context, motorcycleID, ownerID string, evidence *domain.MotorcycleEvidence) (result *domain.MotorcycleEvidence, err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogEvidenceInteractorCreateStart, "motorcycle_id", motorcycleID, "owner_id", ownerID)

	// Step 1: Validate motorcycle exists and ownership
	if _, err = i.evidenceSvc.ValidateMotorcycleOwnership(ctx, motorcycleID, ownerID); err != nil {
		return nil, err
	}

	// Step 2: Check evidence limit
	if err = i.evidenceSvc.CheckEvidenceLimit(ctx, motorcycleID); err != nil {
		return nil, err
	}

	// Step 3: Begin transaction
	tx, txErr := i.evidenceSvc.BeginTx(ctx)
	if txErr != nil {
		log.Error(logger.LogEvidenceInteractorBeginTxError, "error", txErr)
		return nil, domain.ErrEvidenceCannotSave
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogEvidenceInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogEvidenceInteractorRollbackOK)
			}
		}
	}()

	// Step 4: Create evidence (factory + save)
	result, err = i.evidenceSvc.CreateEvidence(ctx, tx, motorcycleID, evidence)
	if err != nil {
		return nil, err
	}

	// Step 5: Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogEvidenceInteractorCommitError, "error", err)
		return nil, domain.ErrEvidenceCannotSave
	}

	log.Success(logger.LogEvidenceInteractorCreateSuccess, "id", result.ID, "motorcycle_id", motorcycleID)

	err = nil
	return result, nil
}

// GetEvidenceByID retrieves an evidence by its ID (HU18)
func (i *EvidenceInteractor) GetEvidenceByID(ctx context.Context, evidenceID, ownerID string) (*domain.MotorcycleEvidence, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogEvidenceInteractorGetStart, "evidence_id", evidenceID)

	// Step 1: Get evidence
	evidence, err := i.evidenceSvc.GetByID(ctx, evidenceID)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorGetError, "error", err, "evidence_id", evidenceID)
		return nil, err
	}

	// Step 2: Validate ownership through motorcycle
	if _, err = i.evidenceSvc.ValidateMotorcycleOwnership(ctx, evidence.MotorcycleID, ownerID); err != nil {
		log.Warn(logger.LogEvidenceInteractorOwnerError, "evidence_id", evidenceID, "owner_id", ownerID)
		return nil, domain.ErrEvidenceNotFound
	}

	log.Success(logger.LogEvidenceInteractorGetSuccess, "evidence_id", evidenceID)
	return evidence, nil
}

// ListEvidenceByMotorcycle retrieves all evidence for a motorcycle (HU18)
func (i *EvidenceInteractor) ListEvidenceByMotorcycle(ctx context.Context, motorcycleID, ownerID string) ([]domain.MotorcycleEvidence, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogEvidenceInteractorListStart, "motorcycle_id", motorcycleID)

	// Step 1: Validate motorcycle exists and ownership
	if _, err := i.evidenceSvc.ValidateMotorcycleOwnership(ctx, motorcycleID, ownerID); err != nil {
		return nil, err
	}

	// Step 2: Get all evidence
	evidences, err := i.evidenceSvc.GetByMotorcycleID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorListError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}

	log.Success(logger.LogEvidenceInteractorListSuccess, "motorcycle_id", motorcycleID, "count", len(evidences))
	return evidences, nil
}

// DeleteEvidence deletes an evidence from Firebase Storage and DB (HU19)
func (i *EvidenceInteractor) DeleteEvidence(ctx context.Context, evidenceID, ownerID string) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogEvidenceInteractorDeleteStart, "evidence_id", evidenceID)

	// Step 1: Get evidence
	evidence, getErr := i.evidenceSvc.GetByID(ctx, evidenceID)
	if getErr != nil {
		log.Error(logger.LogEvidenceInteractorGetError, "error", getErr, "evidence_id", evidenceID)
		return getErr
	}

	// Step 2: Validate ownership through motorcycle
	if _, err = i.evidenceSvc.ValidateMotorcycleOwnership(ctx, evidence.MotorcycleID, ownerID); err != nil {
		log.Warn(logger.LogEvidenceInteractorOwnerError, "evidence_id", evidenceID, "owner_id", ownerID)
		return domain.ErrEvidenceNotFound
	}

	// Step 3: Delete image from Firebase Storage (best-effort)
	i.evidenceSvc.DeleteStorageFile(ctx, evidence.ImageURL)

	// Step 4: Begin transaction
	tx, txErr := i.evidenceSvc.BeginTx(ctx)
	if txErr != nil {
		log.Error(logger.LogEvidenceInteractorBeginTxError, "error", txErr)
		return domain.ErrEvidenceCannotDelete
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogEvidenceInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogEvidenceInteractorRollbackOK)
			}
		}
	}()

	// Step 5: Delete evidence from DB
	if err = i.evidenceSvc.DeleteEvidence(ctx, tx, evidenceID); err != nil {
		return err
	}

	// Step 6: Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogEvidenceInteractorCommitError, "error", err)
		return domain.ErrEvidenceCannotDelete
	}

	log.Success(logger.LogEvidenceInteractorDeleteSuccess, "evidence_id", evidenceID)

	err = nil
	return nil
}

// UpdateEvidence updates an existing evidence (HU17)
func (i *EvidenceInteractor) UpdateEvidence(ctx context.Context, evidenceID, ownerID string, updates *domain.MotorcycleEvidence) (result *domain.MotorcycleEvidence, err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogEvidenceInteractorUpdateStart, "evidence_id", evidenceID)

	// Step 1: Get existing evidence
	evidence, getErr := i.evidenceSvc.GetByID(ctx, evidenceID)
	if getErr != nil {
		log.Error(logger.LogEvidenceInteractorGetError, "error", getErr, "evidence_id", evidenceID)
		return nil, domain.ErrEvidenceNotFound
	}

	// Step 2: Validate ownership through motorcycle
	if _, err = i.evidenceSvc.ValidateMotorcycleOwnership(ctx, evidence.MotorcycleID, ownerID); err != nil {
		log.Warn(logger.LogEvidenceInteractorOwnerError, "evidence_id", evidenceID, "owner_id", ownerID)
		return nil, domain.ErrEvidenceNotFound
	}

	// Step 3: Apply updates (merges fields + cleans up old storage if URL changed)
	i.evidenceSvc.ApplyUpdatesAndCleanup(ctx, evidence, updates)

	// Step 4: Begin transaction
	tx, txErr := i.evidenceSvc.BeginTx(ctx)
	if txErr != nil {
		log.Error(logger.LogEvidenceInteractorBeginTxError, "error", txErr)
		return nil, domain.ErrEvidenceCannotUpdate
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogEvidenceInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogEvidenceInteractorRollbackOK)
			}
		}
	}()

	// Step 5: Update evidence
	if err = i.evidenceSvc.UpdateEvidence(ctx, tx, evidence); err != nil {
		return nil, err
	}

	// Step 6: Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogEvidenceInteractorCommitError, "error", err)
		return nil, domain.ErrEvidenceCannotUpdate
	}

	log.Success(logger.LogEvidenceInteractorUpdateSuccess, "evidence_id", evidenceID)

	err = nil
	return evidence, nil
}

// LookupEvidence retrieves all evidence for a motorcycle without ownership check.
// Used by workshop representatives during plate lookup to display the evidence gallery.
func (i *EvidenceInteractor) LookupEvidence(ctx context.Context, motorcycleID string) ([]domain.MotorcycleEvidence, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogEvidenceInteractorLookupStart, "motorcycle_id", motorcycleID)

	evidences, err := i.evidenceSvc.GetByMotorcycleID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorLookupError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}

	log.Success(logger.LogEvidenceInteractorLookupSuccess, "motorcycle_id", motorcycleID, "count", len(evidences))
	return evidences, nil
}
