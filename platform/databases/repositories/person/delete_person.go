package person

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
)

func (r *repository) DeletePerson(ctx context.Context, tx output.Tx, id string) error {
	var err error

	if tx != nil {
		// Use transaction if provided
		dbTx, ok := tx.(*common.SQLTx)
		if !ok {
			return domain.ErrInvalidTransaction
		}
		_, err = dbTx.ExecContext(ctx, queryDelete, id)
	} else {
		// Use direct connection if no transaction (HU53 self-delete flow)
		_, err = r.db.ExecContext(ctx, queryDelete, id)
	}

	if err != nil {
		return domain.ErrUserCannotDelete
	}

	return nil
}
