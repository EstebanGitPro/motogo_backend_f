package brand

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

const (
	queryGetAllBrands = "SELECT id, name FROM brands ORDER BY name"

	// Dynamic query for validating brand IDs
	queryValidateBrandIDs = "SELECT id FROM brands WHERE id IN (%s)"
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	db               *sql.DB
	stmtGetAllBrands *sql.Stmt
}

// NewRepository creates a new BrandRepository with prepared statements (fail-fast pattern)
func NewRepository(db *sql.DB) (output.BrandRepository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtGetAllBrands, err := db.Prepare(queryGetAllBrands)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetAllBrands", err)
		return nil, fmt.Errorf("error preparing stmtGetAllBrands: %w", err)
	}

	return &repository{
		db:               db,
		stmtGetAllBrands: stmtGetAllBrands,
	}, nil
}

// GetAllBrands retrieves all brands from the catalog ordered by name
func (r *repository) GetAllBrands(ctx context.Context) ([]domain.Brand, error) {
	rows, err := r.stmtGetAllBrands.QueryContext(ctx)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error", err)
		return nil, err
	}
	defer rows.Close()

	var brands []domain.Brand
	for rows.Next() {
		var brand domain.Brand
		if err := rows.Scan(&brand.ID, &brand.Name); err != nil {
			log.Error(logger.LogDatabaseUnavailable, "error scanning brand", err)
			continue
		}
		brands = append(brands, brand)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return brands, nil
}

// ValidateBrandIDs checks if all provided brand IDs exist in the brands table
func (r *repository) ValidateBrandIDs(ctx context.Context, brandIDs []string) error {
	if len(brandIDs) == 0 {
		return nil
	}

	// Build dynamic query with placeholders
	placeholders := ""
	args := make([]interface{}, len(brandIDs))
	for i, id := range brandIDs {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += "?"
		args[i] = id
	}

	query := fmt.Sprintf(queryValidateBrandIDs, placeholders)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Error("error validating brand IDs", "error", err)
		return err
	}
	defer rows.Close()

	// Collect found brand IDs
	foundBrands := make(map[string]bool)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			continue
		}
		foundBrands[id] = true
	}

	// Check if all brand IDs were found
	for _, id := range brandIDs {
		if !foundBrands[id] {
			log.Warn("brand ID not found", "brand_id", id)
			return domain.ErrBrandNotFound
		}
	}

	return nil
}
