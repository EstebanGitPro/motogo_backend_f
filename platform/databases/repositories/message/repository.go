package message

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
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

	queryGetByType = `
		SELECT id, message_code, type, category, module, message_title, message_content, is_active, created_at, updated_at 
		FROM system_messages 
		WHERE type = ? AND is_active = true`

	queryGetByModule = `
		SELECT id, message_code, type, category, module, message_title, message_content, is_active, created_at, updated_at 
		FROM system_messages 
		WHERE module = ? AND is_active = true`

	queryMessageSave = `INSERT INTO system_messages 
		(message_code, type, category, module, message_title, message_content, is_active) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	queryMessageUpdate = `UPDATE system_messages 
		SET message_code = ?, type = ?, category = ?, module = ?, 
		    message_title = ?, message_content = ?, is_active = ? 
		WHERE id = ?`

	queryMessageDelete = `DELETE FROM system_messages WHERE id = ?`
)

// repository implements output.MessageRepository
type repository struct {
	db *sql.DB
}

// MessageRepository combines both interfaces for the repository
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

// ============================================
// Write operations (transactional)
// ============================================

func (r *repository) SaveMessage(ctx context.Context, tx output.Tx, message domain.Message) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		_, err := r.db.ExecContext(ctx, queryMessageSave,
			message.Code, message.Type, message.Category, message.Module,
			message.Title, message.Content, message.Active,
		)
		return err
	}

	_, err := sqlTx.ExecContext(ctx, queryMessageSave,
		message.Code, message.Type, message.Category, message.Module,
		message.Title, message.Content, message.Active,
	)
	return err
}

func (r *repository) UpdateMessage(ctx context.Context, tx output.Tx, message domain.Message) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		_, err := r.db.ExecContext(ctx, queryMessageUpdate,
			message.Code, message.Type, message.Category, message.Module,
			message.Title, message.Content, message.Active, message.ID,
		)
		return err
	}

	_, err := sqlTx.ExecContext(ctx, queryMessageUpdate,
		message.Code, message.Type, message.Category, message.Module,
		message.Title, message.Content, message.Active, message.ID,
	)
	return err
}

func (r *repository) DeleteMessage(ctx context.Context, tx output.Tx, id int64) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		_, err := r.db.ExecContext(ctx, queryMessageDelete, id)
		return err
	}

	_, err := sqlTx.ExecContext(ctx, queryMessageDelete, id)
	return err
}

// ============================================
// Read operations (for service - returns domain.Message)
// ============================================

func (r *repository) GetAllActive(ctx context.Context) ([]domain.Message, error) {
	rows, err := r.db.QueryContext(ctx, queryGetAllActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var m domain.Message
		err := rows.Scan(
			&m.ID, &m.Code, &m.Type, &m.Category, &m.Module,
			&m.Title, &m.Content, &m.Active, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, rows.Err()
}

func (r *repository) GetByCode(ctx context.Context, code string) (*domain.Message, error) {
	var m domain.Message
	err := r.db.QueryRowContext(ctx, queryGetByCode, code).Scan(
		&m.ID, &m.Code, &m.Type, &m.Category, &m.Module,
		&m.Title, &m.Content, &m.Active, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *repository) GetByType(ctx context.Context, msgType string) ([]domain.Message, error) {
	rows, err := r.db.QueryContext(ctx, queryGetByType, msgType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var m domain.Message
		err := rows.Scan(
			&m.ID, &m.Code, &m.Type, &m.Category, &m.Module,
			&m.Title, &m.Content, &m.Active, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, rows.Err()
}

func (r *repository) GetByModule(ctx context.Context, module string) ([]domain.Message, error) {
	rows, err := r.db.QueryContext(ctx, queryGetByModule, module)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var m domain.Message
		err := rows.Scan(
			&m.ID, &m.Code, &m.Type, &m.Category, &m.Module,
			&m.Title, &m.Content, &m.Active, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, rows.Err()
}

// ============================================
// Cache operations (implements messaging.MessageCacheRepository)
// ============================================

// GetAllActiveForCache returns messages for cache (uses cachetypes.CachedMessage type)
func (r *repository) GetAllActiveForCache(ctx context.Context) ([]cachetypes.CachedMessage, error) {
	rows, err := r.db.QueryContext(ctx, queryGetAllActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []cachetypes.CachedMessage
	for rows.Next() {
		var m cachetypes.CachedMessage
		var createdAt, updatedAt interface{} // ignore timestamps for cache
		err := rows.Scan(
			&m.ID, &m.Code, &m.Type, &m.Category, &m.Module,
			&m.Title, &m.Content, &m.Active, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, rows.Err()
}

// GetByCodeForCache returns a message for cache by code
func (r *repository) GetByCodeForCache(ctx context.Context, code string) (*cachetypes.CachedMessage, error) {
	var m cachetypes.CachedMessage
	var createdAt, updatedAt interface{} // ignore timestamps for cache
	err := r.db.QueryRowContext(ctx, queryGetByCode, code).Scan(
		&m.ID, &m.Code, &m.Type, &m.Category, &m.Module,
		&m.Title, &m.Content, &m.Active, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}
