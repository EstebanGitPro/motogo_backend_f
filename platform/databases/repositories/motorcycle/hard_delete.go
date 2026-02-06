package motorcycle

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

const (
	queryHardDelete = `
		DELETE FROM motorcycles 
		WHERE id = ?
	`

	queryClearProfileImage = `
		UPDATE motorcycles 
		SET profile_image_url = NULL, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`
)

// HardDelete permanently removes a motorcycle and all related data via CASCADE (HU45 hybrid)
// Only use when motorcycle has no service history
func (r *repository) HardDelete(ctx context.Context, tx output.Tx, motorcycleID string) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		log.Error(logger.LogMotorcycleRepoInvalidTx, "expected SQLTx for HardDelete")
		return domain.ErrInvalidTransaction
	}

	result, err := sqlTx.ExecContext(ctx, queryHardDelete, motorcycleID)
	if err != nil {
		log.Error(logger.LogMotorcycleRepoDeleteError, "hard delete error", err, "motorcycle_id", motorcycleID)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error(logger.LogMotorcycleRepoDeleteError, "error getting rows affected", err)
		return err
	}

	if rowsAffected == 0 {
		log.Warn(logger.LogMotorcycleRepoDeleteError, "no rows affected in hard delete", "motorcycle_id", motorcycleID)
		return domain.ErrMotorcycleNotFound
	}

	log.Success(logger.LogMotorcycleRepoDeleteSuccess, "hard deleted motorcycle_id", motorcycleID)
	return nil
}

// ClearProfileImageURL sets the profile_image_url to NULL (HU39)
func (r *repository) ClearProfileImageURL(ctx context.Context, tx output.Tx, motorcycleID string) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		log.Error(logger.LogMotorcycleRepoInvalidTx, "expected SQLTx for ClearProfileImageURL")
		return domain.ErrInvalidTransaction
	}

	result, err := sqlTx.ExecContext(ctx, queryClearProfileImage, motorcycleID)
	if err != nil {
		log.Error(logger.LogMotorcycleRepoUpdateError, "clear profile image error", err, "motorcycle_id", motorcycleID)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error(logger.LogMotorcycleRepoUpdateError, "error getting rows affected", err)
		return err
	}

	if rowsAffected == 0 {
		log.Warn(logger.LogMotorcycleRepoUpdateError, "no rows affected clearing profile image", "motorcycle_id", motorcycleID)
		return domain.ErrMotorcycleNotFound
	}

	log.Success(logger.LogMotorcycleRepoUpdateSuccess, "cleared profile image for motorcycle_id", motorcycleID)
	return nil
}

// HasServiceHistory checks if motorcycle has any completed_services or diagnostics (HU45 hybrid)
func (r *repository) HasServiceHistory(ctx context.Context, motorcycleID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, queryHasServiceHistory, motorcycleID, motorcycleID).Scan(&exists)
	if err != nil {
		log.Error(logger.LogMotorcycleRepoGetByIDError, "error checking service history", err, "motorcycle_id", motorcycleID)
		return false, err
	}

	return exists, nil
}
