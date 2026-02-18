package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// deferRollback returns a function suitable for use with defer that rolls back
// the transaction if *errPtr is non-nil. This centralizes the repeated
// rollback-with-logging pattern used across schedule interactors.
func deferRollback(tx output.Tx, errPtr *error, log logger.Logger, rollbackErrMsg, rollbackOKMsg string) func() {
	return func() {
		if *errPtr != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(rollbackErrMsg,
					"rollback_error", rbErr,
					"original_error", *errPtr)
			} else {
				log.Warn(rollbackOKMsg)
			}
		}
	}
}

// verifyScheduleOwnership checks that the representativeID is the owner of the
// branch associated with the given scheduleID. Returns nil on success, or a
// domain error (ErrBranchNotFound, ErrForbidden) on failure.
func verifyScheduleOwnership(
	ctx context.Context,
	scheduleService input.ScheduleService,
	branchService input.BranchService,
	scheduleID, representativeID string,
) error {
	schedule, schedErr := scheduleService.GetScheduleByID(ctx, scheduleID)
	if schedErr != nil {
		return schedErr
	}

	branch, branchErr := branchService.GetBranchByID(ctx, schedule.BranchID)
	if branchErr != nil {
		return domain.ErrBranchNotFound
	}
	if branch.RepresentativeID != representativeID {
		return domain.ErrForbidden
	}

	return nil
}

// verifyBranchOwnership checks that the representativeID is the owner of the
// branch identified by branchID. Returns the branch on success, or a domain
// error on failure.
func verifyBranchOwnership(
	ctx context.Context,
	branchService input.BranchService,
	branchID, representativeID string,
	log logger.Logger,
	branchErrMsg, ownershipErrMsg string,
) error {
	branch, branchErr := branchService.GetBranchByID(ctx, branchID)
	if branchErr != nil {
		log.Error(branchErrMsg, "error", branchErr)
		return domain.ErrBranchNotFound
	}
	if branch.RepresentativeID != representativeID {
		log.Warn(ownershipErrMsg,
			"branch_id", branchID, "representative_id", representativeID)
		return domain.ErrForbidden
	}
	return nil
}
