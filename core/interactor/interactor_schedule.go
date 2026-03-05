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
func (i *ScheduleInteractor) CreateSchedule(ctx context.Context, branchID, representativeID string) (result *domain.BranchSchedule, err error) {
	scheduleInteractorLog.Info(logger.LogScheduleInteractorCreateStart, "branch_id", branchID, "representative_id", representativeID)

	// 1. Verify ownership of branch
	if ownerErr := verifyBranchOwnership(ctx, i.branchService, branchID, representativeID,
		scheduleInteractorLog, logger.LogScheduleInteractorCreateError, logger.LogBranchInteractorOwnershipError); ownerErr != nil {
		return nil, ownerErr
	}

	// 2. Begin transaction
	tx, txErr := i.scheduleService.BeginTx(ctx)
	if txErr != nil {
		scheduleInteractorLog.Error(logger.LogScheduleInteractorTxError, "error", txErr)
		return nil, txErr
	}

	defer func() { _ = tx.Rollback() }()

	// 3. Create schedule via service
	result, err = i.scheduleService.CreateSchedule(ctx, tx, branchID)
	if err != nil {
		scheduleInteractorLog.Error(logger.LogScheduleInteractorCreateError, "branch_id", branchID, "error", err)
		return nil, err
	}

	// 4. Commit transaction
	if err = tx.Commit(); err != nil {
		scheduleInteractorLog.Error(logger.LogScheduleInteractorCommitError, "error", err)
		return nil, err
	}

	scheduleInteractorLog.Info(logger.LogScheduleInteractorCreateComplete, "schedule_id", result.ID, "branch_id", branchID)

	return result, nil
}

// GetScheduleByBranchID retrieves schedule for a branch with ownership check (HU32, protected)
func (i *ScheduleInteractor) GetScheduleByBranchID(ctx context.Context, branchID, representativeID string) (*domain.BranchSchedule, error) {
	// 1. Verify ownership of branch
	if ownerErr := verifyBranchOwnership(ctx, i.branchService, branchID, representativeID,
		scheduleInteractorLog, logger.LogScheduleInteractorGetError, logger.LogBranchInteractorOwnershipError); ownerErr != nil {
		return nil, ownerErr
	}

	// 2. Get schedule
	schedule, err := i.scheduleService.GetScheduleByBranchID(ctx, branchID)
	if err != nil {
		scheduleInteractorLog.Error(logger.LogScheduleInteractorGetError, "branch_id", branchID, "error", err)
		return nil, err
	}
	scheduleInteractorLog.Info(logger.LogScheduleInteractorGetOK, "schedule_id", schedule.ID, "branch_id", branchID)
	return schedule, nil
}

// GetScheduleByBranchIDPublic retrieves schedule for a branch without ownership check (for public endpoints)
func (i *ScheduleInteractor) GetScheduleByBranchIDPublic(ctx context.Context, branchID string) (*domain.BranchSchedule, error) {
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
func (i *ScheduleInteractor) UpdateSchedule(ctx context.Context, schedule domain.BranchSchedule, representativeID string) (err error) {
	scheduleInteractorLog.Info(logger.LogScheduleInteractorUpdateStart, "schedule_id", schedule.ID)

	// 1. Verify ownership via branch
	branch, branchErr := i.branchService.GetBranchByID(ctx, schedule.BranchID)
	if branchErr != nil {
		return domain.ErrBranchNotFound
	}
	if branch.RepresentativeID != representativeID {
		return domain.ErrForbidden
	}

	// 2. Begin transaction
	tx, txErr := i.scheduleService.BeginTx(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() { _ = tx.Rollback() }()

	// 3. Update schedule with all fields
	if err = i.scheduleService.UpdateSchedule(ctx, tx, schedule); err != nil {
		return err
	}

	// 4. Commit
	if err = tx.Commit(); err != nil {
		scheduleInteractorLog.Error(logger.LogScheduleInteractorCommitError, "error", err)
		return err
	}

	scheduleInteractorLog.Info(logger.LogScheduleInteractorUpdateComplete, "schedule_id", schedule.ID)

	return nil
}

// DeleteSchedule orchestrates schedule deletion (HU33)
func (i *ScheduleInteractor) DeleteSchedule(ctx context.Context, scheduleID, representativeID string) (err error) {
	scheduleInteractorLog.Info(logger.LogScheduleInteractorDeleteStart, "schedule_id", scheduleID)

	// 1. Verify ownership via schedule → branch
	if ownerErr := verifyScheduleOwnership(ctx, i.scheduleService, i.branchService,
		scheduleID, representativeID); ownerErr != nil {
		return ownerErr
	}

	// 2. Begin transaction
	tx, txErr := i.scheduleService.BeginTx(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() { _ = tx.Rollback() }()

	// 3. Delete schedule
	if err = i.scheduleService.DeleteSchedule(ctx, tx, scheduleID); err != nil {
		return err
	}

	// 4. Commit
	if err = tx.Commit(); err != nil {
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
func (i *ScheduleInteractor) setActive(ctx context.Context, scheduleID, representativeID string, active bool) (err error) {
	// 1. Verify ownership via schedule → branch
	if ownerErr := verifyScheduleOwnership(ctx, i.scheduleService, i.branchService,
		scheduleID, representativeID); ownerErr != nil {
		return ownerErr
	}

	// 2. Begin transaction
	tx, txErr := i.scheduleService.BeginTx(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() { _ = tx.Rollback() }()

	// 3. Set active status
	if active {
		err = i.scheduleService.ActivateSchedule(ctx, tx, scheduleID)
	} else {
		err = i.scheduleService.DeactivateSchedule(ctx, tx, scheduleID)
	}
	if err != nil {
		return err
	}

	// 4. Commit
	if err = tx.Commit(); err != nil {
		scheduleInteractorLog.Error(logger.LogScheduleInteractorCommitError, "error", err)
		return err
	}

	return nil
}
