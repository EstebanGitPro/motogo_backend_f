package services

import (
	"context"
	"strings"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/EstebanGitPro/motogo-backend/tools/utils"
)

const (
	// FirebaseStorageHost is the required host for valid Firebase Storage URLs
	FirebaseStorageHost = "firebasestorage.googleapis.com"
	// MaxEvidencePerMotorcycle is the maximum number of evidence allowed per motorcycle
	MaxEvidencePerMotorcycle = 4
)

// evidenceService implements input.EvidenceService
type evidenceService struct {
	evidenceRepo   output.EvidenceRepository
	motorcycleRepo output.MotorcycleRepository
	storageClient  output.StorageClient // Optional: Firebase Storage for image deletion
}

// NewEvidenceService creates a new EvidenceService instance
func NewEvidenceService(
	evidenceRepo output.EvidenceRepository,
	motorcycleRepo output.MotorcycleRepository,
) input.EvidenceService {
	return &evidenceService{
		evidenceRepo:   evidenceRepo,
		motorcycleRepo: motorcycleRepo,
	}
}

// WithStorageClient sets the storage client for image deletion (optional)
func (s *evidenceService) WithStorageClient(client output.StorageClient) {
	s.storageClient = client
}

// BeginTx starts a new database transaction
func (s *evidenceService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.evidenceRepo.BeginTx(ctx)
}

// ValidateMotorcycleOwnership validates that the motorcycle exists and belongs to the given owner
func (s *evidenceService) ValidateMotorcycleOwnership(ctx context.Context, motorcycleID, ownerID string) error {
	motorcycle, err := s.motorcycleRepo.GetByID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorMotorcycleError, "error", err, "motorcycle_id", motorcycleID)
		return domain.ErrMotorcycleNotFound
	}
	if motorcycle.OwnerID != ownerID {
		log.Warn(logger.LogEvidenceInteractorOwnerError, "motorcycle_id", motorcycleID, "owner_id", ownerID)
		return domain.ErrMotorcycleNotFound
	}
	return nil
}

// CheckEvidenceLimit checks if the motorcycle has reached the evidence limit
func (s *evidenceService) CheckEvidenceLimit(ctx context.Context, motorcycleID string) error {
	count, err := s.evidenceRepo.CountByMotorcycleID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorCountError, "error", err, "motorcycle_id", motorcycleID)
		return domain.ErrEvidenceCannotSave
	}
	if count >= MaxEvidencePerMotorcycle {
		log.Warn(logger.LogEvidenceInteractorLimitExceeded, "motorcycle_id", motorcycleID, "count", count)
		return domain.ErrEvidenceLimitExceeded
	}
	return nil
}

// CreateEvidence creates a new evidence record with generated ID and timestamp
func (s *evidenceService) CreateEvidence(ctx context.Context, tx output.Tx, motorcycleID string, evidence *domain.MotorcycleEvidence) (*domain.MotorcycleEvidence, error) {
	newEvidence := &domain.MotorcycleEvidence{
		ID:           utils.Generate(),
		MotorcycleID: motorcycleID,
		ImageURL:     evidence.ImageURL,
		Angle:        evidence.Angle,
		Description:  evidence.Description,
		CreatedAt:    time.Now(),
	}
	log.Debug(logger.LogEvidenceInteractorIDGenerated, "id", newEvidence.ID)

	if err := s.evidenceRepo.Save(ctx, tx, newEvidence); err != nil {
		log.Error(logger.LogEvidenceInteractorSaveError, "error", err)
		return nil, domain.ErrEvidenceCannotSave
	}

	return newEvidence, nil
}

// GetEvidenceByID retrieves an evidence by its ID
func (s *evidenceService) GetEvidenceByID(ctx context.Context, evidenceID string) (*domain.MotorcycleEvidence, error) {
	evidence, err := s.evidenceRepo.GetByID(ctx, evidenceID)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorGetError, "error", err, "evidence_id", evidenceID)
		return nil, err
	}
	return evidence, nil
}

// GetEvidenceByMotorcycleID retrieves all evidence for a motorcycle
func (s *evidenceService) GetEvidenceByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.MotorcycleEvidence, error) {
	evidences, err := s.evidenceRepo.GetByMotorcycleID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogEvidenceInteractorListError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}
	return evidences, nil
}

// ApplyEvidenceUpdates applies partial updates to an existing evidence (field-by-field patching)
// Also handles old image cleanup from Firebase Storage when URL changes
func (s *evidenceService) ApplyEvidenceUpdates(ctx context.Context, existing *domain.MotorcycleEvidence, updates *domain.MotorcycleEvidence) {
	if updates.ImageURL != "" && updates.ImageURL != existing.ImageURL {
		// Delete old image from Firebase Storage when URL changes
		if existing.ImageURL != "" {
			s.DeleteStorageFile(ctx, existing.ImageURL)
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

// UpdateEvidence persists an evidence update to the database
func (s *evidenceService) UpdateEvidence(ctx context.Context, tx output.Tx, evidence *domain.MotorcycleEvidence) error {
	if err := s.evidenceRepo.Update(ctx, tx, evidence); err != nil {
		log.Error(logger.LogEvidenceInteractorUpdateError, "error", err, "evidence_id", evidence.ID)
		return domain.ErrEvidenceCannotUpdate
	}
	return nil
}

// DeleteEvidence deletes an evidence record from the database
func (s *evidenceService) DeleteEvidence(ctx context.Context, tx output.Tx, evidenceID string) error {
	if err := s.evidenceRepo.Delete(ctx, tx, evidenceID); err != nil {
		log.Error(logger.LogEvidenceInteractorDeleteError, "error", err, "evidence_id", evidenceID)
		return domain.ErrEvidenceCannotDelete
	}
	return nil
}

// DeleteStorageFile deletes a file from Firebase Storage (best-effort, logs errors but doesn't fail)
func (s *evidenceService) DeleteStorageFile(ctx context.Context, imageURL string) {
	if imageURL == "" || s.storageClient == nil {
		return
	}
	if err := s.storageClient.DeleteStorageFile(ctx, imageURL); err != nil {
		log.Warn(logger.LogEvidenceInteractorDeleteError, "storage delete failed (continuing)", err, "url", imageURL)
	} else {
		log.Info(logger.LogEvidenceInteractorDeleteSuccess, "action", "storage_file_deleted", "url", imageURL)
	}
}

// NewEvidence creates a new MotorcycleEvidence with generated ID and timestamp
// This is a factory function that encapsulates ID generation and timestamp assignment
func NewEvidence(motorcycleID, imageURL string, angle, description *string) *domain.MotorcycleEvidence {
	return &domain.MotorcycleEvidence{
		ID:           utils.Generate(),
		MotorcycleID: motorcycleID,
		Angle:        angle,
		ImageURL:     imageURL,
		Description:  description,
		CreatedAt:    time.Now(),
	}
}

// IsValidAngle checks if the angle is valid (FRONTAL, LATERAL, REAR)
func IsValidAngle(angle string) bool {
	for _, valid := range domain.ValidEvidenceAngles {
		if angle == valid {
			return true
		}
	}
	return false
}

// IsValidFirebaseURL validates that the URL is from Firebase Storage
func IsValidFirebaseURL(url string) bool {
	return strings.HasPrefix(url, "https://"+FirebaseStorageHost)
}

// IsEvidenceLimitReached checks if the evidence limit per motorcycle is reached
func IsEvidenceLimitReached(count int) bool {
	return count >= MaxEvidencePerMotorcycle
}
