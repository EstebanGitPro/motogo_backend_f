package franchise

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// GetFranchiseByID retrieves a franchise by its ID
func (r *repository) GetFranchiseByID(ctx context.Context, franchiseID string) (*domain.Franchise, error) {
	var franchise domain.Franchise
	var description sql.NullString

	err := r.stmtGetFranchiseByID.QueryRowContext(ctx, franchiseID).Scan(
		&franchise.ID,
		&franchise.Name,
		&description,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrFranchiseNotFound
		}
		log.Error(logger.LogFranchiseRepoGetByIDError, "franchise_id", franchiseID, "error", err)
		return nil, err
	}

	if description.Valid {
		franchise.Description = &description.String
	}

	return &franchise, nil
}