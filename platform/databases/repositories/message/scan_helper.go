package message

import (
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// scanMessages iterates over rows and scans each into a domain.Message.
// Centralizes the shared scan logic for all message queries.
func scanMessages(rows *sql.Rows) ([]domain.Message, error) {
	var messages []domain.Message
	for rows.Next() {
		var m domain.Message
		if err := rows.Scan(
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
		); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, rows.Err()
}
