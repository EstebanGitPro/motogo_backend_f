package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

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
