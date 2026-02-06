package motorcycle

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// Update updates an existing motorcycle in the database (HU44)
func (r *repository) Update(ctx context.Context, tx output.Tx, motorcycle *domain.Motorcycle) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		log.Error(logger.LogMotorcycleRepoInvalidTx, "expected SQLTx")
		return domain.ErrInvalidTransaction
	}

	// Handle nullable reference_id
	var refID interface{}
	if motorcycle.ReferenceID != "" {
		refID = motorcycle.ReferenceID
	} else {
		refID = nil
	}

	_, err := sqlTx.ExecContext(ctx, queryUpdate,
		refID,
		motorcycle.Year,
		motorcycle.CurrentMileage,
		motorcycle.OwnerNotes,
		motorcycle.ProfileImageURL,
		motorcycle.ID,
	)

	if err != nil {
		log.Error(logger.LogMotorcycleRepoUpdateError, "error", err, "motorcycle_id", motorcycle.ID)
		return err
	}

	return nil
}
