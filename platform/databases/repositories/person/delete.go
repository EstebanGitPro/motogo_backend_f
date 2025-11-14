package person

import (
    "context"

    "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
    "github.com/EstebanGitPro/motogo-backend/core/ports/output"
)

func (r *repository) DeletePerson(ctx context.Context, tx output.Tx, id string) error {
    dbTx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }

    _, err = dbTx.ExecContext(ctx, queryDelete, id)
    if err != nil {
        dbTx.Rollback()
        return domain.ErrUserCannotSave
    }

    err = dbTx.Commit()
    if err != nil {
        return err
    }

    return nil
}
