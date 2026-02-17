package completed_service

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// GetByID retrieves a completed service by its ID
func (r *repository) GetByID(ctx context.Context, serviceID string) (*domain.CompletedService, error) {
	var cs CompletedService

	err := r.stmtGetByID.QueryRowContext(ctx, serviceID).Scan(
		&cs.ID, &cs.BranchID, &cs.MotorcycleID, &cs.DiagnosticID,
		&cs.RequestDate, &cs.CompletionDate, &cs.Status,
		&cs.QuotedPrice, &cs.FinalPrice, &cs.RepresentativeNotes,
		&cs.CreatedAt, &cs.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrCompletedServiceNotFound
		}
		log.Error(logger.LogCSRepoGetError, err)
		return nil, err
	}

	result := cs.ToDomain()
	return &result, nil
}

// GetByMotorcycleID retrieves completed services for a motorcycle
func (r *repository) GetByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.CompletedService, error) {
	rows, err := r.stmtGetByMotorcycleID.QueryContext(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogCSRepoGetByMotoErr, err)
		return nil, err
	}
	defer rows.Close()

	return r.scanMultiple(rows)
}

// GetByBranchID retrieves completed services for a branch
func (r *repository) GetByBranchID(ctx context.Context, branchID string) ([]domain.CompletedService, error) {
	rows, err := r.stmtGetByBranchID.QueryContext(ctx, branchID)
	if err != nil {
		log.Error(logger.LogCSRepoGetByBranchErr, err)
		return nil, err
	}
	defer rows.Close()

	return r.scanMultiple(rows)
}

// GetItemsByCompletedServiceID retrieves items for a completed service
func (r *repository) GetItemsByCompletedServiceID(ctx context.Context, completedServiceID string) ([]domain.CompletedServiceItem, error) {
	rows, err := r.stmtGetItemsByCSID.QueryContext(ctx, completedServiceID)
	if err != nil {
		log.Error(logger.LogCSRepoGetItemsErr, err)
		return nil, err
	}
	defer rows.Close()

	var items []domain.CompletedServiceItem
	for rows.Next() {
		var item CompletedServiceItem
		if err := rows.Scan(&item.ID, &item.CompletedServiceID, &item.ServiceID,
			&item.Rating, &item.Comment, &item.RatedAt, &item.IsOffensiveComment); err != nil {
			log.Error(logger.LogCSRepoScanItemErr, err)
			return nil, err
		}
		items = append(items, item.ItemToDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// scanMultiple scans multiple rows into domain entities
func (r *repository) scanMultiple(rows *sql.Rows) ([]domain.CompletedService, error) {
	var services []domain.CompletedService
	for rows.Next() {
		var cs CompletedService
		if err := rows.Scan(
			&cs.ID, &cs.BranchID, &cs.BranchName,
			&cs.MotorcycleID, &cs.DiagnosticID,
			&cs.RequestDate, &cs.CompletionDate, &cs.Status,
			&cs.QuotedPrice, &cs.FinalPrice, &cs.RepresentativeNotes,
			&cs.CreatedAt, &cs.UpdatedAt,
		); err != nil {
			log.Error(logger.LogCSRepoScanError, err)
			return nil, err
		}
		services = append(services, cs.ToDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

// HasActiveService checks if a motorcycle already has an active service (PENDIENTE or EN_PROCESO) at the given branch
func (r *repository) HasActiveService(ctx context.Context, branchID, motorcycleID string) (bool, error) {
	var count int
	err := r.stmtHasActiveService.QueryRowContext(ctx, branchID, motorcycleID).Scan(&count)
	if err != nil {
		log.Error(logger.LogCSRepoHasActiveErr, err)
		return false, err
	}
	return count > 0, nil
}
