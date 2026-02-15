package branch

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

func (r *repository) getBranchDisplacementRanges(ctx context.Context, branchID string) ([]domain.DisplacementRange, error) {
	rows, err := r.stmtGetBranchDisplacementRanges.QueryContext(ctx, branchID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }() // Rows close error intentionally ignored

	var ranges []domain.DisplacementRange
	for rows.Next() {
		var dr domain.DisplacementRange
		if rows.Scan(&dr) != nil {
			continue
		}
		ranges = append(ranges, dr)
	}

	return ranges, nil
}
