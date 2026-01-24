package schedule_detail

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// GetExceptionByID retrieves a specific exception by ID (HU21, HU22, HU24, HU25)
func (r *repository) GetExceptionByID(ctx context.Context, exceptionID string) (*domain.ScheduleDetail, error) {
	var exception domain.ScheduleDetail
	var entryType string

	err := r.stmtGetExceptionByID.QueryRowContext(ctx, exceptionID).Scan(
		&exception.ID,
		&exception.ScheduleID,
		&entryType,
		&exception.DayOfWeek,
		&exception.ExceptionStartDate,
		&exception.ExceptionEndDate,
		&exception.OpeningTime,
		&exception.ClosingTime,
		&exception.IsClosed,
		&exception.Active,
		&exception.CreatedAt,
		&exception.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error(logger.LogScheduleDetailRepoGetByIDError, "exception_id", exceptionID, "error", err)
		return nil, err
	}

	exception.EntryType = domain.EntryType(entryType)
	return &exception, nil
}
