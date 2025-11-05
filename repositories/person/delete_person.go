package person

import (
	"context"

	"fmt"

	domain "github.com/EstebanGitPro/motogo-backend/core/domain"
)

func (r *repository) DeletePerson(ctx context.Context, id string) error {

	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, queryDelete, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete person: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get rows affected: %w", err)
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
