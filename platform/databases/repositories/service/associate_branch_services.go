package service

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) AssociateBranchServices(ctx context.Context, tx output.Tx, branchID string, serviceIDs []string) error {
	log.Info(logger.LogBranchServicesRepoAssociate, "branch_id", branchID, "service_count", len(serviceIDs))

	for _, serviceID := range serviceIDs {
		id := domain.GenerateUUID()
		_, err := r.stmtInsertBranchService.ExecContext(ctx, id, branchID, serviceID)
		if err != nil {
			log.Error(logger.LogBranchServicesRepoAssociateErr, "branch_id", branchID, "service_id", serviceID, "error", err)
			return err
		}
	}

	log.Success(logger.LogBranchServicesRepoAssociateOK, "branch_id", branchID, "service_count", len(serviceIDs))
	return nil
}
