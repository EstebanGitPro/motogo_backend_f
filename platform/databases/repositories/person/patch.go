package person

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

const (
	queryPatch = "UPDATE persons SET keycloak_user_id = ? WHERE id = ?"
)

func (r *repository) PatchPerson(id string, keycloakUserID string) error {
	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(context.Background(), queryPatch, keycloakUserID, id)
	if err != nil {
		tx.Rollback()
		return domain.ErrUserCannotSave
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return domain.ErrPersonNotFound
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
