package person

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
)

func (r *repository) DeletePerson(ctx context.Context, tx output.Tx, id string) error {
	// Type assertion segura
	dbTx, ok := tx.(*sqlTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	// Solo ejecutar el query - NO manejar commit/rollback
	_, err := dbTx.ExecContext(ctx, queryDelete, id)
	if err != nil {
		return domain.ErrUserCannotSave
	}

	return nil
}
