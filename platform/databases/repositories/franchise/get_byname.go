package franchise

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) GetFranchiseByName(ctx context.Context, name string) (*domain.Franchise, error) {
	var franchise domain.Franchise
	var description sql.NullString

	err := r.stmtGetFranchiseByName.QueryRowContext(ctx, name).Scan(
		&franchise.ID,
		&franchise.Name,
		&description,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found is not an error for duplicate check
		}
		log.Error(logger.LogFranchiseRepoGetByNameError, "name", name, "error", err)
		return nil, err
	}

	if description.Valid {
		franchise.Description = &description.String
	}

	return &franchise, nil
}
