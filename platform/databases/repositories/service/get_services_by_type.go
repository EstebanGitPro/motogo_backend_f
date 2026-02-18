package service

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) GetServicesByType(ctx context.Context, serviceType string) ([]domain.Service, error) {
	log.Info(logger.LogServiceRepoGetByType, "type", serviceType)

	rows, err := r.stmtGetServicesByType.QueryContext(ctx, serviceType)
	if err != nil {
		log.Error(logger.LogServiceRepoGetByTypeError, "error", err)
		return nil, err
	}
	defer func() { _ = rows.Close() }() // Rows close error intentionally ignored

	return scanServices(rows)
}
