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
func (i *EvidenceInteractor) CreateEvidence(ctx context.Context, motorcycleID, ownerID string, evidence *domain.MotorcycleEvidence) (*domain.MotorcycleEvidence, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogEvidenceInteractorCreateStart, "motorcycle_id", motorcycleID, "owner_id", ownerID)

	// Step 1: Validate motorcycle exists and ownership
	motorcycle, err := i.motorcycleRepo.GetByID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorMotorcycleError, "error", err, "motorcycle_id", motorcycleID)
		return nil, domain.ErrMotorcycleNotFound
	}

	// Step 2: Validate ownership (security by obscurity - 404 for non-owners)
	if motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogEvidenceInteractorOwnerError, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return nil, domain.ErrMotorcycleNotFound
	}

	// Note: Firebase URL and Angle validations are handled by JSON Schema middleware

	// Step 3: Check evidence limit (business rule - requires DB query)
	count, err := i.evidenceRepo.CountByMotorcycleID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorCountError, "error", err, "motorcycle_id", motorcycleID)
		return nil, domain.ErrEvidenceCannotSave
	}
	if services.IsEvidenceLimitReached(count) {
		log.Warn(logger.LogEvidenceInteractorLimitExceeded, "motorcycle_id", motorcycleID, "count", count)
		return nil, domain.ErrEvidenceLimitExceeded
	}

	// Step 4: Create evidence with ID and timestamp (delegated to service factory)
	newEvidence := services.NewEvidence(motorcycleID, evidence.ImageURL, evidence.Angle, evidence.Description)
	log.Debug(logger.LogEvidenceInteractorIDGenerated, "id", newEvidence.ID)

	// Step 6: Begin transaction
	tx, err := i.evidenceRepo.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorBeginTxError, "error", err)
		return nil, domain.ErrEvidenceCannotSave
	}

	// Step 7: Save evidence
	err = i.evidenceRepo.Save(ctx, tx, newEvidence)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorSaveError, "error", err)
		tx.Rollback()
		return nil, domain.ErrEvidenceCannotSave
	}

	// Step 8: Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error(logger.LogEvidenceInteractorCommitError, "error", err)
		return nil, domain.ErrEvidenceCannotSave
	}

	log.Success(logger.LogEvidenceInteractorCreateSuccess, "id", newEvidence.ID, "motorcycle_id", motorcycleID)
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
func (i *EvidenceInteractor) DeleteEvidence(ctx context.Context, evidenceID, ownerID string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogEvidenceInteractorDeleteStart, "evidence_id", evidenceID)

	// Step 1: Get evidence
	evidence, err := i.evidenceRepo.GetByID(ctx, evidenceID)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorGetError, "error", err, "evidence_id", evidenceID)
		return err
	}

	// Step 2: Validate ownership through motorcycle
	motorcycle, err := i.motorcycleRepo.GetByID(ctx, evidence.MotorcycleID)
	if err != nil || motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogEvidenceInteractorOwnerError, "evidence_id", evidenceID, "owner_id", ownerID)
		return domain.ErrEvidenceNotFound
	}

	// Step 3: Delete image from Firebase Storage (if client configured)
	if evidence.ImageURL != "" && i.storageClient != nil {
		if err := i.storageClient.DeleteStorageFile(ctx, evidence.ImageURL); err != nil {
			// Log warning but don't fail the delete - storage cleanup is best effort
			log.Warn(logger.LogEvidenceInteractorDeleteError, "storage delete failed (continuing)", err, "url", evidence.ImageURL)
		} else {
			log.Info(logger.LogEvidenceInteractorDeleteSuccess, "action", "storage_file_deleted", "url", evidence.ImageURL)
		}
	}

	// Step 4: Begin transaction
	tx, err := i.evidenceRepo.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorBeginTxError, "error", err)
		return domain.ErrEvidenceCannotDelete
	}

	// Step 5: Delete evidence from DB
	err = i.evidenceRepo.Delete(ctx, tx, evidenceID)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorDeleteError, "error", err, "evidence_id", evidenceID)
		tx.Rollback()
		return domain.ErrEvidenceCannotDelete
	}

	// Step 6: Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error(logger.LogEvidenceInteractorCommitError, "error", err)
		return domain.ErrEvidenceCannotDelete
	}

	log.Success(logger.LogEvidenceInteractorDeleteSuccess, "evidence_id", evidenceID)
	return nil
}

// UpdateEvidence updates an existing evidence (HU17)
func (i *EvidenceInteractor) UpdateEvidence(ctx context.Context, evidenceID, ownerID string, updates *domain.MotorcycleEvidence) (*domain.MotorcycleEvidence, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogEvidenceInteractorUpdateStart, "evidence_id", evidenceID)

	// Step 1: Get existing evidence
	evidence, err := i.evidenceRepo.GetByID(ctx, evidenceID)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorGetError, "error", err, "evidence_id", evidenceID)
		return nil, domain.ErrEvidenceNotFound
	}

	// Step 2: Validate ownership through motorcycle
	motorcycle, err := i.motorcycleRepo.GetByID(ctx, evidence.MotorcycleID)
	if err != nil || motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogEvidenceInteractorOwnerError, "evidence_id", evidenceID, "owner_id", ownerID)
		return nil, domain.ErrEvidenceNotFound
	}

	// Step 3: Apply updates (only updatable fields)
	// Note: Angle and ImageURL validations are handled by JSON Schema middleware
	if updates.ImageURL != "" && updates.ImageURL != evidence.ImageURL {
		// Delete old image from Firebase Storage when URL changes
		if evidence.ImageURL != "" && i.storageClient != nil {
			if err := i.storageClient.DeleteStorageFile(ctx, evidence.ImageURL); err != nil {
				log.Warn(logger.LogEvidenceInteractorDeleteError, "old image delete failed", err, "url", evidence.ImageURL)
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
	tx, err := i.evidenceRepo.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorBeginTxError, "error", err)
		return nil, domain.ErrEvidenceCannotUpdate
	}

	// Step 5: Update evidence
	err = i.evidenceRepo.Update(ctx, tx, evidence)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorUpdateError, "error", err, "evidence_id", evidenceID)
		tx.Rollback()
		return nil, domain.ErrEvidenceCannotUpdate
	}

	// Step 6: Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error(logger.LogEvidenceInteractorCommitError, "error", err)
		return nil, domain.ErrEvidenceCannotUpdate
	}

	log.Success(logger.LogEvidenceInteractorUpdateSuccess, "evidence_id", evidenceID)
	return evidence, nil
}
