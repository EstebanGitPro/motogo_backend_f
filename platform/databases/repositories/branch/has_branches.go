package branch

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) HasBranchesByRepresentative(ctx context.Context, representativeID string) (bool, error) {

	var exists bool
	err := r.db.QueryRowContext(ctx, queryHasBranchesByRepresentative, representativeID).Scan(&exists)
	if err != nil {
		log.Error(logger.LogBranchRepoGetByRepError, "error", err, "representative_id", representativeID)
		return false, err
	}

	return exists, nil
}
