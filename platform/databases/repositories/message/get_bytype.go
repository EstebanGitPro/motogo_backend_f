package message

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

func (r *repository) GetByType(ctx context.Context, msgType string) ([]domain.Message, error) {
	rows, err := r.db.QueryContext(ctx, queryGetByType, msgType)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }() // Rows close error intentionally ignored

	return scanMessages(rows)
}
