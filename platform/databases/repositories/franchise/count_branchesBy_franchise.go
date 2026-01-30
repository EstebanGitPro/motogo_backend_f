package franchise

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) CountBranchesByFranchise(ctx context.Context, franchiseID string) (int, error) {
	var count int
	err := r.stmtCountBranchesByFranchise.QueryRowContext(ctx, franchiseID).Scan(&count)
	if err != nil {
		log.Error(logger.LogFranchiseRepoCountBranches, "franchise_id", franchiseID, "error", err)
		return 0, err
	}
	return count, nil
}
