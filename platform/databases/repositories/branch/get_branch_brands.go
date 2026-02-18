package branch

import (
	"context"
)

func (r *repository) getBranchBrands(ctx context.Context, branchID string) ([]string, error) {
	rows, err := r.stmtGetBranchBrands.QueryContext(ctx, branchID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }() // Rows close error intentionally ignored

	var brands []string
	for rows.Next() {
		var brand string
		if rows.Scan(&brand) != nil {
			continue
		}
		brands = append(brands, brand)
	}

	return brands, nil
}
