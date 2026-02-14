package branch

import (
	"context"
)

func (r *repository) getBranchDisplacementRanges(ctx context.Context, branchID string) ([]string, error) {
	rows, err := r.stmtGetBranchDisplacementRanges.QueryContext(ctx, branchID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }() // Rows close error intentionally ignored

	var ranges []string
	for rows.Next() {
		var dr string
		if err := rows.Scan(&dr); err != nil {
			continue
		}
		ranges = append(ranges, dr)
	}

	return ranges, nil
}
