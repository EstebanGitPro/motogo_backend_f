package services

import (
	"context"
	"errors"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

var scheduleLog logger.Logger = logger.NewSlogLogger()

// scheduleService implements input.ScheduleService
type scheduleService struct {
	scheduleRepo output.ScheduleRepository
	branchRepo   output.BranchRepository
}

// NewScheduleService creates a new ScheduleService instance
func NewScheduleService(scheduleRepo output.ScheduleRepository, branchRepo output.BranchRepository) input.ScheduleService {
	return &scheduleService{
		scheduleRepo: scheduleRepo,
		branchRepo:   branchRepo,
	}
}

// BeginTx starts a new database transaction
func (s *scheduleService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.scheduleRepo.BeginTx(ctx)
}

// CreateSchedule creates a new schedule for a branch (HU30)
func (s *scheduleService) CreateSchedule(ctx context.Context, tx output.Tx, branchID string) (*domain.BranchSchedule, error) {
	scheduleLog.Info(logger.LogScheduleServiceCreateStart, "branch_id", branchID)

	// 1. Verify branch exists
	branch, err := s.branchRepo.GetBranchByID(ctx, branchID)
	if err != nil {
		scheduleLog.Error(logger.LogScheduleServiceBranchNotFound, "branch_id", branchID, "error", err)
		return nil, domain.ErrBranchNotFound
	}
	if branch == nil {
		scheduleLog.Warn(logger.LogScheduleServiceBranchNotFound, "branch_id", branchID)
		return nil, domain.ErrBranchNotFound
	}

	// 2. Check if branch already has a schedule
	existing, err := s.scheduleRepo.GetScheduleByBranchID(ctx, branchID)
	if err != nil {
		scheduleLog.Error(logger.LogScheduleServiceGetError, "branch_id", branchID, "error", err)
		return nil, err
	}
	if existing != nil {
		scheduleLog.Warn(logger.LogScheduleServiceAlreadyExists, "branch_id", branchID)
		return nil, domain.ErrScheduleAlreadyExists
	}

	// 3. Create new schedule
	schedule := domain.BranchSchedule{
		BranchID:  branchID,
		Active:    true,                                // Default to active
		StartDate: time.Now().Truncate(24 * time.Hour), // Default to today
		EndDate:   nil,                                 // Indefinite by default
	}
	schedule.SetID()

	// 4. Save schedule
	if err := s.scheduleRepo.SaveSchedule(ctx, tx, schedule); err != nil {
		scheduleLog.Error(logger.LogScheduleServiceSaveError, "branch_id", branchID, "error", err)
		return nil, err
	}

	scheduleLog.Info(logger.LogScheduleServiceCreateOK, "schedule_id", schedule.ID, "branch_id", branchID)
	return &schedule, nil
}

// GetScheduleByBranchID retrieves schedule for a branch (HU32)
func (s *scheduleService) GetScheduleByBranchID(ctx context.Context, branchID string) (*domain.BranchSchedule, error) {
	schedule, err := s.scheduleRepo.GetScheduleByBranchID(ctx, branchID)
	if err != nil {
		scheduleLog.Error(logger.LogScheduleServiceGetError, "branch_id", branchID, "error", err)
		return nil, err
	}
	if schedule == nil {
		return nil, domain.ErrScheduleNotFound
	}
	scheduleLog.Info(logger.LogScheduleServiceGetOK, "schedule_id", schedule.ID, "branch_id", branchID)
	return schedule, nil
}

// GetScheduleByID retrieves schedule by ID
func (s *scheduleService) GetScheduleByID(ctx context.Context, scheduleID string) (*domain.BranchSchedule, error) {
	schedule, err := s.scheduleRepo.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		if errors.Is(err, domain.ErrScheduleNotFound) {
			return nil, err
		}
		scheduleLog.Error(logger.LogScheduleServiceGetError, "schedule_id", scheduleID, "error", err)
		return nil, err
	}
	return schedule, nil
}

// UpdateSchedule updates an existing schedule (HU31)
func (s *scheduleService) UpdateSchedule(ctx context.Context, tx output.Tx, schedule domain.BranchSchedule) error {
	// 1. Verify schedule exists
	existing, err := s.scheduleRepo.GetScheduleByID(ctx, schedule.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrScheduleNotFound
	}

	// 2. Update schedule
	if err := s.scheduleRepo.UpdateSchedule(ctx, tx, schedule); err != nil {
		scheduleLog.Error(logger.LogScheduleServiceUpdateError, "schedule_id", schedule.ID, "error", err)
		return err
	}

	scheduleLog.Info(logger.LogScheduleServiceUpdateOK, "schedule_id", schedule.ID)
	return nil
}

// DeleteSchedule deletes a schedule (HU33)
func (s *scheduleService) DeleteSchedule(ctx context.Context, tx output.Tx, scheduleID string) error {
	// 1. Verify schedule exists
	existing, err := s.scheduleRepo.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrScheduleNotFound
	}

	// 2. Delete schedule (cascade will remove schedule_details)
	if err := s.scheduleRepo.DeleteSchedule(ctx, tx, scheduleID); err != nil {
		scheduleLog.Error(logger.LogScheduleServiceDeleteError, "schedule_id", scheduleID, "error", err)
		return err
	}

	scheduleLog.Info(logger.LogScheduleServiceDeleteOK, "schedule_id", scheduleID)
	return nil
}

// ActivateSchedule activates a schedule (HU34)
func (s *scheduleService) ActivateSchedule(ctx context.Context, tx output.Tx, scheduleID string) error {
	return s.setActive(ctx, tx, scheduleID, true)
}

// DeactivateSchedule deactivates a schedule (HU35)
func (s *scheduleService) DeactivateSchedule(ctx context.Context, tx output.Tx, scheduleID string) error {
	return s.setActive(ctx, tx, scheduleID, false)
}

// setActive is a helper function for activate/deactivate
func (s *scheduleService) setActive(ctx context.Context, tx output.Tx, scheduleID string, active bool) error {
	// 1. Verify schedule exists
	existing, err := s.scheduleRepo.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrScheduleNotFound
	}

	// 2. Set active status
	if err := s.scheduleRepo.SetActive(ctx, tx, scheduleID, active); err != nil {
		scheduleLog.Error(logger.LogScheduleServiceActivateError, "schedule_id", scheduleID, "active", active, "error", err)
		return err
	}

	scheduleLog.Info(logger.LogScheduleServiceActivateOK, "schedule_id", scheduleID, "active", active)
	return nil
}
