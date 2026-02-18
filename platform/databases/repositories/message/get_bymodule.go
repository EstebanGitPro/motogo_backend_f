package message

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

func (r *repository) GetByModule(ctx context.Context, module string) ([]domain.Message, error) {
	rows, err := r.db.QueryContext(ctx, queryGetByModule, module)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }() // Rows close error intentionally ignored

	return scanMessages(rows)
}
