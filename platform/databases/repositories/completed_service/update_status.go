package completed_service

import (
	"context"
	"database/sql"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// UpdateStatus updates the status and optionally the completion date of a completed service (HU74)
func (r *repository) UpdateStatus(ctx context.Context, tx output.Tx, serviceID string, status string, completionDate *time.Time) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	var cd sql.NullTime
	if completionDate != nil {
		cd = sql.NullTime{Time: *completionDate, Valid: true}
	}

	_, err := sqlTx.ExecContext(ctx, queryUpdateStatus, status, cd, serviceID)
	if err != nil {
		log.Error(logger.LogCSRepoUpdateStatusErr, "service_id", serviceID, "error", err)
		return err
	}

	return nil
}

// GetStatusHistory retrieves the status transition history for a completed service (HU73)
func (r *repository) GetStatusHistory(ctx context.Context, serviceID string) ([]domain.ServiceStatusHistory, error) {
	rows, err := r.stmtGetStatusHistory.QueryContext(ctx, serviceID)
	if err != nil {
		log.Error(logger.LogCSRepoGetHistoryErr, err)
		return nil, err
	}
	defer rows.Close()

	var history []domain.ServiceStatusHistory
	for rows.Next() {
		var h ServiceStatusHistory
		if err := rows.Scan(
			&h.ID, &h.CompletedServiceID, &h.PreviousStatus,
			&h.NewStatus, &h.CreatedBy, &h.CreatedAt,
		); err != nil {
			log.Error(logger.LogCSRepoScanHistoryErr, err)
			return nil, err
		}
		history = append(history, h.HistoryToDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}
