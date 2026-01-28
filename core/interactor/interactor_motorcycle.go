package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/google/uuid"
)

// MotorcycleInteractor handles motorcycle-related use cases (HU43-47)
type MotorcycleInteractor struct {
	motorcycleRepo output.MotorcycleRepository
	logger         logger.Logger
}

// NewMotorcycleInteractor creates a new MotorcycleInteractor instance
func NewMotorcycleInteractor(motorcycleRepo output.MotorcycleRepository, log logger.Logger) *MotorcycleInteractor {
	return &MotorcycleInteractor{
		motorcycleRepo: motorcycleRepo,
		logger:         log,
	}
}

// RegisterMotorcycle registers a new motorcycle for the authenticated user (HU43)
func (i *MotorcycleInteractor) RegisterMotorcycle(ctx context.Context, motorcycle *domain.Motorcycle) (*domain.Motorcycle, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleInteractorRegStart, "license_plate", motorcycle.LicensePlate, "owner_id", motorcycle.OwnerID)

	// Step 1: Validate reference exists (only if provided - optional until Release 11)
	if motorcycle.ReferenceID != "" {
		refExists, err := i.motorcycleRepo.ValidateReferenceExists(ctx, motorcycle.ReferenceID)
		if err != nil {
			log.Error(logger.LogMotorcycleInteractorRefError, "error", err, "reference_id", motorcycle.ReferenceID)
			return nil, domain.ErrMotorcycleCannotSave
		}
		if !refExists {
			log.Warn(logger.LogMotorcycleInteractorRefNotFound, "reference_id", motorcycle.ReferenceID)
			return nil, domain.ErrReferenceNotFound
		}
	}

	// Step 2: Validate license plate is unique
	plateExists, err := i.motorcycleRepo.CheckLicensePlateExists(ctx, motorcycle.LicensePlate)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorCheckPlateErr, "error", err, "license_plate", motorcycle.LicensePlate)
		return nil, domain.ErrMotorcycleCannotSave
	}
	if plateExists {
		log.Warn(logger.LogMotorcycleInteractorDupPlate, "license_plate", motorcycle.LicensePlate)
		return nil, domain.ErrDuplicateLicensePlate
	}

	// Step 3: Generate UUID
	motorcycle.ID = uuid.New().String()
	log.Debug(logger.LogMotorcycleInteractorIDGenerated, "id", motorcycle.ID)

	// Step 4: Begin transaction
	tx, err := i.motorcycleRepo.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorBeginTxError, "error", err)
		return nil, domain.ErrMotorcycleCannotSave
	}

	// Step 5: Save motorcycle
	err = i.motorcycleRepo.Save(ctx, tx, motorcycle)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorSaveError, "error", err)
		tx.Rollback()
		return nil, domain.ErrMotorcycleCannotSave
	}

	// Step 6: Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error(logger.LogMotorcycleInteractorCommitError, "error", err)
		return nil, domain.ErrMotorcycleCannotSave
	}

	log.Success(logger.LogMotorcycleInteractorRegSuccess, "id", motorcycle.ID, "license_plate", motorcycle.LicensePlate)
	return motorcycle, nil
}

// GetMotorcycleByID retrieves a motorcycle by its ID (HU46)
func (i *MotorcycleInteractor) GetMotorcycleByID(ctx context.Context, motorcycleID string) (*domain.Motorcycle, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleInteractorGetStart, "motorcycle_id", motorcycleID)

	motorcycle, err := i.motorcycleRepo.GetByID(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorGetError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}

	log.Success(logger.LogMotorcycleInteractorGetSuccess, "motorcycle_id", motorcycleID)
	return motorcycle, nil
}

// GetMotorcyclesByOwner retrieves all motorcycles owned by a person (HU47)
func (i *MotorcycleInteractor) GetMotorcyclesByOwner(ctx context.Context, ownerID string) ([]domain.Motorcycle, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogMotorcycleInteractorGetOwnerStart, "owner_id", ownerID)

	motorcycles, err := i.motorcycleRepo.GetByOwnerID(ctx, ownerID)
	if err != nil {
		log.Error(logger.LogMotorcycleInteractorGetOwnerError, "error", err, "owner_id", ownerID)
		return nil, err
	}

	log.Success(logger.LogMotorcycleInteractorGetOwnerSuccess, "owner_id", ownerID, "count", len(motorcycles))
	return motorcycles, nil
}
