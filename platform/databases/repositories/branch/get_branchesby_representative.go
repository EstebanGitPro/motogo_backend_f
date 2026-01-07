package branch

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// GetBranchesByRepresentative retrieves all branches for a representative
func (r *repository) GetBranchesByRepresentative(ctx context.Context, representativeID string) ([]domain.Branch, error) {
	rows, err := r.stmtGetBranchesByRepresentative.QueryContext(ctx, representativeID)
	if err != nil {
		log.Error(logger.LogBranchRepoGetByRepError, "error", err, "representative_id", representativeID)
		return nil, err
	}
	defer rows.Close()

	var branches []domain.Branch
	for rows.Next() {
		var branch domain.Branch
		var franchiseID, profileImageURL sql.NullString

		err := rows.Scan(
			&branch.ID,
			&branch.RepresentativeID,
			&franchiseID,
			&branch.Name,
			&branch.EstablishmentType,
			&profileImageURL,
			&branch.Status,
		)
		if err != nil {
			log.Error(logger.LogBranchRepoScanError, "error", err)
			continue
		}

		if franchiseID.Valid {
			branch.FranchiseID = &franchiseID.String
		}
		if profileImageURL.Valid {
			branch.ProfileImageURL = &profileImageURL.String
		}

		branches = append(branches, branch)
	}

	return branches, nil
}
