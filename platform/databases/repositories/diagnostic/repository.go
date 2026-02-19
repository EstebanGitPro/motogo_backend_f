package diagnostic

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
		INSERT INTO diagnostics (id, motorcycle_id, branch_id, date, problem_description, possible_solution)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	queryUpdate = `
		UPDATE diagnostics
		SET problem_description = ?, possible_solution = ?
		WHERE id = ?
	`

	queryDelete = `
		DELETE FROM diagnostics WHERE id = ?
	`

	queryGetByID = `
		SELECT id, motorcycle_id, branch_id, date, problem_description, possible_solution
		FROM diagnostics
		WHERE id = ?
	`

	queryGetByMotorcycleID = `
		SELECT id, motorcycle_id, branch_id, date, problem_description, possible_solution
		FROM diagnostics
		WHERE motorcycle_id = ?
		ORDER BY date DESC
	`

	queryGetByMotorcycleAndBranch = `
		SELECT id, motorcycle_id, branch_id, date, problem_description, possible_solution
		FROM diagnostics
		WHERE motorcycle_id = ? AND branch_id = ?
		ORDER BY date DESC
		LIMIT 1
	`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	db                           *sql.DB
	stmtInsert                   *sql.Stmt
	stmtUpdate                   *sql.Stmt
	stmtDelete                   *sql.Stmt
	stmtGetByID                  *sql.Stmt
	stmtGetByMotorcycleID        *sql.Stmt
	stmtGetByMotorcycleAndBranch *sql.Stmt
}

// NewRepository creates a new diagnostic repository with prepared statements
func NewRepository(db *sql.DB) (output.DiagnosticRepository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtInsert, err := db.Prepare(queryInsert)
	if err != nil {
		log.Error(logger.LogDiagnosticRepoPrepareInsertError, err)
		return nil, fmt.Errorf("error preparing stmtInsert: %w", err)
	}

	stmtUpdate, err := db.Prepare(queryUpdate)
	if err != nil {
		log.Error(logger.LogDiagnosticRepoPrepareUpdateError, err)
		return nil, fmt.Errorf("error preparing stmtUpdate: %w", err)
	}

	stmtDelete, err := db.Prepare(queryDelete)
	if err != nil {
		log.Error(logger.LogDiagnosticRepoPrepareDeleteError, err)
		return nil, fmt.Errorf("error preparing stmtDelete: %w", err)
	}

	stmtGetByID, err := db.Prepare(queryGetByID)
	if err != nil {
		log.Error(logger.LogDiagnosticRepoPrepareGetIDError, err)
		return nil, fmt.Errorf("error preparing stmtGetByID: %w", err)
	}

	stmtGetByMotorcycleID, err := db.Prepare(queryGetByMotorcycleID)
	if err != nil {
		log.Error(logger.LogDiagnosticRepoPrepareGetMotoError, err)
		return nil, fmt.Errorf("error preparing stmtGetByMotorcycleID: %w", err)
	}

	stmtGetByMotorcycleAndBranch, err := db.Prepare(queryGetByMotorcycleAndBranch)
	if err != nil {
		log.Error(logger.LogDiagnosticRepoPrepareGetMotoBranchError, err)
		return nil, fmt.Errorf("error preparing stmtGetByMotorcycleAndBranch: %w", err)
	}

	return &repository{
		db:                           db,
		stmtInsert:                   stmtInsert,
		stmtUpdate:                   stmtUpdate,
		stmtDelete:                   stmtDelete,
		stmtGetByID:                  stmtGetByID,
		stmtGetByMotorcycleID:        stmtGetByMotorcycleID,
		stmtGetByMotorcycleAndBranch: stmtGetByMotorcycleAndBranch,
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
