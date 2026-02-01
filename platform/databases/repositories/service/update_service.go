package service

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) UpdateService(ctx context.Context, tx output.Tx, service domain.Service) error {
	log.Info(logger.LogServiceRepoUpdate, "service_id", service.ID)

	result, err := r.stmtUpdateService.ExecContext(ctx,
		service.Name,
		service.Description,
		string(service.ServiceType),
		service.IsActive,
		service.ID,
	)
	if err != nil {
		log.Error(logger.LogServiceRepoUpdateError, "error", err, "service_id", service.ID)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Success(logger.LogServiceRepoUpdateOK, "service_id", service.ID, "rows_affected", rowsAffected)
	return nil
}
