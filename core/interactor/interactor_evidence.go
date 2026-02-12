package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// EvidenceInteractor handles motorcycle evidence-related use cases (HU16-19)
type EvidenceInteractor struct {
	evidenceRepo   output.EvidenceRepository
	motorcycleRepo output.MotorcycleRepository
	storageClient  output.StorageClient // Optional: Firebase Storage for image deletion
}

// NewEvidenceInteractor creates a new EvidenceInteractor instance
func NewEvidenceInteractor(
	evidenceRepo output.EvidenceRepository,
	motorcycleRepo output.MotorcycleRepository,
) *EvidenceInteractor {
	return &EvidenceInteractor{
		evidenceRepo:   evidenceRepo,
		motorcycleRepo: motorcycleRepo,
	}
}

// WithStorageClient sets the storage client for image deletion (optional)
func (i *EvidenceInteractor) WithStorageClient(client output.StorageClient) *EvidenceInteractor {
	i.storageClient = client
	return i
}

// CreateEvidence creates a new photographic evidence for a motorcycle (HU16)
func (i *EvidenceInteractor) CreateEvidence(ctx context.Context, motorcycleID, ownerID string, evidence *domain.MotorcycleEvidence) (result *domain.MotorcycleEvidence, err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogEvidenceInteractorCreateStart, "motorcycle_id", motorcycleID, "owner_id", ownerID)

	// Step 1: Validate motorcycle exists and ownership
	motorcycle, motoErr := i.motorcycleRepo.GetByID(ctx, motorcycleID)
	if motoErr != nil {
		log.Error(logger.LogEvidenceInteractorMotorcycleError, "error", motoErr, "motorcycle_id", motorcycleID)
		return nil, domain.ErrMotorcycleNotFound
	}

	// Step 2: Validate ownership (security by obscurity - 404 for non-owners)
	if motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogEvidenceInteractorOwnerError, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return nil, domain.ErrMotorcycleNotFound
	}

	// Note: Firebase URL and Angle validations are handled by JSON Schema middleware

	// Step 3: Check evidence limit (business rule - requires DB query)
	count, countErr := i.evidenceRepo.CountByMotorcycleID(ctx, motorcycleID)
	if countErr != nil {
		log.Error(logger.LogEvidenceInteractorCountError, "error", countErr, "motorcycle_id", motorcycleID)
		return nil, domain.ErrEvidenceCannotSave
	}
	if services.IsEvidenceLimitReached(count) {
		log.Warn(logger.LogEvidenceInteractorLimitExceeded, "motorcycle_id", motorcycleID, "count", count)
		return nil, domain.ErrEvidenceLimitExceeded
	}

	// Step 4: Create evidence with ID and timestamp (delegated to service factory)
	newEvidence := services.NewEvidence(motorcycleID, evidence.ImageURL, evidence.Angle, evidence.Description)
	log.Debug(logger.LogEvidenceInteractorIDGenerated, "id", newEvidence.ID)

	// Step 5: Begin transaction
	tx, txErr := i.evidenceRepo.BeginTx(ctx)
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

	// Step 6: Save evidence
	if err = i.evidenceRepo.Save(ctx, tx, newEvidence); err != nil {
		log.Error(logger.LogEvidenceInteractorSaveError, "error", err)
		return nil, domain.ErrEvidenceCannotSave
	}

	// Step 7: Commit transaction
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

	// Step 1: Get evidence
	evidence, err := i.evidenceRepo.GetByID(ctx, evidenceID)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorGetError, "error", err, "evidence_id", evidenceID)
		return nil, err
	}

	// Step 2: Validate ownership through motorcycle
	motorcycle, err := i.motorcycleRepo.GetByID(ctx, evidence.MotorcycleID)
	if err != nil || motorcycle.OwnerID != ownerID {
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
	motorcycle, err := i.motorcycleRepo.GetByID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorMotorcycleError, "error", err, "motorcycle_id", motorcycleID)
		return nil, domain.ErrMotorcycleNotFound
	}

	// Step 2: Validate ownership
	if motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogEvidenceInteractorOwnerError, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return nil, domain.ErrMotorcycleNotFound
	}

	// Step 3: Get all evidence
	evidences, err := i.evidenceRepo.GetByMotorcycleID(ctx, motorcycleID)
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
	evidence, getErr := i.evidenceRepo.GetByID(ctx, evidenceID)
	if getErr != nil {
		log.Error(logger.LogEvidenceInteractorGetError, "error", getErr, "evidence_id", evidenceID)
		return getErr
	}

	// Step 2: Validate ownership through motorcycle
	motorcycle, motoErr := i.motorcycleRepo.GetByID(ctx, evidence.MotorcycleID)
	if motoErr != nil || motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogEvidenceInteractorOwnerError, "evidence_id", evidenceID, "owner_id", ownerID)
		return domain.ErrEvidenceNotFound
	}

	// Step 3: Delete image from Firebase Storage (if client configured)
	if evidence.ImageURL != "" && i.storageClient != nil {
		if storageErr := i.storageClient.DeleteStorageFile(ctx, evidence.ImageURL); storageErr != nil {
			// Log warning but don't fail the delete - storage cleanup is best effort
			log.Warn(logger.LogEvidenceInteractorDeleteError, "storage delete failed (continuing)", storageErr, "url", evidence.ImageURL)
		} else {
			log.Info(logger.LogEvidenceInteractorDeleteSuccess, "action", "storage_file_deleted", "url", evidence.ImageURL)
		}
	}

	// Step 4: Begin transaction
	tx, txErr := i.evidenceRepo.BeginTx(ctx)
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
	if err = i.evidenceRepo.Delete(ctx, tx, evidenceID); err != nil {
		log.Error(logger.LogEvidenceInteractorDeleteError, "error", err, "evidence_id", evidenceID)
		return domain.ErrEvidenceCannotDelete
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
	evidence, getErr := i.evidenceRepo.GetByID(ctx, evidenceID)
	if getErr != nil {
		log.Error(logger.LogEvidenceInteractorGetError, "error", getErr, "evidence_id", evidenceID)
		return nil, domain.ErrEvidenceNotFound
	}

	// Step 2: Validate ownership through motorcycle
	motorcycle, motoErr := i.motorcycleRepo.GetByID(ctx, evidence.MotorcycleID)
	if motoErr != nil || motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogEvidenceInteractorOwnerError, "evidence_id", evidenceID, "owner_id", ownerID)
		return nil, domain.ErrEvidenceNotFound
	}

	// Step 3: Apply updates (only updatable fields)
	// Note: Angle and ImageURL validations are handled by JSON Schema middleware
	if updates.ImageURL != "" && updates.ImageURL != evidence.ImageURL {
		// Delete old image from Firebase Storage when URL changes
		if evidence.ImageURL != "" && i.storageClient != nil {
			if storageErr := i.storageClient.DeleteStorageFile(ctx, evidence.ImageURL); storageErr != nil {
				log.Warn(logger.LogEvidenceInteractorDeleteError, "old image delete failed", storageErr, "url", evidence.ImageURL)
			} else {
				log.Info(logger.LogEvidenceInteractorDeleteSuccess, "action", "old_image_deleted", "url", evidence.ImageURL)
			}
		}
		evidence.ImageURL = updates.ImageURL
	}
	if updates.Angle != nil {
		evidence.Angle = updates.Angle
	}
	if updates.Description != nil {
		evidence.Description = updates.Description
	}

	// Step 4: Begin transaction
	tx, txErr := i.evidenceRepo.BeginTx(ctx)
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
	if err = i.evidenceRepo.Update(ctx, tx, evidence); err != nil {
		log.Error(logger.LogEvidenceInteractorUpdateError, "error", err, "evidence_id", evidenceID)
		return nil, domain.ErrEvidenceCannotUpdate
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
