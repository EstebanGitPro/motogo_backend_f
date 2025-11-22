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

	dbTx, ok := tx.(*sqlTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	
	result, err := dbTx.ExecContext(ctx, queryPatch, keycloakUserID, id)
	if err != nil {
		return domain.ErrUserCannotSave
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrPersonNotFound
	}

	return nil
}
