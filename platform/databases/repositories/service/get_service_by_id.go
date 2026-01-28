package service

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) GetServiceByID(ctx context.Context, serviceID string) (*domain.Service, error) {
	log.Info(logger.LogServiceRepoGetByID, "service_id", serviceID)

	var svc domain.Service
	var description sql.NullString
	err := r.stmtGetServiceByID.QueryRowContext(ctx, serviceID).Scan(
		&svc.ID,
		&svc.Name,
		&description,
		&svc.ServiceType,
		&svc.IsActive,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn(logger.LogServiceRepoNotFound, "service_id", serviceID)
			return nil, domain.ErrServiceNotFound
		}
		log.Error(logger.LogServiceRepoGetByIDError, "error", err, "service_id", serviceID)
		return nil, err
	}

	if description.Valid {
		svc.Description = description.String
	}

	log.Success(logger.LogServiceRepoGetByIDOK, "service_id", serviceID)
	return &svc, nil
}
