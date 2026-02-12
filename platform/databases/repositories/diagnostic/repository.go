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
		INSERT INTO diagnostics (id, motorcycle_id, branch_id, date, problem_description, possible_solution, sent_via_whatsapp)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	queryUpdate = `
		UPDATE diagnostics
		SET problem_description = ?, possible_solution = ?, sent_via_whatsapp = ?
		WHERE id = ?
	`

	queryDelete = `
		DELETE FROM diagnostics WHERE id = ?
	`

	queryGetByID = `
		SELECT id, motorcycle_id, branch_id, date, problem_description, possible_solution, sent_via_whatsapp
		FROM diagnostics
		WHERE id = ?
	`

	queryGetByMotorcycleID = `
		SELECT id, motorcycle_id, branch_id, date, problem_description, possible_solution, sent_via_whatsapp
		FROM diagnostics
		WHERE motorcycle_id = ?
		ORDER BY date DESC
	`

	queryInsertEvidence = `
		INSERT INTO diagnostic_evidence (id, diagnostic_id, image_url, description, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	queryGetEvidenceByDiagnosticID = `
		SELECT id, diagnostic_id, image_url, description, created_at
		FROM diagnostic_evidence
		WHERE diagnostic_id = ?
		ORDER BY created_at ASC
	`

	queryGetByMotorcycleAndBranch = `
		SELECT id, motorcycle_id, branch_id, date, problem_description, possible_solution, sent_via_whatsapp
		FROM diagnostics
		WHERE motorcycle_id = ? AND branch_id = ?
		ORDER BY date DESC
		LIMIT 1
	`

	queryDeleteEvidenceByDiagnosticID = `
		DELETE FROM diagnostic_evidence WHERE diagnostic_id = ?
	`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	db                               *sql.DB
	stmtInsert                       *sql.Stmt
	stmtUpdate                       *sql.Stmt
	stmtDelete                       *sql.Stmt
	stmtGetByID                      *sql.Stmt
	stmtGetByMotorcycleID            *sql.Stmt
	stmtInsertEvidence               *sql.Stmt
	stmtGetEvidenceByDiagnosticID    *sql.Stmt
	stmtGetByMotorcycleAndBranch     *sql.Stmt
	stmtDeleteEvidenceByDiagnosticID *sql.Stmt
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

	stmtInsertEvidence, err := db.Prepare(queryInsertEvidence)
	if err != nil {
		log.Error(logger.LogDiagnosticRepoPrepareEvidInsError, err)
		return nil, fmt.Errorf("error preparing stmtInsertEvidence: %w", err)
	}

	stmtGetEvidenceByDiagnosticID, err := db.Prepare(queryGetEvidenceByDiagnosticID)
	if err != nil {
		log.Error(logger.LogDiagnosticRepoPrepareEvidGetError, err)
		return nil, fmt.Errorf("error preparing stmtGetEvidenceByDiagnosticID: %w", err)
	}

	stmtGetByMotorcycleAndBranch, err := db.Prepare(queryGetByMotorcycleAndBranch)
	if err != nil {
		log.Error(logger.LogDiagnosticRepoPrepareGetMotoBranchError, err)
		return nil, fmt.Errorf("error preparing stmtGetByMotorcycleAndBranch: %w", err)
	}

	stmtDeleteEvidenceByDiagnosticID, err := db.Prepare(queryDeleteEvidenceByDiagnosticID)
	if err != nil {
		log.Error(logger.LogDiagnosticRepoPrepareEvidDelError, err)
		return nil, fmt.Errorf("error preparing stmtDeleteEvidenceByDiagnosticID: %w", err)
	}

	return &repository{
		db:                               db,
		stmtInsert:                       stmtInsert,
		stmtUpdate:                       stmtUpdate,
		stmtDelete:                       stmtDelete,
		stmtGetByID:                      stmtGetByID,
		stmtGetByMotorcycleID:            stmtGetByMotorcycleID,
		stmtInsertEvidence:               stmtInsertEvidence,
		stmtGetEvidenceByDiagnosticID:    stmtGetEvidenceByDiagnosticID,
		stmtGetByMotorcycleAndBranch:     stmtGetByMotorcycleAndBranch,
		stmtDeleteEvidenceByDiagnosticID: stmtDeleteEvidenceByDiagnosticID,
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
