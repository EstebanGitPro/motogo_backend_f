package service

import (
	"context"
	"database/sql"

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
	defer rows.Close()

	var services []domain.Service
	for rows.Next() {
		var svc domain.Service
		var description sql.NullString
		if err := rows.Scan(&svc.ID, &svc.Name, &description, &svc.ServiceType, &svc.IsActive); err != nil {
			log.Error(logger.LogServiceRepoScanError, "error scanning service", err)
			continue
		}
		if description.Valid {
			svc.Description = description.String
		}
		services = append(services, svc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}
