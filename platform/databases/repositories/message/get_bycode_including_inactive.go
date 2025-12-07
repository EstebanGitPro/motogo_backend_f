package message

import (
	"context"
	"database/sql"

	cachetypes "github.com/EstebanGitPro/motogo-backend/platform/cache/types"
)

// GetByCodeIncludingInactive returns a message by code without filtering by is_active
// This is used to detect if a message exists but is inactive
func (r *repository) GetByCodeIncludingInactive(ctx context.Context, code string) (*cachetypes.CachedMessage, error) {
	var m cachetypes.CachedMessage
	var createdAt, updatedAt interface{}

	err := r.db.QueryRowContext(ctx, queryGetByCodeIncludingInactive, code).Scan(
		&m.ID,
		&m.Code,
		&m.Type,
		&m.Category,
		&m.Module,
		&m.Title,
		&m.Content,
		&m.Active,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}
