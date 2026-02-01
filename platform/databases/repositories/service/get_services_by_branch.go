package service

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) GetServicesByBranch(ctx context.Context, branchID string) ([]domain.BranchServiceInfo, error) {
	log.Info(logger.LogBranchServicesRepoGetByBranch, "branch_id", branchID)

	rows, err := r.stmtGetServicesByBranch.QueryContext(ctx, branchID)
	if err != nil {
		log.Error(logger.LogBranchServicesRepoGetByBranchErr, "error", err)
		return nil, err
	}
	defer func() { _ = rows.Close() }() // Rows close error intentionally ignored

	var services []domain.BranchServiceInfo
	for rows.Next() {
		var info domain.BranchServiceInfo
		var description sql.NullString
		var createdAt sql.NullTime
		if err := rows.Scan(
			&info.Service.ID,
			&info.Service.Name,
			&description,
			&info.Service.ServiceType,
			&createdAt,
			&info.Active,
		); err != nil {
			log.Error(logger.LogServiceRepoScanError, "error scanning branch service", err)
			continue
		}
		if description.Valid {
			info.Service.Description = description.String
		}
		if createdAt.Valid {
			info.AddedAt = createdAt.Time.Format("2006-01-02T15:04:05Z07:00")
		}
		services = append(services, info)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}
