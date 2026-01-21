package franchise

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// GetFranchisesByRepresentative lists all franchises owned by a representative
func (r *repository) GetFranchisesByRepresentative(ctx context.Context, representativeID string) ([]domain.Franchise, error) {
	rows, err := r.stmtGetFranchisesByRepresentative.QueryContext(ctx, representativeID)
	if err != nil {
		log.Error(logger.LogFranchiseRepoGetByRepError, "representative_id", representativeID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var franchises []domain.Franchise
	for rows.Next() {
		var franchise domain.Franchise
		var description sql.NullString

		if err := rows.Scan(&franchise.ID, &franchise.Name, &description); err != nil {
			log.Error(logger.LogFranchiseRepoScanError, "error", err)
			return nil, err
		}

		if description.Valid {
			franchise.Description = &description.String
		}

		franchises = append(franchises, franchise)
	}

	return franchises, nil
}
