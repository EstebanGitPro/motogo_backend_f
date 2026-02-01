package motorcycle

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// Delete performs a soft delete on a motorcycle (HU45)
// Sets deleted_at timestamp instead of physically removing the record
func (r *repository) Delete(ctx context.Context, tx output.Tx, motorcycleID string) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		log.Error(logger.LogMotorcycleRepoInvalidTx, "expected SQLTx")
		return domain.ErrInvalidTransaction
	}

	result, err := sqlTx.ExecContext(ctx, queryDelete, motorcycleID)
	if err != nil {
		log.Error(logger.LogMotorcycleRepoDeleteError, "error", err, "motorcycle_id", motorcycleID)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error(logger.LogMotorcycleRepoDeleteError, "error getting rows affected", err)
		return err
	}

	if rowsAffected == 0 {
		log.Warn(logger.LogMotorcycleRepoDeleteError, "no rows affected", "motorcycle_id", motorcycleID)
		return domain.ErrMotorcycleNotFound
	}

	log.Success(logger.LogMotorcycleRepoDeleteSuccess, "motorcycle_id", motorcycleID)
	return nil
}
