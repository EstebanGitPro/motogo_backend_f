package branch

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
)

// UpdateLocation updates a location for a branch
func (r *repository) UpdateLocation(ctx context.Context, tx output.Tx, location domain.Location) error {
	sqlTx := tx.(*common.SQLTx)

	_, err := sqlTx.StmtContext(ctx, r.stmtUpdateLocation).ExecContext(ctx,
		location.CityID,
		location.Address,
		location.Latitude,
		location.Longitude,
		location.BranchID,
	)

	if err != nil {
		log.Error("error updating location", "error", err, "branch_id", location.BranchID)
		return domain.ErrLocationCannotSave
	}

	return nil
}