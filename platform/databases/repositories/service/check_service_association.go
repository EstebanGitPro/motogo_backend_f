package service

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) CheckServiceAssociation(ctx context.Context, branchID, serviceID string) (bool, error) {
	log.Info(logger.LogBranchServicesRepoCheckAssoc, "branch_id", branchID, "service_id", serviceID)

	var count int
	err := r.stmtCheckServiceAssociation.QueryRowContext(ctx, branchID, serviceID).Scan(&count)
	if err != nil {
		log.Error(logger.LogBranchServicesRepoCheckAssocErr, "error", err)
		return false, err
	}

	return count > 0, nil
}
