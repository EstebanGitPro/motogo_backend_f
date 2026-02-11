package diagnostic_permission

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

const (
	queryInsert = `
		INSERT INTO motorcycle_diagnostic_permissions (id, motorcycle_id, branch_id, active)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE active = VALUES(active), updated_at = CURRENT_TIMESTAMP
	`

	queryDeactivate = `
		UPDATE motorcycle_diagnostic_permissions 
		SET active = FALSE, updated_at = CURRENT_TIMESTAMP
		WHERE motorcycle_id = ? AND branch_id = ?
	`

	queryGetByMotorcycleAndBranch = `
		SELECT id, motorcycle_id, branch_id, active
		FROM motorcycle_diagnostic_permissions
		WHERE motorcycle_id = ? AND branch_id = ? AND active = TRUE
	`

	queryGetByMotorcycleID = `
		SELECT id, motorcycle_id, branch_id, active
		FROM motorcycle_diagnostic_permissions
		WHERE motorcycle_id = ? AND active = TRUE
		ORDER BY created_at DESC
	`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	db                           *sql.DB
	stmtInsert                   *sql.Stmt
	stmtDeactivate               *sql.Stmt
	stmtGetByMotorcycleAndBranch *sql.Stmt
	stmtGetByMotorcycleID        *sql.Stmt
}

// NewRepository creates a new diagnostic permission repository with prepared statements
func NewRepository(db *sql.DB) (output.DiagnosticPermissionRepository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtInsert, err := db.Prepare(queryInsert)
	if err != nil {
		log.Error(logger.LogDiagPermRepoPrepareSaveError, err)
		return nil, fmt.Errorf("error preparing stmtInsert: %w", err)
	}

	stmtDeactivate, err := db.Prepare(queryDeactivate)
	if err != nil {
		log.Error(logger.LogDiagPermRepoPrepareDeleteError, err)
		return nil, fmt.Errorf("error preparing stmtDeactivate: %w", err)
	}

	stmtGetByMotorcycleAndBranch, err := db.Prepare(queryGetByMotorcycleAndBranch)
	if err != nil {
		log.Error(logger.LogDiagPermRepoPrepareGetError, err)
		return nil, fmt.Errorf("error preparing stmtGetByMotorcycleAndBranch: %w", err)
	}

	stmtGetByMotorcycleID, err := db.Prepare(queryGetByMotorcycleID)
	if err != nil {
		log.Error(logger.LogDiagPermRepoPrepareListError, err)
		return nil, fmt.Errorf("error preparing stmtGetByMotorcycleID: %w", err)
	}

	return &repository{
		db:                           db,
		stmtInsert:                   stmtInsert,
		stmtDeactivate:               stmtDeactivate,
		stmtGetByMotorcycleAndBranch: stmtGetByMotorcycleAndBranch,
		stmtGetByMotorcycleID:        stmtGetByMotorcycleID,
	}, nil
}

// BeginTx starts a new database transaction
func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return common.NewSQLTx(tx), nil
}
