package completed_service

import (
	"context"
	"fmt"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// Delete removes a completed service by ID (HU65)
// Related items and status history are removed by DB cascade.
func (r *repository) Delete(ctx context.Context, tx output.Tx, serviceID string) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	_, err := sqlTx.ExecContext(ctx, queryDelete, serviceID)
	if err != nil {
		log.Error(logger.LogCSRepoDeleteError, "service_id", serviceID, "error", err)
		return fmt.Errorf("error deleting completed service: %w", err)
	}

	return nil
}

// SoftDelete marks a completed service as deleted (sets deleted_at timestamp)
// Used for FINALIZADO/CANCELADO services to preserve ratings and history.
func (r *repository) SoftDelete(ctx context.Context, tx output.Tx, serviceID string) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	_, err := sqlTx.ExecContext(ctx, querySoftDelete, serviceID)
	if err != nil {
		log.Error(logger.LogCSRepoDeleteError, "service_id", serviceID, "error", err)
		return fmt.Errorf("error soft-deleting completed service: %w", err)
	}

	return nil
}
