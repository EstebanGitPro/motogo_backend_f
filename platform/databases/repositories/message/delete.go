package message

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
)

func (r *repository) DeleteMessage(ctx context.Context, tx output.Tx, id string) error {
	dbTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	_, err := dbTx.ExecContext(ctx, queryMessageDelete, id)
	if err != nil {
		return domain.ErrMessageCannotDelete
	}

	return nil
}
