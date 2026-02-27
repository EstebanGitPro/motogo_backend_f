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

// UpdateStatusWithPrice updates the status, completion date, and final price of a completed service
func (r *repository) UpdateStatusWithPrice(ctx context.Context, tx output.Tx, serviceID string, status string, completionDate *time.Time, finalPrice *float64) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	var cd sql.NullTime
	if completionDate != nil {
		cd = sql.NullTime{Time: *completionDate, Valid: true}
	}

	var fp sql.NullFloat64
	if finalPrice != nil {
		fp = sql.NullFloat64{Float64: *finalPrice, Valid: true}
	}

	_, err := sqlTx.ExecContext(ctx, queryUpdateStatusWithPrice, status, cd, fp, serviceID)
	if err != nil {
		log.Error(logger.LogCSRepoUpdateStatusErr, "service_id", serviceID, "error", err)
		return err
	}

	return nil
}

// UpdateDetails updates the quoted price, final price, and representative notes of a completed service
func (r *repository) UpdateDetails(ctx context.Context, tx output.Tx, serviceID string, quotedPrice, finalPrice *float64, notes *string) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	var qp sql.NullFloat64
	if quotedPrice != nil {
		qp = sql.NullFloat64{Float64: *quotedPrice, Valid: true}
	}

	var fp sql.NullFloat64
	if finalPrice != nil {
		fp = sql.NullFloat64{Float64: *finalPrice, Valid: true}
	}

	var n sql.NullString
	if notes != nil {
		n = sql.NullString{String: *notes, Valid: true}
	}

	_, err := sqlTx.ExecContext(ctx, queryUpdateDetails, qp, fp, n, serviceID)
	if err != nil {
		log.Error(logger.LogCSRepoUpdDetailsErr, "service_id", serviceID, "error", err)
		return err
	}

	return nil
}
