package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

var scheduleInteractorLog = logger.NewSlogLogger()

// ScheduleInteractor orchestrates schedule operations (HU30-35)
type ScheduleInteractor struct {
	scheduleService input.ScheduleService
	branchService   input.BranchService
}

// NewScheduleInteractor creates a new ScheduleInteractor
func NewScheduleInteractor(scheduleService input.ScheduleService, branchService input.BranchService) *ScheduleInteractor {
	return &ScheduleInteractor{
		scheduleService: scheduleService,
		branchService:   branchService,
	}
}

// CreateSchedule orchestrates schedule creation for a branch (HU30)
func (i *ScheduleInteractor) CreateSchedule(ctx context.Context, branchID, representativeID string) (*domain.BranchSchedule, error) {
	scheduleInteractorLog.Info(logger.LogScheduleInteractorCreateStart, "branch_id", branchID, "representative_id", representativeID)

	// 1. Verify ownership of branch
	branch, err := i.branchService.GetBranchByID(ctx, branchID)
	if err != nil {
		scheduleInteractorLog.Error(logger.LogScheduleInteractorCreateError, "error", "branch not found", "branch_id", branchID)
		return nil, domain.ErrBranchNotFound
	}
	if branch.RepresentativeID != representativeID {
		scheduleInteractorLog.Warn(logger.LogBranchInteractorOwnershipError, "branch_id", branchID, "representative_id", representativeID)
		return nil, domain.ErrForbidden
	}

	// 2. Begin transaction
	tx, err := i.scheduleService.BeginTx(ctx)
	if err != nil {
		scheduleInteractorLog.Error(logger.LogScheduleInteractorTxError, "error", err)
		return nil, err
	}

	// 3. Create schedule via service
	schedule, err := i.scheduleService.CreateSchedule(ctx, tx, branchID)
	if err != nil {
		tx.Rollback()
		scheduleInteractorLog.Error(logger.LogScheduleInteractorCreateError, "branch_id", branchID, "error", err)
		return nil, err
	}

	// 4. Commit transaction
	if err := tx.Commit(); err != nil {
		scheduleInteractorLog.Error(logger.LogScheduleInteractorCommitError, "error", err)
		return nil, err
	}

	scheduleInteractorLog.Info(logger.LogScheduleInteractorCreateComplete, "schedule_id", schedule.ID, "branch_id", branchID)
	return schedule, nil
}

// GetScheduleByBranchID retrieves schedule for a branch (HU32)
func (i *ScheduleInteractor) GetScheduleByBranchID(ctx context.Context, branchID string) (*domain.BranchSchedule, error) {
	schedule, err := i.scheduleService.GetScheduleByBranchID(ctx, branchID)
	if err != nil {
		scheduleInteractorLog.Error(logger.LogScheduleInteractorGetError, "branch_id", branchID, "error", err)
		return nil, err
	}
	scheduleInteractorLog.Info(logger.LogScheduleInteractorGetOK, "schedule_id", schedule.ID, "branch_id", branchID)
	return schedule, nil
}

// GetScheduleByID retrieves schedule by ID
func (i *ScheduleInteractor) GetScheduleByID(ctx context.Context, scheduleID string) (*domain.BranchSchedule, error) {
	return i.scheduleService.GetScheduleByID(ctx, scheduleID)
}

// UpdateSchedule orchestrates schedule update (HU31)
func (i *ScheduleInteractor) UpdateSchedule(ctx context.Context, schedule domain.BranchSchedule, representativeID string) error {
	scheduleInteractorLog.Info(logger.LogScheduleInteractorUpdateStart, "schedule_id", schedule.ID)

	// 1. Verify ownership via branch
	branch, err := i.branchService.GetBranchByID(ctx, schedule.BranchID)
	if err != nil {
		return domain.ErrBranchNotFound
	}
	if branch.RepresentativeID != representativeID {
		return domain.ErrForbidden
	}

	// 2. Begin transaction
	tx, err := i.scheduleService.BeginTx(ctx)
	if err != nil {
		return err
	}

	// 3. Update schedule with all fields
	if err := i.scheduleService.UpdateSchedule(ctx, tx, schedule); err != nil {
		tx.Rollback()
		return err
	}

	// 4. Commit
	if err := tx.Commit(); err != nil {
		scheduleInteractorLog.Error(logger.LogScheduleInteractorCommitError, "error", err)
		return err
	}

	scheduleInteractorLog.Info(logger.LogScheduleInteractorUpdateComplete, "schedule_id", schedule.ID)
	return nil
}

// DeleteSchedule orchestrates schedule deletion (HU33)
func (i *ScheduleInteractor) DeleteSchedule(ctx context.Context, scheduleID, representativeID string) error {
	scheduleInteractorLog.Info(logger.LogScheduleInteractorDeleteStart, "schedule_id", scheduleID)

	// 1. Get schedule and verify ownership
	schedule, err := i.scheduleService.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		return err
	}

	branch, err := i.branchService.GetBranchByID(ctx, schedule.BranchID)
	if err != nil {
		return domain.ErrBranchNotFound
	}
	if branch.RepresentativeID != representativeID {
		return domain.ErrForbidden
	}

	// 2. Begin transaction
	tx, err := i.scheduleService.BeginTx(ctx)
	if err != nil {
		return err
	}

	// 3. Delete schedule
	if err := i.scheduleService.DeleteSchedule(ctx, tx, scheduleID); err != nil {
		tx.Rollback()
		return err
	}

	// 4. Commit
	if err := tx.Commit(); err != nil {
		scheduleInteractorLog.Error(logger.LogScheduleInteractorCommitError, "error", err)
		return err
	}

	scheduleInteractorLog.Info(logger.LogScheduleInteractorDeleteComplete, "schedule_id", scheduleID)
	return nil
}

// ActivateSchedule orchestrates schedule activation (HU34)
func (i *ScheduleInteractor) ActivateSchedule(ctx context.Context, scheduleID, representativeID string) error {
	return i.setActive(ctx, scheduleID, representativeID, true)
}

// DeactivateSchedule orchestrates schedule deactivation (HU35)
func (i *ScheduleInteractor) DeactivateSchedule(ctx context.Context, scheduleID, representativeID string) error {
	return i.setActive(ctx, scheduleID, representativeID, false)
}

// setActive is a helper for activate/deactivate with ownership check
func (i *ScheduleInteractor) setActive(ctx context.Context, scheduleID, representativeID string, active bool) error {
	// 1. Get schedule and verify ownership
	schedule, err := i.scheduleService.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		return err
	}

	branch, err := i.branchService.GetBranchByID(ctx, schedule.BranchID)
	if err != nil {
		return domain.ErrBranchNotFound
	}
	if branch.RepresentativeID != representativeID {
		return domain.ErrForbidden
	}

	// 2. Begin transaction
	tx, err := i.scheduleService.BeginTx(ctx)
	if err != nil {
		return err
	}

	// 3. Set active status
	var serviceErr error
	if active {
		serviceErr = i.scheduleService.ActivateSchedule(ctx, tx, scheduleID)
	} else {
		serviceErr = i.scheduleService.DeactivateSchedule(ctx, tx, scheduleID)
	}
	if serviceErr != nil {
		tx.Rollback()
		return serviceErr
	}

	// 4. Commit
	if err := tx.Commit(); err != nil {
		scheduleInteractorLog.Error(logger.LogScheduleInteractorCommitError, "error", err)
		return err
	}

	return nil
}
