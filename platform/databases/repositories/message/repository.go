package message

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	cachetypes "github.com/EstebanGitPro/motogo-backend/platform/cache/types"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
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

	queryGetByCodeIncludingInactive = `
		SELECT id, message_code, type, category, module, message_title, message_content, is_active, created_at, updated_at 
		FROM system_messages 
		WHERE message_code = ?
		LIMIT 1`
)

var log logger.Logger = logger.NewSlogLogger()

// repository implements output.MessageRepository
type repository struct {
	stmtGetAllActive               *sql.Stmt
	stmtGetByCode                  *sql.Stmt
	stmtGetByCodeForCache          *sql.Stmt
	stmtGetByID                    *sql.Stmt
	stmtGetByType                  *sql.Stmt
	stmtGetByModule                *sql.Stmt
	stmtMessageSave                *sql.Stmt
	stmtMessageDelete              *sql.Stmt
	stmtGetByCodeIncludingInactive *sql.Stmt
	db                             *sql.DB
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

	// Prepare all statements to fail-fast on application startup if SQL queries are malformed
	stmtGetAllActive, err := db.Prepare(queryGetAllActive)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetAllActive", err)
		return nil, err
	}

	stmtGetByCode, err := db.Prepare(queryGetByCode)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetByCode", err)
		return nil, err
	}

	stmtGetByCodeForCache, err := db.Prepare(queryGetByCodeForCache)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetByCodeForCache", err)
		return nil, err
	}

	stmtGetByID, err := db.Prepare(queryGetByID)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetByID", err)
		return nil, err
	}

	stmtGetByType, err := db.Prepare(queryGetByType)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetByType", err)
		return nil, err
	}

	stmtGetByModule, err := db.Prepare(queryGetByModule)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetByModule", err)
		return nil, err
	}

	stmtMessageSave, err := db.Prepare(queryMessageSave)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtMessageSave", err)
		return nil, err
	}

	stmtMessageDelete, err := db.Prepare(queryMessageDelete)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtMessageDelete", err)
		return nil, err
	}

	stmtGetByCodeIncludingInactive, err := db.Prepare(queryGetByCodeIncludingInactive)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetByCodeIncludingInactive", err)
		return nil, err
	}

	return &repository{
		db:                             db,
		stmtGetAllActive:               stmtGetAllActive,
		stmtGetByCode:                  stmtGetByCode,
		stmtGetByCodeForCache:          stmtGetByCodeForCache,
		stmtGetByID:                    stmtGetByID,
		stmtGetByType:                  stmtGetByType,
		stmtGetByModule:                stmtGetByModule,
		stmtMessageSave:                stmtMessageSave,
		stmtMessageDelete:              stmtMessageDelete,
		stmtGetByCodeIncludingInactive: stmtGetByCodeIncludingInactive,
	}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return common.NewSQLTx(tx), nil
}
