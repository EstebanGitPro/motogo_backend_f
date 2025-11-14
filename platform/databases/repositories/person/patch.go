package person

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
)

const (
	queryPatch = "UPDATE persons SET keycloak_user_id = ? WHERE id = ?"
)

func (r *repository) PatchPerson(ctx context.Context, tx output.Tx, id string, keycloakUserID string) error {
	var dbTx *sqlTx
	var shouldCommit bool

	if tx != nil {
		// Usar la transacción existente
		dbTx = tx.(*sqlTx)
		shouldCommit = false
	} else {
		// Crear nueva transacción
		newTx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		dbTx = &sqlTx{Tx: newTx}
		shouldCommit = true
	}

	result, err := dbTx.ExecContext(ctx, queryPatch, keycloakUserID, id)
	if err != nil {
		if shouldCommit {
			dbTx.Rollback()
		}
		return domain.ErrUserCannotSave
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		if shouldCommit {
			dbTx.Rollback()
		}
		return err
	}

	if rowsAffected == 0 {
		if shouldCommit {
			dbTx.Rollback()
		}
		return domain.ErrPersonNotFound
	}

	if shouldCommit {
		err = dbTx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}
