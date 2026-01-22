package services

import (
	"context"
	"regexp"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/google/uuid"
)

var scheduleDetailLog logger.Logger = logger.NewSlogLogger()

// timeRegex validates HH:MM format (24-hour)
var timeRegex = regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]$`)

// scheduleDetailService implements input.ScheduleDetailService
type scheduleDetailService struct {
	detailRepo   output.ScheduleDetailRepository
	scheduleRepo output.ScheduleRepository
}

// NewScheduleDetailService creates a new ScheduleDetailService instance
func NewScheduleDetailService(
	detailRepo output.ScheduleDetailRepository,
	scheduleRepo output.ScheduleRepository,
) input.ScheduleDetailService {
	return &scheduleDetailService{
		detailRepo:   detailRepo,
		scheduleRepo: scheduleRepo,
	}
}

// BeginTx starts a new database transaction
func (s *scheduleDetailService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.detailRepo.BeginTx(ctx)
}

// CreateDetail creates a new schedule detail (HU6)
func (s *scheduleDetailService) CreateDetail(
	ctx context.Context,
	tx output.Tx,
	detail domain.ScheduleDetail,
) (*domain.ScheduleDetail, error) {
	scheduleDetailLog.Info(logger.LogScheduleDetailServiceCreateStart,
		"schedule_id", detail.ScheduleID,
		"day_of_week", detail.DayOfWeek)

	// 1. Verify schedule exists
	schedule, err := s.scheduleRepo.GetScheduleByID(ctx, detail.ScheduleID)
	if err != nil {
		scheduleDetailLog.Error(logger.LogScheduleDetailServiceScheduleNotFound,
			"schedule_id", detail.ScheduleID, "error", err)
		return nil, domain.ErrScheduleNotFound
	}
	if schedule == nil {
		return nil, domain.ErrScheduleNotFound
	}

	// 2. Validate day of week
	if detail.DayOfWeek == nil || !domain.DayOfWeek(*detail.DayOfWeek).IsValid() {
		scheduleDetailLog.Warn(logger.LogScheduleDetailServiceInvalidDay,
			"day_of_week", detail.DayOfWeek)
		return nil, domain.ErrScheduleDetailInvalidDay
	}

	// 3. Validate time format and range if not closed
	if !detail.IsClosed {
		if err := s.ValidateTimeRange(*detail.OpeningTime, *detail.ClosingTime); err != nil {
			return nil, err
		}

		// 4. Check for time conflicts
		hasConflict, err := s.detailRepo.CheckTimeConflict(
			ctx,
			detail.ScheduleID,
			*detail.DayOfWeek,
			*detail.OpeningTime,
			*detail.ClosingTime,
			"", // No exclude ID for new detail
		)
		if err != nil {
			scheduleDetailLog.Error(logger.LogScheduleDetailServiceConflictCheck,
				"error", err)
			return nil, err
		}
		if hasConflict {
			scheduleDetailLog.Warn(logger.LogScheduleDetailServiceTimeConflict,
				"schedule_id", detail.ScheduleID,
				"day_of_week", *detail.DayOfWeek)
			return nil, domain.ErrScheduleDetailTimeConflict
		}
	}

	// 5. Generate ID and set defaults
	detail.ID = uuid.New().String()
	detail.EntryType = domain.EntryTypeRegular
	detail.Active = true
	detail.CreatedAt = time.Now()
	detail.UpdatedAt = time.Now()

	// 6. Save detail
	if err := s.detailRepo.SaveScheduleDetail(ctx, tx, detail); err != nil {
		scheduleDetailLog.Error(logger.LogScheduleDetailServiceSaveError,
			"schedule_id", detail.ScheduleID, "error", err)
		return nil, err
	}

	scheduleDetailLog.Info(logger.LogScheduleDetailServiceCreateOK,
		"detail_id", detail.ID,
		"schedule_id", detail.ScheduleID,
		"day_of_week", *detail.DayOfWeek)

	return &detail, nil
}

// GetDetailByID retrieves a schedule detail by ID
func (s *scheduleDetailService) GetDetailByID(ctx context.Context, detailID string) (*domain.ScheduleDetail, error) {
	detail, err := s.detailRepo.GetDetailByID(ctx, detailID)
	if err != nil {
		scheduleDetailLog.Error(logger.LogScheduleDetailServiceGetError,
			"detail_id", detailID, "error", err)
		return nil, err
	}
	if detail == nil {
		return nil, domain.ErrScheduleDetailNotFound
	}
	return detail, nil
}

// GetDetailsByScheduleID retrieves all schedule details for a schedule (HU9)
func (s *scheduleDetailService) GetDetailsByScheduleID(
	ctx context.Context,
	scheduleID string,
) ([]domain.ScheduleDetail, error) {
	details, err := s.detailRepo.GetDetailsByScheduleID(ctx, scheduleID)
	if err != nil {
		scheduleDetailLog.Error(logger.LogScheduleDetailServiceListError,
			"schedule_id", scheduleID, "error", err)
		return nil, err
	}

	scheduleDetailLog.Info(logger.LogScheduleDetailServiceListOK,
		"schedule_id", scheduleID,
		"count", len(details))

	return details, nil
}

// UpdateDetail updates an existing schedule detail (HU7)
func (s *scheduleDetailService) UpdateDetail(
	ctx context.Context,
	tx output.Tx,
	detail domain.ScheduleDetail,
) error {
	// 1. Verify detail exists
	existing, err := s.detailRepo.GetDetailByID(ctx, detail.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrScheduleDetailNotFound
	}

	// 2. Validate time range if not closed
	if !detail.IsClosed {
		if err := s.ValidateTimeRange(*detail.OpeningTime, *detail.ClosingTime); err != nil {
			return err
		}

		// 3. Check for time conflicts (excluding this detail)
		hasConflict, err := s.detailRepo.CheckTimeConflict(
			ctx,
			detail.ScheduleID,
			*detail.DayOfWeek,
			*detail.OpeningTime,
			*detail.ClosingTime,
			detail.ID, // Exclude this detail from conflict check
		)
		if err != nil {
			return err
		}
		if hasConflict {
			return domain.ErrScheduleDetailTimeConflict
		}
	}

	// 4. Update detail
	detail.UpdatedAt = time.Now()
	if err := s.detailRepo.UpdateScheduleDetail(ctx, tx, detail); err != nil {
		scheduleDetailLog.Error(logger.LogScheduleDetailServiceUpdateError,
			"detail_id", detail.ID, "error", err)
		return err
	}

	scheduleDetailLog.Info(logger.LogScheduleDetailServiceUpdateOK, "detail_id", detail.ID)
	return nil
}

// DeleteDetail deletes a schedule detail (HU8)
func (s *scheduleDetailService) DeleteDetail(
	ctx context.Context,
	tx output.Tx,
	detailID string,
) error {
	// 1. Verify detail exists
	existing, err := s.detailRepo.GetDetailByID(ctx, detailID)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrScheduleDetailNotFound
	}

	// 2. Delete detail
	if err := s.detailRepo.DeleteScheduleDetail(ctx, tx, detailID); err != nil {
		scheduleDetailLog.Error(logger.LogScheduleDetailServiceDeleteError,
			"detail_id", detailID, "error", err)
		return err
	}

	scheduleDetailLog.Info(logger.LogScheduleDetailServiceDeleteOK, "detail_id", detailID)
	return nil
}

// ValidateTimeRange validates that opening and closing times are in correct format
// and that closing time is after opening time
func (s *scheduleDetailService) ValidateTimeRange(openingTime, closingTime string) error {
	// Validate format
	if !timeRegex.MatchString(openingTime) {
		scheduleDetailLog.Warn(logger.LogScheduleDetailServiceInvalidTime,
			"field", "opening_time", "value", openingTime)
		return domain.ErrScheduleDetailInvalidTime
	}
	if !timeRegex.MatchString(closingTime) {
		scheduleDetailLog.Warn(logger.LogScheduleDetailServiceInvalidTime,
			"field", "closing_time", "value", closingTime)
		return domain.ErrScheduleDetailInvalidTime
	}

	// Parse times and validate range
	opening, _ := time.Parse("15:04", openingTime)
	closing, _ := time.Parse("15:04", closingTime)

	if !closing.After(opening) {
		scheduleDetailLog.Warn(logger.LogScheduleDetailServiceInvalidTimeRange,
			"opening", openingTime, "closing", closingTime)
		return domain.ErrScheduleDetailInvalidTime
	}

	return nil
}

// CheckTimeConflict checks if a time slot conflicts with existing slots for the same day
func (s *scheduleDetailService) CheckTimeConflict(
	ctx context.Context,
	scheduleID string,
	dayOfWeek int,
	openingTime, closingTime string,
	excludeDetailID string,
) (bool, error) {
	return s.detailRepo.CheckTimeConflict(ctx, scheduleID, dayOfWeek, openingTime, closingTime, excludeDetailID)
}
