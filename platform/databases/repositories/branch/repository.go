package branch

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

const (
	querySaveBranch = `
		INSERT INTO branches (id, representative_id, franchise_id, name, establishment_type, profile_image_url, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	queryUpdateBranch = `
		UPDATE branches 
		SET franchise_id = ?, name = ?, establishment_type = ?, profile_image_url = ?, status = ?, updated_at = NOW()
		WHERE id = ?
	`
	queryDeleteBranch = "DELETE FROM branches WHERE id = ?"

	queryGetBranchByID = `
		SELECT b.id, b.representative_id, b.franchise_id, b.name, b.establishment_type, b.profile_image_url, b.status,
			   l.id, l.city_id, l.address, l.latitude, l.longitude, c.department_id,
			   p.phone_number
		FROM branches b
		LEFT JOIN locations l ON b.id = l.branch_id
		LEFT JOIN cities c ON l.city_id = c.id
		LEFT JOIN persons p ON b.representative_id = p.id
		WHERE b.id = ?
	`
	queryGetBranchByFranchiseAndName = `
		SELECT id, representative_id, franchise_id, name, establishment_type, profile_image_url, status
		FROM branches 
		WHERE franchise_id = ? AND name = ?
		LIMIT 1
	`
	queryGetBranchesByRepresentative = `
		SELECT b.id, b.representative_id, b.franchise_id, b.name, b.establishment_type, b.profile_image_url, b.status,
		       l.id, l.city_id, l.address, l.latitude, l.longitude, c.department_id,
		       p.phone_number
		FROM branches b
		LEFT JOIN locations l ON b.id = l.branch_id
		LEFT JOIN cities c ON l.city_id = c.id
		LEFT JOIN persons p ON b.representative_id = p.id
		WHERE b.representative_id = ?
	`
	querySaveBranchBrand = `
		INSERT INTO branch_brands (id, branch_id, brand_id, active, created_at, updated_at)
		VALUES (?, ?, ?, TRUE, NOW(), NOW())
	`
	queryDeleteBranchBrands = "DELETE FROM branch_brands WHERE branch_id = ?"
	queryGetBranchBrands    = "SELECT brand_id FROM branch_brands WHERE branch_id = ? AND active = TRUE"

	queryValidateBrands = `
		SELECT id FROM brands WHERE id IN (%s)
	`
	queryHasBranchesByRepresentative = "SELECT EXISTS(SELECT 1 FROM branches WHERE representative_id = ? LIMIT 1)"

	// Displacement range pivot queries
	querySaveBranchDisplacementRange = `
		INSERT INTO branch_displacement_ranges (id, branch_id, displacement_range, active, created_at, updated_at)
		VALUES (?, ?, ?, TRUE, NOW(), NOW())
	`
	queryDeleteBranchDisplacementRanges = "DELETE FROM branch_displacement_ranges WHERE branch_id = ?"
	queryGetBranchDisplacementRanges    = "SELECT displacement_range FROM branch_displacement_ranges WHERE branch_id = ? AND active = TRUE"

	// Service names for nearby search
	queryGetBranchServiceNames = `
		SELECT s.name
		FROM branch_services bs
		JOIN services s ON s.id = bs.service_id
		WHERE bs.branch_id = ? AND bs.active = TRUE
		ORDER BY s.name
	`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	db                              *sql.DB
	stmtSaveBranch                  *sql.Stmt
	stmtUpdateBranch                *sql.Stmt
	stmtDeleteBranch                *sql.Stmt
	stmtGetBranchByID               *sql.Stmt
	stmtGetBranchByFranchiseAndName *sql.Stmt
	stmtGetBranchesByRepresentative *sql.Stmt
	stmtSaveBranchBrand             *sql.Stmt
	stmtDeleteBranchBrands          *sql.Stmt
	stmtGetBranchBrands             *sql.Stmt

	// Displacement range pivot
	stmtSaveBranchDisplacementRange    *sql.Stmt
	stmtDeleteBranchDisplacementRanges *sql.Stmt
	stmtGetBranchDisplacementRanges    *sql.Stmt

	// Service names for nearby search
	stmtGetBranchServiceNames *sql.Stmt
}

// NewRepository creates a new BranchRepository with prepared statements (fail-fast pattern)
func NewRepository(db *sql.DB) (output.BranchRepository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtSaveBranch, err := db.Prepare(querySaveBranch)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtSaveBranch", err)
		return nil, fmt.Errorf("error preparing stmtSaveBranch: %w", err)
	}

	stmtUpdateBranch, err := db.Prepare(queryUpdateBranch)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtUpdateBranch", err)
		return nil, fmt.Errorf("error preparing stmtUpdateBranch: %w", err)
	}

	stmtDeleteBranch, err := db.Prepare(queryDeleteBranch)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtDeleteBranch", err)
		return nil, fmt.Errorf("error preparing stmtDeleteBranch: %w", err)
	}

	stmtGetBranchByID, err := db.Prepare(queryGetBranchByID)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetBranchByID", err)
		return nil, fmt.Errorf("error preparing stmtGetBranchByID: %w", err)
	}

	stmtGetBranchByFranchiseAndName, err := db.Prepare(queryGetBranchByFranchiseAndName)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetBranchByFranchiseAndName", err)
		return nil, fmt.Errorf("error preparing stmtGetBranchByFranchiseAndName: %w", err)
	}

	stmtGetBranchesByRepresentative, err := db.Prepare(queryGetBranchesByRepresentative)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetBranchesByRepresentative", err)
		return nil, fmt.Errorf("error preparing stmtGetBranchesByRepresentative: %w", err)
	}

	stmtSaveBranchBrand, err := db.Prepare(querySaveBranchBrand)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtSaveBranchBrand", err)
		return nil, fmt.Errorf("error preparing stmtSaveBranchBrand: %w", err)
	}

	stmtDeleteBranchBrands, err := db.Prepare(queryDeleteBranchBrands)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtDeleteBranchBrands", err)
		return nil, fmt.Errorf("error preparing stmtDeleteBranchBrands: %w", err)
	}

	stmtGetBranchBrands, err := db.Prepare(queryGetBranchBrands)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetBranchBrands", err)
		return nil, fmt.Errorf("error preparing stmtGetBranchBrands: %w", err)
	}

	stmtSaveBranchDisplacementRange, err := db.Prepare(querySaveBranchDisplacementRange)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtSaveBranchDisplacementRange", err)
		return nil, fmt.Errorf("error preparing stmtSaveBranchDisplacementRange: %w", err)
	}

	stmtDeleteBranchDisplacementRanges, err := db.Prepare(queryDeleteBranchDisplacementRanges)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtDeleteBranchDisplacementRanges", err)
		return nil, fmt.Errorf("error preparing stmtDeleteBranchDisplacementRanges: %w", err)
	}

	stmtGetBranchDisplacementRanges, err := db.Prepare(queryGetBranchDisplacementRanges)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetBranchDisplacementRanges", err)
		return nil, fmt.Errorf("error preparing stmtGetBranchDisplacementRanges: %w", err)
	}

	stmtGetBranchServiceNames, err := db.Prepare(queryGetBranchServiceNames)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetBranchServiceNames", err)
		return nil, fmt.Errorf("error preparing stmtGetBranchServiceNames: %w", err)
	}

	return &repository{
		db:                              db,
		stmtSaveBranch:                  stmtSaveBranch,
		stmtUpdateBranch:                stmtUpdateBranch,
		stmtDeleteBranch:                stmtDeleteBranch,
		stmtGetBranchByID:               stmtGetBranchByID,
		stmtGetBranchByFranchiseAndName: stmtGetBranchByFranchiseAndName,
		stmtGetBranchesByRepresentative: stmtGetBranchesByRepresentative,
		stmtSaveBranchBrand:             stmtSaveBranchBrand,
		stmtDeleteBranchBrands:          stmtDeleteBranchBrands,
		stmtGetBranchBrands:             stmtGetBranchBrands,

		stmtSaveBranchDisplacementRange:    stmtSaveBranchDisplacementRange,
		stmtDeleteBranchDisplacementRanges: stmtDeleteBranchDisplacementRanges,
		stmtGetBranchDisplacementRanges:    stmtGetBranchDisplacementRanges,

		stmtGetBranchServiceNames: stmtGetBranchServiceNames,
	}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return common.NewSQLTx(tx), nil
}
