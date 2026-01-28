package location

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)



func (r *repository) CheckAddressExists(ctx context.Context, address string) (bool, error) {
	var exists int
	err := r.db.QueryRowContext(ctx, queryCheckAddressExists, address).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Error(logger.LogLocationRepoGetDepartmentsError, "error", err, "address", address)
		return false, err
	}
	return true, nil
}
