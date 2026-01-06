package branch

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// GetBranchByFranchiseAndName retrieves a branch by franchise and name (for duplicate validation)
func (r *repository) GetBranchByFranchiseAndName(ctx context.Context, franchiseID, name string) (*domain.Branch, error) {
	var branch domain.Branch
	var franchiseIDVal, profileImageURL sql.NullString

	err := r.stmtGetBranchByFranchiseAndName.QueryRowContext(ctx, franchiseID, name).Scan(
		&branch.ID,
		&branch.RepresentativeID,
		&franchiseIDVal,
		&branch.Name,
		&branch.EstablishmentType,
		&profileImageURL,
		&branch.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrBranchNotFound
		}
		log.Error("error getting branch by franchise and name", "error", err, "franchise_id", franchiseID, "name", name)
		return nil, err
	}

	if franchiseIDVal.Valid {
		branch.FranchiseID = &franchiseIDVal.String
	}
	if profileImageURL.Valid {
		branch.ProfileImageURL = &profileImageURL.String
	}

	return &branch, nil
}