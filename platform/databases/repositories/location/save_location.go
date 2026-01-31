package location

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	uuid "github.com/EstebanGitPro/motogo-backend/tools/utils"
)

func (r *repository) SaveLocation(ctx context.Context, tx output.Tx, location domain.Location) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	locationID := uuid.Generate()

	_, err := sqlTx.ExecContext(ctx, querySaveLocation,
		locationID,
		location.BranchID,
		location.CityID,
		location.Address,
		location.Latitude,
		location.Longitude,
	)

	if err != nil {
		log.Error(logger.LogLocationRepoSaveError, "error", err, "branch_id", location.BranchID)
		return domain.ErrLocationCannotSave
	}

	return nil
}
