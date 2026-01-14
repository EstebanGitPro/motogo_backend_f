package branch

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// HasBranchesByRepresentative checks if a representative has any branches (HU53)
func (r *repository) HasBranchesByRepresentative(ctx context.Context, representativeID string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM branches WHERE representative_id = ? LIMIT 1)"

	var exists bool
	err := r.db.QueryRowContext(ctx, query, representativeID).Scan(&exists)
	if err != nil {
		log.Error(logger.LogBranchRepoGetByRepError, "error", err, "representative_id", representativeID)
		return false, err
	}

	return exists, nil
}
