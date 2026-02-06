package evidence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

const (
	queryInsert = `
		INSERT INTO motorcycle_evidence (id, motorcycle_id, angle, image_url, upload_date)
		VALUES (?, ?, ?, ?, ?)
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
		SELECT id, motorcycle_id, angle, image_url, upload_date
		FROM motorcycle_evidence
		WHERE id = ?
	`

	queryGetByMotorcycleID = `
		SELECT id, motorcycle_id, angle, image_url, upload_date
		FROM motorcycle_evidence
		WHERE motorcycle_id = ?
		ORDER BY upload_date DESC
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
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtInsert", err)
		return nil, fmt.Errorf("error preparing stmtInsert: %w", err)
	}

	stmtUpdate, err := db.Prepare(queryUpdate)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtUpdate", err)
		return nil, fmt.Errorf("error preparing stmtUpdate: %w", err)
	}

	stmtDelete, err := db.Prepare(queryDelete)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtDelete", err)
		return nil, fmt.Errorf("error preparing stmtDelete: %w", err)
	}

	stmtGetByID, err := db.Prepare(queryGetByID)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetByID", err)
		return nil, fmt.Errorf("error preparing stmtGetByID: %w", err)
	}

	stmtGetByMotorcycleID, err := db.Prepare(queryGetByMotorcycleID)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetByMotorcycleID", err)
		return nil, fmt.Errorf("error preparing stmtGetByMotorcycleID: %w", err)
	}

	stmtCountByMotorcycleID, err := db.Prepare(queryCountByMotorcycleID)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtCountByMotorcycleID", err)
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
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return common.NewSQLTx(tx), nil
}

// Save inserts a new evidence record
func (r *repository) Save(ctx context.Context, tx output.Tx, evidence *domain.MotorcycleEvidence) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	dbEvidence := FromDomain(evidence)
	_, err := sqlTx.ExecContext(ctx, queryInsert,
		dbEvidence.ID,
		dbEvidence.MotorcycleID,
		dbEvidence.Angle,
		dbEvidence.ImageURL,
		dbEvidence.UploadDate,
	)
	if err != nil {
		log.Error(logger.LogDatabaseError, "error saving evidence", err)
		return domain.ErrEvidenceCannotSave
	}

	return nil
}

// Update modifies an existing evidence record
func (r *repository) Update(ctx context.Context, tx output.Tx, evidence *domain.MotorcycleEvidence) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	dbEvidence := FromDomain(evidence)
	result, err := sqlTx.ExecContext(ctx, queryUpdate,
		dbEvidence.Angle,
		dbEvidence.ImageURL,
		dbEvidence.ID,
	)
	if err != nil {
		log.Error(logger.LogDatabaseError, "error updating evidence", err)
		return domain.ErrEvidenceCannotUpdate
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrEvidenceNotFound
	}

	return nil
}

// Delete removes an evidence record
func (r *repository) Delete(ctx context.Context, tx output.Tx, evidenceID string) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	result, err := sqlTx.ExecContext(ctx, queryDelete, evidenceID)
	if err != nil {
		log.Error(logger.LogDatabaseError, "error deleting evidence", err)
		return domain.ErrEvidenceCannotDelete
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrEvidenceNotFound
	}

	return nil
}

// GetByID retrieves an evidence by its ID
func (r *repository) GetByID(ctx context.Context, evidenceID string) (*domain.MotorcycleEvidence, error) {
	var evidence Evidence

	err := r.stmtGetByID.QueryRowContext(ctx, evidenceID).Scan(
		&evidence.ID,
		&evidence.MotorcycleID,
		&evidence.Angle,
		&evidence.ImageURL,
		&evidence.UploadDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrEvidenceNotFound
		}
		log.Error(logger.LogDatabaseError, "error getting evidence by ID", err)
		return nil, err
	}

	result := evidence.ToDomain()
	return &result, nil
}

// GetByMotorcycleID retrieves all evidence for a motorcycle
func (r *repository) GetByMotorcycleID(ctx context.Context, motorcycleID string) ([]domain.MotorcycleEvidence, error) {
	rows, err := r.stmtGetByMotorcycleID.QueryContext(ctx, motorcycleID)
	if err != nil {
		log.Error(logger.LogDatabaseError, "error listing evidence by motorcycle", err)
		return nil, err
	}
	defer rows.Close()

	var evidences []domain.MotorcycleEvidence
	for rows.Next() {
		var evidence Evidence
		if err := rows.Scan(
			&evidence.ID,
			&evidence.MotorcycleID,
			&evidence.Angle,
			&evidence.ImageURL,
			&evidence.UploadDate,
		); err != nil {
			log.Error(logger.LogDatabaseError, "error scanning evidence row", err)
			return nil, err
		}
		evidences = append(evidences, evidence.ToDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return evidences, nil
}

// CountByMotorcycleID counts the number of evidence for a motorcycle
func (r *repository) CountByMotorcycleID(ctx context.Context, motorcycleID string) (int, error) {
	var count int
	err := r.stmtCountByMotorcycleID.QueryRowContext(ctx, motorcycleID).Scan(&count)
	if err != nil {
		log.Error(logger.LogDatabaseError, "error counting evidence", err)
		return 0, err
	}
	return count, nil
}
