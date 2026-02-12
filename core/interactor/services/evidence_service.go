package services

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

var evidenceLog logger.Logger = logger.NewSlogLogger()

// EvidenceServiceImpl implements input.EvidenceService
type EvidenceServiceImpl struct {
	evidenceRepo   output.EvidenceRepository
	motorcycleRepo output.MotorcycleRepository
	storageClient  output.StorageClient // Optional: Firebase Storage for image deletion
}

// NewEvidenceService creates a new EvidenceService instance
func NewEvidenceService(
	evidenceRepo output.EvidenceRepository,
	motorcycleRepo output.MotorcycleRepository,
) *EvidenceServiceImpl {
	return &EvidenceServiceImpl{
		evidenceRepo:   evidenceRepo,
		motorcycleRepo: motorcycleRepo,
	}
}

// WithStorageClient sets the storage client for image deletion (optional)
func (s *EvidenceServiceImpl) WithStorageClient(client output.StorageClient) *EvidenceServiceImpl {
	s.storageClient = client
	return s
}

// BeginTx starts a new transaction
func (s *EvidenceServiceImpl) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.evidenceRepo.BeginTx(ctx)
}

// ValidateMotorcycleOwnership validates that the motorcycle exists and belongs to the owner
func (s *EvidenceServiceImpl) ValidateMotorcycleOwnership(ctx context.Context, motorcycleID, ownerID string) (*domain.Motorcycle, error) {
	motorcycle, err := s.motorcycleRepo.GetByID(ctx, motorcycleID)
	if err != nil {
		evidenceLog.Error(logger.LogEvidenceInteractorMotorcycleError, "error", err, "motorcycle_id", motorcycleID)
		return nil, domain.ErrMotorcycleNotFound
	}

	if motorcycle.OwnerID != ownerID {
		evidenceLog.Warn(logger.LogEvidenceInteractorOwnerError, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return nil, domain.ErrMotorcycleNotFound
	}

	return motorcycle, nil
}

// CheckEvidenceLimit validates that the motorcycle has not exceeded the evidence limit
func (s *EvidenceServiceImpl) CheckEvidenceLimit(ctx context.Context, motorcycleID string) error {
	count, err := s.evidenceRepo.CountByMotorcycleID(ctx, motorcycleID)
	if err != nil {
		evidenceLog.Error(logger.LogEvidenceInteractorCountError, "error", err, "motorcycle_id", motorcycleID)
		return domain.ErrEvidenceCannotSave
	}
	if IsEvidenceLimitReached(count) {
		evidenceLog.Warn(logger.LogEvidenceInteractorLimitExceeded, "motorcycle_id", motorcycleID, "count", count)
		return domain.ErrEvidenceLimitExceeded
	}
	return nil
}

// CreateEvidence creates a new evidence using the factory and saves it
func (s *EvidenceServiceImpl) CreateEvidence(ctx context.Context, tx output.Tx, motorcycleID string, evidence *domain.MotorcycleEvidence) (*domain.MotorcycleEvidence, error) {
	newEvidence := NewEvidence(motorcycleID, evidence.ImageURL, evidence.Angle, evidence.Description)
	evidenceLog.Debug(logger.LogEvidenceInteractorIDGenerated, "id", newEvidence.ID)

	if err := s.evidenceRepo.Save(ctx, tx, newEvidence); err != nil {
		evidenceLog.Error(logger.LogEvidenceInteractorSaveError, "error", err)
		return nil, domain.ErrEvidenceCannotSave
	}

	return newEvidence, nil
}

// GetByID retrieves an evidence by its ID
func (s *EvidenceServiceImpl) GetByID(ctx context.Context, evidenceID string) (*domain.MotorcycleEvidence, error) {
	return s.evidenceRepo.GetByID(ctx, evidenceID)
}

// GetByMotorcycleID retrieves all evidence for a motorcycle
func (s *EvidenceServiceImpl) GetByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.MotorcycleEvidence, error) {
	return s.evidenceRepo.GetByMotorcycleID(ctx, motorcycleID)
}

// ApplyUpdatesAndCleanup applies partial updates to an existing evidence and cleans up old storage if needed
func (s *EvidenceServiceImpl) ApplyUpdatesAndCleanup(ctx context.Context, existing *domain.MotorcycleEvidence, updates *domain.MotorcycleEvidence) {
	// If image URL changed, clean up old file from storage
	if updates.ImageURL != "" && updates.ImageURL != existing.ImageURL {
		if existing.ImageURL != "" && s.storageClient != nil {
			if storageErr := s.storageClient.DeleteStorageFile(ctx, existing.ImageURL); storageErr != nil {
				evidenceLog.Warn(logger.LogEvidenceInteractorDeleteError, "old image delete failed", storageErr, "url", existing.ImageURL)
			} else {
				evidenceLog.Info(logger.LogEvidenceInteractorDeleteSuccess, "action", "old_image_deleted", "url", existing.ImageURL)
			}
		}
		existing.ImageURL = updates.ImageURL
	}
	if updates.Angle != nil {
		existing.Angle = updates.Angle
	}
	if updates.Description != nil {
		existing.Description = updates.Description
	}
}

// UpdateEvidence updates an evidence record in the database
func (s *EvidenceServiceImpl) UpdateEvidence(ctx context.Context, tx output.Tx, evidence *domain.MotorcycleEvidence) error {
	if err := s.evidenceRepo.Update(ctx, tx, evidence); err != nil {
		evidenceLog.Error(logger.LogEvidenceInteractorUpdateError, "error", err, "evidence_id", evidence.ID)
		return domain.ErrEvidenceCannotUpdate
	}
	return nil
}

// DeleteEvidence deletes an evidence record from the database
func (s *EvidenceServiceImpl) DeleteEvidence(ctx context.Context, tx output.Tx, evidenceID string) error {
	if err := s.evidenceRepo.Delete(ctx, tx, evidenceID); err != nil {
		evidenceLog.Error(logger.LogEvidenceInteractorDeleteError, "error", err, "evidence_id", evidenceID)
		return domain.ErrEvidenceCannotDelete
	}
	return nil
}

// DeleteStorageFile deletes an image from Firebase Storage (best-effort)
func (s *EvidenceServiceImpl) DeleteStorageFile(ctx context.Context, imageURL string) {
	if imageURL == "" || s.storageClient == nil {
		return
	}
	if storageErr := s.storageClient.DeleteStorageFile(ctx, imageURL); storageErr != nil {
		evidenceLog.Warn(logger.LogEvidenceInteractorDeleteError, "storage delete failed (continuing)", storageErr, "url", imageURL)
	} else {
		evidenceLog.Info(logger.LogEvidenceInteractorDeleteSuccess, "action", "storage_file_deleted", "url", imageURL)
	}
}
