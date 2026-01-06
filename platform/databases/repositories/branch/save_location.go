package branch

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	uuid "github.com/EstebanGitPro/motogo-backend/tools/utils"
)

// SaveLocation saves a location for a branch
func (r *repository) SaveLocation(ctx context.Context, tx output.Tx, location domain.Location) error {
	sqlTx := tx.(*common.SQLTx)

	locationID := uuid.Generate()

	_, err := sqlTx.StmtContext(ctx, r.stmtSaveLocation).ExecContext(ctx,
		locationID,
		location.BranchID,
		location.CityID,
		location.Address,
		location.Latitude,
		location.Longitude,
	)

	if err != nil {
		log.Error("error saving location", "error", err, "branch_id", location.BranchID)
		return domain.ErrLocationCannotSave
	}

	return nil
}