package message

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	cachetypes "github.com/EstebanGitPro/motogo-backend/platform/cache/types"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
)

const (
	queryGetAllActive = `
		SELECT id, message_code, type, category, module, message_title, message_content, is_active, created_at, updated_at 
		FROM system_messages 
		WHERE is_active = true`

	queryGetByCode = `
		SELECT id, message_code, type, category, module, message_title, message_content, is_active, created_at, updated_at 
		FROM system_messages 
		WHERE message_code = ?
		LIMIT 1`

	queryGetByCodeForCache = `
		SELECT id, message_code, type, category, module, message_title, message_content, is_active, created_at, updated_at 
		FROM system_messages 
		WHERE message_code = ? AND is_active = true
		LIMIT 1`

	queryGetByID = `
		SELECT id, message_code, type, category, module, message_title, message_content, is_active, created_at, updated_at 
		FROM system_messages 
		WHERE id = ?
		LIMIT 1`

	queryGetByType = `
		SELECT id, message_code, type, category, module, message_title, message_content, is_active, created_at, updated_at 
		FROM system_messages 
		WHERE type = ? AND is_active = true`

	queryGetByModule = `
		SELECT id, message_code, type, category, module, message_title, message_content, is_active, created_at, updated_at 
		FROM system_messages 
		WHERE module = ? AND is_active = true`

	queryMessageSave = `INSERT INTO system_messages 
		(id, message_code, type, category, module, message_title, message_content, is_active) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	queryMessageDelete = `DELETE FROM system_messages WHERE id = ?`
)

// repository implements output.MessageRepository
type repository struct {
	db *sql.DB
}

type MessageRepository interface {
	output.MessageRepository
	cachetypes.MessageCacheRepository
}

// NewMessageRepository creates a new message repository
func NewMessageRepository(db *sql.DB) (MessageRepository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}
	return &repository{db: db}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return common.NewSQLTx(tx), nil
}
