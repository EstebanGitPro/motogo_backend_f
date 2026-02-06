package interactor

import (
	"context"
	"strings"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/google/uuid"
)

const (
	// MaxEvidencePerMotorcycle is the maximum number of evidence allowed per motorcycle
	MaxEvidencePerMotorcycle = 5
	// FirebaseStorageHost is the required host for valid Firebase Storage URLs
	FirebaseStorageHost = "firebasestorage.googleapis.com"
)

// EvidenceInteractor handles motorcycle evidence-related use cases (HU16-19)
type EvidenceInteractor struct {
	evidenceRepo   output.EvidenceRepository
	motorcycleRepo output.MotorcycleRepository
	logger         logger.Logger
}

// NewEvidenceInteractor creates a new EvidenceInteractor instance
func NewEvidenceInteractor(
	evidenceRepo output.EvidenceRepository,
	motorcycleRepo output.MotorcycleRepository,
	log logger.Logger,
) *EvidenceInteractor {
	return &EvidenceInteractor{
		evidenceRepo:   evidenceRepo,
		motorcycleRepo: motorcycleRepo,
		logger:         log,
	}
}

// CreateEvidence creates a new photographic evidence for a motorcycle (HU16)
func (i *EvidenceInteractor) CreateEvidence(ctx context.Context, motorcycleID, ownerID string, evidence *domain.MotorcycleEvidence) (*domain.MotorcycleEvidence, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

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

	// Step 3: Validate image URL is Firebase Storage
	if !isValidFirebaseURL(evidence.ImageURL) {
		log.Warn(logger.LogEvidenceInteractorURLInvalid, "image_url", evidence.ImageURL)
		return nil, domain.ErrInvalidEvidenceURL
	}

	// Step 3.5: Validate angle (optional but must be valid if provided)
	if evidence.Angle != nil && !isValidAngle(*evidence.Angle) {
		log.Warn(logger.LogEvidenceInteractorAngleInvalid, "angle", *evidence.Angle)
		return nil, domain.ErrInvalidEvidenceAngle
	}

	// Step 4: Check evidence limit
	count, err := i.evidenceRepo.CountByMotorcycleID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorCountError, "error", err, "motorcycle_id", motorcycleID)
		return nil, domain.ErrEvidenceCannotSave
	}
	if count >= MaxEvidencePerMotorcycle {
		log.Warn(logger.LogEvidenceInteractorLimitExceeded, "motorcycle_id", motorcycleID, "count", count)
		return nil, domain.ErrEvidenceLimitExceeded
	}

	// Step 5: Generate UUID and set fields
	evidence.ID = uuid.New().String()
	evidence.MotorcycleID = motorcycleID
	evidence.UploadDate = time.Now()
	log.Debug(logger.LogEvidenceInteractorIDGenerated, "id", evidence.ID)

	// Step 6: Begin transaction
	tx, err := i.evidenceRepo.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorBeginTxError, "error", err)
		return nil, domain.ErrEvidenceCannotSave
	}

	// Step 7: Save evidence
	err = i.evidenceRepo.Save(ctx, tx, evidence)
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

	log.Success(logger.LogEvidenceInteractorCreateSuccess, "id", evidence.ID, "motorcycle_id", motorcycleID)
	return evidence, nil
}

// GetEvidenceByID retrieves an evidence by its ID (HU18)
func (i *EvidenceInteractor) GetEvidenceByID(ctx context.Context, evidenceID, ownerID string) (*domain.MotorcycleEvidence, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

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
	log := i.logger.WithTraceID(traceID)

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

// DeleteEvidence deletes an evidence (HU19)
func (i *EvidenceInteractor) DeleteEvidence(ctx context.Context, evidenceID, ownerID string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

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

	// Step 3: Begin transaction
	tx, err := i.evidenceRepo.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorBeginTxError, "error", err)
		return domain.ErrEvidenceCannotDelete
	}

	// Step 4: Delete evidence
	err = i.evidenceRepo.Delete(ctx, tx, evidenceID)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorDeleteError, "error", err, "evidence_id", evidenceID)
		tx.Rollback()
		return domain.ErrEvidenceCannotDelete
	}

	// Step 5: Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error(logger.LogEvidenceInteractorCommitError, "error", err)
		return domain.ErrEvidenceCannotDelete
	}

	log.Success(logger.LogEvidenceInteractorDeleteSuccess, "evidence_id", evidenceID)
	return nil
}

// isValidFirebaseURL validates that the URL is from Firebase Storage
func isValidFirebaseURL(url string) bool {
	return strings.HasPrefix(url, "https://"+FirebaseStorageHost)
}

// isValidAngle validates that the angle is one of the valid values
func isValidAngle(angle string) bool {
	for _, valid := range domain.ValidEvidenceAngles {
		if angle == valid {
			return true
		}
	}
	return false
}
