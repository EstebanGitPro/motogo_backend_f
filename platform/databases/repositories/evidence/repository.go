package evidence

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
		INSERT INTO motorcycle_evidence (id, motorcycle_id, angle, image_url, description, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	queryUpdate = `
		UPDATE motorcycle_evidence
		SET angle = ?, image_url = ?
		WHERE id = ?
	`

	queryDelete = `
		DELETE FROM motorcycle_evidence WHERE id = ?
	`

	queryGetByID = `
		SELECT id, motorcycle_id, angle, image_url, description, created_at
		FROM motorcycle_evidence
		WHERE id = ?
	`

	queryGetByMotorcycleID = `
		SELECT id, motorcycle_id, angle, image_url, description, created_at
		FROM motorcycle_evidence
		WHERE motorcycle_id = ?
		ORDER BY created_at DESC
	`

	queryCountByMotorcycleID = `
		SELECT COUNT(*) FROM motorcycle_evidence WHERE motorcycle_id = ?
	`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	db                      *sql.DB
	stmtInsert              *sql.Stmt
	stmtUpdate              *sql.Stmt
	stmtDelete              *sql.Stmt
	stmtGetByID             *sql.Stmt
	stmtGetByMotorcycleID   *sql.Stmt
	stmtCountByMotorcycleID *sql.Stmt
}

// NewRepository creates a new evidence repository with prepared statements
func NewRepository(db *sql.DB) (output.EvidenceRepository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtInsert, err := db.Prepare(queryInsert)
	if err != nil {
		log.Error(logger.LogEvidenceRepoPrepareInsertError, err)
		return nil, fmt.Errorf("error preparing stmtInsert: %w", err)
	}

	stmtUpdate, err := db.Prepare(queryUpdate)
	if err != nil {
		log.Error(logger.LogEvidenceRepoPrepareUpdateError, err)
		return nil, fmt.Errorf("error preparing stmtUpdate: %w", err)
	}

	stmtDelete, err := db.Prepare(queryDelete)
	if err != nil {
		log.Error(logger.LogEvidenceRepoPrepareDeleteError, err)
		return nil, fmt.Errorf("error preparing stmtDelete: %w", err)
	}

	stmtGetByID, err := db.Prepare(queryGetByID)
	if err != nil {
		log.Error(logger.LogEvidenceRepoPrepareGetIDError, err)
		return nil, fmt.Errorf("error preparing stmtGetByID: %w", err)
	}

	stmtGetByMotorcycleID, err := db.Prepare(queryGetByMotorcycleID)
	if err != nil {
		log.Error(logger.LogEvidenceRepoPrepareGetMotoErr, err)
		return nil, fmt.Errorf("error preparing stmtGetByMotorcycleID: %w", err)
	}

	stmtCountByMotorcycleID, err := db.Prepare(queryCountByMotorcycleID)
	if err != nil {
		log.Error(logger.LogEvidenceRepoPrepareCountError, err)
		return nil, fmt.Errorf("error preparing stmtCountByMotorcycleID: %w", err)
	}

	return &repository{
		db:                      db,
		stmtInsert:              stmtInsert,
		stmtUpdate:              stmtUpdate,
		stmtDelete:              stmtDelete,
		stmtGetByID:             stmtGetByID,
		stmtGetByMotorcycleID:   stmtGetByMotorcycleID,
		stmtCountByMotorcycleID: stmtCountByMotorcycleID,
	}, nil
}

// BeginTx starts a new database transaction
func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	return common.BeginSQLTx(ctx, r.db)
}
