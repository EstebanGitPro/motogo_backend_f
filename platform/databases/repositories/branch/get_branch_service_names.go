package branch

import (
	"context"
)

// getBranchServiceNames retrieves the service names associated with a branch
func (r *repository) getBranchServiceNames(ctx context.Context, branchID string) ([]string, error) {
	rows, err := r.stmtGetBranchServiceNames.QueryContext(ctx, branchID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var names []string
	for rows.Next() {
		var name string
		if rows.Scan(&name) != nil {
			continue
		}
		names = append(names, name)
	}

	return names, nil
}
