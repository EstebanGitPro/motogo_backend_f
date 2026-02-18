package service

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) GetAllServices(ctx context.Context) ([]domain.Service, error) {
	log.Info(logger.LogServiceRepoGetAll)

	rows, err := r.stmtGetAllServices.QueryContext(ctx)
	if err != nil {
		log.Error(logger.LogServiceRepoGetAllError, "error", err)
		return nil, err
	}
	defer func() { _ = rows.Close() }() // Rows close error intentionally ignored

	return scanServices(rows)
}
