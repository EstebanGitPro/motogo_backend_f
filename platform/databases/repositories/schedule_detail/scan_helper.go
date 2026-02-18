package schedule_detail

import (
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// scanScheduleDetails iterates over rows and scans each into a domain.ScheduleDetail.
// Centralizes the shared scan logic for all schedule detail queries.
func scanScheduleDetails(rows *sql.Rows) ([]domain.ScheduleDetail, error) {
	var details []domain.ScheduleDetail
	for rows.Next() {
		var d domain.ScheduleDetail
		var entryType string

		if err := rows.Scan(
			&d.ID,
			&d.ScheduleID,
			&entryType,
			&d.DayOfWeek,
			&d.ExceptionStartDate,
			&d.ExceptionEndDate,
			&d.OpeningTime,
			&d.ClosingTime,
			&d.IsClosed,
			&d.Active,
			&d.CreatedAt,
			&d.UpdatedAt,
		); err != nil {
			log.Error(logger.LogScheduleDetailRepoScanError, "error", err)
			return nil, err
		}

		d.EntryType = domain.EntryType(entryType)
		details = append(details, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return details, nil
}
