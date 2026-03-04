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
		SELECT m.id, m.license_plate, m.reference_id, m.owner_id, m.year, m.current_mileage, m.owner_notes, m.profile_image_url,
		       r.id as ref_id, r.brand_id, b.name as brand_name, r.model, r.category, r.engine_displacement
		FROM motorcycles m
		LEFT JOIN motorcycle_references r ON m.reference_id = r.id
		LEFT JOIN brands b ON r.brand_id = b.id
		WHERE m.id = ? AND m.deleted_at IS NULL
	`

	queryGetByOwnerID = `
		SELECT m.id, m.license_plate, m.reference_id, m.owner_id, m.year, m.current_mileage, m.owner_notes, m.profile_image_url,
		       r.id as ref_id, r.brand_id, b.name as brand_name, r.model, r.category, r.engine_displacement
		FROM motorcycles m
		LEFT JOIN motorcycle_references r ON m.reference_id = r.id
		LEFT JOIN brands b ON r.brand_id = b.id
		WHERE m.owner_id = ? AND m.deleted_at IS NULL
		ORDER BY m.created_at DESC
	`

	queryGetByLicensePlate = `
		SELECT m.id, m.license_plate, m.reference_id, m.owner_id, m.year, m.current_mileage, m.owner_notes, m.profile_image_url,
		       r.id as ref_id, r.brand_id, b.name as brand_name, r.model, r.category, r.engine_displacement
		FROM motorcycles m
		LEFT JOIN motorcycle_references r ON m.reference_id = r.id
		LEFT JOIN brands b ON r.brand_id = b.id
		WHERE m.license_plate = ? AND m.deleted_at IS NULL
	`

	queryUpdate = `
		UPDATE motorcycles 
		SET reference_id = ?, year = ?, current_mileage = ?, owner_notes = ?, profile_image_url = ?, updated_at = NOW()
		WHERE id = ?
	`

	queryDelete = `
		UPDATE motorcycles 
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	queryGetAllReferences = `
		SELECT r.id, r.brand_id, b.name as brand_name, r.model, r.category, r.engine_displacement
		FROM motorcycle_references r
		INNER JOIN brands b ON r.brand_id = b.id
		ORDER BY b.name, r.model
	`

	queryGetReferencesByBrandID = `
		SELECT r.id, r.brand_id, b.name as brand_name, r.model, r.category, r.engine_displacement
		FROM motorcycle_references r
		INNER JOIN brands b ON r.brand_id = b.id
		WHERE r.brand_id = ?
		ORDER BY r.model
	`

	// queryHasServiceHistory validates completed_services and diagnostics tables exist at startup (HU45)
	queryHasServiceHistory = `
		SELECT EXISTS(
			SELECT 1 FROM completed_services WHERE motorcycle_id = ?
			UNION
			SELECT 1 FROM diagnostics WHERE motorcycle_id = ?
		)
	`

	// HU41 - Motorcycle Category Catalogs
	queryGetDistinctCategories = `
		SELECT r.category, COUNT(*) as line_count
		FROM motorcycle_references r
		GROUP BY r.category
		ORDER BY r.category
	`

	queryGetLinesByCategory = `
		SELECT r.model, b.name as brand_name, r.engine_displacement
		FROM motorcycle_references r
		INNER JOIN brands b ON r.brand_id = b.id
		WHERE r.category = ?
		ORDER BY b.name, r.model
	`

	// HU49 - Engine Displacement Ranges
	queryGetDistinctDisplacements = `
		SELECT DISTINCT r.displacement_range
		FROM motorcycle_references r
		WHERE r.displacement_range IS NOT NULL
		ORDER BY r.displacement_range
	`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	db                           *sql.DB
	stmtGetByID                  *sql.Stmt
	stmtGetByOwnerID             *sql.Stmt
	stmtGetByLicensePlate        *sql.Stmt
	stmtUpdate                   *sql.Stmt
	stmtDelete                   *sql.Stmt
	stmtGetAllReferences         *sql.Stmt
	stmtGetReferencesByBrandID   *sql.Stmt
	stmtGetDistinctCategories    *sql.Stmt // HU41
	stmtGetLinesByCategory       *sql.Stmt // HU41
	stmtGetDistinctDisplacements *sql.Stmt // HU49
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

	stmtGetAllReferences, err := db.Prepare(queryGetAllReferences)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetAllReferences", err)
		return nil, fmt.Errorf("error preparing stmtGetAllReferences: %w", err)
	}

	stmtGetReferencesByBrandID, err := db.Prepare(queryGetReferencesByBrandID)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetReferencesByBrandID", err)
		return nil, fmt.Errorf("error preparing stmtGetReferencesByBrandID: %w", err)
	}

	// Schema validation: validate completed_services and diagnostics tables exist (HU45 hybrid delete)
	// This statement is only prepared to detect missing tables at startup, then immediately closed
	stmtHasServiceHistory, err := db.Prepare(queryHasServiceHistory)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtHasServiceHistory (completed_services/diagnostics tables missing?)", err)
		return nil, fmt.Errorf("error preparing stmtHasServiceHistory: %w", err)
	}
	_ = stmtHasServiceHistory.Close() // Close immediately - only used for schema validation

	stmtGetDistinctCategories, err := db.Prepare(queryGetDistinctCategories)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetDistinctCategories", err)
		return nil, fmt.Errorf("error preparing stmtGetDistinctCategories: %w", err)
	}

	stmtGetLinesByCategory, err := db.Prepare(queryGetLinesByCategory)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetLinesByCategory", err)
		return nil, fmt.Errorf("error preparing stmtGetLinesByCategory: %w", err)
	}

	stmtGetDistinctDisplacements, err := db.Prepare(queryGetDistinctDisplacements)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetDistinctDisplacements", err)
		return nil, fmt.Errorf("error preparing stmtGetDistinctDisplacements: %w", err)
	}

	return &repository{
		db:                           db,
		stmtGetByID:                  stmtGetByID,
		stmtGetByOwnerID:             stmtGetByOwnerID,
		stmtGetByLicensePlate:        stmtGetByLicensePlate,
		stmtUpdate:                   stmtUpdate,
		stmtDelete:                   stmtDelete,
		stmtGetAllReferences:         stmtGetAllReferences,
		stmtGetReferencesByBrandID:   stmtGetReferencesByBrandID,
		stmtGetDistinctCategories:    stmtGetDistinctCategories,
		stmtGetLinesByCategory:       stmtGetLinesByCategory,
		stmtGetDistinctDisplacements: stmtGetDistinctDisplacements,
	}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	return common.BeginSQLTx(ctx, r.db)
}
