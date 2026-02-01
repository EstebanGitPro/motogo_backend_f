package service

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) DissociateBranchService(ctx context.Context, tx output.Tx, branchID, serviceID string) error {
	log.Info(logger.LogBranchServicesRepoDissociate, "branch_id", branchID, "service_id", serviceID)

	result, err := r.stmtDeleteBranchService.ExecContext(ctx, branchID, serviceID)
	if err != nil {
		log.Error(logger.LogBranchServicesRepoDissociateErr, "error", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Warn(logger.LogBranchServicesRepoNotFound, "branch_id", branchID, "service_id", serviceID)
		return domain.ErrServiceNotFound
	}

	log.Success(logger.LogBranchServicesRepoDissociateOK, "branch_id", branchID, "service_id", serviceID)
	return nil
}
