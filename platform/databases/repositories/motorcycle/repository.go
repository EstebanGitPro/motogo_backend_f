package motorcycle

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

const (
	queryGetByID = `
		SELECT m.id, m.license_plate, m.reference_id, m.owner_id, m.year, m.current_mileage, m.owner_notes,
		       r.id as ref_id, r.brand_id, b.name as brand_name, r.model, r.category, r.engine_displacement
		FROM motorcycles m
		LEFT JOIN motorcycle_references r ON m.reference_id = r.id
		LEFT JOIN brands b ON r.brand_id = b.id
		WHERE m.id = ?
	`

	queryGetByOwnerID = `
		SELECT m.id, m.license_plate, m.reference_id, m.owner_id, m.year, m.current_mileage, m.owner_notes,
		       r.id as ref_id, r.brand_id, b.name as brand_name, r.model, r.category, r.engine_displacement
		FROM motorcycles m
		LEFT JOIN motorcycle_references r ON m.reference_id = r.id
		LEFT JOIN brands b ON r.brand_id = b.id
		WHERE m.owner_id = ?
		ORDER BY m.created_at DESC
	`

	queryGetByLicensePlate = `
		SELECT m.id, m.license_plate, m.reference_id, m.owner_id, m.year, m.current_mileage, m.owner_notes,
		       r.id as ref_id, r.brand_id, b.name as brand_name, r.model, r.category, r.engine_displacement
		FROM motorcycles m
		LEFT JOIN motorcycle_references r ON m.reference_id = r.id
		LEFT JOIN brands b ON r.brand_id = b.id
		WHERE m.license_plate = ?
	`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	db                    *sql.DB
	stmtGetByID           *sql.Stmt
	stmtGetByOwnerID      *sql.Stmt
	stmtGetByLicensePlate *sql.Stmt
}

func NewRepository(db *sql.DB) (output.MotorcycleRepository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtGetByID, err := db.Prepare(queryGetByID)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetByID", err)
		return nil, fmt.Errorf("error preparing stmtGetByID: %w", err)
	}

	stmtGetByOwnerID, err := db.Prepare(queryGetByOwnerID)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetByOwnerID", err)
		return nil, fmt.Errorf("error preparing stmtGetByOwnerID: %w", err)
	}

	stmtGetByLicensePlate, err := db.Prepare(queryGetByLicensePlate)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetByLicensePlate", err)
		return nil, fmt.Errorf("error preparing stmtGetByLicensePlate: %w", err)
	}

	return &repository{
		db:                    db,
		stmtGetByID:           stmtGetByID,
		stmtGetByOwnerID:      stmtGetByOwnerID,
		stmtGetByLicensePlate: stmtGetByLicensePlate,
	}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return common.NewSQLTx(tx), nil
}
