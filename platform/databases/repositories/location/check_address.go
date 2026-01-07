package location

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

const queryCheckAddressExists = `
	SELECT 1 FROM locations WHERE LOWER(TRIM(address)) = LOWER(TRIM(?)) LIMIT 1
`

// CheckAddressExists checks if an address already exists in the locations table
func (r *repository) CheckAddressExists(ctx context.Context, address string) (bool, error) {
	var exists int
	err := r.db.QueryRowContext(ctx, queryCheckAddressExists, address).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // Address doesn't exist
		}
		log.Error(logger.LogLocationRepoGetDepartmentsError, "error", err, "address", address)
		return false, err
	}
	return true, nil // Address exists
}
