package person

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

func (r *repository) DeletePerson(id string) error {
	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(context.Background(), queryDelete, id)
	if err != nil {
		tx.Rollback()
		return domain.ErrUserCannotSave
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
