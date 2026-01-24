package schedule_detail

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// GetDetailByID retrieves a schedule detail by its ID
func (r *repository) GetDetailByID(ctx context.Context, detailID string) (*domain.ScheduleDetail, error) {
	var detail domain.ScheduleDetail
	var entryType string

	err := r.stmtGetDetailByID.QueryRowContext(ctx, detailID).Scan(
		&detail.ID,
		&detail.ScheduleID,
		&entryType,
		&detail.DayOfWeek,
		&detail.ExceptionStartDate,
		&detail.ExceptionEndDate,
		&detail.OpeningTime,
		&detail.ClosingTime,
		&detail.IsClosed,
		&detail.Active,
		&detail.CreatedAt,
		&detail.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error(logger.LogScheduleDetailRepoGetByIDError, "detail_id", detailID, "error", err)
		return nil, err
	}

	detail.EntryType = domain.EntryType(entryType)
	return &detail, nil
}
