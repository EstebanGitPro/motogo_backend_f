package message

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
)

func (r *repository) SaveMessage(ctx context.Context, tx output.Tx, message domain.Message) error {
	messageToUpdate := FromDomain(message)

	dbTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	_, err := dbTx.ExecContext(ctx, queryMessageSave,
		messageToUpdate.ID,
		messageToUpdate.Code,
		messageToUpdate.Type,
		messageToUpdate.Category,
		messageToUpdate.Module,
		messageToUpdate.Title,
		messageToUpdate.Content,
		messageToUpdate.Active,
	)

	if err != nil {
		return domain.ErrMessageCannotSave
	}

	return nil
}
