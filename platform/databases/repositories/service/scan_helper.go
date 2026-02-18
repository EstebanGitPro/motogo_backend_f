package service

import (
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// scanServices iterates over rows and scans each into a domain.Service.
// Handles the nullable description column. Centralizes shared scan logic.
func scanServices(rows *sql.Rows) ([]domain.Service, error) {
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
