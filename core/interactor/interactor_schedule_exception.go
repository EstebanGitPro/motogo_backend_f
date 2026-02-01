package interactor

import (
	"context"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

var scheduleExceptionInteractorLog = logger.NewSlogLogger()

// ScheduleExceptionInteractor orchestrates schedule exception operations (HU20-25)
type ScheduleExceptionInteractor struct {
	detailService   input.ScheduleDetailService
	scheduleService input.ScheduleService
	branchService   input.BranchService
}

// NewScheduleExceptionInteractor creates a new ScheduleExceptionInteractor
func NewScheduleExceptionInteractor(
	detailService input.ScheduleDetailService,
	scheduleService input.ScheduleService,
	branchService input.BranchService,
) *ScheduleExceptionInteractor {
	return &ScheduleExceptionInteractor{
		detailService:   detailService,
		scheduleService: scheduleService,
		branchService:   branchService,
	}
}

// CreateException orchestrates schedule exception creation (HU20)
func (i *ScheduleExceptionInteractor) CreateException(
	ctx context.Context,
	exception domain.ScheduleDetail,
	representativeID, branchID string,
) (*domain.ScheduleDetail, error) {
	scheduleExceptionInteractorLog.Info(logger.LogScheduleDetailInteractorCreateStart,
		"schedule_id", exception.ScheduleID,
		"exception_start_date", exception.ExceptionStartDate,
		"representative_id", representativeID)

	// 1. Verify ownership of branch
	branch, err := i.branchService.GetBranchByID(ctx, branchID)
	if err != nil {
		scheduleExceptionInteractorLog.Error(logger.LogScheduleDetailInteractorBranchError, "error", err)
		return nil, domain.ErrBranchNotFound
	}
	if branch.RepresentativeID != representativeID {
		scheduleExceptionInteractorLog.Warn(logger.LogScheduleDetailInteractorOwnershipError,
			"branch_id", branchID, "representative_id", representativeID)
		return nil, domain.ErrForbidden
	}

	// 2. Begin transaction
	tx, err := i.detailService.BeginTx(ctx)
	if err != nil {
		scheduleExceptionInteractorLog.Error(logger.LogScheduleDetailInteractorTxError, "error", err)
		return nil, err
	}

	// 3. Create exception via service
	createdException, err := i.detailService.CreateException(ctx, tx, exception)
	if err != nil {
		tx.Rollback()
		scheduleExceptionInteractorLog.Error(logger.LogScheduleDetailInteractorCreateError,
			"schedule_id", exception.ScheduleID, "error", err)
		return nil, err
	}

	// 4. Commit transaction
	if err := tx.Commit(); err != nil {
		scheduleExceptionInteractorLog.Error(logger.LogScheduleDetailInteractorCommitError, "error", err)
		return nil, err
	}

	scheduleExceptionInteractorLog.Info(logger.LogScheduleDetailInteractorCreateOK,
		"exception_id", createdException.ID,
		"schedule_id", exception.ScheduleID)

	return createdException, nil
}

// ListExceptions retrieves all exceptions for a schedule (HU23)
func (i *ScheduleExceptionInteractor) ListExceptions(
	ctx context.Context,
	scheduleID string,
) ([]domain.ScheduleDetail, error) {
	exceptions, err := i.detailService.GetExceptionsByScheduleID(ctx, scheduleID)
	if err != nil {
		scheduleExceptionInteractorLog.Error(logger.LogScheduleDetailInteractorListError,
			"schedule_id", scheduleID, "error", err)
		return nil, err
	}

	scheduleExceptionInteractorLog.Info(logger.LogScheduleDetailInteractorListOK,
		"schedule_id", scheduleID, "count", len(exceptions))

	return exceptions, nil
}

// GetException retrieves a schedule exception by ID
func (i *ScheduleExceptionInteractor) GetException(ctx context.Context, exceptionID string) (*domain.ScheduleDetail, error) {
	exception, err := i.detailService.GetExceptionByID(ctx, exceptionID)
	if err != nil {
		return nil, err
	}
	return exception, nil
}

// UpdateException orchestrates schedule exception update (HU21)
func (i *ScheduleExceptionInteractor) UpdateException(
	ctx context.Context,
	exception domain.ScheduleDetail,
	representativeID string,
) error {
	// 1. Get existing exception
	existing, err := i.detailService.GetExceptionByID(ctx, exception.ID)
	if err != nil {
		return err
	}

	// 2. Get schedule and verify ownership via branch
	schedule, err := i.scheduleService.GetScheduleByID(ctx, existing.ScheduleID)
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

	// 3. Begin transaction
	tx, err := i.detailService.BeginTx(ctx)
	if err != nil {
		return err
	}

	// 4. Update exception
	if err := i.detailService.UpdateException(ctx, tx, exception); err != nil {
		tx.Rollback()
		return err
	}

	// 5. Commit
	if err := tx.Commit(); err != nil {
		scheduleExceptionInteractorLog.Error(logger.LogScheduleDetailInteractorCommitError, "error", err)
		return err
	}

	return nil
}

// DeleteException orchestrates schedule exception deletion (HU22)
func (i *ScheduleExceptionInteractor) DeleteException(
	ctx context.Context,
	exceptionID, representativeID string,
) error {
	// 1. Get existing exception
	existing, err := i.detailService.GetExceptionByID(ctx, exceptionID)
	if err != nil {
		return err
	}

	// 2. Get schedule and verify ownership via branch
	schedule, err := i.scheduleService.GetScheduleByID(ctx, existing.ScheduleID)
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

	// 3. Begin transaction
	tx, err := i.detailService.BeginTx(ctx)
	if err != nil {
		return err
	}

	// 4. Delete exception
	if err := i.detailService.DeleteException(ctx, tx, exceptionID); err != nil {
		tx.Rollback()
		return err
	}

	// 5. Commit
	if err := tx.Commit(); err != nil {
		scheduleExceptionInteractorLog.Error(logger.LogScheduleDetailInteractorCommitError, "error", err)
		return err
	}

	return nil
}

// ActivateException activates a schedule exception (HU24)
func (i *ScheduleExceptionInteractor) ActivateException(
	ctx context.Context,
	exceptionID, representativeID string,
) error {
	return i.setExceptionActive(ctx, exceptionID, representativeID, true)
}

// DeactivateException deactivates a schedule exception (HU25)
func (i *ScheduleExceptionInteractor) DeactivateException(
	ctx context.Context,
	exceptionID, representativeID string,
) error {
	return i.setExceptionActive(ctx, exceptionID, representativeID, false)
}

// setExceptionActive is a helper to activate/deactivate an exception
func (i *ScheduleExceptionInteractor) setExceptionActive(
	ctx context.Context,
	exceptionID, representativeID string,
	active bool,
) error {
	// 1. Get existing exception
	existing, err := i.detailService.GetExceptionByID(ctx, exceptionID)
	if err != nil {
		return err
	}

	// 2. Get schedule and verify ownership via branch
	schedule, err := i.scheduleService.GetScheduleByID(ctx, existing.ScheduleID)
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

	// 3. Begin transaction
	tx, err := i.detailService.BeginTx(ctx)
	if err != nil {
		return err
	}

	// 4. Update exception with new active status
	existing.Active = active
	existing.UpdatedAt = time.Now()
	if err := i.detailService.UpdateException(ctx, tx, *existing); err != nil {
		tx.Rollback()
		return err
	}

	// 5. Commit
	if err := tx.Commit(); err != nil {
		scheduleExceptionInteractorLog.Error(logger.LogScheduleDetailInteractorCommitError, "error", err)
		return err
	}

	return nil
}
