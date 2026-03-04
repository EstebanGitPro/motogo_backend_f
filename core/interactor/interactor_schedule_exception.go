package interactor

import (
	"context"

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
) (result *domain.ScheduleDetail, err error) {
	scheduleExceptionInteractorLog.Info(logger.LogScheduleDetailInteractorCreateStart,
		"schedule_id", exception.ScheduleID,
		"exception_start_date", exception.ExceptionStartDate,
		"representative_id", representativeID)

	// 1. Verify ownership of branch
	if ownerErr := verifyBranchOwnership(ctx, i.branchService, branchID, representativeID,
		scheduleExceptionInteractorLog, logger.LogScheduleDetailInteractorBranchError, logger.LogScheduleDetailInteractorOwnershipError); ownerErr != nil {
		return nil, ownerErr
	}

	// 2. Begin transaction
	tx, txErr := i.detailService.BeginTx(ctx)
	if txErr != nil {
		scheduleExceptionInteractorLog.Error(logger.LogScheduleDetailInteractorTxError, "error", txErr)
		return nil, txErr
	}

	defer tx.Rollback()

	// 3. Create exception via service
	result, err = i.detailService.CreateException(ctx, tx, exception)
	if err != nil {
		scheduleExceptionInteractorLog.Error(logger.LogScheduleDetailInteractorCreateError,
			"schedule_id", exception.ScheduleID, "error", err)
		return nil, err
	}

	// 4. Commit transaction
	if err = tx.Commit(); err != nil {
		scheduleExceptionInteractorLog.Error(logger.LogScheduleDetailInteractorCommitError, "error", err)
		return nil, err
	}

	scheduleExceptionInteractorLog.Info(logger.LogScheduleDetailInteractorCreateOK,
		"exception_id", result.ID,
		"schedule_id", exception.ScheduleID)

	err = nil
	return result, nil
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
) (err error) {
	// 1. Get existing exception
	existing, existErr := i.detailService.GetExceptionByID(ctx, exception.ID)
	if existErr != nil {
		return existErr
	}

	// 2. Verify ownership via schedule → branch
	if ownerErr := verifyScheduleOwnership(ctx, i.scheduleService, i.branchService,
		existing.ScheduleID, representativeID); ownerErr != nil {
		return ownerErr
	}

	// 3. Begin transaction
	tx, txErr := i.detailService.BeginTx(ctx)
	if txErr != nil {
		return txErr
	}

	defer tx.Rollback()

	// 4. Update exception
	if err = i.detailService.UpdateException(ctx, tx, exception); err != nil {
		return err
	}

	// 5. Commit
	if err = tx.Commit(); err != nil {
		scheduleExceptionInteractorLog.Error(logger.LogScheduleDetailInteractorCommitError, "error", err)
		return err
	}

	err = nil
	return nil
}

// DeleteException orchestrates schedule exception deletion (HU22)
func (i *ScheduleExceptionInteractor) DeleteException(
	ctx context.Context,
	exceptionID, representativeID string,
) (err error) {
	// 1. Get existing exception
	existing, existErr := i.detailService.GetExceptionByID(ctx, exceptionID)
	if existErr != nil {
		return existErr
	}

	// 2. Verify ownership via schedule → branch
	if ownerErr := verifyScheduleOwnership(ctx, i.scheduleService, i.branchService,
		existing.ScheduleID, representativeID); ownerErr != nil {
		return ownerErr
	}

	// 3. Begin transaction
	tx, txErr := i.detailService.BeginTx(ctx)
	if txErr != nil {
		return txErr
	}

	defer tx.Rollback()

	// 4. Delete exception
	if err = i.detailService.DeleteException(ctx, tx, exceptionID); err != nil {
		return err
	}

	// 5. Commit
	if err = tx.Commit(); err != nil {
		scheduleExceptionInteractorLog.Error(logger.LogScheduleDetailInteractorCommitError, "error", err)
		return err
	}

	err = nil
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
) (err error) {
	// 1. Get existing exception
	existing, existErr := i.detailService.GetExceptionByID(ctx, exceptionID)
	if existErr != nil {
		return existErr
	}

	// 2. Verify ownership via schedule → branch
	if ownerErr := verifyScheduleOwnership(ctx, i.scheduleService, i.branchService,
		existing.ScheduleID, representativeID); ownerErr != nil {
		return ownerErr
	}

	// 3. Begin transaction
	tx, txErr := i.detailService.BeginTx(ctx)
	if txErr != nil {
		return txErr
	}

	defer tx.Rollback()

	// 4. Delegate active status change to service
	if err = i.detailService.SetExceptionActive(ctx, tx, exceptionID, active); err != nil {
		return err
	}

	// 5. Commit
	if err = tx.Commit(); err != nil {
		scheduleExceptionInteractorLog.Error(logger.LogScheduleDetailInteractorCommitError, "error", err)
		return err
	}

	err = nil
	return nil
}
