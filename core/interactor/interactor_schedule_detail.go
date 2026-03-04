package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

var scheduleDetailInteractorLog = logger.NewSlogLogger()

// ScheduleDetailInteractor orchestrates schedule detail operations (HU6-9)
type ScheduleDetailInteractor struct {
	detailService   input.ScheduleDetailService
	scheduleService input.ScheduleService
	branchService   input.BranchService
}

// NewScheduleDetailInteractor creates a new ScheduleDetailInteractor
func NewScheduleDetailInteractor(
	detailService input.ScheduleDetailService,
	scheduleService input.ScheduleService,
	branchService input.BranchService,
) *ScheduleDetailInteractor {
	return &ScheduleDetailInteractor{
		detailService:   detailService,
		scheduleService: scheduleService,
		branchService:   branchService,
	}
}

// CreateDetail orchestrates schedule detail creation (HU6)
func (i *ScheduleDetailInteractor) CreateDetail(
	ctx context.Context,
	detail domain.ScheduleDetail,
	representativeID, branchID string,
) (result *domain.ScheduleDetail, err error) {
	scheduleDetailInteractorLog.Info(logger.LogScheduleDetailInteractorCreateStart,
		"schedule_id", detail.ScheduleID,
		"day_of_week", detail.DayOfWeek,
		"representative_id", representativeID)

	// 1. Verify ownership of branch
	if ownerErr := verifyBranchOwnership(ctx, i.branchService, branchID, representativeID,
		scheduleDetailInteractorLog, logger.LogScheduleDetailInteractorBranchError, logger.LogScheduleDetailInteractorOwnershipError); ownerErr != nil {
		return nil, ownerErr
	}

	// 2. Begin transaction
	tx, txErr := i.detailService.BeginTx(ctx)
	if txErr != nil {
		scheduleDetailInteractorLog.Error(logger.LogScheduleDetailInteractorTxError, "error", txErr)
		return nil, txErr
	}

	defer func() { _ = tx.Rollback() }()

	// 3. Create detail via service
	result, err = i.detailService.CreateDetail(ctx, tx, detail)
	if err != nil {
		scheduleDetailInteractorLog.Error(logger.LogScheduleDetailInteractorCreateError,
			"schedule_id", detail.ScheduleID, "error", err)
		return nil, err
	}

	// 4. Commit transaction
	if err = tx.Commit(); err != nil {
		scheduleDetailInteractorLog.Error(logger.LogScheduleDetailInteractorCommitError, "error", err)
		return nil, err
	}

	scheduleDetailInteractorLog.Info(logger.LogScheduleDetailInteractorCreateOK,
		"detail_id", result.ID,
		"schedule_id", detail.ScheduleID)

	err = nil
	return result, nil
}

// ListDetails retrieves all schedule details for a schedule (HU9)
func (i *ScheduleDetailInteractor) ListDetails(ctx context.Context, scheduleID string) ([]domain.ScheduleDetail, error) {
	details, err := i.detailService.GetDetailsByScheduleID(ctx, scheduleID)
	if err != nil {
		scheduleDetailInteractorLog.Error(logger.LogScheduleDetailInteractorListError,
			"schedule_id", scheduleID, "error", err)
		return nil, err
	}

	scheduleDetailInteractorLog.Info(logger.LogScheduleDetailInteractorListOK,
		"schedule_id", scheduleID, "count", len(details))

	return details, nil
}

// GetDetail retrieves a schedule detail by ID
func (i *ScheduleDetailInteractor) GetDetail(ctx context.Context, detailID string) (*domain.ScheduleDetail, error) {
	detail, err := i.detailService.GetDetailByID(ctx, detailID)
	if err != nil {
		return nil, err
	}
	return detail, nil
}

// UpdateDetail orchestrates schedule detail update (HU7)
func (i *ScheduleDetailInteractor) UpdateDetail(
	ctx context.Context,
	detail domain.ScheduleDetail,
	representativeID string,
) (err error) {
	// 1. Get existing detail
	existing, existErr := i.detailService.GetDetailByID(ctx, detail.ID)
	if existErr != nil {
		return existErr
	}
	if existing == nil {
		return domain.ErrScheduleDetailNotFound
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

	defer func() { _ = tx.Rollback() }()

	// 4. Update detail
	if err = i.detailService.UpdateDetail(ctx, tx, detail); err != nil {
		return err
	}

	// 5. Commit
	if err = tx.Commit(); err != nil {
		scheduleDetailInteractorLog.Error(logger.LogScheduleDetailInteractorCommitError, "error", err)
		return err
	}

	err = nil
	return nil
}

// DeleteDetail orchestrates schedule detail deletion (HU8)
func (i *ScheduleDetailInteractor) DeleteDetail(
	ctx context.Context,
	detailID, representativeID string,
) (err error) {
	// 1. Get existing detail
	existing, existErr := i.detailService.GetDetailByID(ctx, detailID)
	if existErr != nil {
		return existErr
	}
	if existing == nil {
		return domain.ErrScheduleDetailNotFound
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

	defer func() { _ = tx.Rollback() }()

	// 4. Delete detail
	if err = i.detailService.DeleteDetail(ctx, tx, detailID); err != nil {
		return err
	}

	// 5. Commit
	if err = tx.Commit(); err != nil {
		scheduleDetailInteractorLog.Error(logger.LogScheduleDetailInteractorCommitError, "error", err)
		return err
	}

	err = nil
	return nil
}
