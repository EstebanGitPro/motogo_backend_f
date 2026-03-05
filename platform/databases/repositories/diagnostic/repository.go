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
		SELECT d.id, d.motorcycle_id, d.branch_id, b.name AS branch_name, d.date, d.problem_description, d.possible_solution
		FROM diagnostics d
		LEFT JOIN branches b ON b.id = d.branch_id
		WHERE d.id = ?
	`

	queryGetByMotorcycleID = `
		SELECT d.id, d.motorcycle_id, d.branch_id, b.name AS branch_name, d.date, d.problem_description, d.possible_solution
		FROM diagnostics d
		LEFT JOIN branches b ON b.id = d.branch_id
		WHERE d.motorcycle_id = ?
		ORDER BY d.date DESC
	`

	queryGetByMotorcycleAndBranch = `
		SELECT d.id, d.motorcycle_id, d.branch_id, b.name AS branch_name, d.date, d.problem_description, d.possible_solution
		FROM diagnostics d
		LEFT JOIN branches b ON b.id = d.branch_id
		WHERE d.motorcycle_id = ? AND d.branch_id = ?
		ORDER BY d.date DESC
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
	return common.BeginSQLTx(ctx, r.db)
}
