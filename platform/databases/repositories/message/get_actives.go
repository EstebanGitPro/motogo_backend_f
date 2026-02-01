package message

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

func (r *repository) GetAllActive(ctx context.Context) ([]domain.Message, error) {
	rows, err := r.db.QueryContext(ctx, queryGetAllActive)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }() // Rows close error intentionally ignored

	var messages []domain.Message
	for rows.Next() {
		var m domain.Message
		err := rows.Scan(
			&m.ID,
			&m.Code,
			&m.Type,
			&m.Category,
			&m.Module,
			&m.Title,
			&m.Content,
			&m.Active,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, rows.Err()
}
