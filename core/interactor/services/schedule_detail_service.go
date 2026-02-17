package services

import (
	"context"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/constants"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

var scheduleDetailLog logger.Logger = logger.NewSlogLogger()

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

	// 3. Validation R1/R2: Check if day is already marked as closed
	dayIsClosed, err := s.detailRepo.CheckDayIsClosed(ctx, detail.ScheduleID, *detail.DayOfWeek, "")
	if err != nil {
		return nil, err
	}
	if dayIsClosed {
		scheduleDetailLog.Warn(logger.LogScheduleDetailServiceTimeConflict,
			"schedule_id", detail.ScheduleID,
			"day_of_week", *detail.DayOfWeek,
			"reason", "day_already_closed")
		return nil, domain.ErrScheduleDetailDayAlreadyClosed
	}

	// 4. Validation R3: If trying to set is_closed=true, check if day has time slots
	if detail.IsClosed {
		dayHasSlots, err := s.detailRepo.CheckDayHasTimeSlots(ctx, detail.ScheduleID, *detail.DayOfWeek, "")
		if err != nil {
			return nil, err
		}
		if dayHasSlots {
			scheduleDetailLog.Warn(logger.LogScheduleDetailServiceTimeConflict,
				"schedule_id", detail.ScheduleID,
				"day_of_week", *detail.DayOfWeek,
				"reason", "day_has_time_slots")
			return nil, domain.ErrScheduleDetailDayHasSlots
		}
	}

	// 5. Validate time format and range if not closed
	if !detail.IsClosed {
		if err := s.ValidateTimeRange(*detail.OpeningTime, *detail.ClosingTime); err != nil {
			return nil, err
		}

		// 6. Check for time conflicts
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
	detail.SetID()
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

	// 2. Preserve immutable fields from existing (not provided in update request)
	detail.ScheduleID = existing.ScheduleID
	detail.DayOfWeek = existing.DayOfWeek
	detail.EntryType = existing.EntryType

	// 3. Validate time range if not closed
	if !detail.IsClosed {
		if err := s.ValidateTimeRange(*detail.OpeningTime, *detail.ClosingTime); err != nil {
			return err
		}

		// 4. Check for time conflicts (excluding this detail)
		hasConflict, err := s.detailRepo.CheckTimeConflict(
			ctx,
			existing.ScheduleID,
			*existing.DayOfWeek,
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

	// 5. Update detail
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
	if !constants.TimeRegex.MatchString(openingTime) {
		scheduleDetailLog.Warn(logger.LogScheduleDetailServiceInvalidTime,
			"field", "opening_time", "value", openingTime)
		return domain.ErrScheduleDetailInvalidTime
	}
	if !constants.TimeRegex.MatchString(closingTime) {
		scheduleDetailLog.Warn(logger.LogScheduleDetailServiceInvalidTime,
			"field", "closing_time", "value", closingTime)
		return domain.ErrScheduleDetailInvalidTime
	}

	// Parse times and validate range (try HH:mm:ss first, then HH:mm)
	opening := parseTime(openingTime)
	closing := parseTime(closingTime)

	if !closing.After(opening) {
		scheduleDetailLog.Warn(logger.LogScheduleDetailServiceInvalidTimeRange,
			"opening", openingTime, "closing", closingTime)
		return domain.ErrScheduleDetailInvalidTime
	}

	return nil
}

// parseTime tries to parse time in HH:mm:ss format first, then HH:mm
func parseTime(timeStr string) time.Time {
	// Try HH:mm:ss format first (from DB)
	t, err := time.Parse(constants.TimeFormatLong, timeStr)
	if err == nil {
		return t
	}
	// Fallback to HH:mm format (from API)
	t, _ = time.Parse(constants.TimeFormatShort, timeStr)
	return t
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

// ============================================
// Schedule Exception Methods (HU20-25)
// ============================================

// CreateException creates a new schedule exception (HU20)
func (s *scheduleDetailService) CreateException(
	ctx context.Context,
	tx output.Tx,
	exception domain.ScheduleDetail,
) (*domain.ScheduleDetail, error) {
	scheduleDetailLog.Info(logger.LogScheduleDetailServiceCreateStart,
		"schedule_id", exception.ScheduleID,
		"exception_start_date", exception.ExceptionStartDate)

	// 1. Verify schedule exists
	schedule, err := s.scheduleRepo.GetScheduleByID(ctx, exception.ScheduleID)
	if err != nil {
		scheduleDetailLog.Error(logger.LogScheduleDetailServiceScheduleNotFound,
			"schedule_id", exception.ScheduleID, "error", err)
		return nil, domain.ErrScheduleNotFound
	}
	if schedule == nil {
		return nil, domain.ErrScheduleNotFound
	}

	// 2. Validate exception start date is provided
	if exception.ExceptionStartDate == nil {
		return nil, domain.ErrScheduleExceptionDatePast
	}

	// 3. Validate exception start date is not in the past
	today := time.Now().Truncate(24 * time.Hour)
	exceptionDay := exception.ExceptionStartDate.Truncate(24 * time.Hour)
	if exceptionDay.Before(today) {
		scheduleDetailLog.Warn(logger.LogScheduleDetailServiceInvalidDay,
			"exception_start_date", exception.ExceptionStartDate)
		return nil, domain.ErrScheduleExceptionDatePast
	}

	// 4. If end date not set, use start date
	if exception.ExceptionEndDate == nil {
		exception.ExceptionEndDate = exception.ExceptionStartDate
	}

	// 5. Check for date conflict (overlapping dates) - using FOR UPDATE lock to prevent race conditions
	// The tx parameter ensures this query runs within the transaction context
	existingExceptions, err := s.detailRepo.GetExceptionsByScheduleIDForUpdate(ctx, tx, exception.ScheduleID)
	if err != nil {
		scheduleDetailLog.Error(logger.LogScheduleDetailDebugGetError, "error", err)
		return nil, err
	}

	scheduleDetailLog.Info(logger.LogScheduleDetailDebugExplicitCheck,
		"schedule_id", exception.ScheduleID,
		"new_start_date", exception.ExceptionStartDate.Format(constants.DateFormat),
		"new_end_date", exception.ExceptionEndDate.Format(constants.DateFormat),
		"new_is_closed", exception.IsClosed,
		"existing_exceptions_count", len(existingExceptions))

	// Check for overlapping dates with any existing exception
	if err := s.checkDateOverlap(exception, existingExceptions); err != nil {
		return nil, err
	}

	// 6. Validation E1: If is_closed=true, check if day is already closed in REGULAR schedule (redundant)
	if err := s.checkClosedDayRedundancy(ctx, exception); err != nil {
		return nil, err
	}

	// 7. Validate time format if not closed
	if !exception.IsClosed {
		if exception.OpeningTime == nil || exception.ClosingTime == nil {
			return nil, domain.ErrScheduleExceptionInvalidTime
		}
		if err := s.ValidateTimeRange(*exception.OpeningTime, *exception.ClosingTime); err != nil {
			return nil, domain.ErrScheduleExceptionInvalidTime
		}
	}

	// 7. Generate ID and set defaults
	exception.SetID()
	exception.EntryType = domain.EntryTypeException
	exception.DayOfWeek = nil // Exceptions don't use day_of_week
	exception.Active = true
	exception.CreatedAt = time.Now()
	exception.UpdatedAt = time.Now()

	// 7.5 Validation: If is_closed=true, clear time fields to prevent inconsistent data
	if exception.IsClosed {
		exception.OpeningTime = nil
		exception.ClosingTime = nil
	}

	// 8. Save exception
	if err := s.detailRepo.SaveScheduleDetail(ctx, tx, exception); err != nil {
		scheduleDetailLog.Error(logger.LogScheduleDetailServiceSaveError,
			"schedule_id", exception.ScheduleID, "error", err)
		return nil, err
	}

	scheduleDetailLog.Info(logger.LogScheduleDetailServiceCreateOK,
		"exception_id", exception.ID,
		"schedule_id", exception.ScheduleID,
		"exception_start_date", exception.ExceptionStartDate)

	return &exception, nil
}

// GetExceptionsByScheduleID retrieves all exceptions for a schedule (HU23)
func (s *scheduleDetailService) GetExceptionsByScheduleID(
	ctx context.Context,
	scheduleID string,
) ([]domain.ScheduleDetail, error) {
	exceptions, err := s.detailRepo.GetExceptionsByScheduleID(ctx, scheduleID)
	if err != nil {
		scheduleDetailLog.Error(logger.LogScheduleDetailServiceListError,
			"schedule_id", scheduleID, "error", err)
		return nil, err
	}

	scheduleDetailLog.Info(logger.LogScheduleDetailServiceListOK,
		"schedule_id", scheduleID,
		"exceptions_count", len(exceptions))

	return exceptions, nil
}

// GetExceptionByID retrieves a specific exception by ID
func (s *scheduleDetailService) GetExceptionByID(ctx context.Context, exceptionID string) (*domain.ScheduleDetail, error) {
	exception, err := s.detailRepo.GetExceptionByID(ctx, exceptionID)
	if err != nil {
		scheduleDetailLog.Error(logger.LogScheduleDetailServiceGetError,
			"exception_id", exceptionID, "error", err)
		return nil, err
	}
	if exception == nil {
		return nil, domain.ErrScheduleExceptionNotFound
	}
	return exception, nil
}

// UpdateException updates an existing schedule exception (HU21)
func (s *scheduleDetailService) UpdateException(
	ctx context.Context,
	tx output.Tx,
	exception domain.ScheduleDetail,
) error {
	// 1. Verify exception exists
	existing, err := s.detailRepo.GetExceptionByID(ctx, exception.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrScheduleExceptionNotFound
	}

	// 2. Validate time range if not closed
	if !exception.IsClosed {
		if exception.OpeningTime == nil || exception.ClosingTime == nil {
			return domain.ErrScheduleExceptionInvalidTime
		}
		if err := s.ValidateTimeRange(*exception.OpeningTime, *exception.ClosingTime); err != nil {
			return domain.ErrScheduleExceptionInvalidTime
		}
	}

	// 3. Update exception (preserve existing exception dates and schedule_id)
	exception.ScheduleID = existing.ScheduleID
	exception.ExceptionStartDate = existing.ExceptionStartDate
	exception.ExceptionEndDate = existing.ExceptionEndDate
	exception.EntryType = domain.EntryTypeException
	exception.UpdatedAt = time.Now()

	// 3.5 Validation: If is_closed=true, clear time fields to prevent inconsistent data
	if exception.IsClosed {
		exception.OpeningTime = nil
		exception.ClosingTime = nil
	}

	if err := s.detailRepo.UpdateScheduleDetail(ctx, tx, exception); err != nil {
		scheduleDetailLog.Error(logger.LogScheduleDetailServiceUpdateError,
			"exception_id", exception.ID, "error", err)
		return err
	}

	scheduleDetailLog.Info(logger.LogScheduleDetailServiceUpdateOK, "exception_id", exception.ID)
	return nil
}

// DeleteException deletes a schedule exception (HU22)
func (s *scheduleDetailService) DeleteException(
	ctx context.Context,
	tx output.Tx,
	exceptionID string,
) error {
	// 1. Verify exception exists
	existing, err := s.detailRepo.GetExceptionByID(ctx, exceptionID)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrScheduleExceptionNotFound
	}

	// 2. Delete exception
	if err := s.detailRepo.DeleteScheduleDetail(ctx, tx, exceptionID); err != nil {
		scheduleDetailLog.Error(logger.LogScheduleDetailServiceDeleteError,
			"exception_id", exceptionID, "error", err)
		return err
	}

	scheduleDetailLog.Info(logger.LogScheduleDetailServiceDeleteOK, "exception_id", exceptionID)
	return nil
}

// SetExceptionActive toggles the active status of an exception (HU24/HU25)
// This encapsulates the state mutation logic that was previously in the interactor
func (s *scheduleDetailService) SetExceptionActive(
	ctx context.Context,
	tx output.Tx,
	exceptionID string,
	active bool,
) error {
	// 1. Verify exception exists
	existing, err := s.detailRepo.GetExceptionByID(ctx, exceptionID)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrScheduleExceptionNotFound
	}

	// 2. Apply state mutation
	existing.Active = active
	existing.UpdatedAt = time.Now()

	// 3. Persist
	if err := s.detailRepo.UpdateScheduleDetail(ctx, tx, *existing); err != nil {
		scheduleDetailLog.Error(logger.LogScheduleDetailServiceUpdateError,
			"exception_id", exceptionID, "active", active, "error", err)
		return err
	}

	scheduleDetailLog.Info(logger.LogScheduleDetailServiceUpdateOK,
		"exception_id", exceptionID, "active", active)
	return nil
}

// CheckExceptionDateConflict checks if an exception already exists for the given date range
func (s *scheduleDetailService) CheckExceptionDateConflict(
	ctx context.Context,
	scheduleID, excludeExceptionID, startDate, endDate string,
) (bool, error) {
	return s.detailRepo.CheckExceptionDateConflict(ctx, scheduleID, excludeExceptionID, startDate, endDate)
}

// checkDateOverlap checks if the new exception overlaps with any existing exceptions.
func (s *scheduleDetailService) checkDateOverlap(exception domain.ScheduleDetail, existingExceptions []domain.ScheduleDetail) error {
	for _, existing := range existingExceptions {
		if existing.ExceptionStartDate == nil || existing.ExceptionEndDate == nil {
			continue
		}

		// Use date-only strings to compare (YYYY-MM-DD) - avoids timezone truncation issues
		existingStartStr := existing.ExceptionStartDate.Format(constants.DateFormat)
		existingEndStr := existing.ExceptionEndDate.Format(constants.DateFormat)
		newStartStr := exception.ExceptionStartDate.Format(constants.DateFormat)
		newEndStr := exception.ExceptionEndDate.Format(constants.DateFormat)

		// Overlap condition: existing.start <= new.end AND existing.end >= new.start
		hasOverlap := existingStartStr <= newEndStr && existingEndStr >= newStartStr

		scheduleDetailLog.Info(logger.LogScheduleDetailDebugCheckOverlap,
			"existing_id", existing.ID,
			"existing_start", existingStartStr,
			"existing_end", existingEndStr,
			"new_start", newStartStr,
			"new_end", newEndStr,
			"existing_is_closed", existing.IsClosed,
			"has_overlap", hasOverlap)

		if hasOverlap {
			scheduleDetailLog.Warn(logger.LogScheduleDetailServiceTimeConflict,
				"schedule_id", exception.ScheduleID,
				"exception_start_date", exception.ExceptionStartDate,
				"conflicting_exception_id", existing.ID)
			return domain.ErrScheduleExceptionDateConflict
		}
	}
	return nil
}

// checkClosedDayRedundancy checks if a closed exception is redundant (day already closed in regular schedule).
func (s *scheduleDetailService) checkClosedDayRedundancy(ctx context.Context, exception domain.ScheduleDetail) error {
	if !exception.IsClosed {
		return nil
	}

	dayOfWeek := int(exception.ExceptionStartDate.Weekday())
	if dayOfWeek == 0 {
		dayOfWeek = 7 // Sunday is 7 in ISO format
	}

	isRedundant, err := s.detailRepo.CheckExceptionIsRedundant(ctx, exception.ScheduleID, dayOfWeek)
	if err != nil {
		return err
	}
	if isRedundant {
		scheduleDetailLog.Warn(logger.LogScheduleDetailServiceTimeConflict,
			"schedule_id", exception.ScheduleID,
			"exception_start_date", exception.ExceptionStartDate,
			"reason", "exception_redundant_day_already_closed")
		return domain.ErrScheduleExceptionRedundant
	}

	return nil
}
