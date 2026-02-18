package motorcycle

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

const (
	queryInsert = `
		INSERT INTO motorcycles (id, license_plate, reference_id, owner_id, year, current_mileage, owner_notes, profile_image_url)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	queryValidateReference = `SELECT EXISTS(SELECT 1 FROM motorcycle_references WHERE id = ?)`
	queryCheckLicensePlate = `SELECT EXISTS(SELECT 1 FROM motorcycles WHERE license_plate = ?)`
)

func (r *repository) Save(ctx context.Context, tx output.Tx, motorcycle *domain.Motorcycle) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		log.Error(logger.LogMotorcycleRepoInvalidTx, "error", "invalid transaction type")
		return domain.ErrMotorcycleCannotSave
	}

	// NOTE: ReferenceID uses sql.NullString to support optional values.
	// Convert ReferenceID to sql.NullString
	var referenceID sql.NullString
	if motorcycle.ReferenceID != "" {
		referenceID = sql.NullString{String: motorcycle.ReferenceID, Valid: true}
	}

	_, err := sqlTx.ExecContext(ctx, queryInsert,
		motorcycle.ID,
		motorcycle.LicensePlate,
		referenceID,
		motorcycle.OwnerID,
		motorcycle.Year,
		motorcycle.CurrentMileage,
		motorcycle.OwnerNotes,
		motorcycle.ProfileImageURL,
	)
	if err != nil {
		log.Error(logger.LogMotorcycleRepoSaveError, "error", err, "license_plate", motorcycle.LicensePlate)
		return err
	}

	log.Success(logger.LogMotorcycleRepoSaveSuccess, "id", motorcycle.ID, "license_plate", motorcycle.LicensePlate)
	return nil
}

func (r *repository) ValidateReferenceExists(ctx context.Context, referenceID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, queryValidateReference, referenceID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Error(logger.LogMotorcycleRepoValidateRefErr, "error", err, "reference_id", referenceID)
		return false, err
	}
	return exists, nil
}

func (r *repository) CheckLicensePlateExists(ctx context.Context, licensePlate string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, queryCheckLicensePlate, licensePlate).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Error(logger.LogMotorcycleRepoCheckPlateErr, "error", err, "license_plate", licensePlate)
		return false, err
	}
	return exists, nil
}
