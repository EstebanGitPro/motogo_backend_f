package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// EvidenceInteractor handles motorcycle evidence-related use cases (HU16-19)
type EvidenceInteractor struct {
	evidenceService input.EvidenceService
}

// NewEvidenceInteractor creates a new EvidenceInteractor instance
func NewEvidenceInteractor(
	evidenceService input.EvidenceService,
) *EvidenceInteractor {
	return &EvidenceInteractor{
		evidenceService: evidenceService,
	}
}

// WithStorageClient sets the storage client for image deletion (optional)
func (i *EvidenceInteractor) WithStorageClient(client output.StorageClient) *EvidenceInteractor {
	i.evidenceService.WithStorageClient(client)
	return i
}

// CreateEvidence creates a new photographic evidence for a motorcycle (HU16)
func (i *EvidenceInteractor) CreateEvidence(ctx context.Context, motorcycleID, ownerID string, evidence *domain.MotorcycleEvidence) (*domain.MotorcycleEvidence, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogEvidenceInteractorCreateStart, "motorcycle_id", motorcycleID, "owner_id", ownerID)

	// 1. Validate motorcycle ownership
	if err := i.evidenceService.ValidateMotorcycleOwnership(ctx, motorcycleID, ownerID); err != nil {
		return nil, err
	}

	// 2. Check evidence limit
	if err := i.evidenceService.CheckEvidenceLimit(ctx, motorcycleID); err != nil {
		return nil, err
	}

	// 3. Begin transaction
	tx, err := i.evidenceService.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorBeginTxError, "error", err)
		return nil, domain.ErrEvidenceCannotSave
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogEvidenceInteractorCommitError, "rollback_error", rbErr, "original_error", err)
			}
		}
	}()

	// 4. Create evidence via service
	newEvidence, err := i.evidenceService.CreateEvidence(ctx, tx, motorcycleID, evidence)
	if err != nil {
		return nil, err
	}

	// 5. Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogEvidenceInteractorCommitError, "error", err)
		return nil, domain.ErrEvidenceCannotSave
	}

	log.Success(logger.LogEvidenceInteractorCreateSuccess, "id", newEvidence.ID, "motorcycle_id", motorcycleID)
	err = nil
	return newEvidence, nil
}

// GetEvidenceByID retrieves an evidence by its ID (HU18)
func (i *EvidenceInteractor) GetEvidenceByID(ctx context.Context, evidenceID, ownerID string) (*domain.MotorcycleEvidence, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogEvidenceInteractorGetStart, "evidence_id", evidenceID)

	// 1. Get evidence
	evidence, err := i.evidenceService.GetEvidenceByID(ctx, evidenceID)
	if err != nil {
		return nil, err
	}

	// 2. Validate ownership through motorcycle
	if err := i.evidenceService.ValidateMotorcycleOwnership(ctx, evidence.MotorcycleID, ownerID); err != nil {
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

	// 1. Validate motorcycle ownership
	if err := i.evidenceService.ValidateMotorcycleOwnership(ctx, motorcycleID, ownerID); err != nil {
		return nil, err
	}

	// 2. Get evidences via service
	evidences, err := i.evidenceService.GetEvidenceByMotorcycleID(ctx, motorcycleID)
	if err != nil {
		return nil, err
	}

	log.Success(logger.LogEvidenceInteractorListSuccess, "motorcycle_id", motorcycleID, "count", len(evidences))
	return evidences, nil
}

// DeleteEvidence deletes an evidence from Firebase Storage and DB (HU19)
func (i *EvidenceInteractor) DeleteEvidence(ctx context.Context, evidenceID, ownerID string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogEvidenceInteractorDeleteStart, "evidence_id", evidenceID)

	// 1. Get evidence
	evidence, err := i.evidenceService.GetEvidenceByID(ctx, evidenceID)
	if err != nil {
		return err
	}

	// 2. Validate ownership through motorcycle
	if err := i.evidenceService.ValidateMotorcycleOwnership(ctx, evidence.MotorcycleID, ownerID); err != nil {
		log.Warn(logger.LogEvidenceInteractorOwnerError, "evidence_id", evidenceID, "owner_id", ownerID)
		return domain.ErrEvidenceNotFound
	}

	// 3. Delete image from storage (best-effort, via service)
	i.evidenceService.DeleteStorageFile(ctx, evidence.ImageURL)

	// 4. Begin transaction
	tx, err := i.evidenceService.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorBeginTxError, "error", err)
		return domain.ErrEvidenceCannotDelete
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogEvidenceInteractorCommitError, "rollback_error", rbErr, "original_error", err)
			}
		}
	}()

	// 5. Delete evidence from DB via service
	if err = i.evidenceService.DeleteEvidence(ctx, tx, evidenceID); err != nil {
		return err
	}

	// 6. Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogEvidenceInteractorCommitError, "error", err)
		return domain.ErrEvidenceCannotDelete
	}

	log.Success(logger.LogEvidenceInteractorDeleteSuccess, "evidence_id", evidenceID)
	err = nil
	return nil
}

// UpdateEvidence updates an existing evidence (HU17)
func (i *EvidenceInteractor) UpdateEvidence(ctx context.Context, evidenceID, ownerID string, updates *domain.MotorcycleEvidence) (*domain.MotorcycleEvidence, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogEvidenceInteractorUpdateStart, "evidence_id", evidenceID)

	// 1. Get existing evidence
	evidence, err := i.evidenceService.GetEvidenceByID(ctx, evidenceID)
	if err != nil {
		return nil, domain.ErrEvidenceNotFound
	}

	// 2. Validate ownership through motorcycle
	if err := i.evidenceService.ValidateMotorcycleOwnership(ctx, evidence.MotorcycleID, ownerID); err != nil {
		log.Warn(logger.LogEvidenceInteractorOwnerError, "evidence_id", evidenceID, "owner_id", ownerID)
		return nil, domain.ErrEvidenceNotFound
	}

	// 3. Apply updates via service (handles old image cleanup + field patching)
	i.evidenceService.ApplyEvidenceUpdates(ctx, evidence, updates)

	// 4. Begin transaction
	tx, err := i.evidenceService.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorBeginTxError, "error", err)
		return nil, domain.ErrEvidenceCannotUpdate
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogEvidenceInteractorCommitError, "rollback_error", rbErr, "original_error", err)
			}
		}
	}()

	// 5. Update evidence via service
	if err = i.evidenceService.UpdateEvidence(ctx, tx, evidence); err != nil {
		return nil, err
	}

	// 6. Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogEvidenceInteractorCommitError, "error", err)
		return nil, domain.ErrEvidenceCannotUpdate
	}

	log.Success(logger.LogEvidenceInteractorUpdateSuccess, "evidence_id", evidenceID)
	err = nil
	return evidence, nil
}
